package cache

import "fmt"

type redisPrefix string

const(
	USER_CHACHE_MAP redisPrefix = "userCacheMap"
)

func GetPrefixKey(prefixName redisPrefix, value string) string {
	return fmt.Sprintf("%v_%v", prefixName, value)
}