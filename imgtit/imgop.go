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

func isDir(path string) bool{

}

func isFile(path string) bool {

}

func Run(opts Options) (bool, error) {
    // check if InputDir exists and is directory

};
