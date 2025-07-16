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

func ObtenerTurnosEnRecepcion(c *gin.Context) {
	usuarioID, existe := c.Get("usuario_id")
	if !existe {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Usuario no autenticado.")
		return
	}

	nombreRol, existe := c.Get("rol")
	if !existe || (nombreRol != "Root" && nombreRol != "Atención en módulo") {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para esta acción.")
		return
	}

	var actividadUsuario int

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

	query := "SELECT actividad_id FROM usuarios WHERE id = $1"
	err = tx.QueryRow(query, usuarioID).Scan(&actividadUsuario)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	rows, err := database.DB.Query(`
		SELECT
			t.id,
			t.numero_turno,
			a.nombre,
			t.tiempo_recepcion
		FROM
			turnos t
			INNER JOIN actividades a ON t.actividad_id = a.id
		WHERE
			a.id = $1
			AND t.tiempo_recepcion::DATE = NOW()::DATE
			AND t.tiempo_inicio_atencion IS NULL
			AND t.tiempo_fin_atencion IS NULL
		ORDER BY
			t.tiempo_recepcion ASC`,
		actividadUsuario)

	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var turnos []models.Turno
	for rows.Next() {
		var turno models.Turno
		if err := rows.Scan(&turno.ID, &turno.NumeroTurno, &turno.ActividadNombre, &turno.TiempoRecepcion); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Turnos obtenidos exitosamente.", turnos)
}

func ObtenerTurnosEnAtencion(c *gin.Context) {
	usuarioID, existe := c.Get("usuario_id")
	if !existe {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Usuario no autenticado.")
		return
	}

	nombreRol, existe := c.Get("rol")
	if !existe || (nombreRol != "Root" && nombreRol != "Atención en módulo") {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para esta acción.")
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

	rows, err := database.DB.Query(`
		SELECT
			t.id,
			t.numero_turno,
			a.nombre,
			t.tiempo_recepcion
		FROM
			turnos t
			INNER JOIN actividades a ON t.actividad_id = a.id
		WHERE
			t.usuario_inicio_atencion_id = $1
			AND t.tiempo_recepcion IS NOT NULL
			AND t.tiempo_inicio_atencion::DATE = NOW()::DATE
			AND t.tiempo_fin_atencion IS NULL
		ORDER BY
			t.numero_turno ASC`,
		usuarioID)

	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var turnos []models.Turno
	for rows.Next() {
		var turno models.Turno
		if err := rows.Scan(&turno.ID, &turno.NumeroTurno, &turno.ActividadNombre, &turno.TiempoRecepcion); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Turnos obtenidos exitosamente.", turnos)
}

func ObtenerTurnosTodos(c *gin.Context) {
	nombreRol, existe := c.Get("rol")
	if !existe || (nombreRol != "Root" && nombreRol != "Recepcionista") {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para esta acción.")
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

	rows, err := database.DB.Query(`
		SELECT
			t.id,
			t.numero_turno,
			a.nombre,
			t.tiempo_recepcion,
			t.tiempo_inicio_atencion,
			t.tiempo_fin_atencion,
			ur.nombre_completo,
			ui.nombre_completo,
			uf.nombre_completo
		FROM
			turnos t
			INNER JOIN actividades a ON t.actividad_id = a.id
			INNER JOIN usuarios ur ON t.usuario_recepcion_id = ur.id
			LEFT JOIN usuarios ui ON t.usuario_inicio_atencion_id = ui.id
			LEFT JOIN usuarios uf ON t.usuario_fin_atencion_id = uf.id
		WHERE
			t.tiempo_recepcion::DATE = NOW()::DATE
		ORDER BY
			t.numero_turno ASC`)

	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var turnos []models.Turno
	for rows.Next() {
		var turno models.Turno
		if err := rows.Scan(
			&turno.ID,
			&turno.NumeroTurno,
			&turno.ActividadNombre,
			&turno.TiempoRecepcion,
			&turno.TiempoInicioAtencion,
			&turno.TiempoFinAtencion,
			&turno.UsuarioRecepcionNombre,
			&turno.UsuarioInicioAtencionNombre,
			&turno.UsuarioFinAtencionNombre,
		); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Turnos obtenidos exitosamente.", turnos)
}
