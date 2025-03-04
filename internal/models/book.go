package models

type Book struct {
	Id     int64  `json:"id"`
	Name   string `json:"name" binding:"required"`
	Author string `json:"author" binding:"required"`
}
