package main

import (
	"log"

	"github.com/boltdb/bolt"
	// db "github.com/hvnsweeting/meomeo/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"net/http"
)

var (
	db bolt.DB
	dbName = []byte("nambk")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// A Entry presents a blacklist entry.
type Entry struct {
	IP      string
	Comment string
}

type returnMessage struct {
	Message string
}

// AddIP adds and IP to blacklist
func AddIP(c echo.Context) error {
	var e Entry
	err := c.Bind(&e)
	// Input malformed
	if err != nil {
		return c.JSON(http.StatusBadRequest, returnMessage{Message: err.Error()})
	}

	// add entry to database
	if err = db.Update(func(tx *bolt.Tx) error {
		return b.Put([]byte(e.Comment), []byte(e.IP))
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, returnMessage{Message: err.Error()})
	}

	msg := returnMessage{"success"}

	log.Println(msg)

	return c.JSON(http.StatusCreated, msg)
}

// CheckBlacklistHandler checks if an IP in blacklist.
func CheckBlacklistHandler(c echo.Context) error {
	ipCheck := c.QueryParam("ip")
	log.Println("INPUT: ", ipCheck)
	if len(ipCheck) == 0 {
		return c.JSON(http.StatusBadRequest, returnMessage{Message: "failed"})
	}
	// add entry to database
	found := false
	// tt locks when call db.Update so multiple clients call check same time has to wait,
	// switch to db.View.
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(dbName).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == ipCheck {
				found = true
				break
			}
		}
		return nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, returnMessage{Message: err.Error()})	
	}

	var msg returnMessage
	if found {
		msg = returnMessage{"found"}
	} else {
		msg = returnMessage{"not found"}
	}

	log.Println(msg)

	return c.JSON(http.StatusCreated, returnMessage{Message: msg})

}

// Check and create db onetime only, when proram starts.
func initDB() error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(dbName)
		return err
	}
}

func main() {

	/* NOTE: This did not work because shorthand declaration is local to the closet syntatic block {..} 
	* so var `db` here is different from global `db` 
	//	db, err := bolt.Open("my.db", 0600, nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	/* FIX: */
	var err error
	if db, err = bolt.Open("my.db", 0600, nil); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.WithError(closeErr).Error("Cannot close db connection.")
		}
	}()
	log.Println(db.Info())

	if err = initDb(); err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.GET("/", Index)

	e.POST("/blacklist/ip", AddIP)
	e.Get("/blacklist/ip", CheckBlacklistHandler)

	e.GET("/tlht", TaiLieuHocTap)
	e.POST("/login", LoginHandler)
	e.POST("/hang", JSONEndpoint)
	e.POST("/ind", JSONIndustry)
	e.SetDebug(true)

	e.Run(standard.New(":1323"))
}
