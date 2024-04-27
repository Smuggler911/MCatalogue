package api

import (
	"MCatalogue/internal/repository/postgresql"
	"MCatalogue/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger *slog.Logger
}

func (h *Handler) InitRoutes(repo *postgresql.Repository) *gin.Engine {
	router := gin.New()
	s, err := service.New(h.logger, repo)
	if err != nil {
		slog.Error("Failed to initialize service instance: ", err)
	}
	car := router.Group("/car")
	{
		car.GET("/:columnName", s.GetCar)
		car.DELETE("/delete/:car_id", s.DeleteCar)
		car.PUT("/update/:car_id", s.UpdateCar)
		car.POST("/add", s.AddCar)

	}
	return router
}
