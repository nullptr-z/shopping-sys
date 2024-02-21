package forms

type PasswordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required" valid:"matches(^[0-9]{11}$)"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

type RegisterUser struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required" valid:"matches(^[0-9]{11}$)"`
	Password   string `form:"password" json:"password" binding:"required,min=3,max=20"`
	ConfirmPwd string `form:"confirmPwd" json:"confirmPwd" binding:"required,min=3,max=20"`
	NickName   string `form:"nickname" json:"nickname" binding:"min=0,max=10"`
}
