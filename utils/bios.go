package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"grout/constants"
)

// SaveBIOSFile saves BIOS file data to all appropriate CFW-specific directories
func SaveBIOSFile(biosFile constants.BIOSFile, platformSlug string, data []byte) error {
	filePaths := GetBIOSFilePaths(biosFile, platformSlug)

	// Save to all target paths
	for _, filePath := range filePaths {
		// Create parent directories if they don't exist
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

// VerifyBIOSFileMD5 verifies the MD5 hash of a BIOS file
func VerifyBIOSFileMD5(data []byte, expectedMD5 string) (bool, string) {
	if expectedMD5 == "" {
		// No MD5 hash to verify against
		return true, ""
	}

	hash := md5.Sum(data)
	actualMD5 := hex.EncodeToString(hash[:])

	return actualMD5 == expectedMD5, actualMD5
}

// BIOSFileExists checks if a BIOS file exists at any of the expected locations
func BIOSFileExists(biosFile constants.BIOSFile, platformSlug string) bool {
	filePaths := GetBIOSFilePaths(biosFile, platformSlug)

	// Return true if file exists in any of the paths
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); err == nil {
			return true
		}
	}

	return false
}

// GetBIOSFileInfo returns information about an existing BIOS file from the first location where it exists
func GetBIOSFileInfo(biosFile constants.BIOSFile, platformSlug string) (exists bool, size int64, md5Hash string, err error) {
	filePaths := GetBIOSFilePaths(biosFile, platformSlug)

	// Check each path and return info from the first existing file
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return false, 0, "", err
		}

		// File exists, calculate MD5 hash
		file, err := os.Open(filePath)
		if err != nil {
			return true, info.Size(), "", err
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return true, info.Size(), "", err
		}

		md5Hash = hex.EncodeToString(hash.Sum(nil))

		return true, info.Size(), md5Hash, nil
	}

	// File doesn't exist in any of the paths
	return false, 0, "", nil
}

// GetBIOSFilesForPlatform returns all BIOS files required for a given platform slug
func GetBIOSFilesForPlatform(platformSlug string) []constants.BIOSFile {
	var biosFiles []constants.BIOSFile

	// Get all cores for this platform
	coreNames, ok := constants.PlatformToLibretroCores[platformSlug]
	if !ok {
		return biosFiles
	}

	// Collect all BIOS files from these cores (deduplicate by filename)
	seen := make(map[string]bool)
	for _, coreName := range coreNames {
		// Normalize core name by removing _libretro suffix
		normalizedCoreName := strings.TrimSuffix(coreName, "_libretro")
		coreInfo, ok := constants.LibretroCoreToBIOS[normalizedCoreName]
		if !ok {
			continue
		}

		for _, file := range coreInfo.Files {
			// Use filename as unique key to avoid duplicates
			if !seen[file.FileName] {
				biosFiles = append(biosFiles, file)
				seen[file.FileName] = true
			}
		}
	}

	return biosFiles
}

// BIOSStatus represents the status of a BIOS file
type BIOSStatus string

const (
	BIOSStatusMissing        BIOSStatus = "missing"
	BIOSStatusValid          BIOSStatus = "valid"
	BIOSStatusInvalidHash    BIOSStatus = "invalid_hash"
	BIOSStatusNoHashToVerify BIOSStatus = "no_hash"
)

// BIOSFileStatus contains information about the status of a BIOS file
type BIOSFileStatus struct {
	File        constants.BIOSFile
	Status      BIOSStatus
	Exists      bool
	Size        int64
	ActualMD5   string
	ExpectedMD5 string
}

// CheckBIOSFileStatus checks the status of a BIOS file
func CheckBIOSFileStatus(biosFile constants.BIOSFile, platformSlug string) BIOSFileStatus {
	status := BIOSFileStatus{
		File:        biosFile,
		ExpectedMD5: biosFile.MD5Hash,
	}

	exists, size, actualMD5, err := GetBIOSFileInfo(biosFile, platformSlug)
	if err != nil {
		status.Status = BIOSStatusMissing
		return status
	}

	status.Exists = exists
	status.Size = size
	status.ActualMD5 = actualMD5

	if !exists {
		status.Status = BIOSStatusMissing
		return status
	}

	if biosFile.MD5Hash == "" {
		status.Status = BIOSStatusNoHashToVerify
		return status
	}

	if actualMD5 == biosFile.MD5Hash {
		status.Status = BIOSStatusValid
	} else {
		status.Status = BIOSStatusInvalidHash
	}

	return status
}
