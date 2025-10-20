# FutsalBook - Aplikasi Booking Lapangan Futsal

**Platform booking lapangan futsal online yang dibangun dengan Pure Golang dan Clean Architecture**

## Tentang Project

**FutsalBook** adalah aplikasi web untuk booking lapangan futsal yang dibangun dengan pendekatan **Clean Architecture** menggunakan **Pure Golang**. Project ini dibuat sebagai referensi pembelajaran untuk implementasi Clean Architecture, SOLID Principles, dan best practices dalam Go development.

### Untuk Customer (Penyewa)

- ✅ **Registrasi & Login** - Sistem autentikasi aman dengan bcrypt
- ✅ **Browse Lapangan** - Lihat daftar lapangan futsal tersedia
- ✅ **Cek Ketersediaan** - Real-time availability check per jam
- ✅ **Booking Lapangan** - Pesan lapangan dengan auto-calculate harga
- ✅ **Riwayat Booking** - Lihat history booking lengkap
- ✅ **Pembatalan** - Cancel booking dengan business rule H-2 jam
- ✅ **Pembayaran** - Integrasi payment gateway (simulasi/real)

### Untuk Owner (Pemilik Lapangan)

- ✅ **Manajemen Lapangan** - CRUD lapangan futsal
- ✅ **Setup Jadwal** - Atur jam operasional per hari
- ✅ **Set Harga** - Tentukan harga per jam
- ✅ **Lihat Booking** - Monitor semua booking lapangan
- ✅ **Dashboard** - Overview pendapatan dan statistik

## 🛠️ Teknologi yang Digunakan

### Backend

| Teknologi | Version | Fungsi |
|-----------|---------|--------|
| **Go** | 1.21+ | Bahasa pemrograman utama |
| **PostgreSQL** | 14+ | Database relational |
| **httprouter** | 1.3.0 | HTTP routing |
| **bcrypt** | latest | Password hashing |
| **godotenv** | 1.5.1 | Environment variables |
| **pq** | 1.10.9 | PostgreSQL driver |

### Frontend

| Teknologi | Fungsi |
|-----------|--------|
| **HTML5** | Struktur halaman |
| **Tailwind CSS** | Styling modern & responsive |
| **Vanilla JavaScript** | Interaktivitas |
| **Fetch API** | HTTP client |

### Database

| Fitur | Implementasi |
|-------|-------------|
| **Schema Design** | Normalized (3NF) |
| **Constraints** | Foreign keys, Check, Unique |
| **Indexes** | Performance optimization |
| **Triggers** | Auto-update timestamps |