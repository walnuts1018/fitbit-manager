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
	state, verifier, redirectURL, err := h.usecase.SignIn()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to sign in: %v", err),
		})
		return
	}
	session.Set("state", state)
	session.Set("verifier", verifier)
	session.Save()

	c.Redirect(http.StatusFound, redirectURL.String())
}

func (h *Handler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	user_id := c.Query("user_id")

	session := sessions.Default(c)
	savedState, ok := session.Get("state").(string)
	if !ok {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "failed to get state",
		})
		return
	}

	if state != savedState {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "invalid state",
		})
		return
	}

	verifier, ok := session.Get("verifier").(string)
	if !ok {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "failed to get verifier",
		})
		return
	}

	if err := h.usecase.Callback(c, user_id, code, verifier); err != nil {
		slog.Error("failed to callback", "error", err)
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  "failed to callback",
		})
		return
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"result": "success",
	})
}
