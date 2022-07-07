package service

import (
	"fmt"
	"log"

	"github.com/Leakageonthelamp/golang-api/dto"
	"github.com/Leakageonthelamp/golang-api/entity"
	"github.com/Leakageonthelamp/golang-api/repository"
	"github.com/mashingan/smapping"
)

type BookService interface {
	InsertBook(book dto.BookCreateDTO) entity.Book
	UpdateBook(book dto.BookUpdateDTO) entity.Book
	DeleteBook(book entity.Book)
	AllBooks() []entity.Book
	FindBookByID(id uint64) entity.Book
	IsAllowToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookService repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookService: bookRepo,
	}
}

func (service *bookService) InsertBook(book dto.BookCreateDTO) entity.Book {
	bookToInsert := entity.Book{}
	err := smapping.FillStruct(&bookToInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Error while mapping book create dto to book entity: %v", err)
	}

	return service.bookService.InsertBook(bookToInsert)
}

func (service *bookService) UpdateBook(book dto.BookUpdateDTO) entity.Book {
	bookToUpdate := entity.Book{}
	err := smapping.FillStruct(&bookToUpdate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Error while mapping book update dto to book entity: %v", err)
	}

	return service.bookService.UpdateBook(bookToUpdate)
}

func (service *bookService) DeleteBook(book entity.Book) {
	service.bookService.DeleteBook(book)
}

func (service *bookService) AllBooks() []entity.Book {
	return service.bookService.AllBooks()
}

func (service *bookService) FindBookByID(id uint64) entity.Book {
	return service.bookService.FindBookByID(id)
}

func (service *bookService) IsAllowToEdit(userID string, bookID uint64) bool {
	book := service.bookService.FindBookByID(bookID)
	v := fmt.Sprintf("%v", book.UserID)
	return v == userID
}
