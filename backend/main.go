package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xu-jiach/mecari-build-hackathon-2023/backend/db"
	"github.com/xu-jiach/mecari-build-hackathon-2023/backend/handler"
)

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	logfile := os.Getenv("LOGFILE")
	if logfile == "" {
		logfile = "access.log"
	}
	lf, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: logFormat(),
		Output: io.MultiWriter(os.Stdout, lf),
	})
	e.Use(logger)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{"GET", "PUT", "DELETE", "OPTIONS", "POST"},
	}))
	e.Use(middleware.BodyLimit("5M"))

	// jwt
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte(handler.GetSecret()),
	}

	// db
	sqlDB, err := db.PrepareDB(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare DB: %s\n", err)
		return exitError
	}
	defer sqlDB.Close()

	h := handler.Handler{
		DB:       sqlDB,
		UserRepo: db.NewUserRepository(sqlDB),
		ItemRepo: db.NewItemRepository(sqlDB),
	}

	// Routes
	e.POST("/initialize", h.Initialize)
	e.GET("/log", h.AccessLog)

	e.GET("/items", h.GetOnSaleItems)
	e.GET("/items/:itemID", h.GetItem)
	e.GET("/items/:itemID/image", h.GetImage)
	e.GET("/items/categories", h.GetCategories)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	// Login required
	l := e.Group("")
	l.Use(echojwt.WithConfig(config))
	l.GET("/users/:userID/items", h.GetUserItems)
	l.POST("/items", h.AddItem)
	l.POST("/sell", h.Sell)
	l.POST("/purchase/:itemID", h.Purchase)
	l.GET("/balance", h.GetBalance)
	l.POST("/balance", h.AddBalance)

	// Start server
	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return exitOK
}

func logFormat() string {
	// Customize freely: https://echo.labstack.com/guide/customization/
	var format string
	format += "time:${time_rfc3339}\t"
	format += "status:${status}\t"
	format += "method:${method}\t"
	format += "uri:${uri}\t"
	format += "latency:${latency_human}\t"
	format += "error:${error}\t"
	format += "\n"

	// Other log choice
	// - json format
	// `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
	// 	`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
	// 	`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
	// 	`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
	// - structured logging:  https://pkg.go.dev/golang.org/x/exp/slog

	return format
}
