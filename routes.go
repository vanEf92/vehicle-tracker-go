package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repo *VehicleRepository) {

	r.GET("/vehicles/:vehicle_id/location", func(c *gin.Context) {
		vehicleID := c.Param("vehicle_id")
		p, err := repo.GetLatest(vehicleID)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicles not found!"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	r.GET("/vehicles/:vehicle_id/history", func(c *gin.Context) {
		vehicleID := c.Param("vehicle_id")
		start, err1 := strconv.ParseInt(c.Query("start"), 10, 64)
		end, err2 := strconv.ParseInt(c.Query("end"), 10, 64)

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Start or End parameter"})
			return
		}

		history, err := repo.GetHistory(vehicleID, start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, history)
	})
}