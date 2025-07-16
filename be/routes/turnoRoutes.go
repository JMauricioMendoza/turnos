package routes

import (
	"fmt"
	"net/http"
	"turnos-api/database"
	"turnos-api/models"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
)

func CrearTurno(c *gin.Context) {
	usuarioID, existe := c.Get("usuario_id")
	if !existe {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Usuario no autenticado.")
		return
	}

	nombreRol, existe := c.Get("rol")
	if !existe || (nombreRol != "Root" && nombreRol != "Recepcionista") {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para esta acción.")
		return
	}

	var turno models.Turno

	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err := c.ShouldBindJSON(&turno); err != nil {
		utils.RespuestaJSON(c, http.StatusBadRequest, "Los datos proporcionados no son válidos.")
		return
	}

	if err = turno.ValidarTurno(turno); err != nil {
		utils.RespuestaJSON(c, http.StatusBadRequest, "Los datos proporcionados no cumplen con los requisitos de validación.")
		return
	}

	query := "INSERT INTO turnos (numero_turno, actividad_id, usuario_recepcion_id) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, turno.NumeroTurno, turno.ActividadID, usuarioID)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespuestaJSON(c, http.StatusCreated, fmt.Sprintf("Turno #%d creado exitosamente.", turno.NumeroTurno))
}
