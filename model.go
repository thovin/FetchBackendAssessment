package main

import (
	"github.com/google/uuid"
	"math"
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

var receipts map[uuid.UUID]receipt = make(map[uuid.UUID]receipt)

func calculatePoints(r receipt) int {
	var points int

	for _, c := range strings.TrimSpace(r.retailer) { //one point per alphnumeric
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
