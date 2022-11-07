package mysql

import (
	"database/sql"
	"encoding/json"
)

// 获得一个 stmt
func (m *Mysql) GetStmt(sql string) (*sql.Stmt, error) {
	return m.DB.Prepare(sql)
}

// 从数据库读单条 json 数据，并初始化到结构体
func (m *Mysql) ReadOneJson(stmt *sql.Stmt, v interface{}, queryArg ...interface{}) error {
	row := stmt.QueryRow(queryArg...)

	tmp := ""
	row.Scan(&tmp)

	err := json.Unmarshal([]byte(tmp), &v)
	if err != nil {
		return err
	}

	return nil
}

// 数据写到数据库中
func (m *Mysql) Write(stmt *sql.Stmt, queryArg ...interface{}) (int64, error) {
	res, err := stmt.Exec(queryArg...)
	if err != nil {
		return 0, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lid, nil
}
