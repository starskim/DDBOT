package huya

import (
	"github.com/starskim/DDBOT/lsp/concern"
)

type GroupConcernConfig struct {
	concern.IConfig
}

func NewGroupConcernConfig(g concern.IConfig) *GroupConcernConfig {
	return &GroupConcernConfig{g}
}
