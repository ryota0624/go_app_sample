package controllers

import (
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/xerrors"
	"net/http"
	"sync"
	"time"
)


type UserID struct {
	value string
}

func (u UserID) format() string {
	return u.value
}

func (u UserID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.format() + `"`), nil
}

type User struct {
	ID UserID `json:"id" form:"id" query:"id"`
	Password string `json:"-" form:"-" query:"-"`
	Name string `json:"name" form:"name" query:"name"`
}


type AllUsers struct {
	sm sync.Map
}

type Sessions struct {
	sm sync.Map
}

var sessions Sessions

var allUsers AllUsers

func init() {
	allUsers = AllUsers{sync.Map{}}
	sessions = Sessions{sync.Map{}}
}

func (users *AllUsers)findUser(id UserID) (bool, User, error) {
	user, ok := users.sm.Load(id.value)

	if !ok {
		return false, User{}, xerrors.New("UserNotFound in Map")
	}

	return true, user.(User) , nil
}

func (users *AllUsers)store(user User) {
	users.sm.Store(user.ID.value, user)
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
	allUsers.store(user)

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

	allUsers.store(user)

	return context.JSON(http.StatusOK, user.Name)
}

func UserController(context echo.Context) error {
	userTokenCookie, err := context.Cookie(userToken)
	if err != nil {
		return err
	}

	exist, user, err := allUsers.findUser(UserID{userTokenCookie.Value})

	if err != nil {
		return err
	}

	if !exist {
		return context.String(http.StatusUnauthorized, "UserNotFound")
	}

	return context.JSON(http.StatusOK, map[string]string{"name": user.Name})
}

func AllUserController(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Response().WriteHeader(http.StatusOK)
	_, _ = ctx.Response().Write([]byte("[\n"))
	ctx.Response().Flush()

	allUsers.sm.Range(func(key, value interface{}) bool {
		user := value.(User)
		if err := json.NewEncoder(ctx.Response()).Encode(user); err != nil {
			return false
		}
		_, _ = ctx.Response().Write([]byte(",\n"))

		ctx.Response().Flush()
		return true
	})
	_, _ = ctx.Response().Write([]byte("]"))
	ctx.Response().Flush()


	return nil
}

const userToken = "user-token"

func writeCookie(c echo.Context, body *LoginBody) string {
	cookie := new(http.Cookie)
	cookie.Name = userToken
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