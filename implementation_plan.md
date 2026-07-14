# Rancangan Migrasi NeoCentral → Go Microservices (v2 - Final)

## Keputusan Arsitektur

| Keputusan | Pilihan |
|---|---|
| ORM | **GORM** |
| HTTP Framework | **Chi** |
| Message Broker | **NATS (JetStream)** |
| Database | **Database per Service** |
| File Storage | **MinIO / S3** |
| Internal Communication | **gRPC** |
| External API | **REST/JSON** (via API Gateway) |
| Fokus Fase 1 | `auth-service`, `master-data-service`, `document-service`, `notification-service` |

---

## 1. Arsitektur Final

```mermaid
graph TB
    subgraph Clients
        FE["Frontend<br/>(React/Next.js)"]
        MOB["Mobile App"]
    end

    subgraph "API Gateway (Chi - REST)"
        GW["api-gateway<br/>:8080<br/>JWT validation, routing,<br/>rate limit, CORS"]
    end

    subgraph "Microservices (Fase 1)"
        AUTH["auth-service<br/>:8001 REST / :9001 gRPC"]
        MASTER["master-data-service<br/>:8002 REST / :9002 gRPC"]
        DOC["document-service<br/>:8009 REST / :9009 gRPC"]
        NOTIF["notification-service<br/>:8008 REST / :9008 gRPC"]
    end

    subgraph "Microservices (Fase 2 - nanti)"
        THESIS["thesis-service"]
        GUIDANCE["guidance-service"]
        EXAM["examination-service"]
        INTERN["internship-service"]
        YUDISIUM["yudisium-service"]
    end

    subgraph "Infrastructure"
        DB_AUTH[(MySQL<br/>neocentral_auth)]
        DB_MASTER[(MySQL<br/>neocentral_master)]
        DB_DOC[(MySQL<br/>neocentral_doc)]
        DB_NOTIF[(MySQL<br/>neocentral_notif)]
        REDIS[(Redis)]
        NATS_SRV["NATS JetStream"]
        MINIO["MinIO<br/>(S3-Compatible)"]
    end

    FE -->|REST/JSON| GW
    MOB -->|REST/JSON| GW

    GW -->|REST| AUTH
    GW -->|REST| MASTER
    GW -->|REST| DOC
    GW -->|REST| NOTIF

    AUTH <-->|gRPC| MASTER
    AUTH <-->|gRPC| NOTIF
    DOC <-->|gRPC| AUTH
    NOTIF <-->|gRPC| AUTH

    AUTH --> DB_AUTH
    AUTH --> REDIS
    MASTER --> DB_MASTER
    MASTER --> REDIS
    DOC --> DB_DOC
    DOC --> MINIO
    NOTIF --> DB_NOTIF
    NOTIF --> NATS_SRV

    AUTH -->|publish events| NATS_SRV
    MASTER -->|publish events| NATS_SRV
    DOC -->|publish events| NATS_SRV
    NOTIF -->|subscribe events| NATS_SRV
```

### Pola Komunikasi

```mermaid
graph LR
    subgraph "External (Client → Gateway → Service)"
        A["REST/JSON"]
    end
    subgraph "Internal (Service ↔ Service)"
        B["gRPC<br/>(synchronous queries)"]
        C["NATS JetStream<br/>(async events)"]
    end

    A --- B
    B --- C
```

| Pola | Kapan Digunakan | Contoh |
|---|---|---|
| **REST** | Client → API Gateway → Service | Login, CRUD operations |
| **gRPC** | Service perlu data dari service lain **saat handle request** | `document-service` perlu validasi user dari `auth-service` |
| **NATS** | Fire-and-forget events, async processing | User registered → kirim email welcome |

---

## 2. Database per Service

### 2.1 Schema Ownership

Setiap service memiliki **database MySQL sendiri** dan **hanya service owner** yang boleh read/write ke database tersebut.

```
┌─────────────────────────────────────────────────────────────────┐
│                    MySQL Server Instance                         │
├──────────────────┬──────────────────┬──────────────┬────────────┤
│ neocentral_auth  │ neocentral_master│ neocentral_doc│neocentral_ │
│                  │                  │              │ notif      │
│ • users          │ • academic_years │ • documents  │• notifica- │
│ • user_roles     │ • rooms          │ • document_  │  tions     │
│ • user_has_roles │ • science_groups │   types      │• push_     │
│ • students       │ • cpls           │ • document_  │  tokens    │
│ • lecturers      │ • cpmks          │   templates  │            │
│                  │ • assessment_    │              │            │
│                  │   criterias      │              │            │
│                  │ • assessment_    │              │            │
│                  │   rubrics        │              │            │
│                  │ • lecturer_      │              │            │
│                  │   availabilities │              │            │
└──────────────────┴──────────────────┴──────────────┴────────────┘
```

### 2.2 Cross-Service Data Access

Karena database terpisah, service **tidak boleh JOIN** antar database. Gunakan:

| Kebutuhan | Solusi |
|---|---|
| `document-service` perlu nama user | gRPC call ke `auth-service.GetUser(userId)` |
| `notification-service` perlu email user | gRPC call ke `auth-service.GetUserContact(userId)` |
| `master-data-service` perlu cek user role | gRPC call ke `auth-service.ValidateRole(userId, role)` |
| Data yang sering diakses (roles, academic year) | Cache di Redis dengan TTL |

### 2.3 Data Sync Strategy

```mermaid
sequenceDiagram
    participant Auth as auth-service
    participant NATS as NATS JetStream
    participant Master as master-data-service
    participant Cache as Redis

    Auth->>NATS: publish("user.created", {id, name, type})
    NATS->>Master: deliver to subscriber
    Master->>Cache: Cache user summary (TTL 5m)
    
    Note over Auth,Cache: Eventual consistency via events
```

---

## 3. Struktur Project

