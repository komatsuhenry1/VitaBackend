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
// (Estou assumindo que você injeta os handlers em algum lugar)
// var PaymentHandler = NewPaymentHandler(services.NewPaymentService(), services.NewPatientService())


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