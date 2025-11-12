package admin

import (
	"medassist/internal/admin/dto"
	"medassist/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"strings"
	"fmt"
)

type AdminHandler struct {
	adminService AdminService
}

func NewAdminHandler(adminService AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// @Summary Dados do Dashboard Administrativo
// @Description Retorna as principais métricas e KPIs para a tela de dashboard do administrador. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessDashboardResponse "Dados de dashboard carregados com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao carregar dados"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/dashboard [get]
func (h *AdminHandler) AdminDashboard(c *gin.Context) {
	dashboardData, err := h.adminService.GetDashboardData()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Dados de dashboard carregados com sucesso", dashboardData)
}

// @Summary Obtém documentos do enfermeiro para análise
// @Description Retorna uma lista de documentos (com URLs de download) de um enfermeiro específico para aprovação. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Enfermeiro (Nurse ID)"
// @Success 200 {object} utils.SuccessDocumentsResponse "Documentos retornados com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao buscar documentos"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/documents/{id} [get]
func (h *AdminHandler) GetDocuments(c *gin.Context) {
	nurseId := c.Param("id")

	documents, err := h.adminService.GetNurseDocumentsToAnalisys(nurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Documentos retornados com sucesso.", documents)
}

// @Summary Aprova o cadastro de um enfermeiro
// @Description Altera o status do enfermeiro para 'verificado' (verification_seal: true). Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Enfermeiro a ser aprovado"
// @Success 200 {object} utils.SuccessResponseString "Enfermeiro aprovado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro na aprovação"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/approve/{id} [patch]
func (h *AdminHandler) ApproveNurseRegister(c *gin.Context) {
	approvedNurseId := c.Param("id")

	data, err := h.adminService.ApproveNurseRegister(approvedNurseId)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Enfermeiro(a) aprovado(a) com sucesso.", data)
}

// @Summary Baixa um arquivo do GridFS
// @Description Faz o download de um arquivo (como um documento de enfermeiro) com base no seu ObjectID do GridFS.
// @Tags Admin
// @Produce application/octet-stream
// @Param id path string true "ID do Arquivo (File ID)"
// @Success 200 {file} file "O arquivo para download"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou Arquivo não encontrado"
// @Failure 500 {object} utils.ErrorResponse "Erro ao enviar o arquivo"
// @Router /admin/download/{id} [get]
func (h *AdminHandler) DownloadFile(c *gin.Context) {
	// 1. Pega o ID do arquivo a partir do parâmetro da URL.
	fileIDHex := c.Param("id")
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		utils.SendErrorResponse(c, "ID do arquivo inválido", http.StatusBadRequest)
		return
	}

	// 2. Chama a camada de serviço para buscar o stream do arquivo.
	downloadStream, err := h.adminService.GetFileStream(fileID)
	if err != nil {
		// O serviço retornará um erro se o arquivo não for encontrado.
		utils.SendErrorResponse(c, "Arquivo não encontrado", http.StatusBadRequest)
		return
	}
	// Garante que o stream será fechado no final da função.
	defer downloadStream.Close()

	// 3. Pega os metadados do arquivo para configurar a resposta.
	fileInfo := downloadStream.GetFile()

	// 4. Define os Headers HTTP. Content-Type padrão se não houver metadata específica.
	c.Header("Content-Type", "application/octet-stream")
	// Content-Length informa o tamanho do arquivo.
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Length, 10))
	// Content-Disposition com "attachment" força o navegador a abrir a caixa de "Salvar Como...".
	c.Header("Content-Disposition", "attachment; filename=\""+fileInfo.Name+"\"")

	// 5. Copia o conteúdo do stream do GridFS diretamente para o corpo da resposta HTTP.
	// Isso é muito eficiente em termos de memória, pois o arquivo não é totalmente carregado no servidor.
	if _, err := io.Copy(c.Writer, downloadStream); err != nil {
		utils.SendErrorResponse(c, "Erro ao enviar o arquivo", http.StatusInternalServerError)
		return
	}
}

