package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xerror "github.com/pkg/errors"
	"log"
	"time"
)

var (
	DB *sql.DB
)

//user表结构体定义
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func dao(id int) (User, error) {
	var user User
	var err error
	row := DB.QueryRow("select 1 as 'id' , 'abc' as name where 1 = ?", id)
	if err = row.Scan(&user.ID, &user.Name); err != nil {
		err = xerror.Wrapf(err, "cannot find user with id %d", id)
	}
	return user, err
}

func service(id int) (User, error) {
	return dao(id)
}

//查询用户
func Biz(id int) string {
	var user User
	var err error
	user, err = service(id)
	if err != nil {
		log.Print(err)
		return serviceErrorMsg(err)
	}
	return fmt.Sprintf("%+v", user)
}

func main() {
	fmt.Println(Biz(1)) // 正常
	fmt.Println(Biz(2)) // 异常
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(100)
}

// 返回客户端信息包装
func serviceErrorMsg(err error) string {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "no match data"
	default:
		return err.Error()
	}
}
