package cli

import (
	"testing"

	"github.com/tink-crypto/tink-go/v2/signature"
)

// signatureTemplates are one template of each signature family.
var signatureTemplates = []string{
	"ECDSA_P256", "ECDSA_P384_SHA384", "ECDSA_P521", "ED25519", "ED25519_RAW",
	"RSA_SSA_PKCS1_3072_SHA256_F4", "RSA_SSA_PSS_3072_SHA256_F4",
}

func TestSignatureEveryFamilyAndEncoding(t *testing.T) {
	data := []byte("release artifact")
	for _, name := range signatureTemplates {
		priv := newKeyset(t, name)
		pub := publicOf(t, priv)
		privFile, pubFile := keysetFile(t, priv), keysetFile(t, pub)
		v, err := signature.NewVerifier(pub)
		if err != nil {
			t.Fatal(err)
		}
		s, err := signature.NewSigner(priv)
		if err != nil {
			t.Fatal(err)
		}
		for _, enc := range encodings {
			out := mustRun(t, data, "signature", "sign", "--keyset", privFile, "--encoding", enc)
			if err := v.Verify(decodeAs(t, enc, out), data); err != nil {
				t.Errorf("%s, %s: tink-go does not verify circe's signature: %v", name, enc, err)
			}
			sig := dataFile(t, "sig", []byte(out))
			if r := run(t, data, "signature", "verify", "--keyset", pubFile, "--signature", sig, "--encoding", enc); r.code != 0 || r.stdout != "" || r.stderr != "" {
				t.Errorf("%s, %s: verify exit %d, stderr %q; want 0 and nothing", name, enc, r.code, r.stderr)
			}
		}
		// circe verifies tink-go's signature, and not of other data.
		tinkSig, err := s.Sign(data)
		if err != nil {
			t.Fatal(err)
		}
		sig := dataFile(t, "sig", tinkSig)
		mustRun(t, data, "signature", "verify", "--keyset", pubFile, "--signature", sig)
		r := run(t, []byte("tampered artifact"), "signature", "verify", "--keyset", pubFile, "--signature", sig)
		if r.code != 3 || r.stderr != "circe: signature verification failed\n" {
			t.Errorf("%s: other data: exit %d, stderr %q; want 3", name, r.code, r.stderr)
		}
	}
}

func TestSignatureWrongKindOfKeyset(t *testing.T) {
	priv := newKeyset(t, "ED25519")
	privFile, pubFile := keysetFile(t, priv), keysetFile(t, publicOf(t, priv))
	sig := dataFile(t, "sig", []byte(mustRun(t, []byte("x"), "signature", "sign", "--keyset", privFile)))

	r := run(t, []byte("x"), "signature", "verify", "--keyset", privFile, "--signature", sig)
	want := "circe: " + privFile + " is a private keyset; circe signature verify needs the public keyset: make it with circe create-public-keyset\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("verify with the private keyset: exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
	r = run(t, []byte("x"), "signature", "sign", "--keyset", pubFile)
	want = "circe: " + pubFile + " is a public keyset; circe signature sign needs the private keyset\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("sign with the public keyset: exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}

func TestSignatureWrongPrimitive(t *testing.T) {
	ks := keysetFile(t, newKeyset(t, "JWT_ES256"))
	r := run(t, []byte("x"), "signature", "sign", "--keyset", ks)
	want := "circe: " + ks + " is not a signature keyset (type.googleapis.com/google.crypto.tink.JwtEcdsaPrivateKey)\n"
	if r.code != 1 || r.stderr != want {
		t.Errorf("exit %d, stderr %q; want 1 and %q", r.code, r.stderr, want)
	}
}
