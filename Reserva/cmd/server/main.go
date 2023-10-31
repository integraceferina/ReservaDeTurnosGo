package main

/*
@title Api Dentist Office by Agustin Vanetta
@version 1.0
@description This API Handle Dentist Office.
@termsOfService  https://github.com/agvanetta
@contact.name API Support
@contact.url https://github.com/agvanetta
@license.name Apache 2.0
@license.url https://www.apache.org/licenses/LICENSE-2.0.html

*/
import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"reserva/cmd/server/handler"
	"reserva/docs"
	"reserva/internal/dentist"
	"reserva/internal/patient"
	"reserva/internal/turns"
	"reserva/pkg/middleware"
	"reserva/pkg/store"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("user"),
			os.Getenv("pass"),
			os.Getenv("hostdb"),
			os.Getenv("port"),
			os.Getenv("db_name"))
	)
	fmt.Print(ConnectionString)

	db, err := sql.Open("mysql", ConnectionString)
	if err != nil {
		log.Fatal("Error opening database")
	}
	// Store
	storeSQL := store.NewSQLStore(db)

	// Dentist
	repoDentist := dentist.NewRepository(storeSQL)
	serviceDentist := dentist.NewService(repoDentist)
	handlerDentist := handler.NewDentistHandler(serviceDentist)

	r := gin.Default()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	dentists := r.Group("/dentists")
	{
		dentists.GET("", handlerDentist.GetAll())
		dentists.GET(":id", handlerDentist.GetByID())
		dentists.POST("", middleware.AuthenticationMiddleware(), handlerDentist.Post())
		dentists.DELETE(":id", middleware.AuthenticationMiddleware(), handlerDentist.Delete())
		dentists.PUT(":id", middleware.AuthenticationMiddleware(), handlerDentist.Put())
		dentists.PATCH(":id", middleware.AuthenticationMiddleware(), handlerDentist.Patch())

	}
	// Patient
	repoPatient := patient.NewRepository(storeSQL)
	servicePatient := patient.NewService(repoPatient)
	handlerPatient := handler.NewPatientHandler(servicePatient)

	patients := r.Group("/patients")
	{
		patients.GET("", handlerPatient.GetAll())
		patients.GET(":id", handlerPatient.GetByID())
		patients.POST("", middleware.AuthenticationMiddleware(), handlerPatient.Post())
		patients.DELETE(":id", middleware.AuthenticationMiddleware(), handlerPatient.Delete())
		patients.PUT(":id", middleware.AuthenticationMiddleware(), handlerPatient.Put())
		patients.PATCH(":id", middleware.AuthenticationMiddleware(), handlerPatient.Patch())
	}
	// Turns
	repoTurns := turns.NewRepository(storeSQL)
	serviceTurns := turns.NewService(repoTurns)
	handlerTurns := handler.NewTurnHandler(serviceTurns)

	turns := r.Group("/turns")
	{
		turns.GET("", handlerTurns.GetAll())
		turns.POST("", middleware.AuthenticationMiddleware(), handlerTurns.Post())
		turns.GET(":id", handlerTurns.GetByID())
		turns.PUT(":id", middleware.AuthenticationMiddleware(), handlerTurns.Put())
		turns.PATCH(":id", middleware.AuthenticationMiddleware(), handlerTurns.Patch())
		turns.POST("/post", middleware.AuthenticationMiddleware(), handlerTurns.PostxEnrollmentAndDni())
		turns.DELETE(":id", middleware.AuthenticationMiddleware(), handlerTurns.Delete())
		turns.GET("/dni", middleware.AuthenticationMiddleware(), handlerTurns.GetByDNI())
	}
	r.Run(":8080")
}
