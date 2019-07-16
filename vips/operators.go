package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

// Abs executes the 'abs' operation
func Abs(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("abs")
	err = vipsCall("abs", options)
	return out, err
}

// Add executes the 'add' operation
func Add(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("add")
	err = vipsCall("add", options)
	return out, err
}

// Analyzeload executes the 'analyzeload' operation
func Analyzeload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("analyzeload")
	err = vipsCall("analyzeload", options)
	return out, err
}

// Autorot executes the 'autorot' operation
func Autorot(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("autorot")
	err = vipsCall("autorot", options)
	return out, err
}

// Avg executes the 'avg' operation
func Avg(in *C.VipsImage, options ...*Option) (float64, error) {
	var out float64
	var err error
	options = append(options,
		InputImage("in", in),
		OutputDouble("out", &out),
	)
	incOpCounter("avg")
	err = vipsCall("avg", options)
	return out, err
}

// Bandbool executes the 'bandbool' operation
func Bandbool(in *C.VipsImage, boolean OperationBoolean, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("boolean", int(boolean)),
		OutputImage("out", &out),
	)
	incOpCounter("bandbool")
	err = vipsCall("bandbool", options)
	return out, err
}

// Bandfold executes the 'bandfold' operation
func Bandfold(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("bandfold")
	err = vipsCall("bandfold", options)
	return out, err
}

// Bandmean executes the 'bandmean' operation
func Bandmean(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("bandmean")
	err = vipsCall("bandmean", options)
	return out, err
}

// Bandunfold executes the 'bandunfold' operation
func Bandunfold(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("bandunfold")
	err = vipsCall("bandunfold", options)
	return out, err
}

// Black executes the 'black' operation
func Black(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("black")
	err = vipsCall("black", options)
	return out, err
}

// Boolean executes the 'boolean' operation
func Boolean(left *C.VipsImage, right *C.VipsImage, boolean OperationBoolean, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		InputInt("boolean", int(boolean)),
		OutputImage("out", &out),
	)
	incOpCounter("boolean")
	err = vipsCall("boolean", options)
	return out, err
}

// Buildlut executes the 'buildlut' operation
func Buildlut(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("buildlut")
	err = vipsCall("buildlut", options)
	return out, err
}

// Byteswap executes the 'byteswap' operation
func Byteswap(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("byteswap")
	err = vipsCall("byteswap", options)
	return out, err
}

// Cache executes the 'cache' operation
func Cache(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("cache")
	err = vipsCall("cache", options)
	return out, err
}

// Cast executes the 'cast' operation
func Cast(in *C.VipsImage, format BandFormat, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("format", int(format)),
		OutputImage("out", &out),
	)
	incOpCounter("cast")
	err = vipsCall("cast", options)
	return out, err
}

// Cmc2Lch executes the 'CMC2LCh' operation
func Cmc2Lch(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("CMC2LCh")
	err = vipsCall("CMC2LCh", options)
	return out, err
}

// Colourspace executes the 'colourspace' operation
func Colourspace(in *C.VipsImage, space Interpretation, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("space", int(space)),
		OutputImage("out", &out),
	)
	incOpCounter("colourspace")
	err = vipsCall("colourspace", options)
	return out, err
}

// Compass executes the 'compass' operation
func Compass(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("compass")
	err = vipsCall("compass", options)
	return out, err
}

// Complex executes the 'complex' operation
func Complex(in *C.VipsImage, cmplx OperationComplex, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("cmplx", int(cmplx)),
		OutputImage("out", &out),
	)
	incOpCounter("complex")
	err = vipsCall("complex", options)
	return out, err
}

// Complex2 executes the 'complex2' operation
func Complex2(left *C.VipsImage, right *C.VipsImage, cmplx OperationComplex2, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		InputInt("cmplx", int(cmplx)),
		OutputImage("out", &out),
	)
	incOpCounter("complex2")
	err = vipsCall("complex2", options)
	return out, err
}

// Complexform executes the 'complexform' operation
func Complexform(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("complexform")
	err = vipsCall("complexform", options)
	return out, err
}

// Complexget executes the 'complexget' operation
func Complexget(in *C.VipsImage, get OperationComplexGet, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("get", int(get)),
		OutputImage("out", &out),
	)
	incOpCounter("complexget")
	err = vipsCall("complexget", options)
	return out, err
}

// Conv executes the 'conv' operation
func Conv(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("conv")
	err = vipsCall("conv", options)
	return out, err
}

// Conva executes the 'conva' operation
func Conva(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("conva")
	err = vipsCall("conva", options)
	return out, err
}

// Convasep executes the 'convasep' operation
func Convasep(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("convasep")
	err = vipsCall("convasep", options)
	return out, err
}

// Convf executes the 'convf' operation
func Convf(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("convf")
	err = vipsCall("convf", options)
	return out, err
}

// Convi executes the 'convi' operation
func Convi(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("convi")
	err = vipsCall("convi", options)
	return out, err
}

// Convsep executes the 'convsep' operation
func Convsep(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("convsep")
	err = vipsCall("convsep", options)
	return out, err
}

// Copy executes the 'copy' operation
func Copy(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("copy")
	err = vipsCall("copy", options)
	return out, err
}

// Countlines executes the 'countlines' operation
func Countlines(in *C.VipsImage, direction Direction, options ...*Option) (float64, error) {
	var nolines float64
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("direction", int(direction)),
		OutputDouble("nolines", &nolines),
	)
	incOpCounter("countlines")
	err = vipsCall("countlines", options)
	return nolines, err
}

// Csvload executes the 'csvload' operation
func Csvload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("csvload")
	err = vipsCall("csvload", options)
	return out, err
}

// Csvsave executes the 'csvsave' operation
func Csvsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("csvsave")
	err = vipsCall("csvsave", options)
	return err
}

