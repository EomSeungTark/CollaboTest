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

func ToString(rows *sql.Rows, length int) string {
	values := make([]UserInfo, length)

	fmt.Println(length)

	i := 0

	for rows.Next() {
		rows.Scan(&values[i].ID, &values[i].PASSWORD, &values[i].NAME, &values[i].EMAIL, &values[i].DATE)
		i++
	}

	j, _ := json.Marshal(values)

	return string(j)
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
	text := ToString(rows, cnt)

	return text
}
