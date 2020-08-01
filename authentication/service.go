package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"

	"golang.org/x/oauth2"
)

type Service interface {
	Login(ctx context.Context, state, code string) (string, bool, error)
}

type service struct {
	repo              Repository
	logger            log.Logger
	googleOauthConfig *oauth2.Config
	oauthStateString  string
	secretKey         string
}

func InitService(repo Repository, logger log.Logger, googleOauthConfig *oauth2.Config, oauthStateString string, secretKey string) Service {
	return &service{repo, logger, googleOauthConfig, oauthStateString, secretKey}
}

func (s *service) Login(ctx context.Context, state, code string) (string, bool, error) {
	logger := log.With(s.logger, "function", "Login")

	isRegistered := true

	// Get user info
	user, err := s.getUserInfo(state, code)
	if err != nil {
		level.Error(logger).Log("error getUserInfo:", err)
		return "", false, err
	}

	// Check if user is already registered
	id, err := s.repo.GetUserIDByEmail(ctx, user.Email)
	if err != nil {
		level.Error(logger).Log("error GetUserIDByEmail:", err)
		return "", false, err
	}

	// If user is not registered, create new user
	if id == "" {
		isRegistered = false
		uuid, _ := uuid.NewV4()
		id = uuid.String()
		user.ID = id
		err = s.repo.CreateUser(ctx, user)
		if err != nil {
			level.Error(logger).Log("error CreateUser:", err)
			return "", false, err
		}

	}

	// Create a token based on the user id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		level.Error(logger).Log("error token.SignedString:", err)
		return "", false, err
	}

	return tokenString, isRegistered, nil
}

func (s *service) getUserInfo(state string, code string) (*User, error) {
	if state != s.oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := s.googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	user := new(User)

	err = json.Unmarshal(contents, &user)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling response body: %s", err.Error())
	}

	return user, nil
}
