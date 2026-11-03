package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/jwt"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func newJWTMAC() *cobra.Command {
	return newGroup("jwt-mac", "Compute and verify JWTs with a JWT MAC keyset",
		jwtIssueCommand[jwt.MAC]{
			use: "compute", short: "MAC the claims JSON on --in into a compact JWT",
			prim: jwtMACPrimitive, factory: jwt.NewMAC, op: jwt.MAC.ComputeMACAndEncode,
		}.command(),
		jwtVerifyCommand[jwt.MAC]{
			use: "verify", short: "Verify the compact JWT on --in and write its claims; exit 3 when it does not verify",
			prim: jwtMACPrimitive, factory: jwt.NewMAC, op: jwt.MAC.VerifyMACAndDecode,
		}.command(),
	)
}

func newJWTSignature() *cobra.Command {
	return newGroup("jwt-signature", "Sign JWTs with a private JWT signature keyset and verify them with its public keyset",
		jwtIssueCommand[jwt.Signer]{
			use: "sign", short: "Sign the claims JSON on --in into a compact JWT with the private keyset",
			prim: jwtSignaturePrimitive, kind: privateKeyset, factory: jwt.NewSigner, op: jwt.Signer.SignAndEncode,
		}.command(),
		jwtVerifyCommand[jwt.Verifier]{
			use: "verify", short: "Verify the compact JWT on --in with the public keyset and write its claims; exit 3 when it does not verify",
			prim: jwtSignaturePrimitive, kind: publicKeyset, factory: jwt.NewVerifier, op: jwt.Verifier.VerifyAndDecode,
		}.command(),
	)
}

// jwtIssueCommand reads a claims JSON object and writes the compact JWT
// that op makes of it.
type jwtIssueCommand[P any] struct {
	use, short string
	prim       primitive
	kind       keysetKind
	factory    func(*keyset.Handle) (P, error)
	op         func(p P, raw *jwt.RawJWT) (string, error)
}

func (c jwtIssueCommand[P]) command() *cobra.Command {
	cmd := &cobra.Command{Use: c.use, Short: c.short}
	flags := addPrimitiveFlags(cmd, true)
	var (
		typeHeader        string
		withoutExpiration bool
	)
	cmd.Flags().StringVar(&typeHeader, "type-header", "", "the typ header of the token")
	cmd.Flags().BoolVar(&withoutExpiration, "without-expiration", false, "allow claims without exp: a token that never expires")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		h, err := flags.handle(cmd, c.prim, c.kind)
		if err != nil {
			return err
		}
		p, err := c.factory(h)
		if err != nil {
			return err
		}
		claims, name, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		var th *string
		if cmd.Flags().Changed("type-header") {
			th = &typeHeader
		}
		raw, err := rawJWT(claims, name, th, withoutExpiration)
		if err != nil {
			return err
		}
		compact, err := c.op(p, raw)
		if err != nil {
			return err
		}
		return flags.write(cmd, []byte(compact+"\n"))
	}
	return cmd
}

// rawJWT makes the RawJWT of a claims JSON object read from name. Tink checks
// the types of the registered claims; the exp rules of Tink's RawJWTOptions
// are checked here, since NewRawJWTFromJSON does not apply them.
func rawJWT(claims []byte, name string, typeHeader *string, withoutExpiration bool) (*jwt.RawJWT, error) {
	var obj map[string]any
	if err := json.Unmarshal(claims, &obj); err != nil || obj == nil {
		return nil, fmt.Errorf("%s is not a JSON object of claims", name)
	}
	raw, err := jwt.NewRawJWTFromJSON(typeHeader, claims)
	if err != nil {
		return nil, fmt.Errorf("claims in %s: %w", name, err)
	}
	switch {
	case !raw.HasExpiration() && !withoutExpiration:
		return nil, fmt.Errorf("the claims in %s have no exp; give --without-expiration for a token that never expires", name)
	case raw.HasExpiration() && withoutExpiration:
		return nil, fmt.Errorf("the claims in %s have an exp, so --without-expiration cannot be given", name)
	}
	return raw, nil
}

// jwtVerifyCommand verifies a compact JWT with op against the validator
// flags and writes its claims. A token that does not verify exits 3.
type jwtVerifyCommand[P any] struct {
	use, short string
	prim       primitive
	kind       keysetKind
	factory    func(*keyset.Handle) (P, error)
	op         func(p P, compact string, v *jwt.Validator) (*jwt.VerifiedJWT, error)
}

