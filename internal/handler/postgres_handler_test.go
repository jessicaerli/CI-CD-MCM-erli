package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func newMockPostgresHandler(t *testing.T) (*PostgresHandler, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	s := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(s)

	return h, mock, func() {
		db.Close()
	}
}

func TestNewPostgresHandler(t *testing.T) {
	s := &store.PostgresStore{}

	h := NewPostgresHandler(s)

	if h == nil {
		t.Fatal("expected handler, got nil")
	}

	if h.Store != s {
		t.Fatal("handler does not contain the expected store")
	}
}

func TestPostgresHandlerRegisterRoutes(t *testing.T) {
	h, _, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"GET", "/products"},
		{"POST", "/products"},
		{"GET", "/products/1"},
		{"PUT", "/products/1"},
		{"DELETE", "/products/1"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		match := &mux.RouteMatch{}

		if !r.Match(req, match) {
			t.Errorf("route not registered: %s %s", tt.method, tt.path)
		}
	}
}

func TestPostgresHandlerHealth(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectPing()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerHealthError(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectPing().WillReturnError(errors.New("database unavailable"))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerGetProducts(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Tofu", 2.99).
		AddRow(2, "Sushi Rice 10l", 19.99)

	mock.ExpectQuery(`SELECT id, name, price FROM products ORDER BY id`).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()

	h.GetProducts(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerGetProductsError(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT id, name, price FROM products ORDER BY id`).
		WillReturnError(errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()

	h.GetProducts(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerGetProduct(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	row := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Tofu", 2.99)

	mock.ExpectQuery(`SELECT id, name, price FROM products WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(row)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	h.GetProduct(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerGetProductNotFound(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT id, name, price FROM products WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	h.GetProduct(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerCreateProduct(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectQuery(`INSERT INTO products \(name, price\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("Tofu", 2.99).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	body := `{"name":"Tofu","price":2.99}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateProduct(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerCreateProductInvalidJSON(t *testing.T) {
	h, _, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(`{"name":`),
	)
	rec := httptest.NewRecorder()

	h.CreateProduct(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostgresHandlerCreateProductInvalid(t *testing.T) {
	h, _, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	body := `{"name":"","price":-1}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateProduct(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostgresHandlerCreateProductStoreError(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectQuery(`INSERT INTO products \(name, price\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("Tofu", 2.99).
		WillReturnError(errors.New("insert failed"))

	body := `{"name":"Tofu","price":2.99}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateProduct(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerUpdateProduct(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE products SET name = \$1, price = \$2 WHERE id = \$3`).
		WithArgs("Smoky Tofu", 3.99, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"name":"Smoky Tofu","price":3.99}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	h.UpdateProduct(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerUpdateProductInvalidJSON(t *testing.T) {
	h, _, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		strings.NewReader(`{"name":`),
	)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	h.UpdateProduct(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostgresHandlerUpdateProductNotFound(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectExec(`UPDATE products SET name = \$1, price = \$2 WHERE id = \$3`).
		WithArgs("Smoky Tofu", 3.99, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body := `{"name":"Smoky Tofu","price":3.99}`
	req := httptest.NewRequest(http.MethodPut, "/products/999", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	h.UpdateProduct(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerDeleteProduct(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	h.DeleteProduct(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresHandlerDeleteProductNotFound(t *testing.T) {
	h, mock, cleanup := newMockPostgresHandler(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest(http.MethodDelete, "/products/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	h.DeleteProduct(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
