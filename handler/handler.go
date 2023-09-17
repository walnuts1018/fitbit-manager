package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/infra/oauth2"
)

func NewHandler() (*gin.Engine, error) {
	r := gin.Default()
	store, err := redis.NewStore(10, "tcp", config.RedisEndpoint, config.RedisPassword, []byte(config.CookieSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %v", err)
	}
	r.Use(sessions.Sessions("FitbitManager", store))

	r.GET("/signin", signIn)
	r.GET("/callback", callback)
	return r, nil
}

func signIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	state, err := randStr(64)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed to generate state")
		return
	}
	redirect := oauth2.Auth(state)
	session.Set("state", state)
	session.Save()

	ctx.Redirect(http.StatusFound, redirect)
}

func callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	session := sessions.Default(ctx)
	fmt.Printf("gotState: %v, state: %v \n", session.Get("state"), state)
	if session.Get("state") != state {
		ctx.String(http.StatusBadRequest, "invalid state")
		return
	}
	at, err := oauth2.Callback(ctx, code)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed to get access token")
		return
	}
	fmt.Println(at)
}

func randStr(n int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
