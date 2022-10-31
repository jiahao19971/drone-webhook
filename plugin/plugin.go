package plugin

import (
	"context"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/sirupsen/logrus"
)

// New returns a new webhook extension.
func New() webhook.Plugin {
	return &plugin{}
}

type plugin struct {
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	logrus.Infof("Current build status %s", req.Build)
	logrus.Infof("Current action status %s", req.Action)

	return nil
}