package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transactions"
	"bwastartup/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwa_golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	userRepo := user.NewRepository(db)
	campaignRepo := campaign.NewRepository(db)
	transactionRepo := transactions.NewRepository(db)
	userService := user.NewService(userRepo)
	campaingService := campaign.NewService(campaignRepo)
	paymentService := payment.NewService()
	transactionService := transactions.NewService(transactionRepo, campaignRepo, paymentService)
	userJwt := auth.NewSevice()
	userHandler := handler.NewUserHandler(userService, userJwt)
	campaignHandler := handler.NewCampaignHandler(campaingService)
	transactionHandler := handler.NewTTransactionHandler(transactionService)

	userService.SaveAvatar(14, "images/coba.jpg")
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "images")
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(userJwt, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(userJwt, userService), userHandler.FetchUser)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaignDetail)
	api.POST("/campaigns", authMiddleware(userJwt, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(userJwt, userService), campaignHandler.UpdatedCampaign)
	api.POST("/campaign-images", authMiddleware(userJwt, userService), campaignHandler.UploadImage)
	api.GET("/campaigns/:id/transactions", authMiddleware(userJwt, userService), transactionHandler.GetTransactionCampaign)
	api.GET("/transactions", authMiddleware(userJwt, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(userJwt, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		headerTokens := strings.Split(authHeader, " ")
		if len(headerTokens) == 2 {
			tokenString = headerTokens[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userId := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}
}
