package main

import (

    "os"
    "log"
    "strconv"

    "torme/encode/bencode"

    "torme/server"
    "torme/config"
)



func cmd(){
    path := os.Args[2]

    decoder := &bencode.Decoder{}

    file, _ := os.Open(path)
    defer file.Close()

    tf, err := decoder.Decode(file)

    if err != nil{
        log.Println("decode failed")
    }

    log.Println("decode success")

    log.Println("data:")
    decoder.PrintK()

    log.Println("name is:", tf.Name)
    log.Println("torrent encoding is:", tf.Encoding)

    for _, fi := range tf.Files {
        paths := "["
        var ff =  false
        for _, p := range fi.Path{
            if ff {
                paths += ","
            }

            ff = true
            
            paths += p
        
        } 

        paths += "]"
        log.Printf("path:%s, length:%d\n", paths, fi.Length)
    }

    var ff = false
    anlv := "["
    for _, anl := range tf.Announces{
        if ff{
            anlv += ","
        }
        ff = true
        var f2 = false
        anlv += "["

        for _, anv := range anl.Announce{
            
            if f2 {
                anlv += ","
            }
        
            f2 = true

            anlv += anv
        }

        anlv += "]"
    }

    anlv += "]"

    log.Println("announces list:", anlv)
    log.Println("comment:", tf.Comment)

    log.Println("publiser:", tf.Publisher)
    log.Println("publisher url:", tf.PublisherUrl)

    

    // pl := len(tf.Pieces)

    // log.Println("pieces:")
    // ff := false
    // for i:= 0; i < pl; i++{
    //     if ff{
    //         log.Printf(",")
    //     }
    //     ff = true
    //     str := string(tf.Pieces[i])
    //     log.Printf("\"%s\"", str)
    // }

}


func doServer(){
    port := os.Args[2]
    config.TorrentFilePath = os.Args[3]

    dp, err := strconv.Atoi(port)

    if err != nil{
        log.Fatalln("port is invalid")
    }

    server.Start(dp)   
}


func main(){
    method := os.Args[1]

    if method == "cmd"{
        cmd()
    }

    if method == "server"{
        doServer()
    }

    
}