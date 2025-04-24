package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func ptr(s string) *string {
	return &s
}

type ffprobeOutput struct {
	Streams []struct {
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		CodecType string `json:"codec_type"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	var data bytes.Buffer
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmd.Stdout = &data
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var output ffprobeOutput
	if err := json.Unmarshal(data.Bytes(), &output); err != nil {
		return "", err
	}

	for _, stream := range output.Streams {
		if stream.CodecType == "video" {
			w, h := stream.Width, stream.Height
			if w == 0 || h == 0 {
				return "", fmt.Errorf("invalid dimensions")
			}

			ratio := float64(w) / float64(h)

			switch {
			case math.Abs(ratio-16.0/9.0) < 0.01:
				return "16:9", nil
			case math.Abs(ratio-9.0/16.0) < 0.01:
				return "9:16", nil
			default:
				return "other", nil
			}
		}
	}

	return "", fmt.Errorf("no video stream found")
}

func processVideoForFastStart(filePath string) (string, error) {
	// append ".processing" to the file path
	outputPath := filePath + ".processing"
	// create a new command
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)
	// run the command
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// return the output file path
	return outputPath, nil
}
