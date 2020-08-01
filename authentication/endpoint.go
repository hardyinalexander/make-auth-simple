package authentication

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/oauth2"
)

type Endpoints interface {
	MainPageHandler(w http.ResponseWriter, r *http.Request)
	LoginPageHandler(w http.ResponseWriter, r *http.Request)
	MakeLoginEndpoint() endpoint.Endpoint
}

type endpoints struct {
	service           Service
	googleOauthConfig *oauth2.Config
	oauthStateString  string
}

func InitEndpoints(service Service, googleOauthConfig *oauth2.Config, oauthStateString string) Endpoints {
	return &endpoints{service, googleOauthConfig, oauthStateString}
}

func (e *endpoints) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
	<body>
		<a href="/login">Google Log In</a>
	</body>
	</html>`

	fmt.Fprintf(w, htmlIndex)
}

func (e *endpoints) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	url := e.googleOauthConfig.AuthCodeURL(e.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (e *endpoints) MakeLoginEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*LoginRequest)
		token, isRegistered, err := e.service.Login(ctx, req.State, req.Code)
		return LoginResponse{
			Token:        token,
			IsRegistered: isRegistered,
		}, err
	}
}
