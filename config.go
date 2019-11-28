package fftool

//const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -ss %s -to %s -c:v %s -c:a %s -bsf:v h264_mp4toannexb -vsync 0 -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const sliceM3u8FFmpegTemplate = `-y -i %s -strict -2 -c:v %s -c:a %s -bsf:v h264_mp4toannexb -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const sliceM3u8ScaleTemplate = `-y -i %s -strict -2 -c:v %s -c:a %s -bsf:v h264_mp4toannexb %s -f hls -hls_list_size 0 -hls_time %d -hls_segment_filename %s %s`
const scaleOutputTemplate = "-vf scale=-2:%d"
const bitRateOutputTemplate = "-b:v %dK"
const frameRateOutputTemplate = "-r %3.2f"

const defaultTemplate = `-y,-i,%s,-strict,-2,-c:v,%s,-c:a,%s,%s,%s`

// Scale ...
type Scale int

// Scale ...
const (
	Scale480P  Scale = 0
	Scale720P  Scale = 1
	Scale1080P Scale = 2
)

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

// Config ...
type Config struct {
	Scale     Scale
	BitRate   int64
	FrameRate float64
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		Scale: Scale720P,
	}
}

func (c *Config) init() {
	if c.BitRate == 0 {
		c.BitRate = bitRateList[c.Scale]
	}
	if c.FrameRate == 0 {
		c.FrameRate = frameRateList[c.Scale]
	}
}
