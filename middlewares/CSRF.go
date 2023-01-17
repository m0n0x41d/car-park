package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

var csrfMd func(http.Handler) http.Handler

func init() {
	// TODO: secret token from .env of config file.
	csrfMd = csrf.Protect([]byte("32-byte-long-auth-key"),
		csrf.MaxAge(0),
		csrf.Path("/"),
		csrf.Secure(false),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden. Invalid CSRF Token"}`))
		})),
	)
}

func CSRF() gin.HandlerFunc {
	return adapter.Wrap(csrfMd)
}
