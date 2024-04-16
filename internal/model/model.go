package model

import "time"

type User struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
