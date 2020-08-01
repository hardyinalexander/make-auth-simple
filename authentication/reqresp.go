package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	_ "github.com/go-playground/validator/v10"
)

type (
	LoginRequest struct {
		State string `json:"state" validate:"required"`
		Code  string `json:"code" validate:"required"`
	}
	LoginResponse struct {
		Token        string `json:"token"`
		IsRegistered bool   `json:"is_registered"`
	}

	// CompleteProfileRequestResponse struct {
	// 	ID          string    `json:"id" validate:"required"`
	// 	Email       string    `json:"email" validate:"email"`
	// 	Name        string    `json:"name"`
	// 	PhoneNumber string    `json:"phone_number ,omitempty"`
	// 	BirthDate   time.Time `json:"birth_date ,omitempty"`
	// }
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := &LoginRequest{
		State: r.FormValue("state"), // or r.PostForm.Get("state")
		Code:  r.FormValue("code"),
	}
	return req, nil
}
