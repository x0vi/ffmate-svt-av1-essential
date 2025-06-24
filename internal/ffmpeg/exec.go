package ffmpeg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/mattn/go-shellwords"
	"github.com/welovemedia/ffmate/internal/config"
	"github.com/yosev/debugo"
)

var debug = debugo.New("ffmpeg")

// ExecuteFFmpeg runs the ffmpeg command, provides progress updates, and checks the result
func Execute(request *ExecutionRequest) error {
	var args []string
	var err error
	if runtime.GOOS == "windows" {
		args, err = shellwordsUnicodeSafe(request.Command)
	} else {
		args, err = shellwords.NewParser().Parse(request.Command)
	}
	if err != nil {
		return fmt.Errorf("FFMPEG - failed to parse command: %v", err)
	}
	args = append(args, "-progress", "pipe:2")
	config.Config().Mutex.RLock()
	cmd := exec.CommandContext(request.Ctx, config.Config().FFMpeg, args...)
	config.Config().Mutex.RUnlock()

	// Buffers for capturing full stderr
	var stderrBuf bytes.Buffer
	var lastLine string
	var duration float64

	// Stderr pipe for real-time progress parsing
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("FFMPEG - failed to get stderr pipe: %v", err)
	}

	// Start the ffmpeg process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("FFMPEG - failed to start ffmpeg: %v", err)
	}

	// Regex to extract the duration field
	reDuration := regexp.MustCompile(`Duration: (\d+:\d+:\d+\.\d+)`)

	// Parse progress in real-time
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stderrBuf.WriteString(line + "\n")
			lastLine = line
			if match := reDuration.FindStringSubmatch(line); match != nil {
				durationStr := match[1]
				duration = parseDuration(durationStr)
			}
			if progress := parseFFmpegOutput(line, duration); progress != nil {
				p := math.Min(100, math.Round((progress.Time/duration*100)*100)/100) // cap at 100
				debug.Debugf("progress: %f %+v (uuid: %s)", p, progress, request.Task.Uuid)

				// Calculate and log the estimated remaining time
				remainingTime, err := progress.EstimateRemainingTime(duration)
				if err != nil {
					debug.Debugf("failed to estimate remaining time: %v", err)
					remainingTime = -1
				}

				request.UpdateFunc(p, remainingTime)
			}
		}
		if err := scanner.Err(); err != nil {
			request.Logger.Warnf("FFMPEG - error reading progress: %v\n", err)
		}
	}()

	// Wait for the ffmpeg process to complete
	err = cmd.Wait()

	// Gather full output for final reporting
	stderr := stderrBuf.String()

	if err != nil {
		return errors.New(stderr)
	}

	fmt.Sprintf("last line: %s", lastLine)

	return nil
}
