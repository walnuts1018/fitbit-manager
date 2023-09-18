package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

var (
	tokenUsecase *usecase.TokenUsecase
)

func NewHandler(usecase *usecase.TokenUsecase) (*gin.Engine, error) {
	tokenUsecase = usecase

	r := gin.Default()
	store := cookie.NewStore([]byte(config.CookieSecret))
	r.Use(sessions.Sessions("FitbitManager", store))

	r.GET("/signin", signIn)
	r.GET("/callback", callback)
	return r, nil
}

func signIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	state, redirect, err := tokenUsecase.SignIn()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed to sign in: %v", err)
		return
	}
	session.Set("state", state)
	session.Save()

	ctx.Redirect(http.StatusFound, redirect)
}

func callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	session := sessions.Default(ctx)
	if session.Get("state") != state {
		ctx.String(http.StatusBadRequest, "invalid state")
		return
	}
	err := tokenUsecase.Callback(ctx, code)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed to callback: %v", err)
		return
	}
}
