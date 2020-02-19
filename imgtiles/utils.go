package imgtiles

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gographics/gmagick"
)

/*
   For given list, return index for random element
*/
func getRandomItemIndex(total int) int {
	rsrc := rand.NewSource(int64(total))
	return rand.New(rsrc).Intn(total)
}

/*
   Return True if orig colors are within range of compared +/- tolerance
*/
func matchedPixels(orig *Color, compared *Color, tolerance float64) *TileMatch {
	var matched = ((float64(compared.Red)*(1.0+tolerance)) > float64(orig.Red) &&
		float64(orig.Red) > (float64(compared.Red)*(1.0-tolerance))) &&

		((float64(compared.Green)*(1.0+tolerance)) > float64(orig.Green) &&
			float64(orig.Green) > (float64(compared.Green)*(1.0-tolerance))) &&

		((float64(compared.Blue)*(1.0+tolerance)) >
			float64(orig.Blue) && float64(orig.Blue) >
			(float64(compared.Blue)*(1.0-tolerance)))

	var diff = Pixdiff{float64(orig.Red) - float64(compared.Red),
		float64(orig.Blue) - float64(compared.Blue),
		float64(orig.Green) - float64(compared.Green),
		float64(orig.Alpha) - float64(compared.Alpha),
	}
	var out = TileMatch{matched, &diff, orig, compared}
	return &out
}

func writeFile(indata *gmagick.MagickWand, outpath string) bool {
	werr := indata.WriteImage(outpath)
	if werr != nil {
		panic(fmt.Sprintf("Cannot write image to file %v: %v", werr.Error(), outpath))

	}
	return true

}

/* Return true if provided path is a directory. Any other case or error will return false.
 */
func isDir(path string) bool {
	var f, err = os.Open(path)
	if err != nil {
		return false
	}
	var fi, errf = f.Stat()
	if errf != nil {
		return false
	}
	return fi.IsDir()

}

/* Return true if provided path is a regular file. Any other case or error will return false.
 */
func isFile(path string) bool {
	var f, err = os.Open(path)
	if err != nil {
		return false
	}
	var fi, errf = f.Stat()
	if errf != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

/* Return true if path does not exist
 */
func fileExists(path string) bool {
	var _, err = os.Open(path)
	if err != nil {
		return false
	}
	return true
}

/* Run pre-processing checks on provided options
 */
func runChecks(opts *Options) error {
	// check if InputDir exists and is directory
	if opts.InputDir == "" {
		return fmt.Errorf("Empty InputDir path")
	}
	if !isDir(opts.InputDir) {
		return fmt.Errorf("Path for InputDir %v is not a dir", opts.InputDir)
	}
	// check if InputFile exists
	if !isFile(opts.InputFile) {
		return fmt.Errorf("Path for InputFile %v is not a file", opts.InputFile)
	}
	// check if output fie exists - so we won't overwrite it.
	if !opts.OverwriteOutput && fileExists(opts.OutputFile) {
		return fmt.Errorf("Path for OutputFile %v exists", opts.InputDir)
	}
	return nil

}

func getRatio(w int, h int) float64 {
	return float64(w) / float64(h)
}
