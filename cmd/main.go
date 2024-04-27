package main

import (
	"MCatalogue/config"
	"MCatalogue/internal/api"
	"MCatalogue/internal/repository/postgresql"
	"MCatalogue/server"
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"
)

func main() {

	conf, _ := config.LoadConfig()
	dbInf := conf.DB
	pg, err := postgresql.New(dbInf)
	if err != nil {
		slog.Error("failed to init db", err)
	}

	//run migrations on start
	err = pg.RunMigrations()
	if err != nil {
		slog.Error("failed to run migrations", err)
		return
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(l)

	capturePC := log.Flags()&(log.Lshortfile|log.Llongfile) != 0
	log.SetOutput(&handlerWriter{l.Handler(), slog.LevelError, capturePC})

	handler := new(api.Handler)
	srv := new(server.Server)
	if err := srv.Run(conf.Port, handler.InitRoutes(pg)); err != nil {
		slog.Error("Error starting server: ", err.Error())
	}
}

type handlerWriter struct {
	h         slog.Handler
	level     slog.Level
	capturePC bool
}

func (w *handlerWriter) Write(buf []byte) (int, error) {
	if !w.h.Enabled(context.Background(), w.level) {
		return 0, nil
	}
	var pc uintptr
	if w.capturePC {
		var pcs [1]uintptr
		runtime.Callers(4, pcs[:])
		pc = pcs[0]
	}
	origLen := len(buf)
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	r := slog.NewRecord(time.Now(), w.level, string(buf), pc)
	return origLen, w.h.Handle(context.Background(), r)
}
