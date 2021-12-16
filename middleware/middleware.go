package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/memstore"
	"github.com/topmyself/redistore"
)

var storeType string = "redis"
var redisStore *redistore.RediStore
var memoryStore *memstore.MemStore

func GetSession(c *gin.Context) (*sessions.Session, error) {
	if storeType == "memory" {
		return memoryStore.Get(c.Request, "cmmnsmemory")
	} else {
		return redisStore.Get(c.Request, "koa")
	}
}

// AuthMiddleware func
func AuthMiddleware(r *gin.Engine) {

	r.Use(func(c *gin.Context) {

		session, _ := GetSession(c)

		fmt.Println("session: ", session)
		// Save it before we write to the response/return from the handler.
		c.Set("session", session)
		sessions.Save(c.Request, c.Writer)
		c.Next()

		// if handler change session, save it
		s, _ := c.Get("session")
		session2, ok := s.(*sessions.Session)
		fmt.Println("session2: ", session2, ok)
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
			fmt.Println("session not ok")
			c.String(http.StatusUnauthorized, "authenticate is required")
			c.Abort()
			return
		}
		_, ok = session.Values["user"]
		if !ok {
			fmt.Println("user not ok")
			c.String(http.StatusUnauthorized, "authenticate is required")
			c.Abort()
			return
		}
	})
}

// Init init func
func Init() (err error) {
	if os.Getenv("SESSION_STORE") == "memory" {
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		storeType = "memory"
	} else {
		host := fmt.Sprintf("%s:%d", config.GetConfig().Rdb.Host, config.GetConfig().Rdb.Port)
		pass := config.GetConfig().Rdb.Password
		redisStore, err = redistore.NewRediStoreWithDB(2000, "tcp", host, pass, "2")
		if err != nil {
			return err
		}

		redisStore.SetKeyPrefix("koa-sess-")
		redisStore.SetSerializer(redistore.JSONSerializer{})
	}

	return nil
}
