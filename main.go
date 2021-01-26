package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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

func listServe(c echo.Context) error {
	defer c.Request().Body.Close()

	noticeList := comicDB.ListLoad(db)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.String(http.StatusOK, noticeList)
}

func listSize(c echo.Context) error {
	defer c.Request().Body.Close()

	noticeSize := comicDB.ListSize(db)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	return c.String(http.StatusOK, noticeSize)
}

func listContext(c echo.Context) error {
	defer c.Request().Body.Close()

	sid := c.Param("sid")
	noticeContext := comicDB.ListContext(db, sid)

	var li comicDB.NoticeInfo
	json.Unmarshal([]byte(noticeContext), &li)
	fmt.Println(li.SID)

	if li.SID != "" {
		return c.String(http.StatusOK, noticeContext)
	} else {
		return c.String(http.StatusBadRequest, noticeContext)
	}
}

func upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// html파일에 file을 받는 input 태그에 name 값이 키 값과 같아야 한다.
	files := form.File["files"]
	fmt.Println("files is : ", files)

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		mkfilepath := filepath.Join(`C:\savedata`, file.Filename)
		fmt.Println("path is : ", mkfilepath)
		dst, err := os.Create(mkfilepath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}

	return c.String(http.StatusOK, "upload ok")
}

func getFiles(c echo.Context, file string) {
	fmt.Println(file)
	c.Attachment(file, file)
}

func fileDownload(c echo.Context) error {
	directory := `C:\savedata`
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, fileName := range files {
		getFiles(c, fileName)
	}

	return c.String(http.StatusOK, "download ok")
}

func firstPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "")
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

	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = t

	// g := e.Group("/login")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))

	e.GET("/", firstPage)
	e.POST("/login/test", dataReceive)
	e.GET("/login/getTest", dataServe)
	e.GET("/list/getTest", listServe)
	e.GET("/list/getTest/:sid", listContext)
	e.GET("/list/getSize", listSize)
	e.POST("/upload", upload)
	e.GET("/download", fileDownload)

	e.Start(":8000")
}
