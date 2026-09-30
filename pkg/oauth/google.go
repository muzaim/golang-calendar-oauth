package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang-test/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleOAuth struct {
	Config *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func NewGoogleOAuth(cfg *config.Config) *GoogleOAuth {
	return &GoogleOAuth{
		Config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				calendar.CalendarScope,
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuth) GetAuthURL(state string) string {
	// AccessTypeOffline asks for refresh token
	return g.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (g *GoogleOAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.Config.Exchange(ctx, code)
}

func (g *GoogleOAuth) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := g.Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get google user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status code: %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode google user info: %w", err)
	}

	return &userInfo, nil
}

func (g *GoogleOAuth) GetCalendarService(ctx context.Context, token *oauth2.Token) (*calendar.Service, error) {
	// TokenSource will automatically refresh access token if expired when using client!
	client := g.Config.Client(ctx, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}
	return srv, nil
}
