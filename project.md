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

```
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
│    Node.js/Express · Python/FastAPI · Go (per service)      │
└────────┬───────────────┬────────────────┬───────────────────┘
         │               │                │
┌────────▼──────┐ ┌──────▼──────┐ ┌──────▼──────────────────┐
│  DATABASE     │ │  CACHE &    │ │  MESSAGE BROKER          │
│  LAYER        │ │  REALTIME   │ │  LAYER                   │
│               │ │             │ │                           │
│ Supabase(PG)  │ │ Redis       │ │ Apache Kafka             │
│ MongoDB       │ │ Upstash     │ │ (Event-Driven Async)     │
│ Firebase RTDB │ │             │ │                           │
└───────────────┘ └─────────────┘ └───────────────────────────┘
```

### 2.2 Keputusan Database per Service

Pemilihan database didasarkan pada karakteristik data dan pola akses masing-masing service:

ServiceDatabase PilihanAlasan Teknis & Strategi Implementasi di Dunia Nyata1. Authentication Service (Shakila)Supabase (Postgres)Mutlak. Menggunakan fitur native Supabase Auth (JWT, hashing bcrypt otomatis, dan manajemen session token). Aman dan memangkas waktu coding 2 minggu.2. User Service (Shakila)Supabase (Postgres)Data profil tabular (Nama, Email, Telp) sangat terstruktur. Sangat diuntungkan oleh fitur Row-Level Security (RLS) Supabase agar user tidak bisa mengintip data profil user lain.3. Driver Service (Shakila)Supabase (Postgres)Data master driver (KTP, SIM, STNK) bersifat relasional. Dokumen fisik ditarik dan disimpan ke dalam Supabase Storage, lalu link URL-nya dicatat di tabel Postgres ini.4. Order Service (Zahid)MongoDBWajib. Siklus hidup (lifecycle) pesanan ride-hailing itu sangat dinamis. Di MongoDB, struktur data orderan ojek (Ride) dan data orderan kirim paket (Send) yang punya atribut berat/dimensi barang bisa disimpan dalam satu collection tanpa merusak skema.5. Matching Service (Zahid)Supabase (Postgres)Karena kita tidak pakai Redis, kita aktifkan Ekstensi PostGIS di Postgres Supabase. Kita bisa menggunakan query ST_DWithin atau ST_Distance untuk mencari koordinat driver_id terdekat dari titik jemput penumpang secara akurat.6. Pricing Service (Niha)MongoDBAturan tarif (zona macet, tarif per KM, biaya admin, tarif peak-hour) berubah-ubah bentuknya berdasarkan jenis layanan. Format dokumen (BSON/JSON) MongoDB mempermudah perubahan aturan pricing tanpa perlu melakukan migrasi tabel.7. Shipment Service (Najmi)MongoDBSangat cocok menggunakan MongoDB karena manifes paket kiriman memiliki detail yang bervariasi (misal: barang pecah belah butuh catatan khusus, makanan butuh penanganan suhu). Struktur nested array MongoDB sangat natural untuk ini.8. Location/Tracking Service (Zahid)Supabase (Postgres)Karena handphone pengemudi menembak koordinat GPS setiap 3 detik, gunakan tabel khusus di Supabase dengan optimasi kolom bertipe data POINT (PostGIS) dan matikan fitur indeks yang tidak perlu agar write throughput-nya kencang.9. Map & Route Service (Niha)MongoDBKarena kita tidak memakai Redis sebagai cache, rute koordinat jalan (polyline rute komplit yang panjang dan berat jika dihitung ulang) bisa di-cache ke dalam MongoDB sebagai dokumen dokumen JSON statis berdasarkan ID rute asal-tujuan.10. Payment Service (Najmi)Supabase (Postgres)Mutlak. Urusan dompet (wallet), invoice, dan transaksi uang tidak boleh ada toleransi salah hitung. Dukungan transaksi ACID di PostgreSQL menjamin data konsisten, anti-double payment, dan memiliki rollback aman jika koneksi terputus.11. Notification Service (Putri)MongoDBLog sejarah notifikasi (Push Notification, SMS, Email) bervolume sangat besar tetapi tidak bernilai transaksi finansial. MongoDB sangat kuat menampung jutaan dokumen log teks notifikasi tanpa membebani database utama.12. Rating Service (Putri)Supabase (Postgres)Query SQL sangat digdaya untuk agregasi data. Menghitung rata-rata rating driver (misal: SELECT AVG(rating_score) FROM ratings WHERE driver_id = ...) berjalan instan dengan fungsi bawaan SQL.13. Punishment Service (Putri)Supabase (Postgres)Memerlukan relasi asing (Foreign Key) yang ketat ke tabel User atau Driver. Log pembekuan akun (suspend) membutuhkan audit trail terstruktur demi urusan legal/hukum perusahaan.14. Promo Service (Putri)Supabase (Postgres)Pengecekan kuota promo (misal: "Kode promo ini hanya boleh dipakai 1000 kali") membutuhkan fitur Database Lock (SELECT ... FOR UPDATE) milik Postgres agar tidak terjadi race condition (kuota jebol karena diklik berbarengan oleh ribuan user).15. Chat Service (Niha)Supabase (Postgres)Gunakan Fitur Supabase Realtime. Supabase memiliki kapabilitas WebSocket bawaan di atas PostgreSQL. Setiap kali ada baris chat baru masuk ke tabel chat_messages, Supabase otomatis menyiarkan (broadcast) pesan tersebut ke handphone driver dan penumpang secara real-time. Timmu tidak perlu ngoding WebSocket dari nol!

### 2.3 Runtime per Service

| Service | Runtime | Framework | Alasan |
|---|---|---|---|
| Auth Service | **Node.js** | Express + Passport | Ekosistem JWT/OAuth2 matang |
| User Service | **Node.js** | Express | Konsistensi dengan Auth |
| Driver Service | **Node.js** | Express | Konsistensi dengan User |
| Order Service | **Node.js** | Express | Event-driven native dengan EventEmitter |
| Matching Service | **Go** | Gin | Performa tinggi, goroutine untuk concurrent matching |
| Pricing Service | **Python** | FastAPI | Kalkulasi numerik, potensi ML untuk dynamic pricing |
| Shipment Service | **Node.js** | Express | Integrasi Cloud Storage mudah |
| Location Service | **Go** | Gin + WebSocket | Throughput tinggi (10K rps), goroutine per connection |
| Map & Route Service | **Python** | FastAPI | Integrasi Google Maps SDK Python |
| Payment Service | **Node.js** | Express | Integrasi payment gateway (Midtrans) |
| Notification Service | **Node.js** | Express | Firebase Admin SDK |
| Rating Service | **Node.js** | Express | Kalkulasi sederhana |
| Punishment Service | **Node.js** | Express | Integrasi dengan Auth untuk suspend |
| Promo Service | **Node.js** | Express | Redis client native |
| Chat Service | **Node.js** | Socket.io | Real-time WebSocket native |

