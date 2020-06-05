package rbac

import (
	"context"

	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
)

type MockService struct {
	GetProvidersFunc      func(ctx context.Context) (providers []*config.OAuthProvider, err error)
	GetProviderByNameFunc func(ctx context.Context, name string) (provider *config.OAuthProvider, err error)
	GetUserByIdentityFunc func(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error)
	GetUserByIDFunc       func(ctx context.Context, id string) (user *contracts.User, err error)
	CreateUserFunc        func(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error)
	UpdateUserFunc        func(ctx context.Context, user contracts.User) (err error)
}

func (s MockService) GetProviders(ctx context.Context) (providers []*config.OAuthProvider, err error) {
	if s.GetProvidersFunc == nil {
		return
	}
	return s.GetProvidersFunc(ctx)
}

func (s MockService) GetProviderByName(ctx context.Context, name string) (provider *config.OAuthProvider, err error) {
	if s.GetProviderByNameFunc == nil {
		return
	}
	return s.GetProviderByNameFunc(ctx, name)
}

func (s MockService) GetUserByIdentity(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	if s.GetUserByIdentityFunc == nil {
		return
	}
	return s.GetUserByIdentityFunc(ctx, identity)
}

func (s MockService) GetUserByID(ctx context.Context, id string) (user *contracts.User, err error) {
	if s.GetUserByIDFunc == nil {
		return
	}
	return s.GetUserByIDFunc(ctx, id)
}

func (s MockService) CreateUser(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error) {
	if s.CreateUserFunc == nil {
		return
	}
	return s.CreateUserFunc(ctx, identity)
}

func (s MockService) UpdateUser(ctx context.Context, user contracts.User) (err error) {
	if s.UpdateUserFunc == nil {
		return
	}
	return s.UpdateUserFunc(ctx, user)
}
