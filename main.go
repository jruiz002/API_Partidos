package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
    "time"
)

var db *gorm.DB

// Definición de la estructura Match
type Match struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    TeamA     string    `json:"team_a"`
    TeamB     string    `json:"team_b"`
    MatchDate time.Time `json:"match_date"`
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
    var match Match
    if err := c.ShouldBindJSON(&match); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&match)
    c.JSON(http.StatusCreated, match)
}

// Actualizar un partido
func UpdateMatch(c *gin.Context) {
    var match Match
    id := c.Param("id")

    if err := db.First(&match, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }

    var updateData Match
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.Model(&match).Updates(updateData) // Actualiza solo los campos enviados
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
