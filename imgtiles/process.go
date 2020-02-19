package imgtiles

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/gographics/gmagick"
)

/*
GmagickNoError - magic const for gmagick.GetLastError(), which returns "undefined exception" on non-error returns
*/
const GmagickNoError = "undefined exception"

/*  analyze if file is image, process it and return TileImg

 */
func analyzeInputFile(fpath string, opts *Options) *TileImg {
	mw := gmagick.NewMagickWand()
	mw.ReadImage(fpath)
	lastErr := mw.GetLastError()

	if lastErr.Error() != GmagickNoError {
		panic(fmt.Sprintf("Cannot process %v: %v", fpath, lastErr.Error()))
	}
	var processed, color = processInputImage(mw, opts)

	if processed == nil {
		return nil
	}
	var ret = new(TileImg)
	ret.Image = processed
	ret.Color = color
	return ret
}

/* Return thumbnail from image and an average color

   img pointer to input image
   opts Options with configuration

   If img doesn't have the same size as a requested Tile, it will be resized, and resized copy will be returned

   retruns
   tile pointer to image
   average pixel pointer to average pixel
*/
func processInputImage(img *gmagick.MagickWand, opts *Options) (*gmagick.MagickWand, *Color) {
	log.Printf("Processing %v, h: %v, h:%v", img.GetImageFilename(), img.GetImageHeight(), img.GetImageWidth())
	if int(img.GetImageWidth()) != opts.TileWidth || int(img.GetImageHeight()) != opts.TileHeight {
		log.Printf("Img size doesn't match h:%v, w:%v", opts.TileHeight, opts.TileWidth)
		scaled := img.Clone()
		scaled.ScaleImage(uint(opts.TileWidth), uint(opts.TileHeight))
		log.Print("Scaled")

		err := scaled.GetLastError()
		if err.Error() != GmagickNoError {
			log.Println("Cannot scale %v: %v", img, err.Error())
			return nil, nil
		}
		var px = getAverageColor(scaled)
		return scaled, px
	}
	var px = getAverageColor(img)
	return img, px
}

/*
   Analyze input master file, extract tiles, and calculate average color for each tile
   Return array of parsed tiles
*/
func analyzeMasterInputFile(f *os.File, opts *Options) *[]TileImg {
	mw := gmagick.NewMagickWand()
	mw.ReadImageFile(f)
	readErr := mw.GetLastError()
	if readErr.Error() != GmagickNoError {
		log.Printf("Cannot process %v: %v\n", f.Name(), readErr.Error())
		return nil
	}

	// get tile size
	var w, h = opts.TileWidth, opts.TileHeight

	// store tiles per width and height, math.Ceil expects float,
	// but we need int in the end, thus conversions.
	var wx = int(math.Ceil(float64(mw.GetImageWidth()) / float64(w)))
	var hx = int(math.Ceil(float64(mw.GetImageHeight()) / float64(h)))

	out := (make([]TileImg, 0))

	// for each tile, create a rect, and crop it, then process
	log.Printf("Img size: %vx%v, got tiles: %v in width, %v in height\n", mw.GetImageWidth(), mw.GetImageHeight(), wx, hx)
	log.Printf("Tiles size x:%v, y:%v", opts.TileWidth, opts.TileHeight)

	cnt := 1
	for curr_w := 0; curr_w < wx; curr_w++ {
		for curr_h := 0; curr_h < hx; curr_h++ {
			log.Printf("Tile %v from master", cnt)
			cnt++
			cropped := mw.Clone()
			cerror := cropped.CropImage(uint(w), uint(h), curr_w*w, curr_h*h)
			if cerror != nil && cerror.Error() != GmagickNoError {
				panic(fmt.Sprintf("Cannot process x: %v/ y: %v (h: %v, w: %v) rect from img: %v\n", curr_w, curr_h, w, h, cerror.Error()))
			}
			var img_data, t_px = processInputImage(cropped, opts)
			tile_pos := TilePosition{curr_w * opts.TileWidth, curr_h * opts.TileHeight, curr_w, curr_h}
			out = append(out, TileImg{img_data, t_px, &tile_pos})
		}
	}
	return &out
}

/*

 */
func produceOutput(master *[]TileImg, tiles *[]TileImg, opts *Options) *gmagick.MagickWand {
	out := gmagick.NewMagickWand()
	out.SetSize(uint(opts.OutputWidth), uint(opts.OutputHeight))
	// tolerance tresholds
	var tolerances = [8]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
	for tile_r, tile := range *master {
		for _, tolerance := range tolerances {
			// gradually find matching tile
			tcolor := getTileFromColor(tile.Color, tiles, tolerance)
			log.Printf("Get tile for %v (%v) tile, %v tolerance: %v", tile_r, tile.Color, tolerance, tcolor)
			/// aw, crap, no image!
			if tcolor != nil {
				AddImage(out, tile.Position, tcolor)
				break
			}

		}
	}
	return out
}

/*
   AddImage modifies img and places _pimg image starting from pos
*/
func AddImage(inimg *gmagick.MagickWand, pos *TilePosition, _pimg *gmagick.MagickWand) {
	err := inimg.CompositeImage(_pimg, gmagick.COMPOSITE_OP_OVER, pos.Xoffset, pos.Yoffset)
	if err != nil {
		panic(fmt.Sprintf("Couldn't overlay image: %v", err.Error()))
	}
}

/*
   Find a matching tile for given average color
*/
func getTileFromColor(color *Color, tiles *[]TileImg, tolerance float64) *gmagick.MagickWand {
	var matched = make([]int, 0)
	// compare color pixel with tolerance
	for idx, tile := range *tiles {
		var _matched = matchedPixels(color, tile.Color, tolerance)
		if _matched.Matched {
			matched = append(matched, idx)
		}
	}
	if len(matched) > 0 {
		idx := getRandomItemIndex(len(matched))
		var m = matched[idx]
		return (*tiles)[m].Image
	}
	return nil

}

func getAverageColor(inImg *gmagick.MagickWand) *Color {
	pxCount, pxList := inImg.GetImageHistogram()
	log.Printf("Img histogram: %v colors", pxCount)
	var rVal, gVal, bVal float64
	var rCount float64

	for pxIdx := 0; pxIdx < int(pxCount); pxIdx++ {

		pxItem := pxList[pxIdx]
		rVal += pxItem.GetRed()
		gVal += pxItem.GetGreen()
		bVal += pxItem.GetRed()
		rCount += float64(pxItem.GetColorCount())
	}

	pix := Color{ColorRGB,
		rVal / rCount,
		gVal / rCount,
		bVal / rCount,
		0.0}
	log.Printf("Average Color: %v", pix)
	return &pix
}
