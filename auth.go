package fluidpay

import (
	"context"
	"errors"
	"net/http"
)

// AuthService obtains JWT tokens and handles account recovery. Most
// integrations should authenticate with an API key instead (see NewClient);
// JWTs are for user-driven sessions. Access it through Client.Auth.
type AuthService struct{ client *Client }

// JWTRequest exchanges a username and password for a JWT.
type JWTRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWT is a token issued by ObtainJWT. Pass Token to WithBearerToken.
type JWT struct {
	APIResource

	Token string `json:"token"`
	// SID is the session ID.
	SID string `json:"sid"`
}

// ForgotUsernameRequest asks for a username reminder email.
type ForgotUsernameRequest struct {
	Email string `json:"email"`
}

// ForgotPasswordRequest asks for a password reset code by email.
type ForgotPasswordRequest struct {
	Username string `json:"username"`
}

// PasswordResetRequest sets a new password with a reset code.
type PasswordResetRequest struct {
	Username  string `json:"username"`
	ResetCode string `json:"reset_code"`
	Password  string `json:"password"`
}

// ObtainJWT exchanges credentials for a JWT. The token envelope carries the
// token at its top level, so this method decodes it directly.
func (s *AuthService) ObtainJWT(ctx context.Context, req *JWTRequest) (*JWT, error) {
	if req == nil {
		return nil, errNilRequest("jwt")
	}
	var out JWT
	resp, err := s.client.call(ctx, http.MethodPost, "token-auth", nil, req, nil)
	if err != nil {
		return nil, err
	}
	// The token is a sibling of status/msg rather than under data, so we
	// re-read the raw body captured on the response.
	if err := decodeTopLevel(resp, &out); err != nil {
		return nil, err
	}
	if out.Token == "" {
		return nil, errors.New("fluidpay: token-auth response did not include a token")
	}
	out.LastResponse = resp
	return &out, nil
}

// Logout invalidates the JWT the client is currently using.
func (s *AuthService) Logout(ctx context.Context) error {
	return doEmpty(ctx, s.client, http.MethodGet, "logout", nil, nil)
}

// ForgotUsername emails a username reminder.
func (s *AuthService) ForgotUsername(ctx context.Context, req *ForgotUsernameRequest) error {
	if req == nil {
		return errNilRequest("forgot username")
	}
	return doEmpty(ctx, s.client, http.MethodPost, "user/forgot-username", nil, req)
}

// ForgotPassword emails a password reset code.
func (s *AuthService) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	if req == nil {
		return errNilRequest("forgot password")
	}
	return doEmpty(ctx, s.client, http.MethodPost, "user/forgot-password", nil, req)
}

// ResetPassword sets a new password using the emailed reset code.
func (s *AuthService) ResetPassword(ctx context.Context, req *PasswordResetRequest) error {
	if req == nil {
		return errNilRequest("password reset")
	}
	return doEmpty(ctx, s.client, http.MethodPost, "user/forgot-password/reset", nil, req)
}
