# 📘 IoT Product R&D Control Center - REST API Documentation (GET Produk & Katalog)

Dokumentasi resmi API untuk integrasi aplikasi pihak ketiga (ERP, Web Portal, Mobile App, Microservices, Python Scripts) guna mengambil data Master Produk, Pohon Struktur Proyek (RAB/BOM), dan Master Komponen.

---

## 🌐 1. Server & Konfigurasi Dasar

| Parameter | Development | Production (Contoh) |
| :--- | :--- | :--- |
| **Protocol** | `http://` | `https://` |
| **Host / IP** | `localhost:8080` | `api-iot.perusahaan.com` |
| **API Base Path** | `/api/v1` | `/api/v1` |
| **Format Data** | `application/json` | `application/json` |

---

## 🔐 2. Autentikasi (Static API Token & JWT Bearer Token)

Tersedia 2 metode autentikasi untuk mengakses API:

### Metode 1: Static API Token (REKOMENDASI - Langsung Tanpa Perlu Login)
Anda dapat men-generate token statis melalui menu **API Access Tokens** di dashboard web dengan masa berlaku yang fleksibel (**7 Hari, 30 Hari, 90 Hari, 1 Tahun, atau UNLIMITED / Never Expire**).

* **Header Wajib (Langsung Kirim Token):**
```http
Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE
Content-Type: application/json
```
*Dengan metode ini, aplikasi luar (ERP, Mobile, Python) **TIDAK PERLU** memanggil endpoint login `/auth/login` terlebih dahulu.*

---

### Metode 2: Dynamic JWT Login Flow (Untuk Sesi Interaktif)
Jika ingin menggunakan alur login pengguna reguler:
* **Method:** `POST`
* **URL:** `/api/v1/auth/login`
* **Request Body:**
```json
{
  "username": "admin",
  "password": "your_password"
}
```
* **Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "USR-001",
      "username": "admin",
      "fullName": "Administrator",
      "role": "Super Admin"
    }
  }
}
```

* **Header Request:**
```http
Authorization: Bearer <token_dari_login>
Content-Type: application/json
```

---

## 📦 3. Endpoint GET Master Produk

### A. List Semua Produk (Dengan Filter & Pagination)
* **Method:** `GET`
* **URL:** `/api/v1/products`
* **Query Parameters (Opsional):**
  * `page` *(int, default: 1)* : Nomor halaman data.
  * `limit` *(int, default: 20)* : Jumlah data per halaman (maks: 100).
  * `search` *(string)* : Pencarian teks (Nama Produk, Kode MPN, Deskripsi).
  * `product_type` *(string)* : Filter tipe produk:
    * `TRADING` : Produk jadi / instrumen trading komersial (contoh: Sensor Aquas, Sensor Photonic).
    * `PROJECT` : Proyek solusi / master RAB (contoh: RAB Maintenance SPARING).
    * `RND` : Produk riset & development internal.
  * `category` *(string)* : Filter kategori produk.
  * `status` *(string)* : Filter status (`Active`, `Draft`, `Obsolete`).

#### Contoh Request:
```http
GET /api/v1/products?product_type=TRADING&limit=50 HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOi...
```

#### Contoh Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "PRD-6585a3f1",
      "code": "PRJ-CQC08260429-LIP0-0040",
      "name": "Sensor Photonic (COD, TSS)",
      "productType": "TRADING",
      "category": "B. Pergantian Alat (Sensor)",
      "status": "Active",
      "projectLead": "Bapak Atta",
      "targetMarket": "Kawasan Industri Lippo Cikarang",
      "description": "Sensor Photonic Optik COD & TSS Unit",
      "tradingDetail": {
        "purchasePrice": 0,
        "sellingPrice": 157616701,
        "currency": "IDR",
        "purchaseCurrency": "IDR",
        "sellingCurrency": "IDR"
      },
      "createdAt": "2026-09-02T13:55:53Z",
      "updatedAt": "2026-09-02T14:15:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 14,
    "totalPages": 1
  }
}
```

---

### B. Detail Satu Produk By ID
* **Method:** `GET`
* **URL:** `/api/v1/products/:id`

#### Contoh Request:
```http
GET /api/v1/products/PRD-6585a3f1 HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOi...
```

#### Contoh Response:
```json
{
  "success": true,
  "data": {
    "id": "PRD-6585a3f1",
    "code": "PRJ-CQC08260429-LIP0-0040",
    "name": "Sensor Photonic (COD, TSS)",
    "productType": "TRADING",
    "category": "B. Pergantian Alat (Sensor)",
    "status": "Active",
    "projectLead": "Bapak Atta",
    "targetMarket": "Kawasan Industri Lippo Cikarang",
    "description": "Sensor Photonic Optik COD & TSS Unit",
    "longDescription": "Sensor Photonic Optik COD & TSS Complete Sparing Probe",
    "tradingDetail": {
      "id": "TRD-8891ac22",
      "productId": "PRD-6585a3f1",
      "purchasePrice": 0,
      "sellingPrice": 157616701,
      "currency": "IDR"
    }
  }
}
```

---

## 🌳 4. Endpoint Struktur Pohon Proyek (RAB / Hierarchy)

Mengambil struktur berjenjang (*tree hierarchy*) dari produk tipe `PROJECT` lengkap dengan folder seksi A, B, C, sub-group, item sensor, dan Architecture Notes.

* **Method:** `GET`
* **URL:** `/api/v1/products/:id/project-items?format=tree`

#### Contoh Request:
```http
GET /api/v1/products/PRD-2026-3AB86C/project-items?format=tree HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOi...
```

