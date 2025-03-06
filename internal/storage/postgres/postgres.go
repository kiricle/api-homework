package postgres

import (
	"database/sql"
	"fmt"
	"github.com/kiricle/api-homework/internal/models"
	_ "github.com/lib/pq"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStorage(log *slog.Logger) (*Storage, error) {
	const op = "storage.NewStorage.Postgres"
	log.Info(op)

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=%s port=%s", "postgres", "api_homework", "goLANGn1nja", "127.0.0.1", "disable", "5432"))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) GetBooks() ([]models.Book, error) {
	const op = "storage.getBooks"
	s.log.Info(op)

	rows, err := s.db.Query("SELECT * FROM books")
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	books := make([]models.Book, 0)
	for rows.Next() {
		var book models.Book
		err1 := rows.Scan(&book.Id, &book.Name, &book.Author)
		if err1 != nil {
			s.log.Error(err.Error())
			return nil, err1
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *Storage) GetBook(id int64) (models.Book, error) {
	const op = "storage.getBook"
	s.log.Info(op)

	row := s.db.QueryRow("SELECT * FROM books WHERE id = $1", id)

	var book models.Book
	if err := row.Scan(&book.Id, &book.Name, &book.Author); err != nil {
		s.log.Error(err.Error())
		return models.Book{}, err
	}

	return book, nil
}

func (s *Storage) CreateBook(name, author string) error {
	const op = "storage.createBook"
	s.log.Info(op)

	tx, err := s.db.Begin()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	_, execErr := tx.Exec("INSERT INTO books (name, author) VALUES ($1, $2)", name, author)
	if execErr != nil {
		s.log.Error(execErr.Error())
		if err := tx.Rollback(); err != nil {
			s.log.Error(err.Error())
			return err
		}
		return execErr
	}

	if err := tx.Commit(); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Info("storage.createBook success")
	return nil
}
