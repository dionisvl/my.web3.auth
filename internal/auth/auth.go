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

	// nonceTTL bounds how long an issued login challenge stays valid, limiting
	// the replay window for a captured signature.
	nonceTTL = 5 * time.Minute
)

// addressRe matches a checksummed-or-not 0x Ethereum address.
var addressRe = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// Result mirrors the JSON contract returned by the PHP AuthService:
// {"error":0} on success, {"error":1,"errorMessage":"..."} on failure.
type Result struct {
	Error        int    `json:"error"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Service handles Web3 signature auth and session state.
type Service struct {
	store *sessions.CookieStore
}

// New builds a Service backed by a cookie session store keyed by sessionKey.
// secure marks the session cookie Secure (send only over HTTPS) — enable it in
// production behind TLS/Traefik.
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

// Challenge is a server-issued login nonce: the exact message the wallet must
// sign, plus a CSRF token for the subsequent authenticated requests.
type Challenge struct {
	Message   string `json:"message"`
	CSRFToken string `json:"csrfToken"`
}

// IssueChallenge generates a fresh single-use nonce message and CSRF token,
// stores them in the session, and returns them to the client. The client must
// sign Message verbatim; Authenticate then checks it matches and is unexpired.
func (s *Service) IssueChallenge(w http.ResponseWriter, r *http.Request, host string) (Challenge, error) {
	nonce := randomHex(16)
	csrf := randomHex(32)
	issued := time.Now()

	// The message binds host + timestamp + nonce so a signature is meaningful
	// only for this site and this single challenge.
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

// CSRFToken returns the CSRF token stored in the session, or "".
func (s *Service) CSRFToken(r *http.Request) string {
	sess, _ := s.store.Get(r, sessionName)
	if v, ok := sess.Values["csrf_token"].(string); ok {
		return v
	}
	return ""
}

// CheckCSRF constant-time compares the submitted token to the session token.
func (s *Service) CheckCSRF(r *http.Request, token string) bool {
	want := s.CSRFToken(r)
	if want == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(want), []byte(token)) == 1
}

// Authenticate verifies the signature and, on success, stores the wallet in the
// session. Input/behaviour matches the PHP implementation.
func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request, walletAddr, message, signature string) Result {
	if walletAddr == "" || message == "" || signature == "" {
		return Result{Error: 1, ErrorMessage: "Missing required parameters"}
	}

	if !addressRe.MatchString(walletAddr) {
		return Result{Error: 1, ErrorMessage: "Invalid wallet address format"}
	}

	// Replay/CSRF protection: the signed message must be the exact challenge we
	// issued to this session, it must not be expired, and it is single-use.
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

	// Consume the nonce so the same signature cannot be replayed.
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

// IsAuthenticated reports whether the current session has a wallet set.
func (s *Service) IsAuthenticated(r *http.Request) bool {
	return s.GetWallet(r) != ""
}

// GetWallet returns the wallet stored in the session, or "".
func (s *Service) GetWallet(r *http.Request) string {
	sess, _ := s.store.Get(r, sessionName)
	if v, ok := sess.Values["wallet"].(string); ok {
		return v
	}
	return ""
}

// Logout clears the session (equivalent to unsetting the PHP session keys).
func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := s.store.Get(r, sessionName)
	sess.Options.MaxAge = -1
	sess.Values = map[interface{}]interface{}{}
	_ = sess.Save(r, w)
}

// VerifySignature recovers the signer address from an EIP-191 personal_sign
// signature and compares it (case-insensitively) to the claimed address.
//
// This replaces the PHP simplito/elliptic-php + kornrunner/keccak logic. The
// go-ethereum accounts.TextHash applies the exact "\x19Ethereum Signed
// Message:\n<len>" prefix + keccak256 that MetaMask/ethers use.
func VerifySignature(message, sigHex, address string) (bool, error) {
	if !addressRe.MatchString(address) {
		return false, nil
	}

	sig, err := hexutil.Decode(sigHex)
	if err != nil {
		return false, err
	}
	// r(32) + s(32) + v(1)
	if len(sig) != 65 {
		return false, nil
	}

	// Normalize recovery id: Ethereum uses 27/28, secp256k1 recover wants 0/1.
	// Same adjustment the PHP code performed on v.
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	if sig[64] != 0 && sig[64] != 1 {
		return false, nil
	}

	hash := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, err
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	return strings.EqualFold(recovered.Hex(), address), nil
}

// randomHex returns n cryptographically random bytes as a hex string.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failing is fatal-grade; panic surfaces it rather than
		// silently issuing a predictable token.
		panic("crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}
