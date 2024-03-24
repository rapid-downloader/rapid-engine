package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rapid-downloader/rapid/client"
	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/log"
	"github.com/rapid-downloader/rapid/setting"
)

type progress struct {
	onprogress OnProgress
	reader     io.ReadCloser
	index      int
	chunkSize  int64
	prog       *client.Progress
}

func (r *progress) Read(payload []byte) (n int, err error) {
	n, err = r.reader.Read(payload)
	if err != nil {
		return n, err
	}

	r.prog.Chunks[r.index].Downloaded += int64(n)
	r.prog.Chunks[r.index].Progress = float64(100*r.prog.Chunks[r.index].Downloaded) / float64(r.prog.Chunks[r.index].Size)

	if r.onprogress != nil {
		r.onprogress(r.prog)
	}

	return n, err
}

func (r *progress) Close() error {
	return r.reader.Close()
}

type chunk struct {
	entry      entry.Entry
	setting    *setting.Setting
	wg         *sync.WaitGroup
	path       string
	index      int
	start      int64
	end        int64
	size       int64
	onprogress OnProgress
	prog       *client.Progress
}

func calculatePosition(entry entry.Entry, chunkSize int64, index int) (int64, int64) {
	if entry.Size() == -1 {
		return -1, -1
	}

	start := int64(index * int(chunkSize))
	end := start + (chunkSize - 1)

	if index == int(entry.ChunkLen())-1 {
		end = entry.Size()
	}

	return start, end
}

func resumePosition(location string) int64 {
	file, err := os.Stat(location)
	if err != nil {
		return 0
	}

	resumePos := file.Size()
	if err := os.Truncate(location, resumePos); err != nil {
		return 0
	}

	return resumePos
}

func newChunk(entry entry.Entry, index int, setting *setting.Setting, progress *client.Progress, wg *sync.WaitGroup) *chunk {
	chunkSize := entry.Size() / int64(entry.ChunkLen())
	start, end := calculatePosition(entry, chunkSize, index)

	progress.Chunks[index].Size = chunkSize

	return &chunk{
		path:       filepath.Join(setting.DownloadLocation, fmt.Sprintf("%s-%d", entry.ID(), index)),
		entry:      entry,
		setting:    setting,
		wg:         wg,
		index:      index,
		start:      start,
		end:        end,
		size:       chunkSize,
		prog:       progress,
		onprogress: nil,
	}
}

func (c *chunk) download(ctx context.Context) error {
	defer c.wg.Done()
	start := time.Now()

	srcFile, err := c.getDownloadFile(ctx, c.prog)
	if err != nil {
		log.Println("error fetching chunk file:", err.Error())
		return err
	}
	defer srcFile.Close()

	dstFile, err := c.getSaveFile()
	if err != nil {
		log.Println("error creating temp file for chunk:", err.Error())
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println("error downloading chunk:", err.Error())
		return err
	}

	c.prog.Chunks[c.index].Done = true
	c.prog.Chunks[c.index].Downloaded = c.size
	c.prog.Chunks[c.index].Progress = 100

	c.onprogress(c.prog)

	elapsed := time.Since(start)
	log.Println("chunk", c.index, "downloaded in", elapsed.Seconds(), "s")

	return nil
}

func (c *chunk) Execute(ctx context.Context) error {
	return c.download(ctx)
}

func (c *chunk) OnError(ctx context.Context, err error) {
	if c.entry.Context().Err() != nil {
		return
	}

	var e error
	for i := 0; i < c.setting.MaxRetry; i++ {
		c.wg.Add(1)
		log.Println("error downloading file:", err.Error(), ". Retrying...")

		if c.entry.Resumable() {
			c.start += resumePosition(c.path)
		}

		if e = c.download(ctx); e == nil {
			return
		}
	}

	log.Println("error downloading file:", err.Error())
}

func (c *chunk) onProgress(onprogress OnProgress) {
	c.onprogress = onprogress
}

func (c *chunk) getDownloadFile(ctx context.Context, prog *client.Progress) (io.ReadCloser, error) {
	req := c.entry.(entry.RequestClient).Request().Clone(ctx)

	if c.start != -1 && c.end != -1 {
		bytesRange := fmt.Sprintf("bytes=%d-%d", c.start, c.end)
		req.Header.Add("Range", bytesRange)

		log.Println("downloading chunk", c.index, "from", c.start, "to", c.end, fmt.Sprintf("(~%d MB)", (c.end-c.start)/(1024*1024)))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fething chunk body:", err.Error())
		return nil, err
	}

	progressBar := &progress{
		onprogress: c.onprogress,
		reader:     res.Body,
		index:      c.index,
		chunkSize:  c.size,
		prog:       prog,
	}

	return progressBar, nil
}

func (c *chunk) getSaveFile() (io.WriteCloser, error) {
	tmpFilename := filepath.Join(c.setting.DownloadLocation, fmt.Sprintf("%s-%d", c.entry.ID(), c.index))
	file, err := os.OpenFile(tmpFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error creating or appending file:", err.Error())
		return nil, err
	}

	return file, nil
}
