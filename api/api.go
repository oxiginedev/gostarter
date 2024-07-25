package api

import (
	"github/oxiginedev/gostarter/api/handlers"
	"github/oxiginedev/gostarter/api/middleware"
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
	"github/oxiginedev/gostarter/services"
	"github/oxiginedev/gostarter/util"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type APIOptions struct {
	Config *config.Configuration
	DB     *database.Connection
}

type API struct {
	Router http.Handler
	opts   *APIOptions
}

func NewAPI(opts *APIOptions) *API {
	ah := &API{
		opts: opts,
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

	handler := &handlers.Handler{
		DB:     a.opts.DB,
		Config: a.opts.Config,
		JWT:    jwt.NewJWT(&a.opts.Config.JWT),
	}

	v1Router := router.Group("/api/v1")

	v1Router.GET("/health", handler.HandleHealth)

	v1Router.POST("/auth/register", handler.RegisterUser)
	v1Router.POST("/auth/login", handler.LoginUser)

	authV1Router := v1Router.Group("", middleware.Authenticate(handler.JWT, a.opts.DB))

	authV1Router.GET("/user/me", handler.GetUser)

	router.RouteNotFound("/*", func(c echo.Context) error {
		err := util.BuildErrorResponse("The requested URL is invalid", nil)
		return c.JSON(http.StatusNotFound, err)
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
