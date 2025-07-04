package tgdh

import (
	"bytes"
	"testing"

	"github.com/lestrrat-go/jwx/v3/jwe"
)

func TestEncrypt(t *testing.T) {
	root, err := GenerateLeaf()
	if err != nil {
		t.Error(err)
	}
	plainText := []byte("Nacatgunma TGDH encryption test")
	contentType := "text/plain"
	_, cipherText, err := root.Encrypt(plainText, contentType)
	if err != nil {
		t.Error(err)
	}
	protected, plainText1, err := root.Decrypt(cipherText)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(plainText, plainText1) {
		t.Error("plain text does not match")
	}
	var contentType1 string
	err = protected.Get(jwe.ContentTypeKey, &contentType1)
	if err != nil {
		t.Error(err)
	}
	if contentType1 != contentType {
		t.Error("content type does not match")
	}
	var kid1 string
	err = protected.Get(jwe.KeyIDKey, &kid1)
	if err != nil {
		t.Error(err)
	}
	if kid1 != root.Did() {
		t.Error("key ID does not match")
	}
	var type1 string
	err = protected.Get(jwe.TypeKey, &type1)
	if err != nil {
		t.Error(err)
	}
	if type1 != "nacatgunma-tgdh+salt" {
		t.Error("type does not match")
	}
}
