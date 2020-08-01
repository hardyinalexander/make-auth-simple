package authentication

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", endpoints.MainPageHandler).Methods("GET")
	r.HandleFunc("/login", endpoints.LoginPageHandler).Methods("GET")

	r.Methods("GET").Path("/callback").Handler(CommonMiddleware(httptransport.NewServer(
		endpoints.MakeLoginEndpoint(),
		decodeLoginRequest,
		encodeResponse,
	)))

	return r

}
