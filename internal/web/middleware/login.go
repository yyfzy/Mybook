package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(paths ...string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths...)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gob.Register(time.Now())
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
