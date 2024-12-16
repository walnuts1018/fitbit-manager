package router

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/consts"
	"github.com/walnuts1018/fitbit-manager/router/handler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(cookieSecret config.CookieSecret, logLevel slog.Level, handler handler.Handler) (*gin.Engine, error) {
	if logLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(sloggin.NewWithConfig(slog.Default(), sloggin.Config{
		DefaultLevel:     logLevel,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         true,
		WithTraceID:        true,

		Filters: []sloggin.Filter{
			sloggin.IgnorePath("/healthz"),
		},
	}))
	r.Use(otelgin.Middleware(consts.ApplicationName))

	r.GET("/healthz", handler.Health)
	store := cookie.NewStore([]byte(cookieSecret))
	r.Use(sessions.Sessions("FitbitManager", store))
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/signin", handler.SignIn)
	r.GET("/callback", handler.Callback)

	v1 := r.Group("/v1")
	{
		v1.GET("/heart", handler.GetHeart)
	}
	return r, nil
}
