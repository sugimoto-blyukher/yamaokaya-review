package dtos

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserListItem struct {
	ID		uint	`json:"id"`
	Name 	string	`json:"name"`
}