```
neocentral-go/
│
├── proto/                              # Shared protobuf definitions
│   ├── auth/
│   │   └── auth.proto
│   ├── master/
│   │   └── master.proto
│   ├── document/
│   │   └── document.proto
│   └── notification/
│       └── notification.proto
│
├── pkg/                                # Shared Go packages (go module)
│   ├── go.mod                          # module neocentral/pkg
│   ├── config/
│   │   └── config.go                   # Base config loader (viper)
│   ├── auth/
│   │   ├── jwt.go                      # JWT generate/verify
│   │   └── password.go                 # bcrypt hash/verify
│   ├── database/
│   │   └── mysql.go                    # GORM MySQL connection factory
│   ├── grpcclient/
│   │   └── dial.go                     # gRPC client connection helper
│   ├── messaging/
│   │   ├── nats.go                     # NATS connection + JetStream setup
│   │   ├── publisher.go                # Generic event publisher
│   │   └── subscriber.go              # Generic event subscriber
│   ├── middleware/
│   │   ├── cors.go
│   │   ├── logging.go                  # Structured request logging
│   │   ├── recovery.go
│   │   └── request_id.go
│   ├── response/
│   │   └── response.go                # Standardized JSON response
│   ├── errors/
│   │   └── errors.go                  # AppError type with HTTP status
│   ├── validator/
│   │   └── validator.go               # Struct validation wrapper
│   └── storage/
│       └── minio.go                   # MinIO/S3 client wrapper
│
├── auth-service/
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── student.go
│   │   │   └── lecturer.go
│   │   ├── repository/
│   │   │   ├── user_repo.go            # Interface
│   │   │   ├── user_repo_gorm.go       # GORM implementation
│   │   │   ├── role_repo.go
│   │   │   └── role_repo_gorm.go
│   │   ├── service/
│   │   │   ├── auth_service.go
│   │   │   ├── microsoft_auth_service.go
│   │   │   └── profile_service.go
│   │   ├── handler/                    # REST handlers (Chi)
│   │   │   ├── auth_handler.go
│   │   │   ├── profile_handler.go
│   │   │   └── router.go
│   │   ├── grpcserver/                 # gRPC server implementation
│   │   │   └── auth_grpc.go
│   │   ├── dto/
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   ├── event/                      # NATS event publishers
│   │   │   └── user_events.go
│   │   └── middleware/
│   │       └── auth_guard.go
│   ├── migrations/
│   │   ├── 000001_create_users.up.sql
│   │   ├── 000001_create_users.down.sql
│   │   ├── 000002_create_roles.up.sql
│   │   └── 000002_create_roles.down.sql
│   └── Dockerfile
│
├── master-data-service/
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── academic_year.go
│   │   │   ├── room.go
│   │   │   ├── science_group.go
│   │   │   ├── cpl.go
│   │   │   ├── cpmk.go
│   │   │   ├── assessment_criteria.go
│   │   │   └── lecturer_availability.go
│   │   ├── repository/
│   │   ├── service/
│   │   ├── handler/
│   │   ├── grpcserver/
│   │   ├── grpcclient/                 # Calls to auth-service
│   │   │   └── auth_client.go
│   │   ├── dto/
│   │   └── event/
│   ├── migrations/
│   └── Dockerfile
│
├── document-service/
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── document.go
│   │   │   ├── document_type.go
│   │   │   └── document_template.go
│   │   ├── repository/
│   │   ├── service/
│   │   │   ├── document_service.go
│   │   │   └── pdf_service.go          # Gotenberg integration
│   │   ├── handler/
│   │   ├── grpcserver/
│   │   ├── grpcclient/
│   │   │   └── auth_client.go
│   │   ├── storage/                    # MinIO upload/download
│   │   │   └── minio_storage.go
│   │   ├── dto/
│   │   └── event/
│   ├── migrations/
│   └── Dockerfile
│
├── notification-service/
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── notification.go
│   │   │   └── push_token.go
│   │   ├── repository/
│   │   ├── service/
│   │   │   └── notification_service.go
│   │   ├── handler/
│   │   ├── grpcserver/
│   │   ├── grpcclient/
│   │   │   └── auth_client.go
│   │   ├── provider/                   # Notification delivery providers
│   │   │   ├── fcm.go                  # Firebase Cloud Messaging
│   │   │   ├── smtp.go                 # Email via SMTP
│   │   │   └── websocket.go            # Real-time push
│   │   ├── subscriber/                 # NATS event subscribers
│   │   │   ├── user_subscriber.go
│   │   │   └── thesis_subscriber.go
│   │   ├── dto/
│   │   └── event/
│   ├── migrations/
│   └── Dockerfile
│
├── api-gateway/
│   ├── go.mod
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── middleware/
│   │   │   ├── jwt_auth.go             # Validate JWT, extract claims
│   │   │   ├── rate_limit.go
│   │   │   └── cors.go
│   │   ├── proxy/
│   │   │   └── reverse_proxy.go        # Route to downstream services
│   │   └── router/
│   │       └── routes.go               # Route definitions
│   └── Dockerfile
│
├── deployments/
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   ├── .env.example
│   └── k8s/                            # Future Kubernetes manifests
│
├── scripts/
│   ├── migrate.sh                      # Run migrations for all services
│   ├── proto-gen.sh                    # Generate Go code from .proto files
│   └── seed/
│       └── seed_master_data.go
│
├── Makefile
└── README.md
```

---

## 4. gRPC Contract Definitions

### 4.1 Auth Service Proto

```protobuf
// proto/auth/auth.proto
syntax = "proto3";

package auth;

option go_package = "neocentral/proto/auth";

service AuthService {
  // Validate JWT token and return user claims
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  
  // Get user by ID
  rpc GetUser(GetUserRequest) returns (UserResponse);
  
  // Get user contact info (for notifications)
  rpc GetUserContact(GetUserContactRequest) returns (UserContactResponse);
  
  // Check if user has specific role
  rpc HasRole(HasRoleRequest) returns (HasRoleResponse);
  
  // Get users by role (e.g., all lecturers, all admins)
  rpc GetUsersByRole(GetUsersByRoleRequest) returns (GetUsersByRoleResponse);
  
  // Batch get users
  rpc BatchGetUsers(BatchGetUsersRequest) returns (BatchGetUsersResponse);
}

message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string identity_number = 3;
  string full_name = 4;
  repeated string roles = 5;
}

message GetUserRequest {
  string user_id = 1;
}

message UserResponse {
  string id = 1;
  string full_name = 2;
  string identity_number = 3;
  string identity_type = 4;
  string email = 5;
  string phone_number = 6;
  string avatar_url = 7;
  bool is_verified = 8;
  repeated RoleInfo roles = 9;
  
  // Optional student/lecturer info
  StudentInfo student = 10;
  LecturerInfo lecturer = 11;
}

message RoleInfo {
  string id = 1;
  string name = 2;
  string status = 3;
}

message StudentInfo {
  int32 enrollment_year = 1;
  int32 sks_completed = 2;
  string status = 3;
  int32 current_semester = 4;
}

message LecturerInfo {
  string science_group_id = 1;
}

message GetUserContactRequest {
  string user_id = 1;
}

message UserContactResponse {
  string user_id = 1;
  string email = 2;
  string phone_number = 3;
  string full_name = 4;
}

message HasRoleRequest {
  string user_id = 1;
  string role_name = 2;
}

message HasRoleResponse {
  bool has_role = 1;
}

message GetUsersByRoleRequest {
  string role_name = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message GetUsersByRoleResponse {
  repeated UserResponse users = 1;
  int32 total = 2;
}

message BatchGetUsersRequest {
  repeated string user_ids = 1;
}

message BatchGetUsersResponse {
  repeated UserResponse users = 1;
}
```

### 4.2 Master Data Service Proto

```protobuf
// proto/master/master.proto
syntax = "proto3";

package master;

option go_package = "neocentral/proto/master";

service MasterDataService {
  rpc GetAcademicYear(GetAcademicYearRequest) returns (AcademicYearResponse);
  rpc GetActiveAcademicYear(EmptyRequest) returns (AcademicYearResponse);
  rpc ListAcademicYears(ListRequest) returns (AcademicYearListResponse);
  
  rpc GetRoom(GetRoomRequest) returns (RoomResponse);
  rpc ListRooms(ListRequest) returns (RoomListResponse);
  
  rpc GetCpl(GetCplRequest) returns (CplResponse);
  rpc ListCpls(ListRequest) returns (CplListResponse);
  
  rpc GetCpmk(GetCpmkRequest) returns (CpmkResponse);
  rpc ListCpmks(ListCpmkRequest) returns (CpmkListResponse);
  
  rpc ListScienceGroups(ListRequest) returns (ScienceGroupListResponse);
  
  rpc GetAssessmentCriteria(GetAssessmentCriteriaRequest) returns (AssessmentCriteriaResponse);
  rpc ListAssessmentCriteria(ListAssessmentCriteriaRequest) returns (AssessmentCriteriaListResponse);
}

message EmptyRequest {}

message ListRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message GetAcademicYearRequest {
  string id = 1;
}

message AcademicYearResponse {
  string id = 1;
  string semester = 2;
  string year = 3;
  string start_date = 4;
  string end_date = 5;
  bool is_active = 6;
}

message AcademicYearListResponse {
  repeated AcademicYearResponse academic_years = 1;
  int32 total = 2;
}

message GetRoomRequest {
  string id = 1;
}

message RoomResponse {
  string id = 1;
  string name = 2;
  string location = 3;
  int32 capacity = 4;
}

message RoomListResponse {
  repeated RoomResponse rooms = 1;
  int32 total = 2;
}

message GetCplRequest {
  string id = 1;
}

message CplResponse {
  string id = 1;
  string code = 2;
  string description = 3;
  int32 minimal_score = 4;
  bool is_active = 5;
}

message CplListResponse {
  repeated CplResponse cpls = 1;
  int32 total = 2;
}

message GetCpmkRequest {
  string id = 1;
}

message ListCpmkRequest {
  string academic_year_id = 1;
  string type = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message CpmkResponse {
  string id = 1;
  string academic_year_id = 2;
  string code = 3;
  string description = 4;
  string type = 5;
}

message CpmkListResponse {
  repeated CpmkResponse cpmks = 1;
  int32 total = 2;
}

message ScienceGroupListResponse {
  repeated ScienceGroupResponse science_groups = 1;
}

message ScienceGroupResponse {
  string id = 1;
  string name = 2;
}

message GetAssessmentCriteriaRequest {
  string id = 1;
}

message ListAssessmentCriteriaRequest {
  string cpmk_id = 1;
  string applies_to = 2;
}

message AssessmentCriteriaResponse {
  string id = 1;
  string cpmk_id = 2;
  string name = 3;
  string applies_to = 4;
  string role = 5;
  int32 max_score = 6;
  int32 display_order = 7;
  repeated AssessmentRubricResponse rubrics = 8;
}

message AssessmentCriteriaListResponse {
  repeated AssessmentCriteriaResponse criterias = 1;
}

message AssessmentRubricResponse {
  string id = 1;
  int32 min_score = 2;
  int32 max_score = 3;
  string description = 4;
  int32 display_order = 5;
}
```

