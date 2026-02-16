package aur

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sjsanc/pacviz/v3/internal/domain"
)

const (
	baseURL        = "https://aur.archlinux.org/rpc/v5"
	defaultTimeout = 5 * time.Second
	defaultTTL     = 5 * time.Minute
	infoBatchSize  = 150
)

// AURResponse represents the JSON response from the AUR RPC API.
type AURResponse struct {
	Version     int          `json:"version"`
	Type        string       `json:"type"`
	ResultCount int          `json:"resultcount"`
	Results     []AURPackage `json:"results"`
	Error       string       `json:"error"`
}

// AURPackage represents a single package from the AUR RPC API.
type AURPackage struct {
	Name           string   `json:"Name"`
	Version        string   `json:"Version"`
	Description    string   `json:"Description"`
	URL            string   `json:"URL"`
	PackageBase    string   `json:"PackageBase"`
	Maintainer     string   `json:"Maintainer"`
	NumVotes       int      `json:"NumVotes"`
	Popularity     float64  `json:"Popularity"`
	OutOfDate      *int64   `json:"OutOfDate"`
	FirstSubmitted int64    `json:"FirstSubmitted"`
	LastModified   int64    `json:"LastModified"`
	License        []string `json:"License"`
	Depends        []string `json:"Depends"`
	MakeDepends    []string `json:"MakeDepends"`
	OptDepends     []string `json:"OptDepends"`
	Conflicts      []string `json:"Conflicts"`
	Provides       []string `json:"Provides"`
	Replaces       []string `json:"Replaces"`
	Groups         []string `json:"Groups"`
}

type cacheEntry struct {
	data      any
	expiresAt time.Time
}

// Client is an HTTP client for the AUR RPC API.
type Client struct {
	httpClient *http.Client
	cacheTTL   time.Duration
	cache      map[string]cacheEntry
	mu         sync.RWMutex
}

// NewClient creates a new AUR RPC client.
func NewClient(timeout time.Duration, cacheTTL time.Duration) *Client {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	if cacheTTL == 0 {
		cacheTTL = defaultTTL
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		cacheTTL:   cacheTTL,
		cache:      make(map[string]cacheEntry),
	}
}

// Search queries the AUR for packages matching the given query.
func (c *Client) Search(query string) ([]*domain.Package, error) {
	cacheKey := "search:" + query
	if cached := c.getCache(cacheKey); cached != nil {
		return cached.([]*domain.Package), nil
	}

	reqURL := fmt.Sprintf("%s/search/%s", baseURL, url.PathEscape(query))
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("AUR search request failed: %w", err)
	}
	defer resp.Body.Close()

	var aurResp AURResponse
	if err := json.NewDecoder(resp.Body).Decode(&aurResp); err != nil {
		return nil, fmt.Errorf("failed to decode AUR response: %w", err)
	}

	if aurResp.Error != "" {
		return nil, fmt.Errorf("AUR API error: %s", aurResp.Error)
	}

	packages := make([]*domain.Package, 0, len(aurResp.Results))
	for _, ap := range aurResp.Results {
		packages = append(packages, convertToDomainPackage(ap))
	}

	c.setCache(cacheKey, packages)
	return packages, nil
}

// Info queries the AUR for specific package names and returns a set of names that exist.
func (c *Client) Info(names []string) (map[string]bool, error) {
	if len(names) == 0 {
		return map[string]bool{}, nil
	}

	result := make(map[string]bool)

	// Batch in groups of infoBatchSize
	for i := 0; i < len(names); i += infoBatchSize {
		end := min(i+infoBatchSize, len(names))
		batch := names[i:end]

		found, err := c.infoBatch(batch)
		if err != nil {
			return nil, err
		}
		for name := range found {
			result[name] = true
		}
	}

	return result, nil
}

func (c *Client) infoBatch(names []string) (map[string]bool, error) {
	cacheKey := "info:" + strings.Join(names, ",")
	if cached := c.getCache(cacheKey); cached != nil {
		return cached.(map[string]bool), nil
	}

	params := url.Values{}
	for _, name := range names {
		params.Add("arg[]", name)
	}

	reqURL := fmt.Sprintf("%s/info?%s", baseURL, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("AUR info request failed: %w", err)
	}
	defer resp.Body.Close()

	var aurResp AURResponse
	if err := json.NewDecoder(resp.Body).Decode(&aurResp); err != nil {
		return nil, fmt.Errorf("failed to decode AUR info response: %w", err)
	}

	if aurResp.Error != "" {
		return nil, fmt.Errorf("AUR API error: %s", aurResp.Error)
	}

	found := make(map[string]bool, len(aurResp.Results))
	for _, pkg := range aurResp.Results {
		found[pkg.Name] = true
	}

	c.setCache(cacheKey, found)
	return found, nil
}

func (c *Client) getCache(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}
	return entry.data
}

func (c *Client) setCache(key string, data any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.cacheTTL),
	}
}

func convertToDomainPackage(ap AURPackage) *domain.Package {
	deps := make([]string, len(ap.Depends))
	copy(deps, ap.Depends)

	optDeps := make(map[string]string)
	for _, od := range ap.OptDepends {
		parts := strings.SplitN(od, ": ", 2)
		if len(parts) == 2 {
			optDeps[parts[0]] = parts[1]
		} else {
			optDeps[parts[0]] = ""
		}
	}

	return &domain.Package{
		Name:         ap.Name,
		Version:      ap.Version,
		Description:  ap.Description,
		URL:          ap.URL,
		Licenses:     ap.License,
		Groups:       ap.Groups,
		Dependencies: deps,
		OptDepends:   optDeps,
		Conflicts:    ap.Conflicts,
		Provides:     ap.Provides,
		Replaces:     ap.Replaces,
		Repository:   "aur",
		IsAUR:        true,
		Packager:     ap.Maintainer,
		BuildDate:    time.Unix(ap.LastModified, 0),
	}
}
