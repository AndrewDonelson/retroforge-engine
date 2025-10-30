package cart

import "errors"

// Metadata represents basic information embedded in a cart.
type Metadata struct {
    Title       string
    Author      string
    Version     string
    Description string
    Genre       string
    Tags        []string // up to 5
}

func (m Metadata) Validate() error {
    if len(m.Tags) > 5 {
        return errors.New("too many tags (max 5)")
    }
    return nil
}


