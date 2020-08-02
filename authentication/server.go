package authentication

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, e Endpoints, h Handler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.MainPageHandler).Methods("GET")
	r.HandleFunc("/login", h.LoginPageHandler).Methods("GET")

	r.Methods("GET").Path("/callback").Handler(CommonMiddleware(httptransport.NewServer(
		e.LoginEndpoint(),
		decodeLoginRequest,
		encodeResponse,
	)))

	return r

}
