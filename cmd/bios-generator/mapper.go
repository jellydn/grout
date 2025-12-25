package main

import (
	"fmt"
	"strings"
)

// Mapper handles mapping core names to platform slugs
type Mapper struct {
	coreToSlug map[string][]string
}

// NewMapper creates a new Mapper with predefined mappings
func NewMapper() *Mapper {
	return &Mapper{
		coreToSlug: map[string][]string{
			// Game Boy / Game Boy Color / Game Boy Advance
			"mgba":     {"gba"},
			"gpsp":     {"gba"},
			"vba_next": {"gba", "gb", "gbc"},
			"vbam":     {"gba", "gb", "gbc"},
			"gambatte": {"gb", "gbc"},
			"sameboy":  {"gb", "gbc"},
			"gearboy":  {"gb", "gbc"},
			"tgbdual":  {"gb", "gbc"},
			"mesen-s":  {"snes", "sfam", "gb", "gbc"},

			// PlayStation / PS2
			"pcsx_rearmed":  {"psx", "ps"},
			"beetle_psx":    {"psx", "ps"},
			"beetle_psx_hw": {"psx", "ps"},
			"duckstation":   {"psx", "ps"},
			"swanstation":   {"psx", "ps"},
			"play":          {"ps2"},
			"pcsx2":         {"ps2"},

			// SNES / Super Famicom
			"snes9x":                    {"snes", "sfam"},
			"snes9x2002":                {"snes", "sfam"},
			"snes9x2005":                {"snes", "sfam"},
			"snes9x2010":                {"snes", "sfam"},
			"bsnes":                     {"snes", "sfam"},
			"bsnes_hd_beta":             {"snes", "sfam"},
			"bsnes_mercury_accuracy":    {"snes", "sfam"},
			"bsnes_mercury_balanced":    {"snes", "sfam"},
			"bsnes_mercury_performance": {"snes", "sfam"},
			"mednafen_snes":             {"snes", "sfam"},

			// Genesis / Mega Drive / Sega CD / Master System / Game Gear
			"genesis_plus_gx":      {"genesis", "sms", "gamegear", "sg1000", "segacd"},
			"genesis_plus_gx_wide": {"genesis", "sms", "gamegear", "sg1000", "segacd"},
			"picodrive":            {"genesis", "sega32", "segacd", "sms"},
			"blastem":              {"genesis"},

			// NES / Famicom / Famicom Disk System
			"fceumm":   {"nes", "famicom", "fds"},
			"mesen":    {"nes", "famicom"},
			"nestopia": {"nes", "famicom"},
			"quicknes": {"nes", "famicom"},
			"bnes":     {"nes", "famicom"},

			// Nintendo 64
			"mupen64plus_next": {"n64"},
			"parallel_n64":     {"n64"},

			// Dreamcast / Naomi / Atomiswave
			"flycast": {"dc", "naomi", "atomiswave"},
			"redream": {"dc"},
			"reicast": {"dc", "naomi"},

			// Saturn
			"beetle_saturn": {"saturn"},
			"yabasanshiro":  {"saturn"},
			"yabause":       {"saturn"},
			"kronos":        {"saturn"},

			// Atari systems
			"stella":     {"atari2600"},
			"stella2014": {"atari2600"},
			"atari800":   {"atari5200", "atari800"},
			"prosystem":  {"atari7800"},

			// Neo Geo / Neo Geo Pocket
			"fbneo":         {"arcade", "neogeo", "fba"},
			"fbalpha":       {"arcade", "neogeo", "fba"},
			"fbalpha2012":   {"arcade", "neogeo", "fba"},
			"mame":          {"arcade"},
			"mame2000":      {"arcade"},
			"mame2003":      {"arcade"},
			"mame2003_plus": {"arcade"},
			"mame2010":      {"arcade"},
			"mame2015":      {"arcade"},
			"mame2016":      {"arcade"},
			"mednafen_ngp":  {"ngp", "ngpc"},
			"race":          {"ngp", "ngpc"},

			// PC Engine / TurboGrafx-16 / SuperGrafx
			"beetle_pce":        {"pce", "tg16", "sgfx"},
			"beetle_pce_fast":   {"pce", "tg16", "sgfx"},
			"mednafen_pce_fast": {"pce", "tg16", "sgfx"},
			"beetle_supergrafx": {"sgfx", "pce"},

			// Lynx
			"handy":         {"lynx"},
			"beetle_lynx":   {"lynx"},
			"mednafen_lynx": {"lynx"},

			// WonderSwan
			"beetle_wswan":   {"wswan", "wsc"},
			"mednafen_wswan": {"wswan", "wsc"},

			// Virtual Boy
			"beetle_vb":   {"vb"},
			"mednafen_vb": {"vb"},

			// 3DO
			"opera": {"3do"},

			// Jaguar
			"virtualjaguar": {"jaguar"},

			// Colecovision
			"bluemsx":    {"colecovision", "msx", "msx2"},
			"gearcoleco": {"colecovision"},

			// Intellivision
			"freeintv": {"intellivision"},

			// Vectrex
			"vecx": {"vectrex"},

			// Odyssey 2
			"o2em": {"odyssey2"},

			// ZX Spectrum
			"fuse": {"zxspectrum"},

			// Commodore systems
			"vice_x64":    {"c64"},
			"vice_x64sc":  {"c64"},
			"vice_x128":   {"c128"},
			"vice_xpet":   {"cpet"},
			"vice_xvic":   {"vic20"},
			"vice_xplus4": {"plus4"},

			// Amiga
			"puae":    {"amiga"},
			"uae4arm": {"amiga"},

			// DOS
			"dosbox":      {"dos"},
			"dosbox_pure": {"dos"},
			"dosbox_svn":  {"dos"},

			// MSX
			"fmsx": {"msx", "msx2"},

			// PC-FX
			"beetle_pcfx":   {"pcfx"},
			"mednafen_pcfx": {"pcfx"},

			// GameCube / Wii
			"dolphin": {"gc", "wii"},

			// PSP
			"ppsspp": {"psp"},

			// Nintendo DS
			"desmume":     {"nds"},
			"desmume2015": {"nds"},
			"melonds":     {"nds"},

			// Nintendo 3DS
			"citra": {"3ds"},

			// Game & Watch
			"gw": {"g-and-w"},

			// ColecoVision
			"smsplus": {"sms", "gamegear"},

			// Channel F
			"freechaf": {"fairchild-channel-f"},
		},
	}
}

