package auth

import (
	"errors"
	"net/http"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/user"

	httpAuth "github.com/Artexus/api-widyabhuvana/src/entity/v1/http/auth"
	"github.com/Artexus/api-widyabhuvana/src/repository/v1/user"
	"github.com/Artexus/api-widyabhuvana/src/util/hash"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/Artexus/api-widyabhuvana/src/util/validate"
	"github.com/go-playground/validator/v10"

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

// Login godoc
// @Tags Auth
// @Summary Login User
// @Description Login User
// @Param body body auth.LoginRequest true "Body"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 201 {object} auth.LoginResponse
// @Router /v1/login [POST]
func (ctrl Controller) Login(ctx *gin.Context) {
	req := httpAuth.LoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	if err = validator.New().Struct(&req); err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, err)
		return
	}

	entity, err := ctrl.repo.TakeByEmail(ctx, req.Email)
	if errors.Is(err, constant.ErrNotFound) {
		rest.ResponseOutput(ctx, http.StatusNotFound, map[string]string{
			"user": constant.ErrNotFound.Error(),
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if entity.Password != hash.GenerateHashToken(req.Password) {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"password": constant.ErrInvalid.Error(),
		})
		return
	}

	token, err := jwt.GenerateToken(jwt.GenerateTokenPayload{
		EncUserID: entity.EncID(),
		Email:     entity.Email,
		Username:  entity.Username,
	})
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	response := httpAuth.LoginResponse{
		Username: entity.Username,
		Token:    token,
	}

	ctx.JSON(http.StatusCreated, response)
}

// Register godoc
// @Tags Auth
// @Summary Register User
// @Description Register User
// @Description `talent` request required if `role_id` is for talent. vice versa
// @Param body body auth.RegisterRequest true "Body"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Success 201 {object} auth.RegisterResponse
// @Router /v1/register [POST]
func (ctrl Controller) Register(ctx *gin.Context) {
	req := httpAuth.RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	if err = validator.New().Struct(&req); err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, err)
		return
	}

	err = validate.ValidateEmail(req.Email)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"email": constant.ErrInvalid.Error(),
		})
		return
	}

	err = validate.ValidatePassword(req.Password)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusBadRequest, map[string]string{
			"password": constant.ErrInvalid.Error(),
		})
		return
	}

	user, err := ctrl.repo.TakeByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, constant.ErrNotFound) {
		constant.Error.Println("db: take by email: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	if user.ID != "" {
		rest.ResponseOutput(ctx, http.StatusConflict, map[string]string{
			"email": constant.ErrRegistered.Error(),
		})
		return
	}

	entity := db.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hash.GenerateHashToken(req.Password),
	}

	err = ctrl.repo.Create(ctx, &entity)
	if err != nil {
		constant.Error.Println("db: create: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	token, err := jwt.GenerateToken(jwt.GenerateTokenPayload{
		EncUserID: entity.EncID(),
		Email:     entity.Email,
		Username:  entity.Username,
	})
	if err != nil {
		constant.Error.Println("db: generate token: ", err.Error())
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		return
	}

	response := httpAuth.RegisterResponse{
		EncID: entity.EncID(),
		Token: token,
	}

	rest.ResponseOutput(ctx, http.StatusCreated, response)
}
