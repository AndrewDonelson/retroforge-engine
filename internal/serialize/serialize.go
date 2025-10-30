package serialize

import (
    "encoding/json"
    "github.com/AndrewDonelson/retroforge-engine/internal/cart"
)

// MetadataToJSON encodes cart metadata to bytes.
func MetadataToJSON(m cart.Metadata) ([]byte, error) {
    if err := m.Validate(); err != nil { return nil, err }
    return json.MarshalIndent(m, "", "  ")
}

// MetadataFromJSON decodes metadata from bytes.
func MetadataFromJSON(b []byte) (cart.Metadata, error) {
    var m cart.Metadata
    if err := json.Unmarshal(b, &m); err != nil { return m, err }
    return m, m.Validate()
}


