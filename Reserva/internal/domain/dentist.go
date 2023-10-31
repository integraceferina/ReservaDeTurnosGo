package domain

type Dentist struct {
	Id         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	LastName   string `json:"lastname,omitempty"`
	Enrollment string `json:"enrollment,omitempty"`
}
