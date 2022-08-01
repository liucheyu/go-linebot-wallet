package service

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/cache"
)

type UserActionService interface {
	GetUserCacheMap (userID string) (map[string]string, error)
	CacheUserDataMap(key string, cacheMap map[string]string, expireTime int) error
	DeleteUserDataMapCache(userID string)
	SaveUserDataMapToDB(map[string]string)(int,error) 
}

type UserActionServiceImp struct {
	Context context.Context
	Redis *redis.Client
	DB *sql.DB
}

func (service *UserActionServiceImp) GetUserCacheMap(userID string) (map[string]string, error) {
	userDataMap,err := service.Redis.HGetAll(service.Context, cache.GetPrefixKey(cache.USER_CHACHE_MAP, userID)).Result()
	return userDataMap,err
}

func (service *UserActionServiceImp) CacheUserDataMap(userID string, cacheMap map[string]string, expireTime int) error {
	service.Redis.HSet(service.Context, cache.GetPrefixKey(cache.USER_CHACHE_MAP, userID), cacheMap)
	return nil
}

func (service *UserActionServiceImp) DeleteUserDataMapCache(userID string) {
	service.Redis.Del(service.Context, cache.GetPrefixKey(cache.USER_CHACHE_MAP, userID))
}

func (service *UserActionServiceImp) SaveUserDataMapToDB(userDataMap map[string]string) (int,error) {
	sql := `insert into pay_record (record_id, item_type_id, pay_mehod_id, item_name) 
	values((SELECT REPLACE (UUID(), "-", "")), ?, ?, ?) `	
	stmt,err := service.DB.Prepare(sql)

	if err != nil {
		return 0, err
	}

	res,err := stmt.Exec(userDataMap["itemType"], userDataMap["payMethod"], userDataMap["itemName"])
	affexted,_ :=res.RowsAffected()
	return int(affexted),err
}

