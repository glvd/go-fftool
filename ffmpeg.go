package fftool

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/goextension/log"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

// SplitArgs ...
type SplitArgs struct {
	StreamFormat    *StreamFormat
	Auto            bool
	Scale           int64
	Start           string
	End             string
	Output          string
	Video           string
	Audio           string
	M3U8            string
	SegmentFileName string
	HLSTime         int
	probe           func(string) (*StreamFormat, error)
	BitRate         int64
	FrameRate       float64
}

// FFmpegContext ...
type ffmpegContext struct {
	once   sync.Once
	mu     sync.RWMutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	done   chan bool
}

// Context ...
func (c *ffmpegContext) Context() context.Context {
	return c.ctx
}

// Add ...
func (c *ffmpegContext) Add(i int) {
	c.wg.Add(i)
}

// Wait ...
func (c *ffmpegContext) Wait() {
	select {
	case <-c.Waiting():
		return
	}
}

// Waiting ...
func (c *ffmpegContext) Waiting() <-chan bool {
	c.once.Do(func() {
		go func() {
			c.wg.Wait()
			c.done <- true
		}()
	})

	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan bool)
	}
	d := c.done
	c.mu.Unlock()
	return d
}

// Done ...
func (c *ffmpegContext) Done() {
	c.wg.Done()
}

// Cancel ...
func (c *ffmpegContext) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Context ...
type Context interface {
	Cancel()
	Add(int)
	Waiting() <-chan bool
	Wait()
	Done()
	Context() context.Context
}

// FFmpegContext ...
func FFmpegContext() Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &ffmpegContext{
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

// SplitOptions ...
type SplitOptions func(args *SplitArgs)

// HLSTimeOption ...
func HLSTimeOption(i int) SplitOptions {
	return func(args *SplitArgs) {
		args.HLSTime = i
	}
}

// ScaleOption ...
func ScaleOption(s int64, v ...string) SplitOptions {
	return func(args *SplitArgs) {
		args.Video = "libx264"
		for _, value := range v {
			args.Video = value
		}
		args.Scale = s
	}
}

// OutputOption ...
func OutputOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Output = s
	}
}

// AutoOption ...
func AutoOption(s bool) SplitOptions {
	return func(args *SplitArgs) {
		args.Auto = s
	}
}

// VideoOption ...
func VideoOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Video = s
	}
}

// AudioOption ...
func AudioOption(s string) SplitOptions {
	return func(args *SplitArgs) {
		args.Video = s
	}
}

// StreamFormatOption ...
func StreamFormatOption(s *StreamFormat) SplitOptions {
	return func(args *SplitArgs) {
		args.StreamFormat = s
	}
}

// BitRateOption ...
func BitRateOption(b int64) SplitOptions {
	return func(args *SplitArgs) {
		args.BitRate = b
	}
}

// ProbeInfoOption ...
func ProbeInfoOption(f func(string) (*StreamFormat, error)) SplitOptions {
	return func(args *SplitArgs) {
		args.probe = f
	}
}

// FFMpegSplitToM3U8WithProbe ...
func FFMpegSplitToM3U8WithProbe(ctx context.Context, file string, args ...SplitOptions) (sa *SplitArgs, e error) {
	args = append(args, ProbeInfoOption(FFProbeStreamFormat))
	return FFMpegSplitToM3U8(ctx, file, args...)
}

func scaleIndex(scale int64) int {
	if scale == 480 {
		return int(Scale480P)
	} else if scale == 1080 {
		return int(Scale1080P)
	}
	return int(Scale720P)
}

func outputScale(sa *SplitArgs) string {
	outputs := []string{fmt.Sprintf(scaleOutputTemplate, sa.Scale)}

	if sa.BitRate != 0 {
		outputs = append(outputs, fmt.Sprintf(bitRateOutputTemplate, sa.BitRate/1024))
	}
	log.Info(sa.FrameRate)
	if sa.FrameRate > 0 {
		outputs = append(outputs, fmt.Sprintf(frameRateOutputTemplate, sa.FrameRate))
	}
	log.Info("output:", strings.Join(outputs, " "))
	return strings.Join(outputs, " ")
}

