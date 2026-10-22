// Package templates maps Tink key template names to tink-go key templates.
package templates

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/proto"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/daead"
	"github.com/tink-crypto/tink-go/v2/hybrid"
	"github.com/tink-crypto/tink-go/v2/jwt"
	"github.com/tink-crypto/tink-go/v2/mac"
	"github.com/tink-crypto/tink-go/v2/prf"
	tinkpb "github.com/tink-crypto/tink-go/v2/proto/tink_go_proto"
	"github.com/tink-crypto/tink-go/v2/signature"
	"github.com/tink-crypto/tink-go/v2/streamingaead"
)

// raw returns a template that is t with output prefix RAW, for the _RAW
// names tink-go has no template function for.
func raw(t func() *tinkpb.KeyTemplate) func() *tinkpb.KeyTemplate {
	return func() *tinkpb.KeyTemplate {
		kt := proto.Clone(t()).(*tinkpb.KeyTemplate)
		kt.OutputPrefixType = tinkpb.OutputPrefixType_RAW
		return kt
	}
}

// table holds Tink's template names. A function, not a value, so every
// lookup gets its own template.
var table = map[string]func() *tinkpb.KeyTemplate{
	// AEAD
	"AES128_GCM":                           aead.AES128GCMKeyTemplate,
	"AES128_GCM_RAW":                       aead.AES128GCMNoPrefixKeyTemplate,
	"AES256_GCM":                           aead.AES256GCMKeyTemplate,
	"AES256_GCM_RAW":                       aead.AES256GCMNoPrefixKeyTemplate,
	"AES128_GCM_SIV":                       aead.AES128GCMSIVKeyTemplate,
	"AES128_GCM_SIV_RAW":                   raw(aead.AES128GCMSIVKeyTemplate),
	"AES256_GCM_SIV":                       aead.AES256GCMSIVKeyTemplate,
	"AES256_GCM_SIV_RAW":                   aead.AES256GCMSIVNoPrefixKeyTemplate,
	"AES128_CTR_HMAC_SHA256":               aead.AES128CTRHMACSHA256KeyTemplate,
	"AES128_CTR_HMAC_SHA256_RAW":           raw(aead.AES128CTRHMACSHA256KeyTemplate),
	"AES256_CTR_HMAC_SHA256":               aead.AES256CTRHMACSHA256KeyTemplate,
	"AES256_CTR_HMAC_SHA256_RAW":           raw(aead.AES256CTRHMACSHA256KeyTemplate),
	"CHACHA20_POLY1305":                    aead.ChaCha20Poly1305KeyTemplate,
	"CHACHA20_POLY1305_RAW":                raw(aead.ChaCha20Poly1305KeyTemplate),
	"XCHACHA20_POLY1305":                   aead.XChaCha20Poly1305KeyTemplate,
	"XCHACHA20_POLY1305_RAW":               raw(aead.XChaCha20Poly1305KeyTemplate),
	"XAES_256_GCM_192_BIT_NONCE":           aead.XAES256GCM192BitNonceKeyTemplate,
	"XAES_256_GCM_192_BIT_NONCE_NO_PREFIX": aead.XAES256GCM192BitNonceNoPrefixKeyTemplate,
	"XAES_256_GCM_160_BIT_NONCE":           aead.XAES256GCM160BitNonceKeyTemplate,
	"XAES_256_GCM_160_BIT_NONCE_NO_PREFIX": aead.XAES256GCM160BitNonceNoPrefixKeyTemplate,

	// Deterministic AEAD
	"AES256_SIV":     daead.AESSIVKeyTemplate,
	"AES256_SIV_RAW": raw(daead.AESSIVKeyTemplate),

	// Streaming AEAD
	"AES128_GCM_HKDF_4KB":        streamingaead.AES128GCMHKDF4KBKeyTemplate,
	"AES128_GCM_HKDF_1MB":        streamingaead.AES128GCMHKDF1MBKeyTemplate,
	"AES256_GCM_HKDF_4KB":        streamingaead.AES256GCMHKDF4KBKeyTemplate,
	"AES256_GCM_HKDF_1MB":        streamingaead.AES256GCMHKDF1MBKeyTemplate,
	"AES128_CTR_HMAC_SHA256_4KB": streamingaead.AES128CTRHMACSHA256Segment4KBKeyTemplate,
	"AES128_CTR_HMAC_SHA256_1MB": streamingaead.AES128CTRHMACSHA256Segment1MBKeyTemplate,
	"AES256_CTR_HMAC_SHA256_4KB": streamingaead.AES256CTRHMACSHA256Segment4KBKeyTemplate,
	"AES256_CTR_HMAC_SHA256_1MB": streamingaead.AES256CTRHMACSHA256Segment1MBKeyTemplate,

	// MAC
	"HMAC_SHA256_128BITTAG":     mac.HMACSHA256Tag128KeyTemplate,
	"HMAC_SHA256_128BITTAG_RAW": raw(mac.HMACSHA256Tag128KeyTemplate),
	"HMAC_SHA256_256BITTAG":     mac.HMACSHA256Tag256KeyTemplate,
	"HMAC_SHA256_256BITTAG_RAW": raw(mac.HMACSHA256Tag256KeyTemplate),
	"HMAC_SHA512_256BITTAG":     mac.HMACSHA512Tag256KeyTemplate,
	"HMAC_SHA512_256BITTAG_RAW": raw(mac.HMACSHA512Tag256KeyTemplate),
	"HMAC_SHA512_512BITTAG":     mac.HMACSHA512Tag512KeyTemplate,
	"HMAC_SHA512_512BITTAG_RAW": raw(mac.HMACSHA512Tag512KeyTemplate),
	"AES_CMAC":                  mac.AESCMACTag128KeyTemplate,
	"AES256_CMAC":               mac.AESCMACTag128KeyTemplate,
	"AES256_CMAC_RAW":           raw(mac.AESCMACTag128KeyTemplate),

	// PRF
	"HMAC_SHA256_PRF": prf.HMACSHA256PRFKeyTemplate,
	"HMAC_SHA512_PRF": prf.HMACSHA512PRFKeyTemplate,
	"HKDF_SHA256":     prf.HKDFSHA256PRFKeyTemplate,
	"AES_CMAC_PRF":    prf.AESCMACPRFKeyTemplate,
	"AES256_CMAC_PRF": prf.AESCMACPRFKeyTemplate,

	// Signature
	"ECDSA_P256":                                  signature.ECDSAP256KeyTemplate,
	"ECDSA_P256_NO_PREFIX":                        signature.ECDSAP256KeyWithoutPrefixTemplate,
	"ECDSA_P256_RAW":                              signature.ECDSAP256RawKeyTemplate,
	"ECDSA_P256_IEEE_P1363_WITHOUT_PREFIX":        signature.ECDSAP256RawKeyTemplate,
	"ECDSA_P384":                                  signature.ECDSAP384SHA512KeyTemplate,
	"ECDSA_P384_SHA512":                           signature.ECDSAP384SHA512KeyTemplate,
	"ECDSA_P384_NO_PREFIX":                        signature.ECDSAP384KeyWithoutPrefixTemplate,
	"ECDSA_P384_SHA384":                           signature.ECDSAP384SHA384KeyTemplate,
	"ECDSA_P384_SHA384_NO_PREFIX":                 signature.ECDSAP384SHA384KeyWithoutPrefixTemplate,
	"ECDSA_P521":                                  signature.ECDSAP521KeyTemplate,
	"ECDSA_P521_NO_PREFIX":                        signature.ECDSAP521KeyWithoutPrefixTemplate,
	"ED25519":                                     signature.ED25519KeyTemplate,
	"ED25519_RAW":                                 signature.ED25519KeyWithoutPrefixTemplate,
	"RSA_SSA_PKCS1_3072_SHA256_F4":                signature.RSA_SSA_PKCS1_3072_SHA256_F4_Key_Template,
	"RSA_SSA_PKCS1_3072_SHA256_F4_RAW":            signature.RSA_SSA_PKCS1_3072_SHA256_F4_RAW_Key_Template,
	"RSA_SSA_PKCS1_3072_SHA256_F4_WITHOUT_PREFIX": signature.RSA_SSA_PKCS1_3072_SHA256_F4_RAW_Key_Template,
	"RSA_SSA_PKCS1_4096_SHA512_F4":                signature.RSA_SSA_PKCS1_4096_SHA512_F4_Key_Template,
	"RSA_SSA_PKCS1_4096_SHA512_F4_RAW":            signature.RSA_SSA_PKCS1_4096_SHA512_F4_RAW_Key_Template,
	"RSA_SSA_PSS_3072_SHA256_F4":                  signature.RSA_SSA_PSS_3072_SHA256_32_F4_Key_Template,
	"RSA_SSA_PSS_3072_SHA256_SHA256_32_F4":        signature.RSA_SSA_PSS_3072_SHA256_32_F4_Key_Template,
	"RSA_SSA_PSS_3072_SHA256_F4_RAW":              signature.RSA_SSA_PSS_3072_SHA256_32_F4_Raw_Key_Template,
	"RSA_SSA_PSS_4096_SHA512_F4":                  signature.RSA_SSA_PSS_4096_SHA512_64_F4_Key_Template,
	"RSA_SSA_PSS_4096_SHA512_SHA512_64_F4":        signature.RSA_SSA_PSS_4096_SHA512_64_F4_Key_Template,
	"RSA_SSA_PSS_4096_SHA512_F4_RAW":              signature.RSA_SSA_PSS_4096_SHA512_64_F4_Raw_Key_Template,

	// Hybrid
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_GCM":                     hybrid.ECIESHKDFAES128GCMKeyTemplate,
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_GCM_RAW":                 raw(hybrid.ECIESHKDFAES128GCMKeyTemplate),
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_CTR_HMAC_SHA256":         hybrid.ECIESHKDFAES128CTRHMACSHA256KeyTemplate,
	"ECIES_P256_HKDF_HMAC_SHA256_AES128_CTR_HMAC_SHA256_RAW":     raw(hybrid.ECIESHKDFAES128CTRHMACSHA256KeyTemplate),
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM":           hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_Key_Template,
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_RAW":       hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_Raw_Key_Template,
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM":           hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_Key_Template,
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_RAW":       hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_Raw_Key_Template,
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_CHACHA20_POLY1305":     hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_CHACHA20_POLY1305_Key_Template,
	"DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_CHACHA20_POLY1305_RAW": hybrid.DHKEM_X25519_HKDF_SHA256_HKDF_SHA256_CHACHA20_POLY1305_Raw_Key_Template,
	"DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_128_GCM":             hybrid.DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_Key_Template,
	"DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_RAW":         hybrid.DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_128_GCM_Raw_Key_Template,
	"DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_256_GCM":             hybrid.DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_Key_Template,
	"DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_RAW":         hybrid.DHKEM_P256_HKDF_SHA256_HKDF_SHA256_AES_256_GCM_Raw_Key_Template,

	// JWT
	"JWT_HS256":             jwt.HS256Template,
	"JWT_HS256_RAW":         jwt.RawHS256Template,
	"JWT_HS384":             jwt.HS384Template,
	"JWT_HS384_RAW":         jwt.RawHS384Template,
	"JWT_HS512":             jwt.HS512Template,
	"JWT_HS512_RAW":         jwt.RawHS512Template,
	"JWT_ES256":             jwt.ES256Template,
	"JWT_ES256_RAW":         jwt.RawES256Template,
	"JWT_ES384":             jwt.ES384Template,
	"JWT_ES384_RAW":         jwt.RawES384Template,
	"JWT_ES512":             jwt.ES512Template,
	"JWT_ES512_RAW":         jwt.RawES512Template,
	"JWT_RS256_2048_F4":     jwt.RS256_2048_F4_Key_Template,
	"JWT_RS256_2048_F4_RAW": jwt.RawRS256_2048_F4_Key_Template,
	"JWT_RS256_3072_F4":     jwt.RS256_3072_F4_Key_Template,
	"JWT_RS256_3072_F4_RAW": jwt.RawRS256_3072_F4_Key_Template,
	"JWT_RS384_3072_F4":     jwt.RS384_3072_F4_Key_Template,
	"JWT_RS384_3072_F4_RAW": jwt.RawRS384_3072_F4_Key_Template,
	"JWT_RS512_4096_F4":     jwt.RS512_4096_F4_Key_Template,
	"JWT_RS512_4096_F4_RAW": jwt.RawRS512_4096_F4_Key_Template,
	"JWT_PS256_2048_F4":     jwt.PS256_2048_F4_Key_Template,
	"JWT_PS256_2048_F4_RAW": jwt.RawPS256_2048_F4_Key_Template,
	"JWT_PS256_3072_F4":     jwt.PS256_3072_F4_Key_Template,
	"JWT_PS256_3072_F4_RAW": jwt.RawPS256_3072_F4_Key_Template,
	"JWT_PS384_3072_F4":     jwt.PS384_3072_F4_Key_Template,
	"JWT_PS384_3072_F4_RAW": jwt.RawPS384_3072_F4_Key_Template,
	"JWT_PS512_4096_F4":     jwt.PS512_4096_F4_Key_Template,
	"JWT_PS512_4096_F4_RAW": jwt.RawPS512_4096_F4_Key_Template,
}

// Lookup returns the key template Tink names name.
func Lookup(name string) (*tinkpb.KeyTemplate, error) {
	t, ok := table[name]
	if !ok {
		return nil, fmt.Errorf("unknown key template %q; see circe list-key-templates", name)
	}
	return t(), nil
}

// Names returns every template name, sorted.
func Names() []string {
	names := make([]string, 0, len(table))
	for n := range table {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
