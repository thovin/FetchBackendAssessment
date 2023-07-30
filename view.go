package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
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
	dateLayout := "2006-01-02 15:04"

	if err := c.BindJSON(&in); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	if in.Retailer == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no retailer"})
		return
	} else if _, err := time.Parse(dateLayout, in.PurchaseDate+" "+in.PurchaseTime); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid purchase date or time"})
		return
	} else if _, err := strconv.ParseFloat(in.Total, 64); in.Total == "" || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "total missing or invalid"})
		return
	} else {
		for _, item := range in.Items {
			if item.ShortDescription == "" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "item missing description"})
				return
			} else if _, err := strconv.ParseFloat(item.Price, 64); item.Price == "" || err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "item: '" + item.ShortDescription + "' missing or invalid price"})
				return
			}
		}
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
