package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
)

// signChallenge signs msg with the fixture private key and returns the 0x sig
// with Ethereum's 27/28 recovery id, exactly as a wallet would.
func signChallenge(t *testing.T, msg string) (string, string) {
	t.Helper()
	key, err := crypto.HexToECDSA("4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	if err != nil {
		t.Fatal(err)
	}
	addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
	sig, err := crypto.Sign(accounts.TextHash([]byte(msg)), key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	return addr, "0x" + toHex(sig)
}

func toHex(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2] = hexdigits[c>>4]
		out[i*2+1] = hexdigits[c&0x0f]
	}
	return string(out)
}

// carryCookies copies Set-Cookie from a recorded response into the next request
// so the session round-trips like a browser.
func carryCookies(t *testing.T, rec *httptest.ResponseRecorder, req *http.Request) {
	t.Helper()
	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
}

func TestChallengeAuthFlow(t *testing.T) {
	svc := New([]byte("0123456789abcdef0123456789abcdef"), false)

	// 1. Issue a challenge.
	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/api/nonce", nil)
	ch, err := svc.IssueChallenge(rec1, req1, "localhost")
	if err != nil {
		t.Fatalf("IssueChallenge: %v", err)
	}
	if ch.Message == "" || ch.CSRFToken == "" {
		t.Fatal("expected non-empty challenge and CSRF token")
	}

	// 2. Sign the issued message and authenticate carrying the session cookie.
	addr, sig := signChallenge(t, ch.Message)
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth", nil)
	carryCookies(t, rec1, req2)
	if res := svc.Authenticate(rec2, req2, addr, ch.Message, sig); res.Error != 0 {
		t.Fatalf("expected auth success, got: %+v", res)
	}

	// 3. Session is now authenticated.
	req3 := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	carryCookies(t, rec2, req3)
	if !svc.IsAuthenticated(req3) {
		t.Fatal("expected authenticated session after successful auth")
	}
}

func TestAuth_RejectsWithoutChallenge(t *testing.T) {
	svc := New([]byte("0123456789abcdef0123456789abcdef"), false)
	// Sign an arbitrary message with no prior challenge issued.
	addr, sig := signChallenge(t, "some message")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", nil)
	if res := svc.Authenticate(rec, req, addr, "some message", sig); res.Error == 0 {
		t.Fatal("expected rejection when no challenge was issued")
	}
}

func TestAuth_RejectsReplay(t *testing.T) {
	svc := New([]byte("0123456789abcdef0123456789abcdef"), false)

	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/api/nonce", nil)
	ch, _ := svc.IssueChallenge(rec1, req1, "localhost")
	addr, sig := signChallenge(t, ch.Message)

	// First auth consumes the nonce.
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth", nil)
	carryCookies(t, rec1, req2)
	if res := svc.Authenticate(rec2, req2, addr, ch.Message, sig); res.Error != 0 {
		t.Fatalf("first auth should succeed: %+v", res)
	}

	// Replaying the same signature in a *fresh* session (no pending challenge)
	// must fail.
	rec3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost, "/api/auth", nil)
	if res := svc.Authenticate(rec3, req3, addr, ch.Message, sig); res.Error == 0 {
		t.Fatal("expected replay to be rejected")
	}
}

func TestCSRF_Roundtrip(t *testing.T) {
	svc := New([]byte("0123456789abcdef0123456789abcdef"), false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/nonce", nil)
	ch, _ := svc.IssueChallenge(rec, req, "localhost")

	check := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	carryCookies(t, rec, check)
	if !svc.CheckCSRF(check, ch.CSRFToken) {
		t.Fatal("expected valid CSRF token to pass")
	}
	if svc.CheckCSRF(check, "wrong-token") {
		t.Fatal("expected wrong CSRF token to fail")
	}
}