### 4.3 Document Service Proto

```protobuf
// proto/document/document.proto
syntax = "proto3";

package document;

option go_package = "neocentral/proto/document";

service DocumentService {
  rpc UploadDocument(UploadDocumentRequest) returns (DocumentResponse);
  rpc GetDocument(GetDocumentRequest) returns (DocumentResponse);
  rpc GetDocumentURL(GetDocumentRequest) returns (DocumentURLResponse);
  rpc DeleteDocument(GetDocumentRequest) returns (DeleteResponse);
  rpc GeneratePDF(GeneratePDFRequest) returns (DocumentResponse);
}

message UploadDocumentRequest {
  string user_id = 1;
  string document_type_id = 2;
  string file_name = 3;
  bytes file_content = 4;
}

message GetDocumentRequest {
  string document_id = 1;
}

message DocumentResponse {
  string id = 1;
  string user_id = 2;
  string document_type_id = 3;
  string file_path = 4;
  string file_name = 5;
  string file_hash = 6;
  string created_at = 7;
}

message DocumentURLResponse {
  string url = 1;
  int64 expires_in = 2;
}

message DeleteResponse {
  bool success = 1;
}

message GeneratePDFRequest {
  string template_name = 1;
  string html_content = 2;
  string user_id = 3;
  string document_type_id = 4;
}
```

### 4.4 Notification Service Proto

```protobuf
// proto/notification/notification.proto
syntax = "proto3";

package notification;

option go_package = "neocentral/proto/notification";

service NotificationService {
  rpc SendNotification(SendNotificationRequest) returns (NotificationResponse);
  rpc SendBulkNotification(SendBulkNotificationRequest) returns (BulkNotificationResponse);
  rpc GetUserNotifications(GetUserNotificationsRequest) returns (NotificationListResponse);
  rpc MarkAsRead(MarkAsReadRequest) returns (NotificationResponse);
  rpc RegisterPushToken(RegisterPushTokenRequest) returns (RegisterPushTokenResponse);
}

message SendNotificationRequest {
  string user_id = 1;
  string title = 2;
  string message = 3;
  string channel = 4;  // "in_app", "push", "email", "all"
}

message SendBulkNotificationRequest {
  repeated string user_ids = 1;
  string title = 2;
  string message = 3;
  string channel = 4;
}

message NotificationResponse {
  string id = 1;
  string user_id = 2;
  string title = 3;
  string message = 4;
  bool is_read = 5;
  string created_at = 6;
}

message NotificationListResponse {
  repeated NotificationResponse notifications = 1;
  int32 total = 2;
  int32 unread_count = 3;
}

message GetUserNotificationsRequest {
  string user_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  bool unread_only = 4;
}

message MarkAsReadRequest {
  string notification_id = 1;
  string user_id = 2;
}

message BulkNotificationResponse {
  int32 sent_count = 1;
  int32 failed_count = 2;
}

message RegisterPushTokenRequest {
  string user_id = 1;
  string token = 2;
  string device_type = 3;  // "android", "ios", "web"
}

message RegisterPushTokenResponse {
  bool success = 1;
}
```

---

## 5. Contoh Implementasi Lengkap

### 5.1 Shared Package - Database (GORM)

```go
// pkg/database/mysql.go
package database

import (
    "fmt"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    MaxOpen  int
    MaxIdle  int
    MaxLife  time.Duration
}

func NewMySQLConnection(cfg Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxOpenConns(cfg.MaxOpen)
    sqlDB.SetMaxIdleConns(cfg.MaxIdle)
    sqlDB.SetConnMaxLifetime(cfg.MaxLife)

    return db, nil
}
```

### 5.2 Shared Package - MinIO Storage

```go
// pkg/storage/minio.go
package storage

import (
    "context"
    "fmt"
    "io"
    "time"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
    Endpoint  string
    AccessKey string
    SecretKey string
    Bucket    string
    UseSSL    bool
}

type MinIOStorage struct {
    client *minio.Client
    bucket string
}

func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
    client, err := minio.New(cfg.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
        Secure: cfg.UseSSL,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create MinIO client: %w", err)
    }

    // Ensure bucket exists
    ctx := context.Background()
    exists, err := client.BucketExists(ctx, cfg.Bucket)
    if err != nil {
        return nil, err
    }
    if !exists {
        if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
            return nil, fmt.Errorf("failed to create bucket: %w", err)
        }
    }

    return &MinIOStorage{client: client, bucket: cfg.Bucket}, nil
}

func (s *MinIOStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
    _, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
        ContentType: contentType,
    })
    return err
}

func (s *MinIOStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
    obj, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
    if err != nil {
        return nil, err
    }
    return obj, nil
}

func (s *MinIOStorage) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
    url, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
    if err != nil {
        return "", err
    }
    return url.String(), nil
}

func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
    return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
```

### 5.3 Shared Package - NATS JetStream

```go
// pkg/messaging/nats.go
package messaging

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/nats-io/nats.go"
)

type NATSClient struct {
    conn *nats.Conn
    js   nats.JetStreamContext
}

func NewNATSClient(url string) (*NATSClient, error) {
    nc, err := nats.Connect(url,
        nats.RetryOnFailedConnect(true),
        nats.MaxReconnects(-1),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to NATS: %w", err)
    }

    js, err := nc.JetStream()
    if err != nil {
        return nil, fmt.Errorf("failed to create JetStream context: %w", err)
    }

    return &NATSClient{conn: nc, js: js}, nil
}

// EnsureStream creates a stream if it doesn't exist
func (c *NATSClient) EnsureStream(name string, subjects []string) error {
    _, err := c.js.StreamInfo(name)
    if err != nil {
        _, err = c.js.AddStream(&nats.StreamConfig{
            Name:     name,
            Subjects: subjects,
            Storage:  nats.FileStorage,
            MaxAge:   0, // unlimited
        })
        if err != nil {
            return fmt.Errorf("failed to create stream %s: %w", name, err)
        }
        log.Printf("📡 Created NATS stream: %s", name)
    }
    return nil
}

// Publish sends a message to a NATS subject
func (c *NATSClient) Publish(subject string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    _, err = c.js.Publish(subject, data)
    return err
}

// Subscribe creates a durable consumer for a subject
func (c *NATSClient) Subscribe(subject, durableName string, handler func(msg *nats.Msg)) error {
    _, err := c.js.Subscribe(subject, handler,
        nats.Durable(durableName),
        nats.DeliverAll(),
        nats.AckExplicit(),
    )
    return err
}

func (c *NATSClient) Close() {
    c.conn.Close()
}
```

