package application

import (
	"fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/email_template"
	"strconv"
)

func BuildEmailIncomeCancelationTemplate(locale common.Locale, firstName, lastName, landingpageEndpoint, apartmentName string, booking booking.Booking) email_template.IEmailTemplate {
	subjectTranslations := map[string]string{
		"deDE": "Stornierung deiner Buchung",
		"enGB": "Cancellation of your booking",
		"frFR": "Annulation de votre réservation",
		"itIT": "Annullamento della tua prenotazione",
	}

	messageTranslations := map[string]string{
		"deDE": `
				<!DOCTYPE html>
				<html lang="de">
				<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>FeWoServ Stornierung</title>
				</head>
				<body style="font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;">
				<div class="email-container" style="max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">
					<img src="{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png" class="logo" style="width: 45%; margin-left: -1.5em;" />
					<h1 style="color: #333333;">Stornierung deiner Buchung</h1>
					<p>Hallo <span class="name-placeholder" style="font-weight: bold;">{{.NAME}}</span>,</p>
					<p>mit bedauern müssen wir dir Mitteilen, dass wir deine Buchung ({{.ID}}) für die Ferienwohnung</p>
					<h4><a href="{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}" style="color: #007BFF; text-decoration: none;">{{.APARTMENT_NAME}}</a></h4> 
					<p>stornieren mussten.</p>
					<br/>
					<p>Solltest du Fragen zu der Stornierung haben, stehe wir dir gerne zur Verfügung.</p>
					<p>Beste Grüße,</p>
					<br/>
					<p>Das FeWoServ Team</p>
				</div>
				
				<!-- Footer -->
				<div class="footer" style="text-align: center; color: #a5a5a5; font-size: small;">
					<h4>FeWoServ</h4>
					<p>FeWoServ Andreas Jakob</p>
					<p>Gribelierstrasse 12</p>
					<p>3954 Leukerbad</p>
					<p>info@fewoserv.com</p>
					<hr/>
					<div class="footerLinks">
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
					<title>FeWoServ Cancellation</title>
				</head>
				<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
					<div class=\"email-container\" style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
					<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" class=\"logo\" style=\"width: 45%; margin-left: -1.5em;\" />
					<h1 style=\"color: #333333;\">Cancellation of your booking</h1>
					<p>Hello <span class=\"name-placeholder\" style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
					<p>We regret to inform you that we had to cancel your booking ({{.ID}}) for the holiday apartment</p>
					<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a>.</h4>
					<br/>
					<p>If you have any questions regarding the cancellation, please feel free to contact us.</p>
					<p>Best regards,</p>
					<br/>
					<p>The FeWoServ Team</p>
				</div>
				
					<!-- Footer -->
					<div class=\"footer\" style=\"text-align: center; color: #a5a5a5; font-size: small;\">
					<h4>FeWoServ</h4>
					<p>FeWoServ Andreas Jakob</p>
					<p>Gribelierstrasse 12</p>
					<p>3954 Leukerbad</p>
					<p>info@fewoserv.com</p>
					<hr/>
					<div class=\"footerLinks\">
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Buchungsinformationen</a>
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
					<title>Annulation FeWoServ</title>
				</head>
				<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
					<div class=\"email-container\" style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
					<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" class=\"logo\" style=\"width: 45%; margin-left: -1.5em;\" />
					<h1 style=\"color: #333333;\">Annulation de votre réservation</h1>
					<p>Bonjour <span class=\"name-placeholder\" style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
					<p>Nous sommes au regret de vous informer que nous avons dû annuler votre réservation ({{.ID}}) pour l'appartement de vacances</p>
					<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a>.</h4>
					<br/>
					<p>Si vous avez des questions concernant l'annulation, n'hésitez pas à nous contacter.</p>
					<p>Cordialement,</p>
					<br/>
					<p>L'équipe FeWoServ</p>
				</div>
				
					<!-- Footer -->
					<div class=\"footer\" style=\"text-align: center; color: #a5a5a5; font-size: small;\">
					<h4>FeWoServ</h4>
					<p>FeWoServ Andreas Jakob</p>
					<p>Gribelierstrasse 12</p>
					<p>3954 Leukerbad</p>
					<p>info@fewoserv.com</p>
					<hr/>
					<div class=\"footerLinks\">
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Buchungsinformationen</a>
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
					<title>Annulamento FeWoServ</title>
				</head>
				<body style=\"font-family: 'Arial', sans-serif; margin: 0; padding: 0; background-color: #f7f7f7;\">
					<div class=\"email-container\" style=\"max-width: 600px; margin: 20px auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">
					<img src=\"{{.LANDINGPAGE_ENDPOINT}}/FeWoServ_logo.png\" class=\"logo\" style=\"width: 45%; margin-left: -1.5em;\" />
					<h1 style=\"color: #333333;\">Cancellazione della tua prenotazione</h1>
					<p>Ciao <span class=\"name-placeholder\" style=\"font-weight: bold;\">{{.NAME}}</span>,</p>
					<p>Ci dispiace informarti che abbiamo dovuto cancellare la tua prenotazione ({{.ID}}) per l'appartamento vacanze</p>
					<h4><a href=\"{{.LANDINGPAGE_ENDPOINT}}/apartment/{{.APARTMENT_ID}}\" style=\"color: #007BFF; text-decoration: none;\">{{.APARTMENT_NAME}}</a>.</h4>
					<br/>
					<p>Se hai domande riguardanti la cancellazione, non esitare a contattarci.</p>
					<p>Cordiali saluti,</p>
					<br/>
					<p>Il Team FeWoServ</p>
				</div>
				
					<!-- Footer -->
					<div class=\"footer\" style=\"text-align: center; color: #a5a5a5; font-size: small;\">
					<h4>FeWoServ</h4>
					<p>FeWoServ Andreas Jakob</p>
					<p>Gribelierstrasse 12</p>
					<p>3954 Leukerbad</p>
					<p>info@fewoserv.com</p>
					<hr/>
					<div class=\"footerLinks\">
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/impressum\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Impressum</a>
						<a href=\"{{.LANDINGPAGE_ENDPOINT}}/termsOfCondition\" style=\"font-size: x-small; text-decoration: none; padding-right: 1em; color: black;\">Buchungsinformationen</a>
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
		"ID":                   strconv.Itoa(booking.BookingNumber),
		"APARTMENT_NAME":       apartmentName,
	}

	subjectVariables := map[string]string{
		"APARTMENT_NAME": apartmentName,
	}

	emailTemplate := email_template.New(locale, subjectTranslations, subjectVariables, messageTranslations, messageVariables)

	return emailTemplate
}
