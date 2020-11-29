package handlers

type getResponse struct {
	Rating float64 `json:"rating"`
}

type reviewResponse struct {
	Author string `json:"author"`
	Rating int `json:"rating"`
	Commentary string `json:"commentary"`
}

type reportResponse struct {
	ReportedBy string `json:"reportedBy"`
	ReportedUser string `json:"reportedUser"`
	PackageName string `json:"packageName"`
	Commentary string `json:"commentary"`
	Reviewed bool `json:"reviewed"`
	ReviewedBy string `json:"reviewedBy"`
	ReviewedDate string `json:"reviewedDate"`
	Review string `json:"review"`
}

type rateRequest struct {
	Token string `json:"token"`
	Name string `json:"name"`
	Rating float64 `json:"rating"`
	Comment string `json:"commentary"`
}

type reviewReportRequest struct {
	Token string `json:"token"`
	PackageName string `json:"packageName"`
	Review string `json:"revision"`
	ReviewedBy string `json:"reviewedBy"`
	ReviewedDate string `json:"reviewedDate"`
	ReportedBy string `json:"reportedBy"`
	Ban bool `json:"ban"`
	DeleteReview bool `json:"deleteReview"`
	ReportedUser string `json:"reportedUser"`
}

type deleteRequest struct {
	Token string `json:"token"`
	Package string `json:"package"`
	Author string `json:"author"`
}

type reportsRequest struct {
	Token string `json:"token"`
}

type registerRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type reportRequest struct {
	Token string `json:"token"`
	Commentary string `json:"commentary"`
	ReportedUser string `json:"reportedUser"`
	PackageName string `json:"packageName"`
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
