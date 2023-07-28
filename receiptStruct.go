package main

import ( //TODO do I need these since everything is one package?
	"github.com/google/uuid"
	"time"
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
