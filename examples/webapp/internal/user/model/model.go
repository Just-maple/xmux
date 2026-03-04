package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserRequest struct {
	ID string `json:"id"`
}

type UpdateUserRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
