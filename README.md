# NeoCentral Microservices (Go)

Proyek ini adalah migrasi sistem informasi NeoCentral dari arsitektur *Monolith* (Node.js/Express) menuju arsitektur **Microservices** menggunakan Golang. 

Arsitektur microservices ini mengadopsi pola **Monorepo**, di mana seluruh service (Auth, Master Data, Gateway, dll) berada dalam satu repositori, namun memiliki struktur direktori, `Makefile`, dan *environment variables* yang terisolasi.

---

## 🏗️ Struktur Proyek

```text
neocentral-go/
├── api-gateway/         # Single Entry Point (REST API Proxy, CORS, Rate Limiting)
├── auth-service/        # Microservice khusus untuk Autentikasi (REST & gRPC)
├── master-data-service/ # Microservice Data Induk (Tahun Ajaran, CPL, Ruangan) [WIP]
├── pkg/                 # Kode/Library yang di-share antar service (Database, Auth, AppError)
├── proto/               # File Protocol Buffers (.proto) dan generated gRPC code
├── docker-compose.yml   # Infrastruktur (MySQL, Redis, NATS)
└── Makefile             # Command global untuk ekosistem (Build proto, tidy, dsb)
```

## 🔄 Alur Komunikasi

1. **Client Request:** Aplikasi Frontend (React/Next.js) atau Mobile mengirimkan HTTP Request (REST) ke **API Gateway** (Port `8080`).
2. **API Gateway (Proxy):** Gateway menerapkan CORS dan meneruskan request tersebut ke microservice yang bersangkutan (misal: `/api/v1/auth/*` diteruskan ke `auth-service` di port `8001`).
3. **Internal Microservices (REST & gRPC):**
   - Tiap service dapat melayani request REST (Port `800x`).
   - Jika satu service butuh data dari service lain secara internal, mereka akan berkomunikasi menggunakan **gRPC** (Port `900x`) yang jauh lebih cepat daripada REST.
4. **Asynchronous Messaging:** Operasi *background* (seperti mengirim email/notifikasi) dilakukan menggunakan message broker **NATS JetStream**.

---

## 🛠️ Prasyarat (Prerequisites) & Instalasi Tools

Sebelum menjalankan proyek ini, pastikan sistem Anda (khususnya Windows) sudah terinstal *tools* berikut:

### 1. Golang (Min. v1.23)
- Unduh dan instal dari web resmi: [https://go.dev/dl/](https://go.dev/dl/)
- Pastikan Go sudah masuk ke dalam PATH (`go version`).

### 2. Make (Untuk Windows)
Karena ekosistem proyek ini sangat bergantung pada `Makefile`, pengguna Windows wajib menginstal `make`.
- **Menggunakan Chocolatey:** Buka PowerShell sebagai Administrator, lalu jalankan: 
  `choco install make`
- **Menggunakan Scoop:** `scoop install make`

### 3. Docker & Docker Compose
Infrastruktur seperti MySQL, Redis, dan NATS dibungkus dalam Docker.
- Instal **Docker Desktop** untuk Windows: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
- Pastikan Docker Desktop dalam keadaan menyala (Running) sebelum melanjutkan.

### 4. Golang Migrate CLI
Tool ini digunakan untuk menjalankan script migrasi database (`make migrate-up`).
- Buka terminal, lalu jalankan perintah ini untuk menginstal versi terbaru dengan dukungan MySQL:
  ```bash
  go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```
- *Catatan:* Pastikan path `C:\Users\<NamaUser>\go\bin` sudah terdaftar di *Environment Variables (PATH)* Windows Anda.

### 5. Protobuf Compiler (Opsional - Jika ingin modifikasi gRPC)
Jika Anda berencana mengubah file `.proto`, Anda perlu menginstal `protoc` dan plugin golang-nya.
1. Unduh `protoc` untuk Windows dari [Releases Protocol Buffers](https://github.com/protocolbuffers/protobuf/releases), ekstrak, dan masukkan folder `/bin`-nya ke PATH Windows.
2. Instal plugin Go:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

---

## 🚀 Tutorial Menjalankan Proyek (Getting Started)

### 1. Menjalankan Infrastruktur Database & Broker
Proyek ini membutuhkan MySQL, Redis, dan NATS. Semua sudah dikonfigurasi di dalam Docker Compose.
Buka terminal di root `neocentral-go` dan jalankan:
```bash
make docker-up
```
*(Tunggu beberapa saat hingga seluruh container menyala).*

### 2. Setup Environment Variables
Setiap microservice dan API Gateway memiliki file `.env.example`.
Anda harus meng-copy file tersebut menjadi `.env` di masing-masing folder:
```bash
# Untuk Auth Service
cp auth-service/.env.example auth-service/.env

# Untuk API Gateway
cp api-gateway/.env.example api-gateway/.env
```

### 3. Migrasi Database (Auth Service)
Karena *database-per-service* diterapkan, setiap service memiliki migrasi mandirinya sendiri.
Untuk memigrasi tabel *users* pada `auth-service`:
```bash
cd auth-service
make migrate-up
```
*(Catatan: Ini akan membuat tabel users dan mengisi data seed secara otomatis).*

### 4. Menjalankan Service

Untuk kemudahan pengembangan, Anda perlu menjalankan service di terminal terpisah.

**Terminal 1 (Menjalankan Auth Service):**
```bash
cd auth-service
make run
```
*Service ini akan berjalan melayani REST (Port 8001) dan gRPC (Port 9001).*

**Terminal 2 (Menjalankan API Gateway):**
```bash
cd api-gateway
make run
```
*Gateway akan berjalan di port 8080 dan mem-proxy request Anda ke service yang sesuai.*

---

## 🧪 Pengetesan & Penggunaan (Testing)

Anda dapat langsung melakukan pengetesan endpoints menggunakan **Swagger UI**.
Karena API Gateway sudah berjalan, Anda dapat mengakses Swagger melalui browser:
👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

Cobalah jalankan request login pada endpoint `/auth/login`. Jika berhasil, Anda tidak perlu lagi menambahkan awalan `Bearer ` saat menempelkan *AccessToken* di tombol **Authorize** Swagger.

---

## 📜 Perintah Berguna (Useful Commands)

**Di Root (`/neocentral-go/`):**
- `make docker-up`: Menyalakan infrastruktur background.
- `make docker-down`: Mematikan infrastruktur background.
- `make proto`: Meng-generate ulang kode gRPC jika Anda mengubah file `.proto`.
- `make tidy`: Merapikan seluruh dependensi `go.mod`.

**Di dalam Service (misal: `/auth-service/`):**
- `make run`: Menjalankan service spesifik.
- `make build`: Membangun file binary service (`/bin/`).
- `make migrate-up`: Menjalankan script database migration ke atas.
- `make migrate-down`: Merollback 1 step script migration.