#### Contoh Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "PRJ-SEC-01",
      "name": "A. Tanpa Pergantian Alat (sensor)",
      "itemType": "SUB_ASSEMBLY",
      "costPrice": 0,
      "sellingPrice": 0,
      "currency": "IDR",
      "subItems": [
        {
          "id": "PRJ-ITM-01",
          "name": "1. Jasa Preventive Maintenance (Lingkup Pekerjaan)",
          "itemType": "SUB_ASSEMBLY",
          "quantity": 1,
          "unitName": "Kunjungan",
          "costPrice": 0,
          "sellingPrice": 5000000,
          "currency": "IDR",
          "notes": "a. Cleaning unit & area\nb. Cleaning sensor\nc. Pengecekkan & pengukuran power sistem\nd. Pengecekkan status software aplikasi\ne. Pengecekkan Hardware Datalogger\nf. Pengecekkan dan pembersihan bak sampling\ng. Pengecekkan status paket data"
        }
      ]
    },
    {
      "id": "PRJ-SEC-02",
      "name": "B. Pergantian Alat (Sensor)",
      "itemType": "SUB_ASSEMBLY",
      "costPrice": 0,
      "sellingPrice": 0,
      "currency": "IDR",
      "subItems": [
        {
          "id": "PRJ-GRP-01",
          "name": "1. Penggantian & repair sparepart yang rusak (Corrective Maintenance)",
          "itemType": "SUB_ASSEMBLY",
          "subItems": [
            {
              "id": "PRJ-ITM-03",
              "name": "Sensor Aquas (pH)",
              "itemType": "PRODUCT",
              "productCode": "PRJ-CQC08260429-LIP0-0039",
              "quantity": 1,
              "unitName": "Pcs",
              "costPrice": 0,
              "sellingPrice": 15625311,
              "currency": "IDR"
            },
            {
              "id": "PRJ-ITM-04",
              "name": "Sensor Photonic (COD, TSS)",
              "itemType": "PRODUCT",
              "productCode": "PRJ-CQC08260429-LIP0-0040",
              "quantity": 1,
              "unitName": "Pcs",
              "costPrice": 0,
              "sellingPrice": 157616701,
              "currency": "IDR"
            }
          ]
        }
      ]
    }
  ]
}
```

---

## 🧩 5. Endpoint GET Master Komponen Elektronik

* **Method:** `GET`
* **URL:** `/api/v1/components?page=1&limit=50&search=`

#### Contoh Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "CMP-001",
      "partNumber": "STM32H753BIT6",
      "internalCode": "IPN-MCU-001",
      "name": "Arm Cortex-M7 480MHz MCU",
      "category": "Microcontrollers & ICs",
      "estimatedUnitCost": 18.50,
      "currency": "USD",
      "status": "Approved"
    }
  ]
}
```

---

## 💻 6. Contoh Kode Pemanggilan dari Aplikasi Lain

### A. JavaScript / TypeScript (Node.js / React / Vue / Next.js)
```javascript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';
const JWT_TOKEN = 'YOUR_JWT_TOKEN_HERE';

async function fetchTradingProducts() {
  try {
    const response = await axios.get(`${API_BASE_URL}/products`, {
      headers: {
        Authorization: `Bearer ${JWT_TOKEN}`
      },
      params: {
        product_type: 'TRADING',
        limit: 100
      }
    });

    if (response.data.success) {
      console.log('Daftar Produk:', response.data.data);
    }
  } catch (error) {
    console.error('Gagal mengambil produk:', error.response?.data || error.message);
  }
}

fetchTradingProducts();
```

### B. Python (Requests Library)
```python
import requests

API_URL = "http://localhost:8080/api/v1/products"
TOKEN = "YOUR_JWT_TOKEN_HERE"

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json"
}

params = {
    "product_type": "TRADING",
    "limit": 50
}

response = requests.get(API_URL, headers=headers, params=params)
res_json = response.json()

if res_json.get("success"):
    for item in res_json["data"]:
        sell_price = item.get("tradingDetail", {}).get("sellingPrice", 0)
        print(f"[{item['code']}] {item['name']} - Rp {sell_price:,.0f}")
```

### C. PHP (cURL Native)
```php
<?php
$token = "YOUR_JWT_TOKEN_HERE";
$ch = curl_init("http://localhost:8080/api/v1/products?product_type=TRADING&limit=50");

curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    "Authorization: Bearer " . $token,
    "Content-Type: application/json"
]);

$response = curl_exec($ch);
curl_close($ch);

$data = json_decode($response, true);
print_r($data['data']);
?>
```

### D. cURL Command Line (Terminal / CMD)
```bash
curl -X GET "http://localhost:8080/api/v1/products?product_type=TRADING&limit=20" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
     -H "Content-Type: application/json"
```

---

## 🚦 7. Referensi Kode Status HTTP

| Status Code | Arti | Keterangan |
| :--- | :--- | :--- |
| **`200 OK`** | Success | Permintaan berhasil dan data dikembalikan dalam objek `data`. |
| **`400 Bad Request`** | Parameter Salah | Format parameter query atau request tidak valid. |
| **`401 Unauthorized`** | Token Tidak Valid / Expired | Token tidak disertakan atau sudah kedaluwarsa. |
| **`403 Forbidden`** | Akses Ditolak | Akun tidak memiliki hak akses (`permission`) `products.view`. |
| **`404 Not Found`** | Data Tidak Ditemukan | ID produk tidak ditemukan di database. |
| **`500 Internal Error`** | Server Error | Terjadi kegagalan koneksi database atau internal sistem. |
