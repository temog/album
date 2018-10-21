package models

import(
  "time"
  // "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/temog/album/backend/util"
)

/*
# mongodb index
db.token.createIndex({token:1, userId:1}, {background:true})
db.token.createIndex({created_at: 1}, {expireAfterSeconds: 7200})
*/

//
type Token struct {
  ID       bson.ObjectId  `json:"_id" bson:"_id"`
  UserId   string         `json:"userId" bson:"userId"`
  Token    string         `json:"token"`
  Mode     string         `json:"mode"`
  CreatedAt time.Time     `json:"created_at" bson:"created_at"`
  UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func CreateToken(userId string) (token string, err error) {

  token = util.Sha512String(userId + util.CastString(util.Random(99999)))

  collection := GetCollection("token")
  err = collection.Insert(&Token{
    bson.NewObjectId(),
    userId,
    token,
    "default",
    time.Now(),
    time.Now(),
  })
  if err != nil {
    LogError(err)
    return token, err
  }

  return token, err
}

func setMode(token string, mode string) bool {

  where := bson.M{"token": token}
  data := bson.M{"$set": bson.M{"mode": mode}}

  collection := GetCollection("token")
  err := collection.Update(where, data)
  if err != nil {
    LogError(err)
    return false
  }
  return true
}

func removeToken(token string) bool {

  where := bson.M{"token": token}
  collection := GetCollection("token")
  _, err := collection.RemoveAll(where)
  if err != nil {
    LogError(err)
    return false
  }
  return true
}

var currentMode string
func FindToken(token string) (userId string) {

  data := Token{}
  col := GetCollection("token")
  err := col.Find(bson.M{"token":token}).One(&data)
  if err != nil {
    LogError(err)
  } else {
    userId = data.UserId
    currentMode = data.Mode
  }
  return userId
}

func IsModeSecret() bool {
  return currentMode == "admin"
}
