package user

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	httpUser "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/user"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	repo *user.Repository
}

func NewController(repo *user.Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// Get godoc
// @Tags User
// @Summary Get Auth User
// @Description Get user detail by its bearer token
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {object} user.GetResponse
// @Router /v1/users [GET]
func (ctrl Controller) Get(ctx *gin.Context) {
	id, _ := jwt.ExtractIDToken(ctx.GetHeader("Authorization"))
	user, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		constant.Error.Println("db: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	resp := httpUser.GetResponse{}
	copier.Copy(&resp, user)

	rest.ResponseData(ctx, http.StatusOK, resp)
}

// Update godoc
// @Tags User
// @Summary Update User
// @Description Update user
// @Param Authorization header string true "Bearer Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 200 {object} user.GetResponse
// @Router /v1/users [PATCH]
func (ctrl Controller) Update(ctx *gin.Context) {
	id, _ := jwt.ExtractIDToken(ctx.GetHeader("Authorization"))
	req := httpUser.UpdateRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	updates := []firestore.Update{}
	if req.DateOfBirth != nil {
		updates = append(updates, firestore.Update{
			Path:  "dob",
			Value: *req.DateOfBirth,
		})
	}

	if req.Email != nil {
		updates = append(updates, firestore.Update{
			Path:  "email",
			Value: *req.Email,
		})
	}

	if req.Name != nil {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: *req.Name,
		})
	}

	err = ctrl.repo.Update(ctx, id, updates)
	if err != nil {
		constant.Error.Println("db: update ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	rest.ResponseOutput(ctx, http.StatusOK, nil)
}
