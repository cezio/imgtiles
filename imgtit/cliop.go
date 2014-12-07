package imgtit

import (
	"flag"
)

type Options struct {
	OutputHeight int //h of output image
	OutputWidth  int //w of output image
	TileHeight   int //h of a tile (it shoudl be scaled to match modulo from OutputHeight
	TileWidth    int // w of a tile (should be scaled/matched modulo from OutputWidth)
	InputDir     string //path to a directory with images for tiles
	InputFile    string //path to a file with input file
	OutputFile   string //path to resutl file to be created
    ErrorsFile   string //path to result error image (gray-scale black-to-white map of distances from reference color to resulting color)
    ColorMapFile string //path to analyzed color map (color image)
    OverwriteOutput bool //allow to ovewrite output file, default no
}

func ParseArgs() *Options {
	var opts = new(Options)

	flag.StringVar(&opts.InputDir, "dir", ".", "Directory with source images for tiles")
	flag.StringVar(&opts.InputFile, "in", "", "Input file")
	flag.StringVar(&opts.OutputFile, "out", "output.jpg", "Output file")
	flag.IntVar(&opts.OutputHeight, "height", 600, "Output image height in px")
	flag.IntVar(&opts.OutputWidth, "width", 800, "Output width in px")
	flag.IntVar(&opts.TileHeight, "tile_height", 12, "Height of a tile, in px")
	flag.IntVar(&opts.TileWidth, "tile_width", 16, "Width of a tile, in px")
    flag.StringVar(&opts.ErrorsFile, "error", "", "Error file")
    flag.StringVar(&opts.ColorMapFile, "color_map", "", "Color map of master file")
    flag.BoolVar(&opts.OverwriteOutput, "overwrite", false, "Allow to overwrite output file (default: no)")

	flag.Parse()
	return opts
}
