package store

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockPostgresStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return &PostgresStore{DB: db}, mock
}

func TestEnsureTable(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`
		CREATE TABLE IF NOT EXISTS products (
			id    SERIAL PRIMARY KEY,
			name  TEXT NOT NULL,
			price NUMERIC(10,2) NOT NULL DEFAULT 0
		)
	`)).WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.EnsureTable()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestEnsureTableError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	expectedErr := errors.New("database error")

	mock.ExpectExec(regexp.QuoteMeta(`
		CREATE TABLE IF NOT EXISTS products (
			id    SERIAL PRIMARY KEY,
			name  TEXT NOT NULL,
			price NUMERIC(10,2) NOT NULL DEFAULT 0
		)
	`)).WillReturnError(expectedErr)

	err := s.EnsureTable()

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected database error, got %v", err)
	}
}

func TestPostgresGetAll(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Tofu", 2.99).
		AddRow(2, "Sushi Rice 10l", 19.99)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products ORDER BY id",
	)).WillReturnRows(rows)

	products, err := s.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	if products[0].Name != "Tofu" {
		t.Errorf("expected Tofu, got %s", products[0].Name)
	}

	if products[1].Price != 19.99 {
		t.Errorf("expected price 19.99, got %f", products[1].Price)
	}
}

func TestPostgresGetAllEmpty(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	rows := sqlmock.NewRows([]string{"id", "name", "price"})

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products ORDER BY id",
	)).WillReturnRows(rows)

	products, err := s.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if products == nil {
		t.Fatal("expected empty slice, got nil")
	}

	if len(products) != 0 {
		t.Fatalf("expected 0 products, got %d", len(products))
	}
}

func TestPostgresGetAllQueryError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	expectedErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products ORDER BY id",
	)).WillReturnError(expectedErr)

	products, err := s.GetAll()

	if products != nil {
		t.Errorf("expected nil products, got %v", products)
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Tofu", 2.99)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products WHERE id = $1",
	)).
		WithArgs(1).
		WillReturnRows(rows)

	product, err := s.GetByID(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected ID 1, got %d", product.ID)
	}

	if product.Name != "Tofu" {
		t.Errorf("expected Widget, got %s", product.Name)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products WHERE id = $1",
	)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	product, err := s.GetByID(999)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if product.ID != 0 {
		t.Errorf("expected zero-value product, got %+v", product)
	}
}

func TestPostgresGetByIDError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	expectedErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, price FROM products WHERE id = $1",
	)).
		WithArgs(1).
		WillReturnError(expectedErr)

	_, err := s.GetByID(1)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	product := model.Product{
		Name:  "Tofu",
		Price: 2.99,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
	)).
		WithArgs("Tofu", 2.99).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	created, err := s.Create(product)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}

	if created.Name != "Tofu" {
		t.Errorf("expected Tofu, got %s", created.Name)
	}
}

func TestPostgresCreateError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	product := model.Product{
		Name:  "Tofu",
		Price: 2.99,
	}

	expectedErr := errors.New("insert failed")

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
	)).
		WithArgs("Tofu", 2.99).
		WillReturnError(expectedErr)

	_, err := s.Create(product)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected insert error, got %v", err)
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	product := model.Product{
		Name:  "Smoky Tofu",
		Price: 3.99,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
	)).
		WithArgs("Smoky Tofu", 3.99, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updated, err := s.Update(1, product)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.ID != 1 {
		t.Errorf("expected ID 1, got %d", updated.ID)
	}

	if updated.Name != "Smoky Tofu" {
		t.Errorf("expected Smoky Tofu, got %s", updated.Name)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	product := model.Product{
		Name:  "Missing",
		Price: 10,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
	)).
		WithArgs("Missing", 10.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.Update(999, product)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresUpdateError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	product := model.Product{
		Name:  "Tofu",
		Price: 2.99,
	}

	expectedErr := errors.New("update failed")

	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
	)).
		WithArgs("Tofu", 2.99, 1).
		WillReturnError(expectedErr)

	_, err := s.Update(1, product)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected update error, got %v", err)
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM products WHERE id = $1",
	)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.Delete(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM products WHERE id = $1",
	)).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := s.Delete(999)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDeleteError(t *testing.T) {
	s, mock := newMockPostgresStore(t)

	expectedErr := errors.New("delete failed")

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM products WHERE id = $1",
	)).
		WithArgs(1).
		WillReturnError(expectedErr)

	err := s.Delete(1)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected delete error, got %v", err)
	}
}
