package repository

import (
    "os"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/account"
    "github.com/stripe/stripe-go/v76/accountlink"
    "github.com/stripe/stripe-go/v76/transfer"
    "github.com/stripe/stripe-go/v76/paymentintent"
    "fmt"
)

type StripeRepository interface {
    CreateExpressAccount(email string) (string, error)
    CreateAccountLink(accountId string) (string, error)
    CreateTransfer(amountInCents int64, destinationAccountId string, sourceTransactionId string) (*stripe.Transfer, error)
}

type stripeRepository struct {
}

// Construtor
func NewStripeRepository() StripeRepository {
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
    return &stripeRepository{}
}

func (r *stripeRepository) CreateExpressAccount(email string) (string, error) {
    params := &stripe.AccountParams{
        Type:    stripe.String(string(stripe.AccountTypeExpress)),
        Email:   stripe.String(email),
        Country: stripe.String("BR"),
        Capabilities: &stripe.AccountCapabilitiesParams{
            CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{Requested: stripe.Bool(true)},
            Transfers:    &stripe.AccountCapabilitiesTransfersParams{Requested: stripe.Bool(true)},
        },
    }
    
    newAccount, err := account.New(params)
    if err != nil {
        return "", err
    }
    
    return newAccount.ID, nil
}

// Cria um link de onboarding de uso único
func (r *stripeRepository) CreateAccountLink(accountId string) (string, error) {
    // Você pode pegar essas URLs do .env também
    refreshURL := os.Getenv("STRIPE_REFRESH_URL")
    returnURL  := os.Getenv("STRIPE_RETURN_URL") 

    linkParams := &stripe.AccountLinkParams{
        Account:    stripe.String(accountId),
        RefreshURL: stripe.String(refreshURL),
        ReturnURL:  stripe.String(returnURL),
        Type:       stripe.String("account_onboarding"),
    }
    
    link, err := accountlink.New(linkParams)
    if err != nil {
        return "", err
    }
    
    return link.URL, nil // Retorna a URL completa
}


func (r *stripeRepository) CreateTransfer(amountInCents int64, destinationAccountId string, paymentIntentId string) (*stripe.Transfer, error) {
    
    // 1. Buscar o PaymentIntent no Stripe usando o ID (pi_...)
    // (Esta parte estava correta)
    pi, err := paymentintent.Get(paymentIntentId, nil)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar payment intent no Stripe: %w", err)
    }

    // --- MUDANÇA (A CORREÇÃO FINAL) ---
    // 2. O campo correto é 'LatestCharge' (singular)
    
    // Verificação de segurança (Boa Prática) para o ponteiro
    if pi.LatestCharge == nil {
        return nil, fmt.Errorf("payment intent %s não possui uma cobrança (charge) associada (status: %s)", pi.ID, pi.Status)
    }

    // O ID da cobrança (ch_...) está dentro do objeto 'LatestCharge'
    latestChargeID := pi.LatestCharge.ID
    if latestChargeID == "" {
        return nil, fmt.Errorf("ID da cobrança (charge) associada está vazio")
    }

    // 3. USAR O 'latestChargeID' (o 'ch_...') como SourceTransaction
    params := &stripe.TransferParams{
        Amount:      stripe.Int64(amountInCents),
        Currency:    stripe.String(string(stripe.CurrencyBRL)),
        Destination: stripe.String(destinationAccountId),
        
        // Aqui usamos a variável correta que acabamos de pegar
        SourceTransaction: stripe.String(latestChargeID), 
    }

    t, err := transfer.New(params)
    if err != nil {
        return nil, err
    }

    return t, nil
}