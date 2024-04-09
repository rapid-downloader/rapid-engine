package entry

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/rapid-downloader/rapid/log"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	Entry interface {
		ID() string
		Name() string
		SetName(name string)
		Location() string
		SetLocation(location string)
		Size() int64
		Type() string  // document, compressed, audio, video, image, other, etc
		URL() string   // url which the entry downloaded from
		ChunkLen() int // total chunks splitted into
		Resumable() bool
		Context() context.Context
		Cancel()
		Expired() bool
		Refresh() error
		Downloader() string
	}

	Headers map[string]string

	RequestClient interface {
		Request() *http.Request
	}

	entry struct {
		ctx               context.Context    `json:"-"`
		cancel            context.CancelFunc `json:"-"`
		request           *http.Request      `json:"-"`
		Id                string             `json:"id"`
		Name_             string             `json:"name"`
		Location_         string             `json:"location"`
		Size_             int64              `json:"size"`
		Filetype_         string             `json:"filetype"`
		URL_              string             `json:"url"`
		Resumable_        bool               `json:"resumable"`
		ChunkLen_         int                `json:"chunkLen"`
		DownloadProvider_ string             `json:"downloadProvider"`
	}

	option struct {
		setting          *setting.Setting
		cookies          []*http.Cookie
		headers          Headers
		downloadProvider string
	}

	Options func(o *option)
)

func UseSetting(setting *setting.Setting) Options {
	return func(o *option) {
		o.setting = setting
	}
}

func AddCookies(cookies []*http.Cookie) Options {
	return func(o *option) {
		o.cookies = cookies
	}
}

func AddHeaders(headers Headers) Options {
	return func(o *option) {
		o.headers = headers
	}
}

func UseDownloader(provider string) Options {
	return func(o *option) {
		o.downloadProvider = provider
	}
}

func id() string {
	return fmt.Sprint(time.Now().Unix())
}

func Fetch(url string, options ...Options) (Entry, error) {
	opt := &option{
		setting: setting.Get(),
	}

	for _, option := range options {
		option(opt)
	}

	log.Println("fetching url...")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("error preparing request:", err.Error())
		return nil, err
	}

	for _, cookie := range opt.cookies {
		req.AddCookie(cookie)
	}

	for key, value := range opt.headers {
		req.Header.Add(key, value)
	}

	// retry fetch 3x if error
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fetching url:", err.Error(), "Retrying....")
		return nil, err
	}

	resumable := resumable(res)
	filename := filepath.Base(handleDuplicate(filepath.Join(opt.setting.DownloadLocation, filename(res))))
	location := filepath.Join(opt.setting.DownloadLocation, filename)
	filetype := filetype(filename)
	ctx, cancel := context.WithCancel(context.Background())
	chunklen := calculatePartition(res.ContentLength, opt.setting)

	if !resumable {
		chunklen = 1
	}

	size := res.ContentLength
	if size == -1 {
		log.Println("downloading with unknown size...")
	}

	downloadProvider := "default"
	if opt.downloadProvider != "" {
		downloadProvider = opt.downloadProvider
	}

	entry := &entry{
		Id:                id(),
		Name_:             filename,
		Location_:         location,
		Filetype_:         filetype,
		URL_:              res.Request.URL.String(),
		Size_:             size,
		ChunkLen_:         chunklen,
		ctx:               ctx,
		cancel:            cancel,
		Resumable_:        resumable,
		request:           req,
		DownloadProvider_: downloadProvider,
	}

	return entry, nil
}

func (e *entry) ID() string {
	return e.Id
}

func (e *entry) Name() string {
	return e.Name_
}

func (e *entry) SetName(name string) {
	e.Name_ = name
}

func (e *entry) Location() string {
	return e.Location_
}

func (e *entry) SetLocation(location string) {
	e.Location_ = location
}

func (e *entry) Size() int64 {
	return e.Size_
}

func (e *entry) Type() string {
	return e.Filetype_
}

func (e *entry) URL() string {
	return e.URL_
}

func (e *entry) ChunkLen() int {
	return e.ChunkLen_
}

func (e *entry) Resumable() bool {
	return e.Resumable_
}

func (e *entry) Context() context.Context {
	return e.ctx
}

func (e *entry) Cancel() {
	e.cancel()
}

func (e *entry) Expired() bool {
	req := e.request.Clone(context.Background())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fetching expired status:", err.Error())
		return true
	}

	return res.StatusCode != http.StatusOK && res.ContentLength <= 0
}

func (e *entry) Refresh() error {
	e.ctx, e.cancel = context.WithCancel(context.Background())

	// TODO: do something else, such as refresh the link (future feature if browser extenstion is present)

	return nil
}

func (e *entry) Downloader() string {
	return e.DownloadProvider_
}

func (e *entry) Request() *http.Request {
	return e.request
}
