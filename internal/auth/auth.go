package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "web3auth"

	// nonceTTL bounds the replay window for a captured signature.
	nonceTTL = 5 * time.Minute
)

var addressRe = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// Result is the JSON contract: {"error":0} or {"error":1,"errorMessage":...}.
type Result struct {
	Error        int    `json:"error"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type Service struct {
	store *sessions.CookieStore
}

// New builds a Service backed by a cookie session store. secure marks the
// cookie Secure (HTTPS only) — enable in production behind TLS.
func New(sessionKey []byte, secure bool) *Service {
	store := sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
		Secure:   secure,
	}
	return &Service{store: store}
}

type Challenge struct {
	Message   string `json:"message"`
	CSRFToken string `json:"csrfToken"`
}

// IssueChallenge stores a fresh single-use nonce + CSRF token in the session
// and returns them. The client must sign Message verbatim.
func (s *Service) IssueChallenge(w http.ResponseWriter, r *http.Request, host string) (Challenge, error) {
	nonce := randomHex(16)
	csrf := randomHex(32)
	issued := time.Now()

	// Bind host + timestamp + nonce so a signature is valid only for this site
	// and this single challenge.
	message := fmt.Sprintf("Sign this message to authenticate on %s at %d (nonce: %s)",
		host, issued.UnixMilli(), nonce)

	sess, _ := s.store.Get(r, sessionName)
	sess.Values["pending_message"] = message
	sess.Values["pending_issued"] = issued.Unix()
	sess.Values["csrf_token"] = csrf
	if err := sess.Save(r, w); err != nil {
		return Challenge{}, err
	}

	return Challenge{Message: message, CSRFToken: csrf}, nil
}

func (s *Service) CSRFToken(r *http.Request) string {
	sess, _ := s.store.Get(r, sessionName)
	if v, ok := sess.Values["csrf_token"].(string); ok {
		return v
	}
	return ""
}

func (s *Service) CheckCSRF(r *http.Request, token string) bool {
	want := s.CSRFToken(r)
	if want == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(want), []byte(token)) == 1
}

// Authenticate verifies the signature against the pending challenge and, on
// success, stores the wallet in the session.
func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request, walletAddr, message, signature string) Result {
	if walletAddr == "" || message == "" || signature == "" {
		return Result{Error: 1, ErrorMessage: "Missing required parameters"}
	}

	if !addressRe.MatchString(walletAddr) {
		return Result{Error: 1, ErrorMessage: "Invalid wallet address format"}
	}

	// The signed message must match the challenge issued to this session and
	// not be expired; it is consumed below (single-use).
	sess, _ := s.store.Get(r, sessionName)
	pending, _ := sess.Values["pending_message"].(string)
	issued, _ := sess.Values["pending_issued"].(int64)
	if pending == "" || issued == 0 {
		return Result{Error: 1, ErrorMessage: "No pending challenge; request a nonce first"}
	}
	if subtle.ConstantTimeCompare([]byte(pending), []byte(message)) != 1 {
		return Result{Error: 1, ErrorMessage: "Message does not match issued challenge"}
	}
	if time.Since(time.Unix(issued, 0)) > nonceTTL {
		return Result{Error: 1, ErrorMessage: "Challenge expired; request a new nonce"}
	}

	valid, err := VerifySignature(message, signature, walletAddr)
	if err != nil {
		return Result{Error: 1, ErrorMessage: "Signature verification error: " + err.Error()}
	}
	if !valid {
		return Result{Error: 1, ErrorMessage: "Invalid signature"}
	}

	delete(sess.Values, "pending_message")
	delete(sess.Values, "pending_issued")
	sess.Values["wallet"] = walletAddr
	sess.Values["login_time"] = time.Now().Unix()
	sess.Values["auth_method"] = "web3_signature"
	if err := sess.Save(r, w); err != nil {
		return Result{Error: 1, ErrorMessage: "Failed to persist session: " + err.Error()}
	}

	return Result{Error: 0}
}

func (s *Service) IsAuthenticated(r *http.Request) bool {
	return s.GetWallet(r) != ""
}

func (s *Service) GetWallet(r *http.Request) string {
	sess, _ := s.store.Get(r, sessionName)
	if v, ok := sess.Values["wallet"].(string); ok {
		return v
	}
	return ""
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := s.store.Get(r, sessionName)
	sess.Options.MaxAge = -1
	sess.Values = map[interface{}]interface{}{}
	_ = sess.Save(r, w)
}

// VerifySignature recovers the signer from an EIP-191 personal_sign signature
// and compares it (case-insensitively) to the claimed address. accounts.TextHash
// applies the "\x19Ethereum Signed Message:\n<len>" prefix + keccak256.
func VerifySignature(message, sigHex, address string) (bool, error) {
	if !addressRe.MatchString(address) {
		return false, nil
	}

	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		return false, err
	}
	if len(sig) != 65 { // r(32) + s(32) + v(1)
		return false, nil
	}

	// Ethereum uses recovery id 27/28; secp256k1 recover wants 0/1.
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	if sig[64] != 0 && sig[64] != 1 {
		return false, nil
	}

	pubKey, err := crypto.SigToPub(accounts.TextHash([]byte(message)), sig)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(crypto.PubkeyToAddress(*pubKey).Hex(), address), nil
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}
