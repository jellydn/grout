package constants

// BIOSFile represents a single ShowBIOSDownload/firmware file requirement
type BIOSFile struct {
	FileName     string // e.g., "gba_bios.bin"
	RelativePath string // e.g., "gba_bios.bin" or "psx/scph5500.bin"
	MD5Hash      string // e.g., "a860e8c0b6d573d191e4ec7db1b1e4f6" (optional, empty string if unknown)
	Optional     bool   // true if ShowBIOSDownload file is optional for the emulator to function
}

// CoreBIOS represents all ShowBIOSDownload requirements for a Libretro core
type CoreBIOS struct {
	CoreName    string     // e.g., "mgba_libretro"
	DisplayName string     // e.g., "Nintendo - Game Boy Advance (mGBA)"
	Files       []BIOSFile // List of ShowBIOSDownload files for this core
}

// CoreBIOSSubdirectories maps Libretro core names (without _libretro suffix)
// to their required ShowBIOSDownload subdirectory within the system ShowBIOSDownload directory.
// Cores not in this map use the root ShowBIOSDownload directory.
var CoreBIOSSubdirectories = mustLoadJSONMap[string, string]("bios/core_subdirectories.json")
