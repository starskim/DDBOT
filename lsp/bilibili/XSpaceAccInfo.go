package bilibili

import (
	"fmt"
	"github.com/starskim/DDBOT/proxy_pool"
	"github.com/starskim/DDBOT/requests"
	"github.com/starskim/DDBOT/utils"
	"go.uber.org/atomic"
	"io"
	"net/http/cookiejar"
	"time"
)

const (
	PathXSpaceAccInfo = "/x/space/wbi/acc/info"
)

type XSpaceAccInfoRequest struct {
	Mid         int64  `json:"mid"`
	Platform    string `json:"platform"`
	Token       string `json:"token"`
	WebLocation string `json:"web_location"`
}

var cj atomic.Pointer[cookiejar.Jar]

func refreshCookieJar() {
	j, _ := cookiejar.New(nil)
	err := requests.Get("https://bilibili.com", nil, io.Discard,
		requests.WithCookieJar(j),
		AddUAOption(),
		requests.RequestAutoHostOption(),
		requests.HeaderOption("accept", "application/json"),
		requests.HeaderOption("accept-language", "zh-CN,zh;q=0.9"),
	)
	if err != nil {
		logger.Errorf("bilibili: refreshCookieJar request error %v", err)
	}
	cj.Store(j)
}

func XSpaceAccInfo(mid int64) (*XSpaceAccInfoResponse, error) {
	st := time.Now()
	defer func() {
		ed := time.Now()
		logger.WithField("FuncName", utils.FuncName()).Tracef("cost %v", ed.Sub(st))
	}()
	url := BPath(PathXSpaceAccInfo)
	params, err := utils.ToDatas(&XSpaceAccInfoRequest{
		Mid:         mid,
		Platform:    "web",
		WebLocation: "1550101",
	})
	if err != nil {
		return nil, err
	}
	signWbi(params)
	var opts = []requests.Option{
		requests.ProxyOption(proxy_pool.PreferNone),
		requests.TimeoutOption(time.Second * 15),
		AddUAOption(),
		requests.HeaderOption("accept", "application/json"),
		requests.HeaderOption("accept-language", "zh-CN,zh;q=0.9"),
		requests.HeaderOption("origin", "https://space.bilibili.com"),
		requests.HeaderOption("referer", fmt.Sprintf("https://space.bilibili.com/%v", mid)),
		requests.HeaderOption("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		requests.RequestAutoHostOption(),
		requests.WithCookieJar(cj.Load()),
		requests.NotIgnoreEmptyOption(),
		delete412ProxyOption,
	}
	opts = append(opts, GetVerifyOption()...)
	xsai := new(XSpaceAccInfoResponse)
	err = requests.Get(url, params, xsai, opts...)
	if err != nil {
		return nil, err
	}
	return xsai, nil
}
