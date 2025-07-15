package routes

import (
	"net/http"
	"time"
	"turnos-api/middleware"
	"turnos-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/Test", Test)
	r.POST("/IniciarSesion", IniciarSesion)
	r.GET("/CerrarSesion/:usuario_id", CerrarSesion)

	auth := r.Group("/")
	auth.Use(middleware.AutenticacionMiddleware())
	{
		//Sesiones
		auth.GET("/ObtenerUsuarioActual", ObtenerUsuarioActual)

		// Usuarios
		auth.POST("/CrearUsuario", CrearUsuario)
	}

	return r
}

func Test(c *gin.Context) {
	utils.RespuestaJSON(c, http.StatusOK, "Test endpoint funcionando correctamente")
}
