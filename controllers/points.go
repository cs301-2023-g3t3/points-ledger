package controllers

import (
	"github.com/cs301-2023-g3t3/points-ledger/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PointsController struct{
	db *gorm.DB
}

func (s *PointsController) SetDB(db *gorm.DB) {
    s.db = db
}

func (s PointsController) GetAccounts(c *gin.Context) {
	userID := c.Query("userID")

	var accounts []models.PointsAccount
	query := s.db

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Error fetching points accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)

}

func (s PointsController) GetSpecificAccount(c *gin.Context) {

	accountID := c.Param("ID")
	var account models.PointsAccount

	if err := s.db.First(&account, "ID = ?", accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s PointsController) AdjustPoints(c *gin.Context) {
	
	accountID := c.Param("ID")
	var input models.Input

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Invalid payload request"})
		return
	}

	c.Set("input", input)

	if input.Amount < 0 {
		message := "Negative values are not allowed"
		c.Set("message", message)
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: message})
			return
	}
	

	var account models.PointsAccount
	if err := s.db.Where("ID = ?", accountID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, models.HTTPError{Code: http.StatusNotFound, Message: "Account not found"})
		return
	}

	switch input.Action {
	case "add":
		account.Balance += input.Amount
	case "deduct":
		if account.Balance >= input.Amount {
			account.Balance -= input.Amount
		} else {
			message := "Insufficient points to deduct"
			c.Set("message", message)
			c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: message})
			return
		}
	case "override":
		account.Balance = input.Amount
	default:
		message := "Invalid action"
		c.Set("message", message)
		c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: message})
		return
	}

	s.db.Save(&account)
	
	message := "Points adjusted successfully"
	c.Set("message", message)
	c.JSON(http.StatusOK, gin.H{
		"message":     message,
		"new_balance": account.Balance,
	})
}
