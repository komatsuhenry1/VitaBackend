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
	CreatePaymentIntent(patientID string, value float64, visitId string) (string, error)
}

type paymentService struct {
	paymentRepository repository.PaymentRepository 
	userRepository    repository.UserRepository
	visitRepository   repository.VisitRepository
}

func NewPaymentService(paymentRepository repository.PaymentRepository, userRepository repository.UserRepository, visitRepository repository.VisitRepository) PaymentService {
	// Configura a chave secreta do Stripe (NUNCA exponha no código)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &paymentService{paymentRepository: paymentRepository, userRepository: userRepository, visitRepository: visitRepository}
}

func (s *paymentService) CreatePaymentIntent(patientID string, value float64, visitId string) (string, error) {

	patient, err := s.userRepository.FindUserById(patientID)
	if err != nil {
		return "", fmt.Errorf("paciente não encontrado: %w", err)
	}

	stripeCustomerID := patient.GatewayCustomerID

	//cria o cliente no stripe

	if stripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Name:  stripe.String(patient.Name),
			Email: stripe.String(patient.Email),
			Phone: stripe.String(patient.Phone),
			Address: &stripe.AddressParams{
				Line1:      stripe.String(patient.Address),
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

		//salva  o costumer id no meu banco

		stripeCustomerID = newCustomer.ID

		updates := bson.M{
			"gateway_customer_id": stripeCustomerID,
		}

		if _, err := s.userRepository.UpdateUserFields(patientID, updates); err != nil {
			return "", fmt.Errorf("erro ao salvar o ID do cliente Stripe no banco: %w", err)
		}
	}

	amountInCents := int64(math.Round(value * 100))

	// CRIAR A INTENÇÃO DE PAGAMENTO (PAYMENT INTENT)
	params := &stripe.PaymentIntentParams{
		Amount:           stripe.Int64(amountInCents),                                           // quantia
		Currency:         stripe.String(string(stripe.CurrencyBRL)),                             // moeda
		Customer:         stripe.String(stripeCustomerID),                                       // cliente
		SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)), // serve para indicar que o método de pagamento poderá ser reutilizado no futuro, sem precisar da interação do usuário (ou seja, "off-session"

		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", fmt.Errorf("erro ao criar PaymentIntent no Stripe: %w", err)
	}

	// b, _ := json.MarshalIndent(pi, "", "  ")
	// fmt.Println(string(b))

	//salvar o paymentIntent.ID
	visitUpdates := bson.M{
		"payment_intent_id": pi.ID,
	}

	_, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdates)
	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}
