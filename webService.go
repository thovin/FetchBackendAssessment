package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// hosturl := "localhost:8080"	//why does it only work one way?
var hosturl string = "localhost:8080"

// receipts :-  make([]receipt, 1) //move to func (model) class?	//why does it only work one way?
var receipts map[uuid.UUID]receipt = make(map[uuid.UUID]receipt)

// var receipts map[uuid.UUID]receipt

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", getPoints)
	router.POST("/receipts/process", postProcessReceipts)

	router.Run(hosturl)
}

func postProcessReceipts(c *gin.Context) {
	var r receipt
	var in inData

	if err := c.BindJSON(&in); err != nil { //TODO do I need more input validation?
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
	}

	r.retailer = in.Retailer

	dateLayout := "2006-01-02 15:04"
	r.purchaseTime, _ = time.Parse(dateLayout, in.PurchaseDate+" "+in.PurchaseTime)

	// if temp, err := strconv.ParseFloat(in.Total, 64); err != nil {
	// 	log.Println(err)
	// } else {
	// 	r.total = temp
	// }
	r.total, _ = strconv.ParseFloat(in.Total, 64)

	items := make([]item, len(in.Items))
	for i, itemIn := range in.Items {
		var item item
		item.shortDescription = itemIn.ShortDescription
		// if temp, err := strconv.ParseFloat(itemIn.Price, 64); err != nil {
		// 	log.Println(err)
		// } else {
		// 	item.price = temp
		// }
		item.price, _ = strconv.ParseFloat(itemIn.Price, 64)
		items[i] = item
	}

	r.items = items

	r.id = uuid.New() //TODO how do I validate unique? Do I actually have to, or is only one receipt at a time handled?
	r.points = calculatePoints(r)

	receipts[r.id] = r
	c.IndentedJSON(http.StatusCreated, gin.H{"id": r.id})
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

func calculatePoints(r receipt) int { //move to func class?
	var points int

	for _, c := range strings.TrimSpace(r.retailer) { //one points per alphnumeric
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			points++
		}
	}

	totalCents := int(r.total * 100)
	if totalCents%100 == 0 {
		points += 50
	}
	if totalCents%25 == 0 {
		points += 25
	}

	points += int((len(r.items))/2) * 5

	for _, item := range r.items { //short description points
		if len(strings.TrimSpace(item.shortDescription))%3 == 0 {
			points += int(math.Ceil(item.price * .2))
		}
	}

	if r.purchaseTime.Day()%2 != 0 {
		points += 6
	}

	if r.purchaseTime.Hour() >= 14 && r.purchaseTime.Hour() < 16 { //assumes "after 2:00pm" is inclusive
		points += 10
	}

	return points

}
