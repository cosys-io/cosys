package apis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cosys-io/cosys/internal/common"
)

type Controller struct {
	Actions map[string]Action
}

func (c *Controller) Action(uid string) (Action, error) {
	action := c.Actions[uid]
	if action == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return action, nil
}

type Action func(common.Cosys, context.Context) http.HandlerFunc

func NewController(actions map[string]Action) *Controller {
	return &Controller{actions}
}
