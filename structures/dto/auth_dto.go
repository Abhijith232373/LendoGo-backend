package dto

type SendOTPReq struct {
    FullName string `json:"fullName"`
    Email    string `json:"email"`
}

type VerifyOTPReq struct {
    Email string `json:"email"`
    OTP   string `json:"otp"`
}

type SetPasswordReq struct {
    FullName        string `json:"fullName"` 
    Email           string `json:"email"`
    Password        string `json:"password"`
    ConfirmPassword string `json:"confirmPassword"`
}

type RegisterReq struct {
    FullName string `json:"fullName"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginReq struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserRes struct {
    ID          string          `json:"id"` 
    FullName    string          `json:"full_name"`
    Email       string          `json:"email"`
    Avatar      string          `json:"avatar,omitempty"`
    Role        string          `json:"role"` 
    Status      string          `json:"status"`
    Permissions map[string]bool `json:"permissions,omitempty"`
}

type AuthRes struct {
    Token string  `json:"token"`
    User  UserRes `json:"user"`
}



type ForgotPasswordReq struct {
    Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
    Email           string `json:"email" validate:"required,email"`
    OTP             string `json:"otp" validate:"required,len=6"`
    Password        string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}