package rbac

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/estafette/estafette-ci-api/config"
	"github.com/estafette/estafette-ci-api/helpers"
	contracts "github.com/estafette/estafette-ci-contracts"
)

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(s Service) Service {
	return &loggingService{s, "oauth"}
}

type loggingService struct {
	Service
	prefix string
}

func (s *loggingService) GetProviders(ctx context.Context) (providers []*config.OAuthProvider, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GetProviders", err) }()

	return s.Service.GetProviders(ctx)
}

func (s *loggingService) GetProviderByName(ctx context.Context, name string) (provider *config.OAuthProvider, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GetProviderByName", err) }()

	return s.Service.GetProviderByName(ctx, name)
}

func (s *loggingService) GetUserByIdentity(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GetUserByIdentity", err, ErrUserNotFound) }()

	return s.Service.GetUserByIdentity(ctx, identity)
}

func (s *loggingService) GetUserByID(ctx context.Context, id string) (user *contracts.User, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GetUserByID", err, ErrUserNotFound) }()

	return s.Service.GetUserByID(ctx, id)
}

func (s *loggingService) CreateUser(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "CreateUser", err) }()

	return s.Service.CreateUser(ctx, identity)
}

func (s *loggingService) UpdateUser(ctx context.Context, user contracts.User) (err error) {
	defer func() { helpers.HandleLogError(s.prefix, "UpdateUser", err) }()

	return s.Service.UpdateUser(ctx, user)
}

func (s *loggingService) GenerateJWT(ctx context.Context, validDuration time.Duration, optionalClaims jwt.MapClaims) (tokenString string, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GenerateJWT", err) }()

	return s.Service.GenerateJWT(ctx, validDuration, optionalClaims)
}

func (s *loggingService) ValidateJWT(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "ValidateJWT", err) }()

	return s.Service.ValidateJWT(ctx, tokenString)
}

func (s *loggingService) GetClaimsFromJWT(ctx context.Context, tokenString string) (claims jwt.MapClaims, err error) {
	defer func() { helpers.HandleLogError(s.prefix, "GetClaimsFromJWT", err) }()

	return s.Service.GetClaimsFromJWT(ctx, tokenString)
}
