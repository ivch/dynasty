package models

type Request struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	UserID      uint   `json:"user_id"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (Request) TableName() string { return "requests" }
