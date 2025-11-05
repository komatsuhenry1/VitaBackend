package payment

import (
	"fmt"
	"math"
	"medassist/internal/repository"
	"os"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"go.mongodb.org/mongo-driver/bson"
)

// Interface para o serviço
type PaymentService interface {
	CreatePaymentIntent(patientID string, value float64) (string, error)
}

type paymentService struct {
	paymentRepository repository.PaymentRepository // Para buscar/salvar o ID do cliente Stripe
	userRepository    repository.UserRepository    // Para buscar/salvar o ID do cliente Stripe
}

// (Estou assumindo que você tem um repositório para buscar pacientes)
func NewPaymentService(paymentRepository repository.PaymentRepository, userRepository repository.UserRepository) PaymentService {
	// Configura a chave secreta do Stripe (NUNCA exponha no código)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &paymentService{paymentRepository: paymentRepository, userRepository: userRepository}
}

func (s *paymentService) CreatePaymentIntent(patientID string, value float64) (string, error) {

	// 1. OBTER/CRIAR O CLIENTE NO STRIPE
	patient, err := s.userRepository.FindUserById(patientID)
	if err != nil {
		return "", fmt.Errorf("paciente não encontrado: %w", err)
	}

	stripeCustomerID := patient.GatewayCustomerID

	fmt.Println("stripeCostumerID", stripeCustomerID)

	if stripeCustomerID == "" {
		// Se não tiver, criamos um novo Cliente no Stripe
		customerParams := &stripe.CustomerParams{
			Name:  stripe.String(patient.Name),
			Email: stripe.String(patient.Email),
			Phone: stripe.String(patient.Phone),
			Address: &stripe.AddressParams{
				Line1:      stripe.String(patient.Address), // ou os campos reais se você tiver
				City:       stripe.String(patient.City),	
				State:      stripe.String(patient.UF),
				PostalCode: stripe.String(patient.CEP),
				Country:    stripe.String("Brasil"),
			}}

		customerParams.AddMetadata("app_patient_id", patientID)

		newCustomer, err := customer.New(customerParams)
		if err != nil {
			return "", fmt.Errorf("erro ao criar cliente no Stripe: %w", err)
		}

		stripeCustomerID = newCustomer.ID

		updates := bson.M{
			"gateway_customer_id": stripeCustomerID,
		}

		if _, err := s.userRepository.UpdateUserFields(patientID, updates); err != nil {
			return "", fmt.Errorf("erro ao salvar o ID do cliente Stripe no banco: %w", err)
		}
	}

	fmt.Println("stripeCostumerID", stripeCustomerID)

	amountInCents := int64(math.Round(value * 100))

    // 3. CRIAR A INTENÇÃO DE PAGAMENTO (PAYMENT INTENT)
    params := &stripe.PaymentIntentParams{
        Amount:           stripe.Int64(amountInCents),
        Currency:         stripe.String(string(stripe.CurrencyBRL)),
        Customer:         stripe.String(stripeCustomerID),
        SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),

        // 👇 MUDANÇA 1: Removido o 'CaptureMethodManual' (agora é automático)

        // 👇 MUDANÇA 2: Adicionado para limitar ao cartão e evitar o erro do Apple Pay
        PaymentMethodTypes: []*string{
            stripe.String("card"),
        },
    }

    // Código antigo que foi alterado:
    // CaptureMethod:    stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
    
    pi, err := paymentintent.New(params)
    if err != nil {
        return "", fmt.Errorf("erro ao criar PaymentIntent no Stripe: %w", err)
    }

    // 4. RETORNAR O "CLIENT SECRET"
    return pi.ClientSecret, nil
}
