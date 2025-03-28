package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
)

var db *gorm.DB

type Match struct {
    ID     uint   `json:"id" gorm:"primaryKey"`
    TeamA  string `json:"team_a"`
    TeamB  string `json:"team_b"`
    ScoreA int    `json:"score_a"`
    ScoreB int    `json:"score_b"`
}

func main() {
    dsn := "host=db user=postgres password=postgres dbname=matches port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("No se pudo conectar a la base de datos")
    }
    db = database
    db.AutoMigrate(&Match{})

    r := gin.Default()

    r.GET("/api/matches", GetMatches)
    r.GET("/api/matches/:id", GetMatch)
    r.POST("/api/matches", CreateMatch)
    r.PUT("/api/matches/:id", UpdateMatch)
    r.DELETE("/api/matches/:id", DeleteMatch)

    r.Run(":8080")
}

func GetMatches(c *gin.Context) {
    var matches []Match
    db.Find(&matches)
    c.JSON(http.StatusOK, matches)
}

func GetMatch(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }
    c.JSON(http.StatusOK, match)
}

func CreateMatch(c *gin.Context) {
    var match Match
    if err := c.ShouldBindJSON(&match); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&match)
    c.JSON(http.StatusCreated, match)
}

func UpdateMatch(c *gin.Context) {
    var match Match
    if err := db.First(&match, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
        return
    }
    if err := c.ShouldBindJSON(&match); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Save(&match)
    c.JSON(http.StatusOK, match)
}

func DeleteMatch(c *gin.Context) {
    db.Delete(&Match{}, c.Param("id"))
    c.JSON(http.StatusNoContent, nil)
}
