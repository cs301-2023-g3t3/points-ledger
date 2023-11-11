package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

//  @Summary        Get all Point Accounts
//  @Description    Retrieves a list of point accounts
//  @Tags           points
//  @Produce        json
//  @Success        200     {array}     models.PointsAccount
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts   [get]
func (s PointsController) GetAllAccounts(c *gin.Context) {
    accounts, code, err := s.PointsService.GetAllAccounts()
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code: code,
            Message: err.Error(),
        })
        return
    }

	c.JSON(code, *accounts)
}

//  @Summary        Get all Point Accounts by Pagination
//  @Description    Retrieves a list of point accounts
//  @Tags           points
//  @Produce        json
//  @Param          page    query   int     true    "page"
//  @Param          size    query   int     true    "size"
//  @Success        200     {array}     models.PointsAccount
//  @Failure        400     {object}    models.HTTPError    "Invalid parameters"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/paginate   [get]
func (s PointsController) GetPaginatedAccounts(c *gin.Context) {
    page := c.DefaultQuery("page", "1")  
    pageSize := c.DefaultQuery("size", "10")  

    // Convert page and pageSize to integers
    pageInt, err := strconv.Atoi(page)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid page parameter",
        })
        return
    }

    pageSizeInt, err := strconv.Atoi(pageSize)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid pageSize parameter",
        })
        return
    }

    accounts, code, err := s.PointsService.GetPaginatedAccounts(pageInt, pageSizeInt)
    if err != nil {
        c.JSON(code, models.HTTPError{
            Code:    code,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":       accounts,
        "pagination": gin.H{"page": pageInt, "size": pageSizeInt},
    })
}

//  @Summary        Get Points Account by Id
//  @Description    Retrieve an Account By ID
//  @Tags           points
//  @Produce        json
//  @Param          ID      path    string  true    "ID"
//  @Success        200     {object}    models.PointsAccount
//  @Failure        400     {object}    models.HTTPError    "Id cannot be empy"
//  @Failure        404     {object}    models.HTTPError    "Points Account not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [get]
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

//  @Summary        Get Points Account by UserId
//  @Description    Retrieve a list of Points Account By UserID
//  @Tags           points
//  @Produce        json
//  @Param          UserID      path    string  true    "UserID"
//  @Success        200     {array}     models.PointsAccount
//  @Failure        400     {object}    models.HTTPError    "Id cannot be empty"
//  @Failure        404     {object}    models.HTTPError    "Points Account not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [get]
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

//  @Summary        Update Points by Id
//  @Description    Update Points By Id
//  @Tags           points
//  @Produce        json
//  @Param          ID      path    string  true    "ID"
//  @Success        200     {string}    "Points adjusted successfully"
//  @Failure        400     {object}    models.HTTPError    "Bad request due to invalid JSON body"
//  @Failure        404     {object}    models.HTTPError    "User not found with Id"
//  @Failure        500     {object}    models.HTTPError
//  @Router         /accounts/{id}   [put]
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
