package repository

import "github.com/sjsanc/pacviz/v3/internal/domain"

// Repository defines the interface for package data access.
type Repository interface {
	// GetInstalled returns all installed packages.
	GetInstalled() ([]*domain.Package, error)

	// GetExplicit returns explicitly installed packages.
	GetExplicit() ([]*domain.Package, error)

	// GetOrphans returns orphaned packages.
	GetOrphans() ([]*domain.Package, error)

	// GetForeign returns foreign packages (not in sync databases).
	GetForeign() ([]*domain.Package, error)

	// Search searches sync databases for packages matching the query.
	Search(query string) ([]*domain.Package, error)

	// GetPackage retrieves detailed information for a specific package.
	GetPackage(name string) (*domain.Package, error)

	// Install installs the specified packages.
	Install(names []string) error

	// Remove removes the specified packages.
	Remove(names []string, cascade bool) error

	// Refresh refreshes the package database.
	Refresh() error
}
