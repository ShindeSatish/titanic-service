// internal/app/handler/handler.go
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shindesatish/titanic-service/internal/app/dto"
	"github.com/shindesatish/titanic-service/internal/app/service"
)

type PassengerHandler struct {
	PassengerService *service.PassengerService
}

func NewPassengerHandler(passengerService *service.PassengerService) *PassengerHandler {
	return &PassengerHandler{PassengerService: passengerService}
}

// @Summary Get all passengers
// @Description Get a list of all passengers in JSON format
// @Tags passengers
// @Produce json
// @Success 200 {array} model.Passenger "OK"
// @Router /passengers [get]
func (h *PassengerHandler) GetAllPassengersHandler(c *gin.Context) {
	passengers, err := h.PassengerService.GetAllPassengers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, passengers)
}

// @Summary Get passenger by ID
// @Description Get passenger data by PassengerId in JSON format
// @Tags passengers
// @Produce json
// @Param id path int true "Passenger ID"
// @Success 200 {object} model.Passenger "OK"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /passengers/{id} [get]
func (h *PassengerHandler) GetPassengerByIDHandler(c *gin.Context) {
	passengerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passenger, err := h.PassengerService.GetPassengerByID(context.Background(), uint(passengerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, passenger)
}

// @Summary Get selected attributes of passenger by ID
// @Description Get selected attributes of passenger by PassengerId in JSON format
// @Tags passengers
// @Produce json
// @Param id path int true "Passenger ID"
// @Param attributes query array true "List of attributes to retrieve"
// @Success 200 {object} map[string]interface{} "OK"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /passenger-attributes/{id} [get]
func (h *PassengerHandler) GetPassengerAttributesHandler(c *gin.Context) {
	passengerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passenger ID"})
		return
	}

	attributes := c.QueryArray("attributes")
	if len(attributes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No attributes specified"})
		return
	}

	for _, attr := range attributes {
		if !dto.ValidAttribute(attr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attribute specified - " + attr, "allowed_attribute": dto.AllowedAttributes})
			return
		}
	}

	passenger, err := h.PassengerService.GetPassengerAttributes(context.Background(), uint(passengerID), attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, passenger)
}

// @Summary Get fare histogram
// @Description Get a histogram of fare prices in percentiles
// @Tags passengers
// @Produce json
// @Success 200 {object} map[string]int "OK"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fare-histogram [get]
func (h *PassengerHandler) GetFareHistogramHandler(c *gin.Context) {
	// Fetch fare data from the service
	fareData, err := h.PassengerService.GetFareHistogram(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get fare histogram: %v", err)})
		return
	}
	// Convert result map to JSON and send the response
	c.JSON(http.StatusOK, fareData)
}
