package storage

import "errors"

var (
	ErrLinkExists = errors.New("link already exists")
	ErrShortLinkNotFound = errors.New("short link not found")
	ErrLinkNotFound = errors.New("link not found")
)