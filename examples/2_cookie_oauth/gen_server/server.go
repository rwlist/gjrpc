package gen_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/gjrpc/pkg/transport"
)

const cookieAccessToken = "access_token"

func AccessTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(cookieAccessToken)
	if err != nil {
		return ""
	}
	return cookie.Value
}

type OAuthService interface {
	ExchangeCode(code string) (accessToken string, err error)
}

func NewServer(handler jsonrpc.Handler, oauthService OAuthService) chi.Router {
	rpcHandler := &transport.HandlerHTTP{
		Handler: handler,
	}

	r := chi.NewRouter()
	r.Handle("/rpc", rpcHandler)
	r.Get("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		accessToken, err := oauthService.ExchangeCode(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  cookieAccessToken,
			Value: accessToken,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	})

	return r
}
