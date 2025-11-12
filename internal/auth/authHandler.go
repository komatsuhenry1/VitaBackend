package auth

import (
	"fmt"
	"log"
	"medassist/internal/auth/dto"
	"medassist/utils"
	"net/http"
	"strconv"
	_"medassist/internal/model"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Registro de Novo Usuário (Paciente)
// @Description Cria um novo usuário (paciente) no sistema, permitindo o upload de uma imagem de perfil.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Nome completo do usuário"
// @Param email formData string true "Email válido do usuário"
// @Param phone formData string true "Telefone do usuário (com DDD)"
// @Param neighborhood formData string true "Bairro"
// @Param city formData string true "Cidade"
// @Param uf formData string true "Estado (UF)"
// @Param complement formData string false "Complemento do endereço (Opcional)"
// @Param number formData string true "Número do endereço"
// @Param street formData string true "Rua/Avenida"
// @Param cep formData string true "CEP"
// @Param cpf formData string true "CPF do usuário"
// @Param password formData string true "Senha (deve seguir as regras de complexidade)"
// @Param image_profile formData file false "Imagem de perfil (Opcional)"
// @Success 200 {object} model.User "Usuário criado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Dados inválidos (validação falhou, e-mail/CPF duplicado, erro no arquivo)"
// @Router /auth/user [post]// @Description Cria um novo usuário (paciente) no sistema, permitindo o upload de uma imagem de perfil.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// 💡 MUDANÇA AQUI: Apontamos diretamente para o models.User.
// Embora você retorne gin.H{"user": ...}, isso diz ao Swagger
// qual é a estrutura de dados principal retornada.
// @Success 200 {object} model.User "Usuário criado com sucesso"
//
// 💡 MUDANÇA AQUI: Apontamos para a nova struct que acabamos de criar.
// @Failure 400 {object} utils.ErrorResponse "Dados inválidos (validação falhou, e-mail/CPF duplicado, erro no arquivo)"
// @Router /auth/user [post]
func (h *AuthHandler) UserRegister(c *gin.Context) {
	// 1. Criar o DTO e preenchê-lo com os dados do formulário
	var userRequestDTO dto.UserRegisterRequestDTO
	userRequestDTO.Name = c.PostForm("name")
	userRequestDTO.Email = c.PostForm("email")
	userRequestDTO.Phone = c.PostForm("phone")
	userRequestDTO.Neighborhood = c.PostForm("neighborhood")
	userRequestDTO.City = c.PostForm("city")
	userRequestDTO.UF = c.PostForm("uf")
	userRequestDTO.Complement = c.PostForm("complement")
	userRequestDTO.Number = c.PostForm("number")
	userRequestDTO.Street = c.PostForm("street")
	userRequestDTO.CEP = c.PostForm("cep")
	userRequestDTO.Cpf = c.PostForm("cpf")
	userRequestDTO.Password = c.PostForm("password")

	form, err := c.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(c, "Erro ao processar formulário multipart: "+err.Error(), http.StatusBadRequest)
		return
	}
	files := form.File

	createdUser, err := h.authService.UserRegister(userRequestDTO, files)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "usuário criado com sucesso", gin.H{"user": createdUser})
}

