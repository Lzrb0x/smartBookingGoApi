package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserStore struct {
	nextID    int64
	byID      map[int64]*models.User
	byPhone   map[string]*models.User
	createErr error
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{
		nextID:  1,
		byID:    make(map[int64]*models.User),
		byPhone: make(map[string]*models.User),
	}
}

func (s *fakeUserStore) Create(_ context.Context, user *models.User) error {
	if s.createErr != nil {
		return s.createErr
	}

	if _, exists := s.byPhone[user.Phone]; exists {
		return errors.New("duplicate")
	}

	user.ID = s.nextID
	s.nextID++
	user.Active = true
	user.CreatedOn = time.Now()
	copyUser := *user
	s.byPhone[user.Phone] = &copyUser
	s.byID[user.ID] = &copyUser
	return nil
}

func (s *fakeUserStore) FindByPhone(_ context.Context, phone string) (*models.User, error) {
	user, ok := s.byPhone[phone]
	if !ok {
		return nil, sql.ErrNoRows
	}
	copyUser := *user
	return &copyUser, nil
}

func (s *fakeUserStore) FindByID(_ context.Context, id int64) (*models.User, error) {
	user, ok := s.byID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	copyUser := *user
	return &copyUser, nil
}

type fakeRefreshTokenStore struct {
	byTokenID map[string]*models.RefreshToken
}

func newFakeRefreshTokenStore() *fakeRefreshTokenStore {
	return &fakeRefreshTokenStore{byTokenID: make(map[string]*models.RefreshToken)}
}

func (s *fakeRefreshTokenStore) Create(_ context.Context, token *models.RefreshToken) error {
	copyToken := *token
	s.byTokenID[token.TokenID] = &copyToken
	return nil
}

func (s *fakeRefreshTokenStore) FindActiveByTokenID(_ context.Context, tokenID string) (*models.RefreshToken, error) {
	token, ok := s.byTokenID[tokenID]
	if !ok || token.RevokedOn != nil {
		return nil, sql.ErrNoRows
	}
	copyToken := *token
	return &copyToken, nil
}

func (s *fakeRefreshTokenStore) RevokeByTokenID(_ context.Context, tokenID string) error {
	token, ok := s.byTokenID[tokenID]
	if !ok {
		return nil
	}
	now := time.Now()
	token.RevokedOn = &now
	return nil
}

func TestAuthService_RegisterSuccess(t *testing.T) {
	service := NewAuthService(
		newFakeUserStore(),
		newFakeRefreshTokenStore(),
		"test-secret",
		15*time.Minute,
		24*time.Hour,
	)

	response, err := service.Register(context.Background(), dtos.RegisterRequest{
		Name:     "Luiz",
		Phone:    "85999999999",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.User == nil || response.User.ID == 0 {
		t.Fatalf("expected user returned with id, got %+v", response.User)
	}
	if response.AccessToken == "" || response.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	users := newFakeUserStore()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_ = users.Create(context.Background(), &models.User{
		UserIdentifier: "uid-1",
		Name:           "Luiz",
		Phone:          "85999999999",
		Password:       string(hash),
		IsComplete:     true,
	})

	service := NewAuthService(
		users,
		newFakeRefreshTokenStore(),
		"test-secret",
		15*time.Minute,
		24*time.Hour,
	)

	_, err = service.Login(context.Background(), dtos.LoginRequest{
		Phone:    "85999999999",
		Password: "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_RefreshSuccessRotatesToken(t *testing.T) {
	users := newFakeUserStore()
	refreshTokens := newFakeRefreshTokenStore()
	service := NewAuthService(users, refreshTokens, "test-secret", 15*time.Minute, 24*time.Hour)

	registerResponse, err := service.Register(context.Background(), dtos.RegisterRequest{
		Name:     "Luiz",
		Phone:    "85999999999",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	refreshResponse, err := service.Refresh(context.Background(), registerResponse.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	if refreshResponse.AccessToken == "" || refreshResponse.RefreshToken == "" {
		t.Fatal("expected new tokens on refresh")
	}

	if refreshResponse.RefreshToken == registerResponse.RefreshToken {
		t.Fatal("expected refresh rotation with a new refresh token")
	}

	_, err = service.Refresh(context.Background(), registerResponse.RefreshToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected old refresh token to be revoked, got %v", err)
	}
}

func TestAuthService_RefreshFailsWhenTokenRevoked(t *testing.T) {
	users := newFakeUserStore()
	refreshTokens := newFakeRefreshTokenStore()
	service := NewAuthService(users, refreshTokens, "test-secret", 15*time.Minute, 24*time.Hour)

	registerResponse, err := service.Register(context.Background(), dtos.RegisterRequest{
		Name:     "Luiz",
		Phone:    "85999999999",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := service.Logout(context.Background(), registerResponse.RefreshToken); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err = service.Refresh(context.Background(), registerResponse.RefreshToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for revoked refresh token, got %v", err)
	}
}

func TestAuthService_LogoutRevokesToken(t *testing.T) {
	users := newFakeUserStore()
	refreshTokens := newFakeRefreshTokenStore()
	service := NewAuthService(users, refreshTokens, "test-secret", 15*time.Minute, 24*time.Hour)

	registerResponse, err := service.Register(context.Background(), dtos.RegisterRequest{
		Name:     "Luiz",
		Phone:    "85999999999",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := service.Logout(context.Background(), registerResponse.RefreshToken); err != nil {
		t.Fatalf("expected no error on logout, got %v", err)
	}

	_, err = service.Refresh(context.Background(), registerResponse.RefreshToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected revoked token after logout, got %v", err)
	}
}
