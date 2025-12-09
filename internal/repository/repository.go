package repository

import "github.com/sjsanc/pacviz/v3/internal/domain"

// Repository defines the interface for package data access.
type Repository interface {
	// GetInstalled returns all installed packages.
	GetInstalled() ([]*domain.Package, error)

	// Search searches sync databases for packages matching the query.
	Search(query string) ([]*domain.Package, error)

	// Install installs the specified packages and returns the command output.
	Install(names []string, password string) (string, error)

	// Remove removes the specified packages and returns the command output.
	Remove(names []string, cascade bool, password string) (string, error)

	// Refresh refreshes the package database.
	Refresh() error
}
