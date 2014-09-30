package imgtit

import (
    "os";
    "path";
    "io/ioutil";
    "fmt";
    "log";
    "github.com/rainycape/magick";
    "container/list";
)

type TileImg struct {
    // processed, shrinked image
    Image *magick.Image;
    // average color from image
    Color *magick.Pixel;
}

type ProcessedImage struct {
    In *magick.Image;
    // output 
    Out *magick.Image;
    Px *[]magick.Pixel;

    TileWidth int;
    TileHeight int;
    TilesVertically int;
    TilesHorizontally int;
}

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

/* Return true if path does not exist
*/
func fileExists(path string) bool {
    var _, err = os.Open(path);
    if err != nil{
        return false;
    }
    return true;
}

/* Run pre-processing checks on provided options
*/
func runChecks(opts *Options) (error){
    // check if InputDir exists and is directory
    if (opts.InputDir == "") {
        return fmt.Errorf("Empty InputDir path");
    }
    if (!isDir(opts.InputDir)) {
        return fmt.Errorf("Path for InputDir %v is not a dir", opts.InputDir);
    }
    // check if InputFile exists
    if (!isFile(opts.InputFile)) {
        return fmt.Errorf("Path for InputFile %v is not a file", opts.InputFile);
    }
    // check if output fie exists - so we won't overwrite it.
    if (fileExists(opts.OutputFile)){
        return fmt.Errorf("Path for OutputFile %v exists.", opts.InputDir);
    }
    return nil;

}

func Run(opts *Options) (bool, error) {
    var checks = runChecks(opts);
    if (checks != nil){
        return false, checks;
    }

    // analyze source image here, each tile is 1*row + col
    var inImage = list.New();

    

    // get list of files in dir
    var inputDirContents, errd = ioutil.ReadDir(opts.InputDir);
    if (errd != nil){
        return false, errd;
    }

    // images that will be used as tiles
    var inImages = list.New();

    // list input dir and prepare image color matrix from available images
    for dirIdx := range inputDirContents {
        var dirItem = inputDirContents[dirIdx];
        var dirPath = path.Join(opts.InputDir, dirItem.Name())
        // skip non-files
        if (!isFile(dirPath)){
            continue;
        }
        log.Printf("Processing %v item\n", dirPath);
        var f, errf = os.Open(dirPath);
        if (errf != nil){
            log.Println("Cannot process %v item: %v", dirPath, errf);
            continue;
        }
        var img = analyzeInputFile(f, opts);
        if (img == nil){
            continue;
        }
        inImages.PushBack(img);
    }
    log.Printf("Processed %v input files", inImages.Len());

    return true, nil;
};

