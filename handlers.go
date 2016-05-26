package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	// log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Index handles XYZ
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// TaiLieuHocTap handles XYZ
func TaiLieuHocTap(c echo.Context) error {
	return c.String(http.StatusForbidden, "Forbidden")
}

// LoginHandler handles login process
func LoginHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, c.QueryParams())
}

type myData struct {
	Password []string
	User     []string `json:"user_name"`
}

type returnData struct {
	Password []string
	User     []string
}

// JSONEndpoint handles ...
func JSONEndpoint(c echo.Context) error {
	r := c.Request()

	data, err := ioutil.ReadAll(r.Body())
	log.Println(string(data))
	var mData myData
	if err == nil && data != nil {
		err = json.Unmarshal(data, &mData)
		if err != nil {
			log.Println("Marshalling: ", err)
			return c.JSON(http.StatusBadRequest, "Bad JSON")
		}
	}

	retData := returnData{User: mData.User, Password: mData.Password}
	return c.JSON(http.StatusOK, retData)
}

// JSONIndustry handles ...
func JSONIndustry(c echo.Context) error {
	var mData myData
	err := c.Bind(mData)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, mData)
}
