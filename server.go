package sherbet

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)


func StartServer(port uint, h http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h,
	}

	go func() {
		log.Println("(HTTPServer) Starting Middelware on port: ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("(HTTPServer) error: %s", err)
		}
		log.Println("(HTTPServer) Stoped Middelware on port: ", port)
	}()

	return srv
}


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not protected!\n")
}



