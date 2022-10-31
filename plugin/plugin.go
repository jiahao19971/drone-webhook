package plugin

import (
	"fmt"
	"context"
	"github.com/drone/drone-go/plugin/webhook"
)

// New returns a new webhook extension.
func New() webhook.Plugin {
	return &plugin{}
}

type plugin struct {
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Event == "build" {
		fmt.Print(req.Build)
		fmt.Print(req.Action)
	}
	return nil
}