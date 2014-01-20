package imgtit

import (
    "flag"
)


type Options struct {
    OutputHeight int
    OutputWidth int
    TileSizeHeight int
    TitleSizeWidth int
    InputDir string
    InputFile string
    OutputFile string
}



func ParseArgs() Options {
    var opts = Options{};

    flag.StringVar(&opts.InputDir, "dir", ".", "Directory with source images");
    flag.StringVar(&opts.InputFile, "in", "", "Input file");
    flag.StringVar(&opts.OutputFile, "out", "output.jpg", "Output file");
    flag.IntVar(&opts.OutputHeight, "height", 600, "Output image height in px");
    flag.IntVar(&opts.OutputWidth, "width", 800, "Output width in px");
    flag.IntVar(&opts.TileSizeHeight, "tile_size_h", 12, "Height of a tile, in px");
    flag.IntVar(&opts.TileSizeHeight, "tile_size_w", 16, "Width of a tile, in px");

    flag.Parse();
    return opts;
}
