package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type inData struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string
	PurchaseTime string
	Total        string
	Items        []inItem
}

type inItem struct {
	ShortDescription string
	Price            string
}

var hosturl string = "localhost:8080"

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", getPoints)
	router.POST("/receipts/process", postProcessReceipts)

	router.Run(hosturl)
}

func postProcessReceipts(c *gin.Context) {
	var in inData

	if err := c.BindJSON(&in); err != nil { //TODO do I need more input validation?
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"id": addReceipt(in)})
}

func getPoints(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	if r, exists := receipts[id]; exists {
		c.IndentedJSON(http.StatusOK, gin.H{"points": r.points})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id not found"})
		return
	}

}
