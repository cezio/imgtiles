package imgtit

import (
    "os";
    "log";
    "math";
    "github.com/rainycape/magick";
)

/*  analyze if file is image, process it and return TileImg

*/
func analyzeInputFile(f *os.File, opts *Options) (*TileImg) {
    var img, err = magick.Decode(f);
    if (err != nil){
        log.Println("Cannot process %v: %v", f.Name(), err)
        return nil;
    }
    var processed, color = processInputImage(img, opts);
    if (processed == nil){
        return nil;
    }
    var ret = new (TileImg);
    ret.Image = processed; ret.Color = color;
    return ret;
}

/* Return thumbnail from image and average color
*/
func processInputImage(img *magick.Image, opts *Options) (*magick.Image, *magick.Pixel){
    var out *magick.Image;
    var err error;
    if (img.Width() != opts.TileWidth || img.Height() != opts.TileHeight){
        out, err = img.Scale(opts.TileWidth, opts.TileHeight);
        if (err != nil){
            log.Println("Cannot scale %v: %v", img, err);
            return nil, nil;
        }
    } else {
        out = img;
    }

    var px, errp = out.AverageColor();
    if (errp != nil){
        log.Println("Cannot sample average color %v: %v", img, errp);
        return nil, nil;
    }
    return out, (*magick.Pixel)(px);
}


/*
    Analyze input master file, extract tiles, and calculate average color for each tile
    Return array of parsed tiles
*/
func analyzeMasterInputFile(f *os.File, opts *Options) (*[]TileImg){
    var img, err = magick.Decode(f);
    if (err != nil){
        log.Printf("Cannot process %v: %v\n", f.Name(), err)
        return nil;
    }
    // get tile size
    var w, h = opts.TileWidth, opts.TileHeight;

    // store tiles per width and height, math.Ceil expects float, 
    // but we need int in the end, thus conversions.
    var wx = int(math.Ceil(float64(img.Width())/float64(w)));
    var hx = int(math.Ceil(float64(img.Height())/float64(h)));

    // for each tile, create a rect, and crop it, then process
    log.Printf("Img size: %vx%v, got tiles: %v in widht, %v in height\n", img.Width(), img.Height(), wx, hx);
    return nil;
}
