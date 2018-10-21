package models

import(
  "fmt"
  "sync"
  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/temog/album/backend/util"
)

type singleton struct {
  Database *mgo.Database
  Session *mgo.Session
}

var instance *singleton
var once sync.Once

// public method
func GetInstance() *singleton {

  // 必ず1回しか実行されない
  once.Do(func() {

    fmt.Println("============== db connection !!! ==============")
    session, err := mgo.Dial("localhost")
    if err != nil {
      panic(err)
    }

    // read secondaryだけど、書き込み発生後は、read write primary に
    session.SetMode(mgo.Monotonic, true)

    instance = &singleton{session.DB("album"), session}
  })
  return instance
}

func GetCollection(collection string) *mgo.Collection {
  db := GetInstance()
  return db.Database.C(collection)
}

func ImageQuery(q bson.M) bson.M {
  if currentMode != "admin" {
    q["secret"] = false
  }

  return q
}

// 書いてはみたが、ginが動いているならずっとsession維持で良い気がするんだが・・
func CloseSession() {
  db := GetInstance()
  db.Session.Close()
  util.Dump("session closed")
  fmt.Println("session closed")
}
