# Dokumentasi API Pendataan Mahasiswa

Dokumentasi lengkap untuk frontend dalam mengimplementasikan alur bisnis **Pendataan Mahasiswa** — dari pengecekan status, registrasi, verifikasi PIN, hingga pengisian dan pengiriman data pendataan.

---

## Daftar Isi

1. [Overview](#overview)
2. [Alur Bisnis (Business Flow)](#alur-bisnis-business-flow)
3. [Base URL](#base-url)
4. [Endpoints](#endpoints)
   - [1. Cek Status Pendataan](#1-cek-status-pendataan)
   - [2. Registrasi Mahasiswa](#2-registrasi-mahasiswa)
   - [3. Verifikasi PIN](#3-verifikasi-pin)
   - [4. Ambil Data Pendataan](#4-ambil-data-pendataan)
   - [5. Simpan Draft](#5-simpan-draft)
   - [6. Submit Final](#6-submit-final)
7. [Model Data (Blok Pendataan)](#model-data-blok-pendataan)
8. [State Management](#state-management)
9. [Error Handling](#error-handling)
10. [Contoh Implementasi Frontend](#contoh-implementasi-frontend)

---

## Overview

Sistem Pendataan Mahasiswa memungkinkan mahasiswa untuk:
1. Mengecek apakah sudah terdaftar berdasarkan email/NIK
2. Registrasi akun baru dengan email, NIK, dan PIN
3. Memasukkan PIN untuk mengakses form pendataan
4. Mengisi data pendataan dalam beberapa blok (I-VI)
5. Menyimpan draft sementara
6. Mengirim data final (setelah submit, data tidak bisa diubah lagi)

---

## Alur Bisnis (Business Flow)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ALUR PENDAFTARAN MAHASISWA                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   LANGKAH   │    │   LANGKAH   │    │   LANGKAH   │    │   LANGKAH   │  │
│  │      1      │───▶│      2      │───▶│      3      │───▶│      4      │  │
│  │ Cek Status  │    │  Registrasi │    │ Verify PIN  │    │ Isi Form    │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘  │
│        │                                       │                  │         │
│        │ sudah                                 │                  │         │
│        │ terdaftar                             │                  ▼         │
│        └───────────────────────────────────────┘           ┌─────────────┐  │
│                                                            │  LANGKAH 5  │  │
│                                                            │ Simpan Draft│  │
│                                                            └─────────────┘  │
│                                                                   │         │
│                                                                   ▼         │
│                                                            ┌─────────────┐  │
│                                                            │  LANGKAH 6  │  │
│                                                            │Submit Final │  │
│                                                            └─────────────┘  │
│                                                                   │         │
│                                                                   ▼         │
│                                                            ┌─────────────┐  │
│                                                            │   SELESAI   │  │
│                                                            │  Status:    │  │
│                                                            │ SUBMITTED   │  │
│                                                            └─────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Detail Alur:

| Langkah | Kondisi | Aksi | Hasil |
|---------|---------|------|-------|
| 1 | Mahasiswa belum tahu status | POST `/mahasiswa/cek-pendataan` | `registered: true/false` |
| 2a | `registered: false` | POST `/mahasiswa/register` | Akun berhasil dibuat |
| 2b | `registered: true` | Langsung ke langkah 3 | - |
| 3 | Mahasiswa ingin isi form | POST `/pendataan-mahasiswa/verify-pin` | `valid: true/false` |
| 4 | PIN valid | GET `/pendataan-mahasiswa?email=&nik=` | Data pendataan (jika ada) |
| 5 | Mahasiswa isi sebagian | POST `/pendataan-mahasiswa/draft` | Draft tersimpan |
| 6 | Mahasiswa selesai isi | POST `/pendataan-mahasiswa/submit` | Data final (tidak bisa diubah) |

---

## Base URL

```
Development: http://localhost:3000
Production:  https://your-domain.com
```

---

## Endpoints

### 1. Cek Status Pendataan

Mengecek apakah mahasiswa sudah terdaftar berdasarkan email atau NIK.

**Endpoint:** `POST /mahasiswa/cek-pendataan`

**Request Body:**
```json
{
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456"
}
```

**Validasi:**
- `email`: Wajib, format email valid
- `nik`: Wajib, tepat 16 digit

**Response Body (sudah terdaftar):**
```json
{
  "registered": true,
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "message": "Anda sudah terdaftar. Silakan lanjut ke penginputan pendataan."
}
```

**Response Body (belum terdaftar):**
```json
{
  "registered": false,
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "message": "Anda belum terdaftar. Silakan lakukan registrasi."
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 200 | Berhasil mengecek status |
| 400 | Format input tidak valid |

---

### 2. Registrasi Mahasiswa

Mendaftarkan mahasiswa baru dengan email, NIK, dan PIN.

**Endpoint:** `POST /mahasiswa/register`

**Request Body:**
```json
{
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "pin": "123456",
  "pinConfirm": "123456"
}
```

**Validasi:**
- `email`: Wajib, format email valid, harus unik
- `nik`: Wajib, tepat 16 digit, harus unik
- `pin`: Wajib, tepat 6 digit
- `pinConfirm`: Wajib, harus sama dengan `pin`

**Response Body (sukses):**
```json
{
  "success": true,
  "message": "Registrasi berhasil. Silakan lanjut ke penginputan pendataan.",
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456"
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 201 | Registrasi berhasil |
| 400 | Format input tidak valid / PIN tidak cocok |
| 409 | Email atau NIK sudah terdaftar |

---

### 3. Verifikasi PIN

Memverifikasi PIN mahasiswa sebelum mengakses form pendataan.

**Endpoint:** `POST /pendataan-mahasiswa/verify-pin`

**Request Body:**
```json
{
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "pin": "123456"
}
```

**Response Body (PIN valid):**
```json
{
  "valid": true,
  "message": "PIN valid"
}
```

**Response Body (PIN salah):**
```json
{
  "valid": false,
  "message": "PIN salah"
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 200 | PIN valid atau salah (cek field `valid`) |
| 400 | Format input tidak valid |
| 401 | Data mahasiswa tidak ditemukan |

---

### 4. Ambil Data Pendataan

Mengambil data pendataan mahasiswa yang sudah tersimpan (untuk melanjutkan draft).

**Endpoint:** `GET /pendataan-mahasiswa?email={email}&nik={nik}`

**Query Parameters:**
| Param | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| email | string | Ya* | Email mahasiswa |
| nik | string | Ya* | NIK mahasiswa |

*\*Salah satu harus diisi*

**Response Body (data ditemukan):**
```json
{
  "id": 1,
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "status": "DRAFT",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T14:20:00Z",
  "blokI": { ... },
  "blokII": { ... },
  "blokIII": [ ... ],
  "blokIV": { ... },
  "blokV": { ... },
  "blokVI": { ... }
}
```

**Response Body (data tidak ditemukan):**
```json
{
  "error": "Data tidak ditemukan"
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 200 | Data ditemukan |
| 400 | Email atau NIK tidak diisi |
| 404 | Data tidak ditemukan |

---

### 5. Simpan Draft

Menyimpan data pendataan sebagai draft. Draft bisa disimpan berkala-kala (auto-save).

**Endpoint:** `POST /pendataan-mahasiswa/draft?email={email}&nik={nik}`

**Query Parameters:**
| Param | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| email | string | Ya | Email mahasiswa |
| nik | string | Ya | NIK mahasiswa |

**Request Body:**
```json
{
  "blokI": {
    "namaKepalaKeluarga": "John Doe",
    "nikKepalaKeluarga": "1234567890123456",
    "nomorKK": "9876543210987654",
    "jumlahAnggotaKeluargaKK": 4,
    "jumlahAnggotaKeluargaPendataan": 3,
    "provinsi": "Sumatera Utara",
    "kabupatenKota": "Kota Medan",
    "kecamatan": "Medan Kota",
    "desaKelurahan": "Medan Baru",
    "klasifikasiDesaKota": "Perkotaan",
    "kodePos": "20111",
    "kodeSLS": "001",
    "namaSLS": "SLS Medan Baru",
    "alamatLengkap": "Jl. Merdeka No. 10",
    "namaJalan": "Jl. Merdeka",
    "nomorRumah": "10",
    "alamatSesuaiKK": 1
  },
  "blokII": { ... },
  "blokIII": [ ... ],
  "blokIV": { ... },
  "blokV": { ... },
  "blokVI": { ... }
}
```

**Response Body:**
```json
{
  "id": 1,
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "status": "DRAFT",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T14:20:00Z",
  "blokI": { ... },
  "blokII": { ... },
  "blokIII": [ ... ],
  "blokIV": { ... },
  "blokV": { ... },
  "blokVI": { ... }
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 200 | Draft berhasil disimpan |
| 400 | Email atau NIK tidak diisi / data tidak valid |
| 500 | Server error |

**Catatan:**
- Jika belum ada data pendataan, sistem akan membuat record baru dengan status `DRAFT`
- Jika sudah ada draft, sistem akan mengupdate data yang ada
- Field PIN tidak perlu dikirim — diambil dari data registrasi

---

### 6. Submit Final

Mengirim data pendataan final. **Setelah submit, data TIDAK bisa diubah lagi.**

**Endpoint:** `POST /pendataan-mahasiswa/submit?email={email}&nik={nik}`

**Query Parameters:**
| Param | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| email | string | Ya | Email mahasiswa |
| nik | string | Ya | NIK mahasiswa |

**Request Body:**
```json
{
  "pin": "123456",
  "blokI": { ... },
  "blokII": { ... },
  "blokIII": [ ... ],
  "blokIV": { ... },
  "blokV": { ... },
  "blokVI": { ... }
}
```

**Response Body (sukses):**
```json
{
  "id": 1,
  "email": "mahasiswa@example.com",
  "nik": "1234567890123456",
  "status": "SUBMITTED",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T16:00:00Z",
  "blokI": { ... },
  "blokII": { ... },
  "blokIII": [ ... ],
  "blokIV": { ... },
  "blokV": { ... },
  "blokVI": { ... }
}
```

**Response Codes:**
| Code | Keterangan |
|------|------------|
| 200 | Data berhasil di-submit |
| 400 | Email/NIK tidak diisi / data tidak valid |
| 401 | PIN salah / data tidak ditemukan |
| 500 | Server error |

---

## Model Data (Blok Pendataan)

### Blok I — Keterangan Umum Keluarga

| Field | Tipe | Keterangan | Contoh |
|-------|------|------------|--------|
| namaKepalaKeluarga | string | Nama kepala keluarga | "John Doe" |
| nikKepalaKeluarga | string | NIK kepala keluarga | "1234567890123456" |
| nomorKK | string | Nomor Kartu Keluarga | "9876543210987654" |
| jumlahAnggotaKeluargaKK | int | Jumlah anggota KK | 4 |
| jumlahAnggotaKeluargaPendataan | int | Jumlah anggota yang didata | 3 |
| provinsi | string | Nama provinsi | "Sumatera Utara" |
| kabupatenKota | string | Nama kabupaten/kota | "Kota Medan" |
| kecamatan | string | Nama kecamatan | "Medan Kota" |
| desaKelurahan | string | Nama desa/kelurahan | "Medan Baru" |
| klasifikasiDesaKota | string | Klasifikasi daerah | "Perkotaan" |
| kodePos | string | Kode pos | "20111" |
| kodeSLS | string | Kode SLS | "001" |
| namaSLS | string | Nama SLS | "SLS Medan Baru" |
| alamatLengkap | string | Alamat lengkap | "Jl. Merdeka No. 10" |
| namaJalan | string | Nama jalan | "Jl. Merdeka" |
| nomorRumah | string | Nomor rumah | "10" |
| alamatSesuaiKK | *int | Alamat sesuai KK (1: Ya, 2: Tidak) | 1 |

### Blok II — Keterangan Usaha/Perusahaan

Berisi 33 rincian tentang usaha/perusahaan. Field-field utama:

| Field | Tipe | Keterangan |
|-------|------|------------|
| namaUsaha | string | Nama usaha |
| namaKomersial | string | Nama komersial usaha |
| alamatUsaha | string | Alamat usaha |
| rt | string | RT |
| rw | string | RW |
| nomorTelepon | string | Nomor telepon |
| nomorHP | string | Nomor HP |
| emailUsaha | string | Email usaha |
| jenisKawasan | *int | Jenis kawasan (1-10) |
| jenisUsaha | *int | Jenis usaha (1-6) |
| punyaNIB | *int | Punya NIB (1: Ya, 2: Tidak) |
| nomorNIB | string | Nomor NIB |
| statusBadanUsaha | string | Status badan usaha |
| namaPengusaha | string | Nama pengusaha/penanggung jawab |
| jenisKelamin | *int | Jenis kelamin (1: L, 2: P) |
| umurPengusaha | *int | Umur pengusaha |
| nikPengusaha | string | NIK pengusaha |
| kegiatanUtama | string | Kegiatan utama usaha |
| tempatUsaha | *int | Tempat usaha (1-11) |
| kodeKBLI | string | Kode KBLI |
| kategoriLapanganUsaha | string | Kategori lapangan usaha |
| jaringanUsaha | *int | Jaringan usaha (1-6) |
| menggunakanInternet | *int | Menggunakan internet (1: Ya, 2: Tidak) |
| tahunBeroperasi | *int | Tahun beroperasi |
| totalPekerja | *int | Total pekerja |
| totalPekerjaDibayar | *int | Total pekerja dibayar |

**Field keuangan tahun 2025:**
- upahGaji2025, biayaProduksi2025, biayaPembelian2025
- biayaOperasional2025, biayaNonOperasional2025, totalPengeluaran2025
- nilaiProduksi2025, pendapatanLain2025, totalPendapatan2025
- persenOnline2025

**Field keuangan satu bulan terakhir:**
- upahGajiBulan, biayaProduksiBulan, biayaPembelianBulan
- biayaOperasionalBulan, biayaNonOperasionalBulan, totalPengeluaranBulan
- nilaiProduksiBulan, pendapatanLainBulan, totalPendapatanBulan
- persenOnlineBulan

### Blok III — Keterangan Sosial Ekonomi Anggota Keluarga

Array of objects. Setiap object = 1 anggota keluarga.

| Field | Tipe | Keterangan |
|-------|------|------------|
| nomorUrut | string | Nomor urut anggota |
| namaAnggota | string | Nama anggota |
| nikAnggota | string | NIK anggota |
| hubunganKeluarga | *int | Hubungan keluarga (1-9) |
| keberadaanAnggota | *int | Keberadaan (1-7) |
| alamatDomisili | *int | Alamat domisili (1-4) |
| provinsiDomisili | string | Provinsi domisili |
| kabupatenKotaDomisili | string | Kabupaten/kota domisili |
| negaraDomisili | string | Negara domisili (untuk WNA) |
| statusPerkawinan | *int | Status perkawinan (1-4) |
| jenisKelaminAnggota | *int | Jenis kelamin (1: L, 2: P) |
| tanggalLahir | string | Tanggal lahir |
| umurAnggota | *int | Umur |
| partisipasiSekolah | *int | Partisipasi sekolah (0-2) |
| ijazahTertinggi | *int | Ijazah tertinggi (0-6) |
| statusPekerjaan | *int | Status pekerjaan (1-6, 9) |
| pendapatanPekerjaan | *int | Pendapatan dari pekerjaan (1: Ya, 2: Tidak, 9: Tidak tahu) |
| upahGaji | *float64 | Upah/gaji |
| punyaRekening | *int | Punya rekening (1-4, 9) |
| disabilitasFisik | *int | Disabilitas fisik (1: Ya, 2: Tidak, 9: Tidak tahu) |
| hipertensi | *int | Penyakit hipertensi |
| rematik | *int | Penyakit rematik |
| asma | *int | Penyakit asma |
| ... (dan seterusnya) | | |

### Blok IV — Keterangan Sosial Ekonomi Keluarga

| Field | Tipe | Keterangan |
|-------|------|------------|
| jenisBangunan | *int | Jenis bangunan (1-5) |
| statusKepemilikan | *int | Status kepemilikan (1-5) |
| buktiKepemilikan | *int | Bukti kepemilikan (1-4) |
| luasLantai | *float64 | Luas lantai (m²) |
| bahanLantai | *int | Bahan lantai (1-9) |
| bahanDinding | *int | Bahan dinding (1-7) |
| bahanAtap | *int | Bahan atap (1-8) |
| fasilitasBAB | *int | Fasilitas BAB (1-6) |
| jenisKloset | *int | Jenis kloset (1-4) |
| sumberAirMinum | *int | Sumber air minum (1-11) |
| sumberPenerangan | *int | Sumber penerangan (1-4) |
| dayaListrik | *int | Daya listrik (1-5) |
| pengeluaranListrik | *float64 | Pengeluaran listrik/bulan |
| pengeluaranMakanMinggu | *float64 | Pengeluaran makanan/minggu |
| tabungGas3kg | *int | Jumlah tabung gas 3kg |
| lemariEs | *int | Punya lemari es (0/1) |
| ac | *int | Punya AC (0/1) |
| sepedaMotor | *int | Jumlah sepeda motor |
| mobil | *int | Jumlah mobil |

### Blok V — Catatan

| Field | Tipe | Keterangan |
|-------|------|------------|
| catatan | string | Catatan tambahan (text) |

### Blok VI — Keterangan Pemberi Jawaban

| Field | Tipe | Keterangan |
|-------|------|------------|
| namaPPL | string | Nama petugas PPL |
| nipNms | string | NIP/NMS petugas |
| nomorHPPPL | string | Nomor HP petugas |
| emailPPL | string | Email petugas |
| tanggalPelaksanaan | string | Tanggal pelaksanaan |
| tandaTangan | *bool | Tanda tangan (true/false) |

---

## State Management

### Status Pendataan

```
┌─────────────┐     ┌─────────────┐
│   DRAFT     │────▶│  SUBMITTED  │
│  (default)  │     │  (final)    │
└─────────────┘     └─────────────┘
       │                   │
       │ bisa edit          │ TIDAK bisa edit
       │ berulang kali      │ (locked)
       ▼                   ▼
   Save Draft          Selesai
```

| Status | Keterangan | Bisa Edit |
|--------|------------|-----------|
| `DRAFT` | Draft tersimpan, belum final | Ya |
| `SUBMITTED` | Data final sudah dikirim | **Tidak** |

### Alur State:

1. Mahasiswa registrasi → Belum ada data pendataan
2. Mahasiswa mulai isi form → Draft tersimpan otomatis (status: `DRAFT`)
3. Mahasiswa kembali lagi → Load draft yang ada
4. Mahasiswa klik "Submit" → Status berubah `SUBMITTED`
5. Setelah `SUBMITTED` → Form menjadi read-only

---

## Error Handling

Semua error response mengikuti format:

```json
{
  "error": "Deskripsi error"
}
```

### Error Codes:

| HTTP Code | Kemungkinan Penyebab |
|-----------|---------------------|
| 400 | Format input tidak valid (email, NIK, PIN) |
| 401 | PIN salah / autentikasi gagal |
| 404 | Data tidak ditemukan |
| 409 | Email/NIK sudah terdaftar (registrasi) |
| 500 | Server error |

### Validasi Input:

| Field | Rule |
|-------|------|
| email | Format email valid |
| nik | Harus 16 digit |
| pin | Harus 6 digit |
| pinConfirm | Harus sama dengan pin |

---

## Contoh Implementasi Frontend

### Contoh Fetch API (JavaScript/TypeScript):

```typescript
const API_BASE = 'http://localhost:3000';

// 1. Cek status pendataan
async function checkStatus(email: string, nik: string) {
  const res = await fetch(`${API_BASE}/mahasiswa/cek-pendataan`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, nik }),
  });
  return res.json(); // { registered: boolean, message: string }
}

// 2. Registrasi
async function register(email: string, nik: string, pin: string, pinConfirm: string) {
  const res = await fetch(`${API_BASE}/mahasiswa/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, nik, pin, pinConfirm }),
  });
  return res.json();
}

// 3. Verifikasi PIN
async function verifyPIN(email: string, nik: string, pin: string) {
  const res = await fetch(`${API_BASE}/pendataan-mahasiswa/verify-pin`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, nik, pin }),
  });
  return res.json(); // { valid: boolean, message: string }
}

// 4. Ambil data pendataan
async function getPendataan(email: string, nik: string) {
  const res = await fetch(
    `${API_BASE}/pendataan-mahasiswa?email=${email}&nik=${nik}`
  );
  return res.json();
}

// 5. Simpan draft
async function saveDraft(email: string, nik: string, data: any) {
  const res = await fetch(
    `${API_BASE}/pendataan-mahasiswa/draft?email=${email}&nik=${nik}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }
  );
  return res.json();
}

// 6. Submit final
async function submitFinal(email: string, nik: string, pin: string, data: any) {
  const res = await fetch(
    `${API_BASE}/pendataan-mahasiswa/submit?email=${email}&nik=${nik}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ pin, ...data }),
    }
  );
  return res.json();
}
```

### Contoh Alur Lengkap Frontend:

```typescript
async function alurPendataanMahasiswa() {
  const email = 'mahasiswa@example.com';
  const nik = '1234567890123456';
  const pin = '123456';

  // Step 1: Cek status
  const status = await checkStatus(email, nik);
  if (!status.registered) {
    // Step 2: Registrasi jika belum terdaftar
    await register(email, nik, pin, pin);
  }

  // Step 3: Verifikasi PIN
  const pinResult = await verifyPIN(email, nik, pin);
  if (!pinResult.valid) {
    alert('PIN salah!');
    return;
  }

  // Step 4: Ambil data existing (jika ada)
  const existingData = await getPendataan(email, nik);

  // Step 5: Tampilkan form (isi dari existingData jika ada)
  // ... user mengisi form ...

  // Step 6: Auto-save draft setiap perubahan
  await saveDraft(email, nik, formData);

  // Step 7: Submit final
  const result = await submitFinal(email, nik, pin, formData);
  if (result.status === 'SUBMITTED') {
    alert('Pendataan berhasil dikirim!');
  }
}
```

---

## Catatan Penting

1. **PIN disimpan sebagai hash** — Frontend TIDAK boleh menyimpan PIN di localStorage dalam bentuk plain text
2. **Auto-save draft** — Disarankan auto-save setiap 30 detik atau saat user berpindah halaman
3. **Blok III adalah array** — Bisa menambah/menghapus anggota keluarga secara dinamis
4. **Field bertipe `*int`/`*float64`** — Boleh null/kosong jika tidak diisi (gunakan `null`, bukan `0`)
5. **Setelah SUBMITTED** — Frontend harus menampilkan form dalam mode read-only
6. **Email dan NIK** — Dikirim sebagai query parameter untuk endpoint draft dan submit
7. **PIN** — Hanya dikirim saat verify-pin dan submit (untuk keamanan)
