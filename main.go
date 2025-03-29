package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
    "time"
)

var db *gorm.DB

// Definición de la estructura Match que representa un partido
type Match struct {
    ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"` // ID único del partido
    HomeTeam     string    `json:"homeTeam" gorm:"type:varchar(255);not null"` // Equipo local
    AwayTeam     string    `json:"awayTeam" gorm:"type:varchar(255);not null"` // Equipo visitante
    MatchDate    time.Time `json:"matchDate" gorm:"type:date;not null"` // Fecha del partido
    Goals        int       `json:"goals" gorm:"default:0"` // Goles del partido
    YellowCards  int       `json:"yellowCards" gorm:"default:0"` // Tarjetas amarillas
    RedCards     int       `json:"redCards" gorm:"default:0"` // Tarjetas rojas
    ExtraTime    int       `json:"extraTime" gorm:"default:0"` // Tiempo extra
}

func main() {
    // Configuración de la base de datos y conexión usando GORM
    dsn := "host=db user=postgres password=postgres dbname=matches port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("No se pudo conectar a la base de datos") // Error si no se puede conectar a la DB
    }
    db = database
    db.AutoMigrate(&Match{}) // Migrar la estructura de la DB (tabla de partidos)

    r := gin.Default()

    // Configuración de CORS para permitir peticiones de cualquier origen
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Definición de las rutas y sus correspondientes handlers
    r.GET("/api/matches", GetMatches)                 // Obtener todos los partidos
    r.GET("/api/matches/:id", GetMatch)               // Obtener un partido por su ID
    r.POST("/api/matches", CreateMatch)               // Crear un nuevo partido
    r.PUT("/api/matches/:id", UpdateMatch)            // Actualizar un partido existente
    r.DELETE("/api/matches/:id", DeleteMatch)         // Eliminar un partido

    // Endpoints adicionales para actualizar detalles específicos del partido
    r.PATCH("/api/matches/:id/goals", UpdateMatchGoals)       // Actualizar el número de goles
    r.PATCH("/api/matches/:id/yellowcards", AddYellowCard)    // Añadir tarjeta amarilla
    r.PATCH("/api/matches/:id/redcards", AddRedCard)          // Añadir tarjeta roja
    r.PATCH("/api/matches/:id/extratime", UpdateExtraTime)    // Actualizar tiempo extra

    r.Run(":8080") // Iniciar el servidor en el puerto 8080
}

// Obtener todos los partidos desde la base de datos
func GetMatches(c *gin.Context) {
    var matches []Match
    db.Find(&matches) // Consultar todos los partidos
    c.JSON(http.StatusOK, matches) // Retornar los partidos en formato JSON
}

// Obtener un partido por su ID
func GetMatch(c *gin.Context) {
    var match Match
    // Buscar el partido por su ID
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si no se encuentra el partido
        return
    }
    c.JSON(http.StatusOK, match) // Retornar el partido encontrado
}

// Crear un nuevo partido en la base de datos
func CreateMatch(c *gin.Context) {
    var input struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"` // Recibimos la fecha como string para convertirla después
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Error si no se puede parsear el JSON
        return
    }

    // Convertir la fecha de string a tipo time.Time
    parsedDate, err := time.Parse("2006-01-02", input.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."}) // Error si la fecha no tiene el formato correcto
        return
    }

    // Crear el partido en la base de datos
    match := Match{
        HomeTeam:  input.HomeTeam,
        AwayTeam:  input.AwayTeam,
        MatchDate: parsedDate,
    }

    db.Create(&match) // Insertar el partido en la base de datos

    c.JSON(http.StatusCreated, match) // Retornar el partido recién creado
}

// Actualizar un partido existente en la base de datos
func UpdateMatch(c *gin.Context) {
    var match Match
    id := c.Param("id")

    // Buscar el partido por ID
    if err := db.First(&match, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    // Estructura para recibir los datos del partido desde el JSON
    var input struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"` // Fecha como string
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Error si no se puede parsear el JSON
        return
    }

    // Parsear la fecha de string a tipo time.Time
    parsedDate, err := time.Parse("2006-01-02", input.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"}) // Error si la fecha no tiene el formato correcto
        return
    }

    // Actualizar el partido con los nuevos datos
    db.Model(&match).Updates(Match{
        HomeTeam:  input.HomeTeam,
        AwayTeam:  input.AwayTeam,
        MatchDate: parsedDate,
    })

    c.JSON(http.StatusOK, match) // Retornar el partido actualizado
}

// Eliminar un partido de la base de datos
func DeleteMatch(c *gin.Context) {
    var match Match
    id := c.Param("id")

    // Buscar el partido por ID
    if err := db.First(&match, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    db.Delete(&match) // Eliminar el partido de la base de datos
    c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"}) // Confirmación de eliminación
}

// Endpoint para añadir una tarjeta amarilla al partido
func AddYellowCard(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    // Incrementar el contador de tarjetas amarillas
    db.Model(&match).Update("YellowCards", match.YellowCards+1)

    // Recargar los datos actualizados
    db.First(&match, c.Param("id"))

    c.JSON(http.StatusOK, match) // Retornar los datos actualizados del partido
}

// Endpoint para añadir una tarjeta roja al partido
func AddRedCard(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    // Incrementar el contador de tarjetas rojas
    db.Model(&match).Update("RedCards", match.RedCards+1)

    // Recargar los datos actualizados
    db.First(&match, c.Param("id"))

    c.JSON(http.StatusOK, match) // Retornar los datos actualizados del partido
}

// Endpoint para añadir tiempo extra al partido
func UpdateExtraTime(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    // Incrementar el tiempo extra
    db.Model(&match).Update("ExtraTime", match.ExtraTime+1)

    // Recargar los datos actualizados
    db.First(&match, c.Param("id"))

    c.JSON(http.StatusOK, match) // Retornar los datos actualizados del partido
}

// Endpoint para incrementar los goles del partido
func UpdateMatchGoals(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"}) // Error si el partido no se encuentra
        return
    }

    // Incrementar el número total de goles en 1
    db.Model(&match).Update("goals", match.Goals+1)

    // Recargar los datos actualizados
    db.First(&match, c.Param("id"))

    c.JSON(http.StatusOK, match) // Retornar los datos actualizados del partido
}
