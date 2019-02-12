package controllers

import (
	"encoding/hex"
	"errors"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/scrypt"
	"net/http"
	"sync"
	"time"
)


type UserID struct {
	value string
}

type User struct {
	ID UserID
	Password string `json:"password" form:"password" query:"password"`
	Name string `json:"name" form:"name" query:"name"`
}


type AllUsers struct {
	sm sync.Map
}

var allUsers AllUsers

func init() {
	allUsers = AllUsers{sync.Map{}}
}

func (users *AllUsers)findUser(id UserID) (bool, User, error) {
	user, exist := users.sm.Load(id.value)

	if !exist {
		return false, User{}, errors.New("UserNotFound")
	}

	return true, user.(User) , nil
}

func (users *AllUsers)store(user User) {
	users.sm.Store(user.ID, user)
}

func LoginViewController(context echo.Context) error {
	return context.Render(http.StatusOK, "login.html", make(map[string]interface{}))
}

func SignUpViewController(context echo.Context) error {
	return context.Render(http.StatusOK, "signup.html", make(map[string]interface{}))
}

type LoginBody struct {
	Password string `json:"password" form:"password" query:"password"`
	Name string `json:"name" form:"name" query:"name"`
}

func SignUpController(context echo.Context) error {
	loginBody := new(LoginBody)
	if err := context.Bind(loginBody); err != nil {
		return context.String(http.StatusBadRequest, err.Error())
	}


	hashedString := writeCookie(context, loginBody)
	user := User {
		ID: UserID{hashedString},
		Name: loginBody.Name,
		Password: loginBody.Password,
	}

	return context.JSON(http.StatusOK, user.Name)
}

func LoginController(context echo.Context) error {
	loginBody := new(LoginBody)
	if err := context.Bind(loginBody); err != nil {
		return context.String(http.StatusBadRequest, err.Error())
	}

	hashedString := writeCookie(context, loginBody)
	user := User {
		ID: UserID{hashedString},
		Name: loginBody.Name,
		Password: loginBody.Password,
	}


	return context.JSON(http.StatusOK, user.Name)
}

func writeCookie(c echo.Context, body *LoginBody) string {
	cookie := new(http.Cookie)
	cookie.Name = "user-token"
	hash := toHashFromScrypt(body.Name + body.Password)
	cookie.Value = hash
	cookie.Expires = time.Now().Add(time.Minute)
	c.SetCookie(cookie)
	return hash
}

func toHashFromScrypt(pass string) string {
	salt := []byte(viper.GetString("salt"))
	converted, _ := scrypt.Key([]byte(pass), salt, 16384, 8, 1, 32)
	return hex.EncodeToString(converted[:])
}