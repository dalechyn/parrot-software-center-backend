package handlers

type getResponse struct {
	Rating float64 `json:"rating"`
}

type reviewResponse struct {
	Author string `json:"author"`
	Rating int `json:"rating"`
	Commentary string `json:"commentary"`
}

type rateRequest struct {
	Token string `json:"token"`
	Name string `json:"name"`
	Rating float64 `json:"rating"`
	Comment string `json:"commentary"`
}

type deleteRequest struct {
	Token string `json:"token"`
	Package string `json:"package"`
	Author string `json:"author"`
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
	Role string `json:"role"`
}

const RoleUser = "user"
const RoleModerator = "moderator"
