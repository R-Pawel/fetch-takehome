package routes

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/reciepts/process", processReciept)
	router.GET("/reciepts/:id/points", sendPoints)
	return router
}

type reciept struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

type id struct {
	ID string `json:"id"`
}

type points struct {
	Points int64 `json:"points"`
}

type item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var recieptsByID = make(map[string]int64)

func processReciept(c *gin.Context) {
	var newReciept reciept
	var pointTotal int64 = 0
	if err := c.BindJSON(&newReciept); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	// missing data check
	if strings.Trim(newReciept.Retailer, " ") == "" || len(newReciept.Items) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	// date check
	purchaseDate, err := time.Parse("2006-01-02", newReciept.PurchaseDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "The receipt is invalid")
		return
	}
	purchaseDay := purchaseDate.Day()
	if purchaseDay%2 == 1 {
		pointTotal += 6
	}

	// time check
	purchaseTime, err := time.Parse("15:04", newReciept.PurchaseTime)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "The reciept is invalid")
		return
	}
	early, _ := time.Parse("15:04", "14:00")
	late, _ := time.Parse("15:04", "16:00")
	if purchaseTime.Before(late) && purchaseTime.After(early) {
		pointTotal += 10
	}

	// check total
	if !(regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(newReciept.Total)) {
		c.IndentedJSON(http.StatusBadRequest, "The reciept is invalid")
		return
	}
	totalParts := strings.Split(newReciept.Total, ".")
	totalCents := totalParts[1]
	if totalCents == "00" {
		pointTotal += 75
	}
	if totalCents == "25" || totalCents == "50" || totalCents == "75" {
		pointTotal += 25
	}

	//points for pairs
	pointTotal += (int64)(len(newReciept.Items)/2) * 5

	//check items
	for _, item := range newReciept.Items {
		description := strings.Trim(item.ShortDescription, " ")
		if !(regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price)) {
			c.IndentedJSON(http.StatusBadRequest, "The reciept is invalid")
			return
		}
		if len(description) == 0 {
			c.IndentedJSON(http.StatusBadRequest, "The reciept is invalid")
			return
		}
		numberPrice, err := strconv.ParseFloat(item.Price, 64)
		if len(description)%3 == 0 {
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, "The reciept is invalid")
				return
			}
			pointTotal += (int64)(math.Ceil(numberPrice * .2))
		}
	}
	//alphanumeric characters
	alphanumeric := 0

	for _, char := range newReciept.Retailer {
		if unicode.IsDigit(char) || unicode.IsLetter(char) {
			alphanumeric += 1
		}
	}
	pointTotal += int64(alphanumeric)

	newID := uuid.New()
	returnID := id{ID: newID.String()}
	recieptsByID[newID.String()] = pointTotal

	c.IndentedJSON(http.StatusOK, returnID)
}

func sendPoints(c *gin.Context) {
	checkID := c.Param("id")
	if recieptsByID[checkID] != 0 {
		returnPoints := points{Points: recieptsByID[checkID]}
		c.IndentedJSON(http.StatusOK, returnPoints)
	} else {
		c.IndentedJSON(http.StatusNotFound, "No receipt found for that id")
	}
}
