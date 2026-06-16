# 🚀 DEALAN — Ride-Hailing Platform
## Dokumentasi Lengkap Proyek Microservices
### Backend Cloud Engineering Guide · Kelompok 22 Tubes

---

> **Dealan** adalah platform ride-hailing berbasis microservices cloud-native yang mendukung layanan transportasi (GoRide, GoCar) dan pengiriman barang (GoSend) dengan fitur unggulan negosiasi harga dua arah dan rating mutual user-driver.

---
## 1. Executive Summary

### 1.1 Gambaran Platform
Dealan dirancang untuk mengatasi gap di pasar ride-hailing Indonesia dengan menggabungkan kekuatan ekosistem super-app (seperti Gojek) dengan transparansi harga model negosiasi (seperti inDrive), sambil membangun di atas fondasi cloud-native yang benar-benar scalable.

### 1.2 Unique Selling Points
| Fitur | Gojek | inDrive | **Dealan** |
|---|---|---|---|
| Super App Ekosistem | ✅ | ❌ | ✅ |
| Negosiasi Harga | ❌ | ✅ | ✅ (dengan batas min/max) |
| Rating Dua Arah | Terbatas | ❌ | ✅ (Mutual Rating) |
| Punishment Otomatis | Manual | ❌ | ✅ (Automated) |
| Microservices Native | Sebagian | ❌ | ✅ (Full Cloud-Native) |

### 1.3 Target Performa Sistem
| Metrik | Target |
|---|---|
| Uptime | 99.9% (< 8.7 jam downtime/tahun) |
| Latency Matching | < 1 detik |
| Peak Load | 10.000+ concurrent users |
| Data Consistency | Strong consistency untuk Payment |

---

## 2. Tech Stack & Keputusan Arsitektur

### 2.1 Peta Teknologi per Layer
```text
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                              │
│              React Native (Mobile App)                      │
│              Next.js (Web Dashboard Admin)                  │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS / WSS
┌────────────────────────▼────────────────────────────────────┐
│                   API GATEWAY LAYER                          │
│         Kong Gateway (Auth, Rate Limit, Routing)            │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                MICROSERVICES LAYER                           │
│                 Golang (Framework Gin)                       │
└────────┬────────────────────────────────┬───────────────────┘
         │                                │
┌────────▼──────┐                 ┌──────▼──────────────────┐
│  DATABASE     │                 │  MESSAGE BROKER          │
│  LAYER        │                 │  LAYER                   │
│               │                 │                           │
│ PostgreSQL    │                 │ Apache Kafka             │
│ (w/ PostGIS)  │                 │ (Event-Driven Async)     │
└───────────────┘                 └───────────────────────────┘
```

### 2.2 Keputusan Database & Runtime per Service
Demi keseragaman, kemudahan deployment, dan menjamin transaksi ACID yang ketat, kelompok kami menyepakati penggunaan **Golang** untuk seluruh runtime dan **PostgreSQL** untuk seluruh basis data.

