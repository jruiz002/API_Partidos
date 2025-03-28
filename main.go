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

// Definición de la estructura Match
type Match struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    HomeTeam  string    `json:"homeTeam" gorm:"type:varchar(255);not null"`
    AwayTeam  string    `json:"awayTeam" gorm:"type:varchar(255);not null"`
    MatchDate time.Time `json:"matchDate" gorm:"type:date;not null"`
}

func main() {
    dsn := "host=db user=postgres password=postgres dbname=matches port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("No se pudo conectar a la base de datos")
    }
    db = database
    db.AutoMigrate(&Match{}) // Migrar la estructura de la DB

    r := gin.Default()

    // Configurar CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Rutas de la API
    r.GET("/api/matches", GetMatches)
    r.GET("/api/matches/:id", GetMatch)
    r.POST("/api/matches", CreateMatch)
    r.PUT("/api/matches/:id", UpdateMatch)
    r.DELETE("/api/matches/:id", DeleteMatch)

    r.Run(":8080") // Servidor en el puerto 8080
}

// Obtener todos los partidos
func GetMatches(c *gin.Context) {
    var matches []Match
    db.Find(&matches)
    c.JSON(http.StatusOK, matches)
}

// Obtener un partido por ID
func GetMatch(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }
    c.JSON(http.StatusOK, match)
}

// Crear un nuevo partido
func CreateMatch(c *gin.Context) {
    var input struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"` // Recibimos como string para procesarlo
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Convertir string a time.Time
    parsedDate, err := time.Parse("2006-01-02", input.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
        return
    }

    // Crear partido
    match := Match{
        HomeTeam:  input.HomeTeam,
        AwayTeam:  input.AwayTeam,
        MatchDate: parsedDate,
    }

    db.Create(&match)

    c.JSON(http.StatusCreated, match)
}

// Actualizar un partido
// Actualizar un partido
func UpdateMatch(c *gin.Context) {
    var match Match
    id := c.Param("id")

    // Buscar el partido en la base de datos
    if err := db.First(&match, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }

    // Definir una estructura para recibir los datos del JSON
    var input struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"` // Fecha como string
    }

    // Bind JSON a la estructura
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parsear la fecha de string a time.Time
    parsedDate, err := time.Parse("2006-01-02", input.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
        return
    }

    // Actualizar los datos en la base de datos
    db.Model(&match).Updates(Match{
        HomeTeam:  input.HomeTeam,
        AwayTeam:  input.AwayTeam,
        MatchDate: parsedDate,
    })

    c.JSON(http.StatusOK, match)
}


// Eliminar un partido
func DeleteMatch(c *gin.Context) {
    var match Match
    id := c.Param("id")

    if err := db.First(&match, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }

    db.Delete(&match)
    c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}