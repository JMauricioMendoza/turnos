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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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

	query := "INSERT INTO turnos (numero_turno, actividad_id, usuario_recepcion_id, estatus) VALUES ($1, $2, $3, 'recepcion')"
	_, err = tx.Exec(query, turno.NumeroTurno, turno.ActividadID, usuarioID)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al insertar turno: "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error en el commit: "+err.Error())
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al obtener actividad del usuario: "+err.Error())
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
			AND t.estatus = 'recepcion'
		ORDER BY
			t.tiempo_recepcion ASC`,
		actividadUsuario)

	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al consultar turnos en recepción: "+err.Error())
		return
	}
	defer rows.Close()

	var turnos []models.Turno
	for rows.Next() {
		var turno models.Turno
		if err := rows.Scan(&turno.ID, &turno.NumeroTurno, &turno.ActividadNombre, &turno.TiempoRecepcion); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al escanear turno: "+err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iterar sobre los resultados: "+err.Error())
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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
			AND t.tiempo_inicio_atencion::DATE = NOW()::DATE
			AND t.estatus = 'atencion'
		ORDER BY
			t.numero_turno ASC`,
		usuarioID)

	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al consultar turnos en atención: "+err.Error())
		return
	}
	defer rows.Close()

	var turnos []models.Turno
	for rows.Next() {
		var turno models.Turno
		if err := rows.Scan(&turno.ID, &turno.NumeroTurno, &turno.ActividadNombre, &turno.TiempoRecepcion); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al escanear turno: "+err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iterar sobre los resultados: "+err.Error())
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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
			uf.nombre_completo,
			t.estatus
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al consultar todos los turnos: "+err.Error())
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
			&turno.Estatus,
		); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al escanear turno: "+err.Error())
			return
		}
		turnos = append(turnos, turno)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iterar sobre los resultados: "+err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Turnos obtenidos exitosamente.", turnos)
}

func LlamarTurno(c *gin.Context) {
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

	var turno models.Turno

	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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

	query := `
		UPDATE
			turnos
		SET
			tiempo_inicio_atencion = NOW(),
			usuario_inicio_atencion_id = $1,
			estatus = 'atencion',
			actualizado_en = NOW()
		WHERE
			id = $2
		RETURNING
			numero_turno`
	err = tx.QueryRow(query, usuarioID, turno.ID).Scan(&turno.NumeroTurno)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al actualizar turno: "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error en el commit: "+err.Error())
		return
	}

	utils.RespuestaJSON(c, http.StatusOK, fmt.Sprintf("El turno #%d ha sido llamado.", turno.NumeroTurno))
}

func ConcluirTurno(c *gin.Context) {
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

	var turno models.Turno

	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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

	query := `
        UPDATE
            turnos
        SET
            tiempo_fin_atencion = NOW(),
            usuario_fin_atencion_id = $1,
			estatus = 'concluido',
            actualizado_en = NOW()
        WHERE
            id = $2
        RETURNING
            numero_turno`
	err = tx.QueryRow(query, usuarioID, turno.ID).Scan(&turno.NumeroTurno)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al actualizar turno: "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error en el commit: "+err.Error())
		return
	}

	utils.RespuestaJSON(c, http.StatusOK, fmt.Sprintf("El turno #%d ha sido concluido.", turno.NumeroTurno))
}

func EditarTurno(c *gin.Context) {
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
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iniciar transacción en BD: "+err.Error())
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

	var query string

	switch turno.Estatus {
	case "recepcion":
		query = `
        UPDATE
            turnos 
        SET
            estatus = 'recepcion',
            tiempo_recepcion = NOW(),
            actualizado_en = NOW(),			
            actividad_id = $1,
            usuario_recepcion_id = $2
        WHERE id = $3
		RETURNING
			numero_turno
        `
	case "concluido":
		query = `
        UPDATE
            turnos 
        SET
            estatus = 'concluido',
            tiempo_fin_atencion = NOW(),
            actualizado_en = NOW(),			
            actividad_id = $1,
            usuario_fin_atencion_id = $2
        WHERE id = $3
		RETURNING
			numero_turno
        `
	default:
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusBadRequest, "El estatus del turno no es válido para editar.")
		return
	}

	err = tx.QueryRow(query, turno.ActividadID, usuarioID, turno.ID).Scan(&turno.NumeroTurno)
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al actualizar turno: "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error en el commit: "+err.Error())
		return
	}

	utils.RespuestaJSON(c, http.StatusOK, fmt.Sprintf("Se ha editado el turno #%d correctamente.", turno.NumeroTurno))
}
