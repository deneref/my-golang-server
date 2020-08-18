package models

import (
	"time"
)

// SQLDataBase struct
type SQLDataBase struct {
	Server          string   `toml:"Server"`
	Database        string   `toml:"Database"`
	ApplicationName string   `toml:"ApplicationName"`
	MaxIdleConns    int      `toml:"MaxIdleConns"`
	MaxOpenConns    int      `toml:"MaxOpenConns"`
	ConnMaxLifetime duration `toml:"ConnMaxLifetime"`
	UserID          string
	Password        string
}

type Order struct {
	Rv      int       `json:"rv"`
	Content []Content `json:"content"`
}

type Content struct {
	Goods       []Goods   `json:"goods"`
	Status      string    `json:"status"`
	OrderID     int       `json:"order_id"`
	StoreID     int       `json:"store_id"`
	DateCreated time.Time `json:"date_created"`
}

type Goods struct {
	Gid    int    `json:"gid"`
	Price  int    `json:"price"`
	Status string `json:"status"`
	ChrtID int    `json:"chrt_id"`
}
