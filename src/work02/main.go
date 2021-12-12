package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type UserInfo struct {
	Uid        int
	UserName   string
	Department string
	Created    string
}

var db *sql.DB
var err error

func main() {
	//连接数据库
	db, err = sql.Open("mysql", "root:123456ssw@/godb?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//调用测试
	var aUserInfo UserInfo
	aUserInfo, err = GetUserInfoById(1)
	fmt.Println(aUserInfo)
	fmt.Println(err)
}

func GetUserInfoById(uid int) (UserInfo, error) {
	//查询数据
	var userInfo UserInfo
	row := db.QueryRow("SELECT * FROM userinfo WHERE uid = ?", uid)
	err = row.Scan(&userInfo.Uid, &userInfo.UserName, &userInfo.Department, &userInfo.Created)

	//判断err类型是否是sql.ErrNoRows，用error.Wrapf抛出
	if errors.Is(err, sql.ErrNoRows) {
		return userInfo, errors.Wrapf(err, fmt.Sprintf("get userinfo null,uid is : %v", uid))
	}

	//抛出未知错误
	if err != nil {
		return userInfo, errors.Wrap(err, "get userinfo failed")
	}
	
	return userInfo, nil
}