| Service | Runtime | Database | Alasan Teknis Adaptasi PostgreSQL & Golang |
|---|---|---|---|
| 1. Auth (Shakila) | Go | PostgreSQL | Golang sangat cepat memproses JWT. Data kredensial aman di tabel SQL relasional. |
| 2. User (Shakila) | Go | PostgreSQL | Profil pengguna sangat terstruktur, pas dengan skema tabel SQL konvensional. |
| 3. Driver (Shakila) | Go | PostgreSQL | Master data driver (KTP, SIM) direlasikan dengan Foreign Key secara ketat. |
| 4. Order (Zahid) | Go | PostgreSQL | Eksekusi pesanan super cepat dengan Go. Variasi data GoSend diatasi dengan kolom bertipe `JSONB`. |
| 5. Matching (Zahid) | Go | PostgreSQL | Go Goroutines menangani konkurensi. DB menggunakan ekstensi **PostGIS** untuk query koordinat terdekat. |
| 6. Pricing (Niha) | Go | PostgreSQL | Aturan harga dinamis disimpan di kolom `JSONB` agar fleksibel saat ada perubahan tarif. |
| 7. Shipment (Najmi) | Go | PostgreSQL | Detail dan manifes barang logistik disimpan secara terstruktur namun dinamis via `JSONB`. |
| 8. Location (Zahid) | Go | PostgreSQL | Pembaruan lokasi tiap 3 detik dikelola dengan Goroutines dan batch insert ke tabel PostGIS. |
| 9. Map Route (Niha) | Go | PostgreSQL | Data polyline jalan yang panjang disimpan dalam tipe data `TEXT` pada PostgreSQL. |
| 10. Payment (Najmi) | Go | PostgreSQL | Strong consistency (ACID) mutlak untuk dompet digital dan transaksi anti-double payment. |
| 11. Notification (Putri)| Go | PostgreSQL | Log riwayat notifikasi dipartisi (Table Partitioning) berdasarkan bulan agar tidak berat. |
| 12. Rating (Putri) | Go | PostgreSQL | Fungsi agregasi `AVG()` SQL bawaan sangat handal untuk menghitung skor rata-rata instan. |
| 13. Punishment (Putri) | Go | PostgreSQL | Log audit sanksi butuh relasi Foreign Key yang ketat demi keperluan legal. |
| 14. Promo (Putri) | Go | PostgreSQL | `SELECT ... FOR UPDATE` mengunci baris agar kuota promo tidak jebol akibat race condition. |
| 15. Chat (Niha) | Go | PostgreSQL | WebSocket dikelola natif oleh Go. Riwayat chat disimpan persisten ke tabel PostgreSQL. |

---

## 3. Database Schema per Service

### 3.1 Auth Service

