package enum

import(
	"fmt"
)

type redisPrefix string

const(
	PREFIX_ACTION_MAP redisPrefix = "actionmap"
	PREFIX_LAST_ACTION redisPrefix = "lastaction"
)

func GetRedisKeyWithPrefix(prefixName redisPrefix, value string) string {
	return fmt.Sprintf("%v_%v", prefixName, value)
}
