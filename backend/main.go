package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/temog/album/backend/models"
	"github.com/temog/album/backend/util"
	"gopkg.in/go-playground/validator.v8"
)

// gin の validation用 簡単なのはこれでいいが・・・
type InputToken struct {
	Token string `json:"token" binding:"required,token"`
}

const Base = "/"

func main() {

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("token", token)
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://local.album.temo.xyz:8080", "https://album.temo.dev"}
	router.Use(cors.New(config))

	// grouping routes : image
	img := router.Group(Base + "api/image/")
	{
		// api/image/add
		img.OPTIONS("add", ResponseOption)
		img.POST("add", func(c *gin.Context) {

			input := models.InputAddImage{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status, err := models.AddImage(input)
			respError := ""
			if err != nil {
				respError = err.Error()
			}

			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"error":  respError,
				"input":  input,
			})
		})

		// api/image/update
		img.OPTIONS("update", ResponseOption)
		img.POST("update", func(c *gin.Context) {

			input := models.InputUpdateImage{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status, err := models.UpdateImage(input)
			respError := ""
			if err != nil {
				respError = err.Error()
			}

			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"error":  respError,
				"input":  input,
			})
		})

		// api/image/delete
		img.OPTIONS("delete", ResponseOption)
		img.POST("delete", func(c *gin.Context) {

			input := models.InputDeleteImage{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status, err := models.DeleteImage(input)
			respError := ""
			if err != nil {
				respError = err.Error()
			}

			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"error":  respError,
				"input":  input,
			})
		})

		// api/image/get
		img.OPTIONS("get", ResponseOption)
		img.GET("get", func(c *gin.Context) {

			input := models.InputGetImage{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			dataImg := models.GetImageById(input.ImageId)
			if dataImg.Name == "" {
				c.String(http.StatusForbidden, "forbidden")
			}

			// 以下画像パスから画像をデコードしてbufferに突っ込んで出力
			dir, _ := os.Getwd()
			imgPath := dir + "/data/image/" + input.ImageId

			file, err := os.Open(imgPath)
			defer file.Close()
			if err != nil {
				c.String(http.StatusForbidden, "forbidden")
			}
			util.Dump(file)

			img, _, imgErr := image.Decode(file)
			if imgErr != nil {
				c.String(http.StatusForbidden, "forbidden")
			}

			buffer := new(bytes.Buffer)

			contentType := dataImg.ContentType
			if contentType == "image/png" {
				err = png.Encode(buffer, img)
			} else if contentType == "image/jpeg" {
				err = jpeg.Encode(buffer, img, nil)
			}

			if err != nil {
				c.String(http.StatusForbidden, "forbidden")
			}

			c.Writer.Header().Set("Content-Type", contentType)
			c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
			c.Writer.Write(buffer.Bytes())
		})

		img.OPTIONS("getTagged", ResponseOption)
		img.POST("getTagged", func(c *gin.Context) {

			input := models.InputGetTaggedImage{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			images := models.GetTaggedImage(input.Tag, input.Page, input.Limit)
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"images": images,
			})
		})
	}

	// grouping routes : tag
	tag := router.Group(Base + "api/tag/")
	{
		tag.OPTIONS("getIndex", ResponseOption)
		tag.POST("getIndex", func(c *gin.Context) {

			input := InputToken{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			tags := models.GetIndex()
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"tags":   tags,
			})
		})

		tag.OPTIONS("getAll", ResponseOption)
		tag.POST("getAll", func(c *gin.Context) {
			input := InputToken{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			tags := models.GetAllTag()
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"tags":   tags,
			})
		})
	}

	// grouping routes : user
	user := router.Group(Base + "api/user/")
	{
		// count
		user.OPTIONS("count", ResponseOption)
		user.GET("count", func(c *gin.Context) {

			status, count := models.UserCount()
			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"count":  count,
			})
		})

		// signIn
		user.OPTIONS("signIn", ResponseOption)
		user.POST("signIn", func(c *gin.Context) {
			input := models.InputSignIn{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status, userId, nickname, role, token := models.SignIn(input)
			c.JSON(http.StatusOK, gin.H{
				"status":   status,
				"userId":   userId,
				"nickname": nickname,
				"role":     role,
				"token":    token,
			})
		})

		// signOut
		user.OPTIONS("signOut", ResponseOption)
		user.POST("signOut", func(c *gin.Context) {

			input := InputToken{}
			status := models.SignOut(input.Token)
			c.JSON(http.StatusOK, gin.H{
				"status": status,
			})
		})

		user.OPTIONS("secret", ResponseOption)
		user.POST("secret", func(c *gin.Context) {
			input := models.InputSecret{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status := models.Secret(input)
			c.JSON(http.StatusOK, gin.H{
				"status": status,
			})
		})

		user.OPTIONS("create", ResponseOption)
		user.POST("create", func(c *gin.Context) {
			// gin の validation GET でも POSTでも動く
			input := models.InputCreateUser{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status := models.CreateUser(input)
			c.JSON(http.StatusOK, gin.H{"status": status})
		})

		user.OPTIONS("edit", ResponseOption)
		user.POST("edit", func(c *gin.Context) {
			input := models.InputEditUser{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			status := models.EditUser(input)
			c.JSON(http.StatusOK, gin.H{"status": status})
		})

		user.OPTIONS("list", ResponseOption)
		user.POST("list", func(c *gin.Context) {
			input := InputToken{}
			if err := c.ShouldBind(&input); err != nil {
				ResponseValidationError(c, err)
				return
			}

			userList := models.UserList()
			c.JSON(http.StatusOK, gin.H{"status": true, "userList": userList})
		})

		user.GET("get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": models.FindUser(),
			})
		})

		user.GET("getAll", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": models.FindAllUser(),
			})
		})
	}

	// 指定なしだと 8080
	router.Run()

	// port指定もできるぽい
	// router.Run(":3000")
}

func ResponseOption(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func ResponseValidationError(c *gin.Context, err error) {
	errStr := err.Error()
	util.Dump(errStr)
	resp := gin.H{}
	if strings.Index(errStr, "'Token'") != -1 {
		resp = gin.H{"status": "authError"}
	} else {
		resp = gin.H{
			"status":  false,
			"message": errStr,
		}
	}
	c.JSON(http.StatusBadRequest, resp)
}

// custom validation
func token(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	token := field.String()
	return models.AuthToken(token)
}
