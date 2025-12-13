package romm

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	URLCover    string    `json:"url_cover"`
	HasCover    bool      `json:"has_cover"`
	IsPublic    bool      `json:"is_public"`
	UserID      int       `json:"user_id"`
	ROMs        []Rom     `json:"roms"`
	ROMCount    int       `json:"rom_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SmartCollection struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	URLCover    string      `json:"url_cover"`
	HasCover    bool        `json:"has_cover"`
	IsPublic    bool        `json:"is_public"`
	UserID      int         `json:"user_id"`
	Criteria    interface{} `json:"criteria"`
	ROMs        []Rom       `json:"roms"`
	ROMCount    int         `json:"rom_count"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type VirtualCollection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	URLCover    string    `json:"url_cover"`
	HasCover    bool      `json:"has_cover"`
	Type        string    `json:"type"`
	ROMs        []Rom     `json:"roms"`
	ROMCount    int       `json:"rom_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) GetCollections() ([]Collection, error) {
	var collections []Collection
	err := c.doRequest("GET", EndpointCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetCollection(id int) (*Collection, error) {
	var collection Collection
	path := fmt.Sprintf(EndpointCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}

func (c *Client) GetSmartCollections() ([]SmartCollection, error) {
	var collections []SmartCollection
	err := c.doRequest("GET", EndpointSmartCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetSmartCollection(id int) (*SmartCollection, error) {
	var collection SmartCollection
	path := fmt.Sprintf(EndpointSmartCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}

func (c *Client) GetVirtualCollections() ([]VirtualCollection, error) {
	var collections []VirtualCollection
	err := c.doRequest("GET", EndpointVirtualCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetVirtualCollection(id int) (*VirtualCollection, error) {
	var collection VirtualCollection
	path := fmt.Sprintf(EndpointVirtualCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}
