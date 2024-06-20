package user

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	httpUser "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/user"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/Artexus/api-widyabhuvana/src/util/storage"
	"github.com/google/uuid"
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

	resp.PhotoURL = fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", constant.Bucket, url.QueryEscape(user.PhotoURL))

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
	req.DateOfBirth = ctx.PostForm("dob")
	req.Email = ctx.PostForm("email")
	req.Name = ctx.PostForm("name")
	err := ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		constant.Error.Println("parse multipart form ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	updates := []firestore.Update{}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		constant.Error.Println("file form request ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpeg" && ext != ".jpg" {
			rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
				"file": constant.ErrInvalid.Error(),
			})
			return
		}

		uuid := uuid.New()
		path := fmt.Sprintf("images/%s/%s%s", id, uuid.String(), ext)
		err = storage.Upload(file, uuid.String(), path)
		if err != nil {
			constant.Error.Println("upload ", err)
			rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
			return
		}

		updates = append(updates, firestore.Update{
			Path:  "photo_url",
			Value: path,
		})
	}

	if req.DateOfBirth != "" {
		updates = append(updates, firestore.Update{
			Path:  "dob",
			Value: req.DateOfBirth,
		})
	}

	if req.Email != "" {
		updates = append(updates, firestore.Update{
			Path:  "email",
			Value: req.Email,
		})
	}

	if req.Name != "" {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: req.Name,
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
