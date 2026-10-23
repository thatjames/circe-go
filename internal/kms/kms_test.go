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

// awsEnv points the AWS SDK at a temp HOME whose shared files hold only
// synthetic profiles, and clears the environment credentials.
func awsEnv(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("AWS_CONFIG_FILE", filepath.Join(home, "config"))
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join(home, "credentials"))
	for _, k := range []string{"AWS_PROFILE", "AWS_DEFAULT_PROFILE", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN", "AWS_WEB_IDENTITY_TOKEN_FILE", "AWS_ROLE_ARN"} {
		t.Setenv(k, "")
	}
	write := func(name, content string) {
		if err := os.WriteFile(filepath.Join(home, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("credentials", "[default]\naws_access_key_id = AKIDHOMEDEFAULT\naws_secret_access_key = "+secret+"\n"+
		"[other]\naws_access_key_id = AKIDHOMEOTHER\naws_secret_access_key = "+secret+"\n")
	write("config", "[profile cfgonly]\naws_access_key_id = AKIDHOMECONFIG\naws_secret_access_key = "+secret+"\n")
}

// keyID returns the access key ID the AWS config for uri resolves to; the
// profiles hold static keys, so nothing contacts AWS.
func keyID(t *testing.T, credential string) string {
	t.Helper()
	cfg, err := awsConfig(context.Background(), awsURI, credential)
	if err != nil {
		t.Fatalf("awsConfig: %v", err)
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("region %q, want us-east-1", cfg.Region)
	}
	c, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	return c.AccessKeyID
}

func TestAWSProfileDefaultChain(t *testing.T) {
	awsEnv(t)
	for profile, want := range map[string]string{"": "AKIDHOMEDEFAULT", "other": "AKIDHOMEOTHER", "cfgonly": "AKIDHOMECONFIG"} {
		t.Setenv("AWS_PROFILE", profile)
		if got := keyID(t, ""); got != want {
			t.Errorf("AWS_PROFILE=%q: key %s, want %s", profile, got, want)
		}
	}
	t.Setenv("AWS_PROFILE", "missing")
	_, err := AEAD(context.Background(), awsURI, "")
	if err == nil || !strings.Contains(err.Error(), "missing") || strings.Contains(err.Error(), secret) {
		t.Errorf("missing profile: error %v", err)
	}
}

func TestAWSProfileCredentialINI(t *testing.T) {
	awsEnv(t)
	// Environment keys do not override the profile of the given file.
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIDENV")
	t.Setenv("AWS_SECRET_ACCESS_KEY", secret)
	ini := secretFile(t, "[default]\naws_access_key_id = AKIDFILEDEFAULT\naws_secret_access_key = "+secret+"\n"+
		"[other]\naws_access_key_id = AKIDFILEOTHER\naws_secret_access_key = "+secret+"\n"+
		"[profile cfgstyle]\naws_access_key_id = AKIDFILECONFIG\naws_secret_access_key = "+secret+"\n")
	for profile, want := range map[string]string{"": "AKIDFILEDEFAULT", "other": "AKIDFILEOTHER", "cfgstyle": "AKIDFILECONFIG"} {
		t.Setenv("AWS_PROFILE", profile)
		if got := keyID(t, ini); got != want {
			t.Errorf("AWS_PROFILE=%q: key %s, want %s", profile, got, want)
		}
	}
	if a, err := AEAD(context.Background(), awsURI, ini); err != nil || a == nil {
		t.Errorf("AEAD: %v", err)
	}

	// A profile only in ~/.aws is not taken from the given file.
	for _, profile := range []string{"cfgonly", "missing"} {
		t.Setenv("AWS_PROFILE", profile)
		_, err := AEAD(context.Background(), awsURI, ini)
		if err == nil || !strings.Contains(err.Error(), `"`+profile+`"`) || !strings.Contains(err.Error(), ini) ||
			!strings.Contains(err.Error(), awsURI) || strings.Contains(err.Error(), secret) {
			t.Errorf("AWS_PROFILE=%s: error %v", profile, err)
		}
	}

	// Without AWS_PROFILE, a file with no default profile is an error too.
	t.Setenv("AWS_PROFILE", "")
	_, err := AEAD(context.Background(), awsURI, secretFile(t, "[other]\naws_access_key_id = AKIDX\naws_secret_access_key = "+secret+"\n"))
	if err == nil || !strings.Contains(err.Error(), `"default"`) || strings.Contains(err.Error(), secret) {
		t.Errorf("no default profile: error %v", err)
	}
}

func TestAWSCredentialCSV(t *testing.T) {
	awsEnv(t)
	t.Setenv("AWS_PROFILE", "other")
	csvFile := secretFile(t, "User name,Password,Access key ID,Secret access key,Console login link\n"+
		"circe-test,,AKIDCSV,"+secret+",https://example.com\n")
	if a, err := AEAD(context.Background(), awsURI, csvFile); err != nil || a == nil {
		t.Errorf("CSV: %v", err)
	}
}

func TestAWSRegion(t *testing.T) {
	for uri, want := range map[string]string{
		awsURI: "us-east-1",
		"aws-kms://arn:aws-us-gov:kms:us-gov-west-1:111122223333:alias/k": "us-gov-west-1",
		"aws-kms://arn:aws:kms::111122223333:key/k":                       "",
		"aws-kms://not-an-arn":                                            "",
	} {
		got, err := awsRegion(uri)
		if got != want || (want == "") != (err != nil) {
			t.Errorf("awsRegion(%s) = %q, %v; want %q", uri, got, err, want)
		}
	}
}
