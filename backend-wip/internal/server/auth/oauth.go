package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// oauthHTTPClient is the shared HTTP client for all OAuth operations with a 30s timeout.
var oauthHTTPClient = &http.Client{Timeout: 30 * time.Second}

// OAuthProvider implements OAuth 2.0 / OIDC authentication.
type OAuthProvider struct {
	cfg config.OAuthProvider
}

// ValidateOAuthProviderURLs checks that all custom OAuth URLs use HTTPS.
// Returns an error if any URL uses an insecure scheme.
func ValidateOAuthProviderURLs(providers []config.OAuthProvider) error {
	for _, p := range providers {
		for _, pair := range []struct{ name, val string }{
			{"AuthURL", p.AuthURL},
			{"TokenURL", p.TokenURL},
			{"UserInfoURL", p.UserInfoURL},
			{"Issuer", p.Issuer},
			{"RedirectURL", p.RedirectURL},
		} {
			if pair.val == "" {
				continue
			}
			u, err := url.Parse(pair.val)
			if err != nil {
				return fmt.Errorf("OAuth provider %q: invalid %s: %w", p.Name, pair.name, err)
			}
			if u.Scheme != "https" {
				return fmt.Errorf("OAuth provider %q: %s must use HTTPS (got %q)", p.Name, pair.name, pair.val)
			}
		}
	}
	return nil
}

// NewOAuthProvider creates a new OAuth provider.
func NewOAuthProvider(cfg config.OAuthProvider) *OAuthProvider {
	return &OAuthProvider{cfg: cfg}
}

// Name returns the provider identifier.
func (p *OAuthProvider) Name() string {
	return p.cfg.Name
}

// Authenticate is not used for OAuth (redirects are used instead).
func (p *OAuthProvider) Authenticate(ctx context.Context, credentials map[string]string) (*UserInfo, error) {
	return nil, fmt.Errorf("use GetAuthURL and HandleCallback for OAuth")
}

// GetAuthURL returns the URL to redirect for OAuth login, with PKCE challenge.
func (p *OAuthProvider) GetAuthURL(state, codeChallenge string) string {
	authURL := p.getAuthEndpoint()

	params := url.Values{}
	params.Set("client_id", p.cfg.ClientID)
	params.Set("redirect_uri", p.cfg.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(p.cfg.Scopes, " "))
	params.Set("state", state)

	if codeChallenge != "" {
		params.Set("code_challenge", codeChallenge)
		params.Set("code_challenge_method", "S256")
	}

	return authURL + "?" + params.Encode()
}

// HandleCallback processes the OAuth callback and returns user info.
// codeVerifier is the PKCE verifier to send with the token exchange.
func (p *OAuthProvider) HandleCallback(ctx context.Context, code, codeVerifier string) (*UserInfo, error) {
	// Exchange code for token
	token, err := p.exchangeCode(ctx, code, codeVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := p.getUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

// getAuthEndpoint returns the authorization endpoint URL.
func (p *OAuthProvider) getAuthEndpoint() string {
	if p.cfg.AuthURL != "" {
		return p.cfg.AuthURL
	}

	// Known provider defaults
	switch p.cfg.Name {
	case "google":
		return "https://accounts.google.com/o/oauth2/v2/auth"
	case "github":
		return "https://github.com/login/oauth/authorize"
	case "cognito":
		// Cognito uses issuer-based URLs
		return p.cfg.Issuer + "/oauth2/authorize"
	case "auth0":
		return p.cfg.Issuer + "/authorize"
	default:
		// For generic OIDC, use discovery
		return p.cfg.Issuer + "/authorize"
	}
}

// getTokenEndpoint returns the token endpoint URL.
func (p *OAuthProvider) getTokenEndpoint() string {
	if p.cfg.TokenURL != "" {
		return p.cfg.TokenURL
	}

	switch p.cfg.Name {
	case "google":
		return "https://oauth2.googleapis.com/token"
	case "github":
		return "https://github.com/login/oauth/access_token"
	case "cognito":
		return p.cfg.Issuer + "/oauth2/token"
	case "auth0":
		return p.cfg.Issuer + "/oauth/token"
	default:
		return p.cfg.Issuer + "/token"
	}
}

// getUserInfoEndpoint returns the userinfo endpoint URL.
func (p *OAuthProvider) getUserInfoEndpoint() string {
	if p.cfg.UserInfoURL != "" {
		return p.cfg.UserInfoURL
	}

	switch p.cfg.Name {
	case "google":
		return "https://www.googleapis.com/oauth2/v2/userinfo"
	case "github":
		return "https://api.github.com/user"
	case "cognito":
		return p.cfg.Issuer + "/oauth2/userInfo"
	case "auth0":
		return p.cfg.Issuer + "/userinfo"
	default:
		return p.cfg.Issuer + "/userinfo"
	}
}

// exchangeCode exchanges authorization code for access token, with optional PKCE verifier.
func (p *OAuthProvider) exchangeCode(ctx context.Context, code, codeVerifier string) (string, error) {
	tokenURL := p.getTokenEndpoint()

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.cfg.RedirectURL)
	data.Set("client_id", p.cfg.ClientID)
	data.Set("client_secret", p.cfg.ClientSecret)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		IDToken     string `json:"id_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// getUserInfo fetches user information from the provider.
func (p *OAuthProvider) getUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfoURL := p.getUserInfoEndpoint()

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed: %s", string(body))
	}

	// Parse response based on provider
	return p.parseUserInfo(resp.Body)
}

// parseUserInfo parses the userinfo response based on provider.
func (p *OAuthProvider) parseUserInfo(body io.Reader) (*UserInfo, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}

	user := &UserInfo{
		Provider: p.cfg.Name,
	}

	// Extract fields based on provider
	switch p.cfg.Name {
	case "google":
		user.ProviderID = getString(data, "id")
		user.Email = getString(data, "email")
		user.Name = getString(data, "name")
		user.AvatarURL = getString(data, "picture")
		user.EmailVerified = getBool(data, "verified_email")

	case "github":
		user.ProviderID = fmt.Sprintf("%v", data["id"])
		user.Email = getString(data, "email")
		user.Name = getString(data, "name")
		if user.Name == "" {
			user.Name = getString(data, "login")
		}
		user.AvatarURL = getString(data, "avatar_url")
		user.EmailVerified = true // GitHub emails are verified

	case "cognito", "auth0":
		user.ProviderID = getString(data, "sub")
		user.Email = getString(data, "email")
		user.Name = getString(data, "name")
		user.AvatarURL = getString(data, "picture")
		user.EmailVerified = getBool(data, "email_verified")

	default:
		// Generic OIDC
		user.ProviderID = getString(data, "sub")
		user.Email = getString(data, "email")
		user.Name = getString(data, "name")
		user.AvatarURL = getString(data, "picture")
		user.EmailVerified = getBool(data, "email_verified")
	}

	if user.Email == "" {
		return nil, fmt.Errorf("email not provided by OAuth provider")
	}

	return user, nil
}

// Helper functions for parsing JSON

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if v, ok := data[key].(bool); ok {
		return v
	}
	return false
}
