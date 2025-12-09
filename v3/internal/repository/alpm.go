package repository

import (
	"fmt"
	"strings"

	"github.com/Jguer/go-alpm/v2"
	"github.com/Morganamilo/go-pacmanconf"
	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// AlpmRepository is the production implementation using go-alpm.
type AlpmRepository struct {
	handle  *alpm.Handle
	localDB alpm.IDB
	syncDBs alpm.IDBList
}

// NewAlpmRepository creates a new ALPM repository.
func NewAlpmRepository() (*AlpmRepository, error) {
	// Load pacman configuration
	pacmanConf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse pacman.conf: %w", err)
	}

	// Initialize ALPM handle
	handle, err := alpm.Initialize(pacmanConf.RootDir, pacmanConf.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ALPM: %w", err)
	}

	// Get local database
	localDB, err := handle.LocalDB()
	if err != nil {
		handle.Release()
		return nil, fmt.Errorf("failed to get local database: %w", err)
	}

	// Register sync databases from pacman.conf
	for _, repo := range pacmanConf.Repos {
		// Register each repository as a sync database
		_, err := handle.RegisterSyncDB(repo.Name, 0) // 0 = no signature check for now
		if err != nil {
			handle.Release()
			return nil, fmt.Errorf("failed to register sync database %s: %w", repo.Name, err)
		}
	}

	// Get sync databases
	syncDBs, err := handle.SyncDBs()
	if err != nil {
		handle.Release()
		return nil, fmt.Errorf("failed to get sync databases: %w", err)
	}

	return &AlpmRepository{
		handle:  handle,
		localDB: localDB,
		syncDBs: syncDBs,
	}, nil
}

// GetInstalled returns all installed packages.
func (r *AlpmRepository) GetInstalled() ([]*domain.Package, error) {
	pkgs := r.localDB.PkgCache()
	result := make([]*domain.Package, 0)

	pkgs.ForEach(func(pkg alpm.IPackage) error {
		result = append(result, r.convertPackage(pkg))
		return nil
	})

	// Compute orphan status (packages with no reverse dependencies)
	r.computeOrphans(result)

	// Compute foreign status (packages not in sync databases)
	r.computeForeign(result)

	return result, nil
}

// convertPackage converts an ALPM package to our domain package.
func (r *AlpmRepository) convertPackage(pkg alpm.IPackage) *domain.Package {
	installReason := domain.ReasonDependency
	if pkg.Reason() == alpm.PkgReasonExplicit {
		installReason = domain.ReasonExplicit
	}

	// Extract dependency names
	deps := make([]string, 0)
	pkg.Depends().ForEach(func(dep *alpm.Depend) error {
		deps = append(deps, dep.Name)
		return nil
	})

	return &domain.Package{
		Name:          pkg.Name(),
		Version:       pkg.Version(),
		Description:   pkg.Description(),
		Architecture:  pkg.Architecture(),
		URL:           pkg.URL(),
		Licenses:      pkg.Licenses().Slice(),
		Groups:        pkg.Groups().Slice(),
		Dependencies:  deps,
		Installed:     true,
		InstallDate:   pkg.InstallDate(),
		InstallReason: installReason,
		InstalledSize: pkg.ISize(),
		Packager:      pkg.Packager(),
		BuildDate:     pkg.BuildDate(),
	}
}

// computeOrphans calculates which packages are orphans (dependencies with no dependents).
func (r *AlpmRepository) computeOrphans(packages []*domain.Package) {
	// Build a map of package names to packages
	pkgMap := make(map[string]*domain.Package)
	for _, pkg := range packages {
		pkgMap[pkg.Name] = pkg
	}

	// Build reverse dependency map
	reverseDeps := make(map[string][]string)
	for _, pkg := range packages {
		for _, dep := range pkg.Dependencies {
			reverseDeps[dep] = append(reverseDeps[dep], pkg.Name)
		}
	}

	// Mark orphans: dependency packages with no reverse dependencies
	for _, pkg := range packages {
		if pkg.InstallReason == domain.ReasonDependency {
			pkg.IsOrphan = len(reverseDeps[pkg.Name]) == 0
		}
		pkg.Required = reverseDeps[pkg.Name]
	}
}

