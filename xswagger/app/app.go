package app

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

// App is app object
type App struct {
	config *Config
	e      *echo.Echo
}

// Run start server
func (app *App) Run(address string) {
	app.e.Logger.Fatal(app.e.Start(address))
}

// Destroy destroy this app.
func (app *App) Destroy() {
}

// CreateApp create a app object
func CreateApp(debug bool, config *Config) *App {
	if config.Static {
		if _, err := swag.ReadDoc(); err != nil {
			panic("no swagger registered, can't use static mode")
		}
	}

	app := &App{config, echo.New()}

	app.e.Debug = debug

	app.e.Use(middleware.Logger())
	app.e.Use(middleware.Recover())
	app.e.Use(middleware.CORS())

	// routers
	app.e.GET("/", redirectToSwaggerIndex)

	app.e.GET(config.DocPath, func(c echo.Context) error {
		f, err := os.Open(config.DocFile)
		if err != nil {
			return err
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		data := make(map[string]interface{}, 0)
		if err := dec.Decode(&data); err != nil {
			return err
		}
		data["host"] = config.Host
		return c.JSON(http.StatusOK, data)
	})

	wrapHandler := echoSwagger.EchoWrapHandler(echoSwagger.URL(config.DocPath))
	app.e.GET(config.SwaggerPath+"*", func(c echo.Context) error {
		if config.Static {
			return echoSwagger.WrapHandler(c)
		}
		return wrapHandler(c)
	})

	return app
}

func redirectToSwaggerIndex(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
