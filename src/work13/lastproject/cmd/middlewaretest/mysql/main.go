package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	////todo ip修改
	dbhostsip  = "localhost:3306"
	dbusername = "shantest"
	dbpassword = "123456ssw"
	dbname     = "mydb"
)

type mysqlDB struct {
	db *sql.DB
}

func (f *mysqlDB) mysqlOpen() {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname)
	if err != nil {
		fmt.Println("mysql连接失败！")
		panic(err)
	}
	fmt.Println("mysql连接成功.....")
	f.db = db
}

func (f *mysqlDB) mysqlClose() {
	defer f.db.Close()
}

func (f *mysqlDB) mysqlSelect(sql_data string) { //select 查询数据
	fmt.Println("sql:", sql_data)
	rows, err := f.db.Query(sql_data)
	if err != nil {
		fmt.Println("查询失败")
	}
	for rows.Next() {
		var runoob_id int
		var runoob_title string
		var runoob_author string
		err = rows.Scan(&runoob_id, &runoob_title, &runoob_author)
		if err != nil {
			panic(err)
		}
		fmt.Println("runoob_id:", runoob_id)
		fmt.Println("runoob_title:", runoob_title)
		fmt.Println("runoob_author:", runoob_author)
	}

}

func (f *mysqlDB) mysqlInsert() { //insert  添加数据
	fmt.Println("开始插入......")
	stmt, err := f.db.Prepare("INSERT INTO runoob_tbl(runoob_title,runoob_author) VALUES(?,?)")
	//defer stmt.Close()
	if err != nil {
		fmt.Println("插入失败......")
		return
	}
	stmt.Exec("c++", "hhhssddd")
	fmt.Println("插入成功......")
}

func (f *mysqlDB) mysqlUpdate() { //update  修改数据
	stmt, err := f.db.Prepare("update runoob_tbl set runoob_title=?,runoob_author =? where runoob_id=?")
	//defer stmt.Close()
	if err != nil {
		fmt.Println("更新失败......")
		return
	}
	result, _ := stmt.Exec("c/c++", "asdfadsadsfa", 5)
	if result == nil {
		fmt.Println("修改失败")
	}
	affectCount, _ := result.RowsAffected() //返回影响的条数,注意有两个返回值
	fmt.Println(affectCount)
}

func (f *mysqlDB) mysqlDelete() { //delete  删除数据
	stmt, err := f.db.Prepare("delete from runoob_tbl where runoob_id=?")
	//defer stmt.Close()
	if err != nil {
		fmt.Println("删除失败......")
		return
	}
	stmt.Exec(5) //不返回任何结果
	fmt.Println("删除成功")
}

func (f *mysqlDB) mysqlTran() {
	//事务
	tx, err := f.db.Begin() //声明一个事务的开始
	if err != nil {
		fmt.Println(err)
		return
	}
	insertSql := "insert into runoob_tbl(runoob_title,runoob_author) VALUES(?,?)"
	insertStmt, insertErr := tx.Prepare(insertSql)
	if insertErr != nil {
		fmt.Println(insertErr)
		return
	}
	insertRes, insertErr := insertStmt.Exec("js", "ff点点点")
	last_insert_id, _ := insertRes.LastInsertId()
	fmt.Println(last_insert_id)
	// defer tx.Rollback()            //回滚之前上面的last_login_id是有的，但在回滚后该操作没有被提交，被回滚了，所以上面打印的Last_login_id的这条数据是不存在与数据库表中的
	tx.Commit() //这里提交了上面的操作，所以上面的执行的sql 会在数据库中产生一条数据
}

func main() {
	db := &mysqlDB{}
	db.mysqlOpen()
	db.mysqlInsert()
	db.mysqlUpdate()
	db.mysqlDelete()
	db.mysqlTran()
	db.mysqlSelect("select runoob_id,runoob_title,runoob_author from runoob_tbl")
	db.mysqlClose() //关闭
}
