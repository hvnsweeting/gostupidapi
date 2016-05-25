package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	// log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TaiLieuHocTap(c echo.Context) error {
	return c.String(http.StatusForbidden, "Forbidden")
}

func LoginHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, c.QueryParams())
}

type MyData struct {
	Password []string
	User     []string `json:"user_name"`
}

type ReturnData struct {
	Password []string
	User     []string
}

func JsonEndpoint(c echo.Context) error {
	r := c.Request()

	data, err := ioutil.ReadAll(r.Body())
	log.Println(string(data))
	var myData MyData
	if err == nil && data != nil {
		err = json.Unmarshal(data, &myData)
		if err != nil {
			log.Println("Marshalling: ", err)
			return c.JSON(http.StatusBadRequest, "Bad JSON")
		}
	}

	retData := ReturnData{User: myData.User, Password: myData.Password}
	return c.JSON(http.StatusOK, retData)
}

func JsonIndustry(c echo.Context) error {
	var myData MyData
	err := c.Bind(myData)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, myData)
}

func main() {
	e := echo.New()
	e.GET("/", Index)
	e.GET("/tlht", TaiLieuHocTap)
	e.POST("/login", LoginHandler)
	e.POST("/hang", JsonEndpoint)
	e.POST("/ind", JsonIndustry)
	e.SetDebug(true)
	e.Run(standard.New(":1323"))
}
