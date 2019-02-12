package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go_app/controllers"
	"io"
	"net/http"
	"html/template"
)

func initConfig() {
	configPath := viper.Get("config_path")
	if configPath == nil {
		configPath = "dev.toml"
	}
	viper.AutomaticEnv()
	viper.SetConfigFile(fmt.Sprintf("./configs/%s", configPath))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func init() {
	initConfig()
}


// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}


func main() {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	e.GET("/config", func(context echo.Context) error {
		return context.JSONPretty(http.StatusOK, viper.AllSettings(), "\t")
	})
	e.GET("/config/:keyName", func(context echo.Context) error {
		keyName := context.Param("keyName")
		return context.String(http.StatusOK, viper.GetString(keyName))
	})

	e.GET("/login", controllers.LoginViewController)
	e.POST("/login", controllers.LoginController)

	e.GET("/signup", controllers.SignUpViewController)
	e.POST("/signup", controllers.SignUpController)
	e.GET("/user", controllers.UserController)
	e.GET("/user/all", controllers.AllUserController)
	e.Static("/static", "public")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", viper.Get("port"))))
}