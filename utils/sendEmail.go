package utils

import (
	"fmt"
	"os"
	"medassist/internal/user/dto"

	"gopkg.in/gomail.v2"
)

func SendEmailNurseRegister(email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	m.SetHeader("Subject", "🔑 Análise de cadastro - Bem-vindo à Plataforma")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<title>Senha de Acesso</title>
		<style>
			body {
				background-color: #f9f9f9;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				color: #333333;
				padding: 0;
				margin: 0;
			}
			.container {
				max-width: 600px;
				margin: 40px auto;
				background-color: #ffffff;
				border-radius: 10px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
				padding: 30px 40px;
			}
			h2 {
				color: #1E88E5;
				text-align: center;
			}
			p {
				line-height: 1.6;
				font-size: 15px;
			}
			.code-box {
				background-color: #f1f1f1;
				border-radius: 6px;
				padding: 10px;
				font-family: monospace;
				font-size: 16px;
				color: #333333;
				margin: 15px 0;
				text-align: center;
				font-weight: bold;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #999999;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🔑 Sua conta está em analise para ser cadastrada no sistema como enfermeiro(a).</h2>
			<p>Olá,</p>
			<p><strong>E-mail cadastrado:</strong></p>
			<div class="code-box">%s</div>

			<p><strong>Sua conta está em analise para ser cadastrada no sistema como enfermeiro(a).</strong></p>

			<p>⚠️ Caso necessário, você pode alterar sua senha assim que fizer o primeiro login.</p>

			<div class="footer">
				<p>Se você não solicitou esta conta, apenas ignore este e-mail.</p>
				<p>Este é um e-mail automático. Por favor, não responda.</p>
			</div>
		</div>
	</body>
	</html>
	`, email)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailUserRegister(email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	m.SetHeader("Subject", "🔑 Cadastro de conta - Bem-vindo à Plataforma")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<title>Senha de Acesso</title>
		<style>
			body {
				background-color: #f9f9f9;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				color: #333333;
				padding: 0;
				margin: 0;
			}
			.container {
				max-width: 600px;
				margin: 40px auto;
				background-color: #ffffff;
				border-radius: 10px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
				padding: 30px 40px;
			}
			h2 {
				color: #1E88E5;
				text-align: center;
			}
			p {
				line-height: 1.6;
				font-size: 15px;
			}
			.code-box {
				background-color: #f1f1f1;
				border-radius: 6px;
				padding: 10px;
				font-family: monospace;
				font-size: 16px;
				color: #333333;
				margin: 15px 0;
				text-align: center;
				font-weight: bold;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #999999;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🔑 Cadastro de conta</h2>
			<p>Olá,</p>
			<p>Seja bem-vindo! Sua conta foi criada com sucesso.</p>
			<p><strong>E-mail cadastrado:</strong></p>
			<div class="code-box">%s</div>

			<p>⚠️ Caso necessário, você pode alterar sua senha assim que fizer o primeiro login.</p>

			<div class="footer">
				<p>Se você não solicitou esta conta, apenas ignore este e-mail.</p>
				<p>Este é um e-mail automático. Por favor, não responda.</p>
			</div>
		</div>
	</body>
	</html>
	`, email)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendAuthCode(email string, code int) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	m.SetHeader("Subject", "🔑 Código de Acesso - Bem-vindo à Plataforma")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<title>Senha de Acesso</title>
		<style>
			body {
				background-color: #f9f9f9;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				color: #333333;
				padding: 0;
				margin: 0;
			}
			.container {
				max-width: 600px;
				margin: 40px auto;
				background-color: #ffffff;
				border-radius: 10px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
				padding: 30px 40px;
			}
			h2 {
				color: #1E88E5;
				text-align: center;
			}
			p {
				line-height: 1.6;
				font-size: 15px;
			}
			.code-box {
				background-color: #f1f1f1;
				border-radius: 6px;
				padding: 10px;
				font-family: monospace;
				font-size: 16px;
				color: #333333;
				margin: 15px 0;
				text-align: center;
				font-weight: bold;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #999999;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🔑 Seu código de acesso</h2>

			<p><strong>Code:</strong></p>
			<div class="code-box">%s</div>

			<p>⚠️ Por motivos de segurança, recomendamos que você altere sua senha no menu de segurança.</p>

			<div class="footer">
				<p>Se você não solicitou esta conta, apenas ignore este e-mail.</p>
				<p>Este é um e-mail automático. Por favor, não responda.</p>
			</div>
		</div>
	</body>
	</html>
	`, code)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailForAdmin(email string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	m.SetHeader("Subject", "🔑 Sua senha de acesso - Bem-vindo à Plataforma")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<title>Senha de Acesso</title>
		<style>
			body {
				background-color: #f9f9f9;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				color: #333333;
				padding: 0;
				margin: 0;
			}
			.container {
				max-width: 600px;
				margin: 40px auto;
				background-color: #ffffff;
				border-radius: 10px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
				padding: 30px 40px;
			}
			h2 {
				color: #1E88E5;
				text-align: center;
			}
			p {
				line-height: 1.6;
				font-size: 15px;
			}
			.code-box {
				background-color: #f1f1f1;
				border-radius: 6px;
				padding: 10px;
				font-family: monospace;
				font-size: 16px;
				color: #333333;
				margin: 15px 0;
				text-align: center;
				font-weight: bold;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #999999;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🔑 Sua Senha de Acesso (ADMINISTRADOR)</h2>
			<p>Olá,</p>
			<p>Seja bem-vindo! Sua conta de administrador foi criada com sucesso.</p>
			<p><strong>E-mail cadastrado:</strong></p>
			<div class="code-box">%s</div><br />


			<p><strong>Sua senha de acesso é a mesma que solicitou a nossa equipe na criação da conta.</strong></p>

			<div class="footer">
				<p>Se você não solicitou esta conta, apenas ignore este e-mail.</p>
				<p>Este é um e-mail automático. Por favor, não responda.</p>
			</div>
		</div>
	</body>
	</html>
	`, email)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailForgotPassword(email, id, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	// Link agora inclui o token no botão
	link := os.Getenv("LOCAL_FRONTEND_URL") + "?token=" + token

	m.SetHeader("Subject", "🔐 Recuperação de senha - MEDASSIST")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
	<meta charset="UTF-8">
	<title>Recuperação de Senha - CTF ARENA</title>
	<style>
	body {
		background-color: #f9f9f9;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		color: #333333;
		padding: 0;
		margin: 0;
	}
	.container {
		max-width: 600px;
		margin: 40px auto;
		background-color: #ffffff;
		border-radius: 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		padding: 30px 40px;
	}
	h2 {
		color: #1E88E5;
		text-align: center;
	}
	p {
		line-height: 1.6;
		font-size: 15px;
	}
	.button {
		display: inline-block;
		padding: 12px 20px;
		margin: 20px 0;
		background-color: #1E88E5;
		color: #ffffff !important;
		text-decoration: none;
		border-radius: 6px;
		font-weight: 600;
		text-align: center;
	}
	.code-box {
		background-color: #f1f1f1;
		border-radius: 6px;
		padding: 10px;
		font-family: monospace;
		font-size: 14px;
		color: #333333;
		margin: 10px 0;
	}
	.footer {
		margin-top: 30px;
		font-size: 12px;
		color: #999999;
		text-align: center;
	}
	</style>
	</head>
	<body>
	<div class="container">
		<h2>🔐 Recuperação de Senha</h2>
		<p>Olá,</p>
		<p>Recebemos uma solicitação para redefinir a senha da sua conta associada ao e-mail:</p>
		<div class="code-box">%s</div>

		<p>Para criar uma nova senha, clique no botão abaixo:</p>
		<a href="%s" class="button">Redefinir Senha</a>

		<p>Se você não solicitou essa alteração, apenas ignore este e-mail. Nenhuma ação será realizada.</p>

		<div class="footer">
			<p>CTF ARENA - Este é um e-mail automático, por favor não responda.</p>
		</div>
	</div>
	</body>
	</html>
	`, email, link)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailRegistrationRejected(email, description string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "❌ Cadastro Rejeitado - MEDASSIST")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
	<meta charset="UTF-8">
	<title>Cadastro Rejeitado - MEDASSIST</title>
	<style>
	body {
		background-color: #f9f9f9;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		color: #333333;
		padding: 0;
		margin: 0;
	}
	.container {
		max-width: 600px;
		margin: 40px auto;
		background-color: #ffffff;
		border-radius: 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		padding: 30px 40px;
	}
	h2 {
		color: #E53935;
		text-align: center;
	}
	p {
		line-height: 1.6;
		font-size: 15px;
	}
	.code-box {
		background-color: #f1f1f1;
		border-radius: 6px;
		padding: 10px;
		font-family: monospace;
		font-size: 14px;
		color: #333333;
		margin: 10px 0;
	}
	.footer {
		margin-top: 30px;
		font-size: 12px;
		color: #999999;
		text-align: center;
	}
	</style>
	</head>
	<body>
	<div class="container">
		<h2>❌ Cadastro Rejeitado</h2>
		<p>Olá,</p>
		<p>Infelizmente, sua solicitação de cadastro no sistema foi rejeitada.</p>

		<p>Motivo:</p>
		<div class="code-box">%s</div>

		<p>Se você acredita que isso foi um engano, entre em contato com o suporte para mais informações.</p>

		<div class="footer">
			<p>MEDASSIST - Este é um e-mail automático, por favor não responda.</p>
		</div>
	</div>
	</body>
	</html>
	`, description)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailApprovedNurse(email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "✅ Cadastro Aprovado - MEDASSIST")

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
	<meta charset="UTF-8">
	<title>Cadastro Aprovado - MEDASSIST</title>
	<style>
	body {
		background-color: #f9f9f9;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		color: #333333;
		padding: 0;
		margin: 0;
	}
	.container {
		max-width: 600px;
		margin: 40px auto;
		background-color: #ffffff;
		border-radius: 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		padding: 30px 40px;
	}
	h2 {
		color:rgb(53, 229, 82);
		text-align: center;
	}
	p {
		line-height: 1.6;
		font-size: 15px;
	}
	.code-box {
		background-color: #f1f1f1;
		border-radius: 6px;
		padding: 10px;
		font-family: monospace;
		font-size: 14px;
		color: #333333;
		margin: 10px 0;
	}
	.footer {
		margin-top: 30px;
		font-size: 12px;
		color: #999999;
		text-align: center;
	}
	</style>
	</head>
	<body>
	<div class="container">
		<h2>Cadastro Aprovado</h2>
		<p>Olá,</p>
		<p>Sua solicitação de cadastro no sistema, foi analisada e aprovada.</p>

		<p>Se você acredita que isso foi um engano, entre em contato com o suporte para mais informações.</p>

		<div class="footer">
			<p>MEDASSIST - Este é um e-mail automático, por favor não responda.</p>
		</div>
	</div>
	</body>
	</html>
	`)

	m.SetBody("text/html", html)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendContactUsEmail(contactUsDto dto.ContactUsDTO) error {	
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", os.Getenv("EMAIL_CENTRAL_CONTACT"))



	m.SetHeader("Reply-To", contactUsDto.Email)

	m.SetHeader("Subject", fmt.Sprintf("Novo Contato: %s", contactUsDto.Subject))

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<title>Novo Contato Recebido</title>
		<style>
			body {
				background-color: #f9f9f9;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				color: #333333;
				padding: 0;
				margin: 0;
			}
			.container {
				max-width: 600px;
				margin: 40px auto;
				background-color: #ffffff;
				border-radius: 10px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
				padding: 30px 40px;
			}
			h2 {
				color: #1E88E5;
				text-align: center;
			}
			p {
				line-height: 1.6;
				font-size: 15px;
			}
			strong {
				color: #555555;
			}
			.message-box {
				background-color: #f1f1f1;
				border-left: 4px solid #1E88E5;
				border-radius: 4px;
				padding: 15px;
				margin-top: 10px;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #999999;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>📧 Novo Contato Recebido</h2>
			<p>Você recebeu uma nova mensagem através do formulário de contato.</p>
			
			<p><strong>Nome:</strong> %s</p>
			<p><strong>E-mail (para resposta):</strong> %s</p>
			<p><strong>Telefone:</strong> %s</p>
			<p><strong>Assunto:</strong> %s</p>
			
			<p><strong>Mensagem:</strong></p>
			<div class="message-box">
				%s
			</div>

			<div class="footer">
				<p>Este é um e-mail automático enviado pelo sistema.</p>
			</div>
		</div>
	</body>
	</html>
	`, contactUsDto.Name, contactUsDto.Email, contactUsDto.Phone, contactUsDto.Subject, contactUsDto.Message)

	m.SetBody("text/html", html)

	// Configuração do discador SMTP
	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	// Envio do e-mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmailVisitSolicitation(email string, patientName string, visitDate string, visitValue string, address string) error {
	// Cria a mensagem de email
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", email)

	// Tema
	m.SetHeader("Subject", "🔔 Nova Solicitação de Visita Recebida")

	// Conteúdo do email
	html := createVisitSolicitationHTML(patientName, visitDate, visitValue, address)
	m.SetBody("text/html", html)

	// Configuração do Dial and Send
	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	// Envio
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("erro ao enviar email de solicitação de visita: %w", err) // Boa prática: enriquecer o erro
	}

	return nil
}

