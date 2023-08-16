package util

import (
	"encoding/hex"
	"testing"
)

func TestEncrypt(t *testing.T) {

}

func TestDecrypt(t *testing.T) {
	//s := "437e510c7f6e12893df9a1b42f3871b705a1d4a6905c3c33efc808db68211b6ae5211950446e3c27bdd0168ae93bc2925e2e577e13006d44ee158b49b024453c7e9c1786c3907df638d10065dd0162355aafe715aff6f854f17a2dd5c8edb63854c0ea45d8700d78d31acbe9d07b95a4"
	s := "a6213f77de7e53e16bf642e6c942773477e682b60f8a95544a1159137abb9e1ca15c493c4d1bf65617a426c56435b03f"
	decodeString, _ := hex.DecodeString(s)
	decrypt := AesDecrypt(decodeString, []byte("justdoittobegood"))
	t.Log(string(decrypt))
}
