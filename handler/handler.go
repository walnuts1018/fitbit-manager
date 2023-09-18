package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

var (
	uc *usecase.Usecase
)

func NewHandler(usecase *usecase.Usecase) (*gin.Engine, error) {
	uc = usecase

	r := gin.Default()
	store := cookie.NewStore([]byte(config.CookieSecret))
	r.Use(sessions.Sessions("FitbitManager", store))
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/signin", signIn)
	r.GET("/callback", callback)

	v1 := r.Group("/v1")
	{
		v1.GET("/heart", getHeart)
	}
	return r, nil
}

func signIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	state, redirect, err := uc.SignIn()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to sign in: %v", err),
		})
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
		ctx.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "invalid state",
		})
		return
	}
	err := uc.Callback(ctx, code)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to callback: %v", err),
		})
		return
	}

	err = uc.NewFitbitClient(ctx)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to create fitbit client: %v", err),
		})
		return
	}

	ctx.HTML(http.StatusOK, "result.html", gin.H{
		"result": "success",
	})
}

func getHeart(ctx *gin.Context) {
	heart, t, err := uc.GetHeartNow(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get heart: %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"heart": heart,
		"time":  t,
	})
}
