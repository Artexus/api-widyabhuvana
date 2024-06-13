package subcategory

import (
	"errors"
	"net/http"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	httpSubCategory "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/subcategory"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subcategory"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/useractivity"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo         *subcategory.Repository
	useractivity *useractivity.Repository
}

func NewController(repo *subcategory.Repository, useractivity *useractivity.Repository) *Controller {
	return &Controller{
		repo:         repo,
		useractivity: useractivity,
	}
}

// Get godoc
// @Tags User
// @Summary Get Sub Category List
// @Description Get category list
// @Param category_id query string true "Category ID"
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {array} category.GetResponse
// @Router /v1/sub-categories [GET]
func (ctrl Controller) Get(ctx *gin.Context) {
	req := httpSubCategory.GetRequest{}
	err := ctx.BindQuery(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	req.CategoryID, err = aes.DecryptID(req.EncCategoryID)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"category_id": constant.ErrInvalid.Error(),
		})
		return
	}

	req.UserID, _ = jwt.ExtractIDToken(ctx.GetHeader("Authorization"))
	categories, err := ctrl.repo.Get(ctx, req.CategoryID)
	if err != nil {
		constant.Error.Println("db: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	_, entity, err := ctrl.useractivity.Get(ctx, req.UserID, req.CategoryID)
	if err != nil && !errors.Is(err, constant.ErrNotFound) {
		constant.Error.Println("useractivity: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	responses := []httpSubCategory.GetResponse{}
	for _, category := range categories {
		if errors.Is(err, constant.ErrNotFound) {
			entity.RemainingTask = len(category.Tasks)
		}

		resp := httpSubCategory.GetResponse{
			EncID:     aes.EncryptID(category.CategoryID),
			Name:      category.Name,
			MaxPoint:  category.MaxPoint,
			Tasks:     aes.EncryptIDs(category.Tasks[len(category.Tasks)-entity.RemainingTask:]),
			TotalTask: len(category.Tasks),
		}

		responses = append(responses, resp)
	}

	rest.ResponseData(ctx, http.StatusOK, responses)
}
