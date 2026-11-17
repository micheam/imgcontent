package main

import (
	"context"
	"io"
	"net/url"

	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ContentPath string

var regxNewline = regexp.MustCompile(`\r\n|\r|\n`) //throw panic if fail

type Filename struct {
	Value string
	Valid bool
}

func NewFilename(raw string) (fname *Filename, err error) {

	if raw == "" {
		return fname, ErrIllegalFilename("file name must not be empty")
	}

	rawName := filepath.Base(raw)

	if regxNewline.MatchString(rawName) {
		return fname, ErrIllegalFilename("file name must not contain new-line.")
	}

	name := strings.TrimSpace(rawName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "#", "_")

	return &Filename{Value: name, Valid: true}, err
}

type ErrIllegalFilename string

func (e ErrIllegalFilename) Error() string {
	return string(e)
}

type UploadUsecase struct {
	PathBuilder ContentPathBuilder
	Writer      ContentWriter
	Presenter   UploadResultPresenter
}

type UploadInput struct {
	Filename
	Reader io.Reader
}

type UploadOutput struct {
	Filename
	url.URL
}

func (u UploadUsecase) Handle(ctx context.Context, input UploadInput) error {

	path, err := u.PathBuilder.Build(ctx, input.Filename)
	if err != nil {
		return errors.Cause(err)
	}

	url, err := u.Writer.Write(ctx, input.Reader, path)
	if err != nil {
		return errors.Wrap(err, "failed to write content")
	}

	return u.Presenter.Complete(ctx, UploadOutput{
		Filename: input.Filename,
		URL:      url,
	})
}

type ContentPathBuilder interface {
	Build(ctx context.Context, filename Filename) (ContentPath, error)
}

type ContentWriter interface {
	Write(ctx context.Context, reader io.Reader, path ContentPath) (url.URL, error)
}

type UploadResultPresenter interface {
	Complete(ctx context.Context, data UploadOutput) error
}
