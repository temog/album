package models

import (
	"github.com/temog/album/backend/util"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

/*
# mongodb index
db.tag.createIndex({name:1}, {background: true, unique: true})
*/

type Tag struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type ResponseGetIndex struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Name         string        `json:"name"`
	LastModified string        `json:"lastModified"`
	NumImage     int           `json:"numImage"`  // 画像数
	Thumbnail    string        `json:"thumbnail"` // サムネイル画像
}

/*
トップページに表示するタグ一覧を返す
*/
func GetIndex() []*ResponseGetIndex {

	data := []*ResponseGetIndex{}

	tags := GetAllTag()
	for _, tag := range tags {
		tagId := tag.ID.Hex()
		img := GetLatestImageByTagId(tagId)
		if img.ID.Hex() == "" {
			continue
		}
		index := &ResponseGetIndex{
			ID:           tag.ID,
			Name:         tag.Name,
			LastModified: img.UpdatedAt.Format("2006-01-02 15:04:05"),
			NumImage:     CountImageByTagId(tagId),
			Thumbnail:    img.ID.Hex(),
		}
		data = append(data, index)
	}

	util.Dump(data)
	return data
}

func GetTagId(name string) []string {

	var tagIds []string
	tags := strings.Split(name, " ")
	for _, tagName := range tags {
		// タグがあるか
		tag := FindTag(tagName)

		// ないので作成
		var id string
		var err error
		if tag.Name == "" {
			id, err = CreateTag(tagName)
			if err != nil {
				LogError(err)
			}
		} else {
			id = tag.ID.Hex()
		}
		tagIds = append(tagIds, id)
	}

	return tagIds
}

func GetAllTag() []Tag {

	tags := []Tag{}
	col := GetCollection("tag")
	err := col.Find(bson.M{}).Sort("-updated_at").All(&tags)
	if err != nil {
		LogError(err)
	}
	return tags
}

func FindTag(name string) Tag {
	tag := Tag{}
	col := GetCollection("tag")
	err := col.Find(bson.M{"name": name}).One(&tag)
	if err != nil {
		LogError(err)
	}
	return tag
}

func CreateTag(name string) (id string, err error) {

	collection := GetCollection("tag")
	objectId := bson.NewObjectId()
	err = collection.Insert(&Tag{
		objectId,
		name,
		time.Now(),
		time.Now(),
	})
	if err != nil {
		LogError(err)
		return id, err
	}

	id = objectId.Hex()
	return id, err
}
