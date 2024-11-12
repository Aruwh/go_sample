package application

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
)

var (
	subjectVariables    = map[string]string{}
	subjectTranslations = map[string]string{
		"deDE": "Schritte zum Zurücksetzen Ihres Passworts",
		"enGB": "Steps to reset your password",
		"frFR": "Étapes pour réinitialiser votre mot de passe",
		"itIT": "Passaggi per reimpostare la tua password",
	}

	messageTranslations = map[string]string{
		"deDE": `
		<!DOCTYPE html>
			<html lang=\"de\">
			<head>
				<meta charset=\"UTF-8\">
				<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
				<title>FeWoServ Einladung</title>
			</head>
			<body style=\"font-family:'Arial',sans-serif;margin:0;padding:0;background-color:white;\">
				<div style=\"max-width:600px;margin:20px auto;background-color:#ffffff;padding:20px;border-radius:8px;box-shadow:0 0 10px rgba(0,0,0,0.1);\"><img src=\"sandbox.fewoserv.com/FeWoServ_logo.png\" style=\"width:45%;margin-left:-1.5em;\">
				<h1 style=\"color:#333333;\">Passwort vergessen ?</h1>
				<p>Hallo <span style=\"font-weight:bold;\">{{.NAME}}</span>,</p>
				<p>Es scheint, als hätten Sie Ihr Passwort vergessen oder möchten es aus Sicherheitsgründen zurücksetzen. Keine Sorge, wir sind hier, um Ihnen zu helfen!</p>
				<p>Folgen Sie einfach den unten stehenden Schritten, um Ihr Passwort zurückzusetzen:</p>
				<ol>
					<li>Klicken Sie auf den folgenden Link, um zur Passwort-Zurücksetzen-Seite zu gelangen: <a href=\"{{.FE_DESTINATION}}/pwdReset/{{.TOKEN}}/{{.LOCALE}}\" style=\"color:#007BFF;text-decoration:none;\">Passwort zurücksetzen</a></li>
					<li>Auf der Seite können Sie ein neues Passwort hinterlegen.</li>
				</ol>
				<p>Bitte beachten Sie, dass der Link nur für eine begrenzte Zeit von 10 Minuten gültig ist, um die Sicherheit Ihres Kontos zu gewährleisten.</p>
				<p>Wenn Sie diese Anfrage nicht gestellt haben oder Hilfe benötigen, zögern Sie nicht, sich mit unserem Support-Team in Verbindung zu setzen.</p>
				<br/>
				<p>Beste Grüße,<br>Das FeWoServ Team</p>
				</div>
			</body>
		</html>
		`,

		"enGB": `
		<!DOCTYPE html>
			<html lang=\"en\">
			<head>
				<meta charset=\"UTF-8\">
				<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
				<title>FeWoServ Invitation</title>
			</head>
			<body style=\"font-family:'Arial',sans-serif;margin:0;padding:0;background-color:white;\">
				<div style=\"max-width:600px;margin:20px auto;background-color:#ffffff;padding:20px;border-radius:8px;box-shadow:0 0 10px rgba(0,0,0,0.1);\"><img src=\"sandbox.fewoserv.com/FeWoServ_logo.png\" style=\"width:45%;margin-left:-1.5em;\">
				<h1 style=\"color:#333333;\">Forgot Password?</h1>
				<p>Hello <span style=\"font-weight:bold;\">{{.NAME}}</span>,</p>
				<p>It seems you have forgotten your password or wish to reset it for security reasons. Don't worry, we're here to help!</p>
				<p>Simply follow the steps below to reset your password:</p>
				<ol>
					<li>Click on the following link to go to the password reset page: <a href=\"{{.FE_DESTINATION}}/pwdReset/{{.TOKEN}}/{{.LOCALE}}\" style=\"color:#007BFF;text-decoration:none;\">Reset Password</a></li>
					<li>On the page, you can set a new password.</li>
				</ol>
				<p>Please note that the link is valid for a limited time of 10 minutes to ensure the security of your account.</p>
				<p>If you did not make this request or need assistance, feel free to contact our support team.</p>
				<br/>
				<p>Best regards,<br>The FeWoServ Team</p>
				</div>
			</body>
		</html>
		`,

		"frFR": `
		<!DOCTYPE html>
			<html lang=\"fr\">
			<head>
				<meta charset=\"UTF-8\">
				<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
				<title>Invitation à FeWoServ</title>
			</head>
			<body style=\"font-family:'Arial',sans-serif;margin:0;padding:0;background-color:white;\">
				<div style=\"max-width:600px;margin:20px auto;background-color:#ffffff;padding:20px;border-radius:8px;box-shadow:0 0 10px rgba(0,0,0,0.1);\"><img src=\"sandbox.fewoserv.com/FeWoServ_logo.png\" style=\"width:45%;margin-left:-1.5em;\">
				<h1 style=\"color:#333333;\">Mot de passe oublié ?</h1>
				<p>Bonjour <span style=\"font-weight:bold;\">{{.NAME}}</span>,</p>
				<p>Il semble que vous avez oublié votre mot de passe ou souhaitez le réinitialiser pour des raisons de sécurité. Ne vous inquiétez pas, nous sommes là pour vous aider !</p>
				<p>Suivez simplement les étapes ci-dessous pour réinitialiser votre mot de passe :</p>
				<ol>
					<li>Cliquez sur le lien suivant pour accéder à la page de réinitialisation du mot de passe : <a href=\"{{.FE_DESTINATION}}/pwdReset/{{.TOKEN}}/{{.LOCALE}}\" style=\"color:#007BFF;text-decoration:none;\">Réinitialiser le mot de passe</a></li>
					<li>Sur la page, vous pouvez définir un nouveau mot de passe.</li>
				</ol>
				<p>Veuillez noter que le lien est valide pendant une durée limitée de 10 minutes pour garantir la sécurité de votre compte.</p>
				<p>Si vous n'avez pas effectué cette demande ou si vous avez besoin d'aide, n'hésitez pas à contacter notre équipe d'assistance.</p>
				<br/>
				<p>Cordialement,<br>L'équipe FeWoServ</p>
				</div>
			</body>
		</html>
		`,

		"itIT": `
		"<!DOCTYPE html>
			<html lang=\"it\">
			<head>
				<meta charset=\"UTF-8\">
				<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
				<title>Invito a FeWoServ</title>
			</head>
			<body style=\"font-family:'Arial',sans-serif;margin:0;padding:0;background-color:white;\">
				<div style=\"max-width:600px;margin:20px auto;background-color:#ffffff;padding:20px;border-radius:8px;box-shadow:0 0 10px rgba(0,0,0,0.1);\"><img src=\"sandbox.fewoserv.com/FeWoServ_logo.png\" style=\"width:45%;margin-left:-1.5em;\">
				<h1 style=\"color:#333333;\">Hai dimenticato la password?</h1>
				<p>Ciao <span style=\"font-weight:bold;\">{{.NAME}}</span>,</p>
				<p>Sembra che tu abbia dimenticato la tua password o desideri reimpostarla per motivi di sicurezza. Non preoccuparti, siamo qui per aiutarti!</p>
				<p>Segui semplicemente i passaggi seguenti per reimpostare la tua password:</p>
				<ol>
					<li>Fai clic sul seguente link per andare alla pagina di reimpostazione della password: <a href=\"{{.FE_DESTINATION}}/pwdReset/{{.TOKEN}}/{{.LOCALE}}\" style=\"color:#007BFF;text-decoration:none;\">Reimposta la password</a></li>
					<li>Sulla pagina, puoi impostare una nuova password.</li>
				</ol>
				<p>Si prega di notare che il link è valido per un periodo limitato di 10 minuti per garantire la sicurezza del tuo account.</p>
				<p>Se non hai fatto questa richiesta o hai bisogno di assistenza, non esitare a contattare il nostro team di supporto.</p>
				<br/>
				<p>Cordiali saluti,<br>Il Team FeWoServ</p>
				</div>
			</body>
		</html>
		`,
	}
)

func BuildPasswordForgottenTemplate(locale common.Locale, firstName, lastName, feDestination, token string) email_template.IEmailTemplate {
	messageVariables := map[string]string{
		"NAME":           lastName + " " + firstName,
		"FE_DESTINATION": feDestination,
		"TOKEN":          token,
		"LOCALE":         string(locale),
	}

	emailTemplate := email_template.New(locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables)

	return emailTemplate
}
