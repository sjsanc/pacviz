package main

import (
	"github.com/Jguer/go-alpm/v2"
	paconf "github.com/Morganamilo/go-pacmanconf"
)

type Pacman struct {
	DB      alpm.IDB
	SyncDBs alpm.IDBList
	Pkgs    []alpm.IPackage
}

func NewPacman() *Pacman {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		panic(err)
	}

	conf, _, err := paconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		panic(err)
	}

	for _, repo := range conf.Repos {
		db, err := h.RegisterSyncDB(repo.Name, 0)
		if err != nil {
			panic(err)
		}
		db.SetServers(repo.Servers)

		if len(repo.Usage) == 0 {
			db.SetUsage(alpm.UsageAll)
		}
		for _, usage := range repo.Usage {
			switch usage {
			case "Sync":
				db.SetUsage(alpm.UsageSync)
			case "Search":
				db.SetUsage(alpm.UsageSearch)
			case "Install":
				db.SetUsage(alpm.UsageInstall)
			case "Upgrade":
				db.SetUsage(alpm.UsageUpgrade)
			case "All":
				db.SetUsage(alpm.UsageAll)
			}
		}
	}

	db, err := h.LocalDB()
	if err != nil {
		panic(err)
	}

	pkgs := db.PkgCache().Slice()

	syncdbs, err := h.SyncDBs()
	if err != nil {
		panic(err)
	}

	return &Pacman{
		DB:      db,
		SyncDBs: syncdbs,
		Pkgs:    pkgs,
	}
}

// TODO: PkgCache is an expensive operation. It can be cached, and this function debounced.
func (pm *Pacman) SearchSyncDBs(terms ...string) []alpm.IPackage {
	var result []alpm.IPackage
	for _, db := range pm.SyncDBs.Slice() {
		if len(terms) == 0 {
			result = append(result, db.PkgCache().Slice()...)
		} else {
			pkgs := db.Search(terms)
			result = append(result, pkgs.Slice()...)
		}
	}
	return result
}

func (pm *Pacman) GetInstalledPkgs() []alpm.IPackage {
	return pm.Pkgs
}

func (pm *Pacman) GetExplicitPkgs() []alpm.IPackage {
	var result []alpm.IPackage
	for _, pkg := range pm.Pkgs {
		if pkg.Reason() == alpm.PkgReasonExplicit {
			result = append(result, pkg)
		}
	}
	return result
}

func (pm *Pacman) GetOrphanPkgs() []alpm.IPackage {
	var result []alpm.IPackage
	for _, pkg := range pm.Pkgs {
		if pkg.Reason() == alpm.PkgReasonDepend && len(pkg.ComputeRequiredBy()) == 0 {
			result = append(result, pkg)
		}
	}
	return result
}
