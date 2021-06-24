xconfig:
	@echo "Building e(x)periment for loading config executeable"
	@env GOOS=darwin GOARCH=amd64 go build -o ./output/xconfig/mac ./experiments/config/load_config.go
	@env GOOS=windows GOARCH=amd64 go build -o ./output/xconfig/win ./experiments/config/load_config.go
	@env GOOS=linux GOARCH=amd64 go build -o ./output/xconfig/lin ./experiments/config/load_config.go
	@echo "Finish building"

xcamprev:
	@echo "Building e(x)periment for camera preview"
	@env GOOS=darwin GOARCH=amd64 go build -o ./output/xcamprev/mac ./experiments/camera/camera_preview/camera_preview.go
	@env GOOS=windows GOARCH=amd64 go build -o ./output/xcamprev/win ./experiments/camera/camera_preview/camera_preview.go
	@env GOOS=linux GOARCH=amd64 go build -o ./output/xcamprev/lin ./experiments/camera/camera_preview/camera_preview.go
	@echo "Finish building"

run:
	go run ./cmd/brainfreeze/main.go

roi:
	go run cmd/cameraroi/main.go

exe:
	go build -o ./output ./cmd/brainfreeze 
	go build -o ./output ./cmd/cameraroi 
	go build -o ./output ./cmd/camerasrc
	