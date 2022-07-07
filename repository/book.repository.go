package repository

import (
	"github.com/Leakageonthelamp/golang-api/entity"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	DeleteBook(book entity.Book)
	AllBooks() []entity.Book
	FindBookByID(id uint64) entity.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (db *bookConnection) InsertBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) UpdateBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) DeleteBook(book entity.Book) {
	db.connection.Delete(&book)
}

func (db *bookConnection) AllBooks() []entity.Book {
	var books []entity.Book
	db.connection.Preload("User").Find(&books)
	return books
}

func (db *bookConnection) FindBookByID(id uint64) entity.Book {
	var book entity.Book
	db.connection.Preload("User").Find(&book, id)
	return book
}
