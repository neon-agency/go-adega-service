package email

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Client interface {
	SendDriverWelcome(toEmail, toName, tempPassword, loginURL string) error
}

type SendGridClient struct {
	apiKey    string
	fromEmail string
	fromName  string
}

func NewSendGridClient(apiKey, fromEmail, fromName string) Client {
	return &SendGridClient{apiKey: apiKey, fromEmail: fromEmail, fromName: fromName}
}

func (s *SendGridClient) SendDriverWelcome(toEmail, toName, tempPassword, loginURL string) error {
	if s.apiKey == "" || s.fromEmail == "" {
		return fmt.Errorf("sendgrid não configurado")
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(toName, toEmail)
	subject := "Seu acesso de motoboy na Adega Flow"
	plain := fmt.Sprintf("Olá %s,\n\nSeu acesso de motoboy foi criado.\n\nE-mail: %s\nSenha provisória: %s\n\nAcesse: %s\n\nNo primeiro login você deverá criar uma nova senha.", toName, toEmail, tempPassword, loginURL)
	html := fmt.Sprintf(`<!doctype html>
<html lang="pt-BR">
<body style="margin:0;background:#eef3f7;font-family:Arial,Helvetica,sans-serif;color:#13293d;">
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="padding:28px 12px;background:#eef3f7;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:640px;background:#ffffff;border:1px solid #dbe3ef;border-radius:24px;overflow:hidden;box-shadow:0 18px 50px rgba(15,23,42,.12);">
          <tr>
            <td style="background:#13293d;padding:28px 32px;color:#ffffff;">
              <div style="font-size:13px;font-weight:700;letter-spacing:1.8px;text-transform:uppercase;color:#8bd3c7;">Adega Flow</div>
              <div style="margin-top:10px;font-size:32px;line-height:38px;font-weight:800;">Seu acesso de entrega está pronto.</div>
              <div style="margin-top:10px;font-size:15px;line-height:24px;color:#d6e2ee;">Entre com a senha provisória e crie sua senha definitiva no primeiro acesso.</div>
            </td>
          </tr>
          <tr>
            <td style="padding:28px 32px;">
              <div style="font-size:18px;font-weight:700;">Olá, %s.</div>
              <div style="margin-top:16px;border:1px solid #dbe3ef;border-radius:18px;padding:18px;background:#f8fafc;">
                <div style="font-size:12px;color:#64748b;text-transform:uppercase;font-weight:700;letter-spacing:.8px;">E-mail de login</div>
                <div style="margin-top:6px;font-size:18px;font-weight:700;color:#0f172a;">%s</div>
                <div style="margin-top:18px;font-size:12px;color:#64748b;text-transform:uppercase;font-weight:700;letter-spacing:.8px;">Senha provisória</div>
                <div style="margin-top:6px;font-size:22px;font-weight:800;color:#0f766e;">%s</div>
              </div>
              <div style="margin-top:24px;">
                <a href="%s" style="display:inline-block;background:#0f766e;color:#ffffff;text-decoration:none;font-weight:800;padding:15px 22px;border-radius:14px;">Entrar na área do motoboy</a>
              </div>
              <div style="margin-top:22px;font-size:13px;line-height:21px;color:#64748b;">Se o botão não abrir, acesse manualmente: <br><a href="%s" style="color:#0f766e;">%s</a></div>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, toName, toEmail, tempPassword, loginURL, loginURL, loginURL)

	message := mail.NewSingleEmail(from, subject, to, plain, html)
	response, err := sendgrid.NewSendClient(s.apiKey).Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("sendgrid status %d: %s", response.StatusCode, response.Body)
	}
	return nil
}
