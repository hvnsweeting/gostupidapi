package main

import (
	"log"

	"github.com/boltdb/bolt"
	// db "github.com/hvnsweeting/meomeo/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"net/http"
)

var db bolt.DB

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Entry struct {
	IP      string
	Comment string
}

type ReturnMessage struct {
	Message string
}

func AddIP(c echo.Context) error {
	var e Entry
	err := c.Bind(&e)
	if err != nil {
		log.Println("error", err)
	}

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// add entry to database
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nambk"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(e.Comment), []byte(e.IP))
		return err
	})

	if err != nil {
		log.Print("DB update failed", err)
	}

	msg := ReturnMessage{"success"}

	log.Println(msg)

	return c.JSON(http.StatusCreated, msg)
}

func CheckBlacklistHandler(c echo.Context) error {
	ipCheck := c.QueryParam("ip")
	log.Println("INPUT: ", ipCheck)
	if len(ipCheck) == 0 {
		c.JSON(http.StatusBadRequest, ReturnMessage{"failed"})
	}

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// add entry to database
	found := false
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nambk"))
		if err != nil {
			return err
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == ipCheck {
				found = true
			}
		}

		return err
	})

	if err != nil {
		log.Print("DB update failed", err)
	}

	var msg ReturnMessage
	if found {
		msg = ReturnMessage{"found"}
	} else {
		msg = ReturnMessage{"not found"}
	}

	log.Println(msg)

	return c.JSON(http.StatusCreated, msg)

}

func main() {

	//	db, err := bolt.Open("my.db", 0600, nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer db.Close()
	//	log.Println(db.Info())

	e := echo.New()
	e.GET("/", Index)

	e.POST("/blacklist/ip", AddIP)
	e.Get("/blacklist/ip", CheckBlacklistHandler)

	e.GET("/tlht", TaiLieuHocTap)
	e.POST("/login", LoginHandler)
	e.POST("/hang", JsonEndpoint)
	e.POST("/ind", JsonIndustry)
	e.SetDebug(true)

	e.Run(standard.New(":1323"))
}
