package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User is the minimal authenticated user payload returned to the frontend.
type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// PortalRole determines which distributor console the authenticated user can access.
type PortalRole string

const (
	// PortalRoleDistributor grants access to distributor-facing portal routes.
	PortalRoleDistributor PortalRole = "distributor"
	// PortalRoleOperator grants access to operator-facing console routes.
	PortalRoleOperator PortalRole = "operator"
)

// Claims is the JWT payload signed by the distributor system.
type Claims struct {
	UserID     int64      `json:"user_id"`
	Email      string     `json:"email"`
	PortalRole PortalRole `json:"portal_role"`
	jwt.RegisteredClaims
}

// Service authenticates users against the shared main-system database and
// manages the distributor portal JWT lifecycle.
type Service struct {
	db        *sql.DB
	jwtSecret []byte
}

// NewService creates an authentication service backed by the shared users table.
func NewService(db *sql.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// Login validates credentials, resolves the caller's portal role, and returns
// a signed JWT for subsequent API access.
func (s *Service) Login(ctx context.Context, email string, password string) (*User, PortalRole, string, error) {
	var user User
	var passwordHash string
	var status string
	err := s.db.QueryRowContext(ctx, `
SELECT id, email, password_hash, role, status
FROM users
WHERE email = $1 AND deleted_at IS NULL
`, email).Scan(&user.ID, &user.Email, &passwordHash, &user.Role, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", fmt.Errorf("query user: %w", err)
	}
	if status != "active" {
		return nil, "", "", ErrInactiveUser
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	portalRole, err := s.resolvePortalRole(ctx, user.ID, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	token, err := s.issueToken(user, portalRole)
	if err != nil {
		return nil, "", "", err
	}

	return &user, portalRole, token, nil
}

// resolvePortalRole maps a main-system user to either operator access or an
// enabled distributor profile.
func (s *Service) resolvePortalRole(ctx context.Context, userID int64, role string) (PortalRole, error) {
	if role == "admin" {
		return PortalRoleOperator, nil
	}

	var distributorStatus string
	err := s.db.QueryRowContext(ctx, `
SELECT status
FROM distributor_profiles
WHERE user_id = $1
`, userID).Scan(&distributorStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrDistributorNotEnabled
		}
		return "", fmt.Errorf("query distributor profile: %w", err)
	}
	if distributorStatus != "active" {
		return "", ErrDistributorNotEnabled
	}
	return PortalRoleDistributor, nil
}

// issueToken signs a short-lived JWT containing the portal role and user identity.
func (s *Service) issueToken(user User, portalRole PortalRole) (string, error) {
	claims := Claims{
		UserID:     user.ID,
		Email:      user.Email,
		PortalRole: portalRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken verifies a bearer token and returns its typed distributor claims.
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(_ *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

var (
	// ErrInvalidCredentials indicates the submitted email or password did not match.
	ErrInvalidCredentials = errors.New("账号或密码错误")
	// ErrInactiveUser blocks logins for users that exist but are not active.
	ErrInactiveUser = errors.New("当前账号未启用")
	// ErrDistributorNotEnabled means a non-admin user lacks an active distributor profile.
	ErrDistributorNotEnabled = errors.New("当前账号尚未开通分销商")
	// ErrInvalidToken indicates the presented JWT is missing, malformed, or unsigned by this service.
	ErrInvalidToken = errors.New("登录状态无效，请重新登录")
)
