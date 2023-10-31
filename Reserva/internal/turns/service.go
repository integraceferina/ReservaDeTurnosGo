package turns

import (
	"reserva/internal/domain"
	"reserva/internal/domain/dto"
)

type Service interface {
	ReadAll() ([]domain.Turns, error)
	Create(turn dto.TurnInsert) error
	Read(id int) (domain.Turns, error)
	Update(id int, turn dto.TurnInsert) error
	Delete(id int) error
	CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error)
	ReadbyDni(dni string) ([]dto.TurnGet, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ReadAll() ([]domain.Turns, error) {
	l, err := s.r.ReadAll()
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *service) Create(turn dto.TurnInsert) error {
	err := s.r.Create(turn)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Read(id int) (domain.Turns, error) {
	t, err := s.r.Read(id)
	if err != nil {
		return domain.Turns{}, err
	}
	return t, nil
}
func (s *service) Update(id int, turn dto.TurnInsert) error {
	err := s.r.Update(id, turn)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Delete(id int) error {
	err := s.r.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error) {
	t, err := s.r.CreateTurnByDniAndEnrollment(turn)
	if err != nil {
		return dto.TurnInsert{}, err
	}
	return t, nil
}
func (s *service) ReadbyDni(dni string) ([]dto.TurnGet, error) {
	l, err := s.r.ReadTurnbyDNI(dni)
	if err != nil {
		return []dto.TurnGet{}, err
	}
	return l, nil
}
