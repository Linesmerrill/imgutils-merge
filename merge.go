// Package merge provides image merging and compositing utilities.
package merge

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

// Direction specifies the merge direction.
type Direction int

const (
	// Horizontal merges images side by side.
	Horizontal Direction = iota
	// Vertical merges images top to bottom.
	Vertical
)

// Alignment specifies how to align images when merging.
type Alignment int

const (
	// AlignStart aligns to top (vertical) or left (horizontal).
	AlignStart Alignment = iota
	// AlignCenter centers the images.
	AlignCenter
	// AlignEnd aligns to bottom (vertical) or right (horizontal).
	AlignEnd
)

// Options configures the merge operation.
type Options struct {
	Direction  Direction
	Alignment  Alignment
	Gap        int         // Space between images
	Background color.Color // Background color for gaps and alignment padding
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Direction:  Horizontal,
		Alignment:  AlignCenter,
		Gap:        0,
		Background: color.Transparent,
	}
}

// Merge combines multiple images into one.
func Merge(images []image.Image, opts Options) image.Image {
	if len(images) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	if len(images) == 1 {
		bounds := images[0].Bounds()
		dst := image.NewRGBA(bounds)
		draw.Draw(dst, bounds, images[0], bounds.Min, draw.Src)
		return dst
	}

	// Calculate total dimensions
	var totalW, totalH, maxW, maxH int
	for _, img := range images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		if opts.Direction == Horizontal {
			totalW += w
			if h > maxH {
				maxH = h
			}
		} else {
			totalH += h
			if w > maxW {
				maxW = w
			}
		}
	}

	// Add gaps
	gapTotal := opts.Gap * (len(images) - 1)
	if opts.Direction == Horizontal {
		totalW += gapTotal
		totalH = maxH
	} else {
		totalH += gapTotal
		totalW = maxW
	}

	dst := image.NewRGBA(image.Rect(0, 0, totalW, totalH))

	// Fill background
	if opts.Background != nil {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{opts.Background}, image.Point{}, draw.Src)
	}

	// Draw images
	offset := 0
	for _, img := range images {
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()

		var x, y int
		if opts.Direction == Horizontal {
			x = offset
			switch opts.Alignment {
			case AlignStart:
				y = 0
			case AlignCenter:
				y = (maxH - h) / 2
			case AlignEnd:
				y = maxH - h
			}
			offset += w + opts.Gap
		} else {
			y = offset
			switch opts.Alignment {
			case AlignStart:
				x = 0
			case AlignCenter:
				x = (maxW - w) / 2
			case AlignEnd:
				x = maxW - w
			}
			offset += h + opts.Gap
		}

		rect := image.Rect(x, y, x+w, y+h)
		draw.Draw(dst, rect, img, b.Min, draw.Over)
	}

	return dst
}

// Overlay places one image on top of another at the specified position.
func Overlay(base, overlay image.Image, x, y int, opacity float64) image.Image {
	baseBounds := base.Bounds()
	dst := image.NewRGBA(baseBounds)
	draw.Draw(dst, baseBounds, base, baseBounds.Min, draw.Src)

	overlayBounds := overlay.Bounds()

	for oy := 0; oy < overlayBounds.Dy(); oy++ {
		for ox := 0; ox < overlayBounds.Dx(); ox++ {
			dx := x + ox
			dy := y + oy
			if dx < 0 || dx >= baseBounds.Dx() || dy < 0 || dy >= baseBounds.Dy() {
				continue
			}

			baseColor := dst.At(dx, dy)
			overlayColor := overlay.At(overlayBounds.Min.X+ox, overlayBounds.Min.Y+oy)

			blended := blendWithOpacity(baseColor, overlayColor, opacity)
			dst.Set(dx, dy, blended)
		}
	}

	return dst
}

// blendWithOpacity blends two colors with the given opacity.
func blendWithOpacity(base, overlay color.Color, opacity float64) color.Color {
	br, bg, bb, ba := base.RGBA()
	or, og, ob, oa := overlay.RGBA()

	if oa == 0 {
		return base
	}

	alpha := float64(oa) / 65535.0 * opacity

	r := uint8((float64(br>>8)*(1-alpha) + float64(or>>8)*alpha))
	g := uint8((float64(bg>>8)*(1-alpha) + float64(og>>8)*alpha))
	b := uint8((float64(bb>>8)*(1-alpha) + float64(ob>>8)*alpha))

	return color.RGBA{r, g, b, uint8(ba >> 8)}
}

// Grid arranges images in a grid pattern.
func Grid(images []image.Image, cols int, gap int, bg color.Color) image.Image {
	if len(images) == 0 || cols <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	rows := (len(images) + cols - 1) / cols

	// Find max cell dimensions
	var maxW, maxH int
	for _, img := range images {
		b := img.Bounds()
		if b.Dx() > maxW {
			maxW = b.Dx()
		}
		if b.Dy() > maxH {
			maxH = b.Dy()
		}
	}

	totalW := cols*maxW + (cols-1)*gap
	totalH := rows*maxH + (rows-1)*gap

	dst := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	if bg != nil {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	}

	for i, img := range images {
		row := i / cols
		col := i % cols

		x := col * (maxW + gap)
		y := row * (maxH + gap)

		b := img.Bounds()
		// Center in cell
		x += (maxW - b.Dx()) / 2
		y += (maxH - b.Dy()) / 2

		rect := image.Rect(x, y, x+b.Dx(), y+b.Dy())
		draw.Draw(dst, rect, img, b.Min, draw.Over)
	}

	return dst
}

// MergeFromFiles loads images from files and merges them.
func MergeFromFiles(paths []string, opts Options) (image.Image, error) {
	images := make([]image.Image, 0, len(paths))

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return Merge(images, opts), nil
}

// SaveJPEG saves the merged image as JPEG.
func SaveJPEG(img image.Image, w io.Writer, quality int) error {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// SavePNG saves the merged image as PNG.
func SavePNG(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}
