package handlers

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/dionisvl/my.web3.auth/internal/auth"
	"github.com/dionisvl/my.web3.auth/internal/wallet"
	"github.com/dionisvl/my.web3.auth/web"
)

// Handlers wires the auth/wallet services and parsed templates.
type Handlers struct {
	auth      *auth.Service
	wallet    *wallet.Service
	templates *template.Template
}

// New parses the embedded templates and returns a Handlers.
func New(authSvc *auth.Service, walletSvc *wallet.Service) (*Handlers, error) {
	tmpl, err := template.ParseFS(web.Templates, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handlers{auth: authSvc, wallet: walletSvc, templates: tmpl}, nil
}

// Register attaches all routes to the given mux. Requires Go 1.22+ for the
// method+path pattern syntax.
func (h *Handlers) Register(mux *http.ServeMux) error {
	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /{$}", h.index)
	mux.HandleFunc("GET /dashboard", h.dashboard)
	mux.HandleFunc("POST /api/auth", h.apiAuth)
	mux.HandleFunc("GET /api/wallet", h.apiWallet)
	mux.HandleFunc("GET /api/logout", h.logout)
	mux.Handle("GET /js/", http.FileServer(http.FS(staticFS)))
	mux.Handle("GET /favicon.ico", http.FileServer(http.FS(staticFS)))
	return nil
}

// GET / — redirect to /dashboard if authed, else render login.
func (h *Handlers) index(w http.ResponseWriter, r *http.Request) {
	if h.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	h.render(w, "login.html", map[string]any{
		"Title":   "Login with Web3 Wallet",
		"Network": h.wallet.GetAPIConfig().Network,
	})
}

// GET /dashboard — require auth, else redirect to /.
func (h *Handlers) dashboard(w http.ResponseWriter, r *http.Request) {
	if !h.auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	addr := h.auth.GetWallet(r)
	h.render(w, "dashboard.html", map[string]any{
		"Title":         "Web3 Wallet Dashboard",
		"Wallet":        addr,
		"WalletDetails": h.wallet.GetWalletDetails(addr),
		"ApiConfig":     h.wallet.GetAPIConfig(),
	})
}

// POST /api/auth — verify signature, set session.
func (h *Handlers) apiAuth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusOK, auth.Result{Error: 1, ErrorMessage: "Invalid form data"})
		return
	}
	result := h.auth.Authenticate(
		w, r,
		r.PostFormValue("wallet"),
		r.PostFormValue("message"),
		r.PostFormValue("signature"),
	)
	writeJSON(w, http.StatusOK, result)
}

// GET /api/wallet — JSON wallet details, 401 if not authed.
func (h *Handlers) apiWallet(w http.ResponseWriter, r *http.Request) {
	if !h.auth.IsAuthenticated(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error":        1,
			"errorMessage": "Not authenticated",
		})
		return
	}
	addr := h.auth.GetWallet(r)
	writeJSON(w, http.StatusOK, map[string]any{
		"error":         0,
		"wallet":        addr,
		"walletDetails": h.wallet.GetWalletDetails(addr),
		"apiConfig":     h.wallet.GetAPIConfig(),
	})
}

// GET /api/logout — clear session, redirect to /.
func (h *Handlers) logout(w http.ResponseWriter, r *http.Request) {
	h.auth.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
