package store

import (
	"testing"
	
	"github.com/mrckurz/CI-CD-MCM/internal/model"

)

var testProducts = []model.Product{
	{Name: "Laptop", Price: 1000},
	{Name: "Phone", Price: 500},
	{Name: "Tablet", Price: 300},
	{Name: "Headphones", Price: 100},
	{Name: "Mouse", Price: 50},
}


func TestCreateAndGet(t *testing.T) {
	_ = NewMemoryStore()
		tests := []struct {
		name    string
		product model.Product
	}{
		{"product 1", testProducts[0]},
		{"product 2", testProducts[1]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()

			created := s.Create(tt.product)

			got, err := s.GetByID(created.ID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Name != tt.product.Name || got.Price != tt.product.Price {
				t.Errorf("got %+v, want %+v", got, tt.product)
			}
		})
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

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name    string
		initial model.Product
		update  model.Product
	}{
		{
			"update laptop",
			testProducts[0],
			testProducts[2],
		},
		{
			"update phone",
			testProducts[1],
			testProducts[3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()

			created := s.Create(tt.initial)

			updated := created
			updated.Name = tt.update.Name
			updated.Price = tt.update.Price

			_, err := s.Update(created.ID, updated)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := s.GetByID(created.ID)

			if got.Name != tt.update.Name || got.Price != tt.update.Price {
				t.Errorf("update failed, got %+v, want %+v", got, tt.update)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name    string
		product model.Product
	}{
		{"delete laptop", testProducts[0]},
		{"delete phone", testProducts[1]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()

			created := s.Create(tt.product)

			err := s.Delete(created.ID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, err = s.GetByID(created.ID)
			if err != ErrNotFound {
				t.Errorf("expected ErrNotFound, got %v", err)
			}
		})
	}
}

func TestGetByIDNotFound(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{"non-existent high ID", 999},
		{"negative ID", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()

			_, err := s.GetByID(tt.id)
			if err != ErrNotFound {
				t.Errorf("expected ErrNotFound, got %v", err)
			}
		})
	}
}