package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var nonNumericRegex = regexp.MustCompile(`[^\d]`)

var emailPattern = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func EmailRegex(email string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))

	if !emailPattern.MatchString(normalized) {
		return "", errors.New("Formato de e-mail inválido")
	}

	return normalized, nil
}
func ValidateCPF(cpf string) (string, error) {
	normalized := nonNumericRegex.ReplaceAllString(cpf, "")

	if len(normalized) != 11 {
		return "", errors.New("o CPF deve conter exatamente 11 dígitos")
	}

	return normalized, nil
}

func ValidateCEP(cep string) (string, error) {
	normalized := nonNumericRegex.ReplaceAllString(cep, "")

	if len(normalized) != 8 {
		return "", errors.New("o CEP deve conter exatamente 8 dígitos")
	}

	return normalized, nil
}

func ValidatePhone(phone string) (string, error) {
	normalized := nonNumericRegex.ReplaceAllString(phone, "")

	if len(normalized) < 10 || len(normalized) > 11 {
		return "", errors.New("o telefone deve conter 10 ou 11 dígitos (incluindo DDD)")
	}

	return normalized, nil
}

func ValidatePasswordRegex(password string) error {
	if len(password) < 8 {
		return errors.New("a senha deve ter no mínimo 8 caracteres")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("a senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	}

	return nil
}

func ValidateCoren(coren string) (string, error) {
	upperCoren := strings.ToUpper(strings.TrimSpace(coren))

	corenPattern := regexp.MustCompile(`^COREN[- ]*([A-Z]{2})[- ]*(\d{6})$`)

	matches := corenPattern.FindStringSubmatch(upperCoren)

	if len(matches) != 3 {
		return "", errors.New("formato de COREN inválido. Use 'COREN-UF 123456'.")
	}

	normalized := matches[1] + matches[2]
	return normalized, nil
}