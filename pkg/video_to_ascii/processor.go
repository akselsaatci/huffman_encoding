package video_to_ascii

import (
	"fmt"
	"image/png"
	"io"
	"log"
	"os/exec"

	"github.com/akselsaatci/huffman/pkg/image_to_ascii"
)

type VideoToAsciiProcessor struct {
	videoPath string
	frameRate int
	resuliton string
	converter image_to_ascii.AsciiConverter
	stdIn     *io.ReadCloser
}

func NewVideoToFrameProcessor(v, r string, fps int, conv image_to_ascii.AsciiConverter, stdIn *io.ReadCloser) *VideoToAsciiProcessor {
	return &VideoToAsciiProcessor{
		videoPath: v,
		resuliton: r,
		frameRate: fps,
		converter: conv,
		stdIn:     stdIn,
	}
}
func (v *VideoToAsciiProcessor) Process() {
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", v.videoPath,
		"-vf", fmt.Sprintf("fps=%d", v.frameRate),
		"-vsync", "vfr",
		"-f", "image2pipe",
		"-vcodec", "png",
		"-s", v.resuliton,
		"-",
	)
	if v.stdIn != nil {
		ffmpegCmd.Stdin = *v.stdIn
	}

	stdout, err := ffmpegCmd.StdoutPipe()
	//	stdErr, err := ffmpegCmd.StderrPipe()

	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	if err := ffmpegCmd.Start(); err != nil {
		fmt.Println("Error starting FFmpeg command:", err)
		return
	}

	//	var wg sync.WaitGroup
	//
	//	// Goroutine to read and log stderr
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		scanner := bufio.NewScanner(stdErr)
	//		for scanner.Scan() {
	//			fmt.Println("FFmpeg stderr:", scanner.Text())
	//		}
	//		if err := scanner.Err(); err != nil {
	//			fmt.Println("Error reading stderr:", err)
	//		}
	//	}()
	//
	for {
		// Decode PNG frame from stdout
		img, err := png.Decode(stdout)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error decoding PNG:", err)
			break
		}

		if img != nil {
			//			fmt.Print("\033[H\033[2J")
			out, err := v.converter.Convert(img)
			if err != nil {
				log.Println(err.Error())
			}

			fmt.Print(out)

		}

	}

	if err := ffmpegCmd.Wait(); err != nil {
		fmt.Println("Error waiting for FFmpeg command:", err)
		return
	}

	fmt.Println("Frames extracted and processed.")

}