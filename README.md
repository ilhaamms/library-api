# LIBRARY API

ini adalah API RESTful sederhana untuk mengelola buku dan penulis dengan autentikasi dan otorisasi pengguna menggunakan token JWT. API ini mendukung operasi CRUD (Create, Read, Update, Delete) untuk kedua entitas yaitu authors dan books dan menjalankan keseluruhan proyek menggunakan Docker Compose

# AUTH
ini adalah auth untuk register dan login, setiap user yang mau akses enpoint pada authors dan book wajib sudah login, kemudian gunakan token JWT yang didapat pada response ketika berhasil login, dan apabila ingin mengakses endpoint pada author & books, maka pada bagian header masukan key dengan Authorization kemudian untuk valuenya adalah token yang didapan ketika berhasil login, untuk token hanya berlaku selama 1 jam saja

# Dokumentasi LIBRARY API

Dokumentasi API RESTful LIBRARY bisa dilihat : [disini](https://documenter.getpostman.com/view/26190643/2sAXxMgu6A)