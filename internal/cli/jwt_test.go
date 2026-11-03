package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tink-crypto/tink-go/v2/jwt"
)

// claimsJSON is a claims object for issuer "circe" and audience "svc" that
// expires at exp.
func claimsJSON(exp time.Time) []byte {
	return fmt.Appendf(nil, `{"iss":"circe","aud":"svc","sub":"user-42","exp":%d,"admin":true}`, exp.Unix())
}

// validatorFor accepts tokens of claimsJSON.
func validatorFor(t *testing.T) *jwt.Validator {
	t.Helper()
	v, err := jwt.NewValidator(&jwt.ValidatorOpts{ExpectedIssuer: ptr("circe"), ExpectedAudience: ptr("svc")})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// expected are the verify flags that accept tokens of claimsJSON.
var expected = []string{"--expected-issuer", "circe", "--expected-audience", "svc"}

// checkClaims fails unless out is the JSON of claims, and one line.
func checkClaims(t *testing.T, out string, claims []byte) {
	t.Helper()
	var got, want map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil || !strings.HasSuffix(out, "}\n") {
		t.Fatalf("verify wrote %q, want the claims JSON and a newline", out)
	}
	if err := json.Unmarshal(claims, &want); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("claims %v, want %v", got, want)
	}
}

func TestJWTMACRoundTripAndInterop(t *testing.T) {
	h := newKeyset(t, "JWT_HS256")
	ks := keysetFile(t, h)
	m, err := jwt.NewMAC(h)
	if err != nil {
		t.Fatal(err)
	}
	claims := claimsJSON(time.Now().Add(time.Hour))
	out := mustRun(t, claims, "jwt-mac", "compute", "--keyset", ks, "--type-header", "JWT")
	if !strings.HasSuffix(out, "\n") || strings.Count(out, ".") != 2 {
		t.Fatalf("compute wrote %q, want a compact JWT and a newline", out)
	}
	if _, err := m.VerifyMACAndDecode(strings.TrimSpace(out), validatorFor(t)); err == nil {
		t.Error("tink-go accepted a typ header the validator did not expect")
	}
	v, err := jwt.NewValidator(&jwt.ValidatorOpts{ExpectedIssuer: ptr("circe"), ExpectedAudience: ptr("svc"), ExpectedTypeHeader: ptr("JWT")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.VerifyMACAndDecode(strings.TrimSpace(out), v); err != nil {
		t.Errorf("tink-go does not verify circe's token: %v", err)
	}
	got := mustRun(t, []byte(out), append([]string{"jwt-mac", "verify", "--keyset", ks, "--expected-type-header", "JWT"}, expected...)...)
	checkClaims(t, got, claims)

	// circe verifies tink-go's token.
	raw, err := jwt.NewRawJWTFromJSON(nil, claims)
	if err != nil {
		t.Fatal(err)
	}
	compact, err := m.ComputeMACAndEncode(raw)
	if err != nil {
		t.Fatal(err)
	}
	checkClaims(t, mustRun(t, []byte(compact), append([]string{"jwt-mac", "verify", "--keyset", ks}, expected...)...), claims)
}

func TestJWTSignatureRoundTripAndInterop(t *testing.T) {
	for _, name := range []string{"JWT_ES256", "JWT_RS256_2048_F4", "JWT_PS256_2048_F4_RAW"} {
		priv := newKeyset(t, name)
		pub := publicOf(t, priv)
		privFile, pubFile := keysetFile(t, priv), keysetFile(t, pub)
		ver, err := jwt.NewVerifier(pub)
		if err != nil {
			t.Fatal(err)
		}
		claims := claimsJSON(time.Now().Add(time.Hour))
		out := mustRun(t, claims, "jwt-signature", "sign", "--keyset", privFile)
		if _, err := ver.VerifyAndDecode(strings.TrimSpace(out), validatorFor(t)); err != nil {
			t.Errorf("%s: tink-go does not verify circe's token: %v", name, err)
		}
		checkClaims(t, mustRun(t, []byte(out), append([]string{"jwt-signature", "verify", "--keyset", pubFile}, expected...)...), claims)

		signer, err := jwt.NewSigner(priv)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := jwt.NewRawJWTFromJSON(nil, claims)
		if err != nil {
			t.Fatal(err)
		}
		compact, err := signer.SignAndEncode(raw)
		if err != nil {
			t.Fatal(err)
		}
		checkClaims(t, mustRun(t, []byte(compact), append([]string{"jwt-signature", "verify", "--keyset", pubFile}, expected...)...), claims)
	}
}

func TestJWTVerifyFailuresExit3(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "JWT_HS256"))
	exp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expired := mustRun(t, claimsJSON(exp), "jwt-mac", "compute", "--keyset", ks)
	fresh := mustRun(t, claimsJSON(time.Now().Add(time.Hour)), "jwt-mac", "compute", "--keyset", ks, "--type-header", "JWT")
	other := mustRun(t, claimsJSON(time.Now().Add(time.Hour)), "jwt-mac", "compute", "--keyset", keysetFile(t, newKeyset(t, "JWT_HS256")))
	verify := func(token string, flags ...string) result {
		return run(t, []byte(token), append([]string{"jwt-mac", "verify", "--keyset", ks}, flags...)...)
	}

	// Before its exp, at a --fixed-now, the expired token verifies.
	if r := verify(expired, append(expected, "--fixed-now", "2019-12-31T00:00:00Z")...); r.code != 0 {
		t.Errorf("expired token before its exp: exit %d, stderr %q", r.code, r.stderr)
	}
	for name, r := range map[string]result{
		"expired":                verify(expired, expected...),
		"expired at --fixed-now": verify(expired, append(expected, "--fixed-now", "2021-01-01T00:00:00Z")...),
		"issuer":                 verify(fresh, "--expected-issuer", "someone", "--expected-audience", "svc", "--expected-type-header", "JWT"),
		"audience":               verify(fresh, "--expected-issuer", "circe", "--expected-audience", "other", "--expected-type-header", "JWT"),
		"type header":            verify(fresh, append(expected, "--expected-type-header", "at+jwt")...),
		"other key":              verify(other, expected...),
		"not a JWT":              verify("not.a.jwt!", expected...),
	} {
		if r.code != 3 || !strings.HasPrefix(r.stderr, "circe: JWT verification failed: ") || r.stdout != "" {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 3 and the reason", name, r.code, r.stdout, r.stderr)
		}
		if strings.Contains(r.stderr, "jwt!") {
			t.Errorf("%s: stderr %q holds the token", name, r.stderr)
		}
	}
}

