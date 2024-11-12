package healthz

type (
	// // // // // // // // // // // // // // // // // // // // // //
	// REST-API RESPONSES
	// // // // // // // // // // // // // // // // // // // // // //

	// Healtzh defines user facing defenition of all reachable domain clients
	Healtzh struct {
		Ready bool `json:"ready"`
	}

	// Ping defines user facing definition of the information that the service is reactive
	Ping struct {
		Alive bool `json:"alive"`
	}
)

// ConvertFromValue implements a dto like function which will convert a bool (received from mono. api side) into a user facing response
func (Healtzh) ConvertFromValue(value bool) *Healtzh {
	return &Healtzh{Ready: value}
}

// ConvertFromValue implements a dto like function which will convert a bool into a user facing response
func (Ping) ConvertFromValue(value bool) *Ping {
	return &Ping{Alive: value}
}
