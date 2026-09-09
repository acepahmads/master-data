# 📘 IoT Product R&D Control Center - REST API Documentation (Database & Master Customer)

Dokumentasi resmi API untuk integrasi sistem eksternal (CRM, ERP, Web Portal, Mobile App, Microservices, Python Scripts) guna mengelola dan mengambil data **Master Customer / Klien** (Client Accounts, Buyer Organizations, Procurement Contacts).

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

Semua endpoint `/api/v1/customers` dilindungi oleh autentikasi (`AuthMiddleware`). Tersedia 2 metode autentikasi:

### Metode 1: Static API Token (REKOMENDASI - Langsung Tanpa Perlu Login)
Anda dapat men-generate token statis melalui menu **API Access Tokens** di dashboard web dengan masa berlaku fleksibel (**7 Hari, 30 Hari, 90 Hari, 1 Tahun, atau UNLIMITED / Never Expire**).

* **Header Wajib:**
```http
Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE
Content-Type: application/json
```
*Dengan metode ini, skrip otomatisasi backend / ERP **TIDAK PERLU** memanggil endpoint login `/auth/login` terlebih dahulu.*

---

### Metode 2: Dynamic JWT Login Flow (Sesi Pengguna Interaktif)
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

## 🗄️ 3. Skema Database Customer (`customers` Table)

Berikut adalah struktur tabel database Customer pada backend (GORM / SQLite / PostgreSQL / MySQL):

| Nama Kolom | Tipe Data | Constraint / Index | Default | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| **`id`** | `VARCHAR(64)` | Primary Key | Otomatis generate kode jika kosong | ID unik customer (contoh: `CUST-2026-001`) |
| **`code`** | `VARCHAR(64)` | Unique Index, Not Null | Otomatis: `CUST-YYYY-XXX` | Nomor kode identitas customer unik |
| **`name`** | `VARCHAR(255)` | Index, Not Null | Fallback ke `company` jika kosong | Nama PIC / Kontak Customer / Perorangan |
| **`company`** | `VARCHAR(255)` | Index | `""` | Nama Perusahaan / Organisasi Pembeli |
| **`phone`** | `VARCHAR(64)` | Nullable | `""` | Nomor Telepon / WhatsApp PIC |
| **`email`** | `VARCHAR(255)` | Nullable | `""` | Alamat Email resmi PIC / Perusahaan |
| **`descr`** | `TEXT` | Nullable | `""` | Catatan, profil proyek, atau riwayat quotation |
| **`status`** | `VARCHAR(32)` | Nullable | `'Active'` | Status hubungan: `Active`, `Lead`, `Inactive` |
| **`source`** | `VARCHAR(64)` | Nullable | `'MANUAL'` | Asal data: `MANUAL` (input form) atau `IMPORT_XLS` (otomatis import excel quotation) |
| **`created_at`**| `TIMESTAMP` | Default Now | Waktu pembuatan data | Tanggal & waktu record dibuat |
| **`updated_at`**| `TIMESTAMP` | Default Now | Waktu perubahan data | Tanggal & waktu record terakhir diperbarui |
| **`deleted_at`**| `TIMESTAMP` | Index (GORM Soft Delete) | `NULL` | Soft delete marker (tidak dihapus fisik dari DB) |

---

## 📡 4. Ringkasan Endpoint API Customer

