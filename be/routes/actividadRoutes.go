package routes

import (
	"net/http"
	"turnos-api/database"
	"turnos-api/models"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
)

func ObtenerActividadesActivas(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, nombre FROM actividades WHERE estatus IS TRUE ORDER BY nombre ASC")
	if err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al obtener las actividades: "+err.Error())
		return
	}
	defer rows.Close()

	var actividades []models.Actividad
	for rows.Next() {
		var actividad models.Actividad
		if err := rows.Scan(&actividad.ID, &actividad.Nombre); err != nil {
			utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al escanear la actividad: "+err.Error())
			return
		}
		actividades = append(actividades, actividad)
	}
	if err = rows.Err(); err != nil {
		utils.RespuestaJSON(c, http.StatusInternalServerError, "Error al iterar sobre las actividades: "+err.Error())
		return
	}
	utils.RespuestaJSON(c, http.StatusOK, "Actividades obtenidas exitosamente.", actividades)
}
