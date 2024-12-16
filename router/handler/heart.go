package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHeart(c *gin.Context) {
	heart, t, err := h.usecase.GetHeartNow(c)
	if err != nil {
		slog.Error("failed to get heart", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get heart",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"heart": heart,
		"time":  t,
	})
}
