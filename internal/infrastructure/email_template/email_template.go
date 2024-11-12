package email_template

import (
	"bytes"
	"encoding/json"
	"fewoserv/internal/infrastructure/common"
	"fmt"
	"html/template"
	"strings"
)

// sample json subject
// {
//     "subject": {
//         "variables": {"PLACEHOLDER": "sample value" },
//         "translations": {"deDE": "","enGB": ""}
//     },

// sample json message
//     "message": {
//         "variables": {"PLACEHOLDER": "sample value" },
//         "translations": {"deDE": "","enGB": ""}
//     }
// }

type (
	IEmailTemplate interface {
		Build(fromName, fromEmail, toEmail string) ([]byte, error)
	}

	ContentData struct {
		Variables    map[string]string `json:"variables"`
		Translations map[string]string `json:"translations"`
	}

	Content struct {
		Subject ContentData `json:"subject"`
		Message ContentData `json:"message"`
	}

	EmailTemplate struct {
		Locale  common.Locale
		Content Content
	}
)

func New(locale common.Locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables map[string]string) IEmailTemplate {
	subject := ContentData{
		Variables:    subjectVariables,
		Translations: subjectTranslations,
	}

	message := ContentData{
		Variables:    messageVariables,
		Translations: messageTranslations,
	}

	return &EmailTemplate{
		Locale:  locale,
		Content: Content{Subject: subject, Message: message},
	}
}

func buildVariablesFromJson(jsonString string) map[string]string {
	variables := make(map[string]string)
	err := json.Unmarshal([]byte(jsonString), &variables)
	if err != nil {
		fmt.Print(err)
	}

	return variables
}

func buildTranslationsFromJson(jsonString string) map[string]string {
	template := make(map[string]string)
	err := json.Unmarshal([]byte(jsonString), &template)
	if err != nil {
		fmt.Print(err)
	}

	return template
}

func transform(usedTemplate string, variables map[string]string) (*string, error) {
	tmpl, err := template.New("tmp").Parse(usedTemplate)
	if err != nil {
		return nil, err
	}
	tmpl.Option()

	var transformedTemplate bytes.Buffer
	err = tmpl.Execute(&transformedTemplate, variables)
	if err != nil {
		return nil, err
	}

	stringifiedTemplate := transformedTemplate.String()
	return &stringifiedTemplate, nil
}

func (eth *EmailTemplate) getSubjectTemplate() string {
	value, doesValueExists := eth.Content.Subject.Translations[string(eth.Locale)]

	if doesValueExists {
		return value
	} else {
		return ""
	}
}

func (eth *EmailTemplate) getMessageTemplate() string {
	value, doesValueExists := eth.Content.Message.Translations[string(eth.Locale)]

	if doesValueExists {
		return value
	} else {
		return ""
	}
}

func (eth *EmailTemplate) createMimeMessage(subject, body, fromName, fromEmail, toEmail string) []byte {
	var message strings.Builder

	// Mime-Version-Header
	message.WriteString("MIME-version: 1.0;\r\n")
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// Content-Type-Header
	contentType := "text/html; charset=UTF-8"
	message.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))

	// From-Header
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, fromEmail))

	// To-Header
	message.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))

	message.WriteString("\r\n")

	// HTML
	message.WriteString(body)

	return []byte(message.String())
}

func (eth *EmailTemplate) Build(fromName, fromEmail, toEmail string) ([]byte, error) {
	messageTemplate := eth.getMessageTemplate()
	transformedHtml, err := transform(messageTemplate, eth.Content.Message.Variables)
	if err != nil {
		return nil, err
	}

	subject := eth.getSubjectTemplate()
	transformedSubject, err := transform(subject, eth.Content.Subject.Variables)
	if err != nil {
		return nil, err
	}

	compiledTemplate := eth.createMimeMessage(*transformedSubject, *transformedHtml, fromName, fromEmail, toEmail)
	return compiledTemplate, nil
}
