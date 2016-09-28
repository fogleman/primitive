package primitive

import "github.com/golang/freetype/raster"

var theRasterizerPool rasterizerPool
var theScanlinePool scanlinePool

func InitPools(n int) {
	theRasterizerPool = newRasterizerPool(n)
	theScanlinePool = newScanlinePool(n)
}

type rasterizerPool chan *raster.Rasterizer

func newRasterizerPool(n int) rasterizerPool {
	ch := make(chan *raster.Rasterizer, n)
	for i := 0; i < n; i++ {
		ch <- raster.NewRasterizer(1, 1)
	}
	return ch
}

func (p rasterizerPool) get() *raster.Rasterizer {
	return <-p
}

func (p rasterizerPool) put(r *raster.Rasterizer) {
	p <- r
}

type scanlinePool chan []Scanline

func newScanlinePool(n int) scanlinePool {
	ch := make(chan []Scanline, n)
	for i := 0; i < n; i++ {
		ch <- make([]Scanline, 0, 8192)
	}
	return ch
}

func (p scanlinePool) get() []Scanline {
	return <-p
}

func (p scanlinePool) put(r []Scanline) {
	p <- r
}
