package models

import (
	"github.com/go-playground/validator/v10"
)

type Turno struct {
	ID              int    `json:"id"`
	NumeroTurno     int    `json:"numero_turno" validate:"gte=0"`
	ActividadID     int    `json:"actividad_id"`
	ActividadNombre string `json:"actividad_nombre"`
	TiempoRecepcion string `json:"tiempo_recepcion"`
}

func (u *Turno) ValidarTurno(Turno) error {
	validate := validator.New()
	return validate.Struct(u)
}
