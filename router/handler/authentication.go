package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(c *gin.Context) {
	session := sessions.Default(c)
	state, redirect, err := h.usecase.SignIn()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to sign in: %v", err),
		})
		return
	}
	session.Set("state", state)
	session.Save()

	c.Redirect(http.StatusFound, redirect)
}

func (h *Handler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	session := sessions.Default(c)

	if session.Get("state") != state {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "invalid state",
		})
		return
	}
	err := h.usecase.Callback(c, code)
	if err != nil {
		slog.Error("failed to callback", "error", err)
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  "failed to callback",
		})
		return
	}

	err = h.usecase.NewFitbitClient(c)
	if err != nil {
		slog.Error("failed to create fitbit client", "error", err)
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  "failed to create fitbit client",
		})
		return
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"result": "success",
	})
}
