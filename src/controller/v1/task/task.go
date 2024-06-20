package task

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	dbAnswer "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/answer"
	dbSubTask "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/subtask"
	dbUserActivity "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/useractivity"
	httpTask "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/task"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/answer"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/category"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subcategory"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subtask"

	"github.com/Artexus/api-widyabhuvana/src/repository/v1/task"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/useractivity"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/array"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Controller struct {
	repo         *task.Repository
	user         *user.Repository
	answer       *answer.Repository
	subtask      *subtask.Repository
	useractivity *useractivity.Repository
	category     *category.Repository
	subcategory  *subcategory.Repository
}

func NewController(repo *task.Repository, user *user.Repository, answer *answer.Repository, subtask *subtask.Repository, useractivity *useractivity.Repository, category *category.Repository, subcategory *subcategory.Repository) *Controller {
	return &Controller{
		repo:         repo,
		answer:       answer,
		user:         user,
		subtask:      subtask,
		useractivity: useractivity,
		category:     category,
		subcategory:  subcategory,
	}
}

// Get godoc
// @Tags Task
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

	p := task.Payload()
	if task.Type == constant.MultipleChoice || task.Type == constant.Essay {
		copier.Copy(&resp.QnAs, p)
	} else if task.Type == constant.Matching {
		answer := []string{}
		questions := []string{}
		for _, t := range p {
			questions = append(questions, t.Question)
			answer = append(answer, t.Answer)
		}

		for i := 0; i < 10; i++ {
			rnd := rand.New(rand.NewSource(time.Now().Unix()))

			x := rnd.Intn(len(answer))
			y := rnd.Intn(len(answer))

			answer[x], answer[y] = answer[y], answer[x]
		}

		resp.Matches = httpTask.Matches{
			Questions: questions,
			Choices:   answer,
		}

	} else if task.Type == constant.Detective {
		resp.Detective.SubTasks = aes.EncryptIDs(task.SubTasks)
	} else if task.Type == constant.Level {
		p := task.LevelPayload()
		copier.Copy(&resp.Levels, p)
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}

// Submit godoc
// @Tags Task
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

	lp := task.LevelPayload()
	if task.Type == constant.Detective {
		req.SubTaskID, err = aes.DecryptID(req.EncSubTaskID)
		if err != nil {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"sub_task_id": constant.ErrInvalid.Error(),
			})
			return
		}
	} else if task.Type == constant.Level {
		if len(lp.Level1.QnA)+len(lp.Level2.QnA) != len(req.Answer) {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"answer": constant.ErrInvalid.Error(),
			})
			return
		}
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

	var resp *httpTask.SubmitResponse
	minPoint := 0
	if task.Type == constant.MultipleChoice || task.Type == constant.Matching {
		if len(req.Answer) != len(task.QnAs) {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"answers": constant.ErrInvalid.Error(),
			})
			return
		}

		payloads := task.Payload()
		for i, payload := range payloads {
			if !strings.EqualFold(req.Answer[i][0], payload.Answer) {
				minPoint += (task.Point / len(payloads))
			}
		}
	} else if task.Type == constant.Essay {
		if len(req.Answer) <= 0 {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"answers": constant.ErrInvalid.Error(),
			})
			return
		}

		id, err := ctrl.answer.Create(ctx, dbAnswer.Answer{
			TaskID: task.ID,
			UserID: req.UserID,
			Answer: req.Answer[0][0],
		})
		if err != nil {
			constant.Error.Println("answer: create ", err)
			rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
			return
		}

		resp = &httpTask.SubmitResponse{
			EncID: aes.EncryptID(id),
		}
	} else if task.Type == constant.Detective {
		var subTask dbSubTask.SubTask
		subTask, err = ctrl.subtask.Get(ctx, req.SubTaskID)
		if err != nil {
			constant.Error.Println("subtask: get ", err)
			rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
			return
		}

		counter := 0
		qnas := subTask.Payload()
		for _, qna := range qnas {
			isFound := false
			for _, answer := range qna.Answer {
				if array.In(qna.Answer, answer) {
					isFound = true
				}
			}

			if !isFound {
				counter++
			}
		}

		minPoint = task.Point/len(qnas) - counter
	} else if task.Type == constant.Level {
		counter := 0
		lp := task.LevelPayload()
		for i, lv := range lp.Level1.QnA {
			isFound := false
			for _, answer := range req.Answer[i] {
				if array.In(lv.Answer, answer) {
					isFound = true
				}
			}

			if !isFound {
				counter++
			}
		}

		minPoint += lp.Level1.Total/len(lp.Level1.QnA) - counter

		for i, lv := range lp.Level2.QnA {
			isFound := false
			for _, answer := range req.Answer[i+len(lp.Level1.QnA)-1] {
				if array.In(lv.Answer, answer) {
					isFound = true
				}
			}

			if !isFound {
				counter++
			}
		}

		minPoint += lp.Level2.Total/len(lp.Level2.QnA) - counter
	}

	totalPoint := task.Point - minPoint
	err = ctrl.user.UpdatePoint(ctx, req.UserID, totalPoint)
	if err != nil {
		constant.Error.Println("update point ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	userActivity = dbUserActivity.UserActivity{
		UserID:            req.UserID,
		CategoryID:        task.CategoryID,
		LastSubCategoryID: task.SubCategoryID,
		LastTaskID:        task.ID,
		SubCategoryPoint:  userActivity.SubCategoryPoint + totalPoint,
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
	} else if userActivity.RemainingTask == 0 {
		resp.SubCategoryPoint, userActivity.SubCategoryPoint = userActivity.SubCategoryPoint, 0
		user, err := ctrl.user.Get(ctx, req.UserID)
		if err != nil {
			constant.Error.Println("user: get ", err)
			rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
			return
		}

		resp.TotalPoint = user.TotalPoint
	}

	err = ctrl.useractivity.Set(ctx, userActivity, docID)
	if err != nil {
		constant.Error.Println("user activity: set ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, resp)
}
