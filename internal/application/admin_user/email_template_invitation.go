package application

import (
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
)

var (
	subjectVariables    = map[string]string{}
	subjectTranslations = map[string]string{
		"deDE": "Einladung zur FeWoServ",
		"enGB": "Invitation to FeWoServ",
		"frFR": "Invitation à FeWoServ",
		"itIT": "Invito a FeWoServ",
	}

	messageTranslations = map[string]string{
		"deDE": `
		<!DOCTYPE html>
			<html lang="de">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>FeWoServ Einladung</title>
			</head>
			<body style="font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">
			<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
				<img src="{{.FE_DESTINATION}}/FeWoServ_logo.png" style="width: 45%; margin-left: -1.5em;">
				<h1 style="color: #333333;">Einladung zur FeWoServ WebApp</h1>
				<p>Sehr geehrte/r <span style="font-weight: bold;">{{.NAME}}</span>,</p>
				<p>Wir freuen uns, Sie zur Nutzung der FeWoServ WebApp einzuladen!</p>
				<p>Die Plattform erleichtert die Verwaltung Ihrer Ferienwohnungen, bietet einen übersichtlichen Kalender, einfache Gästekommunikation und zuverlässiges Buchungsmanagement.</p>
				<p>Um Ihre Registrierung abschließen zu können, bitten wir Sie den folgenden Link zu verwenden und die auf der Seite befindlichen Anweisungen zu befolgen.</p>
				<p>Link: <a href="{{.FE_DESTINATION}}/invite/{{.TOKEN}}/{{.LOCALE}}" style="color: #007BFF; text-decoration: none;">FeWoServ Registrierung</a></p>
				<br/>
				<p>Bei Fragen stehen wir gerne zur Verfügung.</p>
				<p>Beste Grüße,<br>Das FeWoServ Team</p>
			</div>
			</body>
		</html>
		`,

		"enGB": `
		<!DOCTYPE html>
			<html lang='en'>
			<head>
				<meta charset='UTF-8'>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'>
				<title>FeWoServ Invitation</title>
			</head>
			<body style='font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;'>
			<div style='max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);'>
				<img src='{{.FE_DESTINATION}}/FeWoServ_logo.png' style='width: 45%; margin-left: -1.5em;'>
				<h1 style='color: #333333;'>Invitation to the FeWoServ WebApp</h1>
				<p>Dear <span style='font-weight: bold;'>{{.NAME}}</span>,</p>
				<p>We are pleased to invite you to use the FeWoServ WebApp!</p>
				<p>The platform facilitates the management of your holiday apartments, provides a clear calendar, simple guest communication, and reliable booking management.</p>
				<p>To complete your registration, please use the following link and follow the instructions on the page.</p>
				<p>Link: <a href='{{.FE_DESTINATION}}/invite/{{.TOKEN}}/{{.LOCALE}}' style='color: #007BFF; text-decoration: none;'>FeWoServ Registration</a></p>
				<br/>
				<p>If you have any questions, please feel free to contact us.</p>
				<p>Best regards,<br>The FeWoServ Team</p>
			</div>
			</body>
		</html>
		`,

		"frFR": `
		<!DOCTYPE html>
			<html lang='fr'>
			<head>
				<meta charset='UTF-8'>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'>
				<title>Invitation à FeWoServ</title>
			</head>
			<body style='font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;'>
			<div style='max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);'>
				<img src='{{.FE_DESTINATION}}/FeWoServ_logo.png' style='width: 45%; margin-left: -1.5em;'>
				<h1 style='color: #333333;'>Invitation à l'application web FeWoServ</h1>
				<p>Cher/Chère <span style='font-weight: bold;'>{{.NAME}}</span>,</p>
				<p>Nous avons le plaisir de vous inviter à utiliser l'application web FeWoServ !</p>
				<p>La plateforme facilite la gestion de vos appartements de vacances, offre un calendrier clair, une communication simple avec les invités et une gestion fiable des réservations.</p>
				<p>Pour finaliser votre inscription, veuillez utiliser le lien suivant et suivre les instructions sur la page.</p>
				<p>Lien : <a href='{{.FE_DESTINATION}}/invite/{{.TOKEN}}/{{.LOCALE}}' style='color: #007BFF; text-decoration: none;'>Inscription à FeWoServ</a></p>
				<br/>
				<p>Si vous avez des questions, n'hésitez pas à nous contacter.</p>
				<p>Cordialement,<br>L'équipe FeWoServ</p>
			</div>
			</body>
		</html>
		`,

		"itIT": `
		<!DOCTYPE html>
			<html lang='it'>
			<head>
				<meta charset='UTF-8'>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'>
				<title>Invito a FeWoServ</title>
			</head>
			<body style='font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;'>
			<div style='max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);'>
				<img src='{{.FE_DESTINATION}}/FeWoServ_logo.png' style='width: 45%; margin-left: -1.5em;'>
				<h1 style='color: #333333;'>Invito alla WebApp FeWoServ</h1>
				<p>Gentile <span style='font-weight: bold;'>{{.NAME}}</span>,</p>
				<p>Siamo lieti di invitarti a utilizzare la WebApp FeWoServ!</p>
				<p>La piattaforma facilita la gestione dei tuoi appartamenti vacanza, fornisce un calendario chiaro, una semplice comunicazione con gli ospiti e una gestione delle prenotazioni affidabile.</p>
				<p>Per completare la tua registrazione, utilizza il seguente link e segui le istruzioni sulla pagina.</p>
				<p>Link: <a href='{{.FE_DESTINATION}}/invite/{{.TOKEN}}/{{.LOCALE}}' style='color: #007BFF; text-decoration: none;'>Registrazione a FeWoServ</a></p>
				<br/>
				<p>Se hai domande, non esitare a contattarci.</p>
				<p>Cordiali saluti,<br>Il Team FeWoServ</p>
			</div>
			</body>
		</html>
		`,
	}
)

func BuildEmailInvitationTemplate(locale common.Locale, firstName, lastName, feDestination, token string) email_template.IEmailTemplate {
	messageVariables := map[string]string{
		"NAME":           lastName + " " + firstName,
		"FE_DESTINATION": feDestination,
		"TOKEN":          token,
		"LOCALE":         string(locale),
	}

	emailTemplate := email_template.New(locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables)

	return emailTemplate
}
