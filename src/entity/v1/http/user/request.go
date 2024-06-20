package user

type UpdateRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	DateOfBirth *string `json:"dob"`
}