// De00 executes the 'dE00' operation
func De00(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("dE00")
	err = vipsCall("dE00", options)
	return out, err
}

// De76 executes the 'dE76' operation
func De76(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("dE76")
	err = vipsCall("dE76", options)
	return out, err
}

// Decmc executes the 'dECMC' operation
func Decmc(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("dECMC")
	err = vipsCall("dECMC", options)
	return out, err
}

// Deviate executes the 'deviate' operation
func Deviate(in *C.VipsImage, options ...*Option) (float64, error) {
	var out float64
	var err error
	options = append(options,
		InputImage("in", in),
		OutputDouble("out", &out),
	)
	incOpCounter("deviate")
	err = vipsCall("deviate", options)
	return out, err
}

// Divide executes the 'divide' operation
func Divide(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("divide")
	err = vipsCall("divide", options)
	return out, err
}

// DrawImage executes the 'draw_image' operation
func DrawImage(image *C.VipsImage, sub *C.VipsImage, x int, y int, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("image", image),
		InputImage("sub", sub),
		InputInt("x", x),
		InputInt("y", y),
	)
	incOpCounter("draw_image")
	err = vipsCall("draw_image", options)
	return err
}

// DrawSmudge executes the 'draw_smudge' operation
func DrawSmudge(image *C.VipsImage, left int, top int, width int, height int, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("image", image),
		InputInt("left", left),
		InputInt("top", top),
		InputInt("width", width),
		InputInt("height", height),
	)
	incOpCounter("draw_smudge")
	err = vipsCall("draw_smudge", options)
	return err
}

// Dzsave executes the 'dzsave' operation
func Dzsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("dzsave")
	err = vipsCall("dzsave", options)
	return err
}

// Embed executes the 'embed' operation
func Embed(in *C.VipsImage, x int, y int, width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("x", x),
		InputInt("y", y),
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("embed")
	err = vipsCall("embed", options)
	return out, err
}

// ExtractArea executes the 'extract_area' operation
func ExtractArea(input *C.VipsImage, left int, top int, width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("input", input),
		InputInt("left", left),
		InputInt("top", top),
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("extract_area")
	err = vipsCall("extract_area", options)
	return out, err
}

// ExtractBand executes the 'extract_band' operation
func ExtractBand(in *C.VipsImage, band int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("band", band),
		OutputImage("out", &out),
	)
	incOpCounter("extract_band")
	err = vipsCall("extract_band", options)
	return out, err
}

// Eye executes the 'eye' operation
func Eye(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("eye")
	err = vipsCall("eye", options)
	return out, err
}

// Falsecolour executes the 'falsecolour' operation
func Falsecolour(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("falsecolour")
	err = vipsCall("falsecolour", options)
	return out, err
}

// Fastcor executes the 'fastcor' operation
func Fastcor(in *C.VipsImage, ref *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("ref", ref),
		OutputImage("out", &out),
	)
	incOpCounter("fastcor")
	err = vipsCall("fastcor", options)
	return out, err
}

// FillNearest executes the 'fill_nearest' operation
func FillNearest(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("fill_nearest")
	err = vipsCall("fill_nearest", options)
	return out, err
}

// FindTrim executes the 'find_trim' operation
func FindTrim(in *C.VipsImage, options ...*Option) (int, int, int, int, error) {
	var left int
	var top int
	var width int
	var height int
	var err error
	options = append(options,
		InputImage("in", in),
		OutputInt("left", &left),
		OutputInt("top", &top),
		OutputInt("width", &width),
		OutputInt("height", &height),
	)
	incOpCounter("find_trim")
	err = vipsCall("find_trim", options)
	return left, top, width, height, err
}

// Flatten executes the 'flatten' operation
func Flatten(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("flatten")
	err = vipsCall("flatten", options)
	return out, err
}

// Flip executes the 'flip' operation
func Flip(in *C.VipsImage, direction Direction, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("direction", int(direction)),
		OutputImage("out", &out),
	)
	incOpCounter("flip")
	err = vipsCall("flip", options)
	return out, err
}

// Float2Rad executes the 'float2rad' operation
func Float2Rad(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("float2rad")
	err = vipsCall("float2rad", options)
	return out, err
}

// Fractsurf executes the 'fractsurf' operation
func Fractsurf(width int, height int, fractalDimension float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("fractal-dimension", fractalDimension),
		OutputImage("out", &out),
	)
	incOpCounter("fractsurf")
	err = vipsCall("fractsurf", options)
	return out, err
}

// Freqmult executes the 'freqmult' operation
func Freqmult(in *C.VipsImage, mask *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		OutputImage("out", &out),
	)
	incOpCounter("freqmult")
	err = vipsCall("freqmult", options)
	return out, err
}

// Fwfft executes the 'fwfft' operation
func Fwfft(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("fwfft")
	err = vipsCall("fwfft", options)
	return out, err
}

// Gamma executes the 'gamma' operation
func Gamma(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("gamma")
	err = vipsCall("gamma", options)
	return out, err
}

// Gaussblur executes the 'gaussblur' operation
func Gaussblur(in *C.VipsImage, sigma float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("sigma", sigma),
		OutputImage("out", &out),
	)
	incOpCounter("gaussblur")
	err = vipsCall("gaussblur", options)
	return out, err
}

// Gaussmat executes the 'gaussmat' operation
func Gaussmat(sigma float64, minAmpl float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputDouble("sigma", sigma),
		InputDouble("min-ampl", minAmpl),
		OutputImage("out", &out),
	)
	incOpCounter("gaussmat")
	err = vipsCall("gaussmat", options)
	return out, err
}

// Gaussnoise executes the 'gaussnoise' operation
func Gaussnoise(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("gaussnoise")
	err = vipsCall("gaussnoise", options)
	return out, err
}

