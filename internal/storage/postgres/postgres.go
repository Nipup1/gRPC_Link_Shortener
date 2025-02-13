package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go/link_shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct{
	DB *sql.DB
}

func New (storagePath string) (*Storage, error){
	db, err := sql.Open("postgres", storagePath)
	if err != nil{
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func (s *Storage) CreateShortLink(ctx context.Context, link, shortLink string) (string, error){
	stmt, err := s.DB.Prepare("INSERT INTO links(link, short_link) VALUES($1, $2)")
	if err != nil{
		return "", err
	}
	defer stmt.Close() 

	_, err = stmt.ExecContext(ctx, link, shortLink)
	if err!= nil{
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505"{  
				return "", storage.ErrLinkExists
		}

    	return "", err
	}

	return shortLink, nil
}

func (s *Storage) GetLinkByShortLink(ctx context.Context, shortLink string) (string, error){
	stmt, err := s.DB.Prepare("SELECT link FROM links WHERE short_link = $1")
	if err != nil{
		return "", err
	}
	defer stmt.Close() 

	row := stmt.QueryRowContext(ctx, shortLink)

	var link string
	if err = row.Scan(&link); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return "", storage.ErrLinkNotFound
		}

		return "", err
	}

	return link, nil
}

func (s *Storage) GetShortLinkByLink(ctx context.Context, link string) (string, error){
	stmt, err := s.DB.Prepare("SELECT short_link FROM links WHERE link = $1")
	if err != nil{
		return "", err
	}
	defer stmt.Close() 

	row := stmt.QueryRowContext(ctx, link)

	var shortLink string
	if err = row.Scan(&shortLink); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return "", storage.ErrShortLinkNotFound
		}
		
		return "", err
	}
	
	return shortLink, nil
}