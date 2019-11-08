package server

import (
	"image"
	"sync"

	"github.com/njhanley/rgbmatrix"
)

type Server struct {
	mu sync.Mutex
	ma *rgbmatrix.Matrix
}

func New(ma *rgbmatrix.Matrix) *Server {
	return &Server{ma: ma}
}

func (s *Server) Bounds(_ struct{}, r *image.Rectangle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	*r = s.ma.Rect
	return nil
}

func (s *Server) Swap(m *image.RGBA, _ *struct{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ma.RGBA = m
	s.ma.Swap()
	return nil
}
