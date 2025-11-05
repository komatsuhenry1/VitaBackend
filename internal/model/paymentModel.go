package model

// PaymentIntentRequest é o que o frontend (Next.js) nos envia.
// (O frontend lê o valor do sessionStorage)
type PaymentIntentRequest struct {
    Value float64 `json:"value"`
}

// PaymentIntentResponse é o que nosso backend retorna para o frontend.
// (O frontend usa isso para inicializar o Stripe Elements)
type PaymentIntentResponse struct {
    ClientSecret string `json:"client_secret"`
}