*framework dan runtime dapat disesuaikan sesuai kebutuhan, namun harus konsisten dalam penggunaannya di project ini dan diutamakan memakai go*

---


## 4. Database Schema per Service

### 4.1 Auth Service

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

### 4.2 User Service 


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

### 4.3 Driver Service 

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

### 4.4 Order Service 

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

### 4.5 Matching Service

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

### 4.6 Location / Tracking Service 

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

### 4.7 Payment Service 

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

### 4.8 Rating Service 
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

### 4.9 Punishment Service 

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

### 4.10 Chat Service

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

## 6. Konfigurasi Infrastructure

### 6.1 Docker Compose (Local Development)

```yaml
# docker-compose.yml
version: '3.9'

services:
  # ============ DATABASES ============
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data

  mongodb:
    image: mongo:7
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: dealan
      MONGO_INITDB_ROOT_PASSWORD: dealan_secret
    volumes:
      - mongo_data:/data/db

  # ============ MESSAGE BROKER ============
  
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports:
      - "8090:8080"
    depends_on:
      - kafka
    environment:
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:9092

  # ============ API GATEWAY ============
  
  kong:
    image: kong:3.5-alpine
    ports:
      - "8000:8000"   # Proxy
      - "8001:8001"   # Admin API
    environment:
      KONG_DATABASE: "off"
      KONG_DECLARATIVE_CONFIG: /kong/kong.yml
    volumes:
      - ./api-gateway/kong.yml:/kong/kong.yml

  # ============ MICROSERVICES ============

  auth-service:
    build: ./services/auth-service
    ports:
      - "3001:3001"
    environment:
      - NODE_ENV=development
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_KEY=${SUPABASE_KEY}
      - JWT_SECRET=${JWT_SECRET}
      - REDIS_URL=redis://redis:6379
    depends_on:
      - redis

  order-service:
    build: ./services/order-service
    ports:
      - "3004:3004"
    environment:
      - MONGODB_URI=mongodb://dealan:dealan_secret@mongodb:27017/dealan_orders
      - KAFKA_BROKERS=kafka:9092
      - REDIS_URL=redis://redis:6379
    depends_on:
      - mongodb
      - kafka

  matching-service:
    build: ./services/matching-service
    ports:
      - "3005:3005"
    environment:
      - REDIS_URL=redis://redis:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - redis
      - kafka

  location-service:
    build: ./services/location-service
    ports:
      - "3008:3008"
    environment:
      - REDIS_URL=redis://redis:6379
      - FIREBASE_CRED=/secrets/firebase.json
    volumes:
      - ./secrets/firebase.json:/secrets/firebase.json

  payment-service:
    build: ./services/payment-service
    ports:
      - "3010:3010"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_KEY=${SUPABASE_KEY}
      - MIDTRANS_SERVER_KEY=${MIDTRANS_SERVER_KEY}
      - KAFKA_BROKERS=kafka:9092

  notification-service:
    build: ./services/notification-service
    ports:
      - "3011:3011"
    environment:
      - KAFKA_BROKERS=kafka:9092
      - FCM_SERVER_KEY=${FCM_SERVER_KEY}

volumes:
  redis_data:
  mongo_data:
```

### 6.2 Kong API Gateway Config

```yaml
# api-gateway/kong.yml
_format_version: "3.0"

services:
  - name: auth-service
    url: http://auth-service:3001
    routes:
      - name: auth-routes
        paths: ["/auth"]
        strip_path: false
    plugins:
      - name: rate-limiting
        config:
          minute: 20
          hour: 500
          policy: redis
          redis_host: redis

  - name: order-service
    url: http://order-service:3004
    routes:
      - name: order-routes
        paths: ["/orders"]
        strip_path: false
    plugins:
      - name: jwt
        config:
          secret_is_base64: false
      - name: rate-limiting
        config:
          minute: 60
          policy: local

  - name: location-service
    url: http://location-service:3008
    routes:
      - name: location-routes
        paths: ["/location"]
        strip_path: false
    plugins:
      - name: jwt

  - name: payment-service
    url: http://payment-service:3010
    routes:
      - name: payment-routes
        paths: ["/payments"]
        strip_path: false
    plugins:
      - name: jwt
      - name: ip-restriction
        config:
          allow: ["10.0.0.0/8"]  # Hanya internal untuk webhook
          # Webhook endpoint dibuka via route terpisah

plugins:
  - name: cors
    config:
      origins: ["https://dealan.id", "http://localhost:3000"]
      methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
      headers: ["Authorization", "Content-Type"]
      credentials: true
```

### 6.3 Kubernetes Deployment (Production)

```yaml
# infrastructure/k8s/deployments/order-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: dealan-production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
        - name: order-service
          image: gcr.io/dealan-project/order-service:latest
          ports:
            - containerPort: 3004
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          env:
            - name: NODE_ENV
              value: "production"
            - name: MONGODB_URI
              valueFrom:
                secretKeyRef:
                  name: dealan-secrets
                  key: mongodb-uri
          livenessProbe:
            httpGet:
              path: /health
              port: 3004
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 3004
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: order-service
  namespace: dealan-production
spec:
  selector:
    app: order-service
  ports:
    - port: 3004
      targetPort: 3004
  type: ClusterIP
---
# HPA — Auto-scaling berdasarkan CPU
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: order-service-hpa
  namespace: dealan-production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: order-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

---

## 7. Unit & Functional Testing

### 7.1 Auth Service Tests (Jest)

```javascript
// services/auth-service/tests/unit/auth.service.test.js
const AuthService = require('../../src/services/auth.service');
const { createClient } = require('@supabase/supabase-js');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');

jest.mock('@supabase/supabase-js');
jest.mock('bcrypt');
jest.mock('jsonwebtoken');