// @Summary Rejeita o cadastro de um enfermeiro
// @Description Envia um email ao enfermeiro com o motivo da rejeição. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Enfermeiro a ser rejeitado"
// @Param payload body dto.RejectDescription true "Objeto contendo a descrição/motivo da rejeição"
// @Success 200 {object} utils.SuccessResponseString "Enfermeiro rejeitado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido, corpo da requisição inválido ou erro na rejeição"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/reject/{id} [post]
func (h *AdminHandler) RejectNurseRegister(c *gin.Context) {
	rejectedNurseId := c.Param("id")

	var rejectDescription dto.RejectDescription
	if err := c.ShouldBindJSON(&rejectDescription); err != nil {
		utils.SendErrorResponse(c, "Corpo da requisição inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.adminService.RejectNurseRegister(rejectedNurseId, rejectDescription)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Enfermeiro(a) rejeitado com sucesso.", data)
}

// @Summary Listas de gerenciamento (Usuários, Enfermeiros, Visitas)
// @Description Retorna listas completas de usuários, enfermeiros e visitas para gerenciamento do admin. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessUserListsResponse "Listas retornadas com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Erro ao buscar listas"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/users [get]
func (h *AdminHandler) UsersManagement(c *gin.Context) {

	userLists, err := h.adminService.UserLists()
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Listas de usuários retornadas com sucesso.", userLists)
}

// @Summary Atualiza um usuário (Admin)
// @Description Permite ao administrador atualizar campos de um usuário (Paciente ou Enfermeiro). Requer autenticação de Admin.
// @Description Campos protegidos (id, created_at, updated_at) não podem ser atualizados.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Usuário a ser atualizado"
// @Param payload body object true "Campos para atualizar (JSON arbitrário). Ex: {\"hidden\": true, \"name\": \"Novo Nome\"}"
// @Success 200 {object} utils.SuccessUserTypeResponse "Usuário atualizado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido, ID não encontrado ou campo protegido"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/user/{id} [patch]
func (h *AdminHandler) UpdateUser(c *gin.Context){
	userId := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	protectedFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
	}

	for key := range updates {
		if protectedFields[strings.ToLower(key)] {
			utils.SendErrorResponse(c, fmt.Sprintf("Campo(s) %s não pode ser atualizado.", key), http.StatusBadRequest)
			return
		}
	}
	
	user, err := h.adminService.UpdateUser(userId, updates)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Usuário atualizado com sucesso.", user)

}

// @Summary Deleta um usuário (Admin)
// @Description Permite ao administrador deletar um usuário (Paciente ou Enfermeiro) permanentemente. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID do Usuário a ser deletado"
// @Success 200 {object} utils.SuccessResponseNoData "Usuário deletado com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao deletar"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/user/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context){
	userId := c.Param("id")

	err := h.adminService.DeleteNurseOrUser(userId)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "Usuário deletado com sucesso.", http.StatusOK)
}

// @Summary Atualiza uma visita (Admin)
// @Description Permite ao administrador atualizar campos de uma visita. Requer autenticação de Admin.
// @Description Campos protegidos (id, created_at, updated_at) não podem ser atualizados.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser atualizada"
// @Param payload body object true "Campos para atualizar (JSON arbitrário). Ex: {\"status\": \"COMPLETED\", \"nurse_id\": \"...\"}"
// @Success 200 {object} utils.SuccessVisitTypeResponse "Visita atualizada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "JSON inválido, ID não encontrado ou campo protegido"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/visit/{id} [patch]
func (h *AdminHandler) UpdateVisit(c *gin.Context) {
	visitId := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendErrorResponse(c, "JSON inválido", http.StatusBadRequest)
		return
	}

	protectedFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
	}

	for key := range updates {
		if protectedFields[strings.ToLower(key)] {
			utils.SendErrorResponse(c, fmt.Sprintf("Campo(s) %s não pode ser atualizado.", key), http.StatusBadRequest)
			return
		}
	}
	
	visit, err := h.adminService.UpdateVisit(visitId, updates)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Usuário atualizado com sucesso.", visit)
}

// @Summary Deleta uma visita (Admin)
// @Description Permite ao administrador deletar uma visita permanentemente. Requer autenticação de Admin.
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID da Visita a ser deletada"
// @Success 200 {object} utils.SuccessResponseNoData "Visita deletada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "ID inválido ou erro ao deletar"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Administrador)"
// @Router /admin/visit/{id} [delete]
func (h *AdminHandler) DeleteVisit(c *gin.Context){
	visitId := c.Param("id")

	err := h.adminService.DeleteVisit(visitId)
	if err != nil{
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	utils.SendSuccessResponse(c, "Visita deletada com sucesso.", http.StatusOK)
}
