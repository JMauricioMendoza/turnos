package models

import (
	"github.com/go-playground/validator/v10"
)

type Turno struct {
	ID                          int     `json:"id"`
	NumeroTurno                 int     `json:"numero_turno" validate:"gte=0"`
	ActividadID                 int     `json:"actividad_id"`
	ActividadNombre             string  `json:"actividad_nombre"`
	TiempoRecepcion             string  `json:"tiempo_recepcion"`
	TiempoInicioAtencion        *string `json:"tiempo_inicio_atencion"`
	TiempoFinAtencion           *string `json:"tiempo_fin_atencion"`
	UsuarioRecepcionNombre      string  `json:"usuario_recepcion"`
	UsuarioInicioAtencionNombre *string `json:"usuario_inicio_atencion"`
	UsuarioFinAtencionNombre    *string `json:"usuario_fin_atencion"`
}

func (u *Turno) ValidarTurno(Turno) error {
	validate := validator.New()
	return validate.Struct(u)
}