### 5.4 Auth Service - GORM Domain Models

```go
// auth-service/internal/domain/user.go
package domain

import (
    "time"

    "gorm.io/gorm"
)

type IdentityType string

const (
    IdentityNIM   IdentityType = "NIM"
    IdentityNIP   IdentityType = "NIP"
    IdentityOTHER IdentityType = "OTHER"
)

type User struct {
    ID                string       `gorm:"primaryKey;type:varchar(255)" json:"id"`
    FullName          string       `gorm:"column:full_name;not null" json:"fullName"`
    IdentityNumber    string       `gorm:"column:identity_number;uniqueIndex;not null" json:"identityNumber"`
    IdentityType      IdentityType `gorm:"column:identity_type;not null" json:"identityType"`
    Email             *string      `gorm:"uniqueIndex" json:"email"`
    Password          *string      `json:"-"`
    PhoneNumber       *string      `gorm:"column:phone_number" json:"phoneNumber"`
    IsVerified        bool         `gorm:"default:false" json:"isVerified"`
    Token             *string      `gorm:"type:text" json:"-"`
    RefreshToken      *string      `gorm:"column:refresh_token;type:text" json:"-"`
    OAuthProvider     *string      `gorm:"column:oauth_provider" json:"-"`
    OAuthID           *string      `gorm:"column:oauth_id" json:"-"`
    OAuthRefreshToken *string      `gorm:"column:oauth_refresh_token;type:text" json:"-"`
    AvatarURL         *string      `gorm:"column:avatar_url" json:"avatarUrl"`
    CreatedAt         time.Time    `json:"createdAt"`
    UpdatedAt         time.Time    `json:"updatedAt"`

    // Relations
    UserHasRoles []UserHasRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
    Student      *Student      `gorm:"foreignKey:ID" json:"student,omitempty"`
    Lecturer     *Lecturer     `gorm:"foreignKey:ID" json:"lecturer,omitempty"`
}

func (User) TableName() string { return "users" }

type UserRole struct {
    ID   string `gorm:"primaryKey;type:varchar(255)" json:"id"`
    Name string `gorm:"not null" json:"name"`
}

func (UserRole) TableName() string { return "user_roles" }

type RoleStatus string

const (
    RoleStatusActive    RoleStatus = "active"
    RoleStatusNonActive RoleStatus = "nonActive"
)

type UserHasRole struct {
    UserID string     `gorm:"column:user_id;primaryKey" json:"userId"`
    RoleID string     `gorm:"column:role_id;primaryKey" json:"roleId"`
    Status RoleStatus `gorm:"not null" json:"status"`

    Role UserRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
    User User     `gorm:"foreignKey:UserID" json:"-"`
}

func (UserHasRole) TableName() string { return "user_has_roles" }

type Student struct {
    ID                        string    `gorm:"column:user_id;primaryKey" json:"id"`
    Status                    string    `gorm:"column:student_status;default:active" json:"status"`
    EnrollmentYear            *int      `gorm:"column:enrollment_year" json:"enrollmentYear"`
    SKSCompleted              int       `gorm:"column:skscompleted" json:"sksCompleted"`
    MandatoryCoursesCompleted bool      `gorm:"column:mandatory_courses_completed;default:false" json:"mandatoryCoursesCompleted"`
    MKWUCompleted             bool      `gorm:"column:mkwu_completed;default:false" json:"mkwuCompleted"`
    InternshipCompleted       bool      `gorm:"column:internship_completed;default:false" json:"internshipCompleted"`
    KKNCompleted              bool      `gorm:"column:kkn_completed;default:false" json:"kknCompleted"`
    ResearchMethodCompleted   bool      `gorm:"column:research_method_completed;default:false" json:"researchMethodCompleted"`
    CurrentSemester           *int      `gorm:"column:current_semester" json:"currentSemester"`
    CreatedAt                 time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt                 time.Time `gorm:"column:updated_at" json:"updatedAt"`

    User User `gorm:"foreignKey:ID;references:ID" json:"-"`
}

func (Student) TableName() string { return "students" }

type Lecturer struct {
    ID             string    `gorm:"column:user_id;primaryKey" json:"id"`
    ScienceGroupID *string   `gorm:"column:science_group_id" json:"scienceGroupId"`
    Data           *string   `gorm:"type:json" json:"data"`
    CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`

    User User `gorm:"foreignKey:ID;references:ID" json:"-"`
}

func (Lecturer) TableName() string { return "lecturers" }
```

### 5.5 Auth Service - GORM Repository

```go
// auth-service/internal/repository/user_repo.go
package repository

import (
    "context"
    "auth-service/internal/domain"
)

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindByIdentityNumber(ctx context.Context, identityNumber string) (*domain.User, error)
    Create(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
    UpdateTokens(ctx context.Context, id, token, refreshToken string) error
    GetUserRoles(ctx context.Context, userID string) ([]domain.UserHasRole, error)
    HasRole(ctx context.Context, userID, roleName string) (bool, error)
    GetUsersByRole(ctx context.Context, roleName string, page, pageSize int) ([]domain.User, int64, error)
    BatchGetByIDs(ctx context.Context, ids []string) ([]domain.User, error)
}
```

```go
// auth-service/internal/repository/user_repo_gorm.go
package repository

import (
    "context"
    "errors"

    "auth-service/internal/domain"
    "gorm.io/gorm"
)

type gormUserRepo struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
    return &gormUserRepo{db: db}
}

func (r *gormUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).
        Preload("UserHasRoles.Role").
        Preload("Student").
        Preload("Lecturer").
        First(&user, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *gormUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).
        Preload("UserHasRoles.Role").
        First(&user, "email = ?", email).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *gormUserRepo) FindByIdentityNumber(ctx context.Context, identityNumber string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).
        Preload("UserHasRoles.Role").
        Preload("Student").
        Preload("Lecturer").
        First(&user, "identity_number = ?", identityNumber).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *gormUserRepo) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepo) Update(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormUserRepo) UpdateTokens(ctx context.Context, id, token, refreshToken string) error {
    return r.db.WithContext(ctx).
        Model(&domain.User{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "token":         token,
            "refresh_token": refreshToken,
        }).Error
}

func (r *gormUserRepo) GetUserRoles(ctx context.Context, userID string) ([]domain.UserHasRole, error) {
    var roles []domain.UserHasRole
    err := r.db.WithContext(ctx).
        Preload("Role").
        Where("user_id = ? AND status = ?", userID, domain.RoleStatusActive).
        Find(&roles).Error
    return roles, err
}

func (r *gormUserRepo) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.UserHasRole{}).
        Joins("JOIN user_roles ON user_roles.id = user_has_roles.role_id").
        Where("user_has_roles.user_id = ? AND LOWER(user_roles.name) = LOWER(?) AND user_has_roles.status = ?",
            userID, roleName, domain.RoleStatusActive).
        Count(&count).Error
    return count > 0, err
}

func (r *gormUserRepo) GetUsersByRole(ctx context.Context, roleName string, page, pageSize int) ([]domain.User, int64, error) {
    var users []domain.User
    var total int64

    subquery := r.db.WithContext(ctx).
        Model(&domain.UserHasRole{}).
        Select("user_id").
        Joins("JOIN user_roles ON user_roles.id = user_has_roles.role_id").
        Where("LOWER(user_roles.name) = LOWER(?) AND user_has_roles.status = ?",
            roleName, domain.RoleStatusActive)

    query := r.db.WithContext(ctx).
        Where("id IN (?)", subquery).
        Preload("UserHasRoles.Role")

    query.Model(&domain.User{}).Count(&total)

    offset := (page - 1) * pageSize
    err := query.Offset(offset).Limit(pageSize).Find(&users).Error

    return users, total, err
}

