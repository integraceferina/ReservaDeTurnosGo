package store

import (
	"database/sql"
	"errors"
	"log"

	"reserva/internal/domain"
	"reserva/internal/domain/dto"
)

type store struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &store{db: db}
}

type Store interface {
	// ---------------- DENTIST
	// ReadAllDentist Trae todos los dentistas
	ReadAllDentists() ([]domain.Dentist, error)
	// Read devuelve un dentista por su id
	ReadDentist(id int) (domain.Dentist, error)
	// CreateDentist agrega un nuevo dentista
	CreateDentist(dentist domain.Dentist) error
	// UpdateDentist actualiza un dentista
	UpdateDentist(id int, dentist domain.Dentist) error
	// DeleteDentist elimina un dentista
	DeleteDentist(id int) error

	// ---------------- PATIENT
	// ReadAll Trae todos los pacientes
	ReadAllPatient() ([]domain.Patient, error)
	// Read devuelve un paciente por su id
	ReadPatient(id int) (domain.Patient, error)
	// Create agrega un nuevo paciente
	CreatePatient(dentist domain.Patient) error
	// Update actualiza un paciente
	UpdatePatient(id int, patient domain.Patient) error
	// Delete elimina un paciente
	DeletePatient(id int) error

	// ---------------- TURNS
	// Devuelve todos los turnos
	ReadAllTurns() ([]domain.Turns, error)
	// Create crea un turno
	CreateTurn(turn dto.TurnInsert) error
	// Leer un turno x ID
	ReadTurn(id int) (domain.Turns, error)
	// Update actualiza un turno
	UpdateTurn(id int, turnDTO dto.TurnInsert) error
	// Delete elimina un turno
	DeleteTurn(id int) error
	//Crea un turno x DNI del Paciente y la Matricula del Dentista
	CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error)
	//Devuelve un turno x DNI de paciente
	ReadTurnbyDNI(dni string) ([]dto.TurnGet, error)
}

// 		----------DENTIST--------------

func (s *store) ReadAllDentists() ([]domain.Dentist, error) {
	var list []domain.Dentist
	var dentist domain.Dentist
	rows, err := s.db.Query("SELECT * FROM dentist")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&dentist.Id, &dentist.Name, &dentist.LastName, &dentist.Enrollment); err != nil {
			return nil, err
		}
		list = append(list, dentist)
	}
	rows.Close()
	return list, nil
}

func (s *store) ReadDentist(id int) (domain.Dentist, error) {
	//defer func () {s.db.Close()}()
	var dentist domain.Dentist
	row := s.db.QueryRow("SELECT * FROM dentist WHERE id=?", id)

	if err := row.Scan(&dentist.Id, &dentist.Name, &dentist.LastName, &dentist.Enrollment); err != nil {
		return domain.Dentist{}, errors.New("The dentist with this id not exist")
		//panic(err)
	}

	return dentist, nil
}

func (s *store) CreateDentist(dentist domain.Dentist) error {

	st, err := s.db.Prepare("INSERT INTO dentist (name, lastname, enrollment) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(dentist.Name, dentist.LastName, dentist.Enrollment)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateDentist(id int, dentist domain.Dentist) error {

	st, err := s.db.Prepare("UPDATE dentist SET name = ?, lastName = ?, enrollment = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	_, err = st.Exec(dentist.Name, dentist.LastName, dentist.Enrollment, id)
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func (s *store) DeleteDentist(id int) error {
	//Preguntar si esta bien usar un metodo
	var idselect int
	row := s.db.QueryRow("SELECT id FROM dentist WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The dentist doest exists.")
	}
	query := "DELETE FROM dentist WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil

}

// 		----------PATIENT--------------

func (s *store) ReadAllPatient() ([]domain.Patient, error) {
	var list []domain.Patient
	var patient domain.Patient
	rows, err := s.db.Query("SELECT * FROM patient")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&patient.Id, &patient.Name, &patient.LastName, &patient.Domicile, &patient.DNI, &patient.DateUp); err != nil {
			return nil, err
		}
		list = append(list, patient)
	}
	rows.Close()
	return list, nil

}

func (s *store) ReadPatient(id int) (domain.Patient, error) {
	//defer func () {s.db.Close()}()
	var patient domain.Patient
	row := s.db.QueryRow("SELECT * FROM patient WHERE id=?", id)

	if err := row.Scan(&patient.Id, &patient.Name, &patient.LastName, &patient.Domicile, &patient.DNI, &patient.DateUp); err != nil {
		return domain.Patient{}, err
		//panic(patient
	}
	return patient, nil
}

