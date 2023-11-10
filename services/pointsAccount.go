package services

import (
	"errors"
	"net/http"

	"github.com/cs301-2023-g3t3/points-ledger/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

type PointsService struct {
    DB *gorm.DB
}

func NewPointsService(db *gorm.DB) *PointsService {
    return &PointsService{DB: db}
}

func (s *PointsService) GetAllAccounts() (*[]models.PointsAccount, int, error) {
	var accounts []models.PointsAccount

	if err := s.DB.Find(&accounts).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &accounts, http.StatusOK, nil
}

func (s *PointsService) GetAccountById(id string) (*models.PointsAccount, int, error) {
    if id == "" {
        return nil, http.StatusBadRequest, errors.New("id cannot be empty")
    }

    var account models.PointsAccount

    if err := s.DB.First(&account, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("account not found")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }
    
    return &account, http.StatusOK, nil
}

func (s *PointsService) GetAccountByUserId(userId string) (*[]models.PointsAccount, int, error) {
    if userId == "" {
        return nil, http.StatusBadRequest, errors.New("userId cannot be empty")
    }

    var accounts []models.PointsAccount

    if err := s.DB.Where("user_id = ?", userId).Find(&accounts).Error; err != nil {
        return nil, http.StatusInternalServerError, err
    }

    if len(accounts) == 0 {
        return nil, http.StatusNotFound, errors.New("no accounts found with userId")
    }

    return &accounts, http.StatusOK, nil
}

func (s *PointsService) AdjustPoints (input *models.Input, id string) (*models.PointsAccount, int, error) {
    if id == "" {
        return nil, http.StatusBadRequest, errors.New("id cannot be empty")
    }

    if err := validate.Struct(input); err != nil {
        return nil, http.StatusBadRequest, err
    }

    if input.Amount < 0 {
        return nil, http.StatusBadRequest, errors.New("negative value is not allowed")
    }

    var account models.PointsAccount
    if err := s.DB.Where("id = ?", id).First(&account).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, http.StatusNotFound, errors.New("account not found")
        } else {
            return nil, http.StatusInternalServerError, err
        }
    }

    switch input.Action {
    case "add":
        account.Balance += input.Amount
        break
    case "deduct":
        // if account.Balance >= input.Amount {
            account.Balance -= input.Amount
        // } else {
        //     return nil, http.StatusBadRequest, errors.New("insufficient points to deduct")
        // }
        break
    case "override":
        account.Balance = input.Amount
        break
    default:
        return nil, http.StatusBadRequest, errors.New("invalid action")
    }

    account.Id = id

    err := s.DB.Model(models.PointsAccount{Id: id}).Updates(&account).Error

    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

    return &account, http.StatusOK, nil
}
