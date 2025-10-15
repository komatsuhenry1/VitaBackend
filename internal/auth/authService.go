package auth

import (
	"fmt"
	"medassist/internal/auth/dto"
	"medassist/internal/model"
	"medassist/internal/repository"
	"medassist/utils"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	UserRegister(registerRequestDTO dto.UserRegisterRequestDTO, files map[string][]*multipart.FileHeader) (model.User, error)
	NurseRegister(nurseRequestDTO dto.NurseRegisterRequestDTO, files map[string][]*multipart.FileHeader) (model.Nurse, error)
	LoginUser(loginRequestDTO dto.LoginRequestDTO) (string, dto.AuthUser, error)
	SendCodeToEmail(emailAuthRequestDTO dto.EmailAuthRequestDTO) (dto.CodeResponseDTO, error)
	ValidateUserCode(inputCodeDto dto.InputCodeDto) (string, dto.AuthUser, error)
	FirstLoginAdmin() error
	SendEmailForgotPassword(email dto.ForgotPasswordRequestDTO) error
	ChangePasswordUnlogged(updatedPasswordByNewPassword dto.UpdatedPasswordByNewPassword, id string) error
	ValidateToken(token string) error
	ChangePasswordLogged(changePasswordBothRequestDTO dto.ChangePasswordBothRequestDTO, id string) error
	ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error
}

type authService struct {
	userRepository  repository.UserRepository
	nurseRepository repository.NurseRepository
}

func NewAuthService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository) AuthService {
	return &authService{userRepository: userRepository, nurseRepository: nurseRepository}
}

