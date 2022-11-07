package mysql

import (
	"database/sql"
	"fmt"
)

type Mysql struct {
	Conf struct {
		Ip     string `json:"ip"`
		Port   string `json:"port"`
		Db     string `json:"db"`
		User   string `json:"user"`
		PassWd string `json:"passWd"`
	}
	DB *sql.DB `json:"-"`
}

// 初始化数据库
func (m *Mysql) Init() error {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&timeout=10s&readTimeout=10s",
		m.Conf.User,
		m.Conf.PassWd,
		m.Conf.Ip,
		m.Conf.Port,
		m.Conf.Db,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	m.DB = db
	return nil
}
