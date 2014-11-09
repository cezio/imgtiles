package imgtit

import (
    "os";
    "log";
    "bufio";
    "image";
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
    return  ((compared.Red * (1+tolerance)) >  orig.Red > (compared.Red * (1-tolerance))) &&
            ((compared.Green * (1+tolerance)) >  orig.Green > (compared.Green * (1-tolerance))) &&
             ((compared.Blue * (1+tolerance)) >  orig.Blue > (compared.Blue * (1-tolerance)))
}

func writeFile(indata image.Image, outpath string) (bool){
    var f, ferr = os.OpenFile(outpath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664);
    if (ferr != nil){
        log.Printf("Cannot open %v for writing: %v", outpath, ferr);
        return false;
    }
    var bwriter = bufio.NewWriter(f);

    var info = image.NewInfo();
    info.SetFormat("jpg");
    info.SetQuality(100);
    info.SetColorspace(image.RGB);

    var ierr = indata.Encode(bwriter, info);
    if (ierr != nil){
        log.Printf("Cannot write to %v file: %v", outpath, ierr);
    }
}
