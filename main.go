package main

import "./imgtit"


func main(){
    var opts = imgtit.ParseArgs();
    imgtit.Run(opts);
}
