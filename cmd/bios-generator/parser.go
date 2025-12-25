package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// MD5 hash regex pattern: matches "(!) filename.bin (md5): a860e8c0b6d573d191e4ec7db1b1e4f6"
var md5Pattern = regexp.MustCompile(`\(!\)\s*([^(]+)\s*\(md5\):\s*([0-9a-fA-F]{32})`)

// ParsedBIOSFile represents a parsed BIOS file entry
type ParsedBIOSFile struct {
	FileName     string
	RelativePath string
	MD5Hash      string
	Description  string
	Optional     bool
}

// ParsedCore represents a parsed Libretro core with BIOS requirements
type ParsedCore struct {
	CoreName    string
	DisplayName string
	Files       []ParsedBIOSFile
}

// Parser handles parsing .info file content
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseInfoFile parses a single .info file and extracts BIOS information
func (p *Parser) ParseInfoFile(filename, content string) (*ParsedCore, error) {
	// Remove .info extension to get core name
	coreName := strings.TrimSuffix(filename, ".info")

	lines := strings.Split(content, "\n")
	keyValues := make(map[string]string)

	// Parse key-value pairs
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"`)

		keyValues[key] = value
	}

	// Extract display name
	displayName := keyValues["display_name"]
	if displayName == "" {
		displayName = coreName
	}

	// Extract MD5 hashes from notes
	notes := keyValues["notes"]
	md5Hashes := p.extractMD5Hashes(notes)

	// Parse firmware entries
	var biosFiles []ParsedBIOSFile

	// Count firmware entries
	firmwareCount := 0
	if count, err := strconv.Atoi(keyValues["firmware_count"]); err == nil {
		firmwareCount = count
	}

	// If firmware_count not specified, scan for firmware entries
	if firmwareCount == 0 {
		for key := range keyValues {
			if strings.HasPrefix(key, "firmware") && strings.HasSuffix(key, "_path") {
				firmwareCount++
			}
		}
	}

	// Extract each firmware entry
	for i := 0; i < firmwareCount; i++ {
		pathKey := fmt.Sprintf("firmware%d_path", i)
		descKey := fmt.Sprintf("firmware%d_desc", i)
		optKey := fmt.Sprintf("firmware%d_opt", i)

		path := keyValues[pathKey]
		if path == "" {
			continue // Skip if no path specified
		}

		desc := keyValues[descKey]
		if desc == "" {
			desc = path
		}

		optional := strings.ToLower(keyValues[optKey]) == "true"

		// Extract filename from path
		fileName := filepath.Base(path)

		// Look up MD5 hash for this file
		md5Hash := md5Hashes[fileName]

		biosFiles = append(biosFiles, ParsedBIOSFile{
			FileName:     fileName,
			RelativePath: path,
			MD5Hash:      md5Hash,
			Description:  desc,
			Optional:     optional,
		})
	}

	// Only return cores that have BIOS files
	if len(biosFiles) == 0 {
		return nil, nil
	}

	return &ParsedCore{
		CoreName:    coreName,
		DisplayName: displayName,
		Files:       biosFiles,
	}, nil
}

// extractMD5Hashes extracts MD5 hashes from the notes field
func (p *Parser) extractMD5Hashes(notes string) map[string]string {
	hashes := make(map[string]string)

	if notes == "" {
		return hashes
	}

	// Split by pipe character (notes often contain multiple entries separated by |)
	entries := strings.Split(notes, "|")

	for _, entry := range entries {
		matches := md5Pattern.FindStringSubmatch(entry)
		if len(matches) == 3 {
			fileName := strings.TrimSpace(matches[1])
			md5Hash := strings.ToLower(strings.TrimSpace(matches[2]))
			hashes[fileName] = md5Hash
		}
	}

	return hashes
}

// ParseAllInfoFiles parses all .info files
func (p *Parser) ParseAllInfoFiles(infoFiles map[string]string) ([]*ParsedCore, error) {
	var cores []*ParsedCore

	for filename, content := range infoFiles {
		core, err := p.ParseInfoFile(filename, content)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to parse %s: %v\n", filename, err)
			continue
		}

		if core != nil {
			cores = append(cores, core)
		}
	}

	fmt.Printf("Parsed %d cores with BIOS requirements\n", len(cores))
	return cores, nil
}
