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

	// Contact request model
	ContactRequest struct {
		ID string `json:"id"`
		Name
		Address *Address `json:"address"`
	}

	// Contact response model
	ContactResponse struct {
		// id is the unique id of contact
		ID string `json:"id"`
	}
)
