package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/cs301-2023-g3t3/points-ledger/models"
	assert "github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var columns = []string{"id", "user_id", "balance"}
var gormDB, mock = SetUpDB()
var pointsService = NewPointsService(gormDB)

func SetUpDB() (*gorm.DB, sqlmock.Sqlmock){
    // Create a new GORM DB instance with a mocked SQL database
    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatalf("Error creating mock DB: %v", err)
    }


    mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.34"))
    // Create a GORM DB connection with the MySQL driver
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn:                      db,
		DriverName:                "mysql",
		SkipInitializeWithVersion: false,
    }), &gorm.Config{
            SkipDefaultTransaction: true,
        })
    if err != nil {
        log.Fatalf("Error creating GORM DB: %v", err)
    }

    gormDB.AutoMigrate(&models.PointsAccount{})

    // Insert multiple mock user data into the database
    for i := 1; i <= 2; i++ {
        gormDB.Create(&models.PointsAccount{
            Id: fmt.Sprintf("id%d", i),
            UserId: fmt.Sprintf("uid%d", i),
            Balance: i,
        })
    }
    
    return gormDB, mock
}

func TestGetAllPointsAccount(t *testing.T) {
    rows := sqlmock.NewRows(columns)

    var expectedPointsAccount []models.PointsAccount

    for i := 1; i <= 2; i++ {
        acc := models.PointsAccount{
            Id: fmt.Sprintf("id%d", i),
            UserId: fmt.Sprintf("uid%d", i),
            Balance: i,
        }

        expectedPointsAccount = append(expectedPointsAccount, acc)
    }

    for _, acc := range expectedPointsAccount {
        rows.AddRow(acc.Id, acc.UserId, acc.Balance)
    }

    statement := "SELECT * FROM `points_accounts`"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).WillReturnRows(rows)

    accounts, statusCode, err := pointsService.GetAllAccounts()

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, expectedPointsAccount, *accounts)
}

func TestGetPointsById(t *testing.T) {
    expectedRes := models.PointsAccount{
        Id: "id2",
        UserId: "uid2",
        Balance: 2,
    }

    row := sqlmock.NewRows(columns).AddRow(expectedRes.Id, expectedRes.UserId, expectedRes.Balance)

    statement := "SELECT * FROM `points_accounts` WHERE id = ?"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(expectedRes.Id).
        WillReturnRows(row)

    res, statusCode, err := pointsService.GetAccountById(expectedRes.Id)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, expectedRes, *res)
}

func TestGetPointsById_NotFound(t *testing.T) {
    invalidId := "id123"

    statement:= "SELECT * FROM `points_accounts` WHERE id = ?"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(invalidId).
        WillReturnError(gorm.ErrRecordNotFound)

    res, statusCode, err := pointsService.GetAccountById(invalidId)

    assert.Error(t, err, gorm.ErrRecordNotFound)
    assert.Equal(t, http.StatusNotFound, statusCode)
    assert.Nil(t, res)
}

func TestAdjustPoints_Add(t *testing.T) {
    id, userId, balance := "id1", "uid1", 1000
    
    row := sqlmock.NewRows(columns).AddRow(id, userId, 1)

    statement := "SELECT * FROM `points_accounts` WHERE id = ?"
    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(id).
        WillReturnRows(row)

    statement = "UPDATE `points_accounts` SET `id`=?,`user_id`=?,`balance`=? WHERE `id` = ?"
    mock.ExpectExec(regexp.QuoteMeta(statement)).
        WithArgs(id, userId, balance, id).
        WillReturnResult(sqlmock.NewResult(1,1))
    
    input := &models.Input{
        Action: "add",
        Amount: 999,
    }

    res, statusCode, err := pointsService.AdjustPoints(input, id)
    if err != nil {
        t.Fatal(err)
    }

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, balance, res.Balance)
}

func TestAdjustPoints_DeductFailed(t *testing.T) {
    id, userId := "id1", "uid1" 

    row := sqlmock.NewRows(columns).AddRow(id, userId, 1)

    statement := "SELECT * FROM `points_accounts` WHERE id = ?"
    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(id).
        WillReturnRows(row)

    input := &models.Input{
        Action: "deduct",
        Amount: 1000,
    }

    res, statusCode, err := pointsService.AdjustPoints(input, id)

    expectedErr := errors.New("insufficient points to deduct")
    assert.Error(t, expectedErr, err)
    assert.Equal(t, http.StatusBadRequest, statusCode)
    assert.Nil(t, res)
}
