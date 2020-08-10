package tool

import (
	"encoding/json"
	"strings"
)

const ffprobeStreamFormatTemplate = `-v,quiet,-print_format,json,-show_format,-show_streams,%s`

// StreamFormat ...
type StreamFormat struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

// Format ...
type Format struct {
	Filename       string     `json:"filename"`
	NbStreams      int64      `json:"nb_streams"`
	NbPrograms     int64      `json:"nb_programs"`
	FormatName     string     `json:"format_name"`
	FormatLongName string     `json:"format_long_name"`
	StartTime      string     `json:"start_time"`
	Duration       string     `json:"duration"`
	Size           string     `json:"size"`
	BitRate        string     `json:"bit_rate"`
	ProbeScore     int64      `json:"probe_score"`
	Tags           FormatTags `json:"tags"`
}

// FormatTags ...
type FormatTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	Encoder          string `json:"encoder"`
}

// Stream ...
type Stream struct {
	Index              int64            `json:"index"`
	CodecName          string           `json:"codec_name"`
	CodecLongName      string           `json:"codec_long_name"`
	Profile            string           `json:"profile"`
	CodecType          string           `json:"codec_type"`
	CodecTimeBase      string           `json:"codec_time_base"`
	CodecTagString     string           `json:"codec_tag_string"`
	CodecTag           string           `json:"codec_tag"`
	Width              *int64           `json:"width,omitempty"`
	Height             *int64           `json:"height,omitempty"`
	CodedWidth         *int64           `json:"coded_width,omitempty"`
	CodedHeight        *int64           `json:"coded_height,omitempty"`
	HasBFrames         *int64           `json:"has_b_frames,omitempty"`
	SampleAspectRatio  *string          `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio *string          `json:"display_aspect_ratio,omitempty"`
	PixFmt             *string          `json:"pix_fmt,omitempty"`
	Level              *int64           `json:"level,omitempty"`
	ColorRange         *string          `json:"color_range,omitempty"`
	ColorSpace         *string          `json:"color_space,omitempty"`
	ColorTransfer      *string          `json:"color_transfer,omitempty"`
	ColorPrimaries     *string          `json:"color_primaries,omitempty"`
	ChromaLocation     *string          `json:"chroma_location,omitempty"`
	Refs               *int64           `json:"refs,omitempty"`
	IsAVC              *string          `json:"is_avc,omitempty"`
	NalLengthSize      *string          `json:"nal_length_size,omitempty"`
	RFrameRate         string           `json:"r_frame_rate"`
	AvgFrameRate       string           `json:"avg_frame_rate"`
	TimeBase           string           `json:"time_base"`
	StartPts           int64            `json:"start_pts"`
	StartTime          string           `json:"start_time"`
	DurationTs         int64            `json:"duration_ts"`
	Duration           string           `json:"duration"`
	BitRate            string           `json:"bit_rate"`
	BitsPerRawSample   *string          `json:"bits_per_raw_sample,omitempty"`
	NbFrames           string           `json:"nb_frames"`
	Disposition        map[string]int64 `json:"disposition"`
	Tags               StreamTags       `json:"tags"`
	SampleFmt          *string          `json:"sample_fmt,omitempty"`
	SampleRate         *string          `json:"sample_rate,omitempty"`
	Channels           *int64           `json:"channels,omitempty"`
	ChannelLayout      *string          `json:"channel_layout,omitempty"`
	BitsPerSample      *int64           `json:"bits_per_sample,omitempty"`
	MaxBitRate         *string          `json:"max_bit_rate,omitempty"`
}

// StreamTags ...
type StreamTags struct {
	Language    string `json:"language"`
	HandlerName string `json:"handler_name"`
}

// FileInfo ...
type FileInfo struct {
	Ext       string //扩展名
	Caption   string //字幕
	Language  string //语种
	Audio     string //音频
	Video     string //视频
	Sharpness string //清晰度
	Date      string //年份
	CName     string //中文名
	EName     string //英文名
	Prefix    string //前缀(广告信息)
}

// FFProbe ...
type FFProbe struct {
	cmd  *Command
	name string
}

// Name ...
func (ff *FFProbe) Name() string {
	return ff.name
}

// StreamFormat ...
func (ff *FFProbe) StreamFormat(file string) (*StreamFormat, error) {
	s, e := ff.cmd.Run(formatArgs(ffprobeStreamFormatTemplate, file))
	if e != nil {
		return nil, e
	}
	sf := StreamFormat{}
	e = json.Unmarshal([]byte(s), &sf)
	if e != nil {
		return nil, e
	}
	return &sf, nil
}

// Video ...
func (f *StreamFormat) Video() *Stream {
	for _, s := range f.Streams {
		if s.CodecType == "video" {
			return &s
		}
	}
	return nil
}

// IsVideo ...
func (f *StreamFormat) IsVideo() bool {
	return f.Video() != nil
}

// Audio ...
func (f *StreamFormat) Audio() *Stream {
	for _, s := range f.Streams {
		if s.CodecType == "audio" {
			return &s
		}
	}
	return nil
}

// ToString ...
func (info *FileInfo) ToString() string {
	var infos []string
	if info.Prefix != "" {
		infos = append(infos, info.Prefix)
	}
	infos = append(infos, info.CName)
	if info.EName != "" {
		infos = append(infos, info.EName)
	}
	infos = append(infos, info.Date)
	infos = append(infos, info.Sharpness)
	infos = append(infos, info.Video)
	infos = append(infos, info.Audio)
	infos = append(infos, info.Language)
	infos = append(infos, info.Caption)
	return strings.Join(infos, ".") + info.Ext
}

//var resolution = []int{120, 144, 160, 200, 240, 320, 360, 480, 540, 576, 600, 640, 720, 768, 800, 864, 900, 960, 1024, 1050, 1080, 1152, 1200, 1280, 1440, 1536, 1600, 1620, 1800, 1824, 1920, 2048, 2160, 2400, 2560, 2880, 3072, 3200, 4096, 4320, 4800}
var resolution = []int{240, 360, 480, 720, 1080, 1920, 2560, 4096, 4800}

func getResolution(n int64, sta, end int) int {
	size := len(resolution)
	if end == -1 || end > size {
		end = size
	}

	for {
		idx := (sta + end) / 2
		if idx > sta {
			if int64(resolution[idx]) > n {
				end = idx
			} else {
				sta = idx
			}
			continue
		}
		break
	}
	return resolution[sta]
}

// NewFFProbe ...
func NewFFProbe() *FFProbe {
	ff := &FFProbe{
		name: DefaultProbeName,
	}
	ff.cmd = NewCommand(ff.name)
	return ff
}
