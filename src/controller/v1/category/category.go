package category

import (
	"net/http"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	httpCategory "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/category"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/category"

	"github.com/Artexus/api-widyabhuvana/src/util/pagination"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo *category.Repository
}

func NewController(repo *category.Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// Get godoc
// @Tags User
// @Summary Get Category List
// @Description Get category list
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {array} category.GetResponse
// @Router /v1/categories [GET]
func (ctrl Controller) Get(ctx *gin.Context) {
	pgn := pagination.Pagination{}
	err := ctx.BindQuery(&pgn)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	pgn.Paginate()
	categories, err := ctrl.repo.Get(ctx, pgn)
	if err != nil {
		constant.Error.Println("db: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	responses := []httpCategory.GetResponse{}
	copier.Copy(&responses, categories)

	rest.ResponseData(ctx, http.StatusOK, responses)
}
