package mp3

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/dmulholland/mp3lib"
)

type (
	// Merger mp3合并器
	Merger struct {
		totalFrames  uint32
		totalBytes   uint32
		totalFiles   int
		firstBitRate int
		isVBR        bool
		out          *os.File
		done         bool
	}
)

// NewMerger 新建合并器
func NewMerger() *Merger {
	return &Merger{}
}

// Add 添加一个文件进合并文件
func (m *Merger) Add(in io.Reader) (err error) {
	if m.done {
		return errors.New("merge is done")
	}

	if m.out == nil {
		if m.out, err = ioutil.TempFile(os.TempDir(), "mp3merger_*.mp3"); err != nil {
			return
		}
	}

	isFirstFrame := true
	for {
		// Read the next frame from the input file.
		frame := mp3lib.NextFrame(in)
		if frame == nil {
			break
		}

		// Skip the first frame if it's a VBR header.
		if isFirstFrame {
			isFirstFrame = false
			if mp3lib.IsXingHeader(frame) || mp3lib.IsVbriHeader(frame) {
				continue
			}
		}

		// If we detect more than one bitrate we'll need to add a VBR
		// header to the output file.
		if m.firstBitRate == 0 {
			m.firstBitRate = frame.BitRate
		} else if frame.BitRate != m.firstBitRate {
			m.isVBR = true
		}

		// Write the frame to the output file.
		if _, err = m.out.Write(frame.RawBytes); err != nil {
			return
		}

		m.totalFrames++
		m.totalBytes += uint32(len(frame.RawBytes))
	}
	m.totalFiles++
	return
}

// Done 结束
func (m *Merger) Done() (r io.Reader, err error) {
	if m.done {
		err = errors.New("merge is done")
		return
	}
	m.done = true

	// 回到开头
	m.out.Seek(0, 0)

	readers := make([]io.Reader, 0)
	// If we detected multiple bitrates, prepend a VBR header to the file.
	if m.isVBR {
		readers = append(readers, m.getXingHeaderReader())
	}
	readers = append(readers, m.out)

	r = io.MultiReader(readers...)
	return
}

// Prepend an Xing VBR header to the specified MP3 file.
func (m *Merger) getXingHeaderReader() io.Reader {
	xingHeader := mp3lib.NewXingHeader(m.totalFrames, m.totalBytes)
	return bytes.NewReader(xingHeader.RawBytes)
}

// Close 闭关
func (m *Merger) Close() {
	m.done = true
	if m.out != nil {
		m.out.Close()
	}
}
