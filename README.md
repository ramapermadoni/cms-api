## CMS API Gateway
### **Project Overview**  
CMS API Gateway adalah solusi open-source yang mempermudah pengelolaan dan pengiriman konten website dengan satu titik akses API. Dibangun dengan Golang, Gin, dan GORM, serta menggunakan PostgreSQL sebagai database. GORM Migrate akan digunakan untuk manajemen migrasi database, memastikan skema database dapat dikelola dengan efisien selama proses pengembangan dan deployment. Untuk percobaan, proyek ini akan di-deploy di Railway.  

---

### **Objective**  
- Menyediakan API Gateway modular dan mudah diakses untuk mengelola konten dan media.  
- Mengimplementasikan autentikasi dan otorisasi berbasis JWT.  
- Memanfaatkan GORM Migrate untuk pengelolaan migrasi database.  
- Mendukung deployment otomatis melalui Railway sebagai platform cloud.  

---

### **Technology Stack**  
- **Backend**: Golang dengan Gin Framework  
- **ORM**: GORM  
- **Database**: PostgreSQL  
- **Migration Tool**: GORM Migrate  
- **Authentication**: JWT (JSON Web Tokens)  
- **Deployment Platform**: Railway  

---

### **Database Design**  

#### **Tabel 1: `users` (Pengguna)**  
- **Fields**:  
  - `id`: Primary key, auto-increment  
  - `username`: Username unik  
  - `password`: Kata sandi terenkripsi  
  - `email`: Email unik  
  - `role`: Peran pengguna (admin/editor)  
  - `created_at`: Waktu pembuatan pengguna  

#### **Tabel 2: `categories` (Kategori Konten)**  
- **Fields**:  
  - `id`: Primary key, auto-increment  
  - `name`: Nama kategori  
  - `description`: Deskripsi kategori  
  - `created_at`: Waktu pembuatan kategori  

#### **Tabel 3: `posts` (Konten/Artikel)**  
- **Fields**:  
  - `id`: Primary key, auto-increment  
  - `title`: Judul artikel  
  - `content`: Isi artikel  
  - `category_id`: Foreign key ke `categories`  
  - `author`: Foreign key ke `users`  
  - `status`: Status artikel (draft/published)  
  - `created_at`: Waktu pembuatan  
  - `updated_at`: Waktu pembaruan  

#### **Tabel 4: `media` (File/Gambar)**  
- **Fields**:  
  - `id`: Primary key, auto-increment  
  - `file_name`: Nama file  
  - `file_path`: Path penyimpanan file  
  - `post_id`: Foreign key ke `posts`  
  - `uploaded_by`: Foreign key ke `users`  
  - `uploaded_at`: Waktu upload  

---

### **Database Migration Strategy**  
- **GORM Migrate** akan digunakan untuk:  
  - Membuat dan mengelola migrasi database dengan versi yang terkontrol.  
  - Memastikan kompatibilitas schema database di setiap tahap pengembangan.  
  - Menghindari risiko kesalahan data selama perubahan struktur tabel.  

**Contoh Rancangan Migrasi**:  
1. **Initial Migration**: Membuat tabel `users`, `categories`, `posts`, dan `media`.  
2. **Future Migrations**: Menambahkan kolom baru, memperbarui constraint, atau membuat tabel tambahan berdasarkan kebutuhan.  

---

### **API Structure**  

1. **User Management**  
   - `POST /api/users`: Register pengguna baru  
   - `POST /api/login`: Login dan mendapatkan JWT  

2. **Content Management**  
   - `GET /api/posts`: Mendapatkan daftar artikel  
   - `POST /api/posts`: Membuat artikel baru  
   - `GET /api/posts/:id`: Mendapatkan detail artikel  
   - `POST /api/posts/:id`: Memperbarui artikel  

3. **Category Management**  
   - `GET /api/categories`: Mendapatkan daftar kategori  
   - `POST /api/categories`: Membuat kategori baru  

4. **Media Management**  
   - `POST /api/media`: Mengunggah file media  
   - `GET /api/media/:id`: Mendapatkan file media  

---

### **Security and Authentication**  
- **JWT Tokens**:  
  - Digunakan untuk autentikasi dan otorisasi pengguna.  
- **Role-based Access Control (RBAC)**:  
  - Admin dapat mengelola pengguna dan kategori.  
  - Editor dapat membuat dan memperbarui konten.  

---

### **Deployment Strategy**  
- **Railway Deployment**:  
  - Railway dipilih sebagai platform untuk memastikan deployment cepat dan mudah dengan dukungan otomatisasi.  
  - Database PostgreSQL akan dikelola di Railway untuk integrasi langsung dengan aplikasi.  
  - CI/CD pipeline diaktifkan melalui Railway untuk otomatisasi proses build dan deploy dari GitHub repository.  

---