func (r *gormUserRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
    var users []domain.User
    err := r.db.WithContext(ctx).
        Preload("UserHasRoles.Role").
        Preload("Student").
        Preload("Lecturer").
        Where("id IN ?", ids).
        Find(&users).Error
    return users, err
}
```

### 5.6 Auth Service - gRPC Server Implementation

```go
// auth-service/internal/grpcserver/auth_grpc.go
package grpcserver

import (
    "context"

    "auth-service/internal/domain"
    "auth-service/internal/repository"
    "auth-service/internal/service"
    pb "neocentral/proto/auth"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type AuthGRPCServer struct {
    pb.UnimplementedAuthServiceServer
    authService *service.AuthService
    userRepo    repository.UserRepository
}

func NewAuthGRPCServer(authSvc *service.AuthService, userRepo repository.UserRepository) *AuthGRPCServer {
    return &AuthGRPCServer{
        authService: authSvc,
        userRepo:    userRepo,
    }
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    claims, err := s.authService.VerifyAccessToken(req.Token)
    if err != nil {
        return &pb.ValidateTokenResponse{Valid: false}, nil
    }

    return &pb.ValidateTokenResponse{
        Valid:          true,
        UserId:         claims.UserID,
        IdentityNumber: claims.IdentityNumber,
        FullName:       claims.FullName,
        Roles:          claims.Roles,
    }, nil
}

func (s *AuthGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, req.UserId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
    }
    if user == nil {
        return nil, status.Errorf(codes.NotFound, "user not found: %s", req.UserId)
    }

    return toUserProto(user), nil
}

func (s *AuthGRPCServer) GetUserContact(ctx context.Context, req *pb.GetUserContactRequest) (*pb.UserContactResponse, error) {
    user, err := s.userRepo.FindByID(ctx, req.UserId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
    }
    if user == nil {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }

    resp := &pb.UserContactResponse{
        UserId:   user.ID,
        FullName: user.FullName,
    }
    if user.Email != nil {
        resp.Email = *user.Email
    }
    if user.PhoneNumber != nil {
        resp.PhoneNumber = *user.PhoneNumber
    }
    return resp, nil
}

func (s *AuthGRPCServer) HasRole(ctx context.Context, req *pb.HasRoleRequest) (*pb.HasRoleResponse, error) {
    hasRole, err := s.userRepo.HasRole(ctx, req.UserId, req.RoleName)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to check role: %v", err)
    }
    return &pb.HasRoleResponse{HasRole: hasRole}, nil
}

func (s *AuthGRPCServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
    users, err := s.userRepo.BatchGetByIDs(ctx, req.UserIds)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to batch get users: %v", err)
    }

    var protoUsers []*pb.UserResponse
    for i := range users {
        protoUsers = append(protoUsers, toUserProto(&users[i]))
    }
    return &pb.BatchGetUsersResponse{Users: protoUsers}, nil
}

// Helper: domain.User → proto UserResponse
func toUserProto(u *domain.User) *pb.UserResponse {
    resp := &pb.UserResponse{
        Id:             u.ID,
        FullName:       u.FullName,
        IdentityNumber: u.IdentityNumber,
        IdentityType:   string(u.IdentityType),
        IsVerified:     u.IsVerified,
    }
    if u.Email != nil {
        resp.Email = *u.Email
    }
    if u.PhoneNumber != nil {
        resp.PhoneNumber = *u.PhoneNumber
    }
    if u.AvatarURL != nil {
        resp.AvatarUrl = *u.AvatarURL
    }

    for _, uhr := range u.UserHasRoles {
        resp.Roles = append(resp.Roles, &pb.RoleInfo{
            Id:     uhr.RoleID,
            Name:   uhr.Role.Name,
            Status: string(uhr.Status),
        })
    }

    if u.Student != nil {
        resp.Student = &pb.StudentInfo{
            SksCompleted:    int32(u.Student.SKSCompleted),
            Status:          u.Student.Status,
        }
        if u.Student.EnrollmentYear != nil {
            resp.Student.EnrollmentYear = int32(*u.Student.EnrollmentYear)
        }
        if u.Student.CurrentSemester != nil {
            resp.Student.CurrentSemester = int32(*u.Student.CurrentSemester)
        }
    }

    if u.Lecturer != nil {
        resp.Lecturer = &pb.LecturerInfo{}
        if u.Lecturer.ScienceGroupID != nil {
            resp.Lecturer.ScienceGroupId = *u.Lecturer.ScienceGroupID
        }
    }

    return resp
}
```

### 5.7 Auth Service - main.go (Dual REST + gRPC)

```go
// auth-service/cmd/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    chimw "github.com/go-chi/chi/v5/middleware"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "auth-service/internal/config"
    "auth-service/internal/grpcserver"
    "auth-service/internal/handler"
    "auth-service/internal/repository"
    "auth-service/internal/service"
    "auth-service/internal/event"
    "neocentral/pkg/database"
    "neocentral/pkg/messaging"
    "neocentral/pkg/middleware"
    pb "neocentral/proto/auth"
)

func main() {
    cfg := config.Load()

    // ── Database (GORM) ──────────────────────────
    db, err := database.NewMySQLConnection(database.Config{
        Host:     cfg.DB.Host,
        Port:     cfg.DB.Port,
        User:     cfg.DB.User,
        Password: cfg.DB.Password,
        DBName:   cfg.DB.Name, // "neocentral_auth"
        MaxOpen:  25,
        MaxIdle:  10,
        MaxLife:  5 * time.Minute,
    })
    if err != nil {
        log.Fatalf("❌ Database connection failed: %v", err)
    }
    log.Println("✅ Connected to MySQL (neocentral_auth)")

    // ── NATS ─────────────────────────────────────
    natsClient, err := messaging.NewNATSClient(cfg.NATSURL)
    if err != nil {
        log.Fatalf("❌ NATS connection failed: %v", err)
    }
    defer natsClient.Close()

    _ = natsClient.EnsureStream("AUTH", []string{"auth.>"})
    log.Println("✅ Connected to NATS JetStream")

    // ── Repositories ─────────────────────────────
    userRepo := repository.NewGormUserRepository(db)

    // ── Event Publisher ──────────────────────────
    userEvents := event.NewUserEventPublisher(natsClient)

    // ── Services ─────────────────────────────────
    authSvc := service.NewAuthService(userRepo, cfg.JWT, userEvents)

    // ── REST Server (Chi) ────────────────────────
    r := chi.NewRouter()
    r.Use(chimw.RequestID)
    r.Use(chimw.RealIP)
    r.Use(chimw.Logger)
    r.Use(chimw.Recoverer)
    r.Use(middleware.CORS())

    authHandler := handler.NewAuthHandler(authSvc)
    profileHandler := handler.NewProfileHandler(authSvc, userRepo)

    r.Route("/auth", func(r chi.Router) {
        r.Post("/login", authHandler.Login)
        r.Post("/register", authHandler.Register)
        r.Post("/refresh", authHandler.RefreshToken)
        r.Post("/microsoft/callback", authHandler.MicrosoftCallback)
        r.Post("/logout", authHandler.Logout)
    })

    r.Route("/profile", func(r chi.Router) {
        r.Use(handler.AuthGuardMiddleware(cfg.JWT))
        r.Get("/me", profileHandler.GetMe)
        r.Put("/me", profileHandler.UpdateProfile)
    })

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"status":"ok","service":"auth-service","uptime":"%s"}`,
            time.Since(time.Now()).String())
    })

    httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
    httpServer := &http.Server{Addr: httpAddr, Handler: r}

    // ── gRPC Server ──────────────────────────────
    grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
    lis, err := net.Listen("tcp", grpcAddr)
    if err != nil {
        log.Fatalf("❌ Failed to listen gRPC on %s: %v", grpcAddr, err)
    }

    grpcSrv := grpc.NewServer()
    authGRPC := grpcserver.NewAuthGRPCServer(authSvc, userRepo)
    pb.RegisterAuthServiceServer(grpcSrv, authGRPC)
    reflection.Register(grpcSrv) // for grpcurl/grpcui debugging

    // ── Start servers ────────────────────────────
    go func() {
        log.Printf("🌐 REST server listening on %s", httpAddr)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()

    go func() {
        log.Printf("⚡ gRPC server listening on %s", grpcAddr)
        if err := grpcSrv.Serve(lis); err != nil {
            log.Fatalf("gRPC server error: %v", err)
        }
    }()

    // ── Graceful shutdown ────────────────────────
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("🛑 Shutting down auth-service...")

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    grpcSrv.GracefulStop()
    httpServer.Shutdown(ctx)

    log.Println("✅ auth-service stopped gracefully")
}
```

