package resources

import (
	"encoding/json"
	"fmt"
)

// Asset is one catalogue entry (exact-ID lookup).
type Asset struct {
	ID     string
	Type   string
	Status string
}

// Catalogue is the embedded exact-ID asset table.
type Catalogue struct {
	PackVersion string
	ByID        map[string]Asset
}

type catalogueFile struct {
	AssetPackVersion string `json:"asset_pack_version"`
	Assets           []struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
	} `json:"assets"`
}

// LoadCatalogue returns the embedded exact-ID catalogue.
func LoadCatalogue() (Catalogue, error) {
	data, err := contentFS.ReadFile("catalogue/v1.json")
	if err != nil {
		return Catalogue{}, err
	}
	var file catalogueFile
	if err := json.Unmarshal(data, &file); err != nil {
		return Catalogue{}, err
	}
	out := Catalogue{
		PackVersion: file.AssetPackVersion,
		ByID:        make(map[string]Asset, len(file.Assets)),
	}
	for _, a := range file.Assets {
		if a.ID == "" {
			return Catalogue{}, fmt.Errorf("catalogue asset missing id")
		}
		out.ByID[a.ID] = Asset{
			ID:     a.ID,
			Type:   a.Type,
			Status: a.Status,
		}
	}
	return out, nil
}

// SchemaJSON returns the embedded JSON Schema bytes.
func SchemaJSON() ([]byte, error) {
	return contentFS.ReadFile("schema/v1.json")
}

// CatalogueJSON returns the embedded catalogue bytes.
func CatalogueJSON() ([]byte, error) {
	return contentFS.ReadFile("catalogue/v1.json")
}
