package services

import (
	"errors"
	"fmt"
	"github.com/kiricle/api-homework/internal/models"
	"github.com/kiricle/api-homework/internal/storage/cache"
	"log/slog"
	"time"
)

type BookService struct {
	storage BookRepository
	cache   *cache.Cache
	log     *slog.Logger
}

type BookRepository interface {
	CreateBook(name, author string) error
	GetBook(id int64) (models.Book, error)
	GetBooks() ([]models.Book, error)
}

func NewBookService(storage BookRepository, cache *cache.Cache, log *slog.Logger) *BookService {
	return &BookService{storage: storage, cache: cache, log: log}
}

func (s *BookService) CreateBook(name string, author string) error {
	if err := s.storage.CreateBook(name, author); err != nil {
		return err
	}

	return nil
}

func (s *BookService) GetBooks() ([]models.Book, error) {
	key := "books"

	data := s.cache.Storage.Get(key)
	if data != nil {
		s.log.Info("Cache found for /books")
		return data.([]models.Book), nil
	}
	s.log.Info("No cache found for /books")

	books, err := s.storage.GetBooks()
	if err != nil {
		return nil, err
	}
	s.cache.Storage.Set(key, books, time.Second*30)

	return books, nil
}

func (s *BookService) GetBook(id int64) (models.Book, error) {
	op := fmt.Sprintf("BookService.GetBook %d", id)
	key := fmt.Sprintf("book/%d", id)
	s.log.Info(op)

	data := s.cache.Storage.Get(key)
	if data != nil {
		s.log.Info(fmt.Sprintf("Cache found for /books/%d", id))

		book, ok := data.(models.Book)
		if !ok {
			return models.Book{}, errors.New("could not assert book")
		}
		return book, nil
	}
	s.log.Info("No cache found for /books/%d", id)
	book, err := s.storage.GetBook(id)
	if err != nil {
		return models.Book{}, err
	}
	s.cache.Storage.Set(key, book, time.Second*30)
	return book, nil
}
