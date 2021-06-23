

# Quick Note
CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc420 -lopencv_highgui420 -lopencv_imgcodecs420 -lopencv_objdetect420 -lopencv_features2d420 -lopencv_video420 -lopencv_dnn420 -lopencv_xfeatures2d420 -lopencv_plot420 -lopencv_tracking420 -lopencv_img_hash420 -lopencv_calib3d420


set CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc452 -lopencv_highgui452 -lopencv_imgcodecs452 -lopencv_objdetect452 -lopencv_features2d452 -lopencv_video452 -lopencv_dnn452 -lopencv_xfeatures2d452 -lopencv_plot452 -lopencv_tracking452 -lopencv_img_hash452 -lopencv_calib3d452

go build -ldflags="-extldflags=-static" .\cmd\brainfreeze\

go build -ldflags="-extldflags='-static -static-libgcc -static-libstdc++'" .\cmd\brainfreeze\



set CGO_LDFLAGS=-LC:\opencv\build\install\x64\mingw\lib -lopencv_core452 -lopencv_face452 -lopencv_videoio452 -lopencv_imgproc452 -lopencv_highgui452 -lopencv_imgcodecs452 -lopencv_objdetect452 -lopencv_features2d452 -lopencv_video452 -lopencv_dnn452 -lopencv_xfeatures2d452 -lopencv_plot452 -lopencv_tracking452 -lopencv_img_hash452 -lopencv_calib3d452