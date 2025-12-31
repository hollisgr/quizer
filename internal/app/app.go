package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quizer_server/internal/app/services"
	"quizer_server/internal/config"
	"quizer_server/internal/db"
	"quizer_server/internal/handler"
	"quizer_server/internal/middleware"
	"quizer_server/internal/service/game"
	"quizer_server/internal/service/jwt"
	"quizer_server/internal/service/lobby"
	"quizer_server/internal/service/question"
	"quizer_server/internal/service/user"
	"quizer_server/pkg/postgres"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServices(pool *pgxpool.Pool) services.Services {
	storage := db.New(pool)
	us := user.New(storage)
	gs := game.New(storage)
	ls := lobby.New(storage)
	qs := question.New(storage)
	js := jwt.New(us)
	ua := middleware.NewUserAuthenticator(us, js)

	return services.Services{
		UserSvc:     us,
		JwtSvc:      js,
		UserAuth:    ua,
		GameSvc:     gs,
		LobbySvc:    ls,
		QuestionSvc: qs,
	}
}

// SetupRouter configures and returns a gin.Engine instance with registered wallet handlers.
func SetupRouter(s services.Services) *gin.Engine {
	r := gin.Default()
	h := handler.New(r, s)
	h.Register()
	return r
}

// SetupServer creates and returns an http.Server instance based on configuration and router.
func SetupServer(cfg *config.Config, r *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    cfg.Listen.Addr,
		Handler: r,
	}
	return srv
}

// StartServer starts the server concurrently and logs any fatal errors during its operation.
func StartServer(s *http.Server) {
	go func() {

		log.Println("Server is listening: ", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start err: %v", err)
		}
	}()
}

// HandleQuit gracefully shuts down the server when receiving SIGINT or SIGTERM signals.
func HandleQuit(s *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server shutdown err: %v", err)
	}
	log.Println("Application shutdown complete")
}

// ConnectToDB establishes a connection pool to PostgreSQL database using given configuration.
func ConnectToDB(cfg *config.Config) *pgxpool.Pool {
	pgxPool, err := postgres.NewPool(context.Background(), 5, cfg.Postgresql.DSN)
	if err != nil {
		log.Fatalln("cant connect to db")
	}
	log.Println("Connection to database OK")

	err = pgxPool.Ping(context.Background())
	if err != nil {
		log.Fatalln("cant ping to db, err:", err)
	}
	log.Println("Ping to database OK")
	return pgxPool
}
