package models

type UserModel struct {
	ID           int32  `db:"id"`
	Name         string `db:"name"`
	IsDelinquent bool   `db:"is_delinquent"`
}
