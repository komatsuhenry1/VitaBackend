package payment

import (
	"medassist/internal/model"
	"medassist/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService PaymentService
}

func NewPaymentHandler(paymentService PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// @Summary Cria uma Intenção de Pagamento (Stripe)
// @Description Gera um 'client_secret' do Stripe para o paciente logado poder realizar um pagamento (ex: adicionar fundos à carteira). Requer autenticação de Paciente.
// @Tags Payment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body model.PaymentIntentRequest true "Valor do pagamento a ser criado"
// @Success 200 {object} utils.SuccessPaymentIntentResponse "Intenção de pagamento criada com sucesso"
// @Failure 400 {object} utils.ErrorResponse "Requisição inválida"
// @Failure 401 {object} utils.ErrorResponse "Não autorizado (Token JWT inválido ou ausente)"
// @Failure 403 {object} utils.ErrorResponse "Proibido (Usuário não é Paciente)"
// @Failure 500 {object} utils.ErrorResponse "Não foi possível criar a intenção de pagamento"
// @Router /payment/create-intent [post]
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
    var req model.PaymentIntentRequest

    // 1. Bind do JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
        return
    }

    // 2. Obter o ID do paciente logado (do middleware de autenticação)
    //    Isso é FUNDAMENTAL para vincular o pagamento ao paciente correto.
    patientID := utils.GetUserId(c)

    // 3. Chamar o serviço de pagamento
    clientSecret, err := h.paymentService.CreatePaymentIntent(patientID, req.Value)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar a intenção de pagamento"})
        return
    }

    // 4. Retornar o Client Secret para o frontend
    c.JSON(http.StatusOK, model.PaymentIntentResponse{
        ClientSecret: clientSecret,
    })
}