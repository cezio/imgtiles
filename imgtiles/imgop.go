package imgtiles

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/gographics/gmagick"
)

/*
ColorMode Base type
*/
type ColorMode int

/*
Color modes
*/
const (
	ColorGrayscale ColorMode = iota
	ColorRGB
	ColorRGBA
)

/*
Returns name of ColorMode
*/
func (c ColorMode) String() string {
	return [...]string{"GrayScale", "RGB", "RGBA"}[c]

}

/*
Color structure
*/
type Color struct {
	Mode  ColorMode
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

/*
Position info
*/
type Position struct {
	X int
	Y int
}

/*
Pixel keeps color info
*/
type Pixel struct {
	Color    *Color
	Position Position
}

/*
Pixdiff keeps pixel color difference with information on diff direction
*/
type Pixdiff struct {
	R float64
	G float64
	B float64
	A float64
}

/*
TileMatch keeps information on matching a tile to a given average pixel
*/
type TileMatch struct {
	Matched      bool
	Difference   *Pixdiff
	OrigColor    *Color
	MatchedColor *Color
}

/*
TilePosition keeps position of a Tile within a main image
*/
type TilePosition struct {
	// Frame information
	// x/y offset from 0,0 in px
	Xoffset int
	Yoffset int
	// x/y offset in tiles count
	XTile int
	YTile int
}

/*
TileImg consolidates a tile information
*/
type TileImg struct {
	// processed, shrinked image
	Image *gmagick.MagickWand
	// average color from image
	Color    *Color
	Position *TilePosition
}

/*
ProcessedImage keeps information of processed input image and after tilin
*/
type ProcessedImage struct {
	In *gmagick.MagickWand
	// output
	Out *gmagick.MagickWand
	Px  *[]Pixel

	TileWidth         int
	TileHeight        int
	TilesVertically   int
	TilesHorizontally int
}

/*
Run processing for an image
*/
func Run(opts *Options) (bool, error) {
	gmagick.Initialize()
	defer gmagick.Terminate()

	var checks = runChecks(opts)
	if checks != nil {
		return false, checks
	}

	// get list of files in dir
	var inputDirContents, errd = ioutil.ReadDir(opts.InputDir)
	if errd != nil {
		return false, errd
	}

	// images that will be used as tiles
	var inImages = make([]TileImg, 0)

	// list input dir and prepare image color matrix from available images
	for _, dirItem := range inputDirContents {
		var dirPath = path.Join(opts.InputDir, dirItem.Name())
		// skip non-files
		if !isFile(dirPath) {
			continue
		}
		log.Printf("Processing %v item\n", dirPath)
		var img = analyzeInputFile(dirPath, opts)
		if img == nil {
			continue
		}
		inImages = append(inImages, *img)
	}
	log.Printf("Processed %v input files", len(inImages))

	var inFile, errg = os.Open(opts.InputFile)
	if errg != nil {
		log.Println("Cannot process %v item: %v", opts.InputFile, errg)
		return false, errg
	}

	var master = analyzeMasterInputFile(inFile, opts)
	var outdata = produceOutput(master, &inImages, opts)
	if outdata == nil {
		log.Printf("No output file!!! dying!")
		return false, nil
	}
	writeFile(outdata, opts.OutputFile)

	return true, nil
}
