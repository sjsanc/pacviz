package repository

import (
	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// MockRepository is a mock implementation for testing.
type MockRepository struct {
	Packages []*domain.Package
}

// NewMockRepository creates a new mock repository.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Packages: make([]*domain.Package, 0),
	}
}

// GetInstalled returns mock installed packages.
func (m *MockRepository) GetInstalled() ([]*domain.Package, error) {
	return m.Packages, nil
}

// GetExplicit returns mock explicit packages.
func (m *MockRepository) GetExplicit() ([]*domain.Package, error) {
	result := make([]*domain.Package, 0)
	for _, pkg := range m.Packages {
		if pkg.InstallReason == domain.ReasonExplicit {
			result = append(result, pkg)
		}
	}
	return result, nil
}

// GetOrphans returns mock orphaned packages.
func (m *MockRepository) GetOrphans() ([]*domain.Package, error) {
	result := make([]*domain.Package, 0)
	for _, pkg := range m.Packages {
		if pkg.IsOrphan {
			result = append(result, pkg)
		}
	}
	return result, nil
}

// GetForeign returns mock foreign packages.
func (m *MockRepository) GetForeign() ([]*domain.Package, error) {
	result := make([]*domain.Package, 0)
	for _, pkg := range m.Packages {
		if pkg.IsForeign {
			result = append(result, pkg)
		}
	}
	return result, nil
}

// Search returns mock search results.
func (m *MockRepository) Search(query string) ([]*domain.Package, error) {
	// TODO: Implement simple substring search
	return m.Packages, nil
}

// GetPackage returns a mock package by name.
func (m *MockRepository) GetPackage(name string) (*domain.Package, error) {
	for _, pkg := range m.Packages {
		if pkg.Name == name {
			return pkg, nil
		}
	}
	return nil, nil
}

// Install is a no-op for mock.
func (m *MockRepository) Install(names []string, password string) (string, error) {
	return "Package installed successfully", nil
}

// Remove is a no-op for mock.
func (m *MockRepository) Remove(names []string, cascade bool, password string) (string, error) {
	return "Mock removal output", nil
}

// Refresh is a no-op for mock.
func (m *MockRepository) Refresh() error {
	return nil
}

var _ Repository = (*MockRepository)(nil)
