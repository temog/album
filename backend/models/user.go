package models

import (
	"fmt"
	"github.com/temog/album/backend/util"
	"crypto/sha512"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
# mongodb index
db.user.createIndex({account:1}, {background:true, unique:true})
db.user.createIndex({account:1, password:1}, {background:true})
*/

type InputCreateUser struct {
	Account  string `form:"account" json:"account" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Role     string `form:"role" json:"role"`
	Token    string `json:"token"`
}

type InputEditUser struct {
	ID       string `form:"id" binding:"required"`
	Token    string `json:"token" binding:"required,token"`
	Account  string `form:"account" json:"account" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
	Password string `form:"password" json:"password"`
	Role     string `form:"role" json:"role"`
}

type InputSignIn struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InputSecret struct {
	Token    string `json:"token" binding:"required,token"`
	Password string `json:"password" binding:"required"`
}

type ResponseUserList struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type User struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Account   string        `json:"account"`
	Nickname  string        `json:"nickname"`
	Password  string        `json:"password"`
	Role      string        `json:"role"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func UserCount() (status bool, count int) {

	col := GetCollection("user")
	count, err := col.Count()
	if err != nil {
		LogError(err)
	}
	status = true
	return status, count
}

func SignIn(input InputSignIn) (status bool, userId string, nickname string, role string, token string) {

	status, userId, nickname, role = Auth(input.Account, input.Password)

	// 認証OKでtoken発行
	if status {
		token, _ = CreateToken(userId)
	}

	return status, userId, nickname, role, token
}

func SignOut(token string) bool {
	return removeToken(token)
}

// 認証 boolは初期falseである
func Auth(account string, password string) (status bool, userId string, nickname string, role string) {

	result := User{}
	col := GetCollection("user")
	err := col.Find(bson.M{"account": account, "password": util.Sha512String(password)}).One(&result)
	if err != nil {
		LogError(err)
	}

	// 見つかった
	if result.Account != "" {
		status = true
		userId = result.ID.Hex() // Hex() で string な _id がとれる
		nickname = result.Nickname
		role = result.Role
	}
	return status, userId, nickname, role
}

func AuthToken(token string) bool {
	// tokenあるか
	userId := FindToken(token)
	if userId == "" {
		return false
	}

	// ユーザがいるか
	user := FindUserById(userId)
	if user.Account == "" {
		return false
	}

	return true
}

// パスワードで再認証
func Secret(input InputSecret) bool {

	user := FindUserByToken(input.Token)
	if user.Account == "" {
		return false
	}

	if user.Password != util.Sha512String(input.Password) {
		return false
	}

	return setMode(input.Token, "admin")
}

func FindUserById(id string) User {
	data := User{}
	col := GetCollection("user")
	err := col.FindId(bson.ObjectIdHex(id)).One(&data)
	if err != nil {
		LogError(err)
	}
	return data
}

func FindUserByToken(token string) User {
	data := User{}
	userId := FindToken(token)
	if userId == "" {
		return data
	}

	col := GetCollection("user")
	err := col.FindId(bson.ObjectIdHex(userId)).One(&data)
	if err != nil {
		LogError(err)
	}
	return data
}

// 1ユーザだけ
func FindUser() User {
	result := User{}
	col := GetCollection("user")
	err := col.Find(bson.M{"name": "fuiwara"}).One(&result)
	fmt.Println("find user 01")
	if err != nil {
		LogError(err)
	}

	return result
}

// 全ユーザ
func FindAllUser() []User {
	result := []User{}
	col := GetCollection("user")
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}

	return result
}

func UserList() []*ResponseUserList {
	list := []*ResponseUserList{}
	result := []User{}
	col := GetCollection("user")
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		LogError(err)
	}

	for _, user := range result {
		data := &ResponseUserList{
			ID: user.ID,
			Account: user.Account,
			Nickname: user.Nickname,
			Role: user.Role,
		}
		list = append(list, data)
	}

	return list
}

// ユーザ作成
func CreateUser(input InputCreateUser) bool {

	account := input.Account
	nickname := input.Nickname
	password := input.Password
	role := input.Role
	if role == "" {
		role = "guest"
	}

	util.Dump("create user 1")

	col := GetCollection("user")
	count, err := col.Count()
	util.Dump(count)
	util.Dump(err)
	if err != nil {
		LogError(err)
		return false
	}
	util.Dump("create user 2")

	// 初回ユーザ以外は token 必須
	if count != 0 {
		token := input.Token
		user := FindUserByToken(token)
		if IsAdmin(user) == false {
			return false
		}
	} else {
		role = "admin"
	}
	util.Dump("create user 3")

	err = col.Insert(&User{
		bson.NewObjectId(),
		account,
		nickname,
		sha512String(password),
		role,
		time.Now(),
		time.Now(),
	})
	if err != nil {
		LogError(err)
		return false
	}

	util.Dump("create success")
	return true
}

func EditUser(input InputEditUser) bool {

	userId := input.ID
	account := input.Account
	nickname := input.Nickname
	password := input.Password
	role := input.Role

	// todo
	col := GetCollection("user")
	user := FindUserById(userId)
	if user.Account == "" {
		return false
	}

	objectId := bson.ObjectIdHex(userId)
	where := bson.M{"_id": objectId}
	update := bson.M{"account": account, "nickname": nickname, "updated_at": time.Now()}
	util.Dump(password)
	if password != "" {
		update["password"] = sha512String(password)
	}
	if role != "" {
		update["role"] = role
	}
	util.Dump(update)
	err := col.Update(where, bson.M{"$set": update})
	if err != nil {
		LogError(err)
		return false
	}

	return true
}

func IsAdmin(user User) bool {
	return user.Role == "admin"
}

func sha512String(str string) string {
	bytes := sha512.Sum512([]byte(str))
	return hex.EncodeToString(bytes[:])
}