// createVisitSolicitationHTML gera o corpo HTML do email de solicitação de visita.
func createVisitSolicitationHTML(patientName string, visitDate string, visitValue string, address string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="pt-BR">
    <head>
        <meta charset="UTF-8">
        <title>Nova Visita Solicitada</title>
        <style>
            body {
                background-color: #f9f9f9;
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                color: #333333;
                padding: 0;
                margin: 0;
            }
            .container {
                max-width: 600px;
                margin: 40px auto;
                background-color: #ffffff;
                border-radius: 10px;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
                padding: 30px 40px;
            }
            h2 {
                color: #FFC107; /* Cor de alerta ou atenção */
                text-align: center;
            }
            p {
                line-height: 1.6;
                font-size: 15px;
            }
            .details-box {
                background-color: #FFFDE7; /* Amarelo bem suave */
                border: 1px solid #FFECB3;
                border-radius: 6px;
                padding: 15px;
                margin: 20px 0;
            }
            .detail-item {
                margin-bottom: 8px;
                font-size: 15px;
            }
            .detail-item strong {
                color: #555555;
            }
            .footer {
                margin-top: 30px;
                font-size: 12px;
                color: #999999;
                text-align: center;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h2>🔔 Nova Solicitação de Visita</h2>
            <p>Olá,</p>
            <p>O paciente <strong>%s</strong> acabou de solicitar uma visita em sua área. Por favor, verifique os detalhes abaixo para aceitar ou recusar a solicitação.</p>
            <p>Acesse o painel de visitas para visualizar mais detalhes.</p>
            
            <div class="details-box">
                <div class="detail-item"><strong>Paciente:</strong> %s</div>
                <div class="detail-item"><strong>Data/Hora Solicitada:</strong> %s</div>
                <div class="detail-item"><strong>Valor da Visita:</strong> R$%s</div>
                <div class="detail-item"><strong>Endereço:</strong> %s</div>
            </div>

            <p>Acesse a plataforma para visualizar mais informações sobre o paciente e confirmar seu interesse.</p>

            <div class="footer">
                <p>Este é um e-mail automático. Por favor, não responda.</p>
            </div>
        </div>
    </body>
    </html>
    `, patientName, patientName, visitDate, visitValue, address) // O primeiro %s é o nome no cabeçalho.
}


