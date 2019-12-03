package fftool

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goextension/log"
)

const sliceOutputTemplate = "-bsf:v,h264_mp4toannexb,-f,hls,-hls_list_size,0,-hls_time,%d,-hls_segment_filename,%s,%s"
const scaleOutputTemplate = ",-vf,scale=-2:%d"
const cuvidScaleOutputTemplate = ",-vf,scale_npp=-2:%d"
const bitRateOutputTemplate = ",-b:v,%dK"
const frameRateOutputTemplate = ",-r,%3.2f"
const cudaOutputTemplate = "-hwaccel,cuda"

const cuvidOutputTemplate = "-hwaccel,cuvid,-c:v,h264_cuvid"

const defaultTemplate = `-y%s,-i,%s,-strict,-2,-c:v,%s,-c:a,%s%s,%s`

// None ...
const (
	ProcessNone ProcessCore = -1
	ProcessCPU  ProcessCore = 1
	ProcessCUDA ProcessCore = iota
	ProcessCUVID
)

// Scale ...
const (
	ScaleNone  Scale = -1
	Scale480P  Scale = 0
	Scale720P  Scale = 1
	Scale1080P Scale = 2
	//Scale2K    Scale = 3
	//Scale4K    Scale = 4
	//Scale8K    Scale = 5
)

// ProcessCore ...
type ProcessCore int

// Scale ...
type Scale int

var scaleList = []int64{
	0: 480,
	1: 720,
	2: 1080,
}

var bitRateList = []int64{
	//Scale480P:  1000 * 1024,
	//Scale720P:  2000 * 1024,
	//Scale1080P: 4000 * 1024,
	Scale480P:  500 * 1024,
	Scale720P:  1000 * 1024,
	Scale1080P: 2000 * 1024,
	//Scale2K:    4000 * 1024,
	//Scale4K:    8000 * 1024,
}

var frameRateList = []float64{
	Scale480P:  float64(24000)/1001 - 0.005,
	Scale720P:  float64(24000)/1001 - 0.005,
	Scale1080P: float64(30000)/1001 - 0.005,
	//Scale2K:    float64(30000)/1001 - 0.005,
	//Scale4K:    float64(30000)/1001 - 0.005,
}

//Crypto ...
type Crypto struct {
	Key string
}

// CutOut ...
type CutOut struct {
	StartTime string
	EndTime   string
}

// Config ...
type Config struct {
	videoFormat     string
	audioFormat     string
	Scale           Scale
	ProcessCore     ProcessCore
	NeedSlice       bool
	BitRate         int64
	FrameRate       float64
	OutputPath      string //output path
	OutputName      string
	M3U8Name        string
	SegmentFileName string
	HLSTime         int
}

// DefaultOutputPath ...
var DefaultOutputPath = "video_split_temp"

// DefaultOutputName ...
var DefaultOutputName = "media.mp4"

// DefaultM3U8Name ...
var DefaultM3U8Name = "media.m3u8"

// DefaultSegmentFileName ...
var DefaultSegmentFileName = "media-%05d.ts"

// DefaultProcessCore ...
var DefaultProcessCore = ProcessCUDA

// DefaultSlice ...
var DefaultSlice = false

// DefaultHLSTime ...
var DefaultHLSTime = 10

// DefaultScale ...
var DefaultScale = Scale720P

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Scale:           Scale720P,
		ProcessCore:     DefaultProcessCore,
		NeedSlice:       DefaultSlice,
		BitRate:         0,
		FrameRate:       0,
		OutputPath:      DefaultOutputPath,
		OutputName:      DefaultOutputName,
		videoFormat:     "libx264",
		audioFormat:     "aac",
		M3U8Name:        DefaultM3U8Name,
		SegmentFileName: DefaultSegmentFileName,
		HLSTime:         DefaultHLSTime,
	}
}

func (c *Config) init() {

}

// AbsOutput ...
func (c *Config) AbsOutput() string {
	if filepath.IsAbs(c.OutputPath) {
		return c.OutputPath
	}
	abs, err := filepath.Abs(c.OutputPath)
	if err != nil {
		return DefaultOutputPath
	}
	return abs
}

