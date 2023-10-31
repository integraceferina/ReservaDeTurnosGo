package handler

import (
	"errors"
	"os"
	"strconv"

	"reserva/internal/domain/dto"
	"reserva/internal/turns"
	"reserva/pkg/web"

	"github.com/gin-gonic/gin"
)

type turnhandler struct {
	s turns.Service
}

func NewTurnHandler(s turns.Service) *turnhandler {
	return &turnhandler{s: s}
}

// ListTurns godoc
// @Summary Get All Turns
// @Tags Turns
// @Description Get Turns
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Failure 400 {object} web.response
// @Router /turns [get]
func (h *turnhandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		turns, _ := h.s.ReadAll()
		if len(turns) == 0 {
			web.Failure(c, 400, errors.New("There is no turn"))
		}
		web.Success(c, 200, turns)
	}
}

// PostTurns godoc
// @Summary Post Turns
// @Tags Turns
// @Description Post a Turns
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.TurnInsert true "Turns to data"
// @Success 201 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /turns [post]
func (h *turnhandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}

		var turn dto.TurnInsert
		err := c.ShouldBindJSON(&turn)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid turn"))
			return
		}

		valid, err := validateEmptysTurn(&turn)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		err = h.s.Create(turn)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, turn)
	}
}

// TurnsbyID godoc
// @Summary Get Turns by ID
// @Tags Turns
// @Description Get Turns
// @Accept json
// @Produce json
// @Param id path int true "Turns ID"
// @Success 200 {object} web.response
// @Failure 404 {object} web.response
// @Router /turns/{id} [get]
func (h *turnhandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 404, errors.New("Invalid id"))
			return
		}
		turn, err := h.s.Read(id)
		if err != nil {
			web.Failure(c, 404, errors.New("Turn not found"))
			return
		}
		web.Success(c, 200, turn)
	}
}

// UpdateTurns godoc
// @Summary Put Turns
// @Tags Turns
// @Description Put a Turns
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.TurnInsert true "Turns to data"
// @Param id path int true "Turns ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /turns/{id} [put]
func (h *turnhandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		var turn dto.TurnInsert
		err := c.ShouldBindJSON(&turn)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid patient"))
			return
		}
		valid, err := validateEmptysTurn(&turn)
		if !valid {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		err = h.s.Update(id, turn)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Patient Updated")
	}
}

// UpdateTurns godoc
// @Summary Patch Turns
// @Tags Turns
// @Description Patch a Turns
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.TurnInsert true "Turns to data"
// @Param id path int true "Turns ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Failure 409 {object} web.response
// @Router /turns/{id} [patch]
func (h *turnhandler) Patch() gin.HandlerFunc {
	type TurnRequest struct {
		PatientId   int    `json:"patientid,omitempty"`
		DentistId   int    `json:"dentistid,omitempty"`
		DateUp      string `json:"dateup,omitempty"`
		Hour        string `json:"hour,omitempty"`
		Description string `json:"description,omitempty"`
	}
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		var r TurnRequest
		if err := c.ShouldBindJSON(&r); err != nil {
			web.Failure(c, 400, errors.New("Invalid Request"))
			return
		}
		update := dto.TurnInsert{
			DentistId:   r.DentistId,
			PatientId:   r.PatientId,
			DateUp:      r.DateUp,
			Hour:        r.Hour,
			Description: r.Description,
		}
		err = h.s.Update(id, update)
		if err != nil {
			web.Failure(c, 409, errors.New(err.Error()))
			return
		}
		web.Success(c, 200, "Turn Updated")
	}
}

// PostTurnsWithEnrollmentAndDNI godoc
// @Summary PostWithEnrollmentAndDNI Turns
// @Tags Turns
// @Description Post a Turn
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param dentist body dto.TurnPost true "Turns to data"
// @Success 201 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /turns/post [post]
func (h *turnhandler) PostxEnrollmentAndDni() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}

		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}

		var turn dto.TurnPost
		err := c.ShouldBindJSON(&turn)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid turn"))
			return
		}

		if turn.DentistEnrollment == "" || turn.PatientDni == "" {
			web.Failure(c, 400, errors.New("We need DNI and Enrollment to create a Turn"))
			return
		}
		r, err := h.s.CreateTurnByDniAndEnrollment(turn)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 201, r)
	}
}

// DeleteTurns godoc
// @Summary Delete Turns
// @Tags Turns
// @Description Delete a Turns
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id path int true "Turns ID"
// @Success 200 {object} web.response
// @Failure 401 {object} web.response
// @Failure 400 {object} web.response
// @Router /turns/{id} [delete]
func (h *turnhandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			web.Failure(c, 401, errors.New("Token Not Found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("Invalid Token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid ID"))
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.Failure(c, 400, errors.New(err.Error()))
			return
		}

		web.Success(c, 200, "Turn Deleted")
	}
}

// TurnsbyID godoc
// @Summary Get Turns by ID
// @Tags Turns
// @Description Get Turns
// @Accept json
// @Produce json
// @Param dni query string false "DNI"
// @Success 200 {object} web.response
// @Failure 404 {object} web.response
// @Router /turns/dni [get]
func (h *turnhandler) GetByDNI() gin.HandlerFunc {
	return func(c *gin.Context) {
		dni := c.Query("dni")
		turns, err := h.s.ReadbyDni(dni)
		if err != nil {
			web.Failure(c, 404, errors.New("Turn not found"))
			return
		}
		web.Success(c, 200, turns)
	}
}

func validateEmptysTurn(turn *dto.TurnInsert) (bool, error) {
	if turn.Description == "" || turn.Hour == "" || turn.DateUp == "" || turn.DentistId == 0 || turn.PatientId == 0 {
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
