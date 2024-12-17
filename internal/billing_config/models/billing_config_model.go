package models

type BillingConfig struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"` // The JSON value stored as a string
}

type BillingValueConfig struct {
	IsActive bool  `json:"is_active"`
	Value    int32 `json:"value"`
}
