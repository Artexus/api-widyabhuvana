package middleware

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
	"github.com/Artexus/api-widyabhuvana/src/util/jwt"
	"github.com/Artexus/api-widyabhuvana/src/util/rest"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Middleware struct {
	firestore *firestore.Client
}

func New(firestore *firestore.Client) *Middleware {
	return &Middleware{
		firestore: firestore,
	}
}

func (m Middleware) Auth(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	b := "Bearer "
	if !strings.Contains(token, b) {
		rest.ResponseOutput(ctx, http.StatusUnauthorized, nil)
		ctx.Abort()
		return
	}
	t := strings.Split(token, b)
	if len(t) < 2 {
		rest.ResponseOutput(ctx, http.StatusUnauthorized, nil)
		ctx.Abort()
		return
	}

	claims, err := jwt.ExtractToken(token)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusForbidden, map[string]string{
			"token": err.Error(),
		})
		ctx.Abort()
		return
	}

	id, err := aes.DecryptID(claims.EncID)
	if err != nil {
		rest.ResponseOutput(ctx, http.StatusUnauthorized, nil)
		ctx.Abort()
		return
	}

	_, err = m.firestore.Collection("users").Doc(id).Get(context.Background())
	if err != nil && status.Code(err) != codes.NotFound {
		constant.Error.Println("db: get ", err)
		rest.ResponseOutput(ctx, http.StatusInternalServerError, nil)
		ctx.Abort()
		return
	}

	if status.Code(err) == codes.NotFound {
		rest.ResponseOutput(ctx, http.StatusUnauthorized, nil)
		ctx.Abort()
		return
	}
}
