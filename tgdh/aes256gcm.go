package tgdh

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func (root *Node) DeriveAesKey(salt []byte) ([]byte, error) {
	aesKey := make([]byte, 32)
	err := root.DeriveSeed(aesKey, salt, []byte("nacatgunma-tgdh-aes256gcm"))
	return aesKey, err
}

func (root *Node) Encrypt(plainText []byte, contentType string) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	rand.Read(salt)
	aesKey, err := root.DeriveAesKey(salt)
	if err != nil {
		return nil, nil, err
	}
	did := root.Did()
	protected := jwe.NewHeaders()
	protected.Set(jwe.TypeKey, "nacatgunma-tgdh+salt")
	protected.Set(jwe.KeyIDKey, did)
	protected.Set(jwe.SaltKey, salt)
	if contentType != "" {
		protected.Set(jwe.ContentTypeKey, contentType)
	}
	jweBytes, err := jwe.Encrypt(
		plainText,
		jwe.WithKey(jwa.DIRECT(), aesKey),
		jwe.WithContentEncryption(jwa.A256GCM()),
		jwe.WithProtectedHeaders(protected),
		jwe.WithJSON(),
	)
	if err != nil {
		return nil, nil, err
	}
	jwkKey, err := jwk.Import(aesKey)
	if err != nil {
		return nil, nil, err
	}
	jwkKey.Set(jwk.KeyTypeKey, jwa.OctetSeq())
	jwkKey.Set(jwk.AlgorithmKey, jwa.DIRECT())
	jwkKey.Set(jwk.KeyOpsKey, []jwk.KeyOperation{
		jwk.KeyOpEncrypt,
		jwk.KeyOpDecrypt,
	})
	jwkKey.Set(jwk.KeyIDKey, did)
	jwkKey.Set(jwk.KeyUsageKey, jwk.ForEncryption)
	jwkBytes, err := json.Marshal(jwkKey)
	if err != nil {
		return nil, nil, err
	}
	return jwkBytes, jweBytes, nil
}

func Decrypt(aesKey []byte, cipherText []byte) (jwe.Headers, []byte, error) {
	message, err := jwe.Parse(cipherText)
	if err != nil {
		return nil, nil, err
	}
	protected := message.ProtectedHeaders()
	plainText, err := jwe.Decrypt(
		cipherText,
		jwe.WithKey(jwa.DIRECT(), aesKey),
	)
	return protected, plainText, err
}

func (root *Node) Decrypt(cipherText []byte) (jwe.Headers, []byte, error) {
	message, err := jwe.Parse(cipherText)
	if err != nil {
		return nil, nil, err
	}
	protected := message.ProtectedHeaders()
	var saltString string
	err = protected.Get(jwe.SaltKey, &saltString)
	if err != nil {
		return protected, nil, err
	}
	salt, err := base64.RawURLEncoding.DecodeString(saltString)
	if err != nil {
		return protected, nil, err
	}
	aesKey, err := root.DeriveAesKey(salt)
	if err != nil {
		return protected, nil, err
	}
	plainText, err := jwe.Decrypt(
		cipherText,
		jwe.WithKey(jwa.DIRECT(), aesKey),
	)
	return protected, plainText, err
}
