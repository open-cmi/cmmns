package middleware

import (
	"fmt"

	"github.com/open-cmi/cmmns/model/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/topmyself/redistore"
)

// AuthMiddleware func
func AuthMiddleware(r *gin.Engine) error {
	store, err := redistore.NewRediStoreWithDB(2000, "tcp", "localhost:25431", "k9Gjj,dZ", "2")
	if err != nil {
		return err
	}

	store.SetKeyPrefix("koa-sess-")
	store.SetSerializer(redistore.JSONSerializer{})
	r.Use(func(c *gin.Context) {
		session, _ := store.Get(c.Request, "koa")
		// Set some session values.
		user, ok := session.Values["user"].(map[string]interface{})
		if ok {
			userid, ok := user["id"].(string)
			if ok {
				var u auth.User
				u.ID = userid
				c.Set("user", u)
			}
		}
		c.Next()
		var postuser auth.User
		fmt.Println(postuser)
		// Save it before we write to the response/return from the handler.
		sessions.Save(c.Request, c.Writer)
	})
	return nil
}

// Init middleware init
func Init(r *gin.Engine) {
	AuthMiddleware(r)
}
