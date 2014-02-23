package imgtit

import (
    "os";
    "image";
    "image/color";
)

type TilesImg struct {
    Images []image.Image;
    Colors []color.Color;
}

type DestinationImgParams struct {

    // destination image
    Img image.Image
    Src image.Image

    // title size
    TitleSizeW int
    TitleSizeH int
    TitleCount int

};

/* Return true if provided path is a directory. Any other case or error will return false.
*/
func isDir(path string) bool{
    var f, err = os.Open(path);
    if err != nil{
        return false;
    }
    var fi,errf = f.Stat();
    if errf != nil{
        return false;
    }
    return fi.IsDir();

}

/* Return true if provided path is a regular file. Any other case or error will return false.
*/
func isFile(path string) bool {
    var f, err = os.Open(path);
    if err != nil{
        return false;
    }
    var fi, errf = f.Stat();
    if errf != nil{
        return false;
    }
    return fi.Mode().IsRegular();
}

func Run(opts Options) (bool, error) {
    // check if InputDir exists and is directory
    return true, nil;
};
