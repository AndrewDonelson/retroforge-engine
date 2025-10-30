package serialize

import (
    "testing"
    "github.com/AndrewDonelson/retroforge-engine/internal/cart"
)

func TestMetadataRoundTrip(t *testing.T) {
    in := cart.Metadata{
        Title: "Test",
        Author: "Dev",
        Version: "0.1.0",
        Description: "desc",
        Genre: "Arcade",
        Tags: []string{"retro","demo"},
    }
    b, err := MetadataToJSON(in)
    if err != nil { t.Fatalf("encode err: %v", err) }
    out, err := MetadataFromJSON(b)
    if err != nil { t.Fatalf("decode err: %v", err) }
    if out.Title != in.Title || out.Author != in.Author || len(out.Tags) != 2 { t.Fatalf("mismatch: %#v", out) }
}