| Method | Endpoint Path | Deskripsi | Akses |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/customers` | Mendapatkan daftar customer dengan filter, pencarian, dan pagination | Protected |
| `GET` | `/api/v1/customers/:id` | Mengambil detail spesifik satu customer berdasarkan `id` atau `code` | Protected |
| `POST` | `/api/v1/customers` | Membuat record customer baru secara manual | Protected |
| `PUT` | `/api/v1/customers/:id` | Memperbarui data customer berdasarkan `id` atau `code` | Protected |
| `DELETE` | `/api/v1/customers/:id` | Menghapus record customer (soft delete) | Protected |

---

## 📖 5. Detail Endpoint REST API

### A. List Semua Customer (Dengan Filter & Pagination)
* **Method:** `GET`
* **URL:** `/api/v1/customers`
* **Query Parameters (Opsional):**
  * `page` *(integer, default: 1)* : Nomor halaman data.
  * `limit` *(integer, default: 20, max: 100)* : Jumlah data per halaman.
  * `search` *(string)* : Pencarian teks parsial (mencocokkan `name`, `code`, `company`, `phone`, `email`, atau `descr`).
  * `status` *(string)* : Filter berdasarkan status customer:
    * `Active` : Customer aktif / kontrak berjalan.
    * `Lead` : Prospek calon customer.
    * `Inactive` : Customer non-aktif / suspend.
    * `All` : Tampilkan semua status.
  * `source` *(string)* : Filter asal data:
    * `MANUAL` : Dibuat manual oleh tim via UI dashboard atau API.
    * `IMPORT_XLS` : Otomatis di-ingest dari parser dokumen Excel quotation.
    * `All` : Tampilkan semua asal data.
  * `sort` *(string, default: `created_at desc`)* : Urutan data (contoh: `name asc`, `company asc`, `created_at desc`).

#### Contoh Request:
```http
GET /api/v1/customers?page=1&limit=20&status=Active&search=Lippo HTTP/1.1
Host: localhost:8080
Authorization: Bearer iot_live_YOUR_TOKEN_HERE
Content-Type: application/json
```

#### Contoh Response (200 OK):
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": "CUST-2026-001",
      "code": "CUST-2026-001",
      "name": "Bapak Hendra Wijaya",
      "company": "PT Lippo Cikarang Tbk",
      "phone": "081298765432",
      "email": "hendra.w@lippo-cikarang.com",
      "descr": "Pengadaan IoT Sensor SPARING Kawasan Industri",
      "status": "Active",
      "source": "IMPORT_XLS",
      "createdAt": "2026-09-02T10:15:30Z",
      "updatedAt": "2026-09-02T14:20:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### B. Detail Satu Customer by ID / Code
Mencari customer berdasarkan `id` atau nomor unik `code`.

* **Method:** `GET`
* **URL:** `/api/v1/customers/:id`
* **URL Parameter:**
  * `:id` *(string, required)* : ID customer (contoh: `CUST-2026-001`) atau Code customer.

#### Contoh Request:
```http
GET /api/v1/customers/CUST-2026-001 HTTP/1.1
Host: localhost:8080
Authorization: Bearer iot_live_YOUR_TOKEN_HERE
```

#### Contoh Response Sukses (200 OK):
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": "CUST-2026-001",
    "code": "CUST-2026-001",
    "name": "Bapak Hendra Wijaya",
    "company": "PT Lippo Cikarang Tbk",
    "phone": "081298765432",
    "email": "hendra.w@lippo-cikarang.com",
    "descr": "Pengadaan IoT Sensor SPARING Kawasan Industri",
    "status": "Active",
    "source": "IMPORT_XLS",
    "createdAt": "2026-09-02T10:15:30Z",
    "updatedAt": "2026-09-02T14:20:00Z"
  }
}
```

#### Contoh Response Jika Tidak Ditemukan (404 Not Found):
```json
{
  "success": false,
  "message": "Customer not found"
}
```

---

### C. Tambah Customer Baru (Create)
Membuat entri customer baru. Jika `code` tidak diisi, sistem akan otomatis menghasilkan kode format `CUST-YYYY-XXX`. Jika `name` kosong, nilai `company` akan digunakan sebagai nama kontak.

* **Method:** `POST`
* **URL:** `/api/v1/customers`
* **Request Body (JSON):**

| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `name` | `string` | Ya (atau `company`) | Nama PIC customer |
| `company` | `string` | Ya (atau `name`) | Nama Perusahaan / Instansi |
| `code` | `string` | Tidak | Nomor kode kustom (otomatis di-generate jika kosong) |
| `phone` | `string` | Tidak | Nomor telepon / nomor WhatsApp aktif |
| `email` | `string` | Tidak | Alamat email resmi |
| `descr` | `string` | Tidak | Deskripsi tambahan / catatan procurement |
| `status` | `string` | Tidak | `'Active'`, `'Lead'`, atau `'Inactive'` (default: `'Active'`) |
| `source` | `string` | Tidak | Default: `'MANUAL'` |

#### Contoh Request:
```http
POST /api/v1/customers HTTP/1.1
Host: localhost:8080
Authorization: Bearer iot_live_YOUR_TOKEN_HERE
Content-Type: application/json

{
  "name": "Ibu Ratna Dewi",
  "company": "PT Cikarang Listrindo Tbk",
  "phone": "081122334455",
  "email": "ratna.dewi@cl.co.id",
  "descr": "PIC Tender Pengadaan Monitoring Kualitas Air & Udara",
  "status": "Active"
}
```