// @Summary Registro de Novo Enfermeiro (Nurse)
// @Description Cria uma nova solicitação de cadastro de enfermeiro, com dados e upload de documentos obrigatórios.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Nome completo"
// @Param email formData string true "Email válido"
// @Param phone formData string true "Telefone (com DDD)"
// @Param cep formData string true "CEP"
// @Param street formData string true "Rua"
// @Param number formData string true "Número"
// @Param complement formData string false "Complemento (Opcional)"
// @Param neighborhood formData string true "Bairro"
// @Param city formData string true "Cidade"
// @Param uf formData string true "UF"
// @Param cpf formData string true "CPF"
// @Param pix_key formData string true "Chave PIX"
// @Param password formData string true "Senha"
// @Param coren formData string true "Registro Coren"
// @Param specialization formData string true "Especialização (ex: Pediatria)"
// @Param department formData string true "Departamento (ex: Enfermagem)"
// @Param years_experience formData int true "Anos de Experiência (número)"
// @Param bio formData string true "Biografia / Descrição breve"
// @Param start_time formData string true "Horário de início (ex: 08:00)"
// @Param end_time formData string true "Horário de término (ex: 18:00)"
// @Param license_document formData file true "Documento de Licença (CNH/RG)"
// @Param qualifications formData file true "Comprovante de Qualificações/Diplomas"
// @Param general_register formData file true "Registro Geral (RG)"
// @Param residence_comprovant formData file true "Comprovante de Residência"
// @Param profile_image formData file true "Imagem de Perfil"
// @Success 200 {object} utils.SuccessResponseNurse "Cadastro solicitado com sucesso (retorna o objeto Nurse)"
// @Failure 400 {object} utils.ErrorResponse "Dados inválidos, arquivos faltando ou formato incorreto"
// @Router /auth/nurse [post]// @Description Cria uma nova solicitação de cadastro de enfermeiro, com dados e upload de documentos obrigatórios.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} utils.SuccessResponseNurse "Cadastro solicitado com sucesso (retorna o objeto Nurse)"
// @Failure 400 {object} utils.ErrorResponse "Dados inválidos, arquivos faltando ou formato incorreto"
// @Router /auth/nurse [post]
func (h *AuthHandler) NurseRegister(c *gin.Context) {

	yearsExpStr := c.PostForm("years_experience")
	yearsExp, err := strconv.Atoi(yearsExpStr)
	if err != nil {
		utils.SendErrorResponse(c, "Formato inválido para 'anos de experiência'. Esperado um número.", http.StatusBadRequest)
		return
	}

	var nurseRequestDTO dto.NurseRegisterRequestDTO
	nurseRequestDTO.Name = c.PostForm("name")
	nurseRequestDTO.Email = c.PostForm("email")
	nurseRequestDTO.Phone = c.PostForm("phone")

	nurseRequestDTO.CEP = c.PostForm("cep")
	nurseRequestDTO.Street = c.PostForm("street")
	nurseRequestDTO.Number = c.PostForm("number")
	nurseRequestDTO.Complement = c.PostForm("complement")
	nurseRequestDTO.Neighborhood = c.PostForm("neighborhood")
	nurseRequestDTO.City = c.PostForm("city")
	nurseRequestDTO.UF = c.PostForm("uf")

	nurseRequestDTO.Cpf = c.PostForm("cpf")
	nurseRequestDTO.PixKey = c.PostForm("pix_key")
	nurseRequestDTO.Password = c.PostForm("password")
	nurseRequestDTO.Coren = c.PostForm("coren")
	nurseRequestDTO.Specialization = c.PostForm("specialization")
	nurseRequestDTO.Department = c.PostForm("department")
	nurseRequestDTO.YearsExperience = yearsExp
	nurseRequestDTO.Bio = c.PostForm("bio")
	nurseRequestDTO.StartTime = c.PostForm("start_time")
	nurseRequestDTO.EndTime = c.PostForm("end_time")

	form, err := c.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(c, "Erro ao processar formulário: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := form.File // todos arquivos enviados

	requiredFiles := []string{"license_document", "qualifications", "general_register", "residence_comprovant", "profile_image"}
	for _, fieldName := range requiredFiles {
		fmt.Println(requiredFiles)
		if _, ok := files[fieldName]; !ok || len(files[fieldName]) == 0 {
			utils.SendErrorResponse(c, "Arquivo obrigatório não enviado: "+fieldName, http.StatusBadRequest)
			return
		}
	}

	createdNurse, err := h.authService.NurseRegister(nurseRequestDTO, files) // passa files para poder ser salvo no mongo
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Cadastro de enfermeiro solicitado com sucesso.", gin.H{"nurse": createdNurse})
}

// @Summary Login de Usuário
// @Description Autentica um usuário (paciente ou enfermeiro) e retorna um token JWT e dados do usuário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequestDTO true "Credenciais de Login (email e senha)"
// @Success 200 {object} utils.SuccessValidateCodeResponse "Login bem-sucedido"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida ou credenciais incorretas"
// @Router /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var userLoginRequestDTO dto.LoginRequestDTO
	fmt.Println("userLoginRequestDTO", userLoginRequestDTO)
	if err := c.ShouldBindJSON(&userLoginRequestDTO); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	token, authUser, err := h.authService.LoginUser(userLoginRequestDTO)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário logado com sucesso.",
		gin.H{
			"token": token,
			"user":  authUser,
		})
}

// @Summary Envia código de autenticação por email
// @Description Usado para autenticação de dois fatores ou verificação de email. Envia um código temporário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.EmailAuthRequestDTO true "Email para o qual enviar o código"
// @Success 200 {object} dto.CodeResponseDTO "Código enviado com sucesso (retorna o código ou status)"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida"
// @Router /auth/code [patch]
func (h *AuthHandler) SendCode(c *gin.Context) {

	var emailAuthRequestDTO dto.EmailAuthRequestDTO
	if err := c.ShouldBindJSON(&emailAuthRequestDTO); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	codeResponseDTO, err := h.authService.SendCodeToEmail(emailAuthRequestDTO)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Código enviado com sucesso.", codeResponseDTO)

}