// computeForeign marks packages that are not in any sync database as foreign
// and populates the repository name for all packages.
func (r *AlpmRepository) computeForeign(packages []*domain.Package) {
	// Build a map of package names to their repository
	pkgToRepo := make(map[string]string)
	r.syncDBs.ForEach(func(db alpm.IDB) error {
		repoName := db.Name()
		db.PkgCache().ForEach(func(pkg alpm.IPackage) error {
			pkgToRepo[pkg.Name()] = repoName
			return nil
		})
		return nil
	})

	// Set repository name and mark foreign packages
	for _, pkg := range packages {
		if repo, exists := pkgToRepo[pkg.Name]; exists {
			pkg.Repository = repo
			pkg.IsForeign = false
		} else {
			// Package not in any sync database
			pkg.Repository = "foreign"
			pkg.IsForeign = true
		}
	}
}

// GetExplicit returns explicitly installed packages.
func (r *AlpmRepository) GetExplicit() ([]*domain.Package, error) {
	// TODO: Filter by InstallReason == Explicit
	return nil, nil
}

// GetOrphans returns orphaned packages.
func (r *AlpmRepository) GetOrphans() ([]*domain.Package, error) {
	// TODO: Find packages with no reverse dependencies
	return nil, nil
}

// GetForeign returns foreign packages.
func (r *AlpmRepository) GetForeign() ([]*domain.Package, error) {
	// TODO: Find packages not in any sync database
	return nil, nil
}

// Search searches sync databases for packages matching the query.
func (r *AlpmRepository) Search(query string) ([]*domain.Package, error) {
	result := make([]*domain.Package, 0)

	// Search each sync database
	r.syncDBs.ForEach(func(db alpm.IDB) error {
		repoName := db.Name()
		db.PkgCache().ForEach(func(pkg alpm.IPackage) error {
			// Search in name and description
			name := pkg.Name()
			desc := pkg.Description()
			queryLower := strings.ToLower(query)
			if strings.Contains(strings.ToLower(name), queryLower) ||
				strings.Contains(strings.ToLower(desc), queryLower) {
				p := r.convertSyncPackage(pkg, repoName)
				result = append(result, p)
			}
			return nil
		})
		return nil
	})

	return result, nil
}

// convertSyncPackage converts an ALPM sync package to our domain package.
func (r *AlpmRepository) convertSyncPackage(pkg alpm.IPackage, repoName string) *domain.Package {
	// Extract dependency names
	deps := make([]string, 0)
	pkg.Depends().ForEach(func(dep *alpm.Depend) error {
		deps = append(deps, dep.Name)
		return nil
	})

	// Check if package is installed locally
	localPkg := r.localDB.Pkg(pkg.Name())
	installed := localPkg != nil

	p := &domain.Package{
		Name:          pkg.Name(),
		Version:       pkg.Version(),
		Description:   pkg.Description(),
		Architecture:  pkg.Architecture(),
		URL:           pkg.URL(),
		Licenses:      pkg.Licenses().Slice(),
		Groups:        pkg.Groups().Slice(),
		Dependencies:  deps,
		Installed:     installed,
		InstalledSize: pkg.ISize(),
		Packager:      pkg.Packager(),
		BuildDate:     pkg.BuildDate(),
		Repository:    repoName,
		IsForeign:     false,
	}

	// If installed, get install-specific info
	if localPkg != nil {
		p.InstallDate = localPkg.InstallDate()
		if localPkg.Reason() == alpm.PkgReasonExplicit {
			p.InstallReason = domain.ReasonExplicit
		} else {
			p.InstallReason = domain.ReasonDependency
		}
	}

	return p
}

// GetPackage retrieves package details.
func (r *AlpmRepository) GetPackage(name string) (*domain.Package, error) {
	// TODO: Query package by name
	return nil, nil
}

// Install installs packages.
func (r *AlpmRepository) Install(names []string) error {
	// TODO: Install via pacman command or ALPM transaction
	return nil
}

// Remove removes packages.
func (r *AlpmRepository) Remove(names []string, cascade bool) error {
	// TODO: Remove via pacman command or ALPM transaction
	return nil
}

// Refresh refreshes the package database.
func (r *AlpmRepository) Refresh() error {
	// TODO: Run pacman -Sy or ALPM sync
	return nil
}
