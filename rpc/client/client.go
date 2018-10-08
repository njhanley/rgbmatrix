package client

import (
	"image"
	"image/color"
	"image/draw"
	"net/rpc"

	"github.com/njhanley/rgbmatrix/rpc/internal/types"
)

type Client struct {
	rpc *rpc.Client
}

func Connect(addr string) (*Client, error) {
	rpc, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{rpc}, nil
}

func (cl *Client) Close() error {
	return cl.rpc.Close()
}

func (cl *Client) Size() (image.Point, error) {
	var size image.Point
	err := cl.rpc.Call("Server.Size", types.None{}, &size)
	return size, err
}

func (cl *Client) SwapOnVSync() error {
	return cl.rpc.Call("Server.SwapOnVSync", types.None{}, nil)
}

func (cl *Client) Clear() error {
	return cl.rpc.Call("Server.Clear", types.None{}, nil)
}

func (cl *Client) Fill(c color.Color) error {
	return cl.rpc.Call("Server.Fill", color.RGBAModel.Convert(c).(color.RGBA), nil)
}

func (cl *Client) Set(x, y int, c color.Color) error {
	return cl.rpc.Call("Server.Set", types.Pixel{x, y, color.RGBAModel.Convert(c).(color.RGBA)}, nil)
}

func (cl *Client) Draw(m image.Image) error {
	rgba, ok := m.(*image.RGBA)
	if !ok {
		b := m.Bounds()
		rgba = image.NewRGBA(b)
		draw.Draw(rgba, b, m, b.Min, draw.Src)
	}
	return cl.rpc.Call("Server.Draw", rgba, nil)
}
