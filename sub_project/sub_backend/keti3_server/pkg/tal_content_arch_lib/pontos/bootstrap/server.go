package bootstrap

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type BeforeServerFunc func() error
type AfterServerFunc func()
type OnShutdownFunc func()

type Server struct {
	server     *http.Server
	beforeFunc []BeforeServerFunc
	afterFunc  []AfterServerFunc
	onShutdown []OnShutdownFunc
	opts       *ServerOptions
	exit       chan os.Signal
	mutex      sync.Mutex
	wg         sync.WaitGroup
}

func NewServer() *Server {
	g := gin.New()
	g.ContextWithFallback = true
	s := &Server{
		server: &http.Server{
			Handler: g,
		},
		opts: &ServerOptions{},
		exit: make(chan os.Signal, 2),
	}
	return s
}

func (s *Server) ListenAndServe() error {
	var err error
	for _, fn := range s.beforeFunc {
		err = fn()
		if err != nil {
			return err
		}
	}

	log.Printf("httpserver start and server %s\n", s.server.Addr)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.Server().ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("httpserver start error: %s\n", err)
		}
	}()

	// default start pprof,if close pprof,please config "pprof_close" in toml ([server])
	go func() {
		s.debugServerProcess()
	}()

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(s.exit, syscall.SIGINT, syscall.SIGTERM)
	s.waitShutdown()

	for _, fn := range s.afterFunc {
		fn()
	}

	return err
}

func (s *Server) waitShutdown() {
	<-s.exit
	log.Println("shutting down http server gracefully...")

	//run onShutdown function in goroutine which had registered
	//todo timeout for prevent waiting registered shutdown function unlimited time
	for _, f := range s.onShutdown {
		s.wg.Add(1)

		f := f
		go func() {
			defer s.wg.Done()
			f()
		}()
	}
	s.wg.Wait()

	//wait one more second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second+s.server.WriteTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown http server error:%s\n", err)
	}
}

func (s *Server) Server() *http.Server {
	return s.server
}

func (s *Server) RegisterShutdown(f OnShutdownFunc) {
	s.mutex.Lock()
	s.onShutdown = append(s.onShutdown, f)
	s.mutex.Unlock()
}

func (s *Server) GinEngine() *gin.Engine {
	return s.server.Handler.(*gin.Engine)
}

func (s *Server) UseMiddleware(middleware ...gin.HandlerFunc) {
	s.GinEngine().Use(middleware...)
}

func (s *Server) AddBeforeServerFunc(fns ...BeforeServerFunc) {
	s.beforeFunc = append(s.beforeFunc, fns...)
}

func (s *Server) AddAfterServerFunc(fns ...AfterServerFunc) {
	s.afterFunc = append(s.afterFunc, fns...)
}

func (s *Server) Init(options *ServerOptions) {
	s.server.Addr = options.Addr
	s.server.ReadTimeout = options.ReadTimeout
	s.server.WriteTimeout = options.WriteTimeout
	s.server.IdleTimeout = options.IdleTimeout
	s.opts = options
}