describe('AuthService - Unit Tests', () => {
  
  // ==========================================
  // REGISTER
  // ==========================================
  describe('register()', () => {
    test('✅ Berhasil register pengguna baru', async () => {
      const mockUser = { id: 'uuid-001', nomor_hp: '08123456789', role: 'user' };
      bcrypt.hash.mockResolvedValue('hashed_password');
      
      const supabaseMock = {
        from: jest.fn().mockReturnThis(),
        insert: jest.fn().mockReturnThis(),
        select: jest.fn().mockResolvedValue({ data: [mockUser], error: null })
      };
      createClient.mockReturnValue(supabaseMock);
      
      const result = await AuthService.register({
        nomor_hp: '08123456789',
        password: 'Password123!',
        role: 'user'
      });
      
      expect(result.success).toBe(true);
      expect(result.data.nomor_hp).toBe('08123456789');
    });
    
    test('❌ Gagal jika nomor HP sudah terdaftar', async () => {
      const supabaseMock = {
        from: jest.fn().mockReturnThis(),
        insert: jest.fn().mockReturnThis(),
        select: jest.fn().mockResolvedValue({
          data: null,
          error: { code: '23505', message: 'duplicate key' }
        })
      };
      createClient.mockReturnValue(supabaseMock);
      
      await expect(
        AuthService.register({ nomor_hp: '08123456789', password: 'Pass123!' })
      ).rejects.toThrow('Nomor HP sudah terdaftar');
    });
  });
  
  // ==========================================
  // LOGIN
  // ==========================================
  describe('login()', () => {
    test('✅ Berhasil login dan mendapat JWT token', async () => {
      bcrypt.compare.mockResolvedValue(true);
      jwt.sign.mockReturnValue('mock.jwt.token');
      
      const result = await AuthService.login({
        nomor_hp: '08123456789',
        password: 'Password123!'
      });
      
      expect(result.access_token).toBeDefined();
      expect(result.role).toBe('user');
    });
    
    test('❌ Gagal login jika password salah', async () => {
      bcrypt.compare.mockResolvedValue(false);
      
      await expect(
        AuthService.login({ nomor_hp: '08123456789', password: 'SalahPass' })
      ).rejects.toThrow('Kredensial tidak valid');
    });
  });
  
  // ==========================================
  // OTP
  // ==========================================
  describe('verifyOTP()', () => {
    test('✅ OTP valid dan belum expired', async () => {
      const mockOTP = {
        otp_code: '123456',
        expires_at: new Date(Date.now() + 300000).toISOString(),  // +5 menit
        is_used: false
      };
      // Mock DB return
      const result = await AuthService.verifyOTP('08123456789', '123456');
      expect(result.success).toBe(true);
    });
    
    test('❌ OTP sudah expired', async () => {
      const mockOTP = {
        otp_code: '123456',
        expires_at: new Date(Date.now() - 60000).toISOString(),   // -1 menit (expired)
        is_used: false
      };
      await expect(
        AuthService.verifyOTP('08123456789', '123456')
      ).rejects.toThrow('OTP sudah kadaluarsa');
    });
  });
});
```

```javascript
// services/auth-service/tests/functional/auth.api.test.js
const request = require('supertest');
const app = require('../../src/app');

describe('Auth API - Functional Tests', () => {
  
  describe('POST /auth/register', () => {
    test('✅ 201 - Register berhasil', async () => {
      const res = await request(app)
        .post('/auth/register')
        .send({ nomor_hp: '081234567890', password: 'Test1234!', role: 'user' });
      
      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('data.user_id');
      expect(res.body).not.toHaveProperty('data.password');
    });
    
    test('❌ 422 - Validasi field wajib gagal', async () => {
      const res = await request(app)
        .post('/auth/register')
        .send({ nomor_hp: '08123' });  // password kurang, nomor tidak valid
      
      expect(res.status).toBe(422);
      expect(res.body.errors).toBeDefined();
    });
    
    test('❌ 409 - Conflict jika nomor HP duplikat', async () => {
      // Register pertama kali
      await request(app).post('/auth/register')
        .send({ nomor_hp: '081299998888', password: 'Test1234!', role: 'user' });
      
      // Register kedua dengan nomor sama
      const res = await request(app).post('/auth/register')
        .send({ nomor_hp: '081299998888', password: 'Test1234!', role: 'user' });
      
      expect(res.status).toBe(409);
    });
  });
  
  describe('POST /auth/login', () => {
    test('✅ 200 - Login sukses dengan token', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ nomor_hp: '081234567890', password: 'Test1234!' });
      
      expect(res.status).toBe(200);
      expect(res.body).toHaveProperty('access_token');
      expect(res.body).toHaveProperty('refresh_token');
      expect(res.body.token_type).toBe('Bearer');
    });
    
    test('❌ 401 - Password salah', async () => {
      const res = await request(app)
        .post('/auth/login')
        .send({ nomor_hp: '081234567890', password: 'WrongPass!' });
      
      expect(res.status).toBe(401);
    });
    
    test('❌ 429 - Rate limit setelah 5 kali gagal', async () => {
      for (let i = 0; i < 5; i++) {
        await request(app).post('/auth/login')
          .send({ nomor_hp: '081234567890', password: 'wrong' });
      }
      
      const res = await request(app).post('/auth/login')
        .send({ nomor_hp: '081234567890', password: 'wrong' });
      
      expect(res.status).toBe(429);
    });
  });
  
  describe('GET /auth/health', () => {
    test('✅ 200 - Health check', async () => {
      const res = await request(app).get('/auth/health');
      expect(res.status).toBe(200);
      expect(res.body.status).toBe('healthy');
    });
  });
});
```

### 7.2 Order Service Tests

```javascript
// services/order-service/tests/unit/order.service.test.js
const OrderService = require('../../src/services/order.service');
const Order = require('../../src/models/order.model');
const KafkaProducer = require('../../src/events/producer');

jest.mock('../../src/models/order.model');
jest.mock('../../src/events/producer');

describe('OrderService - Unit Tests', () => {
  
  describe('createOrder()', () => {
    const validOrderPayload = {
      user_id: 'user-uuid-001',
      service_type: 'ride',
      lokasi_jemput: { lat: -6.1754, lng: 106.8272, alamat: 'Jl. Sudirman 1' },
      lokasi_tujuan:  { lat: -6.2012, lng: 106.8456, alamat: 'Jl. Gatot Subroto' },
      metode_pembayaran: 'qris'
    };
    
    test('✅ Order ride berhasil dibuat', async () => {
      Order.create.mockResolvedValue({ _id: 'mongo-id', order_id: 'ORD-001', ...validOrderPayload });
      KafkaProducer.publish.mockResolvedValue(true);
      
      const result = await OrderService.createOrder(validOrderPayload);
      
      expect(result.order_id).toMatch(/^ORD-/);
      expect(result.status).toBe('pending');
      expect(KafkaProducer.publish).toHaveBeenCalledWith('order.created', expect.any(Object));
    });
    
    test('❌ Gagal jika service_type tidak valid', async () => {
      await expect(
        OrderService.createOrder({ ...validOrderPayload, service_type: 'taxi' })
      ).rejects.toThrow('service_type tidak valid');
    });
    
    test('❌ Gagal jika user sedang suspend', async () => {
      // Simulasi user suspended
      await expect(
        OrderService.createOrder({ ...validOrderPayload, user_id: 'suspended-user' })
      ).rejects.toThrow('Akun Anda sedang di-suspend');
    });
  });
  
  describe('negotiatePrice()', () => {
    test('✅ Penawaran dalam batas wajar diterima', async () => {
      // Estimasi: 25000, tawaran: 22000 (masih > min price 15000)
      const result = await OrderService.negotiatePrice('ORD-001', 22000);
      expect(result.harga_negosiasi).toBe(22000);
      expect(result.negosiasi_status).toBe('waiting_driver');
    });
    
    test('❌ Penawaran di bawah harga minimum ditolak', async () => {
      await expect(
        OrderService.negotiatePrice('ORD-001', 5000)  // Di bawah minimum
      ).rejects.toThrow('Harga di bawah minimum yang diizinkan');
    });
  });
  
  describe('cancelOrder()', () => {
    test('✅ Cancel sebelum matched tidak kena penalti', async () => {
      Order.findOne.mockResolvedValue({ status: 'pending', user_id: 'user-001' });
      const result = await OrderService.cancelOrder('ORD-001', 'user-001', 'Berubah pikiran');
      expect(result.penalty_applied).toBe(false);
    });
    
    test('❌ Cancel setelah driver pickup kena tracking cancel', async () => {
      Order.findOne.mockResolvedValue({ status: 'pickup', user_id: 'user-001' });
      const result = await OrderService.cancelOrder('ORD-001', 'user-001', 'Alasan');
      expect(result.cancel_count_incremented).toBe(true);
    });
  });
});
```

### 7.3 Payment Service Tests

```javascript
// services/payment-service/tests/unit/payment.service.test.js
describe('PaymentService - Idempotency Tests', () => {
  
  test('✅ Request pembayaran dengan key sama tidak diproses dua kali', async () => {
    const idempotencyKey = 'idem-key-order-001';
    
    // Panggil pertama kali
    const result1 = await PaymentService.createTransaction({
      order_id: 'ORD-001',
      nominal: 22000,
      idempotency_key: idempotencyKey
    });
    
    // Panggil kedua dengan key sama
    const result2 = await PaymentService.createTransaction({
      order_id: 'ORD-001',
      nominal: 22000,
      idempotency_key: idempotencyKey
    });
    
    // Hasil harus identik, transaksi tidak dibuat dua kali
    expect(result1.transaction_id).toBe(result2.transaction_id);
    expect(Transaction.create).toHaveBeenCalledTimes(1);
  });
  
  test('✅ Bagi hasil driver 80% akurat', async () => {
    const nominal = 25000;
    const result = await PaymentService.calculateSplit(nominal);
    
    expect(result.driver_earnings).toBe(20000);     // 80%
    expect(result.platform_fee).toBe(5000);         // 20%
    expect(result.driver_earnings + result.platform_fee).toBe(nominal);
  });
});
```

### 7.4 Matching Service Tests (Go)

```go
// services/matching-service/tests/matching_test.go
package tests

