package handlers

type getResponse struct {
	Name string `json:"name"`
	Rating float64 `json:"rating"`
}

type rateRequest struct {
	Token string `json:"token"`
	Name string `json:"name"`
	Rating float64 `json:"rating"`
	Comment string `json:"comment"`
}

type registerRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}
