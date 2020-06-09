package cloudsource

import (
	"context"

	"github.com/estafette/estafette-ci-api/clients/cloudsourceapi"
)

type MockService struct {
	CreateJobForCloudSourcePushFunc func(ctx context.Context, notification cloudsourceapi.PubSubNotification) (err error)
	IsWhitelistedProjectFunc        func(notification cloudsourceapi.PubSubNotification) (isWhiteListed bool)
}

func (s MockService) CreateJobForCloudSourcePush(ctx context.Context, notification cloudsourceapi.PubSubNotification) (err error) {
	if s.CreateJobForCloudSourcePushFunc == nil {
		return
	}
	return s.CreateJobForCloudSourcePushFunc(ctx, notification)
}

func (s MockService) IsWhitelistedProject(notification cloudsourceapi.PubSubNotification) (isWhiteListed bool) {
	if s.IsWhitelistedProjectFunc == nil {
		return
	}
	return s.IsWhitelistedProjectFunc(notification)
}