// Gifload executes the 'gifload' operation
func Gifload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("gifload")
	err = vipsCall("gifload", options)
	return out, err
}

// Globalbalance executes the 'globalbalance' operation
func Globalbalance(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("globalbalance")
	err = vipsCall("globalbalance", options)
	return out, err
}

// Grey executes the 'grey' operation
func Grey(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("grey")
	err = vipsCall("grey", options)
	return out, err
}

// Grid executes the 'grid' operation
func Grid(in *C.VipsImage, tileHeight int, across int, down int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("tile-height", tileHeight),
		InputInt("across", across),
		InputInt("down", down),
		OutputImage("out", &out),
	)
	incOpCounter("grid")
	err = vipsCall("grid", options)
	return out, err
}

// HistCum executes the 'hist_cum' operation
func HistCum(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_cum")
	err = vipsCall("hist_cum", options)
	return out, err
}

// HistEntropy executes the 'hist_entropy' operation
func HistEntropy(in *C.VipsImage, options ...*Option) (float64, error) {
	var out float64
	var err error
	options = append(options,
		InputImage("in", in),
		OutputDouble("out", &out),
	)
	incOpCounter("hist_entropy")
	err = vipsCall("hist_entropy", options)
	return out, err
}

// HistEqual executes the 'hist_equal' operation
func HistEqual(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_equal")
	err = vipsCall("hist_equal", options)
	return out, err
}

// HistFind executes the 'hist_find' operation
func HistFind(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_find")
	err = vipsCall("hist_find", options)
	return out, err
}

// HistFindIndexed executes the 'hist_find_indexed' operation
func HistFindIndexed(in *C.VipsImage, index *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("index", index),
		OutputImage("out", &out),
	)
	incOpCounter("hist_find_indexed")
	err = vipsCall("hist_find_indexed", options)
	return out, err
}

// HistFindNdim executes the 'hist_find_ndim' operation
func HistFindNdim(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_find_ndim")
	err = vipsCall("hist_find_ndim", options)
	return out, err
}

// HistIsmonotonic executes the 'hist_ismonotonic' operation
func HistIsmonotonic(in *C.VipsImage, options ...*Option) (bool, error) {
	var monotonic bool
	var err error
	options = append(options,
		InputImage("in", in),
		OutputBool("monotonic", &monotonic),
	)
	incOpCounter("hist_ismonotonic")
	err = vipsCall("hist_ismonotonic", options)
	return monotonic, err
}

// HistLocal executes the 'hist_local' operation
func HistLocal(in *C.VipsImage, width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("hist_local")
	err = vipsCall("hist_local", options)
	return out, err
}

// HistMatch executes the 'hist_match' operation
func HistMatch(in *C.VipsImage, ref *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("ref", ref),
		OutputImage("out", &out),
	)
	incOpCounter("hist_match")
	err = vipsCall("hist_match", options)
	return out, err
}

// HistNorm executes the 'hist_norm' operation
func HistNorm(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_norm")
	err = vipsCall("hist_norm", options)
	return out, err
}

// HistPlot executes the 'hist_plot' operation
func HistPlot(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hist_plot")
	err = vipsCall("hist_plot", options)
	return out, err
}

// HoughCircle executes the 'hough_circle' operation
func HoughCircle(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hough_circle")
	err = vipsCall("hough_circle", options)
	return out, err
}

// HoughLine executes the 'hough_line' operation
func HoughLine(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("hough_line")
	err = vipsCall("hough_line", options)
	return out, err
}

// Hsv2Srgb executes the 'HSV2sRGB' operation
func Hsv2Srgb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("HSV2sRGB")
	err = vipsCall("HSV2sRGB", options)
	return out, err
}

// IccExport executes the 'icc_export' operation
func IccExport(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("icc_export")
	err = vipsCall("icc_export", options)
	return out, err
}

// IccImport executes the 'icc_import' operation
func IccImport(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("icc_import")
	err = vipsCall("icc_import", options)
	return out, err
}

// IccTransform executes the 'icc_transform' operation
func IccTransform(in *C.VipsImage, outputProfile string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("output-profile", outputProfile),
		OutputImage("out", &out),
	)
	incOpCounter("icc_transform")
	err = vipsCall("icc_transform", options)
	return out, err
}

// Identity executes the 'identity' operation
func Identity(options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,

		OutputImage("out", &out),
	)
	incOpCounter("identity")
	err = vipsCall("identity", options)
	return out, err
}

// Ifthenelse executes the 'ifthenelse' operation
func Ifthenelse(cond *C.VipsImage, in1 *C.VipsImage, in2 *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("cond", cond),
		InputImage("in1", in1),
		InputImage("in2", in2),
		OutputImage("out", &out),
	)
	incOpCounter("ifthenelse")
	err = vipsCall("ifthenelse", options)
	return out, err
}

// Insert executes the 'insert' operation
func Insert(main *C.VipsImage, sub *C.VipsImage, x int, y int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("main", main),
		InputImage("sub", sub),
		InputInt("x", x),
		InputInt("y", y),
		OutputImage("out", &out),
	)
	incOpCounter("insert")
	err = vipsCall("insert", options)
	return out, err
}

// Invert executes the 'invert' operation
func Invert(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("invert")
	err = vipsCall("invert", options)
	return out, err
}

// Invertlut executes the 'invertlut' operation
func Invertlut(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("invertlut")
	err = vipsCall("invertlut", options)
	return out, err
}

// Invfft executes the 'invfft' operation
func Invfft(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("invfft")
	err = vipsCall("invfft", options)
	return out, err
}