#### Contoh Response (200 OK):
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": "CUST-2026-015",
    "code": "CUST-2026-015",
    "name": "Ibu Ratna Dewi",
    "company": "PT Cikarang Listrindo Tbk",
    "phone": "081122334455",
    "email": "ratna.dewi@cl.co.id",
    "descr": "PIC Tender Pengadaan Monitoring Kualitas Air & Udara",
    "status": "Active",
    "source": "MANUAL",
    "createdAt": "2026-09-09T15:30:00Z",
    "updatedAt": "2026-09-09T15:30:00Z"
  }
}
```

---

### D. Update Data Customer (Update)
Memperbarui data profil, kontak, atau status customer.

* **Method:** `PUT`
* **URL:** `/api/v1/customers/:id`
* **URL Parameter:**
  * `:id` *(string, required)* : ID atau Code customer yang akan diubah.
* **Request Body (JSON):**

```json
{
  "name": "Ibu Ratna Dewi M.",
  "company": "PT Cikarang Listrindo Tbk",
  "phone": "081122334455",
  "email": "ratna.dm@cl.co.id",
  "descr": "PIC Tender Pengadaan Monitoring Kualitas Air (Kontrak Aktif Fase 2)",
  "status": "Active"
}
```

#### Contoh Response (200 OK):
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": "CUST-2026-015",
    "code": "CUST-2026-015",
    "name": "Ibu Ratna Dewi M.",
    "company": "PT Cikarang Listrindo Tbk",
    "phone": "081122334455",
    "email": "ratna.dm@cl.co.id",
    "descr": "PIC Tender Pengadaan Monitoring Kualitas Air (Kontrak Aktif Fase 2)",
    "status": "Active",
    "source": "MANUAL",
    "createdAt": "2026-09-09T15:30:00Z",
    "updatedAt": "2026-09-09T15:35:10Z"
  }
}
```

---

### E. Hapus Customer (Delete)
Menghapus record customer menggunakan mekanisme **Soft Delete** (GORM `DeletedAt`). Data tidak dihapus permanen dari basis data fisik, sehingga riwayat quotation & transaksi masa lalu tetap aman.

* **Method:** `DELETE`
* **URL:** `/api/v1/customers/:id`
* **URL Parameter:**
  * `:id` *(string, required)* : ID atau Code customer.

#### Contoh Request:
```http
DELETE /api/v1/customers/CUST-2026-015 HTTP/1.1
Host: localhost:8080
Authorization: Bearer iot_live_YOUR_TOKEN_HERE
```

#### Contoh Response (200 OK):
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "deleted": true
  }
}
```

---

## 🔄 6. Otomatisasi Ingestion dari Dokumen Excel (Quotation Parser)

Selain melalui REST API di atas, sistem backend memiliki integrasi internal `UpsertFromImport`:
1. Saat pengguna mengunggah berkas Excel penawaran harga (*Quotation* atau RAB Proyek), parser membaca metadata sel customer (`customerName`, `customerCompany`, `customerPhone`, `customerEmail`).
2. Jika record perusahaan / nama customer sudah ada, backend akan melakukan merge / melengkapi data nomor telepon dan email yang sebelumnya masih kosong.
3. Jika belum pernah terdaftar, sistem membuatkan record customer baru dengan flag `source: "IMPORT_XLS"` dan kode otomatis `CUST-YYYY-XXX`.

---

## 💻 7. Contoh Implementasi Kode (Client Integration)

### A. JavaScript / TypeScript (Node.js / Axios)
```javascript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';
const API_TOKEN = 'iot_live_YOUR_STATIC_TOKEN_HERE';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Authorization': `Bearer ${API_TOKEN}`,
    'Content-Type': 'application/json'
  }
});

// 1. Ambil Semua Customer Aktif
async function getActiveCustomers() {
  try {
    const res = await apiClient.get('/customers', {
      params: {
        status: 'Active',
        limit: 50,
        page: 1
      }
    });
    console.log('Total Customer:', res.data.meta.total);
    console.log('Data:', res.data.data);
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}

// 2. Tambah Customer Baru
async function createNewCustomer(payload) {
  try {
    const res = await apiClient.post('/customers', payload);
    console.log('Customer Berhasil Ditambahkan:', res.data.data);
  } catch (error) {
    console.error('Gagal membuat customer:', error.response?.data || error.message);
  }
}

// Eksekusi
getActiveCustomers();
```

---

### B. Python (Requests Library)
```python
import requests

BASE_URL = "http://localhost:8080/api/v1/customers"
TOKEN = "iot_live_YOUR_STATIC_TOKEN_HERE"

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json"
}

