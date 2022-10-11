package services

import (
	"database/sql"
	"encoding/json"
	"final/server/authentification"
	"final/server/models"
	"final/server/views"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(userid int) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd34tg2y4j7j") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////// HANDLERS FOR USER ////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
func (postgres *HandlersController) Register_User(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	body_string := string(body)
	println(body_string)

	var key_data views.Request_Register

	err := json.Unmarshal([]byte(body_string), &key_data)
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "REGISTER_USER_FAILED",
			Error:   err.Error(),
		})
		return
	}

	println(key_data.Username)
	println(key_data.Email)
	println(key_data.Password)
	println(key_data.Age)

	valid(key_data.Email)

	// Email type Validation
	if valid(key_data.Email) != true {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "EMAIL_FORMAT_NOT_FALID",
		})
		return
	}

	// Check Email registered Database
	email_check := postgres.db.Where("email = ?", key_data.Email).First(&models.User{}).Error
	fmt.Println(email_check)
	if email_check == nil {
		WriteJsonResponse(ctx, &views.Response{
			Message: "EMAIL_ALREADY_REGISTERED",
			Status:  http.StatusInternalServerError,
			Payload: "email already use for another account!",
		})
		return
	}

	// Username Validation
	if key_data.Username == "" {
		WriteJsonResponse(ctx, &views.Response{
			Message: "USERNAME_IS_EMPTY",
			Status:  http.StatusInternalServerError,
			Payload: "username field is empty!",
		})
		return
	}
	username_check := postgres.db.Where("username = ?", key_data.Username).First(&models.User{}).Error
	fmt.Println(email_check)
	if username_check == nil {
		WriteJsonResponse(ctx, &views.Response{
			Message: "USER_NAME_ALREADY_TAKEN",
			Status:  http.StatusInternalServerError,
			Payload: "username already use for another account!",
		})
		return
	}

	// Password Validation :
	if key_data.Password == "" {
		WriteJsonResponse(ctx, &views.Response{
			Message: "PASSWORD_IS_EMPTY",
			Status:  http.StatusInternalServerError,
			Payload: "password field can't be empty!",
		})
		return
	}
	password_length := len(key_data.Password)
	if password_length < 6 {
		WriteJsonResponse(ctx, &views.Response{
			Message: "ERROR_PASSWORD_LENGTH",
			Status:  http.StatusInternalServerError,
			Payload: "password length must be more than 6 character",
		})
		return
	}
	//encript Password with bycript
	hash, _ := HashPassword(key_data.Password)

	// validation Age
	age := key_data.Age
	if age < 9 {
		WriteJsonResponse(ctx, &views.Response{
			Message: "AGE_IS_RESTRICTED",
			Status:  http.StatusInternalServerError,
			Payload: "age field can't be empty or less than 8 years old!",
		})
		return
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	err_create := postgres.db.Create(&models.User{
		ID:        r1.Int(),
		Username:  key_data.Username,
		Email:     key_data.Email,
		Password:  hash,
		Age:       key_data.Age,
		Create_At: time.Time{},
		Update_At: time.Time{},
	}).Error
	if err_create != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "REGISTER_USER_FAILED",
			Error:   err_create.Error(),
		})
		return
	}

	WriteJsonResponse_Succes(ctx, &views.Resp_Register_Success{
		Message: "SUCCESS",
		Status:  http.StatusCreated,
		Data: views.Data_Register{
			Age:      key_data.Age,
			Email:    key_data.Email,
			ID:       r1.Int(),
			Username: key_data.Username,
		},
	})
}

func (postgres *HandlersController) Login_User(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	body_string := string(body)
	println(body_string)

	var key_data views.Request_Login

	err := json.Unmarshal([]byte(body_string), &key_data)
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "LOGIN_USER_FAILED",
			Error:   err.Error(),
		})
		return
	}

	println(key_data.Email)
	println(key_data.Password)

	// Email type Validation
	if valid(key_data.Email) != true {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "EMAIL_FORMAT_NOT_FALID",
		})
		return
	}

	// Password and Email Login Check Acount
	var s sql.NullString

	postgres.db.Select("password").Where("email = ?", key_data.Email).First(&models.User{}).Scan(&s)
	password_from_db := s.String
	fmt.Printf("%s \n", password_from_db)

	match := CheckPasswordHash(key_data.Password, password_from_db)
	fmt.Println("Match:   ", match)
	if match != true {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusNotFound,
			Message: "EMAIL_AND_PASSWORD_NOT_MATCH",
			Payload: "password and email does not match!!",
		})
		return
	}

	// generate JWT-TOKEN
	var result models.User
	postgres.db.Table("users").Select("email", "username").Where("email = ?", key_data.Email).Scan(&result)
	println(result.Age)

	token, err := authentification.GenerateJWT(result.Email, result.Username)
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "FAILED_GENERATED_TOKEN",
		})
		return
	}

	WriteJsonResponse_Login(ctx, &views.Resp_Login{
		Message: "SUCCESS",
		Status:  http.StatusCreated,
		Data: views.Token{
			Token: token,
		},
	})
}

func (postgres *HandlersController) PUT_User(ctx *gin.Context) {
	// Check Authorization
	tokenString := ctx.GetHeader("Authorization")
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

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	body_string := string(body)
	println(body_string)

	var key_data views.Request_Login

	err := json.Unmarshal([]byte(body_string), &key_data)
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "LOGIN_USER_FAILED",
			Error:   err.Error(),
		})
		return
	}

	var result models.User
	postgres.db.Table("users").Select("age", "email", "id", "username").Where("email = ?", key_data.Email).Scan(&result)
	println(result.ID)

	WriteJsonResponse_Put(ctx, &views.Resp_Put{
		Message: "SUCCESS",
		Status:  http.StatusCreated,
		Data: views.Put{
			Age:       result.Age,
			Email:     result.Email,
			ID:        result.ID,
			Username:  result.Username,
			Update_At: time.Now(),
		},
	})
}

func (postgres *HandlersController) Delete_User(ctx *gin.Context) {
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

	err := postgres.db.Where("username = ?", username).Delete(&models.User{}).Error
	if err != nil {
		WriteJsonResponse(ctx, &views.Response{
			Status:  http.StatusInternalServerError,
			Message: "DELETE_USER_FAILED",
			Error:   err.Error(),
		})
		return
	}

	WriteJsonResponse_Delete(ctx, &views.Resp_Delete{
		Message: "SUCCESS",
		Status:  http.StatusOK,
		Data: views.Data_Delete{
			Message: "Your account has been successfully deleted",
		},
	})
}
