package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type myServer struct {
	http.Server
	shutdownReq chan bool
}

// NewServer - this is the init function for the server process
func NewServer(port string) *myServer {

	//create server
	s := &myServer{
		Server: http.Server{
			Addr:         ":" + port,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	//register handlers
	router.HandleFunc("/api/simpleserver", s.RootHandler)
	router.HandleFunc("/nwadmin/simpleserver", s.RootHandlerAdmin)
	// CORS stuff
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "X-Request-Token", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	s.Handler = handlers.CORS(headersOk, originsOk, methodsOk)(router)

	return s
}

func (s *myServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}
	log.Printf("Stopping API server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}

}

func (s *myServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("simpleserver\n"))
}

func (s *myServer) RootHandlerAdmin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("simpleserver\n"))
}

func GetTokenFromRequest(r *http.Request) string {
	var tmptoken string
	tmptoken = r.Header.Get("X-API-KEY")
	if tmptoken != "" {
		return tmptoken
	}
	tmptoken = r.URL.Query().Get("authtoken")
	if tmptoken != "" {
		return tmptoken
	}

	return tmptoken
}
