```mermaid
graph TD
    %% Define Nodes/Components
    User([User / Developer / Public])
    
    subgraph Frontend_Service ["frontend Service (Go)"]
        ViteClient[Client App (Vite/TS/React)]
        AppRouter[Frontend Router/App.tsx]
        Components[Components (Tables, Input, Navigation)]
        ExploreView[Explore/Public Views (Listing Search/View)]
        UserView[User Views (Account, Instance Mgmt, Library)]
        DeveloperView[Developer Views (Listing Editor, Hardware, Prices)]
        BackendRepo[Backend Request Layer (Repositories)]
    end

    subgraph API_Service ["api Service (Go)"]
        HTTPHandlers[HTTP Handlers / Endpoints]
        AuthMiddleware[Auth Middleware]
        
        subgraph Internal_Adapters ["Internal Adapters / Use Cases"]
            AccountQuery[Account Query Adapter]
            MinioStorage[Minio Storage Adapter]
            UserLibAdapter[User Library Adapter]
        end

        subgraph Common_Packages ["common Package (Go Module)"]
            subgraph Domain_Layer ["Domain Layer (DDD)"]
                AccountingDomain[Entities: Account, Ledger, Settlement]
                ContractDomain[Entities: Instance Rental Contract]
                ListingDomain[Entities: Listing, ListingMutator, ListingSearchIndex]
                RatesDomain[Entities: HardwareRate]
                UserDomain[Entities: User, SavedListingQuery]
                EventDomain[Entities: Events]
                RepoInterfaces[Repository Interfaces]
            end

            subgraph Application_Layer ["Application Layer"]
                AuthService[Services: UserAuthentication]
                BillingService[Services: HardwareCostCalculation]
                DeveloperService[Services: DeveloperListing, GithubSetup]
                InstanceService[Services: InstanceContract]
                UserService[Services: UserAccount, UserListing, UserPurchase]
                OutboxPattern[Outbox: Dispatcher, EventPublisher]
                UOW[Unit of Work / Transactions]
            end

            subgraph Infrastructure_Layer ["Infrastructure Layer"]
                MongoRepos[MongoDB Repositories implementation]
                InMemIndex[In-Memory Listing Index]
                PlugsTx[Plugs: Transaction Handling]
            end
        end
    end

    %% External Systems / Persistence
    subgraph Storage ["Storage Layer"]
        MongoDB[(MongoDB)]
        Minio[(Minio - assumed for Screenshots/Assets)]
        GithubAPI[GitHub API]
    end

    %% Define Interconnections

    %% Frontend Interactions
    User -->|Interacts| ViteClient
    ViteClient --> AppRouter
    AppRouter --> ExploreView
    AppRouter --> UserView
    AppRouter --> DeveloperView
    ExploreView --> Components
    UserView --> Components
    DeveloperView --> Components
    Components --> BackendRepo

    %% Frontend-API Communication
    BackendRepo -->|HTTP/REST API| HTTPHandlers

    %% API Internal Flow
    HTTPHandlers --> AuthMiddleware
    AuthMiddleware --> AuthService
    HTTPHandlers --> Internal_Adapters

    %% Application Layer Flow
    Internal_Adapters --> Application_Layer
    Application_Layer --> Domain_Layer
    UOW -->|Manages| Application_Layer

    %% Domain to Infrastructure
    RepoInterfaces -->|Implemented by| Infrastructure_Layer

    %% Infrastructure to Storage
    Infrastructure_Layer -->|Data Persistence| MongoDB
    Application_Layer -->|Screenshots| MinioStorage
    MinioStorage -->|Object Storage| Minio
    DeveloperService -->|Setup Service| GithubAPI

    %% Styling
    classDef service fill:#e1f5fe,stroke:#01579b,stroke-width:2px;
    classDef storage fill:#fff9c4,stroke:#fbc02d,stroke-width:2px,stroke-dasharray: 5 5;
    classDef internal fill:#f1f8e9,stroke:#33691e,stroke-width:1px;
    classDef common fill:#f3e5f5,stroke:#4a148c,stroke-width:1px;
    classDef domain fill:#ffebee,stroke:#b71c1c,stroke-width:1px;
    classDef app fill:#e0f2f1,stroke:#004d40,stroke-width:1px;
    classDef infra fill:#efebe9,stroke:#3e2723,stroke-width:1px;

    class Frontend_Service,API_Service service;
    class Storage,MongoDB,Minio,GithubAPI storage;
    class Internal_Adapters internal;
    class Common_Packages common;
    class Domain_Layer domain;
    class Application_Layer app;
    class Infrastructure_Layer infra;
```