// Join executes the 'join' operation
func Join(in1 *C.VipsImage, in2 *C.VipsImage, direction Direction, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in1", in1),
		InputImage("in2", in2),
		InputInt("direction", int(direction)),
		OutputImage("out", &out),
	)
	incOpCounter("join")
	err = vipsCall("join", options)
	return out, err
}

// Jpegload executes the 'jpegload' operation
func Jpegload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("jpegload")
	err = vipsCall("jpegload", options)
	return out, err
}

// Jpegsave executes the 'jpegsave' operation
func Jpegsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("jpegsave")
	err = vipsCall("jpegsave", options)
	return err
}

// JpegsaveMime executes the 'jpegsave_mime' operation
func JpegsaveMime(in *C.VipsImage, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
	)
	incOpCounter("jpegsave_mime")
	err = vipsCall("jpegsave_mime", options)
	return err
}

// Lab2Labq executes the 'Lab2LabQ' operation
func Lab2Labq(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("Lab2LabQ")
	err = vipsCall("Lab2LabQ", options)
	return out, err
}

// Lab2Labs executes the 'Lab2LabS' operation
func Lab2Labs(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("Lab2LabS")
	err = vipsCall("Lab2LabS", options)
	return out, err
}

// Lab2Lch executes the 'Lab2LCh' operation
func Lab2Lch(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("Lab2LCh")
	err = vipsCall("Lab2LCh", options)
	return out, err
}

// Lab2Xyz executes the 'Lab2XYZ' operation
func Lab2Xyz(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("Lab2XYZ")
	err = vipsCall("Lab2XYZ", options)
	return out, err
}

// Labelregions executes the 'labelregions' operation
func Labelregions(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var mask *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("mask", &mask),
	)
	incOpCounter("labelregions")
	err = vipsCall("labelregions", options)
	return mask, err
}

// Labq2Lab executes the 'LabQ2Lab' operation
func Labq2Lab(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LabQ2Lab")
	err = vipsCall("LabQ2Lab", options)
	return out, err
}

// Labq2Labs executes the 'LabQ2LabS' operation
func Labq2Labs(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LabQ2LabS")
	err = vipsCall("LabQ2LabS", options)
	return out, err
}

// Labq2Srgb executes the 'LabQ2sRGB' operation
func Labq2Srgb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LabQ2sRGB")
	err = vipsCall("LabQ2sRGB", options)
	return out, err
}

// Labs2Lab executes the 'LabS2Lab' operation
func Labs2Lab(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LabS2Lab")
	err = vipsCall("LabS2Lab", options)
	return out, err
}

// Labs2Labq executes the 'LabS2LabQ' operation
func Labs2Labq(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LabS2LabQ")
	err = vipsCall("LabS2LabQ", options)
	return out, err
}

// Lch2Cmc executes the 'LCh2CMC' operation
func Lch2Cmc(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LCh2CMC")
	err = vipsCall("LCh2CMC", options)
	return out, err
}

// Lch2Lab executes the 'LCh2Lab' operation
func Lch2Lab(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("LCh2Lab")
	err = vipsCall("LCh2Lab", options)
	return out, err
}

// Linecache executes the 'linecache' operation
func Linecache(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("linecache")
	err = vipsCall("linecache", options)
	return out, err
}

// Logmat executes the 'logmat' operation
func Logmat(sigma float64, minAmpl float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputDouble("sigma", sigma),
		InputDouble("min-ampl", minAmpl),
		OutputImage("out", &out),
	)
	incOpCounter("logmat")
	err = vipsCall("logmat", options)
	return out, err
}

// Mapim executes the 'mapim' operation
func Mapim(in *C.VipsImage, index *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("index", index),
		OutputImage("out", &out),
	)
	incOpCounter("mapim")
	err = vipsCall("mapim", options)
	return out, err
}

// Maplut executes the 'maplut' operation
func Maplut(in *C.VipsImage, lut *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("lut", lut),
		OutputImage("out", &out),
	)
	incOpCounter("maplut")
	err = vipsCall("maplut", options)
	return out, err
}

// MaskButterworth executes the 'mask_butterworth' operation
func MaskButterworth(width int, height int, order float64, frequencyCutoff float64, amplitudeCutoff float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("order", order),
		InputDouble("frequency-cutoff", frequencyCutoff),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		OutputImage("out", &out),
	)
	incOpCounter("mask_butterworth")
	err = vipsCall("mask_butterworth", options)
	return out, err
}

// MaskButterworthBand executes the 'mask_butterworth_band' operation
func MaskButterworthBand(width int, height int, order float64, frequencyCutoffx float64, frequencyCutoffy float64, radius float64, amplitudeCutoff float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("order", order),
		InputDouble("frequency-cutoff-x", frequencyCutoffx),
		InputDouble("frequency-cutoff-y", frequencyCutoffy),
		InputDouble("radius", radius),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		OutputImage("out", &out),
	)
	incOpCounter("mask_butterworth_band")
	err = vipsCall("mask_butterworth_band", options)
	return out, err
}

// MaskButterworthRing executes the 'mask_butterworth_ring' operation
func MaskButterworthRing(width int, height int, order float64, frequencyCutoff float64, amplitudeCutoff float64, ringwidth float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("order", order),
		InputDouble("frequency-cutoff", frequencyCutoff),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		InputDouble("ringwidth", ringwidth),
		OutputImage("out", &out),
	)
	incOpCounter("mask_butterworth_ring")
	err = vipsCall("mask_butterworth_ring", options)
	return out, err
}

// MaskFractal executes the 'mask_fractal' operation
func MaskFractal(width int, height int, fractalDimension float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("fractal-dimension", fractalDimension),
		OutputImage("out", &out),
	)
	incOpCounter("mask_fractal")
	err = vipsCall("mask_fractal", options)
	return out, err
}

// MaskGaussian executes the 'mask_gaussian' operation
func MaskGaussian(width int, height int, frequencyCutoff float64, amplitudeCutoff float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff", frequencyCutoff),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		OutputImage("out", &out),
	)
	incOpCounter("mask_gaussian")
	err = vipsCall("mask_gaussian", options)
	return out, err
}

// MaskGaussianBand executes the 'mask_gaussian_band' operation
func MaskGaussianBand(width int, height int, frequencyCutoffx float64, frequencyCutoffy float64, radius float64, amplitudeCutoff float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff-x", frequencyCutoffx),
		InputDouble("frequency-cutoff-y", frequencyCutoffy),
		InputDouble("radius", radius),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		OutputImage("out", &out),
	)
	incOpCounter("mask_gaussian_band")
	err = vipsCall("mask_gaussian_band", options)
	return out, err
}

// MaskGaussianRing executes the 'mask_gaussian_ring' operation
func MaskGaussianRing(width int, height int, frequencyCutoff float64, amplitudeCutoff float64, ringwidth float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff", frequencyCutoff),
		InputDouble("amplitude-cutoff", amplitudeCutoff),
		InputDouble("ringwidth", ringwidth),
		OutputImage("out", &out),
	)
	incOpCounter("mask_gaussian_ring")
	err = vipsCall("mask_gaussian_ring", options)
	return out, err
}

// MaskIdeal executes the 'mask_ideal' operation
func MaskIdeal(width int, height int, frequencyCutoff float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff", frequencyCutoff),
		OutputImage("out", &out),
	)
	incOpCounter("mask_ideal")
	err = vipsCall("mask_ideal", options)
	return out, err
}

// MaskIdealBand executes the 'mask_ideal_band' operation
func MaskIdealBand(width int, height int, frequencyCutoffx float64, frequencyCutoffy float64, radius float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff-x", frequencyCutoffx),
		InputDouble("frequency-cutoff-y", frequencyCutoffy),
		InputDouble("radius", radius),
		OutputImage("out", &out),
	)
	incOpCounter("mask_ideal_band")
	err = vipsCall("mask_ideal_band", options)
	return out, err
}

// MaskIdealRing executes the 'mask_ideal_ring' operation
func MaskIdealRing(width int, height int, frequencyCutoff float64, ringwidth float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		InputDouble("frequency-cutoff", frequencyCutoff),
		InputDouble("ringwidth", ringwidth),
		OutputImage("out", &out),
	)
	incOpCounter("mask_ideal_ring")
	err = vipsCall("mask_ideal_ring", options)
	return out, err
}

// Match executes the 'match' operation
func Match(ref *C.VipsImage, sec *C.VipsImage, xr1 int, yr1 int, xs1 int, ys1 int, xr2 int, yr2 int, xs2 int, ys2 int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("ref", ref),
		InputImage("sec", sec),
		InputInt("xr1", xr1),
		InputInt("yr1", yr1),
		InputInt("xs1", xs1),
		InputInt("ys1", ys1),
		InputInt("xr2", xr2),
		InputInt("yr2", yr2),
		InputInt("xs2", xs2),
		InputInt("ys2", ys2),
		OutputImage("out", &out),
	)
	incOpCounter("match")
	err = vipsCall("match", options)
	return out, err
}

// Math executes the 'math' operation
func Math(in *C.VipsImage, math OperationMath, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("math", int(math)),
		OutputImage("out", &out),
	)
	incOpCounter("math")
	err = vipsCall("math", options)
	return out, err
}

// Math2 executes the 'math2' operation
func Math2(left *C.VipsImage, right *C.VipsImage, math2 OperationMath2, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		InputInt("math2", int(math2)),
		OutputImage("out", &out),
	)
	incOpCounter("math2")
	err = vipsCall("math2", options)
	return out, err
}

// Matrixload executes the 'matrixload' operation
func Matrixload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("matrixload")
	err = vipsCall("matrixload", options)
	return out, err
}

// Matrixprint executes the 'matrixprint' operation
func Matrixprint(in *C.VipsImage, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
	)
	incOpCounter("matrixprint")
	err = vipsCall("matrixprint", options)
	return err
}

// Matrixsave executes the 'matrixsave' operation
func Matrixsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("matrixsave")
	err = vipsCall("matrixsave", options)
	return err
}

// Max executes the 'max' operation
func Max(in *C.VipsImage, options ...*Option) (float64, error) {
	var out float64
	var err error
	options = append(options,
		InputImage("in", in),
		OutputDouble("out", &out),
	)
	incOpCounter("max")
	err = vipsCall("max", options)
	return out, err
}

// Measure executes the 'measure' operation
func Measure(in *C.VipsImage, h int, v int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("h", h),
		InputInt("v", v),
		OutputImage("out", &out),
	)
	incOpCounter("measure")
	err = vipsCall("measure", options)
	return out, err
}

// Merge executes the 'merge' operation
func Merge(ref *C.VipsImage, sec *C.VipsImage, direction Direction, dx int, dy int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("ref", ref),
		InputImage("sec", sec),
		InputInt("direction", int(direction)),
		InputInt("dx", dx),
		InputInt("dy", dy),
		OutputImage("out", &out),
	)
	incOpCounter("merge")
	err = vipsCall("merge", options)
	return out, err
}

// Min executes the 'min' operation
func Min(in *C.VipsImage, options ...*Option) (float64, error) {
	var out float64
	var err error
	options = append(options,
		InputImage("in", in),
		OutputDouble("out", &out),
	)
	incOpCounter("min")
	err = vipsCall("min", options)
	return out, err
}

// Morph executes the 'morph' operation
func Morph(in *C.VipsImage, mask *C.VipsImage, morph OperationMorphology, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("mask", mask),
		InputInt("morph", int(morph)),
		OutputImage("out", &out),
	)
	incOpCounter("morph")
	err = vipsCall("morph", options)
	return out, err
}

