package dto

type TurnPost struct {
	DentistEnrollment string `json:"enrollment"`
	PatientDni        string `json:"dni"`
	DateUp            string `json:"dateup"`
	Hour              string `json:"hour"`
	Description       string `json:"description"`
}
