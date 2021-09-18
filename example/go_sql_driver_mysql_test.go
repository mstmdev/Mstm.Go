package example

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func mysqlDemo(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/usersdb2015?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
	}
	rows, err1 := db.Query("select UserId,UserName from userinfo")
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	defer rows.Close()
	fmt.Println("")
	cols, _ := rows.Columns()
	for i := range cols {
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	fmt.Println("")
	var UserId int64
	var UserName string
	for rows.Next() {
		if err := rows.Scan(&UserId, &UserName); err == nil {
			fmt.Print(UserId)
			fmt.Print("\t")
			fmt.Print(UserName)
			fmt.Print("\t\r\n")
		}
	}
}
