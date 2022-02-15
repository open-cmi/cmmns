package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

var storeType string = "redis"
var redisStore *redistore.RediStore
var memoryStore *memstore.MemStore

func GetSession(c *gin.Context) (*sessions.Session, error) {
	if storeType == "memory" {
		return memoryStore.Get(c.Request, "cmmns")
	} else {
		return redisStore.Get(c.Request, "koa")
	}
}

// SessionMiddleware func
func SessionMiddleware(r *gin.Engine) {

	r.Use(func(c *gin.Context) {

		session, _ := GetSession(c)

		// Save it before we write to the response/return from the handler.
		c.Set("session", session)
		sessions.Save(c.Request, c.Writer)
		c.Next()

		// if handler change session, save it
		s, _ := c.Get("session")
		session2, ok := s.(*sessions.Session)
		if ok {
			session2.Save(c.Request, c.Writer)
		}
	})
}

// AuthMiddleware func
func AuthMiddleware(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		s, _ := c.Get("session")
		session, ok := s.(*sessions.Session)
		if !ok {
			c.String(http.StatusUnauthorized, "authenticate is required")
			c.Abort()
			return
		}
		_, ok = session.Values["user"]
		if !ok {
			c.String(http.StatusUnauthorized, "authenticate is required")
			c.Abort()
			return
		}
	})
}

func DefaultMiddleware(r *gin.Engine) {
	r.Use(LoggerWithConfig(LoggerConfig{
		Logger: logger.Logger,
	}), gin.Recovery())
}
