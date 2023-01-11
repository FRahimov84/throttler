package cmd

import (
	"context"
	"fmt"
	"github.com/FRahimov84/throttler/internal/infrastructure/external_service"
	repo "github.com/FRahimov84/throttler/internal/infrastructure/repo/postgres"
	"github.com/FRahimov84/throttler/internal/usecase"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	v1 "github.com/FRahimov84/throttler/internal/controller/http/v1"
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
		filename := "app.log"
		l := logger.InitLogger(filename)
		l.Info("application running...")
		// DB
		dbHost := viper.GetString(`database.host`)
		dbPort := viper.GetString(`database.port`)
		dbUser := viper.GetString(`database.user`)
		dbPass := viper.GetString(`database.pass`)
		dbName := viper.GetString(`database.name`)
		dbSSLMode := viper.GetString(`database.ssl_mode`)
		connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode)
		pg, err := postgres.New(connection, postgres.MaxPoolSize(viper.GetInt("database.pool_max")))
		if err != nil {
			l.Fatal("postgres new", zap.Error(err))
		}
		defer pg.Close()
		err = pg.Pool.Ping(context.Background())
		if err != nil {
			l.Fatal("err on ping database", zap.Error(err))
		}
		// Use case
		throttlerUseCase := usecase.New(
			repo.New(pg),
			external_service.New(
				viper.GetString("external_service.url"),
				viper.GetInt("external_service.n"),
				viper.GetInt("external_service.k"),
				viper.GetInt("external_service.x"),
			),
		)
		//HTTP Server
		handler := gin.New()
		v1.NewRouter(handler, throttlerUseCase, l)
		httpServer := httpserver.New(handler, httpserver.Port(viper.GetString("server.port")))
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
		err = httpServer.Shutdown()
		if err != nil {
			l.Error("httpServer.Shutdown", zap.Error(err))
		}

		l.Info("finished!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	viper.SetConfigFile(`./config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
