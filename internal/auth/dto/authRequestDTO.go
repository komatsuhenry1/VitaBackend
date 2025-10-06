package dto

import "medassist/utils"

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (u *LoginRequestDTO) Validate() error {
	if u.Email == "" {
		return utils.ErrParamIsRequired("email", "string")
	}
	if u.Password == "" {
		return utils.ErrParamIsRequired("password", "string")
	}
	return nil
}

type UserRegisterRequestDTO struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Cpf      string `json:"cpf"`
	Password string `json:"password"`
}

func (u *UserRegisterRequestDTO) Validate() error {
	if u.Email == "" {
		return utils.ErrParamIsRequired("email", "string")
	}
	if u.Name == "" {
		return utils.ErrParamIsRequired("name", "string")
	}
	if u.Phone == "" {
		return utils.ErrParamIsRequired("phone", "string")
	}
	if u.Address == "" {
		return utils.ErrParamIsRequired("address", "string")
	}
	if u.Cpf == "" {
		return utils.ErrParamIsRequired("cpf", "string")
	}
	if u.Password == "" {
		return utils.ErrParamIsRequired("password", "string")
	}
	return nil
}

type NurseRegisterRequestDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	Cpf             string `json:"cpf"`
	PixKey          string `json:"pix_key"`
	Password        string `json:"password"`
	LicenseNumber   string `json:"license_number"`
	Specialization  string `json:"specialization"`
	Shift           string `json:"shift"` // turno
	Department      string `json:"department"`
	YearsExperience int    `json:"years_experience"`
	Bio             string `json:"bio"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

func (u *NurseRegisterRequestDTO) Validate() error {
	if u.Email == "" {
		return utils.ErrParamIsRequired("email", "string")
	}
	if u.Name == "" {
		return utils.ErrParamIsRequired("name", "string")
	}
	if u.Phone == "" {
		return utils.ErrParamIsRequired("phone", "string")
	}
	if u.Address == "" {
		return utils.ErrParamIsRequired("address", "string")
	}
	if u.Cpf == "" {
		return utils.ErrParamIsRequired("cpf", "string")
	}
	if u.Password == "" {
		return utils.ErrParamIsRequired("password", "string")
	}
	if u.LicenseNumber == "" {
		return utils.ErrParamIsRequired("license_number", "string")
	}
	if u.Specialization == "" {
		return utils.ErrParamIsRequired("specialization", "string")
	}
	if u.Shift == "" {
		return utils.ErrParamIsRequired("shift", "string")
	}
	if u.Department == "" {
		return utils.ErrParamIsRequired("department", "string")
	}
	if u.YearsExperience == 0 {
		return utils.ErrParamIsRequired("years_experience", "int")
	}
	return nil
}

type EmailAuthRequestDTO struct {
	Email string `json:"email"`
}

type InputCodeDto struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type ForgotPasswordRequestDTO struct {
	Email string `json:"email"`
}

type ChangePasswordBothRequestDTO struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type UpdatedPasswordByNewPassword struct {
	NewPassword string `json:"new_password" binding:"required"`
}
