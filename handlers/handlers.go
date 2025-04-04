package handlers

import (
	"api-enviar-email-moleniuk/database"
	"api-enviar-email-moleniuk/models"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

type EmailLog struct {
	Sender         string
	Title          string
	ContentBody    string
	AttachmentName string
	Attachment     string
	SentAt         time.Time
}

func sendEmailAsync(req models.EmailRequest, smtpUser, smtpPass, smtpHost, smtpPort string, errChan chan<- error, dbChan chan<- EmailLog) {
	e := email.NewEmail()
	e.From = smtpUser
	e.ReplyTo = []string{"no-reply"}
	e.To = []string{req.Sender}
	e.Subject = req.Title
	e.HTML = []byte(req.ContentBody)

	if req.Attachment != "" {
		var err error
		var filePath string

		if _, err := os.Stat(req.Attachment); err == nil {
			filePath = req.Attachment
		} else {
			if req.NameAttachment == "" {
				errChan <- fmt.Errorf("nome do anexo é obrigatório quando enviado em Base64")
				return
			}
			decodedFile, err := base64.StdEncoding.DecodeString(req.Attachment)
			if err != nil {
				errChan <- fmt.Errorf("erro ao decodificar o anexo Base64: %v", err)
				return
			}
			filePath = filepath.Join(os.TempDir(), req.NameAttachment)
			if err = os.WriteFile(filePath, decodedFile, 0644); err != nil {
				errChan <- fmt.Errorf("erro ao salvar o anexo temporário: %v", err)
				return
			}
		}

		if _, err = e.AttachFile(filePath); err != nil {
			errChan <- fmt.Errorf("falha ao anexar arquivo: %v", err)
			return
		}
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	tlsConfig := &tls.Config{ServerName: smtpHost}
	endereco := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	var err error
	if smtpPort == "587" {
		err = e.SendWithStartTLS(endereco, auth, tlsConfig)
	} else {
		err = e.SendWithTLS(endereco, auth, tlsConfig)
	}

	if err != nil {
		errChan <- fmt.Errorf("erro ao enviar e-mail: %v", err)
		return
	}

	emailLog := EmailLog{
		Sender:         req.Sender,
		Title:          req.Title,
		ContentBody:    req.ContentBody,
		AttachmentName: req.NameAttachment,
		Attachment:     req.Attachment,
		SentAt:         time.Now(),
	}

	dbChan <- emailLog
	errChan <- nil
}

func SendEmail(c *gin.Context) {
	var req models.EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.EmailResponse{
			Status:  false,
			Message: "JSON inválido: " + err.Error(),
		})
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		c.JSON(500, models.EmailResponse{
			Status:  false,
			Message: "Configuração SMTP incompleta",
		})
		return
	}

	errChan := make(chan error, 1)
	dbChan := make(chan EmailLog, 1)

	go sendEmailAsync(req, smtpUser, smtpPass, smtpHost, smtpPort, errChan, dbChan)

	err := <-errChan
	if err != nil {
		c.JSON(500, models.EmailResponse{
			Status:  false,
			Message: "E-mail não enviado: " + err.Error(),
		})
		return
	}

	emailLog := <-dbChan
	_, err = database.DB.Exec(`
        INSERT INTO email_logs (sender, title, content_body, name_attachment, attachment, sent_at) 
        VALUES ($1, $2, $3, $4, $5, $6)`,
		emailLog.Sender,
		emailLog.Title,
		emailLog.ContentBody,
		emailLog.AttachmentName,
		emailLog.Attachment,
		emailLog.SentAt,
	)
	if err != nil {
		log.Printf("Erro ao salvar log do email no banco: %v", err)
	}

	c.JSON(200, models.EmailResponse{
		Status:  true,
		Message: "E-mail enviado com sucesso",
	})
}
