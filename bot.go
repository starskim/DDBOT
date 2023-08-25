package DDBOT

import (
	"fmt"
	"github.com/Sora233/DDBOT/lsp"
	"github.com/Sora233/DDBOT/warn"
	"github.com/starskim/MiraiGo-Template/bot"
	"github.com/starskim/MiraiGo-Template/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "github.com/Sora233/DDBOT/logging"
	_ "github.com/Sora233/DDBOT/lsp/acfun"
	_ "github.com/Sora233/DDBOT/lsp/douyu"
	_ "github.com/Sora233/DDBOT/lsp/huya"
	_ "github.com/Sora233/DDBOT/lsp/twitcasting"
	_ "github.com/Sora233/DDBOT/lsp/weibo"
	_ "github.com/Sora233/DDBOT/lsp/youtube"
	_ "github.com/Sora233/DDBOT/msg-marker"
)

// SetUpLog 使用默认的日志格式配置，会写入到logs文件夹内，日志会保留七天
func SetUpLog() {
	writer, err := rotatelogs.New(
		path.Join("logs", "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
	})
	logrus.AddHook(lfshook.NewHook(writer, &logrus.TextFormatter{
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		ForceQuote:       true,
	}))
}

// Run 启动bot，这个函数会阻塞直到收到退出信号
func Run() {
	if fi, err := os.Stat("device.json"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("警告：没有检测到device.json，正在生成，如果是第一次运行，可忽略")
			bot.GenRandomDevice()
		} else {
			warn.Warn(fmt.Sprintf("检查device.json文件失败 - %v", err))
			os.Exit(1)
		}
	} else {
		if fi.IsDir() {
			warn.Warn("检测到device.json，但目标是一个文件夹！请手动确认并删除该文件夹！")
			os.Exit(1)
		} else {
			fmt.Println("检测到device.json，使用存在的device.json")
		}
	}

	if fi, err := os.Stat("application.yaml"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("警告：没有检测到配置文件application.yaml，正在生成，如果是第一次运行，可忽略")
			if err := ioutil.WriteFile("application.yaml", []byte(exampleConfig), 0755); err != nil {
				warn.Warn(fmt.Sprintf("application.yaml生成失败 - %v", err))
				os.Exit(1)
			} else {
				fmt.Println("最小配置application.yaml已生成，请按需修改，如需高级配置请查看帮助文档")
			}
		} else {
			warn.Warn(fmt.Sprintf("检查application.yaml文件失败 - %v", err))
			os.Exit(1)
		}
	} else {
		if fi.IsDir() {
			warn.Warn("检测到application.yaml，但目标是一个文件夹！请手动确认并删除该文件夹！")
			os.Exit(1)
		} else {
			fmt.Println("检测到application.yaml，使用存在的application.yaml")
		}
	}

	config.GlobalConfig.SetConfigName("application")
	config.GlobalConfig.SetConfigType("yaml")
	config.GlobalConfig.AddConfigPath(".")
	config.GlobalConfig.AddConfigPath("./config")

	err := config.GlobalConfig.ReadInConfig()
	if err != nil {
		warn.Warn(fmt.Sprintf("读取配置文件失败！请检查配置文件格式是否正确 - %v", err))
		os.Exit(1)
	}
	config.GlobalConfig.WatchConfig()

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 登录
	bot.Login()

	// 刷新好友列表，群列表
	bot.RefreshList()

	lsp.Instance.PostStart(bot.Instance)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	bot.Stop()
}

var exampleConfig = func() string {
	s := `
# 注意，填写时请把井号及后面的内容删除，并且冒号后需要加一个空格
sign:
  # 数据包的签名服务器
  # 兼容 https://github.com/fuqiuluo/unidbg-fetch-qsign
  # 如果遇到 登录 45 错误, 或者发送信息风控的话需要填入一个服务器
  # 示例:
  # server: 'http://127.0.0.1:8080' # 本地签名服务器
  # server: 'https://signserver.example.com' # 线上签名服务器
  # 服务器可使用docker在本地搭建或者使用他人开放的服务
  server: ''
  # 签名服务器认证 Bearer Token
  # 使用开放的服务可能需要提供此 Token 进行认证
  server-bearer: ''
  # 如果签名服务器的版本在1.1.0及以下, 请将下面的参数改成true
  is-below-110: false
  # 签名服务器所需要的apikey, 如果签名服务器的版本在1.1.0及以下则此项无效
  # 本地部署的默认为114514
  key: '114514'
  # 在实例可能丢失（获取到的签名为空）时是否尝试重新注册
  # 为 true 时，在签名服务不可用时可能每次发消息都会尝试重新注册并签名。
  # 为 false 时，将不会自动注册实例，在签名服务器重启或实例被销毁后需要重启 go-cqhttp 以获取实例
  # 否则后续消息将不会正常签名。关闭此项后可以考虑开启签名服务器端 auto_register 避免需要重启
  auto-register: false
  # 是否在 token 过期后立即自动刷新签名 token（在需要签名时才会检测到，主要防止 token 意外丢失）
  # 独立于定时刷新
  auto-refresh-token: false
  # 定时刷新 token 间隔时间，单位为分钟, 建议 30~40 分钟, 不可超过 60 分钟
  # 目前丢失token也不会有太大影响，可设置为 0 以关闭，推荐开启
  refresh-interval: 40
bot:
  account:  # 你bot的qq号，不填则使用扫码登陆
  password: # 你bot的qq密码
  onJoinGroup: 
    rename: "【bot】"  # BOT进群后自动改名，默认改名为“【bot】”，如果留空则不自动改名

# 初次运行时将不使用b站帐号方便进行测试
# 如果不使用b站帐号，则推荐订阅数不要超过5个，否则推送延迟将上升
# b站相关的功能推荐配置一个b站账号，建议使用小号
# bot将使用您b站帐号的以下功能：
# 关注用户 / 取消关注用户 / 查看关注列表
# 请注意，订阅一个账号后，此处使用的b站账号将自动关注该账号
bilibili:
  SESSDATA: # 你的b站cookie
  bili_jct: # 你的b站cookie
  interval: 25s

concern:
  emitInterval: 5s

logLevel: info
`
	// win上用记事本打开不会正确换行
	if runtime.GOOS == "windows" {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return s
}()
