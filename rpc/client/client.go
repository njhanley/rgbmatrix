package client

import (
	"image"
	"net/rpc"
)

type Matrix struct {
	*image.RGBA
	rpc *rpc.Client
}

func Connect(addr string) (*Matrix, error) {
	rpc, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	var r image.Rectangle
	err = rpc.Call("Service.Bounds", struct{}{}, &r)
	if err != nil {
		return nil, err
	}
	return &Matrix{
		RGBA: image.NewRGBA(r),
		rpc:  rpc,
	}
}

func (m *Matrix) Close() error {
	return m.rpc.Close()
}

func (m *Matrix) Swap() error {
	return m.rpc.Call("Service.Swap", m.RGBA, nil)
}
