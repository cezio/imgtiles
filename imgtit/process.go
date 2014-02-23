package imgtit

import (
    "os";
    "log";
    "image";
    "image/color";
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
func processInputImage(img *image.Image, opts *Options) (*image.Image, *color.Color){
    var out, err = img.Scale(opts.TileWidth, opts.TileHeight);
    if (err != nil){
        log.Println("Cannot scale %v: %v", img, err);
        return nil, nil;
    }
    var px, errp = out.AverageColor();
    if (errp != nil){
        log.Println("Cannot sample average color %v: %v", img, errp); 
        return nil, nil;
    }
    return out, px;
}