func (s *authService) UserRegister(registerRequestDTO dto.UserRegisterRequestDTO, files map[string][]*multipart.FileHeader) (model.User, error) {
	if err := registerRequestDTO.Validate(); err != nil {
		return model.User{}, err
	}

	normalizedEmail, err := utils.EmailRegex(registerRequestDTO.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("email invalido")
	}

	_, err = s.userRepository.FindUserByEmail(normalizedEmail)
	if err == nil {
		return model.User{}, fmt.Errorf("O usuário com o email '%s' já existe", normalizedEmail)
	}

	_, err = s.userRepository.FindUserByCpf(registerRequestDTO.Cpf)
	if err == nil {
		return model.User{}, fmt.Errorf("O usuário com o CPF '%s' já existe", registerRequestDTO.Cpf)
	}

	hashedPassword, err := utils.HashPassword(registerRequestDTO.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("Erro ao criptografar senha: %w", err)
	}

	user := model.User{
		ID:          primitive.NewObjectID(),
		Name:        registerRequestDTO.Name,
		Cpf:         registerRequestDTO.Cpf,
		Phone:       registerRequestDTO.Phone,
		Address:     registerRequestDTO.Address,
		Email:       normalizedEmail,
		Password:    hashedPassword,
		Role:        "PATIENT",
		Hidden:      false,
		FirstAccess: true,
		TempCode:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if fileHeaders, ok := files["image_profile"]; ok && len(fileHeaders) > 0 {
		fileHeader := fileHeaders[0] // Pegamos apenas o primeiro arquivo

		file, err := fileHeader.Open()
		if err != nil {
			return model.User{}, fmt.Errorf("erro ao abrir a imagem de perfil: %w", err)
		}
		defer file.Close()

		uniqueFileName := fmt.Sprintf("%s_profile_%s", user.ID.Hex(), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")

		// Chamamos a nova função UploadFile no repositório de usuário
		fileID, err := s.userRepository.UploadFile(file, uniqueFileName, contentType)
		if err != nil {
			return model.User{}, fmt.Errorf("erro no upload da imagem de perfil: %w", err)
		}

		// Atribuímos o ID do arquivo ao nosso modelo de usuário
		user.ProfileImageID = fileID
	}
	// ---- FIM DA NOVA LÓGICA ----

	if err := s.userRepository.CreateUser(&user); err != nil {
		return model.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	if err := utils.SendEmailUserRegister(registerRequestDTO.Email); err != nil {
		return model.User{}, fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return user, nil
}

func (s *authService) NurseRegister(nurseRequestDTO dto.NurseRegisterRequestDTO, files map[string][]*multipart.FileHeader) (model.Nurse, error) {
	if err := nurseRequestDTO.Validate(); err != nil { // valida se nao falta nenhum campo
		return model.Nurse{}, err
	}

	normalizedEmail, err := utils.EmailRegex(nurseRequestDTO.Email)
	if err != nil {
		return model.Nurse{}, fmt.Errorf("email invalido")
	}

	// Verifica se usuário existe (sem erro se não achar)
	_, err = s.nurseRepository.FindNurseByEmail(normalizedEmail)
	if err == nil {
		return model.Nurse{}, fmt.Errorf("Por favor, tente outro email.")
	}

	_, err = s.nurseRepository.FindNurseByCpf(nurseRequestDTO.Cpf)
	if err == nil {
		return model.Nurse{}, fmt.Errorf("Por favor, tente outro CPF.")
	}

	hashedPassword, err := utils.HashPassword(nurseRequestDTO.Password)
	if err != nil {
		return model.Nurse{}, fmt.Errorf("erro ao criptografar senha: %w", err)
	}

	// FUNCAO QUE VALIDA O RG / LICENSE_ID / ANTECEDENTES

	nurse := model.Nurse{
		ID:               primitive.NewObjectID(),
		Name:             nurseRequestDTO.Name,
		Cpf:              nurseRequestDTO.Cpf,
		Phone:            nurseRequestDTO.Phone,
		Address:          nurseRequestDTO.Address,
		Email:            normalizedEmail,
		Password:         hashedPassword,
		PixKey:           nurseRequestDTO.PixKey,
		VerificationSeal: false,

		LicenseNumber:   nurseRequestDTO.LicenseNumber,
		Specialization:  nurseRequestDTO.Specialization,
		Department:      nurseRequestDTO.Department,
		YearsExperience: nurseRequestDTO.YearsExperience,
		Bio:             nurseRequestDTO.Bio,

		Role:        "NURSE",
		Hidden:      false,
		Online:      false,
		FirstAccess: true,
		TempCode:    0,
		StartTime:   nurseRequestDTO.StartTime,
		EndTime:     nurseRequestDTO.EndTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// faz o upload de todos os arquivos e preenche os IDs no objeto nurse
	for fieldName, fileHeaders := range files {
		if len(fileHeaders) == 0 {
			continue // pula se não houver arquivo para este campo
		}
		fileHeader := fileHeaders[0] // pegamos apenas o primeiro arquivo por campo

		file, err := fileHeader.Open()
		if err != nil {
			return model.Nurse{}, fmt.Errorf("erro ao abrir o arquivo %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		// cria um nome de arquivo único e descritivo
		uniqueFileName := fmt.Sprintf("%s_%s_%s", nurse.ID.Hex(), fieldName, fileHeader.Filename) // <nurse_id><license_number><image_name>

		contentType := fileHeader.Header.Get("Content-Type")

		// usa o método genérico do repositório
		fileID, err := s.nurseRepository.UploadFile(file, uniqueFileName, contentType) // sobe pro mongodb esse arquivo e gera o registroem fs.files // retorna o object id que foi criado em fs.files
		if err != nil {
			// se um upload falhar, a operação inteira é cancelada
			return model.Nurse{}, fmt.Errorf("erro no upload do arquivo %s: %w", fileHeader.Filename, err)
		}

		// atribui o id ao campo correto no nosso objeto `nurse`
		switch fieldName {
		case "license_document":
			nurse.LicenseDocumentID = fileID
		case "qualifications":
			nurse.QualificationsID = fileID
		case "general_register":
			nurse.GeneralRegisterID = fileID
		case "residence_comprovant":
			nurse.ResidenceComprovantId = fileID
		case "profile_image":
			nurse.ProfileImageID = fileID
		}
	}

	if err := s.nurseRepository.CreateNurse(&nurse); err != nil {
		return model.Nurse{}, fmt.Errorf("erro ao criar o registro final do enfermeiro(a): %w", err)
	}

	if err := utils.SendEmailNurseRegister(nurseRequestDTO.Email); err != nil {
		return model.Nurse{}, fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return nurse, nil
}

func (s *authService) LoginUser(loginRequestDTO dto.LoginRequestDTO) (string, dto.AuthUser, error) {
	if err := loginRequestDTO.Validate(); err != nil {
		return "", dto.AuthUser{}, err
	}

	loginRequestDTO.Email = strings.ToLower(loginRequestDTO.Email)

	authUser, err := s.findAuthUserByEmail(loginRequestDTO.Email)
	if err != nil {
		return "", dto.AuthUser{}, err
	}

	// O resto das validações continua igual
	if authUser.Role == "NURSE" && !authUser.VerificationSeal {
		return "", dto.AuthUser{}, fmt.Errorf("a conta ainda não foi verificada")
	}

	if authUser.Hidden {
		return "", dto.AuthUser{}, fmt.Errorf("usuário não permitido para login")
	}
	if !utils.ComparePassword(authUser.Password, loginRequestDTO.Password) {
		return "", dto.AuthUser{}, fmt.Errorf("credenciais incorretas")
	}

	if authUser.TwoFactor {
		//gera o codigo
		code, err := utils.GenerateAuthCode()
		if err != nil {
			return "", dto.AuthUser{}, fmt.Errorf("erro ao gerar codigo de verificacao: %w", err)
		}

		// atualiza o campo temp_code no db
		if authUser.Role == "PATIENT" {
			err = s.userRepository.UpdateTempCode(authUser.ID.Hex(), code)
			if err != nil {
				return "", dto.AuthUser{}, fmt.Errorf("erro ao atualizar codigo de verificacao de paciente: %w", err)
			}
		} else {
			err = s.nurseRepository.UpdateTempCode(authUser.ID.Hex(), code)
			if err != nil {
				return "", dto.AuthUser{}, fmt.Errorf("erro ao atualizar codigo de verificacao de enfermeiro: %w", err)
			}
		}

		//manda para o email
		err = utils.SendAuthCode(authUser.Email, code)
		if err != nil {
			return "", dto.AuthUser{}, fmt.Errorf("erro ao enviar email com código de verificação: %w", err)
		}

		user := dto.AuthUser{
			Email:     authUser.Email,
			TwoFactor: authUser.TwoFactor,
			Role:      authUser.Role,
		}

		return "", user, nil
	}

	token, err := utils.GenerateToken(authUser.ID.Hex(), authUser.Role, authUser.Name, authUser.Hidden, time.Hour*168)
	if err != nil {
		return "", dto.AuthUser{}, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return token, authUser, nil
}

func (s *authService) SendCodeToEmail(emailAuthRequestDTO dto.EmailAuthRequestDTO) (dto.CodeResponseDTO, error) {

    authUser, err := s.findAuthUserByEmail(emailAuthRequestDTO.Email)
    if err != nil {
        // O erro "email não cadastrado" da sua função findAuthUserByEmail será retornado aqui.
        return dto.CodeResponseDTO{}, err
    }

	fmt.Println("=======")
	fmt.Println(authUser.ID)
	fmt.Println("=======")

    code, err := utils.GenerateAuthCode()
    if err != nil {
        return dto.CodeResponseDTO{}, fmt.Errorf("Erro ao gerar código de verificação: %w", err)
    }

    switch authUser.Role {
    case "NURSE":
        err = s.nurseRepository.UpdateTempCode(authUser.ID.Hex(), code)
    default:
        err = s.userRepository.UpdateTempCode(authUser.ID.Hex(), code)
    }

    if err != nil {
        return dto.CodeResponseDTO{}, fmt.Errorf("Erro ao atualizar código de verificação: %w", err)
    }

    // O envio do email continua o mesmo.
    err = utils.SendAuthCode(emailAuthRequestDTO.Email, code)
    if err != nil {
        // Boa prática: envolver o erro original para não perder o contexto.
        return dto.CodeResponseDTO{}, fmt.Errorf("erro ao enviar codigo de verificacao: %w", err)
    }

    // O retorno do DTO continua o mesmo.
    codeResponseDTO := dto.CodeResponseDTO{
        Code: code,
    }

    return codeResponseDTO, nil
}

func (s *authService) findAuthUserByEmail(email string) (dto.AuthUser, error) {
	authUser, err := s.userRepository.FindUserByEmail(email)

	if err != nil && err.Error() == "usuário não encontrado" {
		authUser, err = s.nurseRepository.FindNurseByEmail(email)
		if err != nil {
			return dto.AuthUser{}, fmt.Errorf("email não cadastrado")
		}
	} else if err != nil {
		return dto.AuthUser{}, err
	}

	return authUser, nil
}

func (s *authService) ValidateUserCode(inputCodeDto dto.InputCodeDto) (string, dto.AuthUser, error) {
	authUser, err := s.findAuthUserByEmail(inputCodeDto.Email)
	if err != nil {
		return "", dto.AuthUser{}, err
	}

	if inputCodeDto.Code == authUser.TempCode {
		hourExp := time.Hour * 168
		token, err := utils.GenerateToken(authUser.ID.Hex(), authUser.Role, authUser.Name, authUser.Hidden, hourExp)
		if err != nil {
			return "", dto.AuthUser{}, fmt.Errorf("erro ao gerar token")
		}
		return token, authUser, nil
	}

	return "", dto.AuthUser{}, fmt.Errorf("erro ao validar código de usuário.")
}

func (s *authService) FirstLoginAdmin() error {

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	exists, err := s.userRepository.UserExistsByEmail(adminEmail)
	if err != nil {
		return fmt.Errorf("erro ao encontrar: %w", err)
	} else if exists {
		return fmt.Errorf("o usuário já existe")
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("erro ao atualizar campos do usuario: %w", err)
	}

	adminUser := model.User{
		Name:        adminName,
		Email:       adminEmail,
		Password:    hashedPassword, // hash da ADMIN_PASSWORD na .env
		FirstAccess: false,
		Role:        "ADMIN",
	}

	err = s.userRepository.CreateUser(&adminUser)
	if err != nil {
		return err
	}

	utils.SendEmailForAdmin(adminEmail)

	return nil
}

func (s *authService) SendEmailForgotPassword(forgotPasswordRequestDTO dto.ForgotPasswordRequestDTO) error {
	authUser, err := s.userRepository.FindUserByEmail(forgotPasswordRequestDTO.Email)
	if err != nil && err.Error() == "usuário não encontrado" {
		authUser, err = s.nurseRepository.FindNurseByEmail(forgotPasswordRequestDTO.Email)

		if err != nil {
			return fmt.Errorf("Erro ao encontrar enfermeiro(a) para enviar email: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("Erro ao encontrar usuario para enviar email: %w", err)
	}

	expiration := time.Minute * 15

	token, err := utils.GenerateToken(authUser.ID.Hex(), authUser.Role, authUser.Name, authUser.Hidden, expiration)
	if err != nil {
		return fmt.Errorf("erro ao gerar token: %w", err)
	}

	if err := utils.SendEmailForgotPassword(authUser.Email, authUser.ID.Hex(), token); err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return nil
}

func (s *authService) ChangePasswordUnlogged(updatedPasswordByNewPassword dto.UpdatedPasswordByNewPassword, id string) error {
	authUser, err := s.userRepository.FindAuthUserByID(id)

	if err != nil {
		if err.Error() == "usuário não encontrado" {
			authUser, err = s.nurseRepository.FindAuthNurseByID(id)
			if err != nil {
				return fmt.Errorf("usuário ou enfermeiro(a) com o ID fornecido não foi encontrado: %w", err)
			}
		} else {
			return fmt.Errorf("erro ao buscar usuário: %w", err)
		}
	}
	// a senha precisa ter caracteres especiais, numeros e letras
	if !utils.ValidatePassword(updatedPasswordByNewPassword.NewPassword) {
		return fmt.Errorf("senha invalida. A senha precisa ter caracteres especiais, numeros e letras")
	}
	hashedNewPassword, err := utils.HashPassword(updatedPasswordByNewPassword.NewPassword)
	if err != nil {
		return fmt.Errorf("Erro ao criptografar senha: %w", err)
	}

	if authUser.Role == "NURSE" {
		return s.nurseRepository.UpdatePasswordByNurseID(id, hashedNewPassword)
	}
	return s.userRepository.UpdatePasswordByUserID(id, hashedNewPassword)
}

func (s *authService) ChangePasswordLogged(changePasswordBothRequestDTO dto.ChangePasswordBothRequestDTO, id string) error {
	authUser, err := s.userRepository.FindAuthUserByID(id)

	if err != nil {
		if err.Error() == "usuário não encontrado" {
			authUser, err = s.nurseRepository.FindAuthNurseByID(id)
			if err != nil {
				return fmt.Errorf("usuário ou enfermeiro(a) com o ID fornecido não foi encontrado: %w", err)
			}
		} else {
			return fmt.Errorf("erro ao buscar usuário: %w", err)
		}
	}
	if !utils.ComparePassword(authUser.Password, changePasswordBothRequestDTO.Password) {
		return fmt.Errorf("Senha atual incorreta.")
	}
	// a senha precisa ter caracteres especiais, numeros e letras
	if !utils.ValidatePassword(changePasswordBothRequestDTO.NewPassword) {
		return fmt.Errorf("senha invalida. A senha precisa ter caracteres especiais, numeros e letras")
	}
	hashedNewPassword, err := utils.HashPassword(changePasswordBothRequestDTO.NewPassword)
	if err != nil {
		return fmt.Errorf("Erro ao criptografar senha: %w", err)
	}

	if authUser.Role == "NURSE" {
		return s.nurseRepository.UpdatePasswordLoggedByNurseID(id, hashedNewPassword, changePasswordBothRequestDTO.TwoFactor)
	}
	return s.userRepository.UpdatePasswordLoggedByUserID(id, hashedNewPassword, changePasswordBothRequestDTO.TwoFactor)
}

func (s *authService) ValidateToken(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return err
	}
	return nil
}

// auth_service.go

// Use este método em vez de ChangePasswordUnlogged
func (s *authService) ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error {
	// 1. Valida o token e extrai os dados (claims)
	claims, err := utils.ValidateToken(resetPasswordDTO.Token)
	if err != nil {
		return err // Retorna "Token inválido ou expirado"
	}

	// 2. Extrai o ID do usuário de dentro do token
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return fmt.Errorf("token inválido: ID do usuário não encontrado")
	}

	// 3. Valida a complexidade da senha
	if !utils.ValidatePassword(resetPasswordDTO.NewPassword) {
		return fmt.Errorf("senha invalida. A senha precisa ter caracteres especiais, numeros e letras")
	}

	// 4. Criptografa a nova senha
	hashedNewPassword, err := utils.HashPassword(resetPasswordDTO.NewPassword)
	if err != nil {
		return fmt.Errorf("erro ao criptografar senha: %w", err)
	}

	// 5. Busca o usuário pelo ID do token para saber o Role
	authUser, err := s.userRepository.FindAuthUserByID(userID)
	if err != nil && err.Error() == "usuário não encontrado" {
		authUser, err = s.nurseRepository.FindAuthNurseByID(userID)
		if err != nil {
			return fmt.Errorf("usuário referenciado no token não foi encontrado")
		}
	} else if err != nil {
		return fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	// 6. Atualiza a senha no repositório correto
	if authUser.Role == "NURSE" {
		return s.nurseRepository.UpdatePasswordByNurseID(userID, hashedNewPassword)
	}
	return s.userRepository.UpdatePasswordByUserID(userID, hashedNewPassword)
}
