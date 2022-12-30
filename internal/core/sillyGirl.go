package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func initSillyGirl() {
	sillyGirl := BoltBucket("sillyGirl")
	id := protect(GenUUID(), "sillyGirl")
	updateSillyGirlValue(sillyGirl, "started_at", time.Now().Format("2006.01.02 15:04:05"), true)

	updateSillyGirlValue(sillyGirl, "api_key", time.Now().UnixNano(), false)
	updateSillyGirlValue(sillyGirl, "machineId", id, false)
	updateSillyGirlValue(sillyGirl, "name", "小小", false)
	updateSillyGirlValue(sillyGirl, "uuid", GenUUID(), false)
	updateSillyGirlValue(sillyGirl, "plugin_subcribe_addresses", "sub://T4EywWN46ztYBhHNdOl6Tkzap0mlmGqEMtMXHYFK3RePfCNugQWCjNHPpLZ8JoasT4VDcT9qG9TQFsqfbcA+SPnJAv0s+kH/KO/AX57Cx5vfJ8VDSI7d5JKehw8dp+GFnMbh2Gt1Dr/SSB304QDL/KPXvBfLvxc0USzCgHYjzVk=", false)
}

func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return hex.EncodeToString(mac.Sum(nil))
}

func updateSillyGirlValue(bucket BoltBucket, key string, value interface{}, isUpdate bool) {
	if len(bucket.GetString(key)) == 0 {
		_ = bucket.Set(key, value)
	} else {
		if isUpdate {
			_ = bucket.Set(key, value)
		}
	}
}
