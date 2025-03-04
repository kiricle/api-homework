package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kiricle/api-homework/internal/models"
	"github.com/kiricle/api-homework/internal/services"
	"log/slog"
	"net/http"
	"strconv"
)

type BookHandler struct {
	log         *slog.Logger
	bookService *services.BookService
}

func NewBookHandler(log *slog.Logger, bookService *services.BookService) *BookHandler {
	return &BookHandler{log: log, bookService: bookService}
}

func (bh *BookHandler) GetBooks(c *gin.Context) {
	bh.log.Info("bookHandler.GetBooks.Start")

	books, err := bh.bookService.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, books)
	bh.log.Info("bookHandler.GetBooks.End")
}

func (bh *BookHandler) CreateBook(c *gin.Context) {
	bh.log.Info("bookHandler.CreateBook.Start")

	var book models.Book
	if err := c.ShouldBind(&book); err != nil {
		bh.log.Info("bookHandler.CreateBook.Bind", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bh.bookService.CreateBook(book.Name, book.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book created"})
}

func (bh *BookHandler) GetBook(c *gin.Context) {
	bh.log.Info("bookHandler.GetBook.Start")
	bookIDString := c.Param("id")

	bookID, err := strconv.ParseInt(bookIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := bh.bookService.GetBook(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}
