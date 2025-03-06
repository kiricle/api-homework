package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kiricle/api-homework/internal/models"
	// gin-swagger middleware
	"log/slog"
	"net/http"
	"strconv"
)

type BookHandler struct {
	log         *slog.Logger
	bookService BookService
}

type BookService interface {
	CreateBook(name, author string) error
	GetBook(id int64) (models.Book, error)
	GetBooks() ([]models.Book, error)
}

func NewBookHandler(log *slog.Logger, bookService BookService) *BookHandler {
	return &BookHandler{log: log, bookService: bookService}
}

// @Summary GetBooks
// @Tags books
// @Get
// @Description Get a list of books
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Router /book [get]
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

// @Summary CreateBook
// @tags books
// @Post
// @Description Create book
// @Accept json
// @Produce json
// @Param book body models.CreateBookInput true "Book data"
// @Success 200
// @Router /book [post]
func (bh *BookHandler) CreateBook(c *gin.Context) {
	bh.log.Info("bookHandler.CreateBook.Start")

	var book models.CreateBookInput
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

// @Summary GetBook
// @Tags books
// @Get
// @Description Get a book by id
// @Accept json
// @Param id path integer true "Book ID"
// @Produce json
// @Success 200 {integer} integer
// @Router /book/{id} [get]
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
