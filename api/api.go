package api

import (
	"github/oxiginedev/gostarter/api/handlers"
	"github/oxiginedev/gostarter/api/middleware"
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/services"
	"github/oxiginedev/gostarter/util"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type API struct {
	DB     database.Database
	Config *config.Configuration
}

func NewAPI(db database.Database, cfg *config.Configuration) *API {
	ah := &API{
		DB:     db,
		Config: cfg,
	}

	return ah
}

func (a *API) buildRouter() *echo.Echo {
	router := echo.New()

	router.Use(echoMiddleware.RequestID())
	router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{}))
	router.Use(echoMiddleware.Recover())
	router.Use(middleware.LogHTTPRequest())

	router.HTTPErrorHandler = HTTPErrorHandler

	return router
}

func (a *API) BuildAPIRoutes() *echo.Echo {
	router := a.buildRouter()
	handler := &handlers.Handler{DB: a.DB, Config: a.Config}

	v1Router := router.Group("/api/v1")

	v1Router.GET("/health", handler.HandleHealth)

	v1Router.POST("/auth/register", handler.RegisterUser)
	v1Router.POST("/auth/login", handler.LoginUser)

	router.RouteNotFound("/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, util.BuildErrorResponse("The requested URL is invalid", nil))
	})

	return router
}

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := http.StatusText(http.StatusInternalServerError)

	se, ok := err.(*services.ServiceError)
	if ok {
		if se.Internal != nil {
			if serr, ok := se.Internal.(*services.ServiceError); ok {
				se = serr
			}
		}

		code = se.ErrCode
		message = se.ErrMsg
	}

	err = c.JSON(code, util.BuildErrorResponse(message, nil))

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(se.ErrCode)
	}

	if err != nil {
		log.Fatal(err)
	}
}
