# Breathe-io

## Overview
Layanan Breathe.io menyediakan dua tingkat layanan utama: **Free Tier** dan **Business Tier**. Setiap layanan menawarkan fitur-fitur khusus untuk membantu pengguna memantau kualitas udara dan mendapatkan rekomendasi terkait kesehatan serta dampak lingkungan. Pada Business Tier, pengguna bisnis dapat memantau emisi karbon perusahaan dan mendapatkan laporan yang lebih rinci tentang kualitas udara serta emisi di berbagai lokasi.

## gdrive with brainstorming
```
https://drive.google.com/drive/folders/1QPS_OQFlzARkl7NMHFsiZhKVE-XqwLk7?usp=drive_link
```

## Coretan
```
https://excalidraw.com/#json=RNy2gegtaJurt5RCrTcl0,GlqvPX5bDQa-FnIVzzjjAA
```

## Free Tier Process

### **User Input:**
- Lokasi yang diharapkan (dalam bentuk koordinat longitude dan latitude).

### **System Output:**
- Sistem akan mengambil data kualitas udara dari API pihak ketiga berdasarkan lokasi yang diinput.
- Sistem memberikan saran kesehatan, seperti:
  - Hindari aktivitas di luar ruangan jika kualitas udara buruk.
  - Gunakan masker jika berada di luar ruangan.
  - Buka jendela untuk ventilasi jika kualitas udara baik.

### **Notification:**
- Saran kesehatan akan dikirimkan ke pengguna melalui **WhatsApp** atau **email**.

## Business Tier Process (Dengan Proses Payment Gateway Menggunakan Xendit)

### **User Input:**
- **Informasi Bisnis**:
  - Tipe perusahaan (misalnya, pabrik, kantor, restoran).
  - Total emisi karbon yang dihasilkan oleh bisnis.
  - Lokasi yang diharapkan (dalam bentuk koordinat longitude dan latitude).
  
### **Multi-location Monitoring for Business:**
- Pengguna bisnis dapat memantau beberapa lokasi secara bersamaan, terutama jika mereka memiliki fasilitas atau pabrik di lokasi yang berbeda. Laporan dapat disusun untuk semua lokasi tersebut.

### **System Output:**
- **Perbandingan emisi**: Sistem akan membandingkan total emisi karbon yang dihasilkan oleh perusahaan dengan emisi lingkungan di lokasi yang diinput (data diambil dari API pihak ketiga).
- **Laporan rekomendasi lokasi**: Menentukan apakah lokasi tersebut sesuai untuk bisnis berdasarkan baseline pencemaran lingkungan.
- **Rekomendasi emisi**: Termasuk data emisi lingkungan dan emisi yang dihasilkan oleh bisnis tersebut.

### **Notification & Reporting:**
Laporan rinci dikirim melalui **WhatsApp** atau **email** menggunakan **Twilio**, mencakup:
- Link unduh **PDF/CSV** yang berisi:
  - Data historis kualitas udara di lokasi tersebut.
  - Laporan rekomendasi fasilitas (jika ada lebih dari satu rekomendasi).
  - Laporan pajak karbon, menghitung total emisi fasilitas dikalikan dengan harga per karbon.
  - Laporan emisi gas rumah kaca yang dihasilkan.

### **Scheduled Notifications:**
- Notifikasi berkala setiap 1 menit (atau sesuai interval yang diatur) akan mengirimkan laporan ke WhatsApp atau email pengguna.

## Pengembangan API Tambahan

### Fitur tambahan:
- **Grafik interaktif** atau **peta visualisasi** yang menampilkan data historis kualitas udara di berbagai lokasi selama periode tertentu.
- Pengguna dapat melihat **tren kualitas udara** dari waktu ke waktu untuk lokasi yang mereka minati.

### **API Input:**
- Lokasi (dalam bentuk koordinat longitude dan latitude).
- Start date dan end date (periode waktu untuk tren kualitas udara).

### **API Output:**
- Responsenya dalam format **JSON**, dengan informasi berupa link untuk mengunduh:
  - PDF/CSV berisi laporan data historis kualitas udara berdasarkan periode waktu yang diminta.

## Dependencies

- **Third-party Air Quality API**: Untuk mengambil data kualitas udara di lokasi tertentu.
- **Twilio**: Untuk mengirimkan notifikasi melalui WhatsApp atau email.
- **Xendit**: Sebagai payment gateway untuk menangani pembayaran dari pengguna Business Tier.