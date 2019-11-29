package fftool

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/goextension/log"
)

//const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -ss %s -to %s -c:v %s -c:a %s -bsf:v h264_mp4toannexb -vsync 0 -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -c:v %s -c:a %s -bsf:v h264_mp4toannexb -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const sliceM3u8ScaleTemplate = `-y -i %s -strict -2 -c:v %s -c:a %s -bsf:v h264_mp4toannexb %s -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const scaleOutputTemplate = ",-vf,scale=-2:%d"
const bitRateOutputTemplate = ",-b:v,%dK"
const frameRateOutputTemplate = ",-r,%3.2f"

const defaultTemplate = `-y,-i,%s,-strict,-2,-c:v,%s,-c:a,%s%s,%s`

// Scale ...
type Scale int

// Scale ...
const (
	Scale480P  Scale = 0
	Scale720P  Scale = 1
	Scale1080P Scale = 2
)

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
}

var frameRateList = []float64{
	Scale480P:  float64(24000)/1001 - 0.005,
	Scale720P:  float64(24000)/1001 - 0.005,
	Scale1080P: float64(30000)/1001 - 0.005,
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
	Scale           Scale
	BitRate         int64
	FrameRate       float64
	Output          string //output path
	VideoFormat     string
	AudioFormat     string
	M3U8Name        string
	SegmentFileName string
	HLSTime         int
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		Scale:           Scale720P,
		BitRate:         0,
		FrameRate:       0,
		Output:          "video_split_temp",
		VideoFormat:     "libx264",
		AudioFormat:     "aac",
		M3U8Name:        "media.m3u8",
		SegmentFileName: "media-%05d.ts",
		HLSTime:         10,
	}
}

func (c *Config) init() {

}

// Args ...
func (c *Config) Args() string {
	panic("args")
}

func scaleVale(scale Scale) int64 {
	i := int(scale)
	if len(scaleList) >= i {
		return 0
	}
	return scaleList[i]
}

func resolutionScale(v int64) Scale {
	r := getResolution(v, 0, -1)
	switch {
	case r <= 480:
		return Scale480P
	case r > 720:
		return Scale1080P
	}
	return Scale720P
}

// OptimizeWithFormat ...
func (c *Config) OptimizeWithFormat(sfmt *StreamFormat) (string, error) {
	video := sfmt.Video()
	if video == nil {
		return "", errors.New("video is null")
	}
	e := c.optimizeBitRate(video)
	if e != nil {
		return "", e
	}
	e = c.optimizeFrameRate(video)
	if e != nil {
		return "", e
	}
	return "", nil
}

func (c *Config) optimizeBitRate(video *Stream) (e error) {
	scale := resolutionScale(*video.Height)
	if c.Scale > scale {
		c.Scale = scale
	}

	i, e := strconv.ParseInt(video.BitRate, 10, 64)
	if e != nil {
		i = math.MaxInt64
		log.Errorw("parse:bitrate", "error", e)
	}

	if c.BitRate == 0 {
		c.BitRate = bitRateList[c.Scale]
	}
	if c.BitRate > i {
		c.BitRate = 0
	}
	return nil
}

func (c *Config) optimizeFrameRate(video *Stream) (e error) {
	fr := strings.Split(video.RFrameRate, "/")
	il := 1
	ir := 1
	if len(fr) == 2 {
		il, e = strconv.Atoi(fr[0])
		if e != nil {
			il = 1
			log.Errorw("parse:il", "error", e, "framerate", video.RFrameRate)
		}
		ir, e = strconv.Atoi(fr[1])
		if e != nil {
			ir = 1
			log.Errorw("parse:ir", "error", e, "framerate", video.RFrameRate)
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
