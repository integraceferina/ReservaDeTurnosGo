package dentist

import (
	"errors"

	"reserva/internal/domain"
	"reserva/pkg/store"
)

type Repository interface {
	ReadAll() ([]domain.Dentist, error)
	Read(id int) (domain.Dentist, error)
	Create(dentist domain.Dentist) error
	Update(id int, dentist domain.Dentist) error
	Delete(id int) error
}

type repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return &repository{store: store}
}

func (r *repository) ReadAll() ([]domain.Dentist, error) {
	list, err := r.store.ReadAllDentists()
	if err != nil {
		return []domain.Dentist{}, err
	}
	return list, nil

}
func (r *repository) Read(id int) (domain.Dentist, error) {
	dentist, err := r.store.ReadDentist(id)
	if err != nil {
		return domain.Dentist{}, errors.New("Dentist not found")
	}
	return dentist, nil
}

func (r *repository) Create(dentist domain.Dentist) error {
	if !r.ValidateEnrollment(dentist.Id, dentist.Enrollment) {
		return errors.New("Enrollment already exists")
	}
	err := r.store.CreateDentist(dentist)
	if err != nil {
		return errors.New("Error creating a new Dentist")
	}
	return nil
}
func (r *repository) Update(id int, dentist domain.Dentist) error {
	if !r.ValidateEnrollment(id, dentist.Enrollment) {
		return errors.New("Enrollment already exists")
	}
	original, err := r.store.ReadDentist(id)
	if err != nil {
		return errors.New("The Dentist does not exists")
	}
	complete := unchangeEmptys(dentist, original)
	err = r.store.UpdateDentist(id, complete)
	if err != nil {
		return errors.New("Error updating a new Dentist")
	}
	return nil
}
func (r *repository) Delete(id int) error {
	err := r.store.DeleteDentist(id)
	if err != nil {
		return errors.New("Error deleting a Dentist, Cause 1:he have still turns. Cause 2: He doest exist.")
	}
	return nil
}

// Validation Functions
func (r *repository) ValidateEnrollment(id int, enrollment string) bool {
	list, err := r.store.ReadAllDentists()
	if err != nil {
		return false
	}
	for _, dentist := range list {
		if dentist.Enrollment == enrollment && dentist.Id != id {
			return false
		}
	}
	return true
}
func unchangeEmptys(dentist domain.Dentist, original domain.Dentist) domain.Dentist {

	if dentist.Name == "" {
		dentist.Name = original.Name
	}
	if dentist.LastName == "" {
		dentist.LastName = original.LastName
	}
	if dentist.Enrollment == "" {
		dentist.Enrollment = original.Enrollment
	}
	return dentist
}
