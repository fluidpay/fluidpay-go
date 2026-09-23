package fluidpay

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// UsersService administers user accounts under the authenticated account.
// Access it through Client.Users.
type UsersService struct{ client *Client }

// User statuses and roles.
const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
	UserRoleAdmin      = "admin"
	UserRoleStandard   = "standard"
)

// UserCreateRequest creates a user.
type UserCreateRequest struct {
	// Username is 7-51 characters starting with a letter.
	Username string `json:"username"`
	Name     string `json:"name"`
	// Phone is digits only.
	Phone string `json:"phone"`
	Email string `json:"email"`
	// Timezone is an IANA name such as "America/Chicago".
	Timezone string `json:"timezone"`
	Password string `json:"password,omitempty"`
	// Status is UserStatusActive or UserStatusDisabled.
	Status string `json:"status"`
	// Role is UserRoleAdmin or UserRoleStandard.
	Role string `json:"role"`
	// SendWelcome emails the new user.
	SendWelcome bool `json:"send_welcome,omitempty"`
	// CreateAPIKey and CreatePubAPIKey create keys and return them on the
	// User as APIKey and PubAPIKey.
	CreateAPIKey    bool `json:"create_api_key,omitempty"`
	CreatePubAPIKey bool `json:"create_pub_api_key,omitempty"`
	// Permissions defaults to false for anything not provided.
	Permissions map[string]bool `json:"permissions,omitempty"`
}

// UserUpdateRequest edits a user.
type UserUpdateRequest struct {
	Name        string          `json:"name,omitempty"`
	Phone       string          `json:"phone,omitempty"`
	Email       string          `json:"email,omitempty"`
	Timezone    string          `json:"timezone,omitempty"`
	Status      string          `json:"status,omitempty"`
	Role        string          `json:"role,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
}

// User is a user account.
type User struct {
	APIResource

	ID            string `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Timezone      string `json:"timezone"`
	Status        string `json:"status"`
	Role          string `json:"role"`
	AccountType   string `json:"account_type"`
	AccountTypeID string `json:"account_type_id"`
	// Permissions maps permission names (for example "process_sale") to
	// whether the user holds them.
	Permissions      map[string]bool `json:"permissions"`
	Flags            map[string]any  `json:"flags"`
	TwoFactorEnabled bool            `json:"two_factor_enabled"`
	// Notifications, Defaults and AccessRestrictions are returned as the
	// gateway sends them.
	Notifications      json.RawMessage `json:"notifications,omitempty"`
	Defaults           json.RawMessage `json:"defaults,omitempty"`
	AccessRestrictions json.RawMessage `json:"access_restrictions,omitempty"`
	// APIKey and PubAPIKey are only populated when the user was just
	// created with CreateAPIKey or CreatePubAPIKey.
	APIKey    string    `json:"api_key,omitempty"`
	PubAPIKey string    `json:"pub_api_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Can reports whether the user holds the named permission.
func (u *User) Can(permission string) bool { return u.Permissions[permission] }

// UserList is a list of users.
type UserList struct {
	APIResource
	Data       []*User
	TotalCount int
}

// ChangePasswordRequest changes the authenticated user's password.
type ChangePasswordRequest struct {
	Username        string `json:"username"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Current retrieves the user that owns the API key or token in use.
func (s *UsersService) Current(ctx context.Context) (*User, error) {
	return doOne[User](ctx, s.client, http.MethodGet, "user", nil, nil)
}

// Get retrieves a user.
func (s *UsersService) Get(ctx context.Context, userID string) (*User, error) {
	if err := requireID("user id", userID); err != nil {
		return nil, err
	}
	return doOne[User](ctx, s.client, http.MethodGet, joinPath("user", userID), nil, nil)
}

// List retrieves all users.
func (s *UsersService) List(ctx context.Context) (*UserList, error) {
	items, resp, err := doList[User](ctx, s.client, http.MethodGet, "users", nil, nil)
	if err != nil {
		return nil, err
	}
	return &UserList{APIResource: APIResource{LastResponse: resp}, Data: items, TotalCount: resp.TotalCount}, nil
}

// Create creates a user.
func (s *UsersService) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	if req == nil {
		return nil, errNilRequest("user")
	}
	return doOne[User](ctx, s.client, http.MethodPost, "user", nil, req)
}

// Update edits a user.
func (s *UsersService) Update(ctx context.Context, userID string, req *UserUpdateRequest) (*User, error) {
	if err := requireID("user id", userID); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errNilRequest("user update")
	}
	return doOne[User](ctx, s.client, http.MethodPost, joinPath("user", userID), nil, req)
}

// Delete removes a user.
func (s *UsersService) Delete(ctx context.Context, userID string) (*APIResponse, error) {
	if err := requireID("user id", userID); err != nil {
		return nil, err
	}
	return doEmpty(ctx, s.client, http.MethodDelete, joinPath("user", userID), nil, nil)
}

// ChangePassword changes the authenticated user's password.
func (s *UsersService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*APIResponse, error) {
	if req == nil {
		return nil, errNilRequest("change password")
	}
	return doEmpty(ctx, s.client, http.MethodPost, "user/change-password", nil, req)
}
