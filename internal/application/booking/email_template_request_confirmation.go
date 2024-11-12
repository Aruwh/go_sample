package application

import (
	"fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
	"strconv"
)

func BuildEmailRequestConfirmationTemplate(locale common.Locale, firstName, lastName, landingpageEndpoint, apartmentName string, booking booking.Booking) email_template.IEmailTemplate {
	subjectTranslations := map[string]string{

		"deDE": "Anfragebestätigung",
		"enGB": "Request Confirmation",
		"frFR": "Confirmation de demande",
		"itIT": "Conferma della richiesta",
	}

	messageTranslations := map[string]string{
		"deDE": `
			<!DOCTYPE html>
				<html lang="de">
				<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>FeWoServ Anfragebestätigung</title>
				</head>
				<body style="font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">
				<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
					<img src="{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png" style="width: 45%; margin-left: -1.5em;" />
					<h1 style="color: #333333;">Anfragebestätigung</h1>
					<p>Hallo <span style="font-weight: bold;">{{.NAME}}</span>,</p>
					<p>gute Neuigkeiten für dich, der Status deiner Anfrage ({{.ID}}) für die Ferienwohnung</p> 
					<h4><a href="{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}" style="color: #007BFF; text-decoration: none;">{{.APARTMENT_NAME}}</a></h4> 
					<p>hat sich geändert.</p>
					<p>Wir freuen uns dir mitteilen zu dürfen, dass die Wohnung zwischen dem {{.ARRIVAL}} - {{.DEPARTURE}} für dich reserviert wurde.</p>
					<br/>
					<p>Solltest du noch Fragen zu deiner Anfrage haben, stehen wir dir gerne zur Verfügung.</p>
					<p>Beste Grüße,</p>
					<br/>
					<p>Das FeWoServ Team</p>
				</div>

				<!-- Footer -->
				<div style="text-align: center; color: #a5a5a5; font-size: small;">
					<p><h4>FeWoServ</h4></p>
					<p>FeWoServ Andreas Jakob</p>
					<p>Gribelierstrasse 12</p>
					<p>3954 Leukerbad</p>
					<p>info@fewoserv.com</p>
					<hr/>
					<div>
					<a href="{{.LANDINGPAGE_ENDPOINT}}/impressum" style="font-size: x-small; text-decoration: none; padding-right: 1em; color: black;">Impressum</a>
					<a href="{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition" style="font-size: x-small; text-decoration: none; color: black;">Buchungsinformationen</a>
					</div>
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
				<title>FeWoServ Request Confirmation</title>
			</head>
			<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
				<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
				<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\" />
				<h1 style=\"color: #333333;\">Request Confirmation</h1>
				<p>Hello <span style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
				<p>good news for you, the status of your request ({{.ID}}) for the holiday apartment</p> 
				<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4> 
				<p>has changed.</p>
				<p>We are pleased to inform you that the apartment has been reserved for you between {{.ARRIVAL}} - {{.DEPARTURE}}.</p>
				<br/>
				<p>If you have any further questions about your request, please feel free to contact us.</p>
				<p>Best regards,</p>
				<br/>
				<p>The FeWoServ Team</p>
			</div>
			
				<!-- Footer -->
				<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
				<p><h4>FeWoServ</h4></p>
				<p>FeWoServ Andreas Jakob</p>
				<p>Gribelierstrasse 12</p>
				<p>3954 Leukerbad</p>
				<p>info@fewoserv.com</p>
				<hr/>
				<div>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; color: black;\">Booking Information</a>
				</div>
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
				<title>FeWoServ Confirmation de demande</title>
			</head>
			<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
				<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
				<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\" />
				<h1 style=\"color: #333333;\">Confirmation de demande</h1>
				<p>Bonjour <span style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
				<p>bonne nouvelle pour vous, le statut de votre demande ({{.ID}}) pour l'appartement de vacances</p> 
				<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4>
				<p>a changé.</p>
				<p>Nous sommes heureux de vous informer que l'appartement vous a été réservé entre le {{.ARRIVAL}} et le {{.DEPARTURE}}.</p>
				<br/>
				<p>Si vous avez d'autres questions concernant votre demande, n'hésitez pas à nous contacter.</p>
				<p>Cordialement,</p>
				<br/>
				<p>L'équipe FeWoServ</p>
			</div>
			
				<!-- Footer -->
				<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
				<p><h4>FeWoServ</h4></p>
				<p>FeWoServ Andreas Jakob</p>
				<p>Gribelierstrasse 12</p>
				<p>3954 Leukerbad</p>
				<p>info@fewoserv.com</p>
				<hr/>
				<div>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; color: black;\">Informations de réservation</a>
				</div>
			</div>
			</body>
			</html>
		`,

		"itIT": `
		<!DOCTYPE html>
			<html lang=\"it\">
			<head>
				<meta charset=\"UTF-8\">
				<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
				<title>FeWoServ Conferma della richiesta</title>
			</head>
			<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
				<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
				<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\" />
				<h1 style=\"color: #333333;\">Conferma della richiesta</h1>
				<p>Ciao <span style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
				<p>buone notizie per te, lo stato della tua richiesta ({{.ID}}) per l'appartamento vacanze</p>
				<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4>
				<p>è cambiato.</p>
				<p>Siamo lieti di informarti che l'appartamento è stato prenotato per te tra il {{.ARRIVAL}} e il {{.DEPARTURE}}.</p>
				<br/>
				<p>Se hai ulteriori domande sulla tua richiesta, non esitare a contattarci.</p>
				<p>Cordiali saluti,</p>
				<br/>
				<p>Il team di FeWoServ</p>
			</div>
			
				<!-- Footer -->
				<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
				<p><h4>FeWoServ</h4></p>
				<p>FeWoServ Andreas Jakob</p>
				<p>Gribelierstrasse 12</p>
				<p>3954 Leukerbad</p>
				<p>info@fewoserv.com</p>
				<hr/>
				<div>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; color: black;\">Informazioni di prenotazione</a>
				</div>
			</div>
			</body>
			</html>
		`,
	}

	messageVariables := map[string]string{
		"NAME":                 lastName + " " + firstName,
		"LANDINGPAGE_ENDPOINT": landingpageEndpoint,
		"APARTMENT_ID":         booking.ApartmentID,
		"APARTMENT_NAME":       apartmentName,
		"ID":                   strconv.Itoa(booking.BookingNumber),
		"ARRIVAL":              booking.FromDate.Format("2006-01-02"),
		"DEPARTURE":            booking.FromDate.Format("2006-01-02"),
	}

	subjectVariables := map[string]string{}

	emailTemplate := email_template.New(locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables)

	return emailTemplate
}
