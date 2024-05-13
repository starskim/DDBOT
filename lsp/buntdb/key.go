package buntdb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type KeyPatternFunc func(...interface{}) string

func Key(keys ...interface{}) string {
	var _keys []string
	for _, ikey := range keys {
		rk := reflect.ValueOf(ikey)
		if !rk.IsValid() {
			panic(fmt.Sprintf("invalid value %T %v", ikey, ikey))
		}
		if rk.Kind() == reflect.Ptr || rk.Kind() == reflect.Interface {
			rk = rk.Elem()
		}
		switch rk.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_keys = append(_keys, strconv.FormatInt(rk.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_keys = append(_keys, strconv.FormatUint(rk.Uint(), 10))
		case reflect.String:
			_keys = append(_keys, rk.String())
		case reflect.Bool:
			_keys = append(_keys, strconv.FormatBool(rk.Bool()))
		default:
			panic("unsupported key type " + reflect.ValueOf(ikey).Type().Name())
		}
	}
	return strings.Join(_keys, ":")
}

func NamedKey(name string, keys []interface{}) string {
	newkey := []interface{}{name}
	for _, key := range keys {
		newkey = append(newkey, key)
	}
	return Key(newkey...)
}

func BilibiliGroupConcernStateKey(keys ...interface{}) string {
	return NamedKey("ConcernState", keys)
}
func BilibiliGroupConcernConfigKey(keys ...interface{}) string {
	return NamedKey("ConcernConfig", keys)
}
func DouyuGroupConcernStateKey(keys ...interface{}) string {
	return NamedKey("DouyuConcernState", keys)
}
func DouyuGroupConcernConfigKey(keys ...interface{}) string {
	return NamedKey("DouyuConcernConfig", keys)
}
func YoutubeGroupConcernStateKey(keys ...interface{}) string {
	return NamedKey("YoutubeConcernState", keys)
}
func YoutubeGroupConcernConfigKey(keys ...interface{}) string {
	return NamedKey("YoutubeConcernConfig", keys)
}
func HuyaGroupConcernStateKey(keys ...interface{}) string {
	return NamedKey("HuyaConcernState", keys)
}
func HuyaGroupConcernConfigKey(keys ...interface{}) string {
	return NamedKey("HuyaConcernConfig", keys)
}

func PermissionKey(keys ...interface{}) string {
	return NamedKey("Permission", keys)
}
func BlockListKey(keys ...interface{}) string {
	return NamedKey("BlockList", keys)
}
func GroupPermissionKey(keys ...interface{}) string {
	return NamedKey("GroupPermission", keys)
}
func GroupEnabledKey(keys ...interface{}) string {
	return NamedKey("GroupEnable", keys)
}
func GlobalEnabledKey(keys ...interface{}) string {
	return NamedKey("GlobalEnable", keys)
}
func GroupMessageImageKey(keys ...interface{}) string {
	return NamedKey("GroupMessageImage", keys)
}
func GroupSilenceKey(keys ...interface{}) string {
	return NamedKey("GroupSilence", keys)
}
func GlobalSilenceKey(keys ...interface{}) string {
	return NamedKey("GlobalSilence", keys)
}
func GroupMuteKey(keys ...interface{}) string {
	return NamedKey("GroupMute", keys)
}
func GroupInvitorKey(keys ...interface{}) string {
	return NamedKey("GroupInventor", keys)
}

func LoliconPoolStoreKey(keys ...interface{}) string {
	return NamedKey("LoliconPoolStore", keys)
}

func ImageCacheKey(keys ...interface{}) string {
	return NamedKey("ImageCache", keys)
}

func ModeKey() string {
	return NamedKey("Mode", nil)
}
func NewFriendRequestKey(keys ...interface{}) string {
	return NamedKey("NewFriendRequest", keys)
}
func GroupInvitedKey(keys ...interface{}) string {
	return NamedKey("GroupInvited", keys)
}

func VersionKey(keys ...interface{}) string {
	return NamedKey("Version", keys)
}

func DDBotReleaseKey(keys ...interface{}) string {
	return NamedKey("DDBotReleaseKey", keys)
}

func DDBotNoUpdateKey(keys ...interface{}) string {
	return NamedKey("DDBotNoUpdateKey", keys)
}

func ParseConcernStateKeyWithInt64(key string) (groupCode int64, id int64, err error) {
	keys := strings.Split(key, ":")
	if len(keys) != 3 {
		return 0, 0, errors.New("invalid key")
	}
	groupCode, err = strconv.ParseInt(keys[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	id, err = strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return groupCode, id, nil
}
func ParseConcernStateKeyWithString(key string) (groupCode int64, id string, err error) {
	keys := strings.Split(key, ":")
	if len(keys) != 3 {
		return 0, "", errors.New("invalid key")
	}
	groupCode, err = strconv.ParseInt(keys[1], 10, 64)
	if err != nil {
		return 0, "", err
	}
	return groupCode, keys[2], nil

}
