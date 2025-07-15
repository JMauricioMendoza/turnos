package routes

import (
	"fmt"
	"net/http"
	"turnos-api/database"
	"turnos-api/models"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
)

func CrearUsuario(c *gin.Context) {
	nombreRol, existe := c.Get("rol")
	if !existe || nombreRol != "Root" {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para crear usuarios.")
		return
	}

	var usuario models.Usuario

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

	if err := c.ShouldBindJSON(&usuario); err != nil {
		utils.RespuestaJSON(c, http.StatusBadRequest, "Los datos proporcionados no son válidos.")
		return
	}

	if err = usuario.ValidarUsuario(usuario); err != nil {
		utils.RespuestaJSON(c, http.StatusBadRequest, "Los datos proporcionados no cumplen con los requisitos de validación.")
		return
	}

	queryCheck := "SELECT EXISTS (SELECT 1 FROM usuarios WHERE usuario = $1)"
	err = tx.QueryRow(queryCheck, usuario.Usuario).Scan(&existe)
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	if existe {
		utils.RespuestaJSON(c, http.StatusConflict, "El usuario ya existe.")
		return
	}

	if err := usuario.HashearPassword(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	query := "INSERT INTO usuarios (usuario, nombre_completo, password, actividad_id, rol_id, mesa ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = tx.QueryRow(query, usuario.Usuario, usuario.NombreCompleto, usuario.Password, usuario.ActividadID, usuario.RolID, usuario.Mesa).Scan(&usuario.ID)
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

	utils.RespuestaJSON(c, http.StatusCreated, fmt.Sprintf("Usuario #%d creado exitosamente", usuario.ID))
}