import (
    "testing"
    "github.com/dealan/matching-service/internal/service"
    "github.com/stretchr/testify/assert"
)

func TestFindNearestDriver(t *testing.T) {
    t.Run("✅ Menemukan driver terdekat dalam radius 5km", func(t *testing.T) {
        mockDrivers := []service.DriverLocation{
            {DriverID: "driver-001", Lat: -6.1760, Lng: 106.8280, Rating: 4.8},
            {DriverID: "driver-002", Lat: -6.1800, Lng: 106.8300, Rating: 4.5},
        }

        result, err := service.FindBestDriver(
            mockDrivers,
            service.Point{Lat: -6.1754, Lng: 106.8272},
            "ride",
        )
        
        assert.NoError(t, err)
        assert.Equal(t, "driver-001", result.DriverID)  // driver-001 lebih dekat
    })
    
    t.Run("❌ Tidak ada driver tersedia", func(t *testing.T) {
        _, err := service.FindBestDriver(
            []service.DriverLocation{},
            service.Point{Lat: -6.1754, Lng: 106.8272},
            "ride",
        )
        
        assert.ErrorIs(t, err, service.ErrNoDriverAvailable)
    })
    
    t.Run("✅ Filter driver berdasarkan jenis kendaraan", func(t *testing.T) {
        mockDrivers := []service.DriverLocation{
            {DriverID: "motor-001", VehicleType: "motor", Lat: -6.1760, Lng: 106.8280},
            {DriverID: "mobil-001", VehicleType: "mobil", Lat: -6.1761, Lng: 106.8281},
        }
        
        result, _ := service.FindBestDriver(mockDrivers,
            service.Point{Lat: -6.1754, Lng: 106.8272}, "car")
        
        assert.Equal(t, "mobil-001", result.DriverID)  // Hanya mobil untuk layanan car
    })
}
```

### 7.5 Test Configuration (Jest)

```javascript
// jest.config.js (root)
module.exports = {
  projects: [
    {
      displayName: 'auth-service',
      testMatch: ['<rootDir>/services/auth-service/tests/**/*.test.js'],
      testEnvironment: 'node',
      coverageThreshold: {
        global: { branches: 80, functions: 80, lines: 80, statements: 80 }
      }
    },
    {
      displayName: 'order-service',
      testMatch: ['<rootDir>/services/order-service/tests/**/*.test.js'],
      testEnvironment: 'node',
    },
    // ... service lainnya
  ],
  collectCoverageFrom: [
    'services/*/src/**/*.js',
    '!services/*/src/app.js',
    '!services/*/src/config/**'
  ]
};
```

---

## 8. CI/CD Pipeline Jenkins

### 8.1 Root Jenkinsfile (Orchestrator)

```groovy
// Jenkinsfile (root) — Orchestrates all microservices
pipeline {
    agent any
    
    environment {
        DOCKER_REGISTRY = 'gcr.io/dealan-project'
        K8S_NAMESPACE   = 'dealan-production'
        SONAR_HOST      = 'http://sonarqube:9000'
    }
    
    stages {
        stage('🔍 Detect Changed Services') {
            steps {
                script {
                    // Hanya build service yang berubah (monorepo optimization)
                    def changedFiles = sh(
                        script: 'git diff --name-only HEAD~1 HEAD',
                        returnStdout: true
                    ).trim().split('\n')
                    
                    env.CHANGED_SERVICES = changedFiles
                        .findAll { it.startsWith('services/') }
                        .collect { it.split('/')[1] }
                        .unique()
                        .join(',')
                    
                    echo "Changed services: ${env.CHANGED_SERVICES}"
                }
            }
        }
        
        stage('🧪 Parallel Test All Changed Services') {
            parallel {
                stage('Auth Service') {
                    when { expression { env.CHANGED_SERVICES.contains('auth-service') } }
                    steps { build job: 'dealan/auth-service', propagate: true }
                }
                stage('Order Service') {
                    when { expression { env.CHANGED_SERVICES.contains('order-service') } }
                    steps { build job: 'dealan/order-service', propagate: true }
                }
                stage('Matching Service') {
                    when { expression { env.CHANGED_SERVICES.contains('matching-service') } }
                    steps { build job: 'dealan/matching-service', propagate: true }
                }
                stage('Payment Service') {
                    when { expression { env.CHANGED_SERVICES.contains('payment-service') } }
                    steps { build job: 'dealan/payment-service', propagate: true }
                }
                // ... semua service lainnya
            }
        }
        
        stage('🚀 Deploy to Production') {
            when {
                allOf {
                    branch 'main'
                    expression { currentBuild.result == null || currentBuild.result == 'SUCCESS' }
                }
            }
            steps {
                script {
                    env.CHANGED_SERVICES.split(',').each { service ->
                        sh """
                            kubectl set image deployment/${service} \
                            ${service}=${DOCKER_REGISTRY}/${service}:${GIT_COMMIT} \
                            -n ${K8S_NAMESPACE}
                            kubectl rollout status deployment/${service} -n ${K8S_NAMESPACE} --timeout=300s
                        """
                    }
                }
            }
        }
    }
    
    post {
        success {
            slackSend(color: 'good', message: "✅ Dealan Pipeline SUCCESS: ${env.JOB_NAME} #${env.BUILD_NUMBER}")
        }
        failure {
            slackSend(color: 'danger', message: "❌ Dealan Pipeline FAILED: ${env.JOB_NAME} #${env.BUILD_NUMBER}")
        }
    }
}
```

### 8.2 Per-Service Jenkinsfile

```groovy
// services/auth-service/Jenkinsfile
pipeline {
    agent { docker { image 'node:20-alpine' } }
    
    environment {
        SERVICE_NAME = 'auth-service'
        IMAGE_TAG    = "${env.GIT_COMMIT?.take(7) ?: 'latest'}"
    }
    
    stages {
        stage('📦 Install Dependencies') {
            steps {
                dir("services/${SERVICE_NAME}") {
                    sh 'npm ci --prefer-offline'
                }
            }
        }
        
        stage('🔍 Lint') {
            steps {
                dir("services/${SERVICE_NAME}") {
                    sh 'npm run lint'
                }
            }
        }
        
        stage('🧪 Unit Tests') {
            steps {
                dir("services/${SERVICE_NAME}") {
                    sh 'npm run test:unit -- --coverage --coverageReporters=lcov'
                }
            }
            post {
                always {
                    junit "services/${SERVICE_NAME}/coverage/junit.xml"
                    publishHTML([
                        reportDir: "services/${SERVICE_NAME}/coverage/lcov-report",
                        reportFiles: 'index.html',
                        reportName: "${SERVICE_NAME} Coverage"
                    ])
                }
            }
        }
        
        stage('🧪 Functional Tests') {
            steps {
                dir("services/${SERVICE_NAME}") {
                    withCredentials([
                        string(credentialsId: 'SUPABASE_TEST_URL', variable: 'SUPABASE_URL'),
                        string(credentialsId: 'SUPABASE_TEST_KEY', variable: 'SUPABASE_KEY')
                    ]) {
                        sh 'npm run test:functional'
                    }
                }
            }
        }
        
        stage('📊 SonarQube Analysis') {
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh """
                        sonar-scanner \
                        -Dsonar.projectKey=dealan-${SERVICE_NAME} \
                        -Dsonar.sources=services/${SERVICE_NAME}/src \
                        -Dsonar.javascript.lcov.reportPaths=services/${SERVICE_NAME}/coverage/lcov.info
                    """
                }
                timeout(time: 1, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        
        stage('🐳 Build & Push Docker') {
            when { branch 'main' }
            steps {
                withCredentials([file(credentialsId: 'GCR_KEY', variable: 'GCR_KEY_FILE')]) {
                    sh """
                        gcloud auth activate-service-account --key-file=${GCR_KEY_FILE}
                        docker build -t gcr.io/dealan-project/${SERVICE_NAME}:${IMAGE_TAG} \
                                     -t gcr.io/dealan-project/${SERVICE_NAME}:latest \
                                     services/${SERVICE_NAME}/
                        docker push gcr.io/dealan-project/${SERVICE_NAME}:${IMAGE_TAG}
                        docker push gcr.io/dealan-project/${SERVICE_NAME}:latest
                    """
                }
            }
        }
    }
}
```

### 8.3 Dockerfile per Service (Node.js)

```dockerfile
# services/auth-service/Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:20-alpine AS runner

# Security: non-root user
RUN addgroup -S dealan && adduser -S dealan -G dealan
USER dealan

WORKDIR /app
COPY --chown=dealan:dealan --from=builder /app/node_modules ./node_modules
COPY --chown=dealan:dealan src ./src
COPY --chown=dealan:dealan package.json ./

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget -qO- http://localhost:3001/auth/health || exit 1

EXPOSE 3001
CMD ["node", "src/app.js"]
```

---

## 9. Deployment ke Cloud

### 9.1 Environment Variables per Service

```bash
# services/auth-service/.env.example
NODE_ENV=development
PORT=3001

# Supabase
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_SERVICE_KEY=your-service-role-key

# JWT
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRES_IN=15m
REFRESH_TOKEN_EXPIRES_IN=30d

# Redis
REDIS_URL=redis://localhost:6379

# OTP
OTP_EXPIRY_MINUTES=5
OTP_PROVIDER=twilio   # twilio | vonage | local (dev)
TWILIO_ACCOUNT_SID=xxx
TWILIO_AUTH_TOKEN=xxx
TWILIO_PHONE=+1234567890

# Rate Limiting
RATE_LIMIT_LOGIN_PER_MINUTE=5
RATE_LIMIT_OTP_PER_HOUR=3
```

```bash
# services/order-service/.env.example
NODE_ENV=development
PORT=3004

# MongoDB
MONGODB_URI=mongodb://localhost:27017/dealan_orders

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=order-service
KAFKA_CONSUMER_GROUP=order-service-group

# Redis
REDIS_URL=redis://localhost:6379

# Internal Service URLs (via API Gateway atau langsung)
PRICING_SERVICE_URL=http://localhost:3006
MATCHING_SERVICE_URL=http://localhost:3005
USER_SERVICE_URL=http://localhost:3002
```

### 9.2 Monitoring & Observability

```yaml
# infrastructure/k8s/monitoring/prometheus-values.yaml
# Prometheus + Grafana Stack untuk monitoring

grafana:
  dashboards:
    dealan:
      - name: "Dealan Overview"
        panels:
          - title: "Request Rate per Service"
            type: graph
          - title: "P99 Latency"
            type: graph
          - title: "Error Rate"
            type: stat
          - title: "Active Orders"
            type: stat

# Alert Rules
alerting:
  rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
      for: 2m
      labels:
        severity: critical
      annotations:
        summary: "Error rate > 5% pada {{ $labels.service }}"
    
    - alert: HighLatency
      expr: histogram_quantile(0.99, http_request_duration_seconds_bucket) > 1.0
      for: 5m
      annotations:
        summary: "P99 latency > 1 detik pada {{ $labels.service }}"
    
    - alert: OrderServiceDown
      expr: up{job="order-service"} == 0
      for: 1m
      labels:
        severity: critical
```

---

## 10. AI Prompt Engineering Guide

> Panduan ini berisi template prompt yang telah dioptimalkan untuk digunakan bersama AI coding assistant (Claude, GPT-4, Copilot) dalam mengerjakan setiap komponen proyek Dealan.

---

### 📌 PRINSIP PROMPT YANG BAIK UNTUK PROYEK INI

Sebelum membuat prompt, selalu sertakan konteks berikut:

```
KONTEKS PROYEK:
- Nama proyek: Dealan (ride-hailing app)
- Arsitektur: Microservices cloud-native
- Runtime: Node.js/Express (auth, user, driver, order, payment, notification, rating, punishment, promo, shipment, chat), Go/Gin (matching, location), Python/FastAPI (pricing, map-route)
- Database: Supabase/PostgreSQL (auth, user, driver, payment, rating, punishment), MongoDB (order, pricing, shipment, notification), Redis (matching, promo cache), Firebase RTDB (location, chat)
- Message Broker: Apache Kafka (event-driven async)
- API Gateway: Kong
- Containerization: Docker + Kubernetes
- CI/CD: Jenkins
- Testing: Jest (Node.js), Testify (Go), Pytest (Python)
```

---

### 10.1 Prompt: Membuat Service dari Scratch

```
TUGAS: Implementasikan [NAMA SERVICE] untuk aplikasi Dealan.

KONTEKS PROYEK:
[Paste konteks di atas]

SPESIFIKASI SERVICE:
- Nama: [nama-service]
- Penanggung jawab: [nama]
- Runtime: Node.js + Express
- Database: [Supabase PostgreSQL / MongoDB / Redis]
- Port: [nomor port]

ENDPOINT YANG HARUS DIBUAT:
1. POST /[endpoint] — [deskripsi]
2. GET  /[endpoint]/:id — [deskripsi]
[sesuaikan dengan API contract di dokumentasi]

INPUT SCHEMA:
[copy dari bagian 4.x dokumentasi ini]

OUTPUT SCHEMA:
[copy dari bagian 4.x dokumentasi ini]

RELASI KE SERVICE LAIN:
- Menerima event Kafka dari: [service]
- Mempublish event Kafka ke topic: [topic name]
- Memanggil REST API: [service] untuk [tujuan]

REQUIREMENTS WAJIB:
1. Struktur folder: controllers/, services/, models/, routes/, middlewares/
2. Validasi input menggunakan Joi atau express-validator
3. Error handling terpusat dengan kode HTTP yang benar
4. Logging menggunakan Winston (format JSON untuk produksi)
5. Health check endpoint GET /[service]/health yang return { status, uptime, version }
6. Unit test menggunakan Jest dengan minimal 80% coverage
7. Dockerfile multi-stage build dengan non-root user
8. Environment variable di .env.example (JANGAN hardcode credentials)
9. Tambahkan JSDoc pada semua fungsi di layer service

JANGAN BUAT:
- Jangan buat frontend/UI
- Jangan gunakan ORM (gunakan query builder atau raw SQL untuk Supabase)
- Jangan commit credentials apapun
- Jangan gunakan callback, gunakan async/await

OUTPUT YANG DIHARAPKAN:
Berikan kode lengkap untuk semua file berikut dalam urutan ini:
1. src/app.js
2. src/routes/[nama].routes.js
3. src/controllers/[nama].controller.js
4. src/services/[nama].service.js
5. src/models/[nama].model.js (jika MongoDB)
6. tests/unit/[nama].service.test.js
7. Dockerfile
8. package.json
9. .env.example
```

---

### 10.2 Prompt: Membuat Unit Test Lengkap

```
TUGAS: Buat unit test lengkap untuk [NAMA SERVICE] di proyek Dealan.

KONTEKS:
[File service yang sudah ada — paste isi file service.js di sini]

FRAMEWORK:
- Jest + Supertest (untuk Node.js)
- Semua dependency eksternal WAJIB di-mock (database, Kafka, Redis, HTTP calls)

TEST YANG HARUS DIBUAT:
Untuk setiap fungsi publik di service, buat minimal 3 test case:
1. Happy path (✅ berhasil dengan input valid)
2. Error path (❌ gagal dengan input tidak valid)
3. Edge case (⚠️ kasus batas seperti empty array, null values, concurrent calls)

TEMPLATE PENAMAAN:
- describe('[nama fungsi]()')
- test('✅ [deskripsi singkat skenario berhasil]')
- test('❌ [deskripsi singkat skenario gagal]')
- test('⚠️ [deskripsi edge case]')

SKENARIO SPESIFIK YANG WAJIB DICAKUP:
[Sesuaikan dengan bisnis logic service, misal untuk Order Service:]
- createOrder: test dengan service_type tidak valid
- negotiatePrice: test dengan harga di bawah minimum
- cancelOrder: test perbedaan penalti sebelum/sesudah pickup
- Concurrent order creation tidak boleh menghasilkan order duplikat

COVERAGE TARGET: Minimal 80% untuk branches, functions, lines, statements

OUTPUT:
Kode test lengkap yang siap dijalankan dengan "npm test"
```

---

### 10.3 Prompt: Implementasi Kafka Event Producer/Consumer

```
TUGAS: Tambahkan Kafka integration ke [NAMA SERVICE] di proyek Dealan.

KONTEKS:
- Broker: Apache Kafka (localhost:9092 di development)
- Library: kafkajs
- Pattern: Event-driven (fire-and-forget untuk notifikasi, saga untuk transaksi)

SERVICE INI HARUS:

MEMPUBLISH event ke topic berikut:
[Sesuaikan dengan event yang relevan, contoh untuk Order Service:]
- Topic: "order.created" saat order baru dibuat
  Payload: { order_id, user_id, service_type, lokasi_jemput, lokasi_tujuan, created_at }
- Topic: "order.cancelled" saat order dibatalkan
  Payload: { order_id, user_id, reason, cancelled_at }

MENGKONSUMSI event dari topic berikut:
- Topic: "payment.completed" → update status order menjadi "paid"
  Consumer group: "order-service-group"
- Topic: "matching.driver_assigned" → update driver_id pada order
  Consumer group: "order-service-group"

REQUIREMENTS:
1. Producer harus mendukung retry otomatis (3x dengan exponential backoff)
2. Consumer harus idempotent (cek apakah event sudah diproses via Redis SET NX)
3. Dead Letter Queue untuk event yang gagal diproses setelah max retry
4. Graceful shutdown: flush pending messages sebelum process exit
5. Logging setiap event yang dipublish dan dikonsumsi

STRUKTUR FILE:
- src/events/producer.js
- src/events/consumers/[topic-name].consumer.js
- src/events/schemas/[event-name].schema.js (validasi payload dengan Joi)

Buat kode lengkap untuk semua file di atas beserta unit testnya.
```

---

### 10.4 Prompt: Database Migration & Schema

```
TUGAS: Buat database migration dan seed data untuk [NAMA SERVICE] di proyek Dealan.

DATABASE: [Supabase (PostgreSQL) / MongoDB]

SCHEMA YANG SUDAH DIRANCANG:
[Paste schema dari bagian 4.x di dokumentasi ini]

UNTUK SUPABASE (PostgreSQL):
1. Buat file migration SQL: migrations/[timestamp]_create_[service]_tables.sql
2. Aktifkan Row Level Security (RLS) untuk semua tabel sensitif
3. Buat policy RLS: user hanya bisa akses data miliknya sendiri
4. Buat index untuk kolom yang sering di-query (foreign keys, status, created_at)
5. Buat fungsi PostgreSQL untuk logika yang kompleks (misal: update rating aggregate)
6. Buat seed file: seeds/[service].seed.sql dengan 10 data sample realistis

UNTUK MONGODB:
1. Buat file: src/config/db.js dengan koneksi dan index setup
2. Buat Mongoose schema di: src/models/[nama].model.js
3. Tambahkan validation pada schema level
4. Buat index untuk semua field yang di-query
5. Buat file seed: seeds/[service].seed.js dengan 10 dokumen sample

NAMING CONVENTION:
- Table/Collection: snake_case, plural (orders, users, driver_status)
- Column/Field: snake_case (created_at, driver_id, service_type)
- Index: idx_[table]_[column] (idx_orders_user_id)

Pastikan data seed mencakup berbagai status dan edge case untuk testing.
```

---

### 10.5 Prompt: Implementasi Real-time Feature (Location/Chat)

```
TUGAS: Implementasikan [Location Tracking / Chat] real-time untuk Dealan.

ARSITEKTUR TARGET:
- Protocol: WebSocket (Socket.io untuk Node.js, gorilla/websocket untuk Go)
- State: Firebase Realtime Database untuk persistence
- Buffer: Redis untuk batching write ke Firebase
- Client update rate: setiap 3 detik (location) / real-time (chat)

UNTUK LOCATION SERVICE (Go + WebSocket):
1. WebSocket server yang menerima koneksi dari ribuan driver concurrent
2. Setiap koneksi mengirim {driver_id, lat, lng, bearing} setiap 3 detik
3. Server memvalidasi JWT token pada WebSocket handshake
4. Update Redis buffer: HSET loc:buffer:{driver_id} lat {lat} lng {lng} bearing {bearing}
5. Background goroutine yang flush Redis buffer ke Firebase setiap 1 detik
6. Endpoint REST GET /location/driver/{driver_id} untuk one-time query
7. Endpoint REST GET /location/nearby?lat=&lng=&radius=5&type=ride untuk query driver terdekat

UNTUK CHAT SERVICE (Node.js + Socket.io):
1. Namespace /chat dengan room per order_id
2. Events yang ditangani:
   - join_chat: user/driver masuk ke room order
   - send_message: kirim pesan teks
   - typing: indikator sedang mengetik
   - read_receipt: tanda pesan sudah dibaca
3. Setiap message disimpan ke Firebase RTDB: /chats/{order_id}/messages/{message_id}
4. Kirim push notification via Notification Service jika penerima offline

KEAMANAN:
- Validasi JWT pada setiap koneksi Socket
- Pastikan user/driver hanya bisa join room order miliknya sendiri
- Rate limit: maksimal 30 pesan per menit per koneksi

Buat implementasi lengkap dengan error handling dan graceful shutdown.
```

---

### 10.6 Prompt: API Gateway Kong Configuration

```
TUGAS: Buat konfigurasi lengkap Kong API Gateway untuk semua 15 microservice Dealan.

DAFTAR SERVICE DAN PORT:
[Sesuaikan dengan daftar di bagian 2.3]

UNTUK SETIAP SERVICE, KONFIGURASI:
1. Route: path prefix /[service-name] → upstream ke http://[service-name]:[port]
2. Plugin jwt: untuk semua endpoint kecuali /auth/register, /auth/login, /*/health
3. Plugin rate-limiting: per consumer, policy via Redis
   - Auth endpoints: 20 req/menit
   - Order/Payment: 60 req/menit
   - Location (WebSocket): dikecualikan
4. Plugin cors: origins [https://dealan.id, http://localhost:3000]
5. Plugin request-transformer: tambah header X-Request-ID
6. Plugin response-transformer: hapus header Server dari response

BUAT FILE:
1. api-gateway/kong.yml (declarative format 3.0, database-less)
2. api-gateway/Dockerfile
3. api-gateway/docker-compose.override.yml untuk local dev

JUGA BUAT:
- Instruksi cara test setiap route dengan curl
- Penjelasan konfigurasi rate limit untuk setiap service
```

---

### 10.7 Prompt: Fix Bug / Debugging

```
TUGAS: Debug dan fix masalah berikut di proyek Dealan.

SERVICE: [nama service]
BUG DESCRIPTION: [deskripsikan bug dengan jelas]

ERROR LOG:
[paste error log di sini]

KODE YANG BERMASALAH:
[paste kode yang relevan]

KONTEKS:
- Ini terjadi saat: [kondisi spesifik, misal: saat 1000 concurrent requests]
- Frekuensi: [selalu terjadi / kadang-kadang]
- Environment: [development / staging / production]

YANG SUDAH DICOBA:
[list solusi yang sudah dicoba]

YANG DIHARAPKAN:
1. Identifikasi root cause masalah
2. Berikan fix minimal yang tidak breaking existing tests
3. Tambahkan unit test yang mereproduksi bug ini (test harus fail sebelum fix, pass setelah fix)
4. Jelaskan mengapa bug terjadi dan bagaimana mencegahnya di masa depan
5. Cek apakah ada service lain yang mungkin memiliki masalah serupa
```

---

### 10.8 Prompt: Code Review

```
TUGAS: Lakukan code review pada kode berikut dari proyek Dealan.

KODE YANG DIREVIEW:
[paste kode di sini]

SERVICE: [nama service]
LAYER: [controller / service / model / middleware]

KRITERIA REVIEW:
1. KEAMANAN: Apakah ada SQL injection, XSS, atau kerentanan lain?
2. PERFORMA: Apakah ada N+1 query? Apakah index digunakan dengan benar?
3. KEANDALAN: Apakah ada race condition? Apakah semua error ditangani?
4. MAINTAINABILITY: Apakah kode mudah dibaca? Apakah ada duplikasi?
5. TESTING: Apakah kode mudah di-unit test? Apakah ada side effect tersembunyi?
6. BEST PRACTICES: Apakah sesuai standar Node.js/Go/Python proyek ini?

FORMAT OUTPUT:
Berikan feedback dalam format:
🚨 CRITICAL: [masalah yang harus difix sebelum merge]
⚠️  WARNING: [masalah yang sebaiknya difix]
💡 SUGGESTION: [saran untuk improvement]
✅ GOOD: [hal yang sudah baik]

Sertakan contoh kode perbaikan untuk setiap CRITICAL dan WARNING item.
```

---

## 11. Roadmap Penyelesaian Proyek

### 11.1 Status Saat Ini

| Komponen | Status |
|---|---|
| ✅ Dokumen Analisis & Perancangan | Selesai |
| ✅ Unit Tests (semua service) | Selesai |
| ✅ Functional Tests | Selesai |
| ✅ Jenkins CI/CD Pipeline | Selesai |
| 🔄 Database Schema Implementation | **Tahap Selanjutnya** |
| ⏳ Service Implementation | Belum |
| ⏳ API Gateway Configuration | Belum |
| ⏳ Kafka Integration | Belum |
| ⏳ Kubernetes Deployment | Belum |
| ⏳ Monitoring & Alerting | Belum |

### 11.2 Sprint Plan

```
SPRINT 1 (Minggu 1-2): FOUNDATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Shakila:
  ├── Setup Supabase project + apply Auth schema migration
  ├── Implementasi Auth Service (register, login, OTP, JWT)
  └── Deploy Auth Service ke local Docker

Zahid:
  ├── Setup MongoDB Atlas + apply Order schema
  ├── Setup Redis (Upstash atau local)
  └── Docker Compose lengkap untuk local dev

Team:
  └── Setup Kong API Gateway + basic routing

SPRINT 2 (Minggu 3-4): CORE SERVICES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Shakila:
  ├── User Service (CRUD profil)
  └── Driver Service (profil + status management)

Zahid:
  ├── Order Service (create, status lifecycle)
  ├── Kafka setup + order.created event
  └── Matching Service (Go) — cari driver terdekat

Najmi:
  ├── Pricing Service (Python) — kalkulasi harga
  ├── Pricing negosiasi feature
  └── Shipment Service — detail paket GoSend

SPRINT 3 (Minggu 5-6): REAL-TIME & PAYMENT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Zahid:
  ├── Location Service (Go + WebSocket + Firebase)
  └── Real-time driver tracking

Niha:
  ├── Map & Route Service (Python + Google Maps API)
  ├── ETA calculation
  └── Chat Service (Socket.io + Firebase)

Najmi:
  ├── Payment Service (Midtrans integration)
  ├── Idempotency implementation
  └── Driver wallet management

SPRINT 4 (Minggu 7-8): ECOSYSTEM & POLISH
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Putri:
  ├── Notification Service (Kafka consumer + FCM)
  ├── Rating Service (mutual rating + aggregation)
  ├── Punishment Service (rule engine)
  └── Promo Service (Redis quota management)

Team:
  ├── Kubernetes manifests semua service
  ├── HPA configuration
  └── Prometheus + Grafana monitoring

SPRINT 5 (Minggu 9-10): HARDENING & DELIVERY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Team:
  ├── End-to-end integration testing
  ├── Load testing (k6) — simulasi 1000 concurrent users
  ├── Security audit (OWASP checklist)
  ├── Final documentation
  └── Demo preparation
```

### 11.3 Checklist Per Service (Template)

Gunakan checklist ini untuk setiap service sebelum dianggap "done":

```
SERVICE: [nama-service]
PENANGGUNG JAWAB: [nama]

IMPLEMENTATION:
[ ] Folder structure sesuai standar proyek
[ ] Semua endpoint dari API Contract diimplementasikan
[ ] Validasi input lengkap (Joi/express-validator)
[ ] Error handling terpusat
[ ] Logging terstruktur (JSON format)
[ ] Health check endpoint /health dan /ready

TESTING:
[ ] Unit tests coverage > 80%
[ ] Functional/integration tests untuk semua happy path
[ ] Functional tests untuk error cases penting
[ ] Semua tests pass di CI (Jenkins)

DATABASE:
[ ] Schema/migration sudah di-apply
[ ] Index sesuai query pattern
[ ] Seed data tersedia untuk development

INTEGRATION:
[ ] Kafka events (publish/consume) bekerja
[ ] Komunikasi ke service lain via REST/gRPC bekerja
[ ] Environment variables di .env.example lengkap

DEPLOYMENT:
[ ] Dockerfile bisa build tanpa error
[ ] Docker image size < 200MB
[ ] Kubernetes manifest tersedia
[ ] HPA dikonfigurasi

DOCUMENTATION:
[ ] README.md di folder service
[ ] JSDoc/comment untuk semua fungsi publik
[ ] Postman collection updated
```

---

## 12. Panduan Kontribusi Tim

### 12.1 Git Workflow

```bash
# Branch naming convention
feature/[service-name]/[fitur-singkat]   # misal: feature/auth-service/otp-verification
fix/[service-name]/[bug-singkat]         # misal: fix/order-service/duplicate-order
chore/[deskripsi]                        # misal: chore/update-dependencies

# Commit message format (Conventional Commits)
feat(auth): tambah endpoint refresh token
fix(order): cegah duplicate order saat concurrent request
test(payment): tambah test idempotency
docs(matching): update API contract dokumentasi
chore(ci): update Jenkins pipeline untuk Go services

# Workflow
1. git checkout main && git pull
2. git checkout -b feature/[service]/[fitur]
3. Kerjakan implementasi
4. npm test (pastikan semua test pass)
5. git add . && git commit -m "feat([service]): deskripsi"
6. git push origin feature/[service]/[fitur]
7. Buat Pull Request ke main
8. Minimal 1 reviewer approval sebelum merge
```

### 12.2 Environment Setup

```bash
# 1. Clone repo
git clone https://github.com/dealan-project/dealan.git
cd dealan

# 2. Copy environment files
for service in services/*/; do
  cp "$service/.env.example" "$service/.env"
done

# 3. Isi credentials di masing-masing .env
# Minta ke anggota tim yang sudah punya akses ke:
# - Supabase (dari Shakila)
# - Firebase (dari Zahid)
# - Midtrans sandbox (dari Najmi)

# 4. Start semua infrastruktur lokal
docker-compose up -d redis mongodb kafka kafka-ui kong

# 5. Start service yang ingin dikerja (contoh auth-service)
cd services/auth-service
npm install
npm run dev

# 6. Jalankan semua test
npm test                          # Dari root — jalankan semua service
npm test -- --testPathPattern auth # Hanya auth-service
```

### 12.3 Kontak & Resource

| Resource | Link / Info |
|---|---|
| Supabase Dashboard | Minta URL ke Shakila |
| MongoDB Atlas | Minta connection string ke Zahid |
| Firebase Console | Minta config ke Zahid |
| Jenkins Dashboard | `http://jenkins.dealan.local:8080` |
| Kafka UI | `http://localhost:8090` |
| Kong Admin API | `http://localhost:8001` |
| SonarQube | `http://sonarqube.dealan.local:9000` |
| Dokumentasi API | `/docs/API_REFERENCE.md` |
| Postman Collection | `/docs/POSTMAN_COLLECTION.json` |

---

## Appendix: Referensi

### Teknologi Utama
- [Supabase Docs](https://supabase.com/docs) — PostgreSQL + Auth + Storage
- [kafkajs](https://kafka.js.org/) — Kafka client untuk Node.js
- [Kong Gateway](https://docs.konghq.com/) — API Gateway
- [Firebase Realtime Database](https://firebase.google.com/docs/database) — Real-time data
- [Midtrans Sandbox](https://simulator.sandbox.midtrans.com/) — Payment gateway testing

### Arsitektur & Patterns
- [Microservices Patterns (Chris Richardson)](https://microservices.io/patterns/)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html) — Untuk distributed transactions
- [Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html) — Reliable event publishing

---

*Dokumen ini dibuat oleh Kelompok 22 Tubes — Program Studi Ilmu Komputer, Universitas Pendidikan Indonesia.*  
*Terakhir diperbarui: Juni 2025*