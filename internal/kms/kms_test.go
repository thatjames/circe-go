package kms

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	awsURI   = "aws-kms://arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
	vaultURI = "hcvault://vault.example.com:8200/transit/keys/circe-test"
	gcpURI   = "gcp-kms://projects/p/locations/global/keyRings/r/cryptoKeys/k"
	secret   = "SECRET-VALUE-NEVER-PRINTED"
)

func secretFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "cred")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestUnsupportedPrefix(t *testing.T) {
	_, err := AEAD(context.Background(), "azure-kv://k", "")
	if err == nil {
		t.Fatal("no error")
	}
	for _, p := range []string{gcpPrefix, awsPrefix, vaultPrefix, "azure-kv://k"} {
		if !strings.Contains(err.Error(), p) {
			t.Errorf("error %q does not mention %s", err, p)
		}
	}
}

func TestVaultToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := AEAD(context.Background(), vaultURI, "")
	if err == nil || !strings.Contains(err.Error(), "Vault") || !strings.Contains(err.Error(), vaultURI) ||
		!strings.Contains(err.Error(), "VAULT_TOKEN") {
		t.Errorf("no token: error %v", err)
	}

	// A token from --credential or VAULT_TOKEN builds the AEAD offline.
	if a, err := AEAD(context.Background(), vaultURI, secretFile(t, "  "+secret+"\n")); err != nil || a == nil {
		t.Errorf("token file: %v", err)
	}
	t.Setenv("VAULT_TOKEN", secret)
	if a, err := AEAD(context.Background(), vaultURI, ""); err != nil || a == nil {
		t.Errorf("VAULT_TOKEN: %v", err)
	}

	// A bad key path is an error that keeps the token out.
	_, err = AEAD(context.Background(), "hcvault://vault.example.com/not-a-key-path", "")
	if err == nil || strings.Contains(err.Error(), secret) {
		t.Errorf("bad path: error %v", err)
	}
}

func TestAWS(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	if a, err := AEAD(context.Background(), awsURI, ""); err != nil || a == nil {
		t.Errorf("default chain: %v", err)
	}
	missing := filepath.Join(t.TempDir(), "missing.csv")
	_, err := AEAD(context.Background(), awsURI, missing)
	if err == nil || !strings.Contains(err.Error(), "AWS KMS") || !strings.Contains(err.Error(), awsURI) {
		t.Errorf("missing credential file: error %v", err)
	}
	_, err = AEAD(context.Background(), awsURI, secretFile(t, "not,a\nvalid,"+secret+"\n"))
	if err != nil && strings.Contains(err.Error(), secret) {
		t.Errorf("error %q prints the credential", err)
	}
}

func TestGCP(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.json")
	_, err := AEAD(context.Background(), gcpURI, missing)
	if err == nil || !strings.Contains(err.Error(), "GCP KMS") || !strings.Contains(err.Error(), gcpURI) {
		t.Errorf("missing credential file: error %v", err)
	}
	_, err = AEAD(context.Background(), gcpURI, secretFile(t, `{"type":"service_account","private_key":"`+secret+`"}`))
	if err != nil && strings.Contains(err.Error(), secret) {
		t.Errorf("error %q prints the credential", err)
	}
}
