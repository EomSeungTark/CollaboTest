package comicDB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type UserInfo struct {
	ID       string `json:"id"`
	PASSWORD string `json:"password"`
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	DATE     string `json:"date"`
}

type NoticeInfo struct {
	TITLE   string `json:"title"`
	CONTEXT string `json:"context"`
	DATE    string `json:"date"`
	SECTION string `json:"section"`
}

func DBToString(rows *sql.Rows, length int, flag string) string {
	i := 0

	if flag == "DATA" {
		values := make([]UserInfo, length)
		for rows.Next() {
			rows.Scan(&values[i].ID, &values[i].PASSWORD, &values[i].NAME, &values[i].EMAIL, &values[i].DATE)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)

	} else if flag == "LIST" {
		values := make([]NoticeInfo, length)
		for rows.Next() {
			fmt.Print(rows)
			rows.Scan(&values[i].TITLE, &values[i].CONTEXT, &values[i].DATE, &values[i].SECTION)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)
	}

	return "없는 플레그 입니다."
}

func DataSave(db *sql.DB, userJson UserInfo) {
	nowTime := time.Now()
	toDbTime := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d", nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), nowTime.Minute(), nowTime.Second())
	dataTest := fmt.Sprintf("INSERT INTO USERS (ID, PASSWORD, NAME, EMAIL, DATE) VALUES ('%s', '%s', '%s', '%s', '%s')", userJson.ID, userJson.PASSWORD, userJson.NAME, userJson.EMAIL, toDbTime)
	row, err := db.Query(dataTest)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
}

func DataLoad(db *sql.DB) string {
	getUserSql := fmt.Sprint("SELECT * FROM USERS")
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	_ = db.QueryRow(`select count(*) from users`).Scan(&cnt)
	text := DBToString(rows, cnt, "DATA")

	return text
}

func ListLoad(db *sql.DB) string {
	getUserSql := fmt.Sprint("SELECT title, context, date, notice FROM NOTICE")
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	_ = db.QueryRow(`select count(*) from NOTICE`).Scan(&cnt)
	text := DBToString(rows, cnt, "LIST")

	return text
}
