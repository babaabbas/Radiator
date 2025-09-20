package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"mma_api/internal/config"
	"mma_api/internal/storage/postgres"
	"mma_api/internal/utils/response"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//config
	cfg := *config.Must_Load()
	pg, err := postgres.New(&cfg)
	if err != nil {
		fmt.Print("yo the postgres is not working", err)
	}
	pg.CreateUser("cool", "worker", "abs@gmail.com", "ols")
	fmt.Println(pg)

	//setup routers
	router := http.NewServeMux()

	router.HandleFunc("GET /api/man", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, this is a test response from /api/man")
		response.WriteJson(w, 200, "yoo")
	})
	//setup server
	server := http.Server{
		Addr:    cfg.Http_Server.Addr,
		Handler: router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("server is not starting")
		}
	}()
	<-done
	slog.Info("server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Info("server not closing ")
	}

}
