package user

// ChangePasswordMsg change password msg
type ChangePasswordMsg struct {
	OldPassword     string `json:"oldpass"`
	NewPassword     string `json:"newpass"`
	ConfirmPassword string `json:"confirmpass"`
}

// RegisterMsg register msg
type RegisterMsg struct {
	UserName      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Description   string `json:"description,omitempty"`
	Confirmpass   string `json:"confirmpass"`
	Captcha       string `json:"Captcha,omitempty"`
	CaptchaID     string `json:"captchaid,omitempty"`
	IgnoreCaptcha bool   `json:"ignorecaptcha"`
}

type CreateMsg struct {
	UserName    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Confirmpass string `json:"confirmpass"`
	Description string `json:"description,omitempty"`
}

// LoginMsg struct
type LoginMsg struct {
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Captcha       string `json:"captcha,omitempty"`
	CaptchaID     string `json:"captchaid,omitempty"`
	IgnoreCaptcha bool   `json:"ignorecaptcha,omitempty"`
}
