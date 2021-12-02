# Development Environment
1. Install Golang di https://go.dev/dl/
2. Ikuti perintah di https://gocv.io/getting-started/windows/, atau ikuti yang dibawah untuk lebih jelasnya
3. Download Mingw64 di https://gocv.io/getting-started/windows/, download yang x86_64-posix-seh
4. Yang di download adalah file 7z, ekstract isinya ke C:/Program Files
5. Cari 'Edit The System Environment Variable', lalu update PATH dan tambahkan C:\Program Files\__x86_64-8.1.0-release-posix-seh-rt_v6-rev0__\mingw64\bin (Perhatikan mungkin versinya berbeda)
6. Download https://cmake.org/download/, lalu install, pastikan centang "Add to Path" pas proses instalasi
7. Jalankan `go get -u -d gocv.io/x/gocv` di terminal
8. Buka folder `C:\Users\__hariangr__\go\pkg\mod\gocv.io\x\gocv@__v0.29.0__`, perhatikan username dan versi gocv mungkin berbeda (Memang seharusnya beda malah)
9. Jalankan file `win_build_opencv` di folder tadi, tunggu sampai selesai, bisa makan waktu sejam
10. Tambahkan `C:\opencv\build\install\x64\mingw\bin` ke PATH (Sama kek langkah 5)


# Probable Error
Jika menemukan error misal `ld cannot find -l opencv_core452` atau yang mirip, penyebabnya adalah GOCV yang digunakan mengekspect opencv versi tertentu terinstall (4.5.2 misalnya untuk error diatas) tapi file tersebut tidak ditemukan

Misal kemarin, September 2021 opencv terbaru 4.5.2 dan GOCV yang dipakai versi v0.27.0, tapi 3 bulan kemudian di bulan Desember 2021, di Windows baru coba diinstalin baru gabisa karena OpenCV yang terinstall adalah 4.5.4 tapi GOCV yang di aplikasi ini masih v0.27.0 (Gak support versi opencv terbaru)

Untuk mengatasinya ganti file di go.mod, dan bump versi gocv ke versi terbaru yang support opencv yang terinstall

# Quick Note
CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc420 -lopencv_highgui420 -lopencv_imgcodecs420 -lopencv_objdetect420 -lopencv_features2d420 -lopencv_video420 -lopencv_dnn420 -lopencv_xfeatures2d420 -lopencv_plot420 -lopencv_tracking420 -lopencv_img_hash420 -lopencv_calib3d420


set CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc452 -lopencv_highgui452 -lopencv_imgcodecs452 -lopencv_objdetect452 -lopencv_features2d452 -lopencv_video452 -lopencv_dnn452 -lopencv_xfeatures2d452 -lopencv_plot452 -lopencv_tracking452 -lopencv_img_hash452 -lopencv_calib3d452

go build -ldflags="-extldflags=-static" .\cmd\brainfreeze\

go build -ldflags="-extldflags='-static -static-libgcc -static-libstdc++'" .\cmd\brainfreeze\



set CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc452 -lopencv_highgui452 -lopencv_imgcodecs452 -lopencv_objdetect452 -lopencv_features2d452 -lopencv_video452 -lopencv_dnn452 -lopencv_xfeatures2d452 -lopencv_plot452 -lopencv_tracking452 -lopencv_img_hash452 -lopencv_calib3d452