package controllers

import (
	"net/http"
	// "strconv"
	
	"github.com/joegb/email-forwarder/internal/database"
	"github.com/joegb/email-forwarder/internal/models"
	
	"github.com/gin-gonic/gin"
)

func CreateTarget(c *gin.Context) {
	var target models.ForwardTarget
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := database.DB.Create(&target).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Name already exists"})
		return
	}
	
	c.JSON(http.StatusCreated, target)
}

func ListTargets(c *gin.Context) {
	var targets []models.ForwardTarget
	database.DB.Find(&targets)
	c.JSON(http.StatusOK, targets)
}

func GetTarget(c *gin.Context) {
	id := c.Param("id")
	
	var target models.ForwardTarget
	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	
	c.JSON(http.StatusOK, target)
}

func UpdateTarget(c *gin.Context) {
	id := c.Param("id")
	
	var target models.ForwardTarget
	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	database.DB.Save(&target)
	c.JSON(http.StatusOK, target)
}

func DeleteTarget(c *gin.Context) {
	id := c.Param("id")
	
	var target models.ForwardTarget
	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	
	database.DB.Delete(&target)
	c.Status(http.StatusNoContent)
}