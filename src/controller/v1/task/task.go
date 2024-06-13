package task

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	dbUserActivity "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/useractivity"
	httpTask "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/task"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/category"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subcategory"

	"github.com/Artexus/api-widyabhuvana/src/repository/v1/task"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/useractivity"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Controller struct {
	repo         *task.Repository
	user         *user.Repository
	useractivity *useractivity.Repository
	category     *category.Repository
	subcategory  *subcategory.Repository
}

func NewController(repo *task.Repository, user *user.Repository, useractivity *useractivity.Repository, category *category.Repository, subcategory *subcategory.Repository) *Controller {
	return &Controller{
		repo:         repo,
		user:         user,
		useractivity: useractivity,
		category:     category,
		subcategory:  subcategory,
	}
}

// Get godoc
// @Tags User
// @Summary Get Tasks
// @Description Get tasks
// @Param id query string true "ID"
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {object} task.GetResponse
// @Router /v1/tasks [GET]
func (ctrl Controller) Get(ctx *gin.Context) {
	req := httpTask.GetRequest{}
	err := ctx.BindQuery(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	req.ID, err = aes.DecryptID(req.EncID)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"id": constant.ErrInvalid.Error(),
		})
		return
	}

	task, err := ctrl.repo.Get(ctx, req.ID)
	if err != nil {
		constant.Error.Println("get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	resp := httpTask.GetResponse{}
	copier.Copy(&resp, task)
	if task.Type == constant.MultipleChoice {
		p := task.Payload()
		copier.Copy(&resp.QnAs, p)
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}

// Submit godoc
// @Tags User
// @Summary Submit Tasks
// @Description Submit tasks
// @Description Please send the answer accordingly
// @Param Authorization header string true "Bearer Token"
// @Param body body task.SubmitRequest true "Body"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {string} string "OK"
// @Router /v1/tasks [POST]
func (ctrl Controller) Submit(ctx *gin.Context) {
	req := httpTask.SubmitRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	req.ID, err = aes.DecryptID(req.EncID)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"id": constant.ErrInvalid.Error(),
		})
		return
	}

	req.UserID, _ = jwt.ExtractIDToken(ctx.GetHeader("Authorization"))
	task, err := ctrl.repo.Get(ctx, req.ID)
	if err != nil {
		constant.Error.Println("task: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	docID, userActivity, err := ctrl.useractivity.Get(ctx, req.UserID, task.CategoryID)
	if err != nil && !errors.Is(err, constant.ErrNotFound) {
		constant.Error.Println("useractivity: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	subCategory, err := ctrl.subcategory.Take(ctx, task.SubCategoryID)
	if err != nil {
		constant.Error.Println("sub category: take ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	if userActivity.RemainingTask == 0 {
		userActivity.RemainingTask = len(subCategory.Tasks)
	}

	isFound := false
	for _, t := range subCategory.Tasks[:len(subCategory.Tasks)-userActivity.RemainingTask] {
		if t == task.ID {
			isFound = true
			break
		}
	}

	if isFound {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"task": constant.ErrInvalid.Error(),
		})
		return
	}

	if !errors.Is(err, constant.ErrNotFound) {
		if userActivity.Status == dbUserActivity.Completed {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"task": constant.ErrCompleted.Error(),
			})
			return
		}

		if userActivity.LastTaskID == task.ID {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"task": constant.ErrInvalid.Error(),
			})
			return
		}
	}

	minPoint := 0
	if task.Type == constant.MultipleChoice {
		if len(req.Answer) != len(task.QnAs) {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"answers": constant.ErrInvalid.Error(),
			})
			return
		}

		payloads := task.Payload()
		for i, payload := range payloads {
			if !strings.EqualFold(req.Answer[i], payload.Answer) {
				minPoint += (task.Point / len(payloads))
			}
		}
	}

	totalPoint := task.Point - minPoint
	err = ctrl.user.UpdatePoint(ctx, req.UserID, totalPoint)
	if err != nil {
		constant.Error.Println("update point ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	// record user activities
	userActivity = dbUserActivity.UserActivity{
		UserID:            req.UserID,
		CategoryID:        task.CategoryID,
		LastSubCategoryID: task.SubCategoryID,
		LastTaskID:        task.ID,
		Status:            dbUserActivity.NotYet,
	}

	category, err := ctrl.category.Take(ctx, task.CategoryID)
	if err != nil {
		constant.Error.Println("category: take ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	i := 1
	for _, s := range category.SubCategories {
		if s == task.SubCategoryID {
			break
		}

		i++
	}

	j := 1
	for _, s := range subCategory.Tasks {
		if s == task.ID {
			break
		}

		j++
	}

	userActivity.RemainingSubCategory = len(category.SubCategories) - i
	userActivity.RemainingTask = len(subCategory.Tasks) - j
	if userActivity.RemainingSubCategory == 0 && userActivity.RemainingTask == 0 {
		userActivity.Status = dbUserActivity.Completed
	}

	err = ctrl.useractivity.Set(ctx, userActivity, docID)
	if err != nil {
		constant.Error.Println("user activity: set ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, nil)
}
