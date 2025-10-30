package cart

import "testing"

func TestMetadataValidate(t *testing.T) {
    m := Metadata{Title: "Game", Tags: []string{"a","b","c","d","e"}}
    if err := m.Validate(); err != nil {
        t.Fatalf("unexpected: %v", err)
    }
    m.Tags = append(m.Tags, "f")
    if err := m.Validate(); err == nil {
        t.Fatalf("expected tag limit error")
    }
}


