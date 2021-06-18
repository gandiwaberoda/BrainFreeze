package diagnostic

// import (
// 	"log"
// 	"net/http"

// 	"github.com/hybridgroup/mjpeg"
// 	"gocv.io/x/gocv"
// 	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
// 	"harianugrah.com/brainfreeze/pkg/models/configuration"
// )

// type StreamOutDiagnostic struct {
// 	topCameraAcquisition *acquisition.TopCameraAcquisition
// 	stopTopChan          chan struct{}
// 	conf                 *configuration.FreezeConfig
// }

// var (
// 	streamTopProcessed *mjpeg.Stream
// )

// func CreateNewStreamOutDiagnostic(topCamera *acquisition.TopCameraAcquisition, conf *configuration.FreezeConfig) *StreamOutDiagnostic {
// 	return &StreamOutDiagnostic{
// 		topCameraAcquisition: topCamera,
// 		conf:                 conf,
// 	}
// }

// func (c *StreamOutDiagnostic) workerTopProcessed() {
// 	streamTopProcessed = mjpeg.NewStream()
// 	http.Handle("/top", streamTopProcessed)

// 	img := gocv.NewMat()
// 	defer img.Close()

// 	for {
// 		select {
// 		case <-c.stopTopChan:
// 			return
// 		default:
// 			img = c.topCameraAcquisition.Read()
// 			buf, _ := gocv.IMEncode(".jpg", img)
// 			streamTopProcessed.UpdateJPEG(buf)
// 		}
// 	}
// }

// func (c *StreamOutDiagnostic) StartTopCameraOutput() {
// 	go c.workerTopProcessed()
// }

// func (c *StreamOutDiagnostic) StopTopCameraOutput() {
// 	c.stopTopChan <- struct{}{}

// }

// func (c *StreamOutDiagnostic) Start() {
// 	log.Fatal(http.ListenAndServe(c.conf.Diagnostic.StreamHost, nil))
// }