CREATE TABLE auth_credentials (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id    UUID NOT NULL UNIQUE,         -- FK ke users atau drivers
    role          VARCHAR(10) NOT NULL CHECK (role IN ('user', 'driver', 'admin')),
    password_hash VARCHAR(255) NOT NULL,        -- bcrypt hash
    is_active     BOOLEAN DEFAULT true,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE otp_codes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nomor_hp    VARCHAR(20) NOT NULL,
    otp_code    VARCHAR(6) NOT NULL,
    purpose     VARCHAR(20) NOT NULL CHECK (purpose IN ('login', 'register', 'reset')),
    expires_at  TIMESTAMPTZ NOT NULL,
    is_used     BOOLEAN DEFAULT false,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_otp_nomor_hp ON otp_codes(nomor_hp);

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID NOT NULL,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,   -- Hash dari refresh token
    device_info JSONB,                          -- {"device": "iPhone 14", "os": "iOS 17"}
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_rt_account_id ON refresh_tokens(account_id);

-- Row Level Security
ALTER TABLE auth_credentials ENABLE ROW LEVEL SECURITY;
ALTER TABLE refresh_tokens ENABLE ROW LEVEL SECURITY;
```

### 3.2 User Service 


CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama          VARCHAR(100) NOT NULL,
    email         VARCHAR(100) UNIQUE,
    nomor_hp      VARCHAR(20) NOT NULL UNIQUE,
    alamat        TEXT,
    foto_profil   VARCHAR(500),                -- URL Cloud Storage
    status        VARCHAR(20) DEFAULT 'active'
                  CHECK (status IN ('active', 'suspended', 'banned')),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    preferensi      JSONB DEFAULT '{}',        -- {"notif_promo": true, "bahasa": "id"}
    total_order     INTEGER DEFAULT 0,
    total_spend     DECIMAL(15,2) DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_ratings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    total_rating    DECIMAL(3,2) DEFAULT 5.00,
    total_review    INTEGER DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 3.3 Driver Service 

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE drivers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama            VARCHAR(100) NOT NULL,
    nomor_hp        VARCHAR(20) NOT NULL UNIQUE,
    email           VARCHAR(100),
    foto_profil     VARCHAR(500),
    foto_ktp        VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'active'
                    CHECK (status IN ('active', 'suspended', 'banned')),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE vehicles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id       UUID NOT NULL REFERENCES drivers(id),
    jenis           VARCHAR(10) NOT NULL CHECK (jenis IN ('motor', 'mobil')),
    merek           VARCHAR(50),
    model           VARCHAR(50),
    plat_nomor      VARCHAR(20) NOT NULL UNIQUE,
    tahun           INTEGER,
    foto_kendaraan  VARCHAR(500),
    is_verified     BOOLEAN DEFAULT false
);

CREATE TABLE driver_status (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id       UUID NOT NULL REFERENCES drivers(id) UNIQUE,
    is_online       BOOLEAN DEFAULT false,
    layanan_aktif   VARCHAR(10)[] DEFAULT '{}', -- ['ride','send']
    current_order_id UUID,
    last_seen       TIMESTAMPTZ DEFAULT NOW(),
    -- PostGIS untuk query geospasial
    lokasi          GEOGRAPHY(POINT, 4326)
);
CREATE INDEX idx_driver_lokasi ON driver_status USING GIST(lokasi);

CREATE TABLE driver_ratings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id       UUID NOT NULL REFERENCES drivers(id) UNIQUE,
    total_rating    DECIMAL(3,2) DEFAULT 5.00,
    total_review    INTEGER DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 3.4 Order Service 

{
  _id: ObjectId(),
  order_id: "ORD-20240601-XXXXX",    // Custom readable ID
  user_id: "uuid",
  driver_id: "uuid",                  // null sampai matched
  
  service_type: "ride",               // ride | car | send
  
  lokasi_jemput: {
    type: "Point",
    coordinates: [106.8272, -6.1754], // [lng, lat]
    alamat: "Jl. Sudirman No. 1, Jakarta"
  },
  
  lokasi_tujuan: {
    type: "Point",
    coordinates: [106.8456, -6.2012],
    alamat: "Jl. Gatot Subroto No. 45, Jakarta"
  },
  
  // Khusus GoSend
  detail_paket: {
    kategori: "makanan",              // null jika bukan send
    berat_kg: 1.5,
    nama_penerima: "Budi Santoso",
    nomor_penerima: "08123456789",
    catatan: "Hati-hati mudah pecah"
  },
  
  // Pricing
  harga_estimasi: 25000,
  harga_negosiasi: 22000,             // null jika tidak negosiasi
  harga_final: 22000,
  
  metode_pembayaran: "qris",          // cash | qris | bank
  promo_id: null,
  diskon: 0,
  
  status: "completed",
  // Status flow: pending → matched → pickup → ongoing → completed | cancelled
  
  status_history: [
    { status: "pending",   timestamp: ISODate(), note: "" },
    { status: "matched",   timestamp: ISODate(), note: "Driver: Budi" },
    { status: "pickup",    timestamp: ISODate(), note: "" },
    { status: "ongoing",   timestamp: ISODate(), note: "" },
    { status: "completed", timestamp: ISODate(), note: "" }
  ],
  
  jarak_km: 3.7,
  durasi_menit: 18,
  
  created_at: ISODate(),
  updated_at: ISODate()
}

// Index untuk performa
db.orders.createIndex({ "lokasi_jemput": "2dsphere" });
db.orders.createIndex({ user_id: 1, status: 1 });
db.orders.createIndex({ driver_id: 1, status: 1 });
db.orders.createIndex({ created_at: -1 });
```

### 3.5 Matching Service

# Geo-index semua driver online per jenis layanan
# GEOADD drivers:online:ride <longitude> <latitude> <driver_id>
GEOADD drivers:online:ride  106.8272 -6.1754 "driver-uuid-001"
GEOADD drivers:online:ride  106.8301 -6.1821 "driver-uuid-002"
GEOADD drivers:online:car   106.8345 -6.1901 "driver-uuid-003"

# Query driver dalam radius 5km
# GEORADIUS drivers:online:ride 106.8272 -6.1754 5 km WITHDIST ASC COUNT 10

# Driver metadata (TTL 60 detik, driver harus update berkala)
# Hash per driver
HSET driver:driver-uuid-001  rating 4.8  is_available 1  vehicle_type motor
EXPIRE driver:driver-uuid-001 60

# Order matching lock (prevent double assignment)
SET matching:lock:order-id-xxx  driver-uuid-001  EX 30  NX

# Driver response tracking (timeout 30 detik per driver)
SET matching:pending:order-id-xxx:driver-uuid-001  "waiting"  EX 30
```

### 3.6 Location / Tracking Service 

{
  "locations": {
    "drivers": {
      "driver-uuid-001": {
        "lat": -6.1754,
        "lng": 106.8272,
        "bearing": 270,      // arah hadap, untuk animasi
        "speed": 25,         // km/h
        "is_online": true,
        "order_id": "ORD-20240601-XXXXX",
        "last_updated": 1717200000000  // Unix ms
      }
    },
    // Lokasi user yang sedang dalam perjalanan
    "users": {
      "user-uuid-001": {
        "lat": -6.1800,
        "lng": 106.8250,
        "last_updated": 1717200000000
      }
    }
  }
}

// ============================================
// Redis buffer untuk write batching
// ============================================
// Driver update lokasi setiap 3 detik
// Redis buffer dikuras ke Firebase setiap 1 detik
// Mengurangi Firebase write operations

// Key: loc:buffer:{driver_id}
HSET loc:buffer:driver-uuid-001  lat -6.1754  lng 106.8272  bearing 270
EXPIRE loc:buffer:driver-uuid-001 10
```

### 3.7 Payment Service 

CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id      VARCHAR(50) NOT NULL UNIQUE,  -- TRX-YYYYMMDD-XXXXX
    order_id            VARCHAR(50) NOT NULL,
    user_id             UUID NOT NULL,
    driver_id           UUID NOT NULL,
    
    gross_amount        DECIMAL(15,2) NOT NULL,
    discount_amount     DECIMAL(15,2) DEFAULT 0,
    net_amount          DECIMAL(15,2) NOT NULL,
    
    metode_pembayaran   VARCHAR(20) NOT NULL CHECK (metode_pembayaran IN ('cash','qris','bank')),
    payment_gateway_ref VARCHAR(100),               -- ID dari Midtrans/Xendit
    
    status              VARCHAR(20) DEFAULT 'pending'
                        CHECK (status IN ('pending','success','failed','refunded','expired')),
    
    -- Bagi hasil
    platform_fee        DECIMAL(15,2),              -- Potongan platform (20%)
    driver_earnings     DECIMAL(15,2),              -- Yang diterima driver (80%)
    
    paid_at             TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_transactions_order ON transactions(order_id);
CREATE INDEX idx_transactions_user  ON transactions(user_id);

CREATE TABLE driver_wallets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id       UUID NOT NULL UNIQUE,
    balance         DECIMAL(15,2) DEFAULT 0,
    total_earned    DECIMAL(15,2) DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE wallet_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id       UUID NOT NULL,
    type            VARCHAR(10) CHECK (type IN ('credit','debit')),
    amount          DECIMAL(15,2) NOT NULL,
    balance_after   DECIMAL(15,2) NOT NULL,
    ref_id          VARCHAR(50),                    -- transaction_id
    note            TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE payment_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  VARCHAR(50) NOT NULL,
    event           VARCHAR(50) NOT NULL,           -- created, webhook_received, etc.
    payload         JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Idempotency keys (mencegah double payment)
CREATE TABLE idempotency_keys (
    key             VARCHAR(255) PRIMARY KEY,
    response        JSONB NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### 3.8 Rating Service 
CREATE TABLE reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        VARCHAR(50) NOT NULL,
    
    -- Reviewer bisa user atau driver
    reviewer_id     UUID NOT NULL,
    reviewer_role   VARCHAR(10) NOT NULL CHECK (reviewer_role IN ('user','driver')),
    
    -- Target yang dirating
    target_id       UUID NOT NULL,
    target_role     VARCHAR(10) NOT NULL CHECK (target_role IN ('user','driver')),
    
    rating_score    SMALLINT NOT NULL CHECK (rating_score BETWEEN 1 AND 5),
    comment         TEXT,
    tags            VARCHAR(50)[],  -- ['tepat waktu', 'ramah', 'aman berkendara']
    
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(order_id, reviewer_id)  -- Satu review per order per reviewer
);

CREATE INDEX idx_reviews_target ON reviews(target_id, target_role);
CREATE INDEX idx_reviews_order  ON reviews(order_id);

-- Materialized view untuk agregasi rating (refresh berkala)
CREATE MATERIALIZED VIEW rating_aggregates AS
SELECT
    target_id,
    target_role,
    ROUND(AVG(rating_score)::numeric, 2) AS average_rating,
    COUNT(*) AS total_review,
    COUNT(CASE WHEN rating_score = 5 THEN 1 END) AS bintang_5,
    COUNT(CASE WHEN rating_score = 4 THEN 1 END) AS bintang_4,
    COUNT(CASE WHEN rating_score = 3 THEN 1 END) AS bintang_3,
    COUNT(CASE WHEN rating_score <= 2 THEN 1 END) AS bintang_rendah,
    MAX(created_at) AS last_reviewed_at
FROM reviews
GROUP BY target_id, target_role;

CREATE UNIQUE INDEX ON rating_aggregates(target_id, target_role);
```

### 3.9 Punishment Service 

CREATE TABLE violation_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL,
    account_role    VARCHAR(10) NOT NULL CHECK (account_role IN ('user','driver')),
    reason_code     VARCHAR(50) NOT NULL,
    -- CANCEL_ABUSE, LOW_RATING, FRAUD, HARASSMENT, LATE_PICKUP
    reason_detail   TEXT,
    reported_by     UUID,                       -- null jika auto-generated sistem
    order_id        VARCHAR(50),
    severity        VARCHAR(10) DEFAULT 'low' CHECK (severity IN ('low','medium','high')),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE penalty_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL,
    account_role    VARCHAR(10) NOT NULL,
    penalty_type    VARCHAR(20) NOT NULL CHECK (penalty_type IN ('warning','suspended','banned')),
    reason_codes    VARCHAR(50)[],
    duration_hours  INTEGER,                    -- null jika banned permanen
    starts_at       TIMESTAMPTZ DEFAULT NOW(),
    ends_at         TIMESTAMPTZ,               -- null jika banned permanen
    is_active       BOOLEAN DEFAULT true,
    lifted_by       UUID,                       -- Admin yang mencabut sanksi
    lifted_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_penalty_account ON penalty_history(account_id, is_active);

-- Rules otomatis punishment
CREATE TABLE punishment_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reason_code     VARCHAR(50) NOT NULL,
    threshold       INTEGER NOT NULL,           -- misal: 3 kali cancel = suspend
    penalty_type    VARCHAR(20) NOT NULL,
    duration_hours  INTEGER,
    window_days     INTEGER DEFAULT 30,         -- hitung dalam 30 hari terakhir
    is_active       BOOLEAN DEFAULT true
);

-- Seed default rules
INSERT INTO punishment_rules (reason_code, threshold, penalty_type, duration_hours, window_days) VALUES
('CANCEL_ABUSE', 3, 'suspended', 24, 7),
('LOW_RATING',   1, 'warning', null, 30),
('LOW_RATING',   3, 'suspended', 48, 30),
('FRAUD',        1, 'banned', null, null),
('HARASSMENT',   1, 'banned', null, null);
```

### 3.10 Chat Service

{
  "chats": {
    "{order_id}": {
      "metadata": {
        "user_id": "uuid",
        "driver_id": "uuid",
        "status": "active",      // active | closed
        "created_at": 1717200000000
      },
      "messages": {
        "{message_id}": {
          "sender_id": "uuid",
          "sender_role": "user",  // user | driver
          "message": "Sudah di lokasi penjemputan ya",
          "type": "text",         // text | image | location_share
          "sent_at": 1717200000000,
          "read_status": false
        }
      }
    }
  },
  
  "typing_indicators": {
    "{order_id}": {
      "user-uuid": false,
      "driver-uuid": true
    }
  }
}


```

---

## 4. API Contract & Kafka Events

### 4.1 Kafka Topics
```json
// Topic: order.created
{
  "event": "ORDER_CREATED",
  "payload": {
    "order_id": "ORD-20240601-XXXXX",
    "user_id": "uuid",
    "service_type": "ride"
  }
}

// Topic: payment.completed
{
  "event": "PAYMENT_COMPLETED",
  "payload": {
    "transaction_id": "TRX-20240601-XXXXX",
    "order_id": "ORD-20240601-XXXXX",
    "status": "success"
  }
}
```

### 4.2 REST API Endpoints (Golang Gin)
```text
# AUTH SERVICE  — Port 3001
POST   /auth/register                 # Daftar akun baru
POST   /auth/login                    # Login & dapat JWT

# ORDER SERVICE — Port 3004
POST   /orders                        # Buat pesanan
POST   /orders/:order_id/negotiate    # Ajukan negosiasi harga

# PAYMENT SERVICE — Port 3010
POST   /payments/create               # Buat transaksi
GET    /payments/:transaction_id      # Status transaksi
```

---

## 5. Konfigurasi Infrastructure

### 5.1 Docker Compose (Local Development)
```yaml
version: '3.9'
services:
  # ============ DATABASES ============
  postgres:
    image: postgis/postgis:15-3.3-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: dealan
      POSTGRES_PASSWORD: dealan_secret
      POSTGRES_DB: dealan_db
    volumes:
      - pg_data:/var/lib/postgresql/data

  # ============ MESSAGE BROKER ============
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092

  # ============ API GATEWAY ============
  kong:
    image: kong:3.5-alpine
    ports:
      - "8000:8000"
    environment:
      KONG_DATABASE: "off"
      KONG_DECLARATIVE_CONFIG: /kong/kong.yml

  # ============ MICROSERVICES ============
  order-service:
    build: ./services/order-service
    ports:
      - "3004:3004"
    environment:
      - DB_URL=postgres://dealan:dealan_secret@postgres:5432/dealan_db
      - KAFKA_BROKERS=kafka:9092
    depends_on: [postgres, kafka]

volumes:
  pg_data:
```

---

## 5. API Contract antar Service

### 5.1 Kafka Event Schema

```javascript
// ============================================
// KAFKA TOPICS & EVENT SCHEMAS
// ============================================

// Topic: order.created
{
  "event": "ORDER_CREATED",
  "version": "1.0",
  "timestamp": "2024-06-01T10:00:00Z",
  "payload": {
    "order_id": "ORD-20240601-XXXXX",
    "user_id": "uuid",
    "service_type": "ride",
    "lokasi_jemput": { "lat": -6.1754, "lng": 106.8272 },
    "lokasi_tujuan": { "lat": -6.2012, "lng": 106.8456 }
  }
}

// Topic: order.matched
{
  "event": "ORDER_MATCHED",
  "payload": {
    "order_id": "ORD-20240601-XXXXX",
    "driver_id": "uuid",
    "estimasi_waktu_menit": 5
  }
}

// Topic: ride.completed
{
  "event": "RIDE_COMPLETED",
  "payload": {
    "order_id": "ORD-20240601-XXXXX",
    "user_id": "uuid",
    "driver_id": "uuid",
    "harga_final": 22000,
    "jarak_km": 3.7
  }
}

// Topic: payment.completed
{
  "event": "PAYMENT_COMPLETED",
  "payload": {
    "transaction_id": "TRX-20240601-XXXXX",
    "order_id": "ORD-20240601-XXXXX",
    "status": "success",
    "driver_earnings": 17600
  }
}

// Topic: rating.submitted
{
  "event": "RATING_SUBMITTED",
  "payload": {
    "review_id": "uuid",
    "target_id": "uuid",
    "target_role": "driver",
    "rating_score": 2,
    "order_id": "ORD-20240601-XXXXX"
  }
}
```

### 5.2 REST API Endpoints Utama

```
# ============================================
# AUTH SERVICE  — Port 3001
# ============================================
POST   /auth/register                 # Daftar akun baru
POST   /auth/login                    # Login dengan email/HP + password
POST   /auth/otp/send                 # Kirim OTP ke nomor HP
POST   /auth/otp/verify               # Verifikasi OTP
POST   /auth/token/refresh            # Refresh JWT token
POST   /auth/logout                   # Revoke refresh token
POST   /auth/password/reset           # Reset password via OTP

# ============================================
# ORDER SERVICE — Port 3004
# ============================================
POST   /orders                        # Buat pesanan baru
GET    /orders/:order_id              # Detail pesanan
GET    /orders/user/:user_id          # Riwayat pesanan user
PATCH  /orders/:order_id/status       # Update status (driver accept, dll)
POST   /orders/:order_id/negotiate    # Ajukan negosiasi harga
POST   /orders/:order_id/cancel       # Batalkan pesanan

# ============================================
# MATCHING SERVICE — Port 3005 (gRPC :50051)
# ============================================
POST   /matching/find                 # Cari driver terdekat
POST   /matching/assign               # Assign driver ke order
POST   /matching/decline              # Driver tolak order

# ============================================
# PRICING SERVICE — Port 3006
# ============================================
POST   /pricing/estimate              # Estimasi harga
POST   /pricing/negotiate             # Submit penawaran harga
GET    /pricing/rules/:service_type   # Lihat aturan harga

# ============================================
# PAYMENT SERVICE — Port 3010
# ============================================
POST   /payments/create               # Buat transaksi
POST   /payments/webhook              # Webhook dari payment gateway
GET    /payments/:transaction_id      # Status transaksi
GET    /payments/driver/:driver_id/wallet  # Saldo dompet driver

# ============================================
# RATING SERVICE — Port 3012
# ============================================
POST   /ratings                       # Submit rating
GET    /ratings/user/:user_id         # Rating user
GET    /ratings/driver/:driver_id     # Rating driver
GET    /ratings/order/:order_id       # Rating untuk order ini
```
---

## 6. Testing (Golang)

Setiap service wajib dites menggunakan `testing` bawaan Go dan `testify` untuk asersi.

```go
// services/order-service/tests/order_test.go
package tests

import (
    "testing"
    "[github.com/stretchr/testify/assert](https://github.com/stretchr/testify/assert)"
    "[github.com/dealan/order-service/internal/service](https://github.com/dealan/order-service/internal/service)"
)

func TestCreateOrder(t *testing.T) {
    t.Run("✅ Order berhasil dibuat", func(t *testing.T) {
        payload := service.OrderPayload{
            UserID:      "user-001",
            ServiceType: "ride",
        }
        res, err := service.CreateOrder(payload)
        
        assert.NoError(t, err)
        assert.Equal(t, "pending", res.Status)
    })
}
```

---

## 7. CI/CD Pipeline Jenkins & Dockerfile

### 7.1 Jenkinsfile
```groovy
pipeline {
    agent { docker { image 'golang:1.22-alpine' } }
    environment {
        SERVICE_NAME = 'order-service'
        IMAGE_TAG    = "${env.GIT_COMMIT?.take(7) ?: 'latest'}"
    }
    stages {
        stage('📦 Test & Build') {
            steps {
                dir("services/${SERVICE_NAME}") {
                    sh 'go mod download'
                    sh 'go test ./... -v'
                    sh 'CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go'
                }
            }
        }
        stage('🐳 Push Docker') {
            steps {
                sh """
                    docker build -t gcr.io/dealan-project/${SERVICE_NAME}:${IMAGE_TAG} services/${SERVICE_NAME}/
                    docker push gcr.io/dealan-project/${SERVICE_NAME}:${IMAGE_TAG}
                """
            }
        }
    }
}
```

### 7.2 Dockerfile Golang
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.18
RUN addgroup -S dealan && adduser -S dealan -G dealan
USER dealan
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3004
CMD ["./main"]
```
