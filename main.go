package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/eom/collabotest/comicDB"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "800326"
	DB_NAME     = "postgres"
)

func dataReceive(c echo.Context) error {
	userinfo := comicDB.UserInfo{}

	defer c.Request().Body.Close()
	byte, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(byte, &userinfo)
	fmt.Println(userinfo.ID)

	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	comicDB.DataSave(db, userinfo)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.String(http.StatusOK, "test ok")
}

func dataServe(c echo.Context) error {
	defer c.Request().Body.Close()

	userList := comicDB.DataLoad(db)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.String(http.StatusOK, userList)
}

func main() {
	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	fmt.Println(reflect.TypeOf(db))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Welcome to the server")

	e := echo.New()

	// g := e.Group("/login")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))

	e.POST("/login/test", dataReceive)
	e.GET("/login/getTest", dataServe)

	e.Start(":8000")
}
