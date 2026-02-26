package router

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"red_packet/backend/internal/config"
	handlers "red_packet/backend/internal/http/handlers"
	appMiddleware "red_packet/backend/internal/http/middleware"
	"red_packet/backend/internal/service"
)

func New(svcs *service.Container, cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Printf("%s %s -> %d err=%v", v.Method, v.URI, v.Status, v.Error)
				return nil
			}
			log.Printf("%s %s -> %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api")
	authHandler := handlers.NewAuthHandler(svcs.Auth)
	configHandler := handlers.NewConfigHandler(svcs.Config)

	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/otp", authHandler.OTP)
	api.GET("/config/bootstrap", configHandler.Bootstrap)

	authGroup := api.Group("")
	authGroup.Use(appMiddleware.JWT(cfg))

	referralHandler := handlers.NewReferralHandler(svcs.Referral)
	rewardHandler := handlers.NewRewardHandler(svcs.Reward)
	taskHandler := handlers.NewTaskHandler(svcs.Task)
	lotteryHandler := handlers.NewLotteryHandler(svcs.Lottery)
	walletHandler := handlers.NewWalletHandler(svcs.Wallet)
	withdrawHandler := handlers.NewWithdrawHandler(svcs.Withdraw)
	adminHandler := handlers.NewAdminHandler(svcs.Withdraw, svcs.Task, svcs.Config, svcs.Risk, svcs.AdminOps)

	authGroup.POST("/referral/bind", referralHandler.Bind)
	authGroup.GET("/referral/status", referralHandler.Status)
	authGroup.GET("/reward/summary", rewardHandler.Summary)
	authGroup.GET("/reward/records", rewardHandler.Records)
	authGroup.POST("/reward/unlock", rewardHandler.Unlock)
	authGroup.GET("/task/list", taskHandler.List)
	authGroup.POST("/task/claim", taskHandler.Claim)
	authGroup.GET("/lottery/status", lotteryHandler.Status)
	authGroup.POST("/lottery/spin", lotteryHandler.Spin)
	authGroup.GET("/lottery/records", lotteryHandler.Records)
	authGroup.GET("/wallet", walletHandler.Get)
	authGroup.POST("/withdraw/apply", withdrawHandler.Apply)
	authGroup.GET("/withdraw/records", withdrawHandler.Records)

	adminGroup := api.Group("/admin")
	adminGroup.Use(appMiddleware.AdminKey(cfg))
	adminGroup.GET("/dashboard", adminHandler.Dashboard)
	adminGroup.GET("/task/list", adminHandler.ListTasks)
	adminGroup.POST("/task/save", adminHandler.SaveTask)
	adminGroup.DELETE("/task/:id", adminHandler.DeleteTask)
	adminGroup.GET("/config/list", adminHandler.ListConfigs)
	adminGroup.POST("/config/upsert", adminHandler.UpsertConfig)
	adminGroup.GET("/risk/flags", adminHandler.ListRiskFlags)
	adminGroup.POST("/risk/flag/add", adminHandler.AddRiskFlag)
	adminGroup.GET("/blacklist/list", adminHandler.ListBlacklists)
	adminGroup.POST("/blacklist/add", adminHandler.AddBlacklist)
	adminGroup.GET("/withdraw/list", adminHandler.ListWithdraw)
	adminGroup.POST("/withdraw/review", adminHandler.ReviewWithdraw)

	return e
}
