package cartio

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
)

// RPIHeader represents the header of an .rpi file
type RPIHeader struct {
	Width   uint16
	Height  uint16
	Flags   uint16 // bit 0: 0=landscape, 1=portrait
	Reserved uint16
}

// LoadRPI loads a .rpi (Raw Palette Indexed) file and converts it to a SpriteData.
// The file can be gzip-compressed.
func LoadRPI(data []byte) (*SpriteData, error) {
	var reader io.Reader
	
	// Try to decompress as gzip first
	br := &bytesReader{data: data}
	gzReader, err := gzip.NewReader(br)
	if err == nil {
		// It's gzip-compressed
		defer gzReader.Close()
		reader = gzReader
	} else {
		// Not compressed, use raw data
		br.offset = 0
		reader = br
	}
	
	// Read header (8 bytes)
	var header RPIHeader
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("failed to read RPI header: %w", err)
	}
	
	width := int(header.Width)
	height := int(header.Height)
	
	// Read packed pixel data
	// 4 pixels = 3 bytes (4 * 6 bits = 24 bits = 3 bytes)
	totalPixels := width * height
	pixelsNeeded := (totalPixels + 3) / 4 * 3 // Round up to nearest 3-byte group
	packedData := make([]byte, pixelsNeeded)
	
	n, err := io.ReadFull(reader, packedData)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("failed to read RPI pixel data: %w", err)
	}
	packedData = packedData[:n]
	
	// Unpack 6-bit values from bytes
	// Each 3 bytes contains 4 pixels (24 bits total)
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
	}
	
	pixelIdx := 0
	for i := 0; i < len(packedData); i += 3 {
		if i+2 >= len(packedData) {
			break
		}
		
		// Extract 4 pixels from 3 bytes
		b0 := packedData[i]
		b1 := packedData[i+1]
		b2 := packedData[i+2]
		
		// Unpack:
		// Pixel 0: Byte 0 bits 0-5
		// Pixel 1: Byte 0 bits 6-7, Byte 1 bits 0-3
		// Pixel 2: Byte 1 bits 4-7, Byte 2 bits 0-1
		// Pixel 3: Byte 2 bits 2-7
		p0 := int(b0 & 0x3F)
		p1 := int((b0>>6)&0x03) | int((b1&0x0F)<<2)
		p2 := int((b1>>4)&0x0F) | int((b2&0x03)<<4)
		p3 := int((b2 >> 2) & 0x3F)
		
		// Map encoded values back: 63 = transparent (-1), 0-49 = palette indices
		unpacked := [4]int{p0, p1, p2, p3}
		for j := 0; j < 4 && pixelIdx < totalPixels; j++ {
			encoded := unpacked[j]
			var pixel int
			if encoded == 63 {
				pixel = -1 // Transparent
			} else if encoded >= 0 && encoded <= 49 {
				pixel = encoded // Palette index
			} else {
				pixel = -1 // Invalid, default to transparent
			}
			
			y := pixelIdx / width
			x := pixelIdx % width
			if y < height && x < width {
				pixels[y][x] = pixel
			}
			pixelIdx++
		}
	}
	
	return &SpriteData{
		Width:        width,
		Height:       height,
		Pixels:       pixels,
		Type:         SpriteTypeStatic, // RPI files are always static
		UseCollision: false,
		IsUI:         true, // Screens/backgrounds are typically UI
		Lifetime:     0,
		MaxSpawn:     0,
		MountPoints:  nil,
	}, nil
}

// bytesReader wraps []byte to implement io.Reader
type bytesReader struct {
	data   []byte
	offset int
}

func (br *bytesReader) Read(p []byte) (n int, err error) {
	if br.offset >= len(br.data) {
		return 0, io.EOF
	}
	n = copy(p, br.data[br.offset:])
	br.offset += n
	return n, nil
}

