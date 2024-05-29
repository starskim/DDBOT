package bilibili

import (
	"github.com/starskim/DDBOT/proxy_pool"
	"github.com/starskim/DDBOT/requests"
	"github.com/starskim/DDBOT/utils"
	"time"
)

const (
	PathGetAttentionList = "/feed/v1/feed/get_attention_list"
)

type GetAttentionListRequest struct {
	Uid int64 `json:"uid"`
}

func GetAttentionList() (*GetAttentionListResponse, error) {
	if !IsVerifyGiven() {
		return nil, ErrVerifyRequired
	}
	st := time.Now()
	defer func() {
		ed := time.Now()
		logger.WithField("FuncName", utils.FuncName()).Tracef("cost %v", ed.Sub(st))
	}()
	url := BPath(PathGetAttentionList)
	params, err := utils.ToParams(&GetAttentionListRequest{
		Uid: accountUid.Load(),
	})
	if err != nil {
		return nil, err
	}
	var opts []requests.Option
	opts = append(opts,
		requests.ProxyOption(proxy_pool.PreferNone),
		AddUAOption(),
		requests.TimeoutOption(time.Second*10),
		delete412ProxyOption,
	)
	opts = append(opts, GetVerifyOption()...)
	getAttentionListResp := new(GetAttentionListResponse)
	err = requests.Get(url, params, getAttentionListResp, opts...)
	if err != nil {
		return nil, err
	}
	return getAttentionListResp, nil
}
