package users_transport_http

import "github.com/nullphosphor/golang-todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID      int `json:"id"`
	Version int `json:"version"`

	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func UserDTOFromDomain(user domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = UserDTOResponse(UserDTOFromDomain(user))
	}

	return usersDTO
}
