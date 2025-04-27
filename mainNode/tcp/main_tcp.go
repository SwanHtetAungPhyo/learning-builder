package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SwanHtetAungPhyo/learning/common"
)

const (
	listenAddr    = "127.0.0.1:8081"
	shutdownDelay = 5 * time.Second
)

type Server struct {
	listener       net.Listener
	wg             sync.WaitGroup
	done           chan struct{}
	acceptedBlocks chan *common.Block
	connections    map[net.Conn]struct{}
	connMutex      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		done:           make(chan struct{}),
		acceptedBlocks: make(chan *common.Block, 100),
		connections:    make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	log.Printf("Server listening on %s", listenAddr)
	s.wg.Add(1)
	go s.processBlocks()
	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.done:
					return
				default:
					log.Printf("Accept error: %v", err)
					continue
				}
			}

			s.connMutex.Lock()
			s.connections[conn] = struct{}{}
			s.connMutex.Unlock()

			s.wg.Add(1)
			go func(c net.Conn) {
				defer s.wg.Done()
				s.handleClient(c)

				s.connMutex.Lock()
				delete(s.connections, c)
				s.connMutex.Unlock()
			}(conn)
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	decoder := json.NewDecoder(conn)
	for {
		select {
		case <-s.done:
			return
		default:
			var block common.Block
			if err := decoder.Decode(&block); err != nil {
				if err.Error() == "EOF" {
					log.Println("Client disconnected normally")
					return
				}

				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				if n > 0 {
					log.Printf("Unexpected data received: %q", string(buf[:n]))
				}
				log.Printf("Decode error: %v", err)
				return
			}

			select {
			case s.acceptedBlocks <- &block:
			case <-s.done:
				return
			}
		}
	}
}

func (s *Server) processBlocks() {
	defer s.wg.Done()

	for {
		select {
		case block := <-s.acceptedBlocks:
			if !block.VerifyBlockByMerkle() {
				log.Printf("Unverified block: %s", block.Hash)
				return
			}
			log.Printf("✅ Accepted block: %s with %d TXs", block.Hash, len(block.Txs))
		case <-s.done:
			for {
				select {
				case block := <-s.acceptedBlocks:
					log.Printf("Processing final block: %s", block.Hash)
				default:
					return
				}
			}
		}
	}
}

func (s *Server) Shutdown() {
	close(s.done)

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}

	s.connMutex.Lock()
	for conn := range s.connections {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}
	s.connMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownDelay)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server shutdown gracefully")
	case <-ctx.Done():
		log.Println("Server shutdown timed out, some connections may not have closed cleanly")
	}
}

func main() {
	server := NewServer()

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	log.Printf("Received signal %v, shutting down...", sig)

	server.Shutdown()
	log.Println("Server shutdown completed")
}