// Mosaic executes the 'mosaic' operation
func Mosaic(ref *C.VipsImage, sec *C.VipsImage, direction Direction, xref int, yref int, xsec int, ysec int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("ref", ref),
		InputImage("sec", sec),
		InputInt("direction", int(direction)),
		InputInt("xref", xref),
		InputInt("yref", yref),
		InputInt("xsec", xsec),
		InputInt("ysec", ysec),
		OutputImage("out", &out),
	)
	incOpCounter("mosaic")
	err = vipsCall("mosaic", options)
	return out, err
}

// Mosaic1 executes the 'mosaic1' operation
func Mosaic1(ref *C.VipsImage, sec *C.VipsImage, direction Direction, xr1 int, yr1 int, xs1 int, ys1 int, xr2 int, yr2 int, xs2 int, ys2 int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("ref", ref),
		InputImage("sec", sec),
		InputInt("direction", int(direction)),
		InputInt("xr1", xr1),
		InputInt("yr1", yr1),
		InputInt("xs1", xs1),
		InputInt("ys1", ys1),
		InputInt("xr2", xr2),
		InputInt("yr2", yr2),
		InputInt("xs2", xs2),
		InputInt("ys2", ys2),
		OutputImage("out", &out),
	)
	incOpCounter("mosaic1")
	err = vipsCall("mosaic1", options)
	return out, err
}

// Msb executes the 'msb' operation
func Msb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("msb")
	err = vipsCall("msb", options)
	return out, err
}

// Multiply executes the 'multiply' operation
func Multiply(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("multiply")
	err = vipsCall("multiply", options)
	return out, err
}

// Pdfload executes the 'pdfload' operation
func Pdfload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("pdfload")
	err = vipsCall("pdfload", options)
	return out, err
}

// Percent executes the 'percent' operation
func Percent(in *C.VipsImage, percent float64, options ...*Option) (int, error) {
	var threshold int
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("percent", percent),
		OutputInt("threshold", &threshold),
	)
	incOpCounter("percent")
	err = vipsCall("percent", options)
	return threshold, err
}

// Perlin executes the 'perlin' operation
func Perlin(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("perlin")
	err = vipsCall("perlin", options)
	return out, err
}

// Phasecor executes the 'phasecor' operation
func Phasecor(in *C.VipsImage, in2 *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("in2", in2),
		OutputImage("out", &out),
	)
	incOpCounter("phasecor")
	err = vipsCall("phasecor", options)
	return out, err
}

// Pngload executes the 'pngload' operation
func Pngload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("pngload")
	err = vipsCall("pngload", options)
	return out, err
}

// Pngsave executes the 'pngsave' operation
func Pngsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("pngsave")
	err = vipsCall("pngsave", options)
	return err
}

// Ppmload executes the 'ppmload' operation
func Ppmload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("ppmload")
	err = vipsCall("ppmload", options)
	return out, err
}

// Ppmsave executes the 'ppmsave' operation
func Ppmsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("ppmsave")
	err = vipsCall("ppmsave", options)
	return err
}

// Premultiply executes the 'premultiply' operation
func Premultiply(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("premultiply")
	err = vipsCall("premultiply", options)
	return out, err
}

// Profile executes the 'profile' operation
func Profile(in *C.VipsImage, options ...*Option) (*C.VipsImage, *C.VipsImage, error) {
	var columns *C.VipsImage
	var rows *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("columns", &columns),
		OutputImage("rows", &rows),
	)
	incOpCounter("profile")
	err = vipsCall("profile", options)
	return columns, rows, err
}

// Project executes the 'project' operation
func Project(in *C.VipsImage, options ...*Option) (*C.VipsImage, *C.VipsImage, error) {
	var columns *C.VipsImage
	var rows *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("columns", &columns),
		OutputImage("rows", &rows),
	)
	incOpCounter("project")
	err = vipsCall("project", options)
	return columns, rows, err
}

// Quadratic executes the 'quadratic' operation
func Quadratic(in *C.VipsImage, coeff *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("coeff", coeff),
		OutputImage("out", &out),
	)
	incOpCounter("quadratic")
	err = vipsCall("quadratic", options)
	return out, err
}

// Rad2Float executes the 'rad2float' operation
func Rad2Float(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("rad2float")
	err = vipsCall("rad2float", options)
	return out, err
}

// Radload executes the 'radload' operation
func Radload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("radload")
	err = vipsCall("radload", options)
	return out, err
}

// Radsave executes the 'radsave' operation
func Radsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("radsave")
	err = vipsCall("radsave", options)
	return err
}

// Rank executes the 'rank' operation
func Rank(in *C.VipsImage, width int, height int, index int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("width", width),
		InputInt("height", height),
		InputInt("index", index),
		OutputImage("out", &out),
	)
	incOpCounter("rank")
	err = vipsCall("rank", options)
	return out, err
}

// Rawload executes the 'rawload' operation
func Rawload(filename string, width int, height int, bands int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		InputInt("width", width),
		InputInt("height", height),
		InputInt("bands", bands),
		OutputImage("out", &out),
	)
	incOpCounter("rawload")
	err = vipsCall("rawload", options)
	return out, err
}

// Rawsave executes the 'rawsave' operation
func Rawsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("rawsave")
	err = vipsCall("rawsave", options)
	return err
}

// RawsaveFd executes the 'rawsave_fd' operation
func RawsaveFd(in *C.VipsImage, fd int, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("fd", fd),
	)
	incOpCounter("rawsave_fd")
	err = vipsCall("rawsave_fd", options)
	return err
}

// Recomb executes the 'recomb' operation
func Recomb(in *C.VipsImage, m *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("m", m),
		OutputImage("out", &out),
	)
	incOpCounter("recomb")
	err = vipsCall("recomb", options)
	return out, err
}

