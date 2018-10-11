package server



import(
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"


)





func Start(port int) {
	_port := fmt.Sprintf(":%d", port)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/{torrent}/{info}", handleQueryTorrentInfo)
	

	err := http.ListenAndServe(_port, router)

	if err != nil{
		log.Fatalln("listen to server err:", err)
	}


}

