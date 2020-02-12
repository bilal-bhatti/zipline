package models

type (
	Name struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode string `json:"zipCode"`
	}

	ContactRequest struct {
		Input string `json:"input"`
		Name
		Address *Address `json:"address"`
	}

	ContactResponse struct {
		Output string `json:"output"`
	}
)
