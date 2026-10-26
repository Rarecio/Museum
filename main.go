package main

import (
	"Museum/internal/auth"
	"Museum/internal/config"
	"Museum/internal/handler"
	"Museum/internal/repository"
	"Museum/internal/service"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
				}
			}
			return a
		},
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	slog.Debug("config", "config", cfg)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		slog.Error("database", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("database ping", err)
		os.Exit(1)
	}

	aiClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	userRepo := repository.NewUserRepository(pool)
	exhibitRepo := repository.NewExhibitRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)

	userSvc := service.NewUserService(userRepo, cfg.Secret)
	exhibitSvc := service.NewExhibitService(exhibitRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	reportSvc := service.NewReportService(exhibitRepo)
	chatSvc := service.NewAiService(exhibitRepo, aiClient, "http://26.86.186.191:4040")

	userHandler := handler.NewUserHandler(userSvc)
	exhibitHandler := handler.NewExhibitHandler(exhibitSvc, cfg.ImagesLocation)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	chatHandler := handler.NewAiHandler(chatSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

	r := gin.Default()
	r.Use(cors.Default())
	//r.Use(gin.Recovery())
	//r.Use(gin.Logger())

	config1 := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000", "https://myfrontend.com"}
	config1.AllowAllOrigins = true // Раскомментируйте, чтобы разрешить все

	r.Use(cors.New(config1))

	r.Static("/images", cfg.ImagesLocation)

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		api.GET("/exhibits", exhibitHandler.List)
		api.GET("/exhibits/:id", exhibitHandler.GetById)

		api.GET("/category", categoryHandler.List)

		api.GET("/text", chatHandler.HandleText)
		api.GET("/textai", chatHandler.AiResponse)
		api.POST("/recognize", chatHandler.RecognizeFromTextImage)

		api.GET("/report", reportHandler.Report)
		authGroup := api.Group("")
		authGroup.Use(auth.JWTAuth(cfg.Secret))
		{
			authGroup.GET("/user", userHandler.GetUserInfo)
			authGroup.PATCH("/user", userHandler.UpdateInfo)

			authGroup.PUT("/exhibits", exhibitHandler.Update)

			authGroup.POST("/exhibits", exhibitHandler.Create)
			authGroup.POST("/category", categoryHandler.Create)
		}
	}
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Info("shutdown", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
