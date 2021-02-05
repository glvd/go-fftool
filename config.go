package tool

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

const sliceOutputTemplate = "-bsf:v,h264_mp4toannexb,-f,hls,-hls_list_size,0%s,-hls_time,%d,-hls_segment_filename,%s,%s"
const cryptoOutputTemplate = ",-hls_key_info_file,%s"
const scaleOutputTemplate = ",-vf,scale=-2:%d"
const cuvidScaleOutputTemplate = ",-vf,scale_npp=-2:%d"
const bitRateOutputTemplate = ",-b:v,%dK"
const frameRateOutputTemplate = ",-r,%3.2f"

const cudaOutputTemplate = ",-hwaccel,cuda"
const cuvidOutputTemplate = ",-hwaccel,cuvid,-c:v,h264_cuvid"
const defaultAccel = ",-c:v,%s"
const defaultTemplate = `-y%s,-i,%s,-strict,-2,-c:v,%s,-c:a,%s%s`

// None ...
const (
	ProcessNone             = ""
	ProcessH264CPU          = "libx264"
	ProcessH264QSV          = "h264_qsv"
	ProcessH264AMF          = "h264_amf"
	ProcessH264NVENC        = "h264_nvenc"
	ProcessH264VideoToolBox = "h264_videotoolbox"
	ProcessHevcCPU          = "libx265"
	ProcessHevcQSV          = "hevc_qsv"
	ProcessHevcAMF          = "hevc_amf"
	ProcessHevcNVENC        = "hevc_nvenc"
	ProcessHevcVideoToolBox = "hevc_videotoolbox"
	ProcessCUVID            = "cuvid"
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

// Scale ...
type Scale int

// Config ...
type Config struct {
	processID       string
	crypto          *Crypto
	output          string
	LogOutput       bool
	VideoFormat     string
	AudioFormat     string
	Scale           Scale
	ProcessCore     string
	BitRate         int64
	FrameRate       float64
	OutputPath      string //output path
	OutputName      string
	M3U8Name        string
	SegmentFileName string
	HLSTime         int
	KeyOutput       bool
	Slice           bool
	KeyPath         string
}

// CutOut ...
type CutOut struct {
	StartTime string
	EndTime   string
}

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

// DefaultOutputPath ...
var DefaultOutputPath = "video_split_temp"

// DefaultOutputName ...
var DefaultOutputName = "media.mp4"

// DefaultM3U8Name ...
var DefaultM3U8Name = "media.m3u8"

// DefaultSegmentFileName ...
var DefaultSegmentFileName = "media-%05d.ts"

// DefaultProcessCore ...
var DefaultProcessCore = ProcessH264CPU

// DefaultSlice ...
var DefaultSlice = false

// DefaultHLSTime ...
var DefaultHLSTime = 10

// DefaultScale ...
var DefaultScale = Scale720P

// DefaultKeyName ...
var DefaultKeyName = "m3u8_key"

// DefaultKeyInfoName ...
var DefaultKeyInfoName = "m3u8_key_info"

// DefaultKeyPath ...
var DefaultKeyPath = "output_key"

// DefaultProbeName ...
var DefaultProbeName = "ffprobe"

// DefaultMpegName ...
var DefaultMpegName = "ffmpeg"

// DefaultOpenSSLName ...
var DefaultOpenSSLName = "openssl"

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		crypto:          nil,
		output:          "",
		VideoFormat:     "libx265",
		AudioFormat:     "aac",
		Slice:           DefaultSlice,
		Scale:           Scale720P,
		ProcessCore:     DefaultProcessCore,
		BitRate:         0,
		FrameRate:       0,
		KeyOutput:       true,
		KeyPath:         DefaultKeyPath,
		OutputPath:      DefaultOutputPath,
		OutputName:      DefaultOutputName,
		M3U8Name:        DefaultM3U8Name,
		SegmentFileName: DefaultSegmentFileName,
		HLSTime:         DefaultHLSTime,
	}
}

func abs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		LogError(err)
		return ""
	}
	return abs
}

// SetCrypt ...
func (c *Config) SetCrypt(crypto *Crypto) {
	c.crypto = crypto
}

// CryptoInfo ...
func (c *Config) CryptoInfo() string {
	if c.crypto != nil {
		return fmt.Sprintf(cryptoOutputTemplate, c.crypto.KeyInfoPath)
	}
	return ""
}

// ProcessID ...
func (c *Config) ProcessID() string {
	return c.processID
}

