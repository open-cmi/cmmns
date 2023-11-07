package middleware

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateAuthToken("frank", "10001", "frank@163.com", 0, 1, 1)
	if err != nil {
		t.Errorf("user generate token failed")
	}
	claims, err := ParseAuthToken(token)
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
	claims, err = ParseAuthToken(token)
	if err != nil {
		t.Errorf("user token parse failed: %s", err.Error())
		return
	}
	if claims == nil {
		t.Errorf("claims is nil")
		return
	}
}
