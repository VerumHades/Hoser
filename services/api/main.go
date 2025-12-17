package main

import (
	"api/internal/handlers"
	"context"
	"time"

	"common/infra/configuration"
	librarymongodb "common/infra/mongodb/library"
	listingmongodb "common/infra/mongodb/listing"
	usermongodb "common/infra/mongodb/user"

	listingmem "common/infra/inmem/listing"

	"common/pkg/app"
	"common/pkg/auth"
	"common/pkg/library"
	"common/pkg/listing"
	"common/pkg/user"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// ---------------------------
	// Load configuration
	// ---------------------------
	runningConfiguration, err := configuration.Load[handlers.Configuration]()
	if err != nil {
		log.Fatal(err)
	}

	// Create MongoDB URI with authentication
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		runningConfiguration.MongoDatabaseUserName,
		runningConfiguration.MongoDatabasePassword,
		runningConfiguration.MongoDatabaseHost,
		runningConfiguration.MongoDatabasePort,
		runningConfiguration.MongoDatabaseName,
	)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Choose a database
	db := client.Database(runningConfiguration.MongoDatabaseName)

	// 3. Initialize collections & repositories
	userRepo := usermongodb.NewUserRepository(db.Collection("users"))
	listingRepo := listingmongodb.NewListingRepository(db.Collection("listings"))
	libraryRepo := librarymongodb.NewSavedListingRepository(db.Collection("libraries"))
	//accountRepo := accountmongodb.NewBillingAccountRepository(db.Collection("billing_accounts"))
	//paymentRepo := paymentmongodb.NewPaymentRepository(db.Collection("payments"), db.Collection("one_time_paymens"), db.Collection("subscription_payments"))
	//instanceRepo := instancemongodb.NewInstanceRepository(db.Collection("instances"))
	//currencyRepo := currencymongodb.NewCurrencyRepository(db.Collection("currencies"))
	//hardwareRateRepo := ratesmongodb.NewHardwareCostRepository(db.Collection("hardware_costs"))

	userAuthentificationService := auth.NewAuthenticationService(userRepo)

	userService := user.NewUserService(userRepo)

	_, err = userService.CreateUser("alice", "$2y$10$lGdmMojygg80QG4DPE2xXeT9ByEJrJVa9JnEKRBDSAnxJzaDY9Hk2", true)
	if err != nil {
		fmt.Println(err)
	}

	listingService := listing.NewListingService(listingRepo)
	libraryService := library.NewLibraryService(libraryRepo)
	//billingAccountService := account.NewBillingAccountService(accountRepo)
	//paymentService := payment.NewPaymentService(paymentRepo)
	//instanceService := instance.NewInstanceService(instanceRepo)
	//currencyService := currency.NewCurrencyService(currencyRepo)
	//hardwareCostService := rates.NewHardwareCostService(hardwareRateRepo)

	listingSearchService := listingmem.NewInMemoryListingSearchService()
	listingFacadeService := app.NewListingFacadeService(listingService, listingSearchService)

	publicListingService := app.NewPublicListingService(listingFacadeService)

	userAppService := app.NewUserAppService(
		userService,
		publicListingService, // could wrap listingService + search
		libraryService,
	)

	app := &handlers.App{
		RunningConfiguration: &runningConfiguration,
		UserAppService:       userAppService,
		UserAuth:             userAuthentificationService,
		ListingService:       publicListingService,
	}

	e := echo.New()

	// ---------------------------
	// Middleware
	// ---------------------------
	e.Use(app.CORSMiddleware())
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// ---------------------------
	// Public routes
	// ---------------------------
	public := e.Group("")
	public.GET("/listings", app.PublicListingsHandler)
	public.GET("/rentals/public/:id", app.PublicGetListingHandler)
	public.POST("/login", app.LoginHandler)
	public.POST("/logout", app.LogoutHandler)

	// ---------------------------
	// Authenticated user routes
	// ---------------------------
	auth := e.Group("/user")

	auth.GET("/data", app.UserDataHandler)
	auth.GET("/library", app.UserLibraryHandler)
	auth.POST("/library", app.UserAddListingToLibraryHandler)
	auth.DELETE("/library", app.UserRemoveListingFromLibraryHandler)

	// ---------------------------
	// Developer-only routes
	// ---------------------------
	e.GET("/developer/listings", app.DeveloperListingsHandler)
	dev := e.Group("/developer/listing")
	dev.Use(app.DeveloperOnlyMiddleware)

	dev.POST("", app.DeveloperAddListingHandler)
	dev.PUT("", app.DeveloperUpdateListingHandler)
	dev.DELETE("", app.DeveloperDeleteListingHandler)

	// ---------------------------
	// Start server
	// ---------------------------
	address := fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	e.Logger.Fatal(e.Start(address))
}
