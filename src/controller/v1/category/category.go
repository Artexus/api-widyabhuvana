package category

import (
	"errors"
	"net/http"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/category"
	dbSub "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/subcategory"

	httpCategory "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/category"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/category"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/subcategory"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/useractivity"

	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/pagination"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo         *category.Repository
	useractivity *useractivity.Repository
	subcategory  *subcategory.Repository
}

func NewController(repo *category.Repository, useractivity *useractivity.Repository, subcategory *subcategory.Repository) *Controller {
	return &Controller{
		repo:         repo,
		useractivity: useractivity,
		subcategory:  subcategory,
	}
}

// Get godoc
// @Tags Category
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

// Get godoc
// @Tags Category
// @Summary Get User's progress list
// @Description Get User's progress list
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {array} category.GetResponse
// @Router /v1/categories/progress [GET]
func (ctrl Controller) GetUserProgress(ctx *gin.Context) {
	id, _ := jwt.ExtractIDToken(ctx.GetHeader("Authorization"))
	userActivities, err := ctrl.useractivity.GetByUserID(ctx, id)
	if err != nil && !errors.Is(err, constant.ErrNotFound) {
		constant.Error.Println("useractivity: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	categoryIDs := []string{}
	for _, userActivity := range userActivities {
		categoryIDs = append(categoryIDs, userActivity.CategoryID)
	}

	categories, err := ctrl.repo.GetByIDs(ctx, categoryIDs)
	if err != nil {
		constant.Error.Println("db: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	mapCategories := map[string]db.Category{}
	for _, category := range categories {
		mapCategories[category.ID] = category
	}

	subCategoryIDs := make([]string, 0)
	responses := []httpCategory.GetUserProgressResponse{}
	mapSubCategoryIDs := map[string]string{}
	for _, userActivity := range userActivities {
		subCategory := userActivity.LastSubCategoryID
		if userActivity.RemainingTask == 0 {
			subCategory = mapCategories[userActivity.CategoryID].SubCategories[len(mapCategories[userActivity.CategoryID].SubCategories)-userActivity.RemainingSubCategory]
		}

		mapSubCategoryIDs[userActivity.CategoryID] = subCategory
		subCategoryIDs = append(subCategoryIDs, subCategory)
	}

	subCategories, err := ctrl.subcategory.GetByIDs(ctx, subCategoryIDs)
	if err != nil {
		constant.Error.Println("subcategory: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	mapSubCategories := map[string]dbSub.SubCategory{}
	for _, subCategory := range subCategories {
		mapSubCategories[subCategory.CategoryID] = subCategory
	}

	for _, userActivity := range userActivities {
		resp := httpCategory.GetUserProgressResponse{
			Category: httpCategory.Category{
				EncID: aes.EncryptID(userActivity.CategoryID),
				Name:  mapCategories[userActivity.CategoryID].Name,
			},

			SubCategory: httpCategory.SubCategory{
				EncID: aes.EncryptID(mapSubCategories[userActivity.CategoryID].ID),
				Name:  mapSubCategories[userActivity.CategoryID].Name,
			},

			TotalSubCategory:     len(mapCategories[userActivity.CategoryID].SubCategories),
			RemainingSubCategory: userActivity.RemainingSubCategory,
		}

		responses = append(responses, resp)
	}

	rest.ResponseData(ctx, http.StatusOK, responses)
}
