package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	product := model.Product{
		Name:  "Sushi Rice",
		Price: 5.99,
	}

	created := s.Create(product)

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "Sushi Rice" {
		t.Errorf("expected Name Sushi Rice, got %s", got.Name)
	}

	if got.Price != 5.99 {
		t.Errorf("expected Price 5.99, got %f", got.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestUpdate(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(
		model.Product{
			Name:  "Ube Jam",
			Price: 6.75,
		})

	updated, err := s.Update(created.ID, model.Product{
		Name:  "Ube Jam v2",
		Price: 7.50,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}

	if updated.Name != "Ube Jam v2" {
		t.Errorf("expected name Ube Jam v2, got name %s", updated.Name)
	}

	if updated.Price != 7.50 {
		t.Errorf("expected price 7.50, got %f", updated.Price)
	}
}

func TestDelete(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(
		model.Product{
			Name:  "Ube Jam",
			Price: 6.75,
		})

	err := s.Delete(created.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.GetByID(created.ID)

	if err != ErrNotFound {
		t.Errorf("expected product to be deleted")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetByID(999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
