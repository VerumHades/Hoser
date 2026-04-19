package main

import (
	"api/internal/adapters"
	"api/internal/http"
	"api/internal/http/authentification"
	developerapi "api/internal/http/developer"
	platformapi "api/internal/http/platform"
	userapi "api/internal/http/user"
	"api/internal/infrastructure"

	"common/pkg/application/outbox"
	"common/pkg/application/services/auth"
	"common/pkg/application/services/billing"
	"common/pkg/application/services/developer"
	"common/pkg/application/services/instanceservice"
	"common/pkg/util"

	userservices "common/pkg/application/services/user"

	"common/pkg/application/unitofwork"
	"common/pkg/domain/entities/accounting"
	"common/pkg/domain/entities/listing"
	"common/pkg/domain/entities/user"
	"common/pkg/domain/repositories"

	"common/pkg/infrastructure/external"
	mongodbinstance "common/pkg/infrastructure/mongodb/instance"
	mongodbledger "common/pkg/infrastructure/mongodb/ledger"
	mongodblisting "common/pkg/infrastructure/mongodb/listing"
	mongodboutbox "common/pkg/infrastructure/mongodb/outbox"
	mongodbrates "common/pkg/infrastructure/mongodb/rates"
	mongodbregistry "common/pkg/infrastructure/mongodb/registry"
	mongodbuser "common/pkg/infrastructure/mongodb/user"

	"common/pkg/infrastructure/plugs"
	"common/pkg/shared"
	"context"
	"encoding/json"
	"time"

	"fmt"
	"log"

	"common/pkg/infrastructure/configuration"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PrettyPrintJSON(
	value any,
) (string, error) {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func main() {
	// ---------------------------
	// Load configuration
	// ---------------------------
	runningConfiguration, err := configuration.Load[infrastructure.Configuration]()
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

	//transactionProvider := mongotransaction.NewMongoTransactionProvider(client)
	transactionProvider := plugs.NewFakeTransactionProvider()
	// Choose a database
	db := client.Database(runningConfiguration.MongoDatabaseName)

	databaseRegistry := mongodbregistry.NewDatabaseRegistry(db)

	userRepo := mongodbuser.NewMongoUserRepository(databaseRegistry)
	listingRepo := mongodblisting.NewMongoListingRepository(databaseRegistry)
	libraryRepo := mongodbuser.NewMongoSavedListingRepository(databaseRegistry)
	accountRepo := mongodbledger.NewMongoAccountRepository(databaseRegistry)
	ledgerTransactionRepo := mongodbledger.NewMongoLedgerTransactionRepository(databaseRegistry)
	settlementRepo := mongodbledger.NewMongoSettlementRepository(databaseRegistry)
	instanceRepo := mongodbinstance.NewMongoInstanceRentalContractRepository(databaseRegistry)
	hardwareRateRepo := mongodbrates.NewMongoHardwareCostRateRepository(databaseRegistry)
	githubSetupRepo := mongodblisting.NewMongoGitHubSetupRepository(databaseRegistry)
	savedListingViewRepo := mongodbuser.NewMongoUserSavedListingViewRepository(databaseRegistry)
	outboxRepo := mongodboutbox.NewMongoOutboxRepository(databaseRegistry)

	repos := []interface {
		EnsureIndexes(ctx context.Context) error
	}{
		userRepo,
		listingRepo,
		libraryRepo,
		accountRepo,
		ledgerTransactionRepo,
		settlementRepo,
		instanceRepo,
		hardwareRateRepo,
		githubSetupRepo,
	}
	for _, repo := range repos {
		if err = repo.EnsureIndexes(ctx); err != nil {
			fmt.Println(err)
		}
	}

	/*for i := range 200000 {
		accessMode := listing.Public
		if i%2 == 0 {
			accessMode = listing.Private
		}
		listing, _, _ := listing.NewListing(
			"30225718-7120-4353-999b-55a1ec8fcd4c",
			fmt.Sprintf("Test Listing %d %d", i, time.Now().Unix()),
			fmt.Sprintf("Test Description %d %d", i, time.Now().Unix()), accessMode, &shared.HardwareSpecification{}, 100*int64(i))

		listingRepo.Create(ctx, listing)
	}*/

	cuser, err := user.NewUser("alice", "$2y$10$lGdmMojygg80QG4DPE2xXeT9ByEJrJVa9JnEKRBDSAnxJzaDY9Hk2", true)
	_, err = userRepo.Create(ctx, cuser)
	if err != nil {
		fmt.Println(err)
	}

	userAuthService := auth.NewAuthenticationService(userRepo)

	listingIndexer := external.NewMeiliListingSearchIndex(
		runningConfiguration.MeilisearchHostname,
		runningConfiguration.MeilisearchApiKey,
	)

	{
		for listing := range util.GenerateInBatches(
			ctx,
			500,
			func(ctx context.Context, request shared.BatchRequest[repositories.ListingCursor]) ([]*listing.Listing, repositories.ListingCursor, error) {
				return listingRepo.FetchNextBatchAll(ctx, request)
			}) {
			listingIndexer.Index(ctx, listing)
		}
	}
	eventPublisher := outbox.NewOutboxEventPublisher(outboxRepo)

	transactionalEventPublisher := unitofwork.NewTransactionalEventPublisher(
		transactionProvider,
		eventPublisher,
	)

	// main.go
	storageAdapter, _ := adapters.NewMinioStorageAdapter(
		"localhost:9000",
		"localadmin",
		"localpassword",
		"listings",
		false,
	)

	developerListingService := developer.NewDeveloperListingService(
		*transactionalEventPublisher,
		storageAdapter,
		listingRepo,
		listingRepo,
		githubSetupRepo,
		listingIndexer,
	)

	userListingService := userservices.NewUserListingService(
		storageAdapter,
		listingRepo,
		listingRepo,
		libraryRepo,
		libraryRepo,
		listingIndexer,
	)

	accountAdapterConfig := adapters.PlatformAccountsConfig{
		ProfitAccountID:         "profit",
		HardwareRentAccountID:   "hardware_rent",
		HardwareRefundAccountID: "hardware_refund",
		PlatformCutPercentage:   10,
	}

	userAccountService := userservices.NewUserAccountService(accountRepo, accountRepo)
	accountAdapter := adapters.NewAccountServiceAdapter(accountAdapterConfig, listingRepo, userAccountService)
	hardwareCostCalculationService := billing.NewHardwareCostCalculationService(hardwareRateRepo)

	contractPaymentBuilder := instanceservice.NewContractPaymentBuilder(accountAdapter, hardwareCostCalculationService, ledgerTransactionRepo, ledgerTransactionRepo)
	contractService := instanceservice.NewInstanceContractService(
		contractPaymentBuilder,
		transactionalEventPublisher,
		listingRepo,
		instanceRepo,
		instanceRepo,
		userRepo,
	)

	userPurchaseService := userservices.NewUserPurchaseService(
		accountAdapter,
		transactionalEventPublisher,
		listingRepo,
		ledgerTransactionRepo,
		ledgerTransactionRepo,
		settlementRepo,
	)
	//githubSetupService := developer.NewDeveloperGitHubSetupService(githubSetupRepo, githubSetupRepo)
	//githubSetupService := listing.NewDeveloperGitHubSetupService(githubSetupRepo, githubSetupRepo)

	userPublicAPI := userapi.NewPublicUserAPI(userListingService)
	userProfileAPI := userapi.NewUserProfileAPI(userRepo)
	userLibraryAPI := userapi.NewUserLibraryAPI(
		adapters.NewUserSavedListingViewAdapter(
			savedListingViewRepo,
			libraryRepo,
		),
		userListingService,
	)
	userInstanceContractAPI := userapi.NewInstanceContractAPI(contractService)

	authAPI := authentification.NewUserAuthentificationAPI(
		authentification.UserAuthentificationAPIConfiguration{
			JWTSecret: runningConfiguration.JWTSecret,
		},
		userAuthService,
	)

	developerAPI := developerapi.NewDeveloperListingAPI(
		developerListingService,
		adapters.NewDeveloperCheckAdapter(userRepo),
	)

	purchaseAPI := userapi.NewUserPurchaseAPI(userPurchaseService)
	//developerSetupAPI := developerapi.NewDeveloperListingSetupAPI(githubSetupService)
	hardwareCostRatesAPI := platformapi.NewHardwareCostRatesAPI(hardwareRateRepo)

	{
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dispatcher := outbox.NewDispatcher(outboxRepo, outboxRepo, outbox.OutboxEventCursor{
			LastOccurredAt: time.Now(),
		}, time.Second*5)

		outbox.RegisterTypedListener(
			dispatcher,
			func(ctx context.Context, payload accounting.LedgerTransactionCreatedEvent, occuredAt time.Time) error {
				s, _ := PrettyPrintJSON(payload)
				fmt.Println(s)

				settlement, err := accounting.NewSettlement(payload.ID, shared.AccountID("a"), 100, "a")
				fmt.Println(err)
				settlement.MarkCompleted()
				settlementRepo.Create(ctx, settlement)

				return nil
			},
		)

		outbox.RegisterTypedListener(
			dispatcher,
			func(ctx context.Context, payload listing.ListingCreatedEvent, occuredAt time.Time) error {
				listingIndexer.Index(ctx, listing.ReconstituteListingFromEvent(payload, occuredAt))
				return nil
			},
		)

		outbox.RegisterTypedListener(
			dispatcher,
			func(ctx context.Context, payload listing.ListingUpdateEvent, occuredAt time.Time) error {
				listingIndexer.Update(ctx, &payload)
				return nil
			},
		)

		go func() {
			if err := dispatcher.Run(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}
	//developerapi.NewDeveloperListingSetupAPI(githubSetupService)

	e := echo.New()
	http.SetErrorHandler(e)

	// ---------------------------
	// Middleware
	// ---------------------------
	e.Use(http.CORSMiddleware(runningConfiguration.AllowedOrigins))
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	hardwareCostRatesAPI.RegisterRoutes(e.Group(""))

	authAPI.RegisterRoutes(e.Group(""))
	userPublicAPI.RegisterRoutes(e.Group(""))

	user := e.Group("/user")
	user.Use(authAPI.AuthenticationMiddleware)

	userProfileAPI.RegisterRoutes(user)
	userLibraryAPI.RegisterRoutes(user)
	userInstanceContractAPI.RegisterRoutes(user)
	purchaseAPI.RegisterRoutes(user)

	developer := e.Group("/developer")
	developer.Use(authAPI.AuthenticationMiddleware)
	developer.Use(developerAPI.DeveloperCheckMiddleware)
	developerAPI.RegisterRoutes(developer)
	//developerSetupAPI.RegisterRoutes(developer)

	// ---------------------------
	// Start server
	// ---------------------------
	address := fmt.Sprintf("%s:%s", runningConfiguration.Address, runningConfiguration.Port)
	fmt.Printf("Server running on http://%s\n", address)
	e.Logger.Fatal(e.Start(address))
}
