package imgtit

import (
    "os";
    "log";
    "bufio";
    "image";
    "image/jpeg";
    "math/rand";
    "github.com/rainycape/magick";
)


/*
    For given list, return index for random element
*/
func getRandomItemIndex(total int) (int){
    rsrc := rand.NewSource(int64(total));
    return rand.New(rsrc).Intn(total);
}



/*
    Return True if orig colors are within range of compared +/- tolerance
*/
func matchedPixels(orig *magick.Pixel, compared *magick.Pixel, tolerance float64) (bool){
    return  ((float64(compared.Red) * (1.0+tolerance)) > float64(orig.Red) &&
             float64(orig.Red) > (float64(compared.Red) * (1.0-tolerance))) &&

            ((float64(compared.Green) * (1.0+tolerance)) > float64(orig.Green) &&
             float64(orig.Green) > (float64(compared.Green) * (1.0-tolerance))) &&

             ((float64(compared.Blue) * (1.0+tolerance)) >
                                               float64(orig.Blue) && float64(orig.Blue) >
               (float64(compared.Blue) * (1.0-tolerance)))
}

func writeFile(indata image.Image, outpath string) (bool){
    var f, ferr = os.OpenFile(outpath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664);
    if (ferr != nil){
        log.Printf("Cannot open %v for writing: %v", outpath, ferr);
        return false;
    }
    var bwriter = bufio.NewWriter(f);
    // use correct type

    var opts = new(jpeg.Options)
    opts.Quality = 100;

    var ierr = jpeg.Encode(bwriter, indata, opts);
    if (ierr != nil){
        log.Printf("Cannot write to %v file: %v", outpath, ierr);
    }
    return true;
}
