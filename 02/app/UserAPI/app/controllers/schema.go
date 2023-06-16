package controllers

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

type encryptPassword struct {
	PassWord string `json:"password"`
}
