package v1

import (
	"cloud.google.com/go/firestore"

	"github.com/Artexus/api-widyabhuvana/src/middleware"
	userRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/user"

	"github.com/Artexus/api-widyabhuvana/src/controller/v1/auth"
	"github.com/Artexus/api-widyabhuvana/src/controller/v1/user"

	"github.com/gin-gonic/gin"
)

type Route struct {
	middleware *middleware.Middleware
}

func (r Route) routeAuth(v1 *gin.RouterGroup, handler *auth.Controller) {
	v1.POST("login", handler.Login)
	v1.POST("register", handler.Register)
}

func (r Route) routeUser(v1 *gin.RouterGroup, handler *user.Controller) {
	user := v1.Group("users", r.middleware.Auth)
	user.GET("", handler.Get)
}

func New(middleware *middleware.Middleware) Route {
	return Route{
		middleware: middleware,
	}
}

func (r Route) InitRouter(router *gin.Engine, client *firestore.Client) {
	userRepo := userRepository.NewRepository(client)

	authCtrl := auth.NewController(userRepo)
	userCtrl := user.NewController(userRepo)

	v1 := router.Group("v1")
	r.routeAuth(v1, authCtrl)
	r.routeUser(v1, userCtrl)
}