---

## 6. Event Catalog (NATS JetStream)

### 6.1 Stream & Subject Design

| Stream | Subjects | Publisher | Subscribers |
|---|---|---|---|
| `AUTH` | `auth.user.created`, `auth.user.updated`, `auth.user.login` | auth-service | notification-service, master-data-service |
| `MASTER` | `master.academic_year.activated`, `master.room.created` | master-data-service | thesis-service*, guidance-service* |
| `DOCUMENT` | `document.uploaded`, `document.deleted` | document-service | thesis-service*, internship-service* |
| `NOTIFICATION` | `notification.send`, `notification.bulk_send` | any service | notification-service |

> \* Services bertanda asterisk adalah subscriber yang diimplementasi di Fase 2.

### 6.2 Contoh Event Publisher

```go
// auth-service/internal/event/user_events.go
package event

import (
    "neocentral/pkg/messaging"
    "time"
)

type UserEventPublisher struct {
    nats *messaging.NATSClient
}

func NewUserEventPublisher(nats *messaging.NATSClient) *UserEventPublisher {
    return &UserEventPublisher{nats: nats}
}

type UserCreatedEvent struct {
    UserID         string `json:"userId"`
    FullName       string `json:"fullName"`
    Email          string `json:"email"`
    IdentityType   string `json:"identityType"`
    IdentityNumber string `json:"identityNumber"`
    Timestamp      string `json:"timestamp"`
}

type UserLoginEvent struct {
    UserID    string `json:"userId"`
    IP        string `json:"ip"`
    Timestamp string `json:"timestamp"`
}

func (p *UserEventPublisher) PublishUserCreated(userID, fullName, email, identityType, identityNumber string) error {
    return p.nats.Publish("auth.user.created", UserCreatedEvent{
        UserID:         userID,
        FullName:       fullName,
        Email:          email,
        IdentityType:   identityType,
        IdentityNumber: identityNumber,
        Timestamp:      time.Now().UTC().Format(time.RFC3339),
    })
}

func (p *UserEventPublisher) PublishUserLogin(userID, ip string) error {
    return p.nats.Publish("auth.user.login", UserLoginEvent{
        UserID:    userID,
        IP:        ip,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    })
}
```

### 6.3 Contoh Event Subscriber (notification-service)

```go
// notification-service/internal/subscriber/user_subscriber.go
package subscriber

import (
    "encoding/json"
    "log"

    "github.com/nats-io/nats.go"
    "notification-service/internal/service"
)

type UserSubscriber struct {
    notifService *service.NotificationService
}

func NewUserSubscriber(svc *service.NotificationService) *UserSubscriber {
    return &UserSubscriber{notifService: svc}
}

func (s *UserSubscriber) HandleUserCreated(msg *nats.Msg) {
    var event struct {
        UserID   string `json:"userId"`
        FullName string `json:"fullName"`
        Email    string `json:"email"`
    }

    if err := json.Unmarshal(msg.Data, &event); err != nil {
        log.Printf("❌ Failed to unmarshal user.created event: %v", err)
        msg.Nak()
        return
    }

    log.Printf("📩 Received user.created: %s (%s)", event.FullName, event.Email)

    // Send welcome email
    if event.Email != "" {
        err := s.notifService.SendWelcomeEmail(event.Email, event.FullName)
        if err != nil {
            log.Printf("❌ Failed to send welcome email: %v", err)
            msg.Nak()
            return
        }
    }

    msg.Ack()
}
```

---

## 7. API Gateway

### 7.1 Routing Configuration

```go
// api-gateway/internal/router/routes.go
package router

import (
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/go-chi/chi/v5"
    "api-gateway/internal/middleware"
)

type ServiceRoute struct {
    Prefix    string
    TargetURL string
    AuthRequired bool
}

func SetupRoutes(r chi.Router, jwtSecret string) {
    services := []ServiceRoute{
        // Public routes (no auth needed)
        {Prefix: "/auth", TargetURL: "http://auth-service:8001", AuthRequired: false},

        // Protected routes
        {Prefix: "/profile", TargetURL: "http://auth-service:8001", AuthRequired: true},
        {Prefix: "/master-data", TargetURL: "http://master-data-service:8002", AuthRequired: true},
        {Prefix: "/documents", TargetURL: "http://document-service:8009", AuthRequired: true},
        {Prefix: "/notifications", TargetURL: "http://notification-service:8008", AuthRequired: true},

        // Fase 2 services (placeholder)
        // {Prefix: "/thesis", TargetURL: "http://thesis-service:8003", AuthRequired: true},
        // {Prefix: "/guidance", TargetURL: "http://guidance-service:8004", AuthRequired: true},
        // {Prefix: "/seminars", TargetURL: "http://examination-service:8005", AuthRequired: true},
        // {Prefix: "/defences", TargetURL: "http://examination-service:8005", AuthRequired: true},
        // {Prefix: "/internship", TargetURL: "http://internship-service:8006", AuthRequired: true},
        // {Prefix: "/yudisium", TargetURL: "http://yudisium-service:8007", AuthRequired: true},
    }

    authMiddleware := middleware.JWTAuth(jwtSecret)

    for _, svc := range services {
        target, _ := url.Parse(svc.TargetURL)
        proxy := httputil.NewSingleHostReverseProxy(target)

        if svc.AuthRequired {
            r.Route(svc.Prefix, func(r chi.Router) {
                r.Use(authMiddleware)
                r.Handle("/*", http.StripPrefix("", proxy))
            })
        } else {
            r.Route(svc.Prefix, func(r chi.Router) {
                r.Handle("/*", http.StripPrefix("", proxy))
            })
        }
    }
}
```

### 7.2 JWT Middleware (Gateway)

