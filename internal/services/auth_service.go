package services

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials    = errors.New("credenciais inválidas")
	ErrInvalidToken          = errors.New("token inválido")
	ErrPhoneAlreadyExists    = errors.New("telefone já cadastrado")
	ErrEmailAlreadyExists    = errors.New("email já cadastrado")
	ErrInactiveOrMissingUser = errors.New("usuário não encontrado")
)

type userStore interface {
	Create(ctx context.Context, user *models.User) error
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
}

type refreshTokenStore interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindActiveByTokenID(ctx context.Context, tokenID string) (*models.RefreshToken, error)
	RevokeByTokenID(ctx context.Context, tokenID string) error
}

type tokenClaims struct {
	TokenType string `json:"typ"`
	UserID    int64  `json:"uid"`
	Phone     string `json:"phone,omitempty"`
	jwt.RegisteredClaims
}

type AuthService struct {
	users         userStore
	refreshTokens refreshTokenStore
	jwtSecret     []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	now           func() time.Time
}

func NewAuthService(
	users userStore,
	refreshTokens refreshTokenStore,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		jwtSecret:     []byte(jwtSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		now:           time.Now,
	}
}

func (s *AuthService) Register(ctx context.Context, req dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createReq := dtos.CreateUserRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: string(hash),
	}
	user := createReq.ToModel()

	if err := s.users.Create(ctx, user); err != nil {
		if mappedErr := mapCreateUserError(err); mappedErr != nil {
			return nil, mappedErr
		}
		return nil, err
	}

	return s.issueAuthResponse(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.AuthResponse, error) {
	user, err := s.users.FindByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.Active {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return s.issueAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dtos.RefreshResponse, error) {
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	storedToken, err := s.refreshTokens.FindActiveByTokenID(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if storedToken.UserID != claims.UserID {
		return nil, ErrInvalidToken
	}

	if !storedToken.ExpiresOn.After(s.now()) {
		return nil, ErrInvalidToken
	}

	hash := hashToken(refreshToken)
	if subtle.ConstantTimeCompare([]byte(storedToken.TokenHash), []byte(hash)) != 1 {
		return nil, ErrInvalidToken
	}

	if err := s.refreshTokens.RevokeByTokenID(ctx, claims.ID); err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInactiveOrMissingUser
		}
		return nil, err
	}

	accessToken, accessExpiresOn, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, refreshTokenID, refreshExpiresOn, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokens.Create(ctx, &models.RefreshToken{
		UserID:    user.ID,
		TokenID:   refreshTokenID,
		TokenHash: hashToken(newRefreshToken),
		ExpiresOn: refreshExpiresOn,
	}); err != nil {
		return nil, err
	}

	return &dtos.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(time.Until(accessExpiresOn).Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return ErrInvalidToken
	}

	if err := s.refreshTokens.RevokeByTokenID(ctx, claims.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) issueAuthResponse(ctx context.Context, user *models.User) (*dtos.AuthResponse, error) {
	accessToken, accessExpiresOn, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenID, refreshExpiresOn, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokens.Create(ctx, &models.RefreshToken{
		UserID:    user.ID,
		TokenID:   refreshTokenID,
		TokenHash: hashToken(refreshToken),
		ExpiresOn: refreshExpiresOn,
	}); err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(accessExpiresOn).Seconds()),
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, time.Time, error) {
	now := s.now()
	expiresOn := now.Add(s.accessTTL)
	claims := tokenClaims{
		TokenType: "access",
		UserID:    user.ID,
		Phone:     user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expiresOn),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresOn, nil
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, string, time.Time, error) {
	now := s.now()
	expiresOn := now.Add(s.refreshTTL)
	tokenID := uuid.NewString()
	claims := tokenClaims{
		TokenType: "refresh",
		UserID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expiresOn),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, tokenID, expiresOn, nil
}

func (s *AuthService) parseRefreshToken(tokenValue string) (*tokenClaims, error) {
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != "refresh" || claims.ID == "" || claims.UserID <= 0 {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func mapCreateUserError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return nil
	}

	if pqErr.Code != "23505" {
		return nil
	}

	constraint := strings.ToLower(pqErr.Constraint)
	switch {
	case strings.Contains(constraint, "phone"):
		return ErrPhoneAlreadyExists
	case strings.Contains(constraint, "email"):
		return ErrEmailAlreadyExists
	default:
		return errors.New("usuário já existe")
	}
}
