package turns

import (
	"errors"
	"strings"
	"time"

	"reserva/internal/domain"
	"reserva/internal/domain/dto"
	"reserva/pkg/store"
)

type Repository interface {
	ReadAll() ([]domain.Turns, error)
	Create(turn dto.TurnInsert) error
	Read(id int) (domain.Turns, error)
	Update(id int, turns dto.TurnInsert) error
	Delete(id int) error
	CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error)
	ReadTurnbyDNI(dni string) ([]dto.TurnGet, error)
}

type repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return &repository{store: store}
}

func (r *repository) ReadAll() ([]domain.Turns, error) {
	list, err := r.store.ReadAllTurns()
	if err != nil {
		return []domain.Turns{}, err
	}
	return list, nil
}

func (r *repository) Create(turn dto.TurnInsert) error {
	if !r.ValidatePatient(turn.PatientId) {
		return errors.New("The patient does not exist")
	}
	if !r.ValidateDentist(turn.DentistId) {
		return errors.New("The dentist does not exist")
	}
	err := r.store.CreateTurn(turn)
	if err != nil {
		return errors.New("Error creating a new turn: " + err.Error())
	}
	return nil
}

func (r *repository) Read(id int) (domain.Turns, error) {
	turno, err := r.store.ReadTurn(id)
	if err != nil {
		return domain.Turns{}, err
	}
	return turno, nil
}

func (r *repository) Update(id int, turn dto.TurnInsert) error {
	original, err := r.store.ReadTurn(id)
	if err != nil {
		return errors.New("The Turn does not exists")
	}
	if !r.ValidateDentist(original.Dentist.Id) {
		return errors.New("The Dentist does not exists")
	}
	if !r.ValidatePatient(original.Patient.Id) {
		return errors.New("The Patient does not exists")
	}
	complete := unchangeEmptysTurn(turn, original)
	err = r.store.UpdateTurn(id, complete)
	if err != nil {
		return errors.New("Error updating Turn")
	}
	return nil
}
func (r *repository) Delete(id int) error {
	err := r.store.DeleteTurn(id)
	if err != nil {
		return errors.New("Error deleting a Dentist")
	}
	return nil
}

func (r *repository) CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error) {
	updateturn := ChargeDataTurn(turn)
	turnok, err := r.store.CreateTurnByDniAndEnrollment(updateturn)
	if err != nil {
		return dto.TurnInsert{}, err
	}
	return turnok, nil
}

func (r *repository) ReadTurnbyDNI(dni string) ([]dto.TurnGet, error) {
	if !r.ValidateDNI(dni) {
		return []dto.TurnGet{}, errors.New("There is no patient with that DNI")
	}
	list, err := r.store.ReadTurnbyDNI(dni)
	if err != nil {
		return []dto.TurnGet{}, errors.New("Error getting a turn by DNI")
	}
	return list, nil
}

// Validate Functions

func (r *repository) ValidatePatient(id int) bool {
	_, err := r.store.ReadPatient(id)
	if err != nil {
		return false
	}
	return true
}

func (r *repository) ValidateDentist(id int) bool {
	_, err := r.store.ReadDentist(id)
	if err != nil {
		return false
	}
	return true
}

func (r *repository) ValidateDNI(dni string) bool {
	list, err := r.store.ReadAllPatient()
	if err != nil {
		return false
	}
	for _, patient := range list {
		if patient.DNI == dni {
			return true
		}
	}
	return false
}

func unchangeEmptysTurn(turnDTO dto.TurnInsert, original domain.Turns) dto.TurnInsert {

	if turnDTO.DateUp == "" {
		turnDTO.DateUp = original.DateUp
	}
	if turnDTO.Hour == "" {
		turnDTO.Hour = original.Hour
	}
	if turnDTO.Description == "" {
		turnDTO.Description = original.Description
	}
	if turnDTO.DentistId == 0 {
		turnDTO.DentistId = original.Dentist.Id
	}
	if turnDTO.PatientId == 0 {
		turnDTO.PatientId = original.Patient.Id
	}

	return turnDTO
}

func ChargeDataTurn(turn dto.TurnPost) dto.TurnPost {
	t := time.Now()
	s := t.Format("2006-01-02 15:04:05")
	fechayhora := strings.Split(s, " ")
	if turn.DateUp == "" {
		turn.DateUp = fechayhora[0]
	}
	if turn.Hour == "" {
		turn.Hour = fechayhora[1]
	}
	if turn.Description == "" {
		turn.Description = "Turno creado la fecha: " + fechayhora[0] + " " + fechayhora[1]
	}
	return turn
}
