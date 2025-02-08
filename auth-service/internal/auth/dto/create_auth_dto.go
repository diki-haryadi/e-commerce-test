package authDto

import (
	"database/sql"
	"errors"
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"regexp"
)

type SignUpRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func isEmailOrPhone(value string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	return emailRegex.MatchString(value) || phoneRegex.MatchString(value)
}

func (caDto *SignUpRequestDto) ValidateSignUpDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.Username,
			validator.Required,
			validator.Length(3, 50),
			validator.By(func(value interface{}) error {
				str, ok := value.(string)
				if !ok {
					return errors.New("username must be string")
				}
				if !isEmailOrPhone(str) {
					return errors.New("username must be valid email or phone number")
				}
				return nil
			}),
		),
		validator.Field(
			&caDto.Password,
			validator.Required,
			validator.Length(6, 100),
			validator.Match(regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{6,}$`)),
		),
	)
}

type CreateSignUpResponseDto struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type SignInRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (caDto *SignInRequestDto) ValidateSignInDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.Username,
			validator.Required,
			validator.Length(3, 50),
			validator.By(func(value interface{}) error {
				str, ok := value.(string)
				if !ok {
					return errors.New("username must be string")
				}
				if !isEmailOrPhone(str) {
					return errors.New("username must be valid email or phone number")
				}
				return nil
			}),
		),
		validator.Field(
			&caDto.Password,
			validator.Required,
			validator.Length(6, 100),
			validator.Match(regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{6,}$`)),
		),
	)
}

type CreateSignInResponseDto struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username"`
	Password sql.NullString `json:"password"`
}

type SignInResponse struct {
	User         CreateSignInResponseDto `json:"user"`
	AccessToken  string                  `json:"access_token"`
	RefreshToken string                  `json:"refresh_token"`
	ExpiresIn    int64                   `json:"expires_in"`
}
