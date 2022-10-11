package services

import (
	"encoding/json"
	"final/server/authentification"
	"final/server/models"
	"final/server/views"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (postgres *HandlersController) Post_Photos(ctx *gin.Context) {
	// Check Authorization
	tokenString := ctx.GetHeader("Authorization")
	jwtString := strings.Split(tokenString, "Bearer ")[1]
	if tokenString == "" {
		ctx.JSON(401, gin.H{"error": "request does not contain an access token"})
		ctx.Abort()
		return
	}
	err1 := authentification.ValidateToken(tokenString)
	if err1 != nil {
		ctx.JSON(401, gin.H{"error": err1.Error()})
		ctx.Abort()
		return
	}
	ctx.Next()

	// decode/Extract JWT
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Verivication"), nil
	})
	username := fmt.Sprintf("%v", claims["username"])
	email := fmt.Sprintf("%v", claims["email"])
	println(username)
	println(email)

	// Get Body Value
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	body_string := string(body)
	println(body_string)

	var key_data views.Request_Photos

	err := json.Unmarshal([]byte(body_string), &key_data)
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "FAILED_POST_FOTO",
			Error:   err.Error(),
		})
		return
	}
	println("Title: %s", key_data.Title)
	println("Caption: %s", key_data.Caption)
	println("Photo_Url: %s", key_data.Photo_Url)

	// Tittle and Photo_Url Validation
	if key_data.Title == "" || key_data.Photo_Url == "" {
		WriteJsonResponse(ctx, &views.Response{
			Message: "TITLE_OR_PHOTO_URL_IS_EMPTY",
			Status:  http.StatusInternalServerError,
			Payload: "photo_url or title field is empty!",
		})
		return
	}

	// query data from table photo
	var result models.User
	postgres.db.Table("users").Select("id").Where("email = ?", email).Scan(&result)
	println(result.ID)

	//generate photo_ID
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	err_photo := postgres.db.Create(&models.Photo{
		ID:        0,
		Title:     key_data.Title,
		Caption:   key_data.Caption,
		Photo_Url: key_data.Photo_Url,
		User_Id:   result.ID,
		Create_At: time.Time{},
		Update_At: time.Time{},
	}).Error
	if err_photo != nil {
		WriteJsonResponse(ctx, &views.Response{
			Message: "POST_PHOTO_FAILED",
			Status:  http.StatusInternalServerError,
			Payload: nil,
			Error:   err_photo.Error(),
		})
		return
	}

	WriteJsonResponse_PostPhoto(ctx, &views.Resp_Post_Photo{
		Message: "SUCCESS",
		Status:  http.StatusOK,
		Data: views.Data_Photo{
			ID:        r1.Int(),
			Title:     key_data.Title,
			Caption:   key_data.Caption,
			Photo_Url: key_data.Photo_Url,
			User_Id:   result.ID,
			Create_At: time.Time{},
		},
	})
}
