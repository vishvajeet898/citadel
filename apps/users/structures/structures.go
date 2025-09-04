package structures

// @swagger:response UserDetail
type UserDetail struct {
	// The id of the user.
	// example: 1
	Id uint `json:"id"`
	// The user name of the user.
	// example: "John Doe"
	UserName string `json:"user_name"`
	// The user type of the user.
	// example: "Pathologoist"
	UserType string `json:"user_type"`
	// The email of the user.
	// example: "shrish.chandra@orangehealth.in"
	Email string `json:"email"`
}

type UserDetailDbFilters struct {
	Limit  int
	Offset int
}
