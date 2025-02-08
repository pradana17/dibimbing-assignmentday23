# Dibimbing Assignment Day 23

Proyek ini adalah REST API berbasis **Golang** menggunakan **Gin** sebagai framework HTTP dan **GORM** sebagai ORM untuk mengelola database.

## 📌 Teknologi yang Digunakan

- Golang
- Gin (Framework HTTP)
- GORM (ORM untuk database)
- MySQL (Database)
- godotenv (untuk mengelola variabel lingkungan)

## ⚙️ Setup dan Instalasi

1. **Clone Repository**
   ```sh
   git clone https://github.com/username/repository.git
   cd repository
   ```

2. **Buat dan Konfigurasi Database**  
   Pastikan MySQL sudah terinstal, lalu buat database baru dengan perintah:
   ```sql
   CREATE DATABASE nama_database;
   ```

3. **Buat File `.env`**  
   Buat file `.env` berdasarkan contoh berikut:
   ```
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   DB_HOST=localhost
   DB_PORT=3306
   ```

4. **Instal Dependensi**
   ```sh
   go mod tidy
   ```

5. **Jalankan Proyek**
   ```sh
   go run main.go
   ```

## 📌 Dokumentasi API

### 1. Produk
- **Menambah Produk**
  - **POST** `/produk`
  - **Body:**
    ```json
    {
      "nama": "Produk A",
      "kategori": "Elektronik",
      "harga": 100000
    }
    ```
  - **Response:**
    ```json
    {
      "id": 1,
      "nama": "Produk A",
      "kategori": "Elektronik",
      "harga": 100000
    }
    ```

- **Mendapatkan Semua Produk**
  - **GET** `/produk`
  - **Response:**
    ```json
    [
      {
        "id": 1,
        "nama": "Produk A",
        "kategori": "Elektronik",
        "harga": 100000
      }
    ]
    ```

- **Mendapatkan Produk Berdasarkan ID**
  - **GET** `/produk/id/:id`

- **Mendapatkan Produk Berdasarkan Kategori**
  - **GET** `/produk/kategori/:kategori`

- **Mengupdate Produk**
  - **PUT** `/produk/:id`

- **Menghapus Produk**
  - **DELETE** `/produk/:id`

---

### Dokumentasi Lebih Lanjut Bisa Dilihat di Bawah ini
## 🔗 Postman Collection
[Download Postman Collection](postman_collection.json)

