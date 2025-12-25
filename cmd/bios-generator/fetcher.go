package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	githubAPIURL = "https://api.github.com/repos/libretro/libretro-super/contents/dist/info"
	rawGitHubURL = "https://raw.githubusercontent.com/libretro/libretro-super/master/dist/info"
)

// GitHubFile represents a file entry from GitHub API
type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

// Fetcher handles downloading .info files from GitHub
type Fetcher struct {
	client *http.Client
}

// NewFetcher creates a new Fetcher instance
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchInfoFileList retrieves the list of .info files from GitHub
func (f *Fetcher) FetchInfoFileList() ([]string, error) {
	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add User-Agent header (GitHub API requires it)
	req.Header.Set("User-Agent", "grout-bios-generator")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching file list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("decoding GitHub API response: %w", err)
	}

	// Filter for .info files only
	var infoFiles []string
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, ".info") {
			infoFiles = append(infoFiles, file.Name)
		}
	}

	return infoFiles, nil
}

// FetchInfoFile downloads a single .info file from GitHub
func (f *Fetcher) FetchInfoFile(filename string) (string, error) {
	url := fmt.Sprintf("%s/%s", rawGitHubURL, filename)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request for %s: %w", filename, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching %s: %w", filename, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching %s returned status %d", filename, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", filename, err)
	}

	return string(body), nil
}

// FetchAllInfoFiles downloads all .info files
func (f *Fetcher) FetchAllInfoFiles() (map[string]string, error) {
	fileList, err := f.FetchInfoFileList()
	if err != nil {
		return nil, fmt.Errorf("fetching file list: %w", err)
	}

	fmt.Printf("Found %d .info files\n", len(fileList))

	infoFiles := make(map[string]string)
	for i, filename := range fileList {
		fmt.Printf("Fetching %d/%d: %s\n", i+1, len(fileList), filename)

		content, err := f.FetchInfoFile(filename)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to fetch %s: %v\n", filename, err)
			continue
		}

		infoFiles[filename] = content

		// Small delay to avoid rate limiting
		if i < len(fileList)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("Successfully fetched %d/%d files\n", len(infoFiles), len(fileList))
	return infoFiles, nil
}
