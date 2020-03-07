package stream

import (
	"io/ioutil"
	"fmt"


	"errors"
	"os/exec"
	"strconv"
	"time"

	"github.com/stunndard/goicy/aac"
	"github.com/stunndard/goicy/mpeg"
	"github.com/stunndard/goicy/network"
)

var totalFramesSent uint64
var totalTimeBegin time.Time
var Abort bool

func (c *StreamClient) StreamFFMPEG(filename string) error {
	var res  error

	if c.cmd != nil {
		c.cmd.Process.Kill()
	}

	cleanUp := func(err error) {
		c.cmd.Process.Kill()
		totalFramesSent = 0
		res = err
	}

	cmdArgs := []string{}
	profile := ""
	if c.config.StreamFormat == "mpeg" {
		profile = "MPEG"
		cmdArgs = []string{
			"-i", filename,
			"-c:a", "libmp3lame",
			"-b:a", strconv.Itoa(c.config.StreamBitrate),
			"-cutoff", "20000",
			"-ar", strconv.Itoa(c.config.StreamSamplerate),
			//"-ac", strconv.Itoa(config.Cfg.StreamChannels),
			"-f", "mp3",
			"-write_xing", "0",
			"-id3v2_version", "0",
			"-loglevel", "fatal",
			"-",
		}
	} else {
		if c.config.StreamAACProfile == "lc" {
			profile = "aac_low"
		} else if c.config.StreamAACProfile == "he" {
			profile = "aac_he"
		} else {
			profile = "aac_he_v2"
		}
		cmdArgs = []string{
			"-i", filename,
			"-c:a", "libfdk_aac",
			"-profile:a", profile, //"aac_low", //
			"-b:a", strconv.Itoa(c.config.StreamBitrate),
			"-cutoff", "20000",
			"-ar", strconv.Itoa(c.config.StreamSamplerate),
			//"-ac", strconv.Itoa(config.Cfg.StreamChannels),
			"-f", "adts",
			"-loglevel", "fatal",
			"-",
		}
	}

	c.cmd = exec.Command(c.config.FFMPEGPath, cmdArgs...)

	f, _ := c.cmd.StdoutPipe()
	stderr, _ := c.cmd.StderrPipe()

	if err := c.cmd.Start(); err != nil {
		return err
	}

	zzz, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", zzz)

	frames := 0

	sr := c.config.StreamSamplerate
	spf := 0
	framesToRead := 1

	for {
		sendBegin := time.Now()

		var lbuf []byte
		var err error
		if c.config.StreamFormat == "mpeg" {
			lbuf, err = mpeg.GetFramesStdin(f, framesToRead)
			if framesToRead == 1 {
				if len(lbuf) < 4 {
					cleanUp(err)
					break
				}
				spf = mpeg.GetSPF(lbuf[0:4])
				framesToRead = (sr / spf) + 1
				mbuf, _ := mpeg.GetFramesStdin(f, framesToRead-1)
				lbuf = append(lbuf, mbuf...)
			}
		} else {
			lbuf, err = aac.GetFramesStdin(f, framesToRead)
			if framesToRead == 1 {
				if len(lbuf) < 7 {
					cleanUp(err)
					break
				}
				if c.config.StreamAACProfile != "lc" {
					sr = sr / 2
				}
				spf = aac.GetSPF(lbuf[0:7])
				framesToRead = (sr / spf) + 1
				mbuf, _ := aac.GetFramesStdin(f, framesToRead-1)
				lbuf = append(lbuf, mbuf...)
			}
		}

		if err != nil {
			cleanUp(err)
			break
		}

		if len(lbuf) <= 0 {
			break
		}

		if totalFramesSent == 0 {
			totalTimeBegin = time.Now()
			//stdoutFramesSent = 0
		}

		if err := network.Send(c.sock, lbuf); err != nil {
			cleanUp(err)
			break
		}

		totalFramesSent = totalFramesSent + uint64(framesToRead)
		frames = frames + framesToRead

		timeElapsed := int(float64((time.Now().Sub(totalTimeBegin)).Seconds()) * 1000)
		timeSent := int(float64(totalFramesSent) * float64(spf) / float64(sr) * 1000)
		//timeFileElapsed := int(float64((time.Now().Sub(timeFileBegin)).Seconds()) * 1000)

		bufferSent := 0
		if timeSent > timeElapsed {
			bufferSent = timeSent - timeElapsed
		}

		// calculate the send lag
		sendLag := int(float64((time.Now().Sub(sendBegin)).Seconds()) * 1000)

		// regulate sending rate
		timePause := 0
		if bufferSent < (c.config.BufferSize - 100) {
			timePause = 900 - sendLag
		} else {
			if bufferSent > c.config.BufferSize {
				timePause = 1100 - sendLag
			} else {
				timePause = 975 - sendLag
			}
		}

		if Abort {
			err := errors.New("Aborted by user")
			cleanUp(err)
			break
		}

		time.Sleep(time.Duration(time.Millisecond) * time.Duration(timePause))
	}
	c.cmd.Wait()

	return res
}
