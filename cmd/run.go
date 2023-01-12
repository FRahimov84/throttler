package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/FRahimov84/throttler/config"
	v1 "github.com/FRahimov84/throttler/internal/controller/http/v1"
	"github.com/FRahimov84/throttler/internal/infrastructure/external_service"
	repo "github.com/FRahimov84/throttler/internal/infrastructure/repo/postgres"
	"github.com/FRahimov84/throttler/internal/infrastructure/tasks"
	"github.com/FRahimov84/throttler/internal/usecase"
	"github.com/FRahimov84/throttler/pkg/httpserver"
	"github.com/FRahimov84/throttler/pkg/logger"
	"github.com/FRahimov84/throttler/pkg/postgres"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().Unix())
		// Logger
		filename := "logs/app.log"
		l := logger.InitLogger(filename)
		l.Info("application running...")
		// Config
		cfg, err := config.LoadConfig("./config/config.json")
		if err != nil {
			l.Fatal("Err on load config", zap.Error(err))
		}
		fmt.Printf("%+v\n", cfg)
		// storage
		var (
			uRepo usecase.ThrottlerRepo
		)

		if cfg.EnableRedis {
			// TODO: implement redis
			fmt.Println("with redis mode")
			return
		} else {
			connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SslMode)
			if cfg.DB.URL != "" {
				connection = cfg.DB.URL
			}
			fmt.Println(connection)
			pg, err := postgres.New(connection, postgres.MaxPoolSize(cfg.DB.PoolMax))
			if err != nil {
				l.Fatal("postgres new", zap.Error(err))
			}
			defer pg.Close()
			err = pg.Pool.Ping(context.Background())
			if err != nil {
				l.Fatal("err to ping database", zap.Error(err))
			}
			uRepo = repo.New(pg)
		}
		// Use case
		throttlerUseCase := usecase.New(
			uRepo,
			external_service.New(cfg.ExternalSvc.Url),
		)
		// Tasks
		task := tasks.New(cfg.ExternalSvc.N, cfg.ExternalSvc.K, cfg.ExternalSvc.X)
		ctx, cancel := context.WithCancel(context.Background())
		ctx = logger.ToCtx(ctx, l)
		go task.Do(ctx, throttlerUseCase)
		//HTTP Server
		handler := gin.New()
		v1.NewRouter(handler, throttlerUseCase, l)
		httpServer := httpserver.New(handler, httpserver.Port(cfg.Server.Port))
		// Waiting signal
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		select {
		case s := <-interrupt:
			l.Info("got signal from OS", zap.Any("signal", s.String()))
		case err = <-httpServer.Notify():
			l.Error("Server notify", zap.Error(err))
		}
		// Shutdown
		l.Info("application stopping...")
		cancel()
		err = httpServer.Shutdown()
		if err != nil {
			l.Error("httpServer.Shutdown", zap.Error(err))
		}
		time.Sleep(3 * time.Second)
		l.Info("finished!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