// Args ...
func (c *Config) Args(input, output string) string {
	var exts []interface{}

	if c.ProcessCore != ProcessCPU && c.videoFormat != "copy" {
		c.videoFormat = "h264_nvenc"
	}

	if c.Scale != -1 {
		log.Infow("scale", "scale", c.Scale, "value", scaleVale(c.Scale))
		if c.ProcessCore != ProcessCUVID {
			exts = append(exts, fmt.Sprintf(scaleOutputTemplate, scaleVale(c.Scale)))
		} else {
			exts = append(exts, fmt.Sprintf(cuvidScaleOutputTemplate, scaleVale(c.Scale)))
		}
	}
	if c.BitRate != 0 {
		exts = append(exts, fmt.Sprintf(bitRateOutputTemplate, c.BitRate/1024))
	}
	if c.FrameRate != 0 {
		exts = append(exts, fmt.Sprintf(frameRateOutputTemplate, c.FrameRate))
	}
	output = filepath.Join(c.AbsOutput(), c.OutputName)
	if c.NeedSlice {
		if filepath.Ext(c.OutputName) != "" {
			panic(fmt.Sprintf("slice cannot output with name %s", c.OutputName))
		}
		output = fmt.Sprintf(sliceOutputTemplate, c.HLSTime, filepath.Join(output, c.SegmentFileName), filepath.Join(output, c.M3U8Name))
	}

	return outputTemplate(c.ProcessCore, input, c.videoFormat, c.audioFormat, output, exts...)
}

func outputTemplate(p ProcessCore, input, cv, ca, output string, exts ...interface{}) string {
	var outExt []string
	for range exts {
		outExt = append(outExt, "%s")
	}
	def := ""
	if p == ProcessCPU {
		def = fmt.Sprintf(defaultTemplate, "", input, cv, ca, strings.Join(outExt, " "), output)
	} else {
		def = fmt.Sprintf(defaultTemplate, cudaOutputTemplate, input, cv, ca, strings.Join(outExt, " "), output)
	}
	log.Infow("format", "def", def)
	return fmt.Sprintf(def, exts...)
}

func scaleVale(scale Scale) int64 {
	i := int(scale)
	if len(scaleList) <= i {
		return 0
	}
	return scaleList[i]
}

func resolutionScale(v int64) Scale {
	//r := getResolution(v, 0, -1)
	switch {
	case v <= 480:
		return Scale480P
	case v > 960:
		return Scale1080P
	}
	return Scale720P
}

// Clone ...
func (c *Config) Clone() Config {
	return *c
}

// OptimizeWithFormat ...
func (c *Config) OptimizeWithFormat(sfmt *StreamFormat) (e error) {
	return optimizeWithFormat(c, sfmt)
}

func optimizeWithFormat(c *Config, sfmt *StreamFormat) (e error) {
	if sfmt == nil {
		return errors.New("format is null")
	}
	video := sfmt.Video()
	if video == nil {
		return errors.New("video is null")
	}

	i, e := strconv.ParseInt(video.BitRate, 10, 64)
	if e != nil {
		i = math.MaxInt64
		log.Errorw("parse:bitrate", "error", e)
	}

	e = optimizeBitRate(c, *video.Height, i)
	if e != nil {
		return e
	}

	e = optimizeFrameRate(c, video.RFrameRate)
	if e != nil {
		return e
	}

	if video.CodecName == "h264" && c.Scale == 0 {
		c.videoFormat = "copy"
	}

	if audio := sfmt.Audio(); audio != nil && audio.CodecName == "aac" {
		c.audioFormat = "copy"
	}

	return nil
}

func optimizeBitRate(c *Config, height int64, bitRate int64) (e error) {
	scale := resolutionScale(height)
	if c.Scale > scale {
		c.Scale = scale
	}

	if c.BitRate == 0 {
		c.BitRate = bitRateList[c.Scale]
	}
	if c.BitRate > bitRate {
		c.BitRate = 0
	}
	return nil
}

func optimizeFrameRate(c *Config, frameRate string) (e error) {
	fr := strings.Split(frameRate, "/")
	il := 1
	ir := 1
	if len(fr) == 2 {
		il, e = strconv.Atoi(fr[0])
		if e != nil {
			il = 1
			log.Errorw("parse:il", "error", e, "framerate", frameRate)
		}
		ir, e = strconv.Atoi(fr[1])
		if e != nil {
			ir = 1
			log.Errorw("parse:ir", "error", e, "framerate", frameRate)
		}
	}
	if c.FrameRate == 0 {
		c.FrameRate = frameRateList[c.Scale]
	}

	if c.FrameRate > float64(il)/float64(ir) {
		c.FrameRate = 0
	}
	log.Infow("info", "framerate", c.FrameRate, "il", il, "ir", ir, "il/ir", il/ir)
	return nil
}
