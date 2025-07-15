package routes

import (
	"net/http"
	"strconv"
	"turnos-api/database"
	"turnos-api/models"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IniciarSesion(c *gin.Context) {
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
		utils.RespuestaJSON(c, http.StatusBadRequest, "Formato de datos no válidos.")
		return
	}

	passwordIngresada := usuario.Password

	query := "SELECT id, password FROM usuarios WHERE usuario = $1 AND estatus IS TRUE"
	err = tx.QueryRow(query, usuario.Usuario).Scan(&usuario.ID, &usuario.Password)
	if err != nil {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "El usuario no existe o está inhabilitado.")
		return
	}

	if !usuario.VerificarPassword(passwordIngresada) {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Contraseña incorrecta.")
		return
	}

	token := uuid.NewString()

	_, err = tx.Exec("DELETE FROM sesiones WHERE usuario_id = $1", usuario.ID)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = tx.Exec("INSERT INTO sesiones (usuario_id, token, expira_en) VALUES ($1, $2, (SELECT now() + interval '1 hour'))", usuario.ID, token)
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

	utils.RespuestaJSON(c, http.StatusOK, "Sesión iniciada correctamente.", token)
}

func ObtenerUsuarioActual(c *gin.Context) {
	usuarioID := c.GetInt("usuario_id")

	var nombre string
	err := database.DB.QueryRow(`
		SELECT nombre_completo
		FROM usuarios
		WHERE id = $1
	`, usuarioID).Scan(&nombre)

	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	respuesta := map[string]interface{}{
		"usuario_id": usuarioID,
		"nombre":     nombre,
	}

	utils.RespuestaJSON(c, http.StatusOK, "Sesión válida", respuesta)
}

func CerrarSesion(c *gin.Context) {
	usuarioId := c.Param("usuario_id")

	if _, err := strconv.Atoi(usuarioId); err != nil {
		utils.RespuestaJSON(c, http.StatusBadRequest, "ID de usuario inválido.")
		return
	}

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

	var existe bool
	queryCheck := "SELECT EXISTS (SELECT 1 FROM sesiones WHERE usuario_id = $1)"
	err = tx.QueryRow(queryCheck, usuarioId).Scan(&existe)
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !existe {
		utils.RespuestaJSON(c, http.StatusOK, "Sesión cerrada correctamente.")
		return
	}

	_, err = tx.Exec("DELETE FROM sesiones WHERE usuario_id = $1", usuarioId)
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

	utils.RespuestaJSON(c, http.StatusOK, "Sesión cerrada correctamente.")
}
