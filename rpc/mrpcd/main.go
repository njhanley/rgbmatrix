package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"

	"github.com/fogleman/gg"
	"github.com/njhanley/rgbmatrix"
	"github.com/njhanley/rgbmatrix/rpc/server"
	"golang.org/x/sys/unix"
)

func main() {
	addr := ":8080"
	cfg := rgbmatrix.DefaultConfig

	flag.StringVar(&addr, "a", addr, "listen address")
	flag.IntVar(&cfg.Rows, "r", cfg.Rows, "rows on each panel")
	flag.IntVar(&cfg.Columns, "c", cfg.Columns, "columns on each panel")
	flag.IntVar(&cfg.ChainLength, "l", cfg.ChainLength, "length of each chain")
	flag.IntVar(&cfg.Parallel, "p", cfg.Parallel, "parallel chains")
	flag.IntVar(&cfg.Brightness, "b", cfg.Brightness, "brightness percentage")
	flag.Parse()

	ma, err := rgbmatrix.New(cfg)
	if err != nil {
		log.Print("rgbmatrix.New: ", err)
		return
	}
	defer ma.Close()

	p := ma.Size()
	width, height := float64(p.X), float64(p.Y)
	dc := gg.NewContext(p.X, p.Y)
	dc.SetRGBA255(0xff, 0x00, 0x00, 0x40)
	dc.SetLineWidth(1)
	dc.DrawRectangle(0, 0, width, height)
	dc.DrawLine(0, 0, width, height)
	dc.DrawLine(width, 0, 0, height)
	dc.Stroke()
	m := dc.Image()

	err = rpc.Register(server.New(ma))
	if err != nil {
		log.Print("rpc.Register: ", err)
		return
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Print("net.Listen: ", err)
		return
	}
	defer l.Close()

	go func() {
		for {
			ma.Clear()
			ma.Draw(m)
			ma.SwapOnVSync()

			conn, err := l.Accept()
			if err != nil {
				log.Print("l.Accept: ", err)
				continue
			}

			// only serve one client at a time
			rpc.ServeConn(conn)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, unix.SIGINT, unix.SIGTERM)
	<-sc
}
