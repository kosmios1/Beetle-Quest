package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/sony/gobreaker/v2"
	"golang.org/x/oauth2"
)

type Oauth2Repo struct {
	cb      *gobreaker.CircuitBreaker[*http.Response]
	authCb  *gobreaker.CircuitBreaker[string]
	tokenCb *gobreaker.CircuitBreaker[*oauth2.Token]
}

var (
	oauth2Config = &oauth2.Config{
		ClientID:     utils.FindEnv("OAUTH2_CLIENT_ID"),
		ClientSecret: utils.FindEnv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  utils.FindEnv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  utils.FindEnv("OAUTH2_AUTH_ENDPOINT"),
			TokenURL: utils.FindEnv("OAUTH2_TOKEN_ENDPOINT"),
		},
	}

	adminOauth2Config = &oauth2.Config{
		ClientID:     utils.FindEnv("OAUTH2_CLIENT_ID"),
		ClientSecret: utils.FindEnv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  utils.FindEnv("OAUTH2_ADMIN_REDIRECT_URL"),
		Scopes:       []string{"admin"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  utils.FindEnv("OAUTH2_AUTH_ENDPOINT"),
			TokenURL: utils.FindEnv("OAUTH2_TOKEN_ENDPOINT"),
		},
	}

	revokeTokenEndpoint string = utils.FindEnv("OAUTH2_REVOKE_TOKEN_ENDPOINT")
	verifyTokenEndpoint string = utils.FindEnv("OAUTH2_VERIFY_TOKEN_ENDPOINT")
)

func NewOauth2Repo() *Oauth2Repo {
	return &Oauth2Repo{
		cb:      gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
		authCb:  gobreaker.NewCircuitBreaker[string](gobreaker.Settings{}),
		tokenCb: gobreaker.NewCircuitBreaker[*oauth2.Token](gobreaker.Settings{}),
	}
}

func (r *Oauth2Repo) AuthCodeURL(stateHex, id string) string {
	url, err := r.authCb.Execute(func() (string, error) {
		var url string
		url = oauth2Config.AuthCodeURL(
			stateHex,
			oauth2.SetAuthURLParam("user_id", id),
		)

		if url == "" {
			return "", errors.New("failed to generate auth code url")
		}

		return url, nil
	})

	if err != nil {
		return ""
	}

	return url
}

func (r *Oauth2Repo) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := r.tokenCb.Execute(func() (*oauth2.Token, error) {
		token, err := oauth2Config.Exchange(ctx, code)

		if err != nil {
			return nil, err
		}

		return token, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *Oauth2Repo) RevokeToken(token string) error {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		reqBody := []byte("token=" + token)
		resp, err := http.Post(
			revokeTokenEndpoint,
			"application/x-www-form-urlencoded",
			bytes.NewBuffer(reqBody),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return models.ErrFailedToRevokeToken
	}

	return nil
}

func (r *Oauth2Repo) VerifyToken(token string) error {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		reqBody := []byte("token=" + token)
		resp, err := http.Post(
			verifyTokenEndpoint,
			"application/x-www-form-urlencoded",
			bytes.NewBuffer(reqBody),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return models.ErrFailedToVerifyToken
	}

	return nil
}

// Admin ==============================================================================================================
func (r *Oauth2Repo) AdminAuthCodeURL(stateHex, id string) string {
	url, err := r.authCb.Execute(func() (string, error) {
		var url string
		url = adminOauth2Config.AuthCodeURL(
			stateHex,
			oauth2.SetAuthURLParam("user_id", id),
		)

		if url == "" {
			return "", errors.New("failed to generate auth code url")
		}

		return url, nil
	})

	if err != nil {
		return ""
	}

	return url
}

func (r *Oauth2Repo) AdminExchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := r.tokenCb.Execute(func() (*oauth2.Token, error) {
		token, err := adminOauth2Config.Exchange(ctx, code)

		if err != nil {
			return nil, err
		}

		return token, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