func TestJWTClaimsErrors(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "JWT_HS256"))
	future := time.Now().Add(time.Hour).Unix()
	for _, c := range []struct {
		claims string
		flags  []string
		want   string // in stderr
	}{
		{`{"iss":"circe"}`, nil, "have no exp; give --without-expiration"},
		{fmt.Sprintf(`{"exp":%d}`, future), []string{"--without-expiration"}, "have an exp, so --without-expiration cannot be given"},
		{`{"exp":"tomorrow"}`, nil, `"exp"`},
		{fmt.Sprintf(`{"exp":%d,"aud":5}`, future), nil, "audience"},
		{fmt.Sprintf(`{"exp":%d,"iss":1}`, future), nil, "iss"},
		{`[1,2]`, nil, "stdin is not a JSON object of claims"},
		{`null`, nil, "stdin is not a JSON object of claims"},
		{`{`, nil, "stdin is not a JSON object of claims"},
	} {
		r := run(t, []byte(c.claims), append([]string{"jwt-mac", "compute", "--keyset", ks}, c.flags...)...)
		if r.code != 1 || !strings.Contains(r.stderr, c.want) || r.stdout != "" {
			t.Errorf("claims %s: exit %d, stderr %q; want 1 and %q", c.claims, r.code, r.stderr, c.want)
		}
	}
}

func TestJWTWithoutExpiration(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "JWT_HS256"))
	claims := []byte(`{"iss":"circe"}`)
	token := mustRun(t, claims, "jwt-mac", "compute", "--keyset", ks, "--without-expiration")
	if r := run(t, []byte(token), "jwt-mac", "verify", "--keyset", ks, "--expected-issuer", "circe"); r.code != 3 {
		t.Errorf("token without exp, no --allow-missing-expiration: exit %d, want 3", r.code)
	}
	checkClaims(t, mustRun(t, []byte(token), "jwt-mac", "verify", "--keyset", ks, "--expected-issuer", "circe", "--allow-missing-expiration"), claims)
}

func TestJWTValidatorUsageErrors(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "JWT_HS256"))
	for _, flags := range [][]string{
		{"--clock-skew", "11m"},
		{"--fixed-now", "yesterday"},
		{"--expected-issuer", "a", "--ignore-issuer"},
		{"--encoding", "hex"},
	} {
		if r := run(t, []byte("x.y.z"), append([]string{"jwt-mac", "verify", "--keyset", ks}, flags...)...); r.code != 2 {
			t.Errorf("verify %v: exit %d, want 2", flags, r.code)
		}
	}
}

func TestJWTSignatureWrongKindOfKeyset(t *testing.T) {
	priv := newKeyset(t, "JWT_ES256")
	privFile, pubFile := keysetFile(t, priv), keysetFile(t, publicOf(t, priv))
	r := run(t, []byte("x.y.z"), "jwt-signature", "verify", "--keyset", privFile)
	want := "circe: " + privFile + " is a private keyset; circe jwt-signature verify needs the public keyset: make it with circe create-public-keyset\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("verify with the private keyset: exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
	claims := claimsJSON(time.Now().Add(time.Hour))
	if r := run(t, claims, "jwt-signature", "sign", "--keyset", pubFile); r.code != 1 || !strings.Contains(r.stderr, "needs the private keyset") {
		t.Errorf("sign with the public keyset: exit %d, stderr %q; want 1", r.code, r.stderr)
	}
	if r := run(t, claims, "jwt-mac", "compute", "--keyset", privFile); r.code != 1 || !strings.Contains(r.stderr, "is not a JWT MAC keyset") {
		t.Errorf("jwt-mac with a JWT signature keyset: exit %d, stderr %q; want 1", r.code, r.stderr)
	}
}
