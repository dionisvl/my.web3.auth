package auth

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/sessions"
)

const sessionName = "web3auth"

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
func New(sessionKey []byte) *Service {
	store := sessions.NewCookieStore(sessionKey)
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
		// Secure is left false so the app works over plain HTTP in local dev
		// (same as the original). Behind TLS/Traefik you can flip this on.
	}
	return &Service{store: store}
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

	valid, err := VerifySignature(message, signature, walletAddr)
	if err != nil {
		return Result{Error: 1, ErrorMessage: "Signature verification error: " + err.Error()}
	}
	if !valid {
		return Result{Error: 1, ErrorMessage: "Invalid signature"}
	}

	sess, _ := s.store.Get(r, sessionName)
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
