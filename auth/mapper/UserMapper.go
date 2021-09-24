package mapper

import (
	"auth/dto"
	"auth/models"
	"database/sql"
)

func ToUser(userDto *dto.RegisterUserDto) *models.User {
	return &models.User{
		Email:     userDto.Email,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Password:  userDto.Password,
		Role:      sql.NullString{userDto.Role, true},
	}
}

func ToUserDto(user *models.User) *dto.RegisterUserDto {
	return &dto.RegisterUserDto{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role.String,
	}
}
