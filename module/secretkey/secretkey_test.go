package system

import "testing"

func TestGenerateSecretKey(t *testing.T) {
	GenerateSecretKey("test_generate", "rsa", 2048, "zhaodl", "")
}