// Reduce executes the 'reduce' operation
func Reduce(in *C.VipsImage, hshrink float64, vshrink float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("hshrink", hshrink),
		InputDouble("vshrink", vshrink),
		OutputImage("out", &out),
	)
	incOpCounter("reduce")
	err = vipsCall("reduce", options)
	return out, err
}

// Reduceh executes the 'reduceh' operation
func Reduceh(in *C.VipsImage, hshrink float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("hshrink", hshrink),
		OutputImage("out", &out),
	)
	incOpCounter("reduceh")
	err = vipsCall("reduceh", options)
	return out, err
}

// Reducev executes the 'reducev' operation
func Reducev(in *C.VipsImage, vshrink float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("vshrink", vshrink),
		OutputImage("out", &out),
	)
	incOpCounter("reducev")
	err = vipsCall("reducev", options)
	return out, err
}

// Relational executes the 'relational' operation
func Relational(left *C.VipsImage, right *C.VipsImage, relational OperationRelational, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		InputInt("relational", int(relational)),
		OutputImage("out", &out),
	)
	incOpCounter("relational")
	err = vipsCall("relational", options)
	return out, err
}

// Remainder executes the 'remainder' operation
func Remainder(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("remainder")
	err = vipsCall("remainder", options)
	return out, err
}

// Replicate executes the 'replicate' operation
func Replicate(in *C.VipsImage, across int, down int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("across", across),
		InputInt("down", down),
		OutputImage("out", &out),
	)
	incOpCounter("replicate")
	err = vipsCall("replicate", options)
	return out, err
}

// Resize executes the 'resize' operation
func Resize(in *C.VipsImage, scale float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("scale", scale),
		OutputImage("out", &out),
	)
	incOpCounter("resize")
	err = vipsCall("resize", options)
	return out, err
}

// Rot executes the 'rot' operation
func Rot(in *C.VipsImage, angle Angle, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("angle", int(angle)),
		OutputImage("out", &out),
	)
	incOpCounter("rot")
	err = vipsCall("rot", options)
	return out, err
}

// Rot45 executes the 'rot45' operation
func Rot45(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("rot45")
	err = vipsCall("rot45", options)
	return out, err
}

// Round executes the 'round' operation
func Round(in *C.VipsImage, round OperationRound, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("round", int(round)),
		OutputImage("out", &out),
	)
	incOpCounter("round")
	err = vipsCall("round", options)
	return out, err
}

// Scale executes the 'scale' operation
func Scale(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("scale")
	err = vipsCall("scale", options)
	return out, err
}

// Scrgb2Bw executes the 'scRGB2BW' operation
func Scrgb2Bw(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("scRGB2BW")
	err = vipsCall("scRGB2BW", options)
	return out, err
}

// Scrgb2Srgb executes the 'scRGB2sRGB' operation
func Scrgb2Srgb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("scRGB2sRGB")
	err = vipsCall("scRGB2sRGB", options)
	return out, err
}

// Scrgb2Xyz executes the 'scRGB2XYZ' operation
func Scrgb2Xyz(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("scRGB2XYZ")
	err = vipsCall("scRGB2XYZ", options)
	return out, err
}

// Sequential executes the 'sequential' operation
func Sequential(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("sequential")
	err = vipsCall("sequential", options)
	return out, err
}

// Sharpen executes the 'sharpen' operation
func Sharpen(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("sharpen")
	err = vipsCall("sharpen", options)
	return out, err
}

// Shrink executes the 'shrink' operation
func Shrink(in *C.VipsImage, hshrink float64, vshrink float64, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputDouble("hshrink", hshrink),
		InputDouble("vshrink", vshrink),
		OutputImage("out", &out),
	)
	incOpCounter("shrink")
	err = vipsCall("shrink", options)
	return out, err
}

// Shrinkh executes the 'shrinkh' operation
func Shrinkh(in *C.VipsImage, hshrink int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("hshrink", hshrink),
		OutputImage("out", &out),
	)
	incOpCounter("shrinkh")
	err = vipsCall("shrinkh", options)
	return out, err
}

// Shrinkv executes the 'shrinkv' operation
func Shrinkv(in *C.VipsImage, vshrink int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("vshrink", vshrink),
		OutputImage("out", &out),
	)
	incOpCounter("shrinkv")
	err = vipsCall("shrinkv", options)
	return out, err
}

// Sign executes the 'sign' operation
func Sign(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("sign")
	err = vipsCall("sign", options)
	return out, err
}

// Similarity executes the 'similarity' operation
func Similarity(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("similarity")
	err = vipsCall("similarity", options)
	return out, err
}

// Sines executes the 'sines' operation
func Sines(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("sines")
	err = vipsCall("sines", options)
	return out, err
}

// Smartcrop executes the 'smartcrop' operation
func Smartcrop(input *C.VipsImage, width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("input", input),
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("smartcrop")
	err = vipsCall("smartcrop", options)
	return out, err
}

// Spcor executes the 'spcor' operation
func Spcor(in *C.VipsImage, ref *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputImage("ref", ref),
		OutputImage("out", &out),
	)
	incOpCounter("spcor")
	err = vipsCall("spcor", options)
	return out, err
}

// Spectrum executes the 'spectrum' operation
func Spectrum(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("spectrum")
	err = vipsCall("spectrum", options)
	return out, err
}

// Srgb2Hsv executes the 'sRGB2HSV' operation
func Srgb2Hsv(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("sRGB2HSV")
	err = vipsCall("sRGB2HSV", options)
	return out, err
}

// Srgb2Scrgb executes the 'sRGB2scRGB' operation
func Srgb2Scrgb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("sRGB2scRGB")
	err = vipsCall("sRGB2scRGB", options)
	return out, err
}