# 1. Fetch data customer dengan pencarian
params = {
    "search": "Cikarang",
    "status": "Active",
    "limit": 10
}

response = requests.get(BASE_URL, headers=headers, params=params)
if response.status_code == 200:
    res_data = response.json()
    print(f"Ditemukan {res_data['meta']['total']} customer:")
    for cust in res_data.get("data", []):
        print(f"- [{cust['code']}] {cust['name']} ({cust.get('company', '-')}) | Phone: {cust.get('phone', '-')}")
else:
    print(f"Error {response.status_code}: {response.text}")

# 2. Membuat customer baru
new_cust_payload = {
    "name": "Bapak Bambang",
    "company": "PT Sinergi Alam Makmur",
    "phone": "081399887766",
    "email": "bambang@sinergi.co.id",
    "descr": "Proyek Pengolahan Limbah WWTP",
    "status": "Active"
}

create_resp = requests.post(BASE_URL, headers=headers, json=new_cust_payload)
print("Create Result:", create_resp.json())
```

---

### C. PHP (cURL Native)
```php
<?php
$token = "iot_live_YOUR_STATIC_TOKEN_HERE";
$endpoint = "http://localhost:8080/api/v1/customers?limit=25&status=Active";

$ch = curl_init($endpoint);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    "Authorization: Bearer " . $token,
    "Content-Type: application/json"
]);

$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($httpCode === 200) {
    $result = json_decode($response, true);
    print_r($result['data']);
} else {
    echo "Request Error: HTTP " . $httpCode;
}
?>
```

---

### D. cURL CLI (Terminal / CMD / Bash Script)
```bash
# 1. GET Daftar Customer
curl -X GET "http://localhost:8080/api/v1/customers?page=1&limit=10&status=Active" \
     -H "Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE" \
     -H "Content-Type: application/json"

# 2. GET Detail Customer by Code / ID
curl -X GET "http://localhost:8080/api/v1/customers/CUST-2026-001" \
     -H "Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE"

# 3. POST Tambah Customer
curl -X POST "http://localhost:8080/api/v1/customers" \
     -H "Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Bapak Irwan",
       "company": "PT Jaya Teknik Mandiri",
       "phone": "08123456789",
       "email": "irwan@jayateknik.com",
       "status": "Active"
     }'

# 4. PUT Update Customer
curl -X PUT "http://localhost:8080/api/v1/customers/CUST-2026-001" \
     -H "Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE" \
     -H "Content-Type: application/json" \
     -d '{
       "phone": "081299998888",
       "descr": "Nomor kontak WhatsApp baru"
     }'

# 5. DELETE Customer (Soft Delete)
curl -X DELETE "http://localhost:8080/api/v1/customers/CUST-2026-001" \
     -H "Authorization: Bearer iot_live_YOUR_STATIC_TOKEN_HERE"
```

---

## 🚦 8. Referensi Kode Status HTTP & Format Respons

Sistem API menggunakan envelope standar `APIResponse`:
```json
{
  "success": true,
  "message": "Deskripsi status operasi",
  "data": { ... },
  "meta": { ... },
  "errors": null
}
```

| Status Code | Arti | Keterangan |
| :--- | :--- | :--- |
| **`200 OK`** | Success | Permintaan berhasil diproses dan mengembalikan data yang diminta. |
| **`201 Created`** | Created | Resource baru berhasil dibuat. |
| **`400 Bad Request`** | Validasi Gagal | Payload tidak lengkap / format JSON salah (contoh: tidak ada `name` dan `company`). |
| **`401 Unauthorized`** | Token Tidak Valid | Header `Authorization: Bearer <token>` tidak disertakan atau kedaluwarsa. |
| **`403 Forbidden`** | Akses Ditolak | Token tidak memiliki hak akses yang memadai. |
| **`404 Not Found`** | Tidak Ditemukan | ID atau Code customer tidak ditemukan di database. |
| **`500 Internal Error`** | Server Error | Terjadi kesalahan teknis internal atau kegagalan query database backend. |
