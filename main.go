package main

import (
    "./imgtit";
    "fmt";
)


func main(){
    var opts = imgtit.ParseArgs();
    var _, err = imgtit.Run(opts);
    if (err != nil){
        fmt.Println("Error while processing", opts, err);
    }
}
