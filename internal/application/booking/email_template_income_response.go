package application

import (
	"fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
	"strconv"
)

func BuildEmailIncomeResponseTemplate(locale common.Locale, firstName, lastName, landingpageEndpoint, apartmentName string, booking booking.Booking) email_template.IEmailTemplate {
	subjectTranslations := map[string]string{
		"deDE": "Anfrage der Ferienwohnung {{.APARTMENT_NAME}}",
		"enGB": "Inquiry about the holiday apartment {{.APARTMENT_NAME}}",
		"frFR": "Demande de renseignements sur l'appartement de vacances {{.APARTMENT_NAME}}",
		"itIT": "Richiesta di informazioni sull'appartamento per le vacanze {{.APARTMENT_NAME}}",
	}

	messageTranslations := map[string]string{
		"deDE": `
			<!DOCTYPE html>
			<html lang="de">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>FeWoServ Anfrage Wohnung</title>
			</head>
			<body style="font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">
				<div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
				<img src="{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png" style="width: 45%; margin-left: -1.5em;">
				<h1 style="color: #333333;">Anfrage der Ferienwohnung {{.APARTMENT_NAME}}</h1>
				<p>Hallo {{.NAME}},</p>
				<p>hiermit bestätigen wir dir den Eingang der Anfrage ({{.ID}}) für die Ferienwohnung</p> 
				<h4><a href="{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}" style="color: #007BFF; text-decoration: none;">{{.APARTMENT_NAME}}</a></h4> 
				<p>Wir kümmern uns jetzt darum, schnellstmöglich zu antworten und dir mitzuteilen, ob deine Wunschwohnung in dem Zeitraum {{.ARRIVAL}} - {{.DEPATURE}} für dich verfügbar ist.</p>
				<p>Solltest du Fragen zu der Anfrage haben, stehe wir dir gerne zur Verfügung.</p>
				<p>Beste Grüße,</p>
				<p>Das FeWoServ Team</p>
				</div>
				<div style="text-align: center; color: #a5a5a5; font-size: small;">
				<p><h4>FeWoServ</h4></p>
				<p>FeWoServ Andreas Jakob</p>
				<p>Gribelierstrasse 12</p>
				<p>3954 Leukerbad</p>
				<p>info@fewoserv.com</p>
				<hr>
				<div>
					<a href="{{.LANDINGPAGE_ENDPOINT}}/impressum" style="font-size: x-small; text-decoration: none; padding-right: 1em; color: black;">Impressum</a>
					<a href="{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition" style="font-size: x-small; text-decoration: none; padding-right: 1em; color: black;">Buchungsinformationen</a>
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
				<title>FeWoServ Apartment Inquiry</title>
			</head>
			<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
				<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
				<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\">
				<h1 style=\"color: #333333;\">Apartment Inquiry: {{.APARTMENT_NAME}}</h1>
				<p>Hello {{.NAME}},</p>
				<p>This is to confirm receipt of your inquiry ({{.ID}}) for the apartment</p> 
				<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4> 
				<p>We will now proceed to respond as soon as possible and inform you whether your desired apartment is available for you during the period {{.ARRIVAL}} - {{.DEPATURE}}.</p>
				<p>If you have any questions regarding the inquiry, we are at your disposal.</p>
				<p>Best regards,</p>
				<p>The FeWoServ Team</p>
				</div>
				<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
				<p><h4>FeWoServ</h4></p>
				<p>FeWoServ Andreas Jakob</p>
				<p>Gribelierstrasse 12</p>
				<p>3954 Leukerbad</p>
				<p>info@fewoserv.com</p>
				<hr>
				<div>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Legal Notice</a>
					<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Booking Information</a>
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
	<title>Demande d'appartement FeWoServ</title>
</head>
<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
	<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
	<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\">
	<h1 style=\"color: #333333;\">Demande d'appartement {{.APARTMENT_NAME}}</h1>
	<p>Bonjour {{.NAME}},</p>
	<p>Nous vous confirmons la réception de votre demande ({{.ID}}) pour l'appartement</p> 
	<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4> 


	<p>Nous nous engageons à répondre dans les meilleurs délais et à vous informer si l'appartement de votre choix est disponible pour vous pendant la période {{.ARRIVAL}} - {{.DEPATURE}}.</p>
	<p>Si vous avez des questions concernant la demande, nous sommes à votre disposition.</p>
	<p>Cordialement,</p>
	<p>L'équipe FeWoServ</p>
	</div>
	<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
	<p><h4>FeWoServ</h4></p>
	<p>FeWoServ Andreas Jakob</p>
	<p>Gribelierstrasse 12</p>
	<p>3954 Leukerbad</p>
	<p>info@fewoserv.com</p>
	<hr>
	<div>
		<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Mentions légales</a>
		<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Informations de réservation</a>
	</div>
	</div>
</body>
</html>
		`,
		"itIT": `<!DOCTYPE html>
		<html lang=\"it\">
		<head>
			<meta charset=\"UTF-8\">
			<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
			<title>Richiesta appartamento FeWoServ</title>
		</head>
		<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
			<div style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
			<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" style=\"width: 45%; margin-left: -1.5em;\">
			<h1 style=\"color: #333333;\">Richiesta appartamento {{.APARTMENT_NAME}}</h1>
			<p>Ciao {{.NAME}},</p>
			<p>Confermiamo la ricezione della tua richiesta ({{.ID}}) per l'appartamento</p> 
			<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a></h4> 
			<p>Ci impegneremo ora a rispondere il prima possibile e informarti se l'appartamento desiderato è disponibile per te durante il periodo {{.ARRIVAL}} - {{.DEPATURE}}.</p>
			<p>Se hai domande riguardanti la richiesta, siamo a tua disposizione.</p>
			<p>Cordiali saluti,</p>
			<p>Il Team FeWoServ</p>
			</div>
			<div style=\"text-align: center; color: #a5a5a5; font-size: small;\">
			<p><h4>FeWoServ</h4></p>
			<p>FeWoServ Andreas Jakob</p>
			<p>Gribelierstrasse 12</p>
			<p>3954 Leukerbad</p>
			<p>info@fewoserv.com</p>
			<hr>
			<div>
				<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Avviso legale</a>
				<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Informazioni sulla prenotazione</a>
			</div>
			</div>
		</body>
		</html>`,
	}

	messageVariables := map[string]string{
		"NAME":                 lastName + " " + firstName,
		"LANDINGPAGE_ENDPOINT": landingpageEndpoint,
		"APARTMENT_ID":         booking.ApartmentID,
		"APARTMENT_NAME":       apartmentName,
		"ID":                   strconv.Itoa(booking.BookingNumber),
		"ARRIVAL":              booking.FromDate.Format("2006-01-02"),
		"DEPATURE":             booking.ToDate.Format("2006-01-02"),
	}

	subjectVariables := map[string]string{
		"APARTMENT_NAME": apartmentName,
	}

	emailTemplate := email_template.New(locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables)

	return emailTemplate
}
