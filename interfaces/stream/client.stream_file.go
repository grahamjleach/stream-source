package stream

import (
	"errors"
	"os"
	"time"

	"github.com/stunndard/goicy/aac"
	"github.com/stunndard/goicy/mpeg"
	"github.com/stunndard/goicy/network"
)

func (c *StreamClient) StreamFile(filename string) error {
	var (
		br                  float64
		spf, sr, frames, ch int
	)

	var err error
	if c.config.StreamFormat == "mpeg" {
		err = mpeg.GetFileInfo(filename, &br, &spf, &sr, &frames, &ch)
	} else {
		err = aac.GetFileInfo(filename, &br, &spf, &sr, &frames, &ch)
	}
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	if c.config.StreamFormat == "mpeg" {
		mpeg.SeekTo1StFrame(*f)
	} else {
		aac.SeekTo1StFrame(*f)
	}

	framesSent := 0

	// get number of frames to read in 1 iteration
	framesToRead := (sr / spf) + 1
	timeBegin := time.Now()

	for framesSent < frames {
		sendBegin := time.Now()

		var lbuf []byte
		if c.config.StreamFormat == "mpeg" {
			lbuf, err = mpeg.GetFrames(*f, framesToRead)
		} else {
			lbuf, err = aac.GetFrames(*f, framesToRead)
		}
		if err != nil {
			return err
		}

		if err := network.Send(c.sock, lbuf); err != nil {
			return err
		}

		framesSent = framesSent + framesToRead

		timeElapsed := int(float64((time.Now().Sub(timeBegin)).Seconds()) * 1000)
		timeSent := int(float64(framesSent) * float64(spf) / float64(sr) * 1000)

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
			err := errors.New("aborted by user")
			return err
		}

		time.Sleep(time.Duration(time.Millisecond) * time.Duration(timePause))
	}

	// pause to clear up the buffer
	timeBetweenTracks := int(((float64(frames)*float64(spf))/float64(sr))*1000) - int(float64((time.Now().Sub(timeBegin)).Seconds())*1000)
	time.Sleep(time.Duration(time.Millisecond) * time.Duration(timeBetweenTracks))

	return nil
}