func (c jwtVerifyCommand[P]) command() *cobra.Command {
	cmd := &cobra.Command{Use: c.use, Short: c.short}
	flags := addPrimitiveFlags(cmd, true)
	vf := addValidatorFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		validator, err := vf.validator(cmd)
		if err != nil {
			return err
		}
		h, err := flags.handle(cmd, c.prim, c.kind)
		if err != nil {
			return err
		}
		p, err := c.factory(h)
		if err != nil {
			return err
		}
		data, _, err := flags.readData(cmd)
		if err != nil {
			return err
		}
		compact := strings.TrimSpace(string(data))
		verified, err := c.op(p, compact, validator)
		if err != nil {
			return failedError{"JWT verification failed: " + jwtReason(err, compact)}
		}
		claims, err := verified.JSONPayload()
		if err != nil {
			return err
		}
		return flags.write(cmd, append(claims, '\n'))
	}
	return cmd
}

// jwtReason is Tink's reason a token did not verify, unless it quotes a
// part of the token: tink-go does when a segment is not base64url.
func jwtReason(err error, compact string) string {
	reason := err.Error()
	for _, seg := range strings.Split(compact, ".") {
		if seg != "" && strings.Contains(reason, seg) {
			return "the token is not a valid JWS compact serialization"
		}
	}
	return reason
}

// validatorFlags are the flags of Tink's jwt.ValidatorOpts.
type validatorFlags struct {
	typeHeader, issuer, audience                    string
	ignoreTypeHeader, ignoreIssuer, ignoreAudiences bool
	allowMissingExpiration, expectIssuedInThePast   bool
	clockSkew                                       time.Duration
	fixedNow                                        string
}

func addValidatorFlags(cmd *cobra.Command) *validatorFlags {
	v := &validatorFlags{}
	f := cmd.Flags()
	f.StringVar(&v.typeHeader, "expected-type-header", "", "the typ header the token must have")
	f.StringVar(&v.issuer, "expected-issuer", "", "the iss claim the token must have")
	f.StringVar(&v.audience, "expected-audience", "", "an audience the token's aud claim must hold")
	f.BoolVar(&v.ignoreTypeHeader, "ignore-type-header", false, "accept any typ header")
	f.BoolVar(&v.ignoreIssuer, "ignore-issuer", false, "accept any iss claim")
	f.BoolVar(&v.ignoreAudiences, "ignore-audiences", false, "accept any aud claim")
	f.BoolVar(&v.allowMissingExpiration, "allow-missing-expiration", false, "accept a token without exp")
	f.BoolVar(&v.expectIssuedInThePast, "expect-issued-in-the-past", false, "refuse a token whose iat is in the future")
	f.DurationVar(&v.clockSkew, "clock-skew", 0, "clock skew to allow for exp, nbf and iat, at most 10m")
	f.StringVar(&v.fixedNow, "fixed-now", "", "validate at this RFC 3339 time instead of now")
	return v
}

// validator makes the jwt.Validator. Flags Tink refuses are usage errors.
func (v *validatorFlags) validator(cmd *cobra.Command) (*jwt.Validator, error) {
	opts := &jwt.ValidatorOpts{
		IgnoreTypeHeader:       v.ignoreTypeHeader,
		IgnoreIssuer:           v.ignoreIssuer,
		IgnoreAudiences:        v.ignoreAudiences,
		AllowMissingExpiration: v.allowMissingExpiration,
		ExpectIssuedInThePast:  v.expectIssuedInThePast,
		ClockSkew:              v.clockSkew,
	}
	if cmd.Flags().Changed("expected-type-header") {
		opts.ExpectedTypeHeader = &v.typeHeader
	}
	if cmd.Flags().Changed("expected-issuer") {
		opts.ExpectedIssuer = &v.issuer
	}
	if cmd.Flags().Changed("expected-audience") {
		opts.ExpectedAudience = &v.audience
	}
	if v.fixedNow != "" {
		t, err := time.Parse(time.RFC3339, v.fixedNow)
		if err != nil {
			return nil, usagef("--fixed-now %q is not an RFC 3339 time", v.fixedNow)
		}
		opts.FixedNow = t
	}
	validator, err := jwt.NewValidator(opts)
	if err != nil {
		return nil, usageError{err}
	}
	return validator, nil
}