// ActionOutput ...
func (c *Config) ActionOutput() string {
	if c.Slice {
		return fmt.Sprintf(sliceOutputTemplate, c.CryptoInfo(), c.HLSTime, filepath.Join(c.ProcessPath(), c.SegmentFileName), filepath.Join(c.ProcessPath(), c.M3U8Name))
	}
	return filepath.Join(c.ProcessPath(), c.OutputName)
}

// SaveKey ...
func (c *Config) SaveKey() error {
	if c.crypto != nil && c.KeyOutput {
		c.crypto.URL = DefaultKeyName
		c.crypto.KeyInfoPath = filepath.Join(abs(c.KeyPath), c.ProcessID(), DefaultKeyInfoName)
		c.crypto.KeyPath = filepath.Join(abs(c.KeyPath), c.ProcessID(), DefaultKeyName)
		if err := c.crypto.SaveKey(); err != nil {
			return err
		}
		if err := c.crypto.SaveKeyInfo(); err != nil {
			return err
		}
	}
	return nil
}

// ConfigOptions ...
func (c *Config) ConfigOptions() ConfigOptions {
	return func(cfg *Config) {
		*cfg = *c
	}
}

// ProcessPath ...
func (c *Config) ProcessPath() string {
	c.OutputPath = abs(c.OutputPath)
	return filepath.Join(c.OutputPath, c.ProcessID())
}

// Action ...
func (c *Config) Action() error {
	if c.Slice {
		return c.SaveKey()
	}
	return nil
}

// ScaleValue ...
func ScaleValue(scale Scale) int64 {
	i := int(scale)
	if len(scaleList) <= i {
		return 0
	}
	return scaleList[i]
}

func resolutionScale(v int64) Scale {
	switch {
	case v <= 480:
		return Scale480P
	case v > 960:
		return Scale1080P
	}
	return Scale720P
}

// OptimizeWithFormat ...
func OptimizeWithFormat(c *Config, sfmt *StreamFormat) (e error) {
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

	if video.CodecName == "h264" && c.Scale == -1 {
		c.VideoFormat = "copy"
	}

	if audio := sfmt.Audio(); audio != nil && audio.CodecName == "aac" {
		c.AudioFormat = "copy"
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

func outputArgs(c *Config, input string) string {
	var exts []interface{}

	//gpu decode config
	if c.VideoFormat != "copy" {
		c.VideoFormat = c.ProcessCore
	}

	//add scale setting
	if c.Scale != -1 {
		log.Infow("scale", "scale", c.Scale, "value", ScaleValue(c.Scale))
		if c.ProcessCore != ProcessCUVID {
			exts = append(exts, fmt.Sprintf(scaleOutputTemplate, ScaleValue(c.Scale)))
		} else {
			exts = append(exts, fmt.Sprintf(cuvidScaleOutputTemplate, ScaleValue(c.Scale)))
		}
	}
	//add bitrate arguments
	if c.BitRate != 0 {
		exts = append(exts, fmt.Sprintf(bitRateOutputTemplate, c.BitRate/1024))
	}
	//add frame rate arguments
	if c.FrameRate != 0 {
		exts = append(exts, fmt.Sprintf(frameRateOutputTemplate, c.FrameRate))
	}

	//generate slice arguments
	//output arguments
	return outputTemplate(c.ProcessCore, input, c.VideoFormat, c.AudioFormat, c.ActionOutput(), exts...)
}

func outputTemplate(core string, input, cv, ca, output string, exts ...interface{}) string {
	var outExt []string
	exts = append(exts, ","+output)
	for range exts {
		outExt = append(outExt, "%s")
	}
	var tmpl string
	//cuda/cpu/cuvid arguments case
	switch core {
	case ProcessH264CPU:
		tmpl = fmt.Sprintf(defaultTemplate, "", input, cv, ca, strings.Join(outExt, ""))
	case ProcessCUVID:
		tmpl = fmt.Sprintf(defaultTemplate, cuvidOutputTemplate, input, cv, ca, strings.Join(outExt, ""))
	case ProcessH264QSV,
		ProcessH264AMF,
		ProcessH264NVENC,
		ProcessH264VideoToolBox,
		ProcessHevcQSV,
		ProcessHevcAMF,
		ProcessHevcNVENC,
		ProcessHevcVideoToolBox:
		tmpl = fmt.Sprintf(defaultTemplate, "", input, cv, ca, strings.Join(outExt, ""))
	default:
		panic(fmt.Sprintf("wrong core type:%v", core))
	}
	log.Infow("format", "tmpl", tmpl, "output", output)
	return fmt.Sprintf(tmpl, exts...)
}
