package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Leakageonthelamp/golang-api/dto"
	"github.com/Leakageonthelamp/golang-api/entity"
	"github.com/Leakageonthelamp/golang-api/helper"
	"github.com/Leakageonthelamp/golang-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type BookController interface {
	All(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) All(ctx *gin.Context) {
	var books []entity.Book = c.bookService.AllBooks()
	res := helper.BuildResponse(true, "Get all books", books)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = c.bookService.FindBookByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "No data", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Get book by id", book)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Invalid Request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)
	convertUserID, err := strconv.ParseUint(userID, 10, 64)

	if err == nil {
		bookDTO.UserID = convertUserID
	}

	result := c.bookService.InsertBook(bookDTO)
	res := helper.BuildResponse(true, "Book Inserted", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) Update(ctx *gin.Context) {
	// Check if send with param id
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Check if valid request
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Invalid Request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Check if book exist
	var book entity.Book = c.bookService.FindBookByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "No data", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Check if valid token
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	findUserId := fmt.Sprintf("%v", claims["user_id"])

	// Check if user is owner of book
	IsAllow := c.bookService.IsAllowToEdit(findUserId, id)
	if !IsAllow {
		response := helper.BuildErrorResponse("You are not allowed to edit this book", "You are not allowed to edit this book", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	// Update book
	userId, errID := strconv.ParseUint(findUserId, 10, 64)
	if errID == nil {
		bookUpdateDTO.ID = id
		bookUpdateDTO.UserID = userId
	}
	result := c.bookService.UpdateBook(bookUpdateDTO)
	response := helper.BuildResponse(true, "Book Updated", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *bookController) Delete(ctx *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	}
	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	IsAllow := c.bookService.IsAllowToEdit(userID, book.ID)
	if !IsAllow {
		response := helper.BuildErrorResponse("You are not allowed to edit this book", "You are not allowed to edit this book", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	c.bookService.DeleteBook(book)
	res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])

	return id
}
