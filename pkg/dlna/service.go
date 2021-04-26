package dlna

import (
	"net"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
)

type dmsConfig struct {
	Path                string
	IfName              string
	Http                string
	FriendlyName        string
	LogHeaders          bool
	StallEventSubscribe bool
	NotifyInterval      time.Duration
}

var dmsServer *Server
var dmsStarted bool

func initDMS() {
	var dmsConfig = &dmsConfig{
		Path:           "",
		IfName:         "",
		Http:           ":1338",
		FriendlyName:   "",
		LogHeaders:     false,
		NotifyInterval: 30 * time.Second,
	}

	dmsServer = &Server{
		txnManager: manager.GetInstance().TxnManager,
		Interfaces: func(ifName string) (ifs []net.Interface) {
			var err error
			if ifName == "" {
				ifs, err = net.Interfaces()
			} else {
				var if_ *net.Interface
				if_, err = net.InterfaceByName(ifName)
				if if_ != nil {
					ifs = append(ifs, *if_)
				}
			}
			if err != nil {
				logger.Fatal(err)
			}
			var tmp []net.Interface
			for _, if_ := range ifs {
				if if_.Flags&net.FlagUp == 0 || if_.MTU <= 0 {
					continue
				}
				tmp = append(tmp, if_)
			}
			ifs = tmp
			return
		}(dmsConfig.IfName),
		HTTPConn: func() net.Listener {
			conn, err := net.Listen("tcp", dmsConfig.Http)
			if err != nil {
				logger.Fatal(err)
			}
			return conn
		}(),
		FriendlyName:   dmsConfig.FriendlyName,
		RootObjectPath: filepath.Clean(dmsConfig.Path),
		LogHeaders:     dmsConfig.LogHeaders,
		// Icons: []Icon{
		// 	{
		// 		Width:    48,
		// 		Height:   48,
		// 		Depth:    8,
		// 		Mimetype: "image/png",
		// 		//ReadSeeker: readIcon(config.Config.Interfaces.DLNA.ServiceImage, 48),
		// 	},
		// 	{
		// 		Width:    128,
		// 		Height:   128,
		// 		Depth:    8,
		// 		Mimetype: "image/png",
		// 		//ReadSeeker: readIcon(config.Config.Interfaces.DLNA.ServiceImage, 128),
		// 	},
		// },
		StallEventSubscribe: dmsConfig.StallEventSubscribe,
		NotifyInterval:      dmsConfig.NotifyInterval,
	}
}

// func getIconReader(fn string) (io.Reader, error) {
// 	b, err := assets.ReadFile("dlna/" + fn + ".png")
// 	return bytes.NewReader(b), err
// }

// func readIcon(path string, size uint) *bytes.Reader {
// 	r, err := getIconReader(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	imageData, _, err := image.Decode(r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return resizeImage(imageData, size)
// }

// func resizeImage(imageData image.Image, size uint) *bytes.Reader {
// 	img := resize.Resize(size, size, imageData, resize.Lanczos3)
// 	var buff bytes.Buffer
// 	png.Encode(&buff, img)
// 	return bytes.NewReader(buff.Bytes())
// }

func Start() {
	initDMS()
	go func() {
		logger.Info("Starting DLNA")
		if err := dmsServer.Serve(); err != nil {
			logger.Fatal(err)
		}
	}()
	dmsStarted = true
}

func Stop() {
	logger.Info("Stopping DLNA")
	err := dmsServer.Close()
	if err != nil {
		logger.Fatal(err)
	}
	dmsStarted = false
}

func IsStarted() bool {
	return dmsStarted
}
