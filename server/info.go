package server

import (
    "os"
    "io"
	"fmt"
    "net/http"
    "log"

    "github.com/gorilla/mux"


    "torme/encode/bencode"
    "torme/torrent"

    "torme/config"

)


func handleQueryTorrentInfo(w http.ResponseWriter, r *http.Request){
    log.Println("handle query torrent info")

	vars := mux.Vars(r)
    tr := vars["torrent"]
    info := vars["info"]


    

    path := fmt.Sprintf("%s/%s", config.TorrentFilePath, tr)

    if _, err := os.Stat(path); err != nil{
        if os.IsNotExist(err){
            io.WriteString(w, "no this torrent")
        }
    }


    file, _ := os.Open(path)
    defer file.Close()

    decoder := &bencode.Decoder{}
    tf, err := decoder.Decode(file)

    if err != nil{
        log.Println("decode failed")
    }

    log.Println("decode success")

    log.Println("data:")
    decoder.PrintK()


    switch(info){
    case "base":
        queryBase(tf, w, r)
    case "files":
        queryFiles(tf, w, r)
    case "announce":
        queryAnnounces(tf, w, r)
    case "pieces":
        queryPieces(tf, w, r)
    }
	
}

func queryBase(tf *torrent.TorrentInfo, w http.ResponseWriter, r *http.Request) {

    str := "<h1>name:" + tf.Name + "</h1>"
    str += "<h1> torrent encoding is:" + tf.Encoding + "</h1>"
    str += "<h1>comment:" + tf.Comment + "</h1>"
    str += "<h1>publisher:" + tf.Publisher + "</h1>"
    str += "<h1>publisher url:" + tf.PublisherUrl + "</h1>"

    io.WriteString(w, str)
}

func queryAnnounces(tf *torrent.TorrentInfo, w http.ResponseWriter, r *http.Request) {
    str := "<h1>announce list</h1>"
    str += "<table>"
    for _, anl := range tf.Announces{
        str += "<tr>"
    
        for _, anv := range anl.Announce{
           str += "<td>" + anv + "</td>"
        }

        str += "</tr>"
    }

    str += "</table>"

    io.WriteString(w, str)
}


func queryFiles(tf *torrent.TorrentInfo, w http.ResponseWriter, r *http.Request){

    str := "<h1>files</h1>"
    str += "<table>"
    for _, fi := range tf.Files {
        str += "<tr>"
        
        for _, p := range fi.Path{
            str += "<td>" + p + "</td>"
        } 
        str += fmt.Sprintf("<td>%d</td></tr>",fi.Length)
    }

    str += "</table>"
    io.WriteString(w, str)
}

func queryPieces(tf *torrent.TorrentInfo, w http.ResponseWriter, r *http.Request){
    str := "<h1>pieces</h1>"
    str += "<table>"

    pl := len(tf.Pieces)
    for i:= 0; i < pl; i++{
        str += "<tr><td>" + string(tf.Pieces[i]) + "</td></tr>"
    }
    str += "</table>"

    io.WriteString(w, str)
}
