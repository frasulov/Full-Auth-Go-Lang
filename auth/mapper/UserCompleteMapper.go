package mapper

import (
	"auth/dto"
	"auth/models"
	"database/sql"
)

func UserCompleteMapper(user *models.User, dto dto.RegisterUserStep2Dto) {
	user.Username = sql.NullString{dto.Username, true}
	user.BirthDate = dto.BirthDate
	user.PhoneNumber = sql.NullString{dto.PhoneNumber, true}
	user.Country = sql.NullString{dto.Country, true}
	user.IsCompleted = true
}
