package auth

import "testing"

func TestIssueAndParseToken(t *testing.T) {
	t.Parallel()

	service := NewService(nil, "unit-test-secret")
	user := User{
		ID:    42,
		Email: "distributor@example.com",
		Role:  "user",
	}

	token, err := service.issueToken(user, PortalRoleDistributor)
	if err != nil {
		t.Fatalf("issueToken() error = %v", err)
	}

	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("claims.UserID = %d, want %d", claims.UserID, user.ID)
	}
	if claims.Email != user.Email {
		t.Fatalf("claims.Email = %q, want %q", claims.Email, user.Email)
	}
	if claims.PortalRole != PortalRoleDistributor {
		t.Fatalf("claims.PortalRole = %q, want %q", claims.PortalRole, PortalRoleDistributor)
	}
}

func TestParseTokenRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	service := NewService(nil, "unit-test-secret")

	if _, err := service.ParseToken("not-a-token"); err != ErrInvalidToken {
		t.Fatalf("ParseToken() error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestParseTokenRejectsTokenSignedWithDifferentSecret(t *testing.T) {
	t.Parallel()

	issuer := NewService(nil, "issuer-secret")
	consumer := NewService(nil, "consumer-secret")

	token, err := issuer.issueToken(User{ID: 7, Email: "ops@example.com", Role: "admin"}, PortalRoleOperator)
	if err != nil {
		t.Fatalf("issueToken() error = %v", err)
	}

	if _, err := consumer.ParseToken(token); err != ErrInvalidToken {
		t.Fatalf("ParseToken() error = %v, want %v", err, ErrInvalidToken)
	}
}
