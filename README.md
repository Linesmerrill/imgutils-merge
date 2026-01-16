# imgutils-merge

[![Go Reference](https://pkg.go.dev/badge/github.com/imgutils-org/imgutils-merge.svg)](https://pkg.go.dev/github.com/imgutils-org/imgutils-merge)
[![Go Report Card](https://goreportcard.com/badge/github.com/imgutils-org/imgutils-merge)](https://goreportcard.com/report/github.com/imgutils-org/imgutils-merge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for merging and compositing images. Part of the [imgutils](https://github.com/imgutils-org) collection.

## Features

- Horizontal and vertical image merging
- Alignment options (start, center, end)
- Gap/spacing between images
- Grid layout for multiple images
- Overlay compositing with opacity
- Configurable background colors

## Installation

```bash
go get github.com/imgutils-org/imgutils-merge
```

## Quick Start

```go
package main

import (
    "image"
    "os"

    "github.com/imgutils-org/imgutils-merge"
)

func main() {
    // Load images
    var images []image.Image
    for _, path := range []string{"img1.jpg", "img2.jpg", "img3.jpg"} {
        file, _ := os.Open(path)
        img, _, _ := image.Decode(file)
        file.Close()
        images = append(images, img)
    }

    // Merge horizontally
    result := merge.Merge(images, merge.DefaultOptions())

    // Save result
    out, _ := os.Create("merged.jpg")
    merge.SaveJPEG(result, out, 85)
    out.Close()
}
```

## Usage Examples

### Horizontal Merge

```go
opts := merge.Options{
    Direction: merge.Horizontal,
    Alignment: merge.AlignCenter,
}
result := merge.Merge(images, opts)
```

### Vertical Merge

```go
opts := merge.Options{
    Direction: merge.Vertical,
    Alignment: merge.AlignCenter,
}
result := merge.Merge(images, opts)
```

### With Gap/Spacing

```go
opts := merge.Options{
    Direction:  merge.Horizontal,
    Gap:        20, // 20px between images
    Background: color.White,
}
result := merge.Merge(images, opts)
```

### Alignment Options

```go
// Align to top (vertical merge) or left (horizontal merge)
opts := merge.Options{
    Direction: merge.Horizontal,
    Alignment: merge.AlignStart,
}

// Center alignment (default)
opts := merge.Options{
    Direction: merge.Horizontal,
    Alignment: merge.AlignCenter,
}

// Align to bottom (vertical merge) or right (horizontal merge)
opts := merge.Options{
    Direction: merge.Horizontal,
    Alignment: merge.AlignEnd,
}
```

### Grid Layout

```go
// Create a 3-column grid
result := merge.Grid(images, 3, 10, color.White)
// 3 = columns
// 10 = gap between cells
// color.White = background color
```

### Overlay Images

```go
// Place one image on top of another
result := merge.Overlay(base, overlay, 100, 50, 0.8)
// 100, 50 = x, y position
// 0.8 = 80% opacity
```

### From File Paths

```go
paths := []string{"photo1.jpg", "photo2.jpg", "photo3.jpg"}
result, err := merge.MergeFromFiles(paths, merge.Options{
    Direction: merge.Horizontal,
    Gap:       10,
})
if err != nil {
    log.Fatal(err)
}
```

## API Reference

### Types

#### Direction

```go
type Direction int

const (
    Horizontal Direction = iota // Side by side
    Vertical                    // Top to bottom
)
```

#### Alignment

```go
type Alignment int

const (
    AlignStart  Alignment = iota // Top or Left
    AlignCenter                  // Center
    AlignEnd                     // Bottom or Right
)
```

#### Options

```go
type Options struct {
    Direction  Direction   // Merge direction
    Alignment  Alignment   // Image alignment
    Gap        int         // Space between images
    Background color.Color // Background color
}
```

### Functions

| Function | Description |
|----------|-------------|
| `DefaultOptions()` | Returns defaults (Horizontal, Center, no gap) |
| `Merge(images, opts)` | Merge multiple images |
| `MergeFromFiles(paths, opts)` | Load and merge from file paths |
| `Overlay(base, overlay, x, y, opacity)` | Overlay image at position |
| `Grid(images, cols, gap, bg)` | Arrange images in a grid |
| `SaveJPEG(img, w, quality)` | Save as JPEG |
| `SavePNG(img, w)` | Save as PNG |

## Common Use Cases

### Photo Collage

```go
// Create a 2x2 photo grid
collage := merge.Grid(photos, 2, 5, color.White)
```

### Before/After Comparison

```go
opts := merge.Options{
    Direction: merge.Horizontal,
    Gap:       2,
    Alignment: merge.AlignCenter,
}
comparison := merge.Merge([]image.Image{before, after}, opts)
```

### Social Media Banner

```go
// Merge profile pic with background
banner := merge.Overlay(background, profilePic, 50, 50, 1.0)
```

### Contact Sheet

```go
// Create thumbnail contact sheet
sheet := merge.Grid(thumbnails, 4, 10, color.Black)
```

## Requirements

- Go 1.16 or later

## Related Packages

- [imgutils-watermark](https://github.com/imgutils-org/imgutils-watermark) - Watermarking
- [imgutils-crop](https://github.com/imgutils-org/imgutils-crop) - Image cropping
- [imgutils-sdk](https://github.com/imgutils-org/imgutils-sdk) - Unified SDK

## License

MIT License - see [LICENSE](LICENSE) for details.
