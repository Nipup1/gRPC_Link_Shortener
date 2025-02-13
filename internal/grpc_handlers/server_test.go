package grpc_handlers_test

import (
	"context"
	"database/sql"
	"fmt"
	"go/link_shortener/internal/app"
	"go/link_shortener/internal/app/grpcapp"
	"go/link_shortener/internal/service/link"
	"go/link_shortener/internal/storage/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	lsv1 "github.com/Nipup1/link_shortener_gRPC/gen/go/link_shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createAppWithMockDB(port int) (*app.App, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	repository := &postgres.Storage{
		DB: database,
	}

	service := link.New(repository)

	grpcApp := grpcapp.New(port, service)

	return &app.App{
		GRPCSrv: grpcApp,
	}, mock, nil
}

func createAppWithInMemory(port int) *app.App{
	application := app.New(port, "", true)

	return application
}

func TestShortenDB(t *testing.T) {
	application, mock, err := createAppWithMockDB(44045)
	if err != nil {
		t.Fatal(err)
	}

	go application.GRPCSrv.MustRun()
	defer application.GRPCSrv.Stop()

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lsv1.NewLinkShortenerClient(conn)

	tests := []struct {
		name      string
		input     *lsv1.ShortenRequest
		expected  *lsv1.ShortenResponse
		expectErr bool
	}{
		{
			name:      "Valid Link",
			input:     &lsv1.ShortenRequest{Link: "http://google.com"},
			expected:  &lsv1.ShortenResponse{ShortLink: "dawAdawaA"},
			expectErr: false,
		},
		{
			name:      "Invalid Link",
			input:     &lsv1.ShortenRequest{Link: "awdawdawdawda"},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Link Already Exists",
			input:     &lsv1.ShortenRequest{Link: "http://google.com"},
			expected:  &lsv1.ShortenResponse{ShortLink: "dawAdawaA"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name{
			case "Valid Link":
				mock.ExpectPrepare("SELECT short_link FROM links WHERE link = \\$1").  
        			ExpectQuery().
					WithArgs(tt.input.Link).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectPrepare("SELECT link FROM links WHERE short_link = \\$1").
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectPrepare("INSERT INTO links\\(link, short_link\\) VALUES\\(\\$1, \\$2\\)").
        			ExpectExec().  
        			WithArgs(tt.input.Link, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			case "Invalid Link":

			case "Link Already Exists":
				mock.ExpectPrepare("SELECT short_link FROM links WHERE link = \\$1").  
        			ExpectQuery().
					WithArgs(tt.input.Link).
					WillReturnRows(sqlmock.NewRows([]string{"short_link"}).AddRow(tt.expected.ShortLink))
			}

			resp, err := client.Shorten(context.Background(), tt.input)

			if tt.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.expectErr && resp != nil {
				t.Errorf("Expected %v but got %v", tt.expected, resp)
			}
		})
	}
}

func TestLengthenDB(t *testing.T) {
	application, mock, err := createAppWithMockDB(44045)
	if err != nil {
		t.Fatal(err)
	}

	go application.GRPCSrv.MustRun()
	defer application.GRPCSrv.Stop()

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lsv1.NewLinkShortenerClient(conn)

	tests := []struct {
		name      string
		input     *lsv1.LengthenRequest
		expected  *lsv1.LengthenResponse
		expectErr bool
	}{
		{
			name:      "Link found",
			input:     &lsv1.LengthenRequest{ShortLink: "dawAdawaA"},
			expected:  &lsv1.LengthenResponse{Link: "http://google.com"},
			expectErr: false,
		},
		{
			name:      "Link not found",
			input:     &lsv1.LengthenRequest{ShortLink: "123ojbn12"},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name{
			case "Link found":
				mock.ExpectPrepare("SELECT link FROM links WHERE short_link = \\$1").
					ExpectQuery().
					WithArgs(tt.input.ShortLink).
					WillReturnRows(sqlmock.NewRows([]string{"link"}).AddRow(tt.expected.Link))

			case "Link not found":
				mock.ExpectPrepare("SELECT link FROM links WHERE short_link = \\$1").  
        			ExpectQuery().
					WithArgs(tt.input.ShortLink).
					WillReturnError(sql.ErrNoRows)
			}


			resp, err := client.Lengthen(context.Background(), tt.input)

			if tt.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.expectErr && resp != nil {
				t.Errorf("Expected %v but got %v", tt.expected, resp)
			}
		})
	}
}

func TestShortenIM(t *testing.T) {
	application := createAppWithInMemory(44045)

	go application.GRPCSrv.MustRun()
	defer application.GRPCSrv.Stop()

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lsv1.NewLinkShortenerClient(conn)

	tests := []struct {
		name      string
		input     *lsv1.ShortenRequest
		expected  *lsv1.ShortenResponse
		expectErr bool
	}{
		{
			name:      "Valid Link",
			input:     &lsv1.ShortenRequest{Link: "http://google.com"},
			expected:  &lsv1.ShortenResponse{ShortLink: "dawAdawaA"},
			expectErr: false,
		},
		{
			name:      "Invalid Link",
			input:     &lsv1.ShortenRequest{Link: "awdawdawdawda"},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Link Already Exists",
			input:     &lsv1.ShortenRequest{Link: "http://google.com"},
			expected:  &lsv1.ShortenResponse{ShortLink: "dawAdawaA"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Shorten(context.Background(), tt.input)
			fmt.Println(resp, err)

			if tt.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.expectErr && resp != nil {
				t.Errorf("Expected %v but got %v", tt.expected, resp)
			}
		})
	}
}

func TestLengthenIM(t *testing.T) {
	application := createAppWithInMemory(44045)

	go application.GRPCSrv.MustRun()
	defer application.GRPCSrv.Stop()

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lsv1.NewLinkShortenerClient(conn)

	tests := []struct {
		name      string
		input     *lsv1.LengthenRequest
		expected  *lsv1.LengthenResponse
		expectErr bool
	}{
		{
			name:      "Link found",
			input:     nil,
			expected:  &lsv1.LengthenResponse{Link: "http://google.com"},
			expectErr: false,
		},
		{
			name:      "Link not found",
			input:     &lsv1.LengthenRequest{ShortLink: "123ojbn12"},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Link found"{
				resp, _ := client.Shorten(context.Background(), &lsv1.ShortenRequest{Link: "http://google.com"})
				tt.input = &lsv1.LengthenRequest{ShortLink: resp.ShortLink}
			}

			resp, err := client.Lengthen(context.Background(), tt.input)

			if tt.expectErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.expectErr && resp != nil {
				t.Errorf("Expected %v but got %v", tt.expected, resp)
			}
		})
	}
}