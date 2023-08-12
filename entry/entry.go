package entry

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	Entry interface {
		ID() string
		Name() string
		Location() string
		Size() int64
		Type() string  // document, compressed, audio, video, image, other, etc
		URL() string   // url which the entry downloaded from
		ChunkLen() int // total chunks splitted into
		Resumable() bool
		Context() context.Context
		Cancel()
		Expired() bool
		Refresh() error
	}

	EntryCookies interface {
		Cookies() []*http.Cookie
	}

	entry struct {
		id        string
		name      string
		location  string
		size      int64
		filetype  string
		url       string
		resumable bool
		chunkLen  int
		logger    logger.Logger
		ctx       context.Context
		cancel    context.CancelFunc
		cookies   []*http.Cookie
	}

	option struct {
		setting setting.Setting
		cookies []*http.Cookie
	}

	Options func(o *option)
)

func UseSetting(setting setting.Setting) Options {
	return func(o *option) {
		o.setting = setting
	}
}

func AddCookies(cookies []*http.Cookie) Options {
	return func(o *option) {
		o.cookies = cookies
	}
}

func Fetch(url string, options ...Options) (Entry, error) {
	opt := &option{
		setting: setting.Default(),
	}

	for _, option := range options {
		option(opt)
	}

	logger := logger.New(opt.setting)
	logger.Print("Fetching url...")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Print("Error preparing request:", err.Error())
		return nil, err
	}

	for _, cookie := range opt.cookies {
		req.AddCookie(cookie)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Print("Error fetching url:", err.Error())
	}

	resumable := resumable(res)
	filename := handleDuplicate(filename(res))
	location := filepath.Join(opt.setting.DownloadLocation(), filename)
	filetype := filetype(filename)
	ctx, cancel := context.WithCancel(context.Background())
	chunklen := calculatePartition(res.ContentLength, opt.setting)

	if !resumable {
		chunklen = 1
	}

	return &entry{
		id:        randID(10),
		name:      filename,
		location:  location,
		filetype:  filetype,
		url:       url,
		size:      res.ContentLength,
		logger:    logger,
		chunkLen:  chunklen,
		ctx:       ctx,
		cancel:    cancel,
		resumable: resumable,
		cookies:   opt.cookies,
	}, nil
}

func (e *entry) ID() string {
	return e.id
}

func (e *entry) Name() string {
	return e.name
}

func (e *entry) Location() string {
	return e.location
}

func (e *entry) Size() int64 {
	return e.size
}

func (e *entry) Type() string {
	return e.filetype
}

func (e *entry) URL() string {
	return e.url
}

func (e *entry) ChunkLen() int {
	return e.chunkLen
}

func (e *entry) Resumable() bool {
	return e.resumable
}

func (e *entry) Context() context.Context {
	return e.ctx
}

func (e *entry) Cancel() {
	e.cancel()
}

func (e *entry) Expired() bool {
	req, err := http.NewRequest("HEAD", e.url, nil)

	if len(e.cookies) > 0 {
		for _, cookie := range e.cookies {
			req.AddCookie(cookie)
		}
	}

	if err != nil {
		e.logger.Print("Could not prepare for checking url expiration:", err.Error())
		return true
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e.logger.Print("Error checking url expiration:", err.Error())
	}

	return res.StatusCode != http.StatusOK && res.ContentLength <= 0
}

func (e *entry) Refresh() error {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	// TODO: do something else, such as refresh the link (future feature if browser extenstion is present)

	return nil
}

func (e *entry) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("ID: %v\n", e.id))
	buffer.WriteString(fmt.Sprintf("Name: %v\n", e.name))
	buffer.WriteString(fmt.Sprintf("Location: %v\n", e.location))
	buffer.WriteString(fmt.Sprintf("Size: %v\n", e.size))
	buffer.WriteString(fmt.Sprintf("Filetype: %v\n", e.filetype))
	buffer.WriteString(fmt.Sprintf("URL: %v\n", e.url))
	buffer.WriteString(fmt.Sprintf("Resumable: %v\n", e.resumable))
	buffer.WriteString(fmt.Sprintf("ChunkLen: %v\n", e.chunkLen))
	buffer.WriteString(fmt.Sprintf("Expired: %v\n", e.Expired()))

	return buffer.String()
}

func (e *entry) Cookies() []*http.Cookie {
	return e.cookies
}
