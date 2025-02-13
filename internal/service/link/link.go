package link

import (
	"context"
	"crypto/rand"
	"errors"
	"go/link_shortener/internal/storage"
	"math/big"
)

const letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const lengthUrl int = 10 

type Repository interface {
	CreateShortLink(ctx context.Context, link, shortLink string) (string, error)
	GetLinkByShortLink(ctx context.Context, shortLink string) (string, error)
	GetShortLinkByLink(ctx context.Context, link string) (string, error)
}

type Link struct{
	repo Repository
}

func New(repo Repository) *Link{
	return &Link{
		repo: repo,
	}
}

func (l *Link) ShortenLink (ctx context.Context, link string)  (string, error){
	shortLink, err := l.repo.GetShortLinkByLink(ctx, link)
	if errors.Is(err, storage.ErrShortLinkNotFound){
		for {
			shortLink, err = randomLink()
			if err != nil{
				return "", err
			}

			_, err := l.repo.GetLinkByShortLink(ctx, shortLink)
			if errors.Is(err, storage.ErrLinkNotFound){
				shortLink, err = l.repo.CreateShortLink(ctx, link, shortLink)
				if err != nil{
					return "", err
				}
				break
			}
		}
	}

	return shortLink, nil
}

func (l *Link) LengthenLink (ctx context.Context, shortLink string) (string, error){
	link, err := l.repo.GetLinkByShortLink(ctx, shortLink)
	if err != nil{
		return "", err
	}

	return link, nil
}

func randomLink() (string, error){
	b := make([]byte, lengthUrl)  
	for i := range b {  
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))  
		if err != nil {  
			return "", err  
		}  
		b[i] = letterRunes[num.Int64()]  
	}  
	return string(b), nil 
}