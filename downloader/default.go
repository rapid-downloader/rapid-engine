package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/log"
	"github.com/rapid-downloader/rapid/setting"
	"github.com/rapid-downloader/rapid/worker"
)

// downloader that save the result into local file
type localDownloader struct {
	setting    *setting.Setting
	onprogress OnProgress
}

var Default = "default"

// TODO: possible change: rethink where the worker should be running (i.e on service init)
func newLocalDownloader(opt *option) Downloader {
	setting := setting.Get()
	if opt.setting != nil {
		setting = opt.setting
	}

	return &localDownloader{
		setting: setting,
	}
}

func (dl *localDownloader) Download(entry entry.Entry) error {
	start := time.Now()

	if entry.Expired() {
		return errUrlExpired
	}

	w, err := worker.New(entry.Context(), entry.ChunkLen(), entry.ChunkLen())
	if err != nil {
		log.Println("Error creating worker", err.Error())
		return err
	}

	var wg sync.WaitGroup
	w.Start()
	defer w.Stop()

	chunks := make([]*chunk, entry.ChunkLen())
	for i := 0; i < entry.ChunkLen(); i++ {
		chunks[i] = newChunk(entry, i, dl.setting, &wg)

		if dl.onprogress != nil {
			chunks[i].onProgress(dl.onprogress)
		}
	}

	for _, chunk := range chunks {
		wg.Add(1)
		w.Add(chunk)
	}

	wg.Wait()

	if entry.Context().Err() != nil {
		return nil
	}

	if err := dl.createFile(entry, dl.setting); err != nil {
		log.Println("Error combining chunks:", err.Error())
		return err
	}

	elapsed := time.Since(start)
	log.Println(entry.Name(), "downloaded  in", elapsed.Seconds(), "s")

	return nil
}

var errUrlExpired = fmt.Errorf("link is expired")

func (dl *localDownloader) Resume(entry entry.Entry) error {
	start := time.Now()

	if entry.Expired() {
		return errUrlExpired
	}

	if err := entry.Refresh(); err != nil {
		return err
	}

	log.Println("Resuming download", entry.Name(), "...")

	if !entry.Resumable() {
		log.Println(entry.Name(), "does not support resume download. Restarting...")
		return dl.Download(entry)
	}

	worker, err := worker.New(entry.Context(), entry.ChunkLen(), entry.ChunkLen())
	if err != nil {
		log.Println("Error creating worker", err.Error())
		return err
	}

	var wg sync.WaitGroup
	worker.Start()
	defer worker.Stop()

	chunks := make([]*chunk, 0)
	for i := 0; i < entry.ChunkLen(); i++ {
		chunk := newChunk(entry, i, dl.setting, &wg)
		if file, err := os.Stat(chunk.path); err == nil && file.Size() == chunk.size {
			continue
		}

		chunk.start += resumePosition(chunk.path)
		if dl.onprogress != nil {
			chunk.onProgress(dl.onprogress)
		}

		chunks = append(chunks, chunk)
	}

	for _, chunk := range chunks {
		wg.Add(1)
		worker.Add(chunk)
	}

	wg.Wait()

	if err := dl.createFile(entry, dl.setting); err != nil {
		log.Println("Error combining chunks:", err.Error())
		return err
	}

	elapsed := time.Since(start)
	log.Println(entry.Name(), "resumed in", elapsed.Seconds(), "s")

	return nil
}

func (dl *localDownloader) Restart(entry entry.Entry) error {
	log.Println("Restarting download", entry.Name(), "...")

	if entry.Expired() {
		return errUrlExpired
	}

	if err := entry.Refresh(); err != nil {
		return err
	}

	// remove the downloaded chunk if any
	for i := 0; i < entry.ChunkLen(); i++ {
		chunkFile := filepath.Join(dl.setting.DownloadLocation, fmt.Sprintf("%s-%d", entry.ID(), i))
		if err := os.Remove(chunkFile); err != nil {
			return err
		}
	}

	return dl.Download(entry)
}

func (dl *localDownloader) Stop(entry entry.Entry) error {
	log.Println("Stopping download", entry.Name(), "...")

	entry.Cancel()
	return nil
}

// Watch will update the id, index, downloaded bytes, and progress in percent of chunks. Watch must be called before Download
func (dl *localDownloader) Watch(update OnProgress) {
	dl.onprogress = update
}

// createFile will combine chunks into single actual file
func (dl *localDownloader) createFile(entry entry.Entry, s *setting.Setting) error {
	// if chunk len is 1, then just rename the chunk into entry filename
	if entry.ChunkLen() == 1 {
		chunkname := filepath.Join(s.DownloadLocation, fmt.Sprintf("%s-%d", entry.ID(), 0))
		return os.Rename(chunkname, entry.Location())
	}

	file, err := os.Create(entry.Location())
	if err != nil {
		log.Println("Error creating downloaded file:", err.Error())
		return err
	}

	defer file.Close()

	for i := 0; i < entry.ChunkLen(); i++ {
		tmpFilename := filepath.Join(s.DownloadLocation, fmt.Sprintf("%s-%d", entry.ID(), i))
		if err := dl.appendChunk(file, tmpFilename); err != nil {
			return err
		}
	}

	return nil
}

func (dl *localDownloader) appendChunk(dst io.Writer, srcName string) error {
	tmpFile, err := os.Open(srcName)
	if err != nil {
		log.Println("Error opening downloaded chunk file:", err.Error())
		return err
	}

	defer tmpFile.Close()

	if _, err := io.Copy(dst, tmpFile); err != nil {
		log.Println("Error copying chunk file into actual file:", err.Error())
		return err
	}

	return os.Remove(srcName)
}

func init() {
	registerDownloader(Default, newLocalDownloader)
}
