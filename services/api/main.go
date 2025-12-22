package main

import (
	"api/internal/handlers"
	"context"
	"time"

	"common/infra/configuration"
	accountmongodb "common/infra/mongodb/billing/account"
	paymentmongodb "common/infra/mongodb/billing/payment"
	instancemongodb "common/infra/mongodb/instance"
	librarymongodb "common/infra/mongodb/library"
	listingmongodb "common/infra/mongodb/listing"
	usermongodb "common/infra/mongodb/user"

	listingmem "common/infra/inmem/listing"
	inmempayments "common/infra/inmem/payments"

	"common/pkg/app"
	"common/pkg/auth"
	"common/pkg/billing/account"
	"common/pkg/billing/payments/payment"
	"common/pkg/instance"
	"common/pkg/library"
	"common/pkg/listing"
	"common/pkg/money"
	"common/pkg/user"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DudReconciliationHandler is a no-op implementation of ReconciliationHandler
type DudReconciliationHandler struct{}

// OnInstanceCreated is called when a new instance is created.
// This implementation does nothing and always returns nil.
func (h *DudReconciliationHandler) OnInstanceCreated(instance *instance.Instance) error {
	// No-op
	return nil
}

// OnInstanceUpdated is called when an instance is updated.
// This implementation does nothing and always returns nil.
func (h *DudReconciliationHandler) OnInstanceUpdated(instance *instance.Instance) error {
	// No-op
	return nil
}

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
	userRepo, err := usermongodb.NewUserRepository(db.Collection("users"))
	if err != nil {
		fmt.Println(err)
		return
	}

	listingRepo := listingmongodb.NewListingRepository(db.Collection("listings"))

	libraryRepo, err := librarymongodb.NewSavedListingRepository(db.Collection("libraries"))
	if err != nil {
		fmt.Println(err)
		return
	}
	accountRepo := accountmongodb.NewBillingAccountRepository(db.Collection("billing_accounts"))
	paymentRepo := paymentmongodb.NewPaymentRepository(db.Collection("payments"), db.Collection("one_time_paymens"), db.Collection("subscription_payments"))
	instanceRepo := instancemongodb.NewInstanceRepository(db.Collection("instances"))
	//currencyRepo := currencymongodb.NewCurrencyRepository(db.Collection("currencies"))
	//hardwareRateRepo := ratesmongodb.NewHardwareCostRepository(db.Collection("hardware_costs"))

	userService := user.NewUserService(userRepo)
	userAuthentificationService := auth.NewAuthenticationService(userService)

	_, err = userService.CreateUser("alice", "$2y$10$lGdmMojygg80QG4DPE2xXeT9ByEJrJVa9JnEKRBDSAnxJzaDY9Hk2", true)
	if err != nil {
		fmt.Println(err)
	}

	listingService := listing.NewListingService(listingRepo)
	libraryService := library.NewLibraryService(libraryRepo)

	billingAccountService := account.NewBillingAccountService(accountRepo)
	paymentService := payment.NewPaymentService(paymentRepo)
	instanceService := instance.NewInstanceService(instanceRepo)
	//currencyService := currency.NewCurrencyService(currencyRepo)
	//hardwareCostService := rates.NewHardwareCostService(hardwareRateRepo)

	listingSearchService := listingmem.NewInMemoryListingSearchService()

	listings, err := listingRepo.ListAll()
	if err == nil {
		for _, listing := range listings {
			listingSearchService.IndexListing(listing)
		}
	}

	listingFacadeService := app.NewListingFacadeService(listingService, listingSearchService)
	publicListingService := app.NewPublicListingService(listingFacadeService)

	userAppService := app.NewUserAppService(
		userService,
		publicListingService, // could wrap listingService + search
		libraryService,
	)

	paymentGatewayResolver := inmempayments.NewInMemoryPaymentGatewayResolver()

	toplevelPaymentService := app.NewPaymentService(
		paymentService,
		billingAccountService,
		paymentGatewayResolver,
	)

	instanceEngineService := *app.NewInstanceEngineService(
		instanceService,
		billingAccountService,
		toplevelPaymentService,
		publicListingService,
		&DudReconciliationHandler{},
	)

	_, err = toplevelPaymentService.Pay("d49138ad-0b17-4299-b0ca-b2655fb2d208", money.Money{Amount: 1000, CurrencyCode: "USD"})
	if err != nil {
		fmt.Print(err)
	}

	app := &handlers.App{
		RunningConfiguration:  &runningConfiguration,
		UserAppService:        userAppService,
		UserAuth:              userAuthentificationService,
		ListingService:        publicListingService,
		BillingAccountService: billingAccountService,
		PaymentService:        paymentService,
		InstanceEngineService: &instanceEngineService,
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
	public.GET("/listing/:id", app.PublicListingHandler)
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

	// ================= Billing account endpoints =================
	auth.GET("/billing", app.UserBillingAccountsHandler)                    // list all accounts
	auth.POST("/billing", app.UserCreateBillingAccountHandler)              // create a new account
	auth.GET("/billing/:id", app.UserGetBillingAccountHandler)              // get single account
	auth.POST("/billing/:id/suspend", app.UserSuspendBillingAccountHandler) // suspend account
	auth.POST("/billing/:id/close", app.UserCloseBillingAccountHandler)     // close account

	auth.GET("/billing/:billingAccountId/payments", app.UserListPaymentsHandler)
	auth.GET("/billing/:billingAccountId/payment/:paymentId", app.UserGetPaymentHandler)
	auth.GET("/billing/:billingAccountId/payment/:paymentId/metadata", app.UserGetPaymentMetadataHandler)

	// ================= Instance endpoints =================
	auth.POST("/instances", app.UserLaunchInstanceHandler)                           // Launch a new instance
	auth.PATCH("/instances/:id/hardware", app.UserUpdateHardwareHandler)             // Update hardware spec
	auth.GET("/instances/:id", app.UserGetInstanceHandler)                           // Get single instance
	auth.GET("/billing/:billingId/instances", app.UserListInstancesByBillingHandler) // List instances by billing account
	auth.GET("/instances", app.UserListInstancesByOwnerHandler)                      // List all instances for the authenticated user

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
