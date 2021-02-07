package main

import (
	"encoding/json"
	"fmt"
	tool "github.com/glvd/go-media-tool"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	tool.DefaultProbeName = "D:\\workspace\\golang\\project\\go-fftool\\bin\\ffprobe.exe"
}

func formatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "format",
		Short: "get the video format for conversion",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("no source input")
				return
			}
			p := tool.NewFFProbe()
			fmt.Println("source", args[0])

			format, err := p.StreamFormat(args[0])
			if err != nil {
				fmt.Println("format error:", err)
				return
			}
			if err := writeToJSON(jsonName(args[0]), format); err != nil {
				fmt.Println("write error:", err)
				return
			}
		},
	}
}

func writeToJSON(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("open file error:%v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode error:%v", err)
	}
	return nil
}

func jsonName(path string) string {
	vol := filepath.VolumeName(path)
	i := len(path) - 1
	for i >= len(vol) && !os.IsPathSeparator(path[i]) {
		i--
	}
	name := path[i+1:]
	path = path[:i+1]
	for i := len(name) - 1; i >= 0 && !os.IsPathSeparator(name[i]); i-- {
		if name[i] == '.' {
			return filepath.Join(path, name[:i-1]+".JSON")
		}
	}
	return ""
}
