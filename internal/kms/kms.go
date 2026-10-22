// Package kms turns a Tink KMS key URI into the AEAD that encrypts keysets:
// GCP KMS, AWS KMS or HashiCorp Vault.
package kms

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/tink-crypto/tink-go-awskms/v2/integration/awskms"
	"github.com/tink-crypto/tink-go-gcpkms/v2/integration/gcpkms"
	"github.com/tink-crypto/tink-go-hcvault/v2/integration/hcvault"
	"github.com/tink-crypto/tink-go/v2/tink"
	"google.golang.org/api/option"
)

const (
	gcpPrefix   = "gcp-kms://"
	awsPrefix   = "aws-kms://"
	vaultPrefix = "hcvault://"
)

// AEAD returns the AEAD of the KMS key uri. credential is the path of the
// backend's credentials file, or "" for its default credentials. Errors
// name the backend and the URI, never credential values.
func AEAD(ctx context.Context, uri, credential string) (tink.AEAD, error) {
	var (
		a       tink.AEAD
		err     error
		backend string
	)
	switch {
	case strings.HasPrefix(uri, gcpPrefix):
		backend = "GCP KMS"
		a, err = gcpAEAD(ctx, uri, credential)
	case strings.HasPrefix(uri, awsPrefix):
		backend = "AWS KMS"
		a, err = awsAEAD(uri, credential)
	case strings.HasPrefix(uri, vaultPrefix):
		backend = "Vault"
		a, err = vaultAEAD(uri, credential)
	default:
		return nil, fmt.Errorf("unsupported master key URI %q: want %s, %s or %s", uri, gcpPrefix, awsPrefix, vaultPrefix)
	}
	if err != nil {
		return nil, fmt.Errorf("%s key %s: %w", backend, uri, err)
	}
	return a, nil
}

func gcpAEAD(ctx context.Context, uri, credential string) (tink.AEAD, error) {
	var opts []gcpkms.Option
	if credential != "" {
		if _, err := os.Stat(credential); err != nil {
			return nil, fmt.Errorf("credential file: %w", err)
		}
		opts = append(opts, gcpkms.WithGoogleAPIClientOptions(option.WithCredentialsFile(credential)))
	}
	c, err := gcpkms.NewClient(ctx, gcpPrefix, opts...)
	if err != nil {
		return nil, err
	}
	return c.GetAEAD(uri)
}

func awsAEAD(uri, credential string) (tink.AEAD, error) {
	var opts []awskms.ClientOption
	if credential != "" {
		if _, err := os.Stat(credential); err != nil {
			return nil, fmt.Errorf("credential file: %w", err)
		}
		opts = append(opts, awskms.WithCredentialPath(credential))
	}
	c, err := awskms.NewClientWithOptions(uri, opts...)
	if err != nil {
		return nil, err
	}
	return c.GetAEAD(uri)
}

func vaultAEAD(uri, credential string) (tink.AEAD, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("malformed URI")
	}
	token := os.Getenv("VAULT_TOKEN")
	if credential != "" {
		b, err := os.ReadFile(credential)
		if err != nil {
			return nil, fmt.Errorf("credential file: %w", err)
		}
		token = strings.TrimSpace(string(b))
	}
	if token == "" {
		return nil, errors.New("no Vault token: pass --credential or set VAULT_TOKEN")
	}
	cfg := vault.DefaultConfig()
	if cfg.Error != nil {
		return nil, cfg.Error
	}
	if u.Host != "" {
		cfg.Address = "https://" + u.Host
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	c, err := hcvault.NewClientWithAEADOptions(vaultPrefix, client.Logical())
	if err != nil {
		return nil, err
	}
	return c.GetAEAD(uri)
}
