package email

import "testing"

func TestSend(t *testing.T) {
	gConf.Server = "smtp.gmail.com"
	gConf.Port = 465
	gConf.User = "test@gmail.com"
	gConf.Password = "FMyc9EJxLd6W8N35"
	gConf.UseTLS = true

	err := Send([]string{"test@163.com"}, "email test", "this is test content", nil)
	if err != nil {
		t.Errorf("send email failed %s\n", err.Error())
	}
}
