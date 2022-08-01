package service

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/liucheyu/go-linebot-wallet/pkg/linebot/utils"
)

var dbUtils utils.DBUtils

func init() {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "123456789",
		Addr:                 "localhost:13306",
		Net:                  "tcp",
		DBName:               "linebot_wallet",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}

	// 釋放連線
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	dbUtils = utils.DBUtils{DB: db}
}

func TestGetDataRowsMap(t *testing.T) {
	rowSets := dbUtils.GetSqlRowSet("select item_type_id, item_type_name from item_type")

	for rowSets.Next() {
		fmt.Println("item_type_id", rowSets.GetInt("item_type_id"), "item_type_name", rowSets.GetString("item_type_name"))
	}
}
