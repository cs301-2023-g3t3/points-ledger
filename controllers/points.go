package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cs301-2023-g3t3/points-ledger/models"
	"github.com/cs301-2023-g3t3/points-ledger/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PointsController struct{
    PointsService *services.PointsService
}

func NewPointsController(db gorm.DB) *PointsController {
    return &PointsController{
        PointsService: services.NewPointsService(&db),
    }
}

func (s PointsController) GetAccounts(c *gin.Context) {
	userID := c.Query("userID")

    accounts, code, err := s.PointsService.GetAccounts(userID)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: err.Error(),
        })
        return
    }

	c.JSON(code, *accounts)
}

func (s PointsController) GetSpecificAccount(c *gin.Context) {
	accountID := c.Param("ID")

    account, code, err := s.PointsService.GetAccountById(accountID)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: err.Error(),
        })
        return
    }

	c.JSON(code, *account)
}

func (s PointsController) GetAccountByUser(c *gin.Context) {
	userID := c.Param("UserID")

    accounts, code, err := s.PointsService.GetAccountByUserId(userID)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: err.Error(),
        })
    }

	c.JSON(code, *accounts)
}

func (s PointsController) AdjustPoints(c *gin.Context) {
	accountID := c.Param("ID")
	var input models.Input
    if err := json.NewDecoder(c.Request.Body).Decode(&input); err != nil {
        c.Set("message", err.Error())
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code: http.StatusBadRequest,
            Message: fmt.Sprintf("Invalid JSON request: %v", err.Error()),
        })
        return
    }

    account, code, err := s.PointsService.AdjustPoints(&input, accountID)
    if err != nil {
        c.Set("message", err.Error())
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: err.Error(),
        })
        return
    }

	c.Set("input", input)

	message := "Points adjusted successfully"
	c.Set("message", message)
	c.JSON(http.StatusOK, gin.H{
		"message":     message,
		"new_balance": account.Balance,
	})
}
