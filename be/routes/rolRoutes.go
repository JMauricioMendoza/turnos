package routes

import (
	"net/http"
	"turnos-api/database"
	"turnos-api/models"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
)

func ObtenerRolesActivos(c *gin.Context) {
	nombreRol, existe := c.Get("rol")
	if !existe || nombreRol != "Root" {
		utils.RespuestaJSON(c, http.StatusUnauthorized, "Rol no autorizado para esta acción.")
		return
	}

	rows, err := database.DB.Query("SELECT id, nombre FROM roles WHERE estatus IS TRUE ORDER BY nombre ASC")
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var roles []models.Rol
	for rows.Next() {
		var rol models.Rol
		if err := rows.Scan(&rol.ID, &rol.Nombre); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		roles = append(roles, rol)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Roles obtenidos exitosamente.", roles)
}
