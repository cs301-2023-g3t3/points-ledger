package models

type PointsAccount struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	Balance int    `json:"balance"`
}
