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

	r.GET("/test", Test)
	r.POST("/sesion/crear", IniciarSesion)
	r.DELETE("/sesion/eliminar/:usuario_id", CerrarSesion)

	auth := r.Group("/")
	auth.Use(middleware.AutenticacionMiddleware())
	{
		//Sesiones
		auth.GET("/sesion/obtener/", ObtenerUsuarioActual)

		// Usuarios
		auth.POST("/usuario/crear", CrearUsuario)

		// Turnos
		auth.POST("/turno/crear", CrearTurno)
		auth.GET("/turno/obtener/recepcion", ObtenerTurnosEnRecepcion)
		auth.GET("/turno/obtener/atencion", ObtenerTurnosEnAtencion)
		auth.GET("/turno/obtener/todos", ObtenerTurnosTodos)
		auth.PATCH("/turno/llamar", LlamarTurno)
		auth.PATCH("/turno/concluir", ConcluirTurno)
		auth.PATCH("/turno/editar", EditarTurno)

		// Roles
		auth.GET("/rol/obtener/activos", ObtenerRolesActivos)

		// Actividades
		auth.GET("/actividades/obtener/activos", ObtenerActividadesActivas)
	}

	return r
}

func Test(c *gin.Context) {
	utils.RespuestaJSON(c, http.StatusOK, "Test endpoint funcionando correctamente")
}
