package inmemory

import (
	"context"
	"go/link_shortener/internal/storage"
)

type Storage struct {
	dbLink map[string]string
	dbShortLink map[string]string
}

func New() *Storage {
	return &Storage{
		dbLink: make(map[string]string),
		dbShortLink: make(map[string]string),
	}
}

func (s *Storage) CreateShortLink(ctx context.Context, link, shortLink string) (string, error) {
	if _, found := s.dbLink[link]; found{
		return "", storage.ErrLinkExists
	} else{
		s.dbLink[link] = shortLink
		s.dbShortLink[shortLink] = link
		return shortLink, nil
	}
}

func (s *Storage) GetLinkByShortLink(ctx context.Context, shortLink string) (string, error) {
	if val, found := s.dbShortLink[shortLink]; found{
		return val, nil
	} else{
		return "", storage.ErrLinkNotFound
	}
}

func (s *Storage) GetShortLinkByLink(ctx context.Context, link string) (string, error) {
	if val, found := s.dbLink[link]; found{
		return val, nil
	} else{
		return "", storage.ErrShortLinkNotFound
	}
}
