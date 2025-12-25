package constants

//go:generate go run ../cmd/bios-generator/main.go

// BIOSFile represents a single BIOS/firmware file requirement
type BIOSFile struct {
	FileName     string // e.g., "gba_bios.bin"
	RelativePath string // e.g., "gba_bios.bin" or "psx/scph5500.bin"
	MD5Hash      string // e.g., "a860e8c0b6d573d191e4ec7db1b1e4f6" (optional)
	Description  string // e.g., "Game Boy Advance BIOS"
	Optional     bool   // true if firmware*_opt = "true"
}

// CoreBIOS represents all BIOS requirements for a Libretro core
type CoreBIOS struct {
	CoreName    string     // e.g., "mgba_libretro"
	DisplayName string     // e.g., "Nintendo - Game Boy Advance (mGBA)"
	Files       []BIOSFile // List of BIOS files for this core
}
