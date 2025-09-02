package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// defaultPprofPort default pprof port
const defaultPprofPort = ":15555"

// debugServerProcess start debug golang code
func (s *Server) debugServerProcess() {
	if !s.opts.PprofClose {
		go func() {
			serverPort := defaultPprofPort
			if s.opts.PprofPort != "" {
				serverPort = s.opts.PprofPort
			}

			log.Printf("[HTTP] pprof server listening at port %s\n", serverPort)
			if err := http.ListenAndServe(fmt.Sprintf("%s", serverPort), nil); err != nil {
				log.Fatalln("[HTTP] pprof server start failed", "error", err)
			}
		}()
	}
}
