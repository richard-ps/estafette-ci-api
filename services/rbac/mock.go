package rbac

import (
	"context"

	"github.com/estafette/estafette-ci-api/config"
	contracts "github.com/estafette/estafette-ci-contracts"
)

type MockService struct {
	GetProvidersFunc       func(ctx context.Context) (providers []*config.OAuthProvider, err error)
	GetProviderByNameFunc  func(ctx context.Context, name string) (provider *config.OAuthProvider, err error)
	CreateUserFunc         func(ctx context.Context, identity contracts.UserIdentity) (user *contracts.User, err error)
	UpdateUserFunc         func(ctx context.Context, user contracts.User) (err error)
	CreateGroupFunc        func(ctx context.Context, group contracts.Group) (insertedGroup *contracts.Group, err error)
	UpdateGroupFunc        func(ctx context.Context, group contracts.Group) (err error)
	CreateOrganizationFunc func(ctx context.Context, organization contracts.Organization) (insertedOrganization *contracts.Organization, err error)
	UpdateOrganizationFunc func(ctx context.Context, organization contracts.Organization) (err error)
	CreateClientFunc       func(ctx context.Context, client contracts.Client) (insertedClient *contracts.Client, err error)
	UpdateClientFunc       func(ctx context.Context, client contracts.Client) (err error)
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

func (s MockService) CreateGroup(ctx context.Context, group contracts.Group) (insertedGroup *contracts.Group, err error) {
	if s.CreateGroupFunc == nil {
		return
	}
	return s.CreateGroupFunc(ctx, group)
}

func (s MockService) UpdateGroup(ctx context.Context, group contracts.Group) (err error) {
	if s.UpdateGroupFunc == nil {
		return
	}
	return s.UpdateGroupFunc(ctx, group)
}

func (s MockService) CreateOrganization(ctx context.Context, organization contracts.Organization) (insertedOrganization *contracts.Organization, err error) {
	if s.CreateOrganizationFunc == nil {
		return
	}
	return s.CreateOrganizationFunc(ctx, organization)
}

func (s MockService) UpdateOrganization(ctx context.Context, organization contracts.Organization) (err error) {
	if s.UpdateOrganizationFunc == nil {
		return
	}
	return s.UpdateOrganizationFunc(ctx, organization)
}

func (s MockService) CreateClient(ctx context.Context, client contracts.Client) (insertedClient *contracts.Client, err error) {
	if s.CreateClientFunc == nil {
		return
	}
	return s.CreateClientFunc(ctx, client)
}

func (s MockService) UpdateClient(ctx context.Context, client contracts.Client) (err error) {
	if s.UpdateClientFunc == nil {
		return
	}
	return s.UpdateClientFunc(ctx, client)
}
