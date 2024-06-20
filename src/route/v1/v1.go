package v1

import (
	"cloud.google.com/go/firestore"

	"github.com/Artexus/api-widyabhuvana/src/middleware"
	answerRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/answer"
	categoryRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/category"
	subCategoryRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/subcategory"
	subTaskRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/subtask"
	taskRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/task"
	userRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	userActivityRepository "github.com/Artexus/api-widyabhuvana/src/repository/v1/useractivity"

	"github.com/Artexus/api-widyabhuvana/src/controller/v1/auth"
	"github.com/Artexus/api-widyabhuvana/src/controller/v1/category"
	"github.com/Artexus/api-widyabhuvana/src/controller/v1/subcategory"
	"github.com/Artexus/api-widyabhuvana/src/controller/v1/subtask"
	"github.com/Artexus/api-widyabhuvana/src/controller/v1/task"
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
	user.PATCH("", handler.Update)
}

func (r Route) routeCategory(v1 *gin.RouterGroup, handler *category.Controller) {
	category := v1.Group("categories", r.middleware.Auth)
	category.GET("", handler.Get)
	category.GET("progress", handler.GetUserProgress)
}

func (r Route) routeTask(v1 *gin.RouterGroup, handler *task.Controller) {
	category := v1.Group("tasks", r.middleware.Auth)
	category.GET("", handler.Get)
	category.POST("", handler.Submit)
}

func (r Route) routeSubTask(v1 *gin.RouterGroup, handler *subtask.Controller) {
	category := v1.Group("sub-tasks", r.middleware.Auth)
	category.GET("", handler.Get)
}

func (r Route) routeSubCategory(v1 *gin.RouterGroup, handler *subcategory.Controller) {
	category := v1.Group("sub-categories", r.middleware.Auth)
	category.GET("", handler.Get)
}

func New(middleware *middleware.Middleware) Route {
	return Route{
		middleware: middleware,
	}
}

func (r Route) InitRouter(router *gin.Engine, client *firestore.Client) {
	userRepo := userRepository.NewRepository(client)
	categoryRepo := categoryRepository.NewRepository(client)
	subCategoryRepo := subCategoryRepository.NewRepository(client)
	taskRepo := taskRepository.NewRepository(client)
	subTaskRepo := subTaskRepository.NewRepository(client)
	userActivityRepo := userActivityRepository.NewRepository(client)
	answerRepo := answerRepository.NewRepository(client)

	authCtrl := auth.NewController(userRepo)
	userCtrl := user.NewController(userRepo)
	categoryCtrl := category.NewController(categoryRepo, userActivityRepo, subCategoryRepo)
	subCategoryCtrl := subcategory.NewController(subCategoryRepo, userActivityRepo)
	taskCtrl := task.NewController(taskRepo, userRepo, answerRepo, subTaskRepo, userActivityRepo, categoryRepo, subCategoryRepo)
	subTaskCtrl := subtask.NewController(subTaskRepo, userRepo)

	v1 := router.Group("v1")
	r.routeAuth(v1, authCtrl)
	r.routeUser(v1, userCtrl)
	r.routeCategory(v1, categoryCtrl)
	r.routeSubCategory(v1, subCategoryCtrl)
	r.routeTask(v1, taskCtrl)
	r.routeSubTask(v1, subTaskCtrl)
}
