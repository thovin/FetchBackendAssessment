package main

import (
	"github.com/google/uuid"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type item struct {
	shortDescription string
	price            float64
}

type receipt struct {
	id           uuid.UUID
	retailer     string
	purchaseTime time.Time
	total        float64
	items        []item
	points       int
}

var receipts map[uuid.UUID]receipt = make(map[uuid.UUID]receipt)

func calculatePoints(r receipt) int {
	var points int

	for _, c := range strings.TrimSpace(r.retailer) { //one point per alphanumeric
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

func addReceipt(in inData) uuid.UUID {
	var r receipt

	r.retailer = in.Retailer
	dateLayout := "2006-01-02 15:04"
	r.purchaseTime, _ = time.Parse(dateLayout, in.PurchaseDate+" "+in.PurchaseTime)
	r.total, _ = strconv.ParseFloat(in.Total, 64)

	items := make([]item, len(in.Items))
	for i, itemIn := range in.Items {
		var item item
		item.shortDescription = itemIn.ShortDescription
		item.price, _ = strconv.ParseFloat(itemIn.Price, 64)
		items[i] = item
	}

	r.items = items
	r.id = uuid.New()
	r.points = calculatePoints(r)

	receipts[r.id] = r
	return r.id

}
