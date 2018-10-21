package models

import (
	"github.com/temog/album/backend/util"
	//"gopkg.in/mgo.v2"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

/*
# mongodb index
db.user.createIndex({account:1}, {background:true, unique:true})
db.user.createIndex({account:1, password:1}, {background:true})
*/

type Image struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	UserId      string        `json:"userId" bson:"userId"`
	Tag         []string      `json:"tag"`
	Memo        string        `json:"memo"`
	Name        string        `json:"name"`
	Secret      bool          `json:"secret"`
	Size        int           `json:"size"`
	ContentType string        `json:"contentType"`
	Url         string        `json:"url"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

// input parameter
type InputAddImage struct {
	Token  string  `json:"token" binding:"required,token"`
	Tag    string  `json:"tag" binding:"required"`
	Images []Image `json:"images" binding:"required"`
}

type InputUpdateImage struct {
	ID     string `json:"id" binding:"required"`
	Token  string `json:"token" binding:"required,token"`
	Tag    string `json:"tag" binding:"required"`
	Memo   string `json:"memo"`
	Secret bool   `json:"secret"`
}

type InputDeleteImage struct {
	ID     string `json:"id" binding:"required"`
	Token  string `json:"token" binding:"required,token"`
}

type InputGetImage struct {
	Token   string `form:"token" json:"token" binding:"required,token"`
	ImageId string `form:"imageId" json:"imageId" binding:"required"`
}

type InputGetTaggedImage struct {
	Token string `form:"token" json:"token" binding:"required,token"`
	Tag   string `json:"tag" binding:"required"`
	Page  int    `json:"page" binding:"required"`
	Limit int    `json:"limit" binding:"required"`
}

/*
  処理フロー
  1. タグのIDを取得
  2. userId を取得
  3. DBにデータ作る
  4. 画像を保存する
*/
func AddImage(input InputAddImage) (result bool, err error) {

	tagId := GetTagId(input.Tag)
	if len(tagId) == 0 {
		return result, err
	}

	user := FindUserByToken(input.Token)
	if user.Account == "" {
		return result, err
	}

	// 画像毎にループ
	userId := user.ID.Hex()
	errorString := ""
	for _, image := range input.Images {

		if checkDuplicate(image.Name, image.Size) == false {
			errorString += image.Name + " は登録済みです,"
			continue
		}

		// content-type 取得
		contentType := "image/png"
		if strings.Index(image.Url, "image/jpeg") != -1 {
			contentType = "image/jpeg"
		}

		// データ生成
		var imageId string
		imageId, err = CreateImage(userId, tagId, image.Memo, image.Name, image.Secret, image.Size, contentType)
		if err != nil {
			return result, err
		}

		// 画像保存
		if SaveImage(imageId, image.Url, contentType) {
			result = true
		}
	}

	if errorString != "" {
		err = errors.New(errorString)
	}

	return result, err
}

func UpdateImage(input InputUpdateImage) (result bool, err error) {

	tagId := GetTagId(input.Tag)
	if len(tagId) == 0 {
		return result, err
	}

	user := FindUserByToken(input.Token)
	if user.Account == "" {
		return result, err
	}

	data := Image{}
	col := GetCollection("image")
	objectId := bson.ObjectIdHex(input.ID)
	err = col.FindId(objectId).One(&data)
	if err != nil {
		return result, err
	}

	where := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"memo": input.Memo, "tag": tagId, "secret": input.Secret}}
	err = col.Update(where, update)
	if err != nil {
		return result, err
	}

	result = true
	return result, err
}

func DeleteImage(input InputDeleteImage) (result bool, err error) {

	user := FindUserByToken(input.Token)

	image := Image{}
	col := GetCollection("image")
	objectId := bson.ObjectIdHex(input.ID)
	err = col.FindId(objectId).One(&image)
	if err != nil {
		return result, err
	}

	// admin じゃない + userが一致しない
	isAdmin := IsAdmin(user)
	if isAdmin == false && user.ID.Hex() != image.UserId {
		err = errors.New("do not have permission")
		return result, err
	}

	info, err := col.RemoveAll(bson.M{"_id": objectId})
	if err != nil {
		return result, err
	}
	util.Dump(info)

	result = true
	return result, err
}

// 画像名、サイズがかぶってたら蹴る
func checkDuplicate(name string, size int) bool {

	data := Image{}
	col := GetCollection("image")
	col.Find(bson.M{"name": name, "size": size}).One(&data)
	util.Dump(data)

	// 存在する
	if data.Name != "" {
		return false
	}
	return true
}

func CreateImage(userId string, tag []string, memo string, name string, secret bool, size int, contentType string) (id string, err error) {

	collection := GetCollection("image")
	objectId := bson.NewObjectId()
	err = collection.Insert(&Image{
		objectId,
		userId,
		tag,
		memo,
		name,
		secret,
		size,
		contentType,
		"",
		time.Now(),
		time.Now(),
	})
	if err != nil {
		LogError(err)
		return id, err
	}

	return objectId.Hex(), err
}

func SaveImage(imageId string, url string, contentType string) bool {

	split := strings.Split(url, "data:"+contentType+";base64,")
	img, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		LogError(err)
		return false
	}

	file, err := os.Create("./data/image/" + imageId)
	defer file.Close()
	file.Write(img)

	return true
}

func CountImageByTagId(tagId string) int {

	col := GetCollection("image")
	query := ImageQuery(bson.M{"tag": tagId})
	num, err := col.Find(query).Count()
	if err != nil {
		LogError(err)
		return 0
	}

	return num
}

func GetImageById(id string) Image {

	data := Image{}
	col := GetCollection("image")
	err := col.FindId(bson.ObjectIdHex(id)).One(&data)
	if err != nil {
		LogError(err)
	}

	return data
}

func GetLatestImageByTagId(tagId string) Image {

	data := Image{}
	col := GetCollection("image")
	query := ImageQuery(bson.M{"tag": tagId})
	if IsModeSecret() == false {
		query["secret"] = false
	}

	err := col.Find(query).Sort("-updated_at").One(&data)
	if err != nil {
		LogError(err)
	}

	return data
}

func GetTaggedImage(tagId string, page int, limit int) []Image {
	page--
	images := []Image{}
	col := GetCollection("image")
	query := ImageQuery(bson.M{"tag": tagId})
	err := col.Find(query).Sort("-updated_at").Skip(page * limit).Limit(limit).All(&images)
	if err != nil {
		LogError(err)
	}
	return images
}
