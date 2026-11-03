package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tink-crypto/tink-go/v2/keyset"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
)

const typeURLPrefix = "type.googleapis.com/google.crypto.tink."

// primitive is a Tink primitive and the key types tink-go has for it.
type primitive struct {
	name     string   // with its article, for errors: "an AEAD"
	keyTypes []string // type URLs without typeURLPrefix
}

var (
	aeadPrimitive = primitive{"an AEAD", []string{
		"AesGcmKey", "AesEaxKey", "AesGcmSivKey", "AesCtrHmacAeadKey", "ChaCha20Poly1305Key",
		"XChaCha20Poly1305Key", "XAesGcmKey", "KmsEnvelopeAeadKey",
	}}
	daeadPrimitive         = primitive{"a deterministic AEAD", []string{"AesSivKey"}}
	streamingAEADPrimitive = primitive{"a streaming AEAD", []string{"AesGcmHkdfStreamingKey", "AesCtrHmacStreamingKey"}}
	macPrimitive           = primitive{"a MAC", []string{"HmacKey", "AesCmacKey"}}
	prfPrimitive           = primitive{"a PRF", []string{"HmacPrfKey", "HkdfPrfKey", "AesCmacPrfKey"}}
	signaturePrimitive     = primitive{"a signature", privateAndPublic(
		"Ecdsa", "Ed25519", "RsaSsaPkcs1", "RsaSsaPss", "MlDsa", "SlhDsa", "CompositeMlDsa",
	)}
	hybridPrimitive       = primitive{"a hybrid", privateAndPublic("Hpke", "EciesAeadHkdf")}
	jwtMACPrimitive       = primitive{"a JWT MAC", []string{"JwtHmacKey"}}
	jwtSignaturePrimitive = primitive{"a JWT signature", privateAndPublic(
		"JwtEcdsa", "JwtRsaSsaPkcs1", "JwtRsaSsaPss", "JwtMlDsa",
	)}
)

// privateAndPublic lists the private and public key types of asymmetric
// key families.
func privateAndPublic(families ...string) []string {
	var types []string
	for _, f := range families {
		types = append(types, f+"PrivateKey", f+"PublicKey")
	}
	return types
}

func (p primitive) has(typeURL string) bool {
	name, ok := strings.CutPrefix(typeURL, typeURLPrefix)
	return ok && slices.Contains(p.keyTypes, name)
}

// keysetKind is the kind of keyset an operation takes.
type keysetKind int

const (
	anyKeyset     keysetKind = iota // symmetric primitives
	privateKeyset                   // sign, hybrid decrypt
	publicKeyset                    // verify, hybrid encrypt
)

// checkKeyset fails unless every enabled key of h, the keyset read from
// name, is a key of prim, and h is of the kind command needs. A keyset is
// private when tink-go can make its public keyset.
func checkKeyset(h *keyset.Handle, name, command string, prim primitive, kind keysetKind) error {
	for _, k := range h.KeysetInfo().GetKeyInfo() {
		if k.GetStatus() == tinkpb.KeyStatusType_ENABLED && !prim.has(k.GetTypeUrl()) {
			return fmt.Errorf("%s is not %s keyset (%s)", name, prim.name, k.GetTypeUrl())
		}
	}
	if kind == anyKeyset {
		return nil
	}
	_, err := h.Public()
	private := err == nil
	switch {
	case kind == publicKeyset && private:
		return fmt.Errorf("%s is a private keyset; %s needs the public keyset: make it with circe create-public-keyset", name, command)
	case kind == privateKeyset && !private:
		return fmt.Errorf("%s is a public keyset; %s needs the private keyset", name, command)
	}
	return nil
}
