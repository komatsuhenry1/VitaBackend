package repository

import (
    "os"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/account"
    "github.com/stripe/stripe-go/v76/accountlink"
)

type StripeRepository interface {
    CreateExpressAccount(email string) (string, error)
    CreateAccountLink(accountId string) (string, error)
}

// Implementação
type stripeRepository struct {
    // Aqui você pode adicionar configs, se necessário
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