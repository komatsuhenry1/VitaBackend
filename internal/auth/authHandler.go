package auth

import (
	"fmt"
	"medassist/internal/auth/dto"
	"medassist/utils"
	"net/http"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) UserRegister(c *gin.Context) {
	// 1. Criar o DTO e preenchê-lo com os dados do formulário
	var userRequestDTO dto.UserRegisterRequestDTO
	userRequestDTO.Name = c.PostForm("name")
	userRequestDTO.Email = c.PostForm("email")
	userRequestDTO.Phone = c.PostForm("phone")
	userRequestDTO.Address = c.PostForm("address")
	userRequestDTO.Cpf = c.PostForm("cpf")
	userRequestDTO.Password = c.PostForm("password")

	// 2. Processar a parte de multipart/form-data para pegar os arquivos
	form, err := c.MultipartForm()
	if err != nil {
		// Se não for um multipart form, pode ser que nenhum arquivo foi enviado, o que é ok.
		// Vamos tratar o `form` como nulo ou vazio na camada de serviço.
		// Se o erro for outro (ex: form mal formatado), ele será pego.
		// Para simplificar, vamos assumir que o erro aqui indica um problema.
		utils.SendErrorResponse(c, "Erro ao processar formulário multipart: "+err.Error(), http.StatusBadRequest)
		return
	}
	files := form.File // Pega o mapa de todos os arquivos

	// 3. Chamar o serviço com o DTO e os arquivos
	createdUser, err := h.authService.UserRegister(userRequestDTO, files)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "usuário criado com sucesso", gin.H{"user": createdUser})
}

func (h *AuthHandler) NurseRegister(c *gin.Context) {

	yearsExpStr := c.PostForm("years_experience")
	fmt.Println("yearsExpStr: ", yearsExpStr)
	fmt.Printf("yearsExpStr type: %T\n", yearsExpStr)
	yearsExp, err := strconv.Atoi(yearsExpStr)
	if err != nil {
		// Se a conversão falhar, retorne um erro claro para o frontend
		utils.SendErrorResponse(c, "Formato inválido para 'anos de experiência'. Esperado um número.", http.StatusBadRequest)
		return // Interrompe a execução
	}

	var nurseRequestDTO dto.NurseRegisterRequestDTO
	nurseRequestDTO.Name = c.PostForm("name")
	nurseRequestDTO.Email = c.PostForm("email")
	nurseRequestDTO.Phone = c.PostForm("phone")
	nurseRequestDTO.Address = c.PostForm("address")
	nurseRequestDTO.Cpf = c.PostForm("cpf")
	nurseRequestDTO.PixKey = c.PostForm("pix_key")
	nurseRequestDTO.Password = c.PostForm("password")
	nurseRequestDTO.LicenseNumber = c.PostForm("license_number")
	nurseRequestDTO.Specialization = c.PostForm("specialization")
	nurseRequestDTO.Shift = c.PostForm("shift")
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

	utils.SendSuccessResponse(c, "usuário criado com sucesso", gin.H{"nurse": createdNurse})
}

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
			"user": authUser,
		})
}

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
				"id":               authUser.ID,
				"name":             authUser.Name,
				"email":            authUser.Email,
				"role":             authUser.Role,
				"two_factor":       authUser.TwoFactor,
				"profile_image_id": authUser.ProfileImageID,
			},
		})
}

func (h *AuthHandler) FirstLoginAdmin(c *gin.Context) {
	err := h.authService.FirstLoginAdmin()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, "Usuário inicial criado com sucesso.", http.StatusOK)
}

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

// Use este handler em vez de ChangePasswordUnlogged
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

func (h *AuthHandler) ChangePasswordLogged(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}
	id, ok := claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId inválido no token"})
		return
	}

	var changePasswordBothRequestDTO dto.ChangePasswordBothRequestDTO
	if err := c.ShouldBindJSON(&changePasswordBothRequestDTO); err != nil {
		utils.SendErrorResponse(c, "Requisição inválida", http.StatusBadRequest)
		return
	}

	err := h.authService.ChangePasswordLogged(changePasswordBothRequestDTO, id)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(
		c,
		"Senha atualizada com sucesso.",
		gin.H{
			"token": "senha atualizada",
		},
	)
}

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

	utils.SendSuccessResponse(c, "Token válido", http.StatusOK)
}