```go
// api-gateway/internal/middleware/jwt_auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "neocentral/pkg/auth"
    "neocentral/pkg/response"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                response.Error(w, http.StatusUnauthorized, "missing authorization header")
                return
            }

            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) != 2 || parts[0] != "Bearer" {
                response.Error(w, http.StatusUnauthorized, "invalid authorization format")
                return
            }

            claims, err := auth.VerifyAccessToken(parts[1], jwtSecret)
            if err != nil {
                response.Error(w, http.StatusUnauthorized, "invalid or expired token")
                return
            }

            // Forward user info as headers to downstream services
            r.Header.Set("X-User-ID", claims.UserID)
            r.Header.Set("X-User-Identity", claims.IdentityNumber)
            r.Header.Set("X-User-Name", claims.FullName)
            r.Header.Set("X-User-Roles", strings.Join(claims.Roles, ","))

            ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

## 8. Docker Compose (Fase 1 - Complete)

```yaml
# deployments/docker-compose.yml
services:

  # ── API Gateway ──────────────────────────────────
  api-gateway:
    build:
      context: ..
      dockerfile: api-gateway/Dockerfile
    container_name: neo-gateway
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      HTTP_PORT: 8080
      JWT_SECRET: ${JWT_SECRET}
      AUTH_SERVICE_URL: http://auth-service:8001
      MASTER_DATA_SERVICE_URL: http://master-data-service:8002
      DOCUMENT_SERVICE_URL: http://document-service:8009
      NOTIFICATION_SERVICE_URL: http://notification-service:8008
    depends_on:
      auth-service:
        condition: service_healthy
      master-data-service:
        condition: service_healthy
      document-service:
        condition: service_healthy
      notification-service:
        condition: service_healthy
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/health"]
      interval: 15s
      timeout: 5s
      retries: 3

  # ── Auth Service ─────────────────────────────────
  auth-service:
    build:
      context: ..
      dockerfile: auth-service/Dockerfile
    container_name: neo-auth
    restart: unless-stopped
    environment:
      HTTP_PORT: 8001
      GRPC_PORT: 9001
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: neocentral_auth
      REDIS_URL: redis://redis:6379
      NATS_URL: nats://nats:4222
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRES_IN: ${JWT_EXPIRES_IN}
      REFRESH_TOKEN_SECRET: ${REFRESH_TOKEN_SECRET}
      CLIENT_ID: ${CLIENT_ID}
      CLIENT_SECRET: ${CLIENT_SECRET}
      TENANT_ID: ${TENANT_ID}
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_started
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8001/health"]
      interval: 15s
      timeout: 5s
      retries: 5
      start_period: 30s

  # ── Master Data Service ──────────────────────────
  master-data-service:
    build:
      context: ..
      dockerfile: master-data-service/Dockerfile
    container_name: neo-master
    restart: unless-stopped
    environment:
      HTTP_PORT: 8002
      GRPC_PORT: 9002
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: neocentral_master
      REDIS_URL: redis://redis:6379
      NATS_URL: nats://nats:4222
      AUTH_GRPC_ADDR: auth-service:9001
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8002/health"]
      interval: 15s
      timeout: 5s
      retries: 5

  # ── Document Service ─────────────────────────────
  document-service:
    build:
      context: ..
      dockerfile: document-service/Dockerfile
    container_name: neo-document
    restart: unless-stopped
    environment:
      HTTP_PORT: 8009
      GRPC_PORT: 9009
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: neocentral_doc
      NATS_URL: nats://nats:4222
      AUTH_GRPC_ADDR: auth-service:9001
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: ${MINIO_ACCESS_KEY}
      MINIO_SECRET_KEY: ${MINIO_SECRET_KEY}
      MINIO_BUCKET: neocentral-documents
      MINIO_USE_SSL: "false"
      GOTENBERG_URL: http://gotenberg:3000
    depends_on:
      mysql:
        condition: service_healthy
      minio:
        condition: service_healthy
      gotenberg:
        condition: service_started
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8009/health"]
      interval: 15s
      timeout: 5s
      retries: 5

  # ── Notification Service ─────────────────────────
  notification-service:
    build:
      context: ..
      dockerfile: notification-service/Dockerfile
    container_name: neo-notification
    restart: unless-stopped
    environment:
      HTTP_PORT: 8008
      GRPC_PORT: 9008
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: neocentral_notif
      REDIS_URL: redis://redis:6379
      NATS_URL: nats://nats:4222
      AUTH_GRPC_ADDR: auth-service:9001
      SMTP_HOST: ${SMTP_HOST}
      SMTP_PORT: ${SMTP_PORT}
      SMTP_USER: ${SMTP_USER}
      SMTP_PASS: ${SMTP_PASS}
      SMTP_FROM: ${SMTP_FROM}
      FCM_PROJECT_ID: ${FCM_PROJECT_ID}
      FCM_CLIENT_EMAIL: ${FCM_CLIENT_EMAIL}
      FCM_PRIVATE_KEY: ${FCM_PRIVATE_KEY}
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_started
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8008/health"]
      interval: 15s
      timeout: 5s
      retries: 5

  # ═══════════════ INFRASTRUCTURE ═══════════════════

  # ── MySQL ────────────────────────────────────────
  mysql:
    image: mysql:8.0
    container_name: neo-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
    volumes:
      - mysql-data:/var/lib/mysql
      - ./init-databases.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "3306:3306"
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s

  # ── Redis ────────────────────────────────────────
  redis:
    image: redis:7-alpine
    container_name: neo-redis
    restart: unless-stopped
    command: redis-server --save 60 1 --loglevel warning
    volumes:
      - redis-data:/data
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ── NATS (JetStream) ────────────────────────────
  nats:
    image: nats:2-alpine
    container_name: neo-nats
    restart: unless-stopped
    command: ["--jetstream", "--store_dir=/data"]
    volumes:
      - nats-data:/data
    ports:
      - "4222:4222"   # Client
      - "8222:8222"   # Monitoring
    networks:
      - neocentral-net

  # ── MinIO (S3-compatible Object Storage) ────────
  minio:
    image: minio/minio
    container_name: neo-minio
    restart: unless-stopped
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    volumes:
      - minio-data:/data
    ports:
      - "9000:9000"   # API
      - "9001:9001"   # Console
    networks:
      - neocentral-net
    healthcheck:
      test: ["CMD", "mc", "ready", "local"]
      interval: 15s
      timeout: 5s
      retries: 5

  # ── Gotenberg (PDF Generation) ──────────────────
  gotenberg:
    image: gotenberg/gotenberg:8
    container_name: neo-gotenberg
    restart: unless-stopped
    command:
      - gotenberg
      - "--chromium-disable-routes=false"
      - "--api-timeout=60s"
    networks:
      - neocentral-net

volumes:
  mysql-data:
  redis-data:
  nats-data:
  minio-data:

networks:
  neocentral-net:
    driver: bridge
```

### 8.1 Database Initialization Script

```sql
-- deployments/init-databases.sql
-- Creates separate databases for each microservice

CREATE DATABASE IF NOT EXISTS neocentral_auth
    CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS neocentral_master
    CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS neocentral_doc
    CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS neocentral_notif
    CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant permissions to service user
