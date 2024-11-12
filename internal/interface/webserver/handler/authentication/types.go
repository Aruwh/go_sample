package authentication

type (
	// // // // // // // // // // // // // // // // // // // // // //
	// REQUESTS
	// // // // // // // // // // // // // // // // // // // // // //

	LoginAdminRequest struct {
		Email       string `json:"email"`
		RawPassword string `json:"password"`
	}

	ForgotPwdRequest struct {
		Email string `json:"email"`
	}

	ResetPwdRequest struct {
		Pwd         string `json:"pwd"`
		RepeatedPwd string `json:"repeatedPwd"`
		Token       string `json:"token"`
	}

	// // // // // // // // // // // // // // // // // // // // // //
	// RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //
	LoginAdminResponse struct {
		Token string `json:"token"`
	}

	ForgotPwdResponse struct {
		WasDelivered bool `json:"wasDelivered"`
	}
)

// ConvertFromValue implements a dto like function which will convert a string user facing response
func (LoginAdminResponse) ConvertFromValue(token string) *LoginAdminResponse {
	return &LoginAdminResponse{Token: token}
}

func (ForgotPwdResponse) ConvertFromValue(wasDelivered bool) *ForgotPwdResponse {
	return &ForgotPwdResponse{WasDelivered: wasDelivered}
}