func (s *store) CreatePatient(patient domain.Patient) error {
	st, err := s.db.Prepare("INSERT INTO patient (name, lastname, domicile, dni, dateup) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(patient.Name, patient.LastName, patient.Domicile, patient.DNI, patient.DateUp)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UpdatePatient(id int, patient domain.Patient) error {
	st, err := s.db.Prepare("UPDATE patient SET name = ?, lastName = ?, domicile = ?, dni = ?, dateup = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(patient.Name, patient.LastName, patient.Domicile, patient.DNI, patient.DateUp, id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) DeletePatient(id int) error {
	//Preguntar si esta bien usar un metodo
	var idselect int
	row := s.db.QueryRow("SELECT id FROM patient WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The patient doest exists.")
	}
	query := "DELETE FROM patient WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil

}

// 		----------TURNS--------------

func (s *store) ReadAllTurns() ([]domain.Turns, error) {
	var list []domain.Turns
	var turn domain.Turns
	rows, err := s.db.Query("SELECT t.id, t.date, t.hour, t.description, p.id, d.id FROM turns AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id")
	if err != nil {
		return nil, err
	}
	for rows.Next() {

		if err := rows.Scan(&turn.Id, &turn.DateUp, &turn.Hour, &turn.Description, &turn.Patient.Id, &turn.Dentist.Id); err != nil {
			return nil, err
		}
		list = append(list, turn)
	}
	rows.Close()
	return list, nil
}

func (s *store) ReadTurn(id int) (domain.Turns, error) {
	var turn domain.Turns
	row := s.db.QueryRow("SELECT t.id, t.date, t.hour, t.description, p.id, d.id FROM turns AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id WHERE t.id = ?", id)

	if err := row.Scan(&turn.Id, &turn.DateUp, &turn.Hour, &turn.Description, &turn.Patient.Id, &turn.Dentist.Id); err != nil {
		return domain.Turns{}, err
	}
	return turn, nil
}

func (s *store) CreateTurn(turnDTO dto.TurnInsert) error {
	st, err := s.db.Prepare("INSERT INTO turns (date, hour, description, patientid, dentistid) VALUES (?, ?, ?, ? ,?)")
	if err != nil {
		return err
	}

	res, err := st.Exec(turnDTO.DateUp, turnDTO.Hour, turnDTO.Description, turnDTO.PatientId, turnDTO.DentistId)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	defer st.Close()
	return nil
}

func (s *store) UpdateTurn(id int, turnDTO dto.TurnInsert) error {
	st, err := s.db.Prepare("UPDATE turns SET date = ?, hour = ?, description = ?, patientid = ?, dentistid = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer st.Close()

	res, err := st.Exec(&turnDTO.DateUp, &turnDTO.Hour, &turnDTO.Description, &turnDTO.PatientId, &turnDTO.DentistId, id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
func (s *store) CreateTurnByDniAndEnrollment(turn dto.TurnPost) (dto.TurnInsert, error) {
	var idpatient int
	row := s.db.QueryRow("SELECT id FROM patient WHERE dni=?", turn.PatientDni)
	if err := row.Scan(&idpatient); err != nil {
		return dto.TurnInsert{}, errors.New("The patient doest exists.")
	}
	var iddentist int
	row2 := s.db.QueryRow("SELECT id FROM dentist WHERE enrollment=?", turn.DentistEnrollment)
	if err := row2.Scan(&iddentist); err != nil {
		return dto.TurnInsert{}, errors.New("The dentist doest exists.")
	}
	newturninsert := dto.TurnInsert{
		DateUp:      turn.DateUp,
		Hour:        turn.Hour,
		Description: turn.Description,
		PatientId:   idpatient,
		DentistId:   iddentist,
	}

	st, err := s.db.Prepare("INSERT INTO turns (date, hour, description, patientid, dentistid) VALUES (?, ?, ?, ? ,?)")
	if err != nil {
		return dto.TurnInsert{}, err
	}

	res, err := st.Exec(newturninsert.DateUp, newturninsert.Hour, newturninsert.Description, newturninsert.PatientId, newturninsert.DentistId)
	if err != nil {
		return dto.TurnInsert{}, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return dto.TurnInsert{}, err
	}
	return newturninsert, nil

}

func (s *store) DeleteTurn(id int) error {
	var idselect int
	row := s.db.QueryRow("SELECT id FROM turns WHERE id = ?", id)
	if err := row.Scan(&idselect); err != nil {
		return errors.New("The patient doest exists.")
	}
	query := "DELETE FROM turns WHERE id = ?"
	_, err := s.db.Exec(query, idselect)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) ReadTurnbyDNI(dni string) ([]dto.TurnGet, error) {
	var turn dto.TurnGet
	var list []dto.TurnGet
	rows, err := s.db.Query("SELECT t.date, t.hour, t.description, p.name, p.dni, d.name  FROM turns AS t JOIN patient AS p ON t.patientid = p.id JOIN dentist AS d ON t.dentistid = d.id WHERE p.dni = ?", dni)
	if err != nil {
		return []dto.TurnGet{}, err
	}
	for rows.Next() {

		if err := rows.Scan(&turn.Date, &turn.Hour, &turn.Description, &turn.PatientName, &turn.DNIPatient, &turn.DentistName); err != nil {
			return nil, err
		}
		list = append(list, turn)
	}
	rows.Close()
	return list, nil
}
