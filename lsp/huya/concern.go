package huya

import (
	"fmt"
	"github.com/Sora233/DDBOT/lsp/concern"
	"github.com/Sora233/DDBOT/lsp/concern_type"
	"github.com/Sora233/DDBOT/lsp/mmsg"
	"github.com/Sora233/DDBOT/lsp/mmsg/mt"
	localutils "github.com/Sora233/DDBOT/utils"
	"github.com/Sora233/MiraiGo-Template/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger = utils.GetModuleLogger("huya-concern")

const (
	Live concern_type.Type = "live"
)

type Concern struct {
	*StateManager
}

func (c *Concern) Site() string {
	return Site
}

func (c *Concern) Types() []concern_type.Type {
	return []concern_type.Type{Live}
}

func (c *Concern) ParseId(s string) (interface{}, error) {
	return s, nil
}

func (c *Concern) GetStateManager() concern.IStateManager {
	return c.StateManager
}

func (c *Concern) Stop() {
	logger.Trace("正在停止huya concern")
	logger.Trace("正在停止huya StateManager")
	c.StateManager.Stop()
	logger.Trace("huya StateManager已停止")
	logger.Trace("huya concern已停止")
}

func (c *Concern) Start() error {
	c.UseEmitQueue()
	c.StateManager.UseNotifyGeneratorFunc(c.notifyGenerator())
	c.StateManager.UseFreshFunc(c.fresh())
	return c.StateManager.Start()
}

func (c *Concern) Add(ctx mmsg.IMsgCtx, target mt.Target, id interface{}, ctype concern_type.Type) (concern.IdentityInfo, error) {
	var err error
	log := logger.WithFields(localutils.TargetFields(target)).WithField("id", id)

	err = c.StateManager.CheckTargetConcern(target, id, ctype)
	if err != nil {
		return nil, err
	}

	liveInfo, err := RoomPage(id.(string))
	if err != nil {
		log.Errorf("RoomPage error %v", err)
		return nil, fmt.Errorf("查询房间信息失败 %v - %v", id, err)
	}
	_, err = c.StateManager.AddTargetConcern(target, id, ctype)
	if err != nil {
		return nil, err
	}
	return liveInfo, nil
}

func (c *Concern) Remove(ctx mmsg.IMsgCtx, target mt.Target, _id interface{}, ctype concern_type.Type) (concern.IdentityInfo, error) {
	id := _id.(string)
	identity, _ := c.Get(id)
	_, err := c.StateManager.RemoveTargetConcern(target, id, ctype)
	_ = c.RWCoverTx(func(tx *buntdb.Tx) error {
		allCtype, err := c.GetConcern(id)
		if err != nil {
			return err
		}
		if allCtype.Empty() {
			err = c.DeleteLiveInfo(id)
		}
		return err
	})
	return identity, err
}

func (c *Concern) Get(id interface{}) (concern.IdentityInfo, error) {
	liveInfo, err := c.FindRoom(id.(string), false)
	if err != nil {
		return nil, err
	}
	return concern.NewIdentity(liveInfo.RoomId, liveInfo.GetName()), nil
}

func (c *Concern) FindRoom(roomId string, load bool) (*LiveInfo, error) {
	var liveInfo *LiveInfo
	if load {
		var err error
		liveInfo, err = RoomPage(roomId)
		if err != nil {
			return nil, err
		}
		_ = c.StateManager.AddLiveInfo(liveInfo)
	}
	if liveInfo != nil {
		return liveInfo, nil
	}
	return c.StateManager.GetLiveInfo(roomId)
}

func (c *Concern) FindOrLoadRoom(roomId string) (*LiveInfo, error) {
	info, _ := c.FindRoom(roomId, false)
	if info == nil {
		return c.FindRoom(roomId, true)
	}
	return info, nil
}

func (c *Concern) notifyGenerator() concern.NotifyGeneratorFunc {
	return func(target mt.Target, event concern.Event) []concern.Notify {
		switch info := event.(type) {
		case *LiveInfo:
			if info.Living() {
				info.Logger().WithFields(localutils.TargetFields(target)).Trace("living notify")
			} else {
				info.Logger().WithFields(localutils.TargetFields(target)).Trace("noliving notify")
			}
			return []concern.Notify{NewConcernLiveNotify(target, info)}
		default:
			logger.Errorf("unknown EventType %+v", event)
			return nil
		}
	}
}

func (c *Concern) fresh() concern.FreshFunc {
	return c.EmitQueueFresher(func(ctype concern_type.Type, id interface{}) ([]concern.Event, error) {
		var result []concern.Event
		roomid := id.(string)
		if ctype.ContainAll(Live) {
			oldInfo, _ := c.FindRoom(roomid, false)
			liveInfo, err := c.FindRoom(roomid, true)
			if err == ErrRoomNotExist || err == ErrRoomBanned {
				logger.WithFields(logrus.Fields{
					"RoomId":   roomid,
					"RoomName": oldInfo.GetName(),
				}).Warn("直播间不存在或被封禁，订阅将失效")
				c.RemoveAllById(id)
				return nil, err
			}
			if err != nil {
				return nil, fmt.Errorf("load liveinfo failed %v", err)
			}
			if oldInfo == nil {
				liveInfo.liveStatusChanged = true
			} else {
				if oldInfo.Living() != liveInfo.Living() {
					liveInfo.liveStatusChanged = true
				}
				if oldInfo.RoomName != liveInfo.RoomName {
					liveInfo.liveTitleChanged = true
				}
			}
			result = append(result, liveInfo)
		}
		return result, nil
	})
}

func NewConcern(notify chan<- concern.Notify) *Concern {
	c := &Concern{
		StateManager: NewStateManager(notify),
	}
	return c
}