CREATE USER IF NOT EXISTS 'neocentral'@'%' IDENTIFIED BY 'change_me_in_production';
GRANT ALL PRIVILEGES ON neocentral_auth.* TO 'neocentral'@'%';
GRANT ALL PRIVILEGES ON neocentral_master.* TO 'neocentral'@'%';
GRANT ALL PRIVILEGES ON neocentral_doc.* TO 'neocentral'@'%';
GRANT ALL PRIVILEGES ON neocentral_notif.* TO 'neocentral'@'%';
FLUSH PRIVILEGES;
```

---

## 9. Migration Files (Contoh)

### 9.1 Auth Service Migrations

```sql
-- auth-service/migrations/000001_create_users.up.sql
CREATE TABLE IF NOT EXISTS users (
    id              VARCHAR(255) PRIMARY KEY,
    full_name       VARCHAR(255) NOT NULL,
    identity_number VARCHAR(255) NOT NULL UNIQUE,
    identity_type   ENUM('NIM', 'NIP', 'OTHER') NOT NULL,
    email           VARCHAR(255) UNIQUE,
    password        VARCHAR(255),
    phone_number    VARCHAR(50),
    is_verified     BOOLEAN DEFAULT FALSE,
    token           TEXT,
    refresh_token   TEXT,
    oauth_provider  VARCHAR(100),
    oauth_id        VARCHAR(255),
    oauth_refresh_token TEXT,
    avatar_url      VARCHAR(500),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_roles (
    id   VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_has_roles (
    user_id VARCHAR(255) NOT NULL,
    role_id VARCHAR(255) NOT NULL,
    status  ENUM('active', 'nonActive') NOT NULL DEFAULT 'active',
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES user_roles(id) ON DELETE CASCADE,
    INDEX idx_user_has_roles_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS students (
    user_id                    VARCHAR(255) PRIMARY KEY,
    student_status             ENUM('dropout','bss','lulus','mengundurkan_diri','active') DEFAULT 'active',
    enrollment_year            INT,
    skscompleted               INT NOT NULL DEFAULT 0,
    mandatory_courses_completed BOOLEAN DEFAULT FALSE,
    mkwu_completed             BOOLEAN DEFAULT FALSE,
    internship_completed       BOOLEAN DEFAULT FALSE,
    kkn_completed              BOOLEAN DEFAULT FALSE,
    research_method_completed  BOOLEAN DEFAULT FALSE,
    current_semester           INT,
    created_at                 TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_students_skscompleted (skscompleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lecturers (
    user_id          VARCHAR(255) PRIMARY KEY,
    science_group_id VARCHAR(255),
    data             JSON,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

```sql
-- auth-service/migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS lecturers;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS user_has_roles;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
```

---

## 10. Makefile

```makefile
# Makefile
.PHONY: proto build test run-dev docker-up docker-down migrate

# ── Proto generation ──────────────────────────────
proto:
	@echo "🔧 Generating protobuf Go code..."
	protoc --go_out=. --go-grpc_out=. proto/auth/auth.proto
	protoc --go_out=. --go-grpc_out=. proto/master/master.proto
	protoc --go_out=. --go-grpc_out=. proto/document/document.proto
	protoc --go_out=. --go-grpc_out=. proto/notification/notification.proto
	@echo "✅ Proto generation complete"

# ── Build all services ────────────────────────────
build:
	@echo "🏗️  Building all services..."
	cd auth-service && go build -o ../bin/auth-service ./cmd/main.go
	cd master-data-service && go build -o ../bin/master-data-service ./cmd/main.go
	cd document-service && go build -o ../bin/document-service ./cmd/main.go
	cd notification-service && go build -o ../bin/notification-service ./cmd/main.go
	cd api-gateway && go build -o ../bin/api-gateway ./cmd/main.go
	@echo "✅ Build complete"

# ── Build individual service ──────────────────────
build-%:
	cd $*-service && go build -o ../bin/$*-service ./cmd/main.go

# ── Test all services ─────────────────────────────
test:
	cd auth-service && go test ./... -v -count=1
	cd master-data-service && go test ./... -v -count=1
	cd document-service && go test ./... -v -count=1
	cd notification-service && go test ./... -v -count=1

test-%:
	cd $*-service && go test ./... -v -count=1

# ── Run locally ───────────────────────────────────
run-%:
	cd $*-service && go run ./cmd/main.go

# ── Docker ────────────────────────────────────────
docker-up:
	cd deployments && docker compose up -d --build

docker-down:
	cd deployments && docker compose down

docker-logs:
	cd deployments && docker compose logs -f

# ── Database migrations ──────────────────────────
migrate-up-%:
	migrate -path ./$*-service/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(localhost:3306)/neocentral_$*" up

migrate-down-%:
	migrate -path ./$*-service/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(localhost:3306)/neocentral_$*" down 1

migrate-up-all:
	$(MAKE) migrate-up-auth
	$(MAKE) migrate-up-master
	$(MAKE) migrate-up-doc
	$(MAKE) migrate-up-notif

# ── Lint ──────────────────────────────────────────
lint:
	golangci-lint run ./...

# ── Tidy all modules ─────────────────────────────
tidy:
	cd pkg && go mod tidy
	cd auth-service && go mod tidy
	cd master-data-service && go mod tidy
	cd document-service && go mod tidy
	cd notification-service && go mod tidy
	cd api-gateway && go mod tidy
```

---

## 11. Contoh Dockerfile (auth-service)

```dockerfile
# auth-service/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared packages
COPY pkg/ ./pkg/
COPY proto/ ./proto/

# Copy service code
COPY auth-service/ ./auth-service/

# Build
WORKDIR /app/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# ── Runtime ──────────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget

COPY --from=builder /auth-service /usr/local/bin/auth-service

EXPOSE 8001 9001

CMD ["auth-service"]
```

---

## 12. Fase Implementasi (Timeline)

```mermaid
gantt
    title Fase 1 - Core Services Implementation
    dateFormat  YYYY-MM-DD

    section Foundation
    Setup project structure & Go modules   :f0, 2026-08-01, 3d
    Shared pkg (database, messaging, auth) :f1, after f0, 5d
    Proto definitions & code generation    :f2, after f0, 3d
    Docker Compose infra                   :f3, after f0, 2d

    section auth-service
    Domain models & migrations             :a1, after f1, 3d
    GORM repositories                      :a2, after a1, 4d
    Auth service logic (login, register)   :a3, after a2, 4d
    Microsoft OAuth integration            :a4, after a3, 3d
    REST handlers (Chi)                    :a5, after a3, 3d
    gRPC server                            :a6, after a5, 3d
    NATS event publisher                   :a7, after a6, 2d
    Unit tests                             :a8, after a7, 3d

    section master-data-service
    Domain models & migrations             :m1, after a5, 3d
    GORM repositories                      :m2, after m1, 4d
    Service logic                          :m3, after m2, 3d
    REST handlers + gRPC server            :m4, after m3, 4d
    gRPC client → auth-service             :m5, after m4, 2d
    Unit tests                             :m6, after m5, 2d

    section document-service
    Domain models & migrations             :d1, after m3, 2d
    MinIO storage integration              :d2, after d1, 3d
    Gotenberg PDF integration              :d3, after d2, 3d
    REST handlers + gRPC                   :d4, after d3, 3d
    Unit tests                             :d5, after d4, 2d

    section notification-service
    Domain models & migrations             :n1, after a7, 2d
    FCM + SMTP providers                   :n2, after n1, 4d
    NATS subscribers                       :n3, after n2, 3d
    REST handlers + gRPC                   :n4, after n3, 3d
    Unit tests                             :n5, after n4, 2d

    section API Gateway
    Routing + JWT middleware               :g1, after a6, 3d
    Rate limiting + CORS                   :g2, after g1, 2d

    section Integration
    End-to-end testing                     :i1, after n5, 5d
    Data migration script (existing → new) :i2, after i1, 4d
    Performance benchmarking               :i3, after i2, 3d
```

---

## 13. Data Migration dari Monolith

> [!WARNING]
> Karena menggunakan **Database per Service**, data dari MySQL existing perlu di-split ke 4 database terpisah. Ini perlu script migrasi khusus.

### Strategy:

```
┌──────────────────────────────────────┐
│     Existing: tka_db (single DB)     │
└─────────────┬────────────────────────┘
              │ Migration Script
              ▼
┌─────────────┬──────────────┬─────────────┬────────────┐
│neocentral_  │neocentral_   │neocentral_  │neocentral_ │
│auth         │master        │doc          │notif       │
│             │              │             │            │
│users,       │academic_     │documents,   │notifica-   │
│roles,       │years,        │doc_types,   │tions       │
│students,    │rooms, cpls,  │doc_templates│            │
│lecturers    │cpmks, etc    │             │            │
└─────────────┴──────────────┴─────────────┴────────────┘
```

Script migrasi akan:
1. Dump tabel-tabel dari `tka_db` berdasarkan ownership
2. Transform data jika perlu (e.g., UUID format)
3. Import ke database target masing-masing service
4. Validate data integrity (row counts, FK checks)

---

## Verification Plan

### Automated Tests
```bash
# Run all tests
make test

# Run specific service tests
make test-auth
make test-master-data
make test-document
make test-notification

# Integration tests (requires docker-compose up)
make docker-up
make test-integration
```

### Manual Verification
- Start semua services via `make docker-up`
- Test auth flow: register → login → refresh token
- Test CRUD master data (academic year, rooms, CPL, CPMK)
- Test document upload/download via MinIO
- Test notification delivery (in-app, email)
- Verify gRPC calls antar service menggunakan `grpcurl`
- Compare response format dengan Node.js existing
