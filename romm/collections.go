package romm

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URLCover    string    `json:"url_cover"`
	HasCover    bool      `json:"has_cover"`
	IsPublic    bool      `json:"is_public"`
	UserID      int       `json:"user_id"`
	ROMs        []Rom     `json:"roms"`
	ROMCount    int       `json:"rom_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) GetCollections() ([]Collection, error) {
	var collections []Collection
	err := c.doRequest("GET", endpointCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetCollection(id int) (Collection, error) {
	var collection Collection
	path := fmt.Sprintf(endpointCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return collection, err
}

func (c *Client) GetSmartCollections() ([]Collection, error) {
	var collections []Collection
	err := c.doRequest("GET", endpointSmartCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetSmartCollection(id int) (Collection, error) {
	var collection Collection
	path := fmt.Sprintf(endpointSmartCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return collection, err
}

func (c *Client) GetVirtualCollections() ([]Collection, error) {
	var collections []Collection
	err := c.doRequest("GET", endpointVirtualCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) GetVirtualCollection(id int) (Collection, error) {
	var collection Collection
	path := fmt.Sprintf(endpointVirtualCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return collection, err
}
