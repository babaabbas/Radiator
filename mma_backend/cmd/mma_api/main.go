package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"mma_api/internal/config"
	"mma_api/internal/http/handlers/auth"
	"mma_api/internal/http/handlers/product"
	"mma_api/internal/storage/postgres"
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
	user, err := pg.GetBoM(1)
	fmt.Println(user)
	if err != nil {
		fmt.Println(err)
	}

	//setup routers
	router := http.NewServeMux()
	router.HandleFunc("POST /api/register", auth.Register_handler(pg))
	router.HandleFunc("POST /api/login", auth.Login_handler(pg))
	router.HandleFunc("GET /api/users", auth.GetUsersHandler(pg))
	router.HandleFunc("GET /api/users/{id}", auth.GetUserByIDHandler(pg))
	router.HandleFunc("DELETE /api/users/{id}", auth.DeleteUserByIDHandler(pg))
	router.HandleFunc("PUT /api/users/{id}", auth.UpdateUserHandler(pg))
	router.HandleFunc("GET /api/products/", product.GetProductsHandler(pg))
	router.HandleFunc("GET /api/products/{id}", product.GetProductByIDHandler(pg))
	router.HandleFunc("POST /api/products/", product.CreateProductHandler(pg))
	router.HandleFunc("POST /api/products/{id}/bom", product.CreateBoMHandler(pg))
	router.HandleFunc("GET /api/products/{id}/bom", product.GetBoMHandler(pg))

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