// @Summary Valida código de autenticação
// @Description Valida o código temporário (enviado por email) e retorna um token de sessão e dados do usuário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.InputCodeDto true "Email e código para validação"
// @Success 200 {object} utils.SuccessValidateCodeResponse "Código validado com sucesso (retorna token e usuário)"
// @Failure 400 {object} utils.ErrorResponse "Código inválido (ou outros erros de requisição)"
// @Router /auth/validate [post]
func (h *AuthHandler) ValidateCode(c *gin.Context) {
	var inputCodeDto dto.InputCodeDto
	if err := c.ShouldBindJSON(&inputCodeDto); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	token, authUser, err := h.authService.ValidateUserCode(inputCodeDto)
	if err != nil {
		utils.SendErrorResponse(c, "Código inválido.", http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Código validado com sucesso.",
		gin.H{
			"token": token,
			"user": gin.H{
				"_id":              authUser.ID,
				"name":             authUser.Name,
				"email":            authUser.Email,
				"role":             authUser.Role,
				"two_factor":       authUser.TwoFactor,
				"profile_image_id": authUser.ProfileImageID,
			},
		})
}

// @Summary Cria o primeiro usuário Administrador
// @Description Endpoint para configuração inicial. Cria o primeiro usuário administrador no banco de dados se ele ainda não existir.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseNoData "Resposta de sucesso"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /auth/adm [post]
func (h *AuthHandler) FirstLoginAdmin(c *gin.Context) {
	err := h.authService.FirstLoginAdmin()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Usuário inicial criado com sucesso.", http.StatusOK)
}

// @Summary Envia email de recuperação de senha
// @Description Inicia o fluxo de "esqueci minha senha" enviando um email com código/token para o usuário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.ForgotPasswordRequestDTO true "Email do usuário que esqueceu a senha"
// @Success 200 {object} utils.SuccessResponseNoData "Email enviado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida (e.g., email mal formatado ou não encontrado)"
// @Router /auth/email [post]
func (h *AuthHandler) SendEmailForgotPassword(c *gin.Context) {
	var email dto.ForgotPasswordRequestDTO
	if err := c.ShouldBindJSON(&email); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Recebida solicitação de redefinição de senha para o DTO: %+v\n", email)

	err := h.authService.SendEmailForgotPassword(email)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendSuccessResponse(c, "Email enviado com sucesso.", nil)
}

func (h *AuthHandler) ChangePasswordUnlogged(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendErrorResponse(c, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	var updatedPasswordByNewPassword dto.UpdatedPasswordByNewPassword
	if err := c.ShouldBindJSON(&updatedPasswordByNewPassword); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	err := h.authService.ChangePasswordUnlogged(updatedPasswordByNewPassword, id)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Senha atualizada com sucesso.", nil)
}

// @Summary Reseta a senha do usuário
// @Description Altera a senha do usuário usando um token de reset (obtido no fluxo "esqueci minha senha").
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.ResetPasswordDTO true "Token de reset e nova senha"
// @Success 200 {object} utils.SuccessResponseNoData "Senha alterada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Dados inválidos (e.g., senha fraca)"
// @Failure 401 {object} utils.ErrorResponse "Token inválido ou expirado"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Dados inválidos: token e nova senha são obrigatórios.", http.StatusBadRequest)
		return
	}

	// Chama o novo método do serviço
	err := h.authService.ResetPassword(req)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.SendSuccessResponse(c, "Senha alterada com sucesso.", nil)
}

// @Summary Altera a senha do usuário logado
// @Description Permite que um usuário (paciente ou enfermeiro) logado altere sua própria senha. Requer autenticação JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.ChangePasswordBothRequestDTO true "Dados da nova senha"
// @Success 200 {object} utils.SuccessResponseNoData "Senha ou configurações atualizadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Paciente ou Enfermeiro)"
// @Router /auth/logged/password [patch]
func (h *AuthHandler) ChangePasswordLogged(c *gin.Context) {
	userId := utils.GetUserId(c)

	var changePasswordBothRequestDTO dto.ChangePasswordBothRequestDTO
	if err := c.ShouldBindJSON(&changePasswordBothRequestDTO); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	// --- MUDANÇA AQUI ---
	// 1. Capturar o novo 'bool' retornado
	passwordWasChanged, err := h.authService.ChangePasswordLogged(changePasswordBothRequestDTO, userId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Definir a mensagem com base no que mudou
	message := "Configurações de segurança atualizadas com sucesso."
	if passwordWasChanged {
		message = "Senha atualizada com sucesso."
	}

	// 3. Usar a nova mensagem
	utils.SendSuccessResponse(
		c,
		message, // <-- usa a variável
		gin.H{
			// O token não muda, enviar isso é desnecessário e pode confundir
			"token": "senha atualizada",
		},
	)
	// --- FIM DA MUDANÇA ---
}

// @Summary Valida um token de reset de senha
// @Description Verifica se um token (enviado por email) é válido e não expirou, antes de permitir a troca da senha.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.ValidateTokenDTO true "Token de reset a ser validado"
// @Success 200 {object} utils.SuccessResponseNoData "Token válido"
// @Failure 400 {object} utils.ErrorResponse "Token é obrigatório"
// @Failure 401 {object} utils.ErrorResponse "Token inválido ou expirado"
// @Router /auth/validate-token [post]
func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	var req dto.ValidateTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Token é obrigatório", http.StatusBadRequest)
		return
	}

	err := h.authService.ValidateToken(req.Token)
	if err != nil {
		utils.SendErrorResponse(c, "Token inválido ou expirado", http.StatusUnauthorized)
		return
	}

	utils.SendSuccessResponse(c, "Token válido", nil)
}
