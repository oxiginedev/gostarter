package api

import (
	"log"
	"net/http"

	"github.com/oxiginedev/gostarter/api/handlers"
	"github.com/oxiginedev/gostarter/api/middleware"
	"github.com/oxiginedev/gostarter/config"
	"github.com/oxiginedev/gostarter/internal/database"
	"github.com/oxiginedev/gostarter/internal/pkg/jwt"
	"github.com/oxiginedev/gostarter/services"
	"github.com/oxiginedev/gostarter/util"

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

	h := &handlers.Handler{
		DB:     a.opts.DB,
		Config: a.opts.Config,
		JWT:    jwt.NewJWT(&a.opts.Config.JWT),
	}

	v1Router := router.Group("/api/v1")

	v1Router.GET("/health", h.HandleHealth)

	v1Router.POST("/auth/register", h.RegisterUser)
	v1Router.POST("/auth/login", h.LoginUser)

	authV1Router := v1Router.Group("", middleware.Authenticate(h.JWT, a.opts.DB))

	authV1Router.GET("/auth/verification/:token", h.VerifyUserEmail)
	authV1Router.POST("/auth/verification/resend", h.ResendVerificationEmail)
	authV1Router.GET("/user/me", h.GetUser)

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
