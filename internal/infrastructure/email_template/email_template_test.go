package email_template

import (
	"encoding/json"
	"fewoserv/internal/infrastructure/common"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	name          = "tim ley"
	feDestination = "http://bla"
	apartmentID   = "alsjdlaksdjlksad"
	apartmentName = "\"Penthaus Constantin\""
	locale        = common.EnGB
)

func TestBuildVariablesFromJson(t *testing.T) {

	msgVars := fmt.Sprintf(`
	{
		"NAME": "%s",
		"FE_DESTINATION": "%s",
		"APARTMENT_ID": "%s",
		"APARTMENT_NAME": "%s",
		"LOCALE": "%s"
	}
	`, name, feDestination, apartmentID, strconv.Quote(apartmentName), locale)

	// msgVars := "{\n\t\t\"APARTMENT_NAME\": \"Wohnung 324 \\\"Penthaus Constantin\\\"\"\n\t}"
	mappedVars := buildVariablesFromJson(msgVars)

	fmt.Sprintln(mappedVars["APARTMENT_NAME"])

	x := mappedVars["APARTMENT_NAME"]
	fmt.Sprintln(x)

	assert.NotNil(t, mappedVars)
}

func TestBla(t *testing.T) {
	data := map[string]string{
        "multiline": `Dies ist ein
                      mehrzeiliger
                      Textinhalt.`,
    }

    // Konvertierung der Map in JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Fehler beim Marshalling:", err)
        return
    }

    // Ausgabe des JSON-Strings
    fmt.Println(string(jsonData))
}
