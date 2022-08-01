package utils

import (
	"database/sql"
	"fmt"
	"strconv"
)

type DBUtils struct {
	DB *sql.DB
}

func (db DBUtils) GetSqlRowSet(sqlText string) SqlRowSet {
	rows, err := db.DB.Query(sqlText)

	if err != nil {
		panic(err)
	}

	cols, err := rows.Columns()

	if err != nil {
		panic(err)
	}

	return SqlRowSet{Rows: rows, Cols: cols}
}

type SqlRowSet struct {
	Rows   *sql.Rows
	Cols   []string
	RowMap map[string]interface{}
}

func (rs *SqlRowSet) Next() bool {
	if !rs.Rows.Next() {
		return false
	}

	columns := make([]interface{}, len(rs.Cols))
	columnPointers := make([]interface{}, len(rs.Cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := rs.Rows.Scan(columnPointers...); err != nil {
		// Todo 細部錯誤處理
		panic(err)
	}

	rs.RowMap = make(map[string]interface{})
	for i, colName := range rs.Cols {
		val := columnPointers[i].(*interface{})
		rs.RowMap[colName] = *val
	}

	//fmt.Println(rs.rowMap)

	return true
}

func (rs *SqlRowSet) GetString(columnName string) string {
	if rs.RowMap[columnName] == nil {
		// Todo 細部錯誤處理
		panic(fmt.Sprintf("columnName %s not found", columnName))
	}

	switch s := rs.RowMap[columnName].(type) {
	case string:
		return s
	case []uint8:
		return string(s)
	}

	panic(fmt.Sprintf("columnName %s not match type of %s", columnName, "string"))
}

func (rs *SqlRowSet) GetInt(columnName string) int {
	if rs.RowMap[columnName] == nil {
		// Todo 細部錯誤處理
		panic(fmt.Sprintf("columnName %s not found", columnName))
	}

	switch s := rs.RowMap[columnName].(type) {
	case int:
		return s
	case []uint8:
		// Todo []uint8轉int其他方式
		res, err := strconv.Atoi(string(s))
		// Todo 細部錯誤處理
		if err != nil {
			panic(err)
		}
		return res
	}

	panic(fmt.Sprintf("columnName %s not match type of %s", columnName, "int"))
}
