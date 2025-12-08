package repository

import (
	"time"

	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// CachedRepository wraps a repository with caching.
type CachedRepository struct {
	repo Repository

	// Cache storage
	installedCache   []*domain.Package
	installedCacheTS time.Time

	searchCache map[string][]*domain.Package
	cacheTTL    time.Duration
}

// NewCachedRepository creates a cached repository wrapper.
func NewCachedRepository(repo Repository, ttl time.Duration) *CachedRepository {
	return &CachedRepository{
		repo:        repo,
		searchCache: make(map[string][]*domain.Package),
		cacheTTL:    ttl,
	}
}

// GetInstalled returns cached installed packages.
func (c *CachedRepository) GetInstalled() ([]*domain.Package, error) {
	// TODO: Check cache validity
	// TODO: Return cached if valid, otherwise query and cache
	return c.repo.GetInstalled()
}

// GetExplicit returns cached explicit packages.
func (c *CachedRepository) GetExplicit() ([]*domain.Package, error) {
	// TODO: Use cached installed packages and filter
	return c.repo.GetExplicit()
}

// GetOrphans returns cached orphan packages.
func (c *CachedRepository) GetOrphans() ([]*domain.Package, error) {
	// TODO: Use cached installed packages and filter
	return c.repo.GetOrphans()
}

// GetForeign returns cached foreign packages.
func (c *CachedRepository) GetForeign() ([]*domain.Package, error) {
	// TODO: Use cached installed packages and filter
	return c.repo.GetForeign()
}

// Search returns cached search results.
func (c *CachedRepository) Search(query string) ([]*domain.Package, error) {
	// TODO: Check search cache
	// TODO: Return cached if valid, otherwise query and cache
	return c.repo.Search(query)
}

// GetPackage retrieves package (with caching).
func (c *CachedRepository) GetPackage(name string) (*domain.Package, error) {
	// TODO: LRU cache for individual packages
	return c.repo.GetPackage(name)
}

// Install invalidates cache and installs.
func (c *CachedRepository) Install(names []string) error {
	// TODO: Invalidate installed cache after install
	return c.repo.Install(names)
}

// Remove invalidates cache and removes.
func (c *CachedRepository) Remove(names []string, cascade bool) error {
	// TODO: Invalidate installed cache after remove
	return c.repo.Remove(names, cascade)
}

// Refresh invalidates all caches and refreshes.
func (c *CachedRepository) Refresh() error {
	// TODO: Clear all caches
	c.searchCache = make(map[string][]*domain.Package)
	c.installedCache = nil
	return c.repo.Refresh()
}
