package service

import (
	"MCatalogue/config"
	"MCatalogue/internal/model"
	"MCatalogue/internal/repository/postgresql"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type Service struct {
	log  *slog.Logger
	repo *postgresql.Repository
}

func New(log *slog.Logger, repo *postgresql.Repository) (*Service, error) {
	return &Service{
		log:  log,
		repo: repo,
	}, nil
}

func (s *Service) GetCar(c *gin.Context) {
	const op = "getting car data"
	slog.Info(op)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	param := c.Query("param")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	columnName := c.Param("columnName")

	offset := (page - 1) * limit

	cars, err := s.repo.GetAllCarData(columnName, param, limit, offset)
	if err != nil {
		slog.Error("Error getting car data: %v", err)
		c.JSON(500, gin.H{
			"error": "failed to get car data" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"cars": cars,
	})

}

func (s *Service) DeleteCar(c *gin.Context) {
	const op = "deleting car row"
	slog.Info(op)

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	err = s.repo.DeleteCarRow(carId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to delete row " + err.Error(),
		})
		return
	}
	c.Status(200)

}
func (s *Service) UpdateCar(c *gin.Context) {
	const op = "updating car row"
	slog.Info(op)

	carId, err := strconv.Atoi(c.Param("car_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	var car model.Car
	err = c.ShouldBindJSON(&car)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	err = s.repo.EditCarRow(carId, model.Car{
		Id:      carId,
		Model:   car.Model,
		Mark:    car.Mark,
		RegNum:  car.RegNum,
		OwnerId: car.OwnerId,
		Year:    car.Year,
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to edit row : " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "updated",
	})

}
func (s *Service) AddCar(c *gin.Context) {
	const op = "adding car"
	slog.Info(op)

	type RequestBody struct {
		RegNums []string `json:"regNums"`
	}
	type Car struct {
		RegNum string       `json:"regNum"`
		Mark   string       `json:"mark "`
		Model  string       `json:"model"`
		Year   int          `json:"year,omitempty"`
		Owner  model.People `json:"owner"`
	}

	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error ": err.Error()})
		return
	}
	reqBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal request body"})
		return
	}

	conf, _ := config.LoadConfig()
	externalApi := conf.ExternalApi

	response, err := http.Post(externalApi, "application/json", bytes.NewReader(reqBodyJson))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to make external API request"})
		return
	}
	defer response.Body.Close()

	var car Car
	if err := json.NewDecoder(response.Body).Decode(&car); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse response from external API"})
		return
	}
	c.JSON(response.StatusCode, car)

}
