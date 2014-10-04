package imgtit

import (
    "math/rand";
    "github.com/rainycape/magick";
)


/*
    For given list, return index for random element
*/
func getRandomItemIndex(total int) (int){
    rsrc := rand.NewSource(int64(total));
    return Rand{rsrc}.Intn(total);
}



/*
    Return True if orig colors are within range of compared +/- tolerance
*/
func matchedPixels(orig *magick.Pixel, compared *magick.Pixel, tolerance float64) (bool){
    return  ((compared.Red * (1+tolerance)) >  orig.Red > (compared.Red * (1-tolerance))) &&
            ((compared.Green * (1+tolerance)) >  orig.Green > (compared.Green * (1-tolerance))) &&
             ((compared.Blue * (1+tolerance)) >  orig.Blue > (compared.Blue * (1-tolerance)))
}
