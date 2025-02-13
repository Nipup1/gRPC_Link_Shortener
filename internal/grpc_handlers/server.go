package grpc_handlers

import (
	"context"
	"errors"
	"go/link_shortener/internal/storage"

	lsv1 "github.com/Nipup1/link_shortener_gRPC/gen/go/link_shortener"
	"github.com/go-playground/validator/v10"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LinkShortener interface{
	ShortenLink (ctx context.Context, link string)  (string, error)
	LengthenLink (ctx context.Context, shortLink string) (string, error)
}

type serverAPI struct{
	lsv1.UnimplementedLinkShortenerServer
	linkShortener LinkShortener
}

func Register(gRPC *grpc.Server, linkShortener LinkShortener){
	lsv1.RegisterLinkShortenerServer(gRPC, &serverAPI{linkShortener: linkShortener})
}

func (s *serverAPI) Shorten (ctx context.Context, req *lsv1.ShortenRequest) (*lsv1.ShortenResponse, error){
	if err := validateEmail(req.Link); err != nil{
		return nil, status.Error(codes.InvalidArgument, "invalid link")
	}
	
	shortLink, err := s.linkShortener.ShortenLink(ctx, req.GetLink())
	if err != nil{
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &lsv1.ShortenResponse{
		ShortLink: shortLink,
	}, nil
}

func (s *serverAPI) Lengthen (ctx context.Context, req *lsv1.LengthenRequest) (*lsv1.LengthenResponse, error){
	link, err := s.linkShortener.LengthenLink(ctx, req.GetShortLink())
	if err != nil{
		if errors.Is(err, storage.ErrLinkNotFound){
			return nil, status.Error(codes.NotFound, "short link not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &lsv1.LengthenResponse{
		Link: link,
	}, nil 
}

func validateEmail(link string) error{
	validate := validator.New() 
	if err := validate.Var(link, "required,url");  err != nil{
		return err
	}

	return nil
}