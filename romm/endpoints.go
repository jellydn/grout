package romm

const (
	endpointLogin = "/api/login"

	endpointPlatforms    = "/api/platforms"
	endpointPlatformByID = "/api/platforms/%d"

	endpointRoms         = "/api/roms"
	endpointRomByID      = "/api/roms/%d"
	endpointRomsDownload = "/api/roms/download"
	endpointRomsByHash   = "/api/roms/by-hash"

	endpointCollections           = "/api/collections"
	endpointCollectionByID        = "/api/collections/%d"
	endpointSmartCollections      = "/api/collections/smart"
	endpointSmartCollectionByID   = "/api/collections/smart/%d"
	endpointVirtualCollections    = "/api/collections/virtual"
	endpointVirtualCollectionByID = "/api/collections/virtual/%d"

	endpointFirmware = "/api/firmware"

	endpointSaves = "/api/saves"
)