// InferPlatformSlugs returns platform slugs for a given core name
func (m *Mapper) InferPlatformSlugs(coreName string) []string {
	// Normalize core name (remove _libretro suffix)
	baseName := strings.TrimSuffix(coreName, "_libretro")

	// Look up in static mapping
	if slugs, ok := m.coreToSlug[baseName]; ok {
		return slugs
	}

	// Pattern-based fallback (try to find platform name in core name)
	knownPlatforms := []string{
		"gba", "gb", "gbc", "nes", "snes", "n64", "gc", "wii",
		"genesis", "saturn", "dc", "psx", "ps2", "psp",
		"arcade", "mame", "fba", "neogeo",
		"atari", "lynx", "jaguar",
		"pce", "tg16", "sgfx",
		"ngp", "ngpc", "wswan", "wsc",
		"3do", "3ds", "nds",
		"amiga", "c64", "msx", "dos",
	}

	for _, platform := range knownPlatforms {
		if strings.Contains(baseName, platform) {
			return []string{platform}
		}
	}

	return nil
}

// BuildPlatformToCoresMap creates a reverse mapping from platform slug to cores
func (m *Mapper) BuildPlatformToCoresMap(cores []*ParsedCore) map[string][]string {
	platformToCores := make(map[string][]string)

	for _, core := range cores {
		slugs := m.InferPlatformSlugs(core.CoreName)

		for _, slug := range slugs {
			// Add core to this platform's list if not already present
			if !contains(platformToCores[slug], core.CoreName) {
				platformToCores[slug] = append(platformToCores[slug], core.CoreName)
			}
		}
	}

	return platformToCores
}

// contains checks if a string slice contains a given string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// NormalizeCoreNames creates a map with normalized core names (without _libretro suffix)
func (m *Mapper) NormalizeCoreNames(cores []*ParsedCore) map[string]*ParsedCore {
	normalized := make(map[string]*ParsedCore)

	for _, core := range cores {
		baseName := strings.TrimSuffix(core.CoreName, "_libretro")
		normalized[baseName] = core
	}

	fmt.Printf("Normalized %d core names\n", len(normalized))
	return normalized
}
