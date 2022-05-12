package huya

import (
	"github.com/Sora233/DDBOT/internal/test"
	"github.com/Sora233/DDBOT/lsp/concern"
	"github.com/Sora233/DDBOT/lsp/mmsg/mt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newGroupLiveInfo(roomId string, live bool, liveStatusChanged bool, liveTitleChanged bool) *ConcernLiveNotify {
	li := &ConcernLiveNotify{
		LiveInfo: &LiveInfo{
			RoomId:            roomId,
			liveStatusChanged: liveStatusChanged,
			liveTitleChanged:  liveTitleChanged,
			IsLiving:          live,
		},
		Target: mt.NewGroupTarget(test.G1),
	}
	return li
}

func TestNewGroupConcernConfig(t *testing.T) {
	g := NewGroupConcernConfig(new(concern.ConcernConfig))
	assert.NotNil(t, g)
}

func TestGroupConcernConfig_ShouldSendHook(t *testing.T) {
	var notify = []concern.Notify{
		// 下播状态 什么也没变 不推
		newGroupLiveInfo(test.NAME1, false, false, false),
		// 下播状态 标题变了 不推
		newGroupLiveInfo(test.NAME1, false, false, true),
		// 下播了 检查配置
		newGroupLiveInfo(test.NAME1, false, true, false),
		// 下播了 检查配置
		newGroupLiveInfo(test.NAME1, false, true, true),
		// 直播状态 什么也没变 不推
		newGroupLiveInfo(test.NAME1, true, false, false),
		// 直播状态 改了标题 检查配置
		newGroupLiveInfo(test.NAME1, true, false, true),
		// 开播了 推
		newGroupLiveInfo(test.NAME1, true, true, false),
		// 开播了改了标题 推
		newGroupLiveInfo(test.NAME1, true, true, true),
	}

	var testCase = []*GroupConcernConfig{
		{
			IConfig: &concern.ConcernConfig{},
		},
		{
			IConfig: &concern.ConcernConfig{
				ConcernNotifyMap: map[mt.TargetType]*concern.ConcernNotifyConfig{
					mt.TargetGroup: {
						TitleChangeNotify: Live,
					},
				},
			},
		},
		{
			IConfig: &concern.ConcernConfig{
				ConcernNotifyMap: map[mt.TargetType]*concern.ConcernNotifyConfig{
					mt.TargetGroup: {
						OfflineNotify: Live,
					},
				},
			},
		},
		{
			IConfig: &concern.ConcernConfig{
				ConcernNotifyMap: map[mt.TargetType]*concern.ConcernNotifyConfig{
					mt.TargetGroup: {
						OfflineNotify:     Live,
						TitleChangeNotify: Live,
					},
				},
			},
		},
	}
	var expected = [][]bool{
		{
			false, false, false, false,
			false, false, true, true,
		},
		{
			false, false, false, false,
			false, true, true, true,
		},
		{
			false, false, true, true,
			false, false, true, true,
		},
		{
			false, false, true, true,
			false, true, true, true,
		},
	}
	assert.Equal(t, len(expected), len(testCase))
	for index1, g := range testCase {
		assert.Equal(t, len(expected[index1]), len(notify))
		for index2, liveInfo := range notify {
			result := g.ShouldSendHook(liveInfo)
			assert.NotNil(t, result)
			assert.Equal(t, expected[index1][index2], result.Pass)
		}
	}
}

func TestGroupConcernConfig_AtBeforeHook(t *testing.T) {
	var notify = []concern.Notify{
		// 下播状态 什么也没变 不推
		newGroupLiveInfo(test.NAME1, false, false, false),
		// 下播状态 标题变了 不推
		newGroupLiveInfo(test.NAME1, false, false, true),
		// 下播了 检查配置
		newGroupLiveInfo(test.NAME1, false, true, false),
		// 下播了 检查配置
		newGroupLiveInfo(test.NAME1, false, true, true),
		// 直播状态 什么也没变 不推
		newGroupLiveInfo(test.NAME1, true, false, false),
		// 直播状态 改了标题 检查配置
		newGroupLiveInfo(test.NAME1, true, false, true),
		// 开播了 推
		newGroupLiveInfo(test.NAME1, true, true, false),
		// 开播了改了标题 推
		newGroupLiveInfo(test.NAME1, true, true, true),
	}
	var expcted = []bool{
		false, false, false, false, false, false, true, true,
	}
	var config = &GroupConcernConfig{IConfig: &concern.ConcernConfig{}}
	for idx, n := range notify {
		hook := config.AtBeforeHook(n)
		assert.EqualValues(t, expcted[idx], hook.Pass)
	}
}
