// Package kms turns a Tink KMS key URI into the AEAD that encrypts keysets:
// GCP KMS, AWS KMS or HashiCorp Vault.
package kms

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awssdkkms "github.com/aws/aws-sdk-go-v2/service/kms"
	vault "github.com/hashicorp/vault/api"
	"github.com/tink-crypto/tink-go-awskms/v3/integration/awskms"
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
		a, err = awsAEAD(ctx, uri, credential)
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

func awsAEAD(ctx context.Context, uri, credential string) (tink.AEAD, error) {
	var opt awskms.ClientOption
	if credential != "" && isCSV(credential) {
		// Static keys from an IAM console CSV: Tink reads them itself.
		opt = awskms.WithCredentialPath(credential)
	} else {
		cfg, err := awsConfig(ctx, uri, credential)
		if err != nil {
			return nil, err
		}
		opt = awskms.WithKMS(awssdkkms.NewFromConfig(cfg))
	}
	c, err := awskms.NewClientWithOptions(ctx, uri, opt)
	if err != nil {
		return nil, err
	}
	return c.GetAEAD(uri)
}

// awsConfig loads the AWS config for uri: the SDK default chain, which
// honours AWS_PROFILE, or with credential the profile AWS_PROFILE (else
// default) of that shared credentials/config INI file. Credentials are
// resolved lazily; nothing here contacts AWS.
func awsConfig(ctx context.Context, uri, credential string) (aws.Config, error) {
	region, err := awsRegion(uri)
	if err != nil {
		return aws.Config{}, err
	}
	if credential == "" {
		return config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}
	if _, err := os.Stat(credential); err != nil {
		return aws.Config{}, fmt.Errorf("credential file: %w", err)
	}
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
		config.WithSharedConfigFiles([]string{credential}),
		config.WithSharedCredentialsFiles([]string{credential}),
	)
	var notExist config.SharedConfigProfileNotExistError
	if errors.As(err, &notExist) {
		return aws.Config{}, fmt.Errorf("no AWS profile %q in credential file %s", profile, credential)
	}
	if err != nil {
		return aws.Config{}, fmt.Errorf("credential file %s, AWS profile %q: %w", credential, profile, err)
	}
	return cfg, nil
}

// isCSV reports whether path parses as a CSV of more than one column, as
// the IAM console writes; an INI file is one column or no CSV at all.
// Tink's own CSV check decides the same way.
func isCSV(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	rec, err := csv.NewReader(f).Read()
	return err == nil && len(rec) > 1
}

// awsRegion returns the region of an aws-kms://arn:<partition>:kms:<region>:... URI.
func awsRegion(uri string) (string, error) {
	f := strings.Split(strings.TrimPrefix(uri, awsPrefix), ":")
	if len(f) < 5 || f[0] != "arn" || !strings.HasPrefix(f[1], "aws") || f[2] != "kms" || f[3] == "" {
		return "", errors.New("no region in key URI: want aws-kms://arn:<partition>:kms:<region>:...")
	}
	return f[3], nil
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
