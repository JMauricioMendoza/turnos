package middleware

import (
	"net/http"
	"turnos-api/database"
	"turnos-api/utils"

	"github.com/gin-gonic/gin"
)

func AutenticacionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			utils.RespuestaJSON(c, http.StatusUnauthorized, "Token requerido")
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		var usuarioID int
		var nombreRol string

		query := `
            SELECT s.usuario_id, r.nombre
            FROM sesiones s
            JOIN usuarios u ON s.usuario_id = u.id
			JOIN roles r ON u.rol_id = r.id
            WHERE s.token = $1 AND s.expira_en > NOW() AND u.estatus IS TRUE
        `
		err := database.DB.QueryRow(query, token).Scan(&usuarioID, &nombreRol)

		if err != nil {
			utils.RespuestaJSON(c, http.StatusUnauthorized, "No se encontró una sesión válida")
			c.Abort()
			return
		}

		c.Set("usuario_id", usuarioID)
		c.Set("rol", nombreRol)
		c.Next()
	}
}
