package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Post struct {
	Date     string    `json:"date"`
	Country  string    `json:"country"`
	Location string    `json:"location"`
	Activity string    `json:"activity"`
	Name     string    `json:"name"`
	Age      string    `json:"age"`
	Injury   string    `json:"injury"`
	Species  string    `json:"species"`
	Id       uuid.UUID `json:"id"`
}

type PostRequest struct {
	Date     string `form:"date" binding:"required"`
	Country  string `form:"country" binding:"required"`
	Location string `form:"location" binding:"required"`
	Activity string `form:"activity" binding:"required"`
	Name     string `form:"name" binding:"required"`
	Age      string `json:"age"`
	Injury   string `form:"injury" binding:"required"`
	Species  string `form:"species" binding:"required"`
}

type UpdatePostRequest struct {
	Date     string `form:"date"`
	Country  string `form:"country"`
	Location string `form:"location"`
	Activity string `form:"activity"`
	Name     string `form:"name"`
	Age      string `json:"age"`
	Injury   string `form:"injury"`
	Species  string `form:"species"`
}

var (
	data   []Post
	dataMu sync.Mutex
)

func main() {
	data = importData()
	data = getRandomItems(data, 10)

	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/shark", func(c *gin.Context) {
		dataMu.Lock()
		defer dataMu.Unlock()
		c.JSON(http.StatusOK, data)
	})

	r.POST("/shark", func(c *gin.Context) {
		var requestBody PostRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newPost := Post{
			Date:     requestBody.Date,
			Country:  requestBody.Country,
			Location: requestBody.Location,
			Activity: requestBody.Activity,
			Name:     requestBody.Name,
			Age:      requestBody.Age,
			Injury:   requestBody.Injury,
			Species:  requestBody.Species,
			Id:       uuid.New(),
		}

		dataMu.Lock()
		data = append(data, newPost)
		dataMu.Unlock()

		c.JSON(http.StatusOK, newPost)
	})

	r.PUT("/shark/:id", func(c *gin.Context) {
		id := c.Param("id")
		var requestBody UpdatePostRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dataMu.Lock()
		defer dataMu.Unlock()
		for i := range data {
			if data[i].Id.String() == id {
				if requestBody.Date != "" {
					data[i].Date = requestBody.Date
				}
				if requestBody.Country != "" {
					data[i].Country = requestBody.Country
				}
				if requestBody.Location != "" {
					data[i].Location = requestBody.Location
				}
				if requestBody.Activity != "" {
					data[i].Activity = requestBody.Activity
				}
				if requestBody.Name != "" {
					data[i].Name = requestBody.Name
				}
				if requestBody.Age != "" {
					data[i].Age = requestBody.Age
				}
				if requestBody.Injury != "" {
					data[i].Injury = requestBody.Injury
				}
				if requestBody.Species != "" {
					data[i].Species = requestBody.Species
				}
				c.JSON(http.StatusOK, data[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
	})

	r.DELETE("/shark/:id", func(c *gin.Context) {
		id := c.Param("id")

		dataMu.Lock()
		defer dataMu.Unlock()
		for i := range data {
			if data[i].Id.String() == id {
				data = remove(data, i)
				c.JSON(http.StatusOK, gin.H{"status": "Post deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
	})

	return r
}

func remove(slice []Post, s int) []Post {
	return append(slice[:s], slice[s+1:]...)
}

func importData() []Post {
	content, err := os.ReadFile("./global-shark-attack.json")
	if err != nil {
		log.Fatal(err)
	}

	var data []Post
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatal("JSON Unmarshal error:", err)
	}

	for i := range data {
		data[i].Id = uuid.New()
	}
	return data
}

func getRandomItems(arr []Post, numItems int) []Post {
	if numItems > len(arr) {
		numItems = len(arr)
	}
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr[:numItems]
}
