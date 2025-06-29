package token

import (
	"testing"
	"time"

	"github.com/open-cmi/gobase/essential/webserver/middleware"
)

func TestGenerateToken(t *testing.T) {
	err := CreateToken("test", "frank", "10001", "frank@163.com", 0, 1, 1)
	if err != nil {
		t.Errorf("user generate token failed")
	}
	m := GetTokenRecord("frank")
	if m == nil {
		t.Errorf("get token failed")
	}
	claims, err := middleware.ParseAuthToken(m.Token)
	if claims == nil || err != nil {
		t.Errorf("user token parse failed")
		return
	}
	if claims.Username != "frank" {
		t.Errorf("username is invalid")
	}
	if claims.UserID != "10001" {
		t.Errorf("id is invalid")
	}
	if claims.Email != "frank@163.com" {
		t.Errorf("email is invalid")
	}
	if claims.Status != 1 {
		t.Errorf("status is invalid")
	}
	if claims.Role != 0 {
		t.Errorf("role is invalid")
	}
	time.Sleep(2 * time.Second)
	claims, err = middleware.ParseAuthToken(m.Token)
	if err != nil {
		t.Errorf("user token parse failed: %s", err.Error())
		return
	}
	if claims == nil {
		t.Errorf("claims is nil")
		return
	}
}
