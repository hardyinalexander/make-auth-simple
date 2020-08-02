package authentication

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Handler interface {
	MainPageHandler(w http.ResponseWriter, r *http.Request)
	LoginPageHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	googleOauthConfig *oauth2.Config
	oauthStateString  string
}

func InitHandler(googleOauthConfig *oauth2.Config, oauthStateString string) Handler {
	return &handler{googleOauthConfig, oauthStateString}
}

func (h *handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
	<body>
		<a href="/login">Google Log In</a>
	</body>
	</html>`

	fmt.Fprintf(w, htmlIndex)
}

func (h *handler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	url := h.googleOauthConfig.AuthCodeURL(h.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