## **Deskripsi Role dan Hak Akses pada Sistem CMS***
### 1. **Admin**  
- **Deskripsi:** Pengguna dengan hak akses tertinggi di sistem.  
- **Hak Akses:**  
  - **User Management:**  
    - Menambahkan, memperbarui, dan menghapus pengguna (user).
  - **Category Management:**  
    - Menambahkan, mengedit, dan menghapus kategori.
  - **Post Management:**  
    - Membuat, mengedit, mempublikasikan, atau menghapus post.
  - **Media Management:**  
    - Mengunggah dan menghapus file media.
  - **Role Assignment:**  
    - Mengubah peran pengguna lainnya (misalnya mengubah pengguna menjadi admin atau editor).  
  - **Audit:**  
    - Dapat melihat data log seperti siapa yang membuat, memperbarui, atau menghapus suatu item.

---

### 2. **Editor**  
- **Deskripsi:** Pengguna yang bertanggung jawab untuk mengelola konten (kategori, post, dan media).  
- **Hak Akses:**  
  - **Post Management:**  
    - Membuat, mengedit, dan mempublikasikan post (kecuali post milik admin, tergantung aturan).  
  - **Category Management:**  
    - Menambahkan, mengedit, dan menghapus kategori.  
  - **Media Management:**  
    - Mengunggah dan menghapus file media.
  - **Terbatas pada user:**  
    - Tidak bisa mengubah data user atau peran pengguna lainnya.  

---

### 3. **Author (Penulis)**  
- **Deskripsi:** Pengguna yang hanya dapat membuat dan mengelola post miliknya sendiri.  
- **Hak Akses:**  
  - **Post Management:**  
    - Membuat dan mengedit post yang dimiliki, tetapi tidak bisa mempublikasikan (kecuali dengan izin editor/admin).  
  - **Media Management:**  
    - Dapat mengunggah media untuk post miliknya.  
  - **Terbatas:**  
    - Tidak bisa mengedit atau menghapus kategori.  
    - Tidak memiliki akses untuk melihat, menambah, atau menghapus pengguna lain.

---

### **Ringkasan Hak Akses untuk Setiap Role**  

| **Aksi**                     | **Admin** | **Editor** | **Author**         |
|------------------------------|-----------|------------|---------------------|
| **User Management**          |           |            |                     |
| Menambah user                | ✔️        | ❌         | ❌                  |
| Mengubah role pengguna       | ✔️        | ❌         | ❌                  |
| **Category Management**      |           |            |                     |
| Menambah kategori             | ✔️        | ✔️        | ❌                  |
| Mengedit kategori             | ✔️        | ✔️        | ❌                  |
| Menghapus kategori            | ✔️        | ✔️        | ❌                  |
| **Post Management**          |           |            |                     |
| Membuat post                 | ✔️        | ✔️        | ✔️                  |
| Mengedit post                | ✔️        | ✔️        | ✔️ (milik sendiri)  |
| Menghapus post               | ✔️        | ✔️        | ❌                  |
| Mempublikasikan post         | ✔️        | ✔️        | ❌                  |
| **Media Management**         |           |            |                     |
| Mengunggah media             | ✔️        | ✔️        | ✔️                  |
| Menghapus media              | ✔️        | ✔️        | ❌                  |

### Penjelasan Tabel:
- **User Management**: Mengelola pengguna, termasuk menambah pengguna dan mengubah peran.
- **Category Management**: Mengelola kategori konten yang dapat digunakan dalam pos.
- **Post Management**: Mengelola konten pos, termasuk membuat, mengedit, dan menghapus pos.
- **Media Management**: Mengelola media yang digunakan dalam pos atau kategori
---

## 🚀 Get Started Today!

Terima kasih telah menjelajahi **CMS API Gateway**!  
Kami berharap proyek ini mempercepat pengembangan Anda dan membuat API CMS menjadi lebih mudah. Jangan ragu untuk berkontribusi, menyarankan fitur, atau melaporkan masalah. Mari kita bangun sesuatu yang luar biasa bersama-sama!

---

### 📄 Documentation & Resources

- 📖 [![Postman Documentation](https://img.shields.io/badge/Postman-Documentation-orange?logo=postman)](https://documenter.getpostman.com/view/15292179/2sAXxY2T1J)  
- 🌐 Hosted on [cms-api-production-6d2c.up.railway.app](https://cms-api-production-6d2c.up.railway.app)

---

### 🛠️ Contributing

Want to make this project better? We welcome contributions!  
1. Fork the repository  
2. Create your feature branch: `git checkout -b feature/AmazingFeature`  
3. Commit your changes: `git commit -m 'Add some AmazingFeature'`  
4. Push to the branch: `git push origin feature/AmazingFeature`  
5. Open a Pull Request

---

### 👥 Contact & Support

Need help or have questions? Reach out to the team:  
📧 **ramaprogramming@gmail.com**  
🐙 [GitHub Issues](https://github.com/ramapermadoni/cms-api/issues) for bug reports and feature requests.

---

<p align="center">
  Built with ❤️ by Rama Permadoni
</p>
