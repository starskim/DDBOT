package huya

import (
	"github.com/starskim/DDBOT/lsp/concern"
)

func init() {
	concern.RegisterConcern(NewConcern(concern.GetNotifyChan()))
}
