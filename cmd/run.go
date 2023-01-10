package cmd

import (
	"context"
	"fmt"
	"github.com/FRahimov84/throttler/pkg/logger"
	"github.com/FRahimov84/throttler/pkg/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		// Logger
		filename := "app.log"
		l := logger.InitLogger(filename)
		ctx := context.WithValue(context.Background(), "logger", l)

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

		err = pg.Pool.Ping(ctx)
		if err != nil {
			l.Fatal("err on ping database", zap.Error(err))
		}
		l = logger.FromCtx(ctx)
		l.Info("from ctx")

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
