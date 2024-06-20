package subtask

import (
	"net/http"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	httpSubTask "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/subtask"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subtask"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Controller struct {
	repo *subtask.Repository
	user *user.Repository
}

func NewController(repo *subtask.Repository, user *user.Repository) *Controller {
	return &Controller{
		repo: repo,
		user: user,
	}
}

// Get godoc
// @Tags Sub Task
// @Summary Get Sub Tasks
// @Description Get sub tasks
// @Param id query string true "ID"
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {object} task.GetResponse
// @Router /v1/sub-tasks [GET]
func (ctrl Controller) Get(ctx *gin.Context) {
	req := httpSubTask.GetRequest{}
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

	resp := httpSubTask.GetResponse{}
	copier.Copy(&resp, task)

	qnas := task.Payload()
	resp.QnAs = make([]httpSubTask.QnAPayload, 0)
	for _, qna := range qnas {
		resp.QnAs = append(resp.QnAs, httpSubTask.QnAPayload{
			Question: qna.Question,
			Choices:  qna.Choices,
		})
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}
