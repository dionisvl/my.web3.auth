package auth

import "testing"

// Fixture generated from the deterministic private key
// 4c0883...362318 signing the message below with EIP-191 personal_sign.
const (
	fxAddress = "0x2c7536E3605D9C16a7a3D7b1898e529396a65c23"
	fxMessage = "Sign this message to authenticate on localhost at 1700000000000"
	fxSig     = "0x94eb90e35b230b90369a3f8b9bd8685d9cab506aeaa1f7d583e6983c01a6e6247dade37731d331dfc961d1a86b52017f2aae1e92b194a48bc6e89229d0feba201c"
)

func TestVerifySignature_Valid(t *testing.T) {
	ok, err := VerifySignature(fxMessage, fxSig, fxAddress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected signature to be valid")
	}
}

func TestVerifySignature_CaseInsensitiveAddress(t *testing.T) {
	lower := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	ok, err := VerifySignature(fxMessage, fxSig, lower)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected valid match regardless of address casing")
	}
}

func TestVerifySignature_WrongAddress(t *testing.T) {
	other := "0x000000000000000000000000000000000000dEaD"
	ok, _ := VerifySignature(fxMessage, fxSig, other)
	if ok {
		t.Fatal("expected mismatch for a different address")
	}
}

func TestVerifySignature_TamperedMessage(t *testing.T) {
	ok, _ := VerifySignature(fxMessage+"!", fxSig, fxAddress)
	if ok {
		t.Fatal("expected failure when the signed message is altered")
	}
}

func TestVerifySignature_BadAddressFormat(t *testing.T) {
	ok, _ := VerifySignature(fxMessage, fxSig, "not-an-address")
	if ok {
		t.Fatal("expected failure for malformed address")
	}
}

func TestVerifySignature_WrongLength(t *testing.T) {
	// 64 bytes instead of 65.
	short := "0x" + "ab" + fxSig[4:]
	short = short[:130]
	ok, _ := VerifySignature(fxMessage, short, fxAddress)
	if ok {
		t.Fatal("expected failure for wrong signature length")
	}
}

func TestVerifySignature_BadHex(t *testing.T) {
	ok, err := VerifySignature(fxMessage, "0xZZZZ", fxAddress)
	if ok {
		t.Fatal("expected failure for non-hex signature")
	}
	if err == nil {
		t.Fatal("expected an error for malformed hex")
	}
}
