package user

// ChangePasswordMsg change password msg
type ChangePasswordMsg struct {
	OldPassword     string `json:"oldpass"`
	NewPassword     string `json:"newpass"`
	ConfirmPassword string `json:"confirmpass"`
}

type CreateMsg struct {
	UserName    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirmpass"`
	Role        string `json:"role"`
	Description string `json:"description,omitempty"`
}

// LoginMsg struct
type LoginMsg struct {
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Captcha       string `json:"captcha,omitempty"`
	CaptchaID     string `json:"captchaid,omitempty"`
	IgnoreCaptcha bool   `json:"ignore_captcha,omitempty"`
}

type EditMsg struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

type ResetPasswdRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpass"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}
