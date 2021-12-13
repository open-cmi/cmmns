package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/config"
	"github.com/topmyself/redistore"
)

var store *redistore.RediStore

// AuthMiddleware func
func AuthMiddleware(r *gin.Engine) {

	r.Use(func(c *gin.Context) {

		session, _ := store.Get(c.Request, "koa")

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

// UserPermMiddleware func
func UserPermMiddleware(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		s, ok := c.Get("session")
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

// Init init func
func Init() (err error) {
	host := fmt.Sprintf("%s:%d", config.GetConfig().Rdb.Host, config.GetConfig().Rdb.Port)
	pass := config.GetConfig().Rdb.Password
	store, err = redistore.NewRediStoreWithDB(2000, "tcp", host, pass, "2")
	if err != nil {
		return err
	}

	store.SetKeyPrefix("koa-sess-")
	store.SetSerializer(redistore.JSONSerializer{})
	return nil
}
