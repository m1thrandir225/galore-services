package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
)

//go:embed *.gohtml
var EmbeddedTemplate embed.FS

type templateData struct {
	Title   string
	Message string
}

type templateRenderer struct {
	tmpl *template.Template
}

func newEmbeddedTemplateRenderer() (*templateRenderer, error) {
	tmpl, err := template.ParseFS(EmbeddedTemplate, "*.gohtml")
	if err != nil {
		return nil, err
	}

	return &templateRenderer{
		tmpl: tmpl,
	}, nil
}

func generateTemplateWithMessage(title, message string) string {
	renderer, err := newEmbeddedTemplateRenderer()
	if err != nil {
		log.Fatal(err)
	}

	var renderedHTML bytes.Buffer
	data := templateData{
		Title:   title,
		Message: message,
	}

	err = renderer.tmpl.ExecuteTemplate(&renderedHTML, "mail-template", data)
	if err != nil {
		log.Fatal(err)
	}

	return renderedHTML.String()
}

func GeneratePasswordResetSuccessfullyMail() string {
	title := fmt.Sprintf("Password Reset Successfully")
	message := fmt.Sprintf("You have successfully reset your password.")
	return generateTemplateWithMessage(title, message)
}

func GeneratePasswordOTPMail(otp string) string {
	title := fmt.Sprintf("Password Reset Request")
	message := fmt.Sprintf("You requested a password reset. Here is your OTP Code: <strong>%s</strong>", otp)
	return generateTemplateWithMessage(title, message)
}

func GenerateWelcomeMail(email string) string {
	title := fmt.Sprintf("Welcome to Galore")
	message := fmt.Sprintf("Hi, %s. Thank you for registering to the Galore app. Hope you enjoy our services.", email)
	return generateTemplateWithMessage(title, message)
}
