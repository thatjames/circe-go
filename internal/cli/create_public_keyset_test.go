package cli

import (
	"testing"
	"time"

	"github.com/tink-crypto/tink-go/v2/hybrid"
	"github.com/tink-crypto/tink-go/v2/jwt"
	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
	"github.com/tink-crypto/tink-go/v2/signature"
)

// publicOf runs create-public-keyset on priv and checks the result holds
// only public keys with priv's IDs and primary.
func publicOf(t *testing.T, priv *keyset.Handle) *keyset.Handle {
	t.Helper()
	pub := parseJSON(t, []byte(mustRun(t, jsonKeyset(t, priv), "create-public-keyset")))
	for _, k := range material(pub).GetKey() {
		if k.GetKeyData().GetKeyMaterialType() != tinkpb.KeyData_ASYMMETRIC_PUBLIC {
			t.Errorf("key %d is %v, want ASYMMETRIC_PUBLIC", k.GetKeyId(), k.GetKeyData().GetKeyMaterialType())
		}
	}
	if pub.KeysetInfo().GetPrimaryKeyId() != priv.KeysetInfo().GetPrimaryKeyId() {
		t.Error("primary changed")
	}
	return pub
}

func TestPublicSignatureKeyset(t *testing.T) {
	for _, name := range []string{"ED25519", "ECDSA_P256"} {
		priv := newKeyset(t, name)
		pub := publicOf(t, priv)
		s, err := signature.NewSigner(priv)
		if err != nil {
			t.Fatal(err)
		}
		sig, err := s.Sign([]byte("data"))
		if err != nil {
			t.Fatal(err)
		}
		v, err := signature.NewVerifier(pub)
		if err != nil {
			t.Fatal(err)
		}
		if err := v.Verify(sig, []byte("data")); err != nil {
			t.Errorf("%s: public keyset does not verify: %v", name, err)
		}
	}
}

func TestPublicHybridKeyset(t *testing.T) {
	priv := newKeyset(t, "DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM")
	pub := publicOf(t, priv)
	enc, err := hybrid.NewHybridEncrypt(pub)
	if err != nil {
		t.Fatal(err)
	}
	ct, err := enc.Encrypt([]byte("data"), nil)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := hybrid.NewHybridDecrypt(priv)
	if err != nil {
		t.Fatal(err)
	}
	if pt, err := dec.Decrypt(ct, nil); err != nil || string(pt) != "data" {
		t.Errorf("decrypt: %q, %v", pt, err)
	}
}

func TestPublicJWTKeyset(t *testing.T) {
	priv := newKeyset(t, "JWT_ES256")
	pub := publicOf(t, priv)
	s, err := jwt.NewSigner(priv)
	if err != nil {
		t.Fatal(err)
	}
	exp := time.Now().Add(time.Hour)
	raw, err := jwt.NewRawJWT(&jwt.RawJWTOptions{Issuer: ptr("circe"), ExpiresAt: &exp})
	if err != nil {
		t.Fatal(err)
	}
	token, err := s.SignAndEncode(raw)
	if err != nil {
		t.Fatal(err)
	}
	v, err := jwt.NewVerifier(pub)
	if err != nil {
		t.Fatal(err)
	}
	val, err := jwt.NewValidator(&jwt.ValidatorOpts{ExpectedIssuer: ptr("circe")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.VerifyAndDecode(token, val); err != nil {
		t.Errorf("public keyset does not verify the token: %v", err)
	}
}

func TestPublicKeysetOfSymmetricKeyset(t *testing.T) {
	r := run(t, jsonKeyset(t, newKeyset(t, "AES256_GCM")), "create-public-keyset")
	if r.code != 1 || r.stdout != "" {
		t.Errorf("exit %d, stdout %q; want 1 and nothing written", r.code, r.stdout)
	}
	pub, err := newKeyset(t, "ED25519").Public()
	if err != nil {
		t.Fatal(err)
	}
	if r := run(t, jsonKeyset(t, pub), "create-public-keyset"); r.code != 1 {
		t.Errorf("public input: exit %d, want 1", r.code)
	}
}

func ptr[T any](v T) *T { return &v }
