package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/app/api"
	"github.com/princecee/go_chat/app/websocket"
)

func StartApp(conn *pgxpool.Pool) {
	r := gin.Default()

	env := os.Getenv("ENVIRONMENT")
	port := os.Getenv("PORT")

	if env == "production" {
		gin.SetMode("release")
	}

	// setup websocket and api handlers
	websocket.SetupWebsocket(r, conn)
	api.SetupAPI(r, conn)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := http.Server{
		Handler: r.Handler(),
		Addr:    fmt.Sprintf(":%s", port),
	}

	errChan := make(chan error)
	log.Printf("server started on port :%s", port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-errChan:
	case <-ctx.Done():
		stop()
		log.Println("server shutting down in 5s")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server shut down successfully")
}
