package main

import (
	"net/http"
	"gorm.io/gorm"
)

type ToDo struct {
	ID   int    `gorm:"primaryKey"`
	Task string `gorm:"not null;default:null" json:"task"`
	Done bool   `json:"done"`
}

type database struct {
	err error
	DB  *gorm.DB
}

type Response struct {
	Response string `json:"response"`
}

type App struct { 
	server *http.Server
	db *database
}