// optimizeScale ...
func optimizeScale(sa *SplitArgs, video *Stream) {
	if sa.Scale != 0 {
		if video.Height != nil && *video.Height < sa.Scale {
			//pass when video is smaller then input
			sa.Scale = 0
			return
		}

		idx := scaleIndex(sa.Scale)
		i, e := strconv.ParseInt(video.BitRate, 10, 64)
		if e != nil {
			log.Error(e)
			i = math.MaxInt64
		}

		if sa.BitRate == 0 {
			sa.BitRate = bitRateList[idx]
		}
		if sa.BitRate != 0 {
			if sa.BitRate > i {
				sa.BitRate = 0
			}
		}
		log.Info(video.RFrameRate)
		fr := strings.Split(video.RFrameRate, "/")
		il := 1
		ir := 1
		if len(fr) == 2 {
			il, e = strconv.Atoi(fr[0])
			if e != nil {
				il = 1
				log.Error(e)
			}
			ir, e = strconv.Atoi(fr[1])
			if e != nil {
				ir = 1
				log.Error(e)
			}
		}
		if sa.FrameRate == 0 {
			sa.FrameRate = frameRateList[idx]
		}
		log.Info(sa.FrameRate, il, ir, il/ir)
		if sa.FrameRate > 0 {
			if sa.FrameRate > float64(il)/float64(ir) {
				sa.FrameRate = 0
			}
		}
	}
}

// FFMpegSplitToM3U8WithOptimize ...
func FFMpegSplitToM3U8WithOptimize(ctx context.Context, file string, args ...SplitOptions) (sa *SplitArgs, e error) {
	args = append(args, ProbeInfoOption(FFProbeStreamFormat))
	return FFMpegSplitToM3U8(ctx, file, args...)
}

// FFMpegSplitToM3U8 ...
func FFMpegSplitToM3U8(ctx context.Context, file string, args ...SplitOptions) (sa *SplitArgs, e error) {
	if strings.Index(file, " ") != -1 {
		return nil, xerrors.New("file name cannot have spaces")
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	sa = &SplitArgs{
		Output:          "",
		Auto:            true,
		Video:           "libx264",
		Audio:           "aac",
		M3U8:            "media.m3u8",
		SegmentFileName: "media-%05d.ts",
		HLSTime:         10,
	}
	for _, o := range args {
		o(sa)
	}

	if sa.probe != nil {
		sa.StreamFormat, e = sa.probe(file)
		if e != nil {
			return nil, e
		}
	}
	if sa.StreamFormat != nil {
		video := sa.StreamFormat.Video()
		audio := sa.StreamFormat.Audio()
		if !sa.StreamFormat.IsVideo() || audio == nil || video == nil {
			return nil, xerrors.New("open file failed with ffprobe")
		}

		//check scale before codec check
		optimizeScale(sa, video)

		if video.CodecName == "h264" && sa.Scale == 0 {
			sa.Video = "copy"
		}

		if audio.CodecName == "aac" {
			sa.Audio = "copy"
		}
	}

	sa.Output, e = filepath.Abs(sa.Output)
	if e != nil {
		return nil, e
	}
	if sa.Auto {
		sa.Output = filepath.Join(sa.Output, uuid.New().String())
		_ = os.MkdirAll(sa.Output, os.ModePerm)
	}

	sfn := filepath.Join(sa.Output, sa.SegmentFileName)
	m3u8 := filepath.Join(sa.Output, sa.M3U8)

	tpl := fmt.Sprintf(sliceM3u8FFmpegTemplate, file, sa.Video, sa.Audio, sa.HLSTime, sfn, m3u8)
	if sa.Scale != 0 {
		tpl = fmt.Sprintf(sliceM3u8ScaleTemplate, file, sa.Video, sa.Audio, outputScale(sa), sa.HLSTime, sfn, m3u8)
	}

	if err := FFMpegRun(ctx, tpl); err != nil {
		return nil, err
	}
	return sa, nil
}

// FFMpegRun ...
func FFMpegRun(ctx context.Context, args string) (e error) {
	ffmpeg := NewFFMpeg()
	ffmpeg.SetArgs(args)
	info := make(chan string, 1024)
	go func() {
		e = ffmpeg.RunContext(ctx, info)
	}()
	for v := range info {
		select {
		case <-ctx.Done():
			return
		default:
		}
		log.Infow("running", "proc", v)
	}
	return
}
