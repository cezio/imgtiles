package imgtit

import (
    "os";
    "log";
    "math";
    "image";
    "image/draw";
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

    out := (make([]TileImg,0));

    // for each tile, create a rect, and crop it, then process
    log.Printf("Img size: %vx%v, got tiles: %v in widht, %v in height\n", img.Width(), img.Height(), wx, hx);
    for curr_w := 1; curr_w < wx; curr_w++ {
        for curr_h :=1; curr_h < hx; curr_h++ {

        curr_rect := magick.Rect{(curr_w-1) * w, (curr_h-1)* h, uint(w), uint(h)};
        cropped, cerror := img.Crop(curr_rect);
        if (cerror != nil){
            log.Printf("Cannot process %v rect from img: %v\n", curr_rect, cerror)
            return nil;
        }
        var img_data, t_px = processInputImage(cropped, opts);
        tile_pos := TilePosition{curr_w * opts.TileWidth, curr_h * opts.TileHeight, curr_w, curr_h};
        out = append(out, TileImg{img_data, t_px,&tile_pos });
        }
    }
    return &out;
}


/*

*/
func produceOutput(master *[]TileImg, tiles *[]TileImg, opts *Options) (image.Image){
    var out = image.NewRGBA(image.Rect(0, 0, opts.OutputWidth, opts.OutputHeight));
    // tolerance tresholds
    var tolerances = [8]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
    for tile_r, tile := range *master {
        for _, tolerance := range tolerances {
            // gradually find matching tile
            tcolor := getTileFromColor(tile.Color, tiles, tolerance);
            log.Printf("Get tile for %v (%v) tile, %v tolerance: %v", tile_r, tile.Color, tolerance, tcolor);
            /// aw, crap, no image!
            if (tcolor != nil){
                AddImage(out, tile.Position, tcolor);
                break;
            }

        }
    }
    return out;
}


/*
    AddImage modifies img and places nimg image starting from pos
*/
func AddImage(inimg draw.Image, pos *TilePosition, _pimg *magick.Image) {
    if (pos == nil){
        return;
    }
    var pimg, perr = _pimg.GoImage();
    if (perr != nil){
        log.Printf("Cannot extract GoImage from %v: %v", _pimg, perr);
        return;
    }


    nframe_min := image.Point{pos.Xoffset, pos.Yoffset};
    nframe_max := image.Point{int(inimg.Bounds().Max.X) + nframe_min.X, int(inimg.Bounds().Max.Y) + nframe_min.Y};
    nframe := image.Rectangle{nframe_min, nframe_max};
    draw.Draw(inimg, nframe, pimg, image.ZP, draw.Src);

}


/*
    Find a matching tile for given average color
*/
func getTileFromColor(color *magick.Pixel, tiles *[]TileImg, tolerance float64) (*magick.Image){
    var matched = make([]int, 0);
    // compare color pixel with tolerance
    for idx, tile := range *tiles {
        var _matched = matchedPixels(color, tile.Color, tolerance);
        if (_matched.Matched) {
            matched = append(matched, idx);
        }
    }
    if (len(matched)>0){
        idx := getRandomItemIndex(len(matched));
        var m = matched[idx];
        return (*tiles)[m].Image;
    }
    return nil;

}
