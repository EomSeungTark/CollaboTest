package comicDB

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
	SID       string `json:"sid"`
	TITLE     string `json:"title"`
	CONTEXT   string `json:"context"`
	USERID    string `json:"user_id"`
	DATE      string `json:"date"`
	VIEWCOUNT string `json:"view_count"`
	SECTION   string `json:"section"`
}

func DBToString(rows *sql.Rows, length int, flag string) string {
	var i int = 0

	if flag == "DATA" {
		values := make([]UserInfo, length)
		for rows.Next() {
			rows.Scan(&values[i].ID, &values[i].PASSWORD, &values[i].NAME, &values[i].EMAIL, &values[i].DATE)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)

	} else if flag == "NOTICE" {
		values := make([]NoticeInfo, length)
		for rows.Next() {
			rows.Scan(&values[i].SID, &values[i].TITLE, &values[i].CONTEXT, &values[i].USERID, &values[i].DATE, &values[i].VIEWCOUNT, &values[i].SECTION)
			i++
		}
		j, _ := json.Marshal(values)

		return string(j)
	} else if flag == "NOTICE_ONE" {
		value := NoticeInfo{}
		for rows.Next() {
			rows.Scan(&value.SID, &value.TITLE, &value.CONTEXT, &value.USERID, &value.DATE, &value.VIEWCOUNT, &value.SECTION)
			i++
		}
		j, _ := json.Marshal(value)

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
	getUserSql := fmt.Sprint("SELECT * FROM NOTICE ORDER BY SECTION DESC, DATE DESC")
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	_ = db.QueryRow(`select count(*) from NOTICE`).Scan(&cnt)
	text := DBToString(rows, cnt, "NOTICE")

	return text
}

func ListSize(db *sql.DB) string {
	var cnt int
	_ = db.QueryRow(`select count(*) from NOTICE`).Scan(&cnt)

	return strconv.Itoa(cnt)
}

func ListContext(db *sql.DB, sid string) string {
	sqlState := fmt.Sprintf("update notice set VIEW_COUNT=VIEW_COUNT+1 where SID=%s", sid)
	_, _ = db.Query(sqlState)

	getUserSql := fmt.Sprintf("SELECT * FROM NOTICE WHERE SID=%s", sid)
	rows, err := db.Query(getUserSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	text := DBToString(rows, 1, "NOTICE_ONE")

	return text
}
