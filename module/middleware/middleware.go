package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/pkg/memstore"
	"github.com/open-cmi/cmmns/pkg/redistore"
)

var redisStore *redistore.RediStore
var memoryStore *memstore.MemStore

func GetSession(c *gin.Context) (*sessions.Session, error) {
	if gConf.Store == "memory" {
		return memoryStore.Get(c.Request, "koa")
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
		user, ok := session.Values["user"]
		if ok {
			c.Set("user", user)
		}
		c.Next()

		// if handler change session, save it
		s, _ := c.Get("session")
		session2, ok := s.(*sessions.Session)
		if ok {
			session2.Save(c.Request, c.Writer)
		}
	})
}

func JWTMiddleware(r *gin.Engine) {

	r.Use(func(c *gin.Context) {
		// 这里校验token
		tokenstr := c.Request.Header.Get("Authorization")
		if strings.HasPrefix(tokenstr, "Bearer ") {
			token := strings.TrimPrefix(tokenstr, "Bearer ")

			// 这里做验证
			claims, err := ParseAuthToken(token)
			if err != nil {
				logger.Errorf("parse auth token: %s\n", err.Error())
			} else if claims != nil {
				user := make(map[string]interface{})
				user["username"] = claims.Username
				user["id"] = claims.UserID
				user["email"] = claims.Email
				user["role"] = claims.Role
				user["status"] = claims.Status
				c.Set("user", user)
			}
		}
		c.Next()
	})
}

// AuthMiddleware func
func AuthMiddleware(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		_, exist := c.Get("user")
		if !exist {
			c.String(http.StatusUnauthorized, "authenticate is required")
			c.Abort()
		}
		c.Next()
	})
}

func DefaultMiddleware(r *gin.Engine) {
	r.Use(LoggerWithConfig(LoggerConfig{
		Logger: logger.Logger,
	}), gin.Recovery())
}

type UserClaims struct {
	UserID   string
	Username string
	Email    string
	Role     int
	Status   int
	jwt.RegisteredClaims
}

func ParseAuthToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("cmmns"), nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			r := GetTokenRecordByToken(token)
			if r == nil {
				return nil, errors.New("token not exist")
			}
			return claims, nil
		}
	}

	return nil, err
}

func GenerateAuthToken(name string, username string, id string, email string, role int, status int, expireDay int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3600 * 24 * (time.Duration)(expireDay) * time.Second)
	issuer := "cmmns"
	claims := UserClaims{
		UserID:   id,
		Username: username,
		Role:     role,
		Email:    email,
		Status:   status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("cmmns"))
	if err != nil {
		return "", err
	}

	t := NewTokenRecord()
	t.ExpireDay = expireDay
	t.Token = token
	t.Name = name
	err = t.Save()
	if err != nil {
		return "", err
	}
	return token, err
}

func DeleteAuthToken(name string) error {
	t := GetTokenRecord(name)
	if t == nil {
		return errors.New("token not existed")
	}
	err := t.Remove()
	if err != nil {
		return err
	}
	return nil
}
