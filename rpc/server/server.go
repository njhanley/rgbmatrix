package server

import (
	"image"
	"image/color"
	"sync"

	"github.com/njhanley/rgbmatrix"
	"github.com/njhanley/rgbmatrix/rpc/internal/types"
)

type Server struct {
	mu sync.Mutex
	ma *rgbmatrix.Matrix
}

func New(ma *rgbmatrix.Matrix) *Server {
	return &Server{ma: ma}
}

func (s *Server) Size(_ types.None, size *image.Point) error {
	*size = s.ma.Size()
	return nil
}

func (s *Server) SwapOnVSync(types.None, *types.None) error {
	s.mu.Lock()
	s.ma.SwapOnVSync()
	s.mu.Unlock()
	return nil
}

func (s *Server) Clear(types.None, *types.None) error {
	s.mu.Lock()
	s.ma.Clear()
	s.mu.Unlock()
	return nil
}

func (s *Server) Fill(c color.RGBA, _ *types.None) error {
	s.mu.Lock()
	s.ma.Fill(c)
	s.mu.Unlock()
	return nil
}

func (s *Server) Set(p types.Pixel, _ *types.None) error {
	s.mu.Lock()
	s.ma.Set(p.X, p.Y, p.C)
	s.mu.Unlock()
	return nil
}

func (s *Server) Draw(m *image.RGBA, _ *types.None) error {
	s.mu.Lock()
	s.ma.Draw(m)
	s.mu.Unlock()
	return nil
}
