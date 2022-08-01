package service

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/utils"
)

type BaseDataCache struct {
	Context context.Context
	Redis   *redis.Client
	DB      *sql.DB
}

func (base BaseDataCache) GetItemTypesMap() map[string]string {
	itemMap, err := base.Redis.HGetAll(base.Context, "itemTypes_map").Result()

	if err != nil || len(itemMap) == 0 {
		itemMap = map[string]string{}
		sqlText := "select item_type_id, item_type_name from item_type"
		dbUtils := utils.DBUtils{DB: base.DB}
		rowSets := dbUtils.GetSqlRowSet(sqlText)

		for rowSets.Next() {
			itemMap["item_type_id"] = rowSets.GetString("item_type_id")
			itemMap["item_type_name"] = rowSets.GetString("item_type_name")
		}

		base.Redis.HSet(base.Context, "itemTypes_map", itemMap)
	}

	return itemMap
}

func (base BaseDataCache) GetPayMethodMap() map[string]string {
	itemMap, err := base.Redis.HGetAll(base.Context, "paymethod_map").Result()

	if err != nil || len(itemMap) == 0 {
		itemMap = map[string]string{}
		sqlText := "select pay_method_id, pay_method_name from pay_method"
		dbUtils := utils.DBUtils{DB: base.DB}
		rowSets := dbUtils.GetSqlRowSet(sqlText)

		for rowSets.Next() {
			itemMap["pay_method_id"] = rowSets.GetString("pay_method_id")
			itemMap["pay_method_name"] = rowSets.GetString("pay_method_name")
		}

		base.Redis.HSet(base.Context, "paymethod_map", itemMap)
	}

	return itemMap
}
