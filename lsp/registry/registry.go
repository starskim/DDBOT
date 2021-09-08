package registry

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Sora233/DDBOT/lsp/concern"
	"golang.org/x/sync/errgroup"
)

var logger = utils.GetModuleLogger("registry")

type option struct {
}

type OptFunc func(opt *option) *option

type ConcernCenter struct {
	M map[string]map[concern.Type]concern.Concern
}

var globalCenter = newConcernCenter()

func newConcernCenter() *ConcernCenter {
	cc := new(ConcernCenter)
	cc.M = make(map[string]map[concern.Type]concern.Concern)
	return cc
}

func RegisterConcernManager(c concern.Concern, concernType []concern.Type, opts ...OptFunc) {
	site := c.Site()
	for _, ctype := range concernType {
		if !ctype.IsTrivial() {
			panic(fmt.Sprintf("Concern %v: Type %v IsTrivial() must be True", site, ctype))
		}
	}
	if _, found := globalCenter.M[site]; !found {
		globalCenter.M[site] = make(map[concern.Type]concern.Concern)
	}
	for _, ctype := range concernType {
		if lastC, found := globalCenter.M[site][ctype]; !found {
			globalCenter.M[site][ctype] = c
		} else {
			logger.Errorf("Concern %v: Type %v is already registered by Concern %v, skip.", site, ctype, lastC.Site())
		}
	}
}

func StartAll() error {
	all := ListConcernManager()
	errG := errgroup.Group{}
	for _, c := range all {
		errG.Go(func() error {
			return c.Start()
		})
	}
	return errG.Wait()
}

func StopAll() {
	all := ListConcernManager()
	for _, c := range all {
		c.Stop()
	}
}

func ListConcernManager() []concern.Concern {
	var resultMap = make(map[concern.Concern]interface{})
	for _, cmap := range globalCenter.M {
		for _, c := range cmap {
			resultMap[c] = struct{}{}
		}
	}
	var result []concern.Concern
	for k := range resultMap {
		result = append(result, k)
	}
	return result
}

func GetConcernManager(site string, ctype concern.Type) concern.Concern {
	if sub, found := globalCenter.M[site]; !found {
		return nil
	} else {
		return sub[ctype]
	}
}