// Stats executes the 'stats' operation
func Stats(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("stats")
	err = vipsCall("stats", options)
	return out, err
}

// Stdif executes the 'stdif' operation
func Stdif(in *C.VipsImage, width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("stdif")
	err = vipsCall("stdif", options)
	return out, err
}

// Subsample executes the 'subsample' operation
func Subsample(input *C.VipsImage, xfac int, yfac int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("input", input),
		InputInt("xfac", xfac),
		InputInt("yfac", yfac),
		OutputImage("out", &out),
	)
	incOpCounter("subsample")
	err = vipsCall("subsample", options)
	return out, err
}

// Subtract executes the 'subtract' operation
func Subtract(left *C.VipsImage, right *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("left", left),
		InputImage("right", right),
		OutputImage("out", &out),
	)
	incOpCounter("subtract")
	err = vipsCall("subtract", options)
	return out, err
}

// Svgload executes the 'svgload' operation
func Svgload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("svgload")
	err = vipsCall("svgload", options)
	return out, err
}

// System executes the 'system' operation
func System(cmdFormat string, options ...*Option) error {
	var err error
	options = append(options,
		InputString("cmd-format", cmdFormat),
	)
	incOpCounter("system")
	err = vipsCall("system", options)
	return err
}

// Text executes the 'text' operation
func Text(text string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("text", text),
		OutputImage("out", &out),
	)
	incOpCounter("text")
	err = vipsCall("text", options)
	return out, err
}

// Thumbnail executes the 'thumbnail' operation
func Thumbnail(filename string, width int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		InputInt("width", width),
		OutputImage("out", &out),
	)
	incOpCounter("thumbnail")
	err = vipsCall("thumbnail", options)
	return out, err
}

// ThumbnailImage executes the 'thumbnail_image' operation
func ThumbnailImage(in *C.VipsImage, width int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		InputInt("width", width),
		OutputImage("out", &out),
	)
	incOpCounter("thumbnail_image")
	err = vipsCall("thumbnail_image", options)
	return out, err
}

// Tiffload executes the 'tiffload' operation
func Tiffload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("tiffload")
	err = vipsCall("tiffload", options)
	return out, err
}

// Tiffsave executes the 'tiffsave' operation
func Tiffsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("tiffsave")
	err = vipsCall("tiffsave", options)
	return err
}

// Tilecache executes the 'tilecache' operation
func Tilecache(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("tilecache")
	err = vipsCall("tilecache", options)
	return out, err
}

// Tonelut executes the 'tonelut' operation
func Tonelut(options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,

		OutputImage("out", &out),
	)
	incOpCounter("tonelut")
	err = vipsCall("tonelut", options)
	return out, err
}

// Unpremultiply executes the 'unpremultiply' operation
func Unpremultiply(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("unpremultiply")
	err = vipsCall("unpremultiply", options)
	return out, err
}

// Vipsload executes the 'vipsload' operation
func Vipsload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("vipsload")
	err = vipsCall("vipsload", options)
	return out, err
}

// Vipssave executes the 'vipssave' operation
func Vipssave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("vipssave")
	err = vipsCall("vipssave", options)
	return err
}

// Webpload executes the 'webpload' operation
func Webpload(filename string, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputString("filename", filename),
		OutputImage("out", &out),
	)
	incOpCounter("webpload")
	err = vipsCall("webpload", options)
	return out, err
}

// Webpsave executes the 'webpsave' operation
func Webpsave(in *C.VipsImage, filename string, options ...*Option) error {
	var err error
	options = append(options,
		InputImage("in", in),
		InputString("filename", filename),
	)
	incOpCounter("webpsave")
	err = vipsCall("webpsave", options)
	return err
}

// Worley executes the 'worley' operation
func Worley(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("worley")
	err = vipsCall("worley", options)
	return out, err
}

// Wrap executes the 'wrap' operation
func Wrap(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("wrap")
	err = vipsCall("wrap", options)
	return out, err
}

// Xyz executes the 'xyz' operation
func Xyz(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("xyz")
	err = vipsCall("xyz", options)
	return out, err
}

// Xyz2Lab executes the 'XYZ2Lab' operation
func Xyz2Lab(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("XYZ2Lab")
	err = vipsCall("XYZ2Lab", options)
	return out, err
}

// Xyz2Scrgb executes the 'XYZ2scRGB' operation
func Xyz2Scrgb(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("XYZ2scRGB")
	err = vipsCall("XYZ2scRGB", options)
	return out, err
}

// Xyz2Yxy executes the 'XYZ2Yxy' operation
func Xyz2Yxy(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("XYZ2Yxy")
	err = vipsCall("XYZ2Yxy", options)
	return out, err
}

// Yxy2Xyz executes the 'Yxy2XYZ' operation
func Yxy2Xyz(in *C.VipsImage, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("in", in),
		OutputImage("out", &out),
	)
	incOpCounter("Yxy2XYZ")
	err = vipsCall("Yxy2XYZ", options)
	return out, err
}

// Zone executes the 'zone' operation
func Zone(width int, height int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputInt("width", width),
		InputInt("height", height),
		OutputImage("out", &out),
	)
	incOpCounter("zone")
	err = vipsCall("zone", options)
	return out, err
}

// Zoom executes the 'zoom' operation
func Zoom(input *C.VipsImage, xfac int, yfac int, options ...*Option) (*C.VipsImage, error) {
	var out *C.VipsImage
	var err error
	options = append(options,
		InputImage("input", input),
		InputInt("xfac", xfac),
		InputInt("yfac", yfac),
		OutputImage("out", &out),
	)
	incOpCounter("zoom")
	err = vipsCall("zoom", options)
	return out, err
}
