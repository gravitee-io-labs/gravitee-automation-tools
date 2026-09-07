package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am-mock-server/server"
	"github.com/spf13/cobra"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := newCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	var port int
	var basePath string

	cmd := &cobra.Command{
		Use:          "am-mock-server",
		Short:        "Start the AM automation mock server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return serve(cmd.Context(), port, basePath)
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "HTTP listen port")
	cmd.Flags().StringVar(&basePath, "base-path", server.BasePath, "API base path")
	cmd.CompletionOptions.DisableDefaultCmd = true

	return cmd
}

func serve(ctx context.Context, port int, basePath string) error {
	handler := server.NewWithPath(server.NewMockAM(context.Background()), basePath)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("AM mock server listening on http://localhost:%d%s", port, basePath)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		return srv.Shutdown(context.Background())
	}
}
