package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/infra/oauth2"
)

func NewHandler() (*gin.Engine, error) {
	r := gin.Default()
	store, err := redis.NewStore(10, "tcp", *config.RedisEndpoint, *config.RedisPassword, []byte(*config.CookieSecret))
	if err != nil {
		return nil, err
	}
	r.Use(sessions.Sessions("FitbitManager", store))

	r.GET("/signin", SignIn)
	return r, nil
}

func SignIn(c *gin.Context) {
	session := sessions.Default(c)
	state := "test"
	redirect := oauth2.Auth(state)
	session.Set("state", state)
	session.Save()

	c.Redirect(http.StatusFound, redirect)
}
