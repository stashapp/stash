package dlna

import (
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
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

type sceneServer interface {
	StreamSceneDirect(scene *models.Scene, w http.ResponseWriter, r *http.Request)
	ServeScreenshot(scene *models.Scene, w http.ResponseWriter, r *http.Request)
}

type Service struct {
	txnManager  models.TransactionManager
	config      *config.Instance
	sceneServer sceneServer

	server  *Server
	running bool
	mutex   sync.Mutex

	startTimer *time.Timer
	stopTimer  *time.Timer
}

func (s *Service) init() {
	var dmsConfig = &dmsConfig{
		Path:           "",
		IfName:         "",
		Http:           ":1338",
		FriendlyName:   "",
		LogHeaders:     false,
		NotifyInterval: 30 * time.Second,
	}

	s.server = &Server{
		txnManager:  s.txnManager,
		sceneServer: s.sceneServer,
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

// NewService initialises and returns a new DLNA service.
func NewService(txnManager models.TransactionManager, cfg *config.Instance, sceneServer sceneServer) *Service {
	ret := &Service{
		txnManager:  txnManager,
		sceneServer: sceneServer,
		config:      cfg,
		mutex:       sync.Mutex{},
	}

	ret.init()
	return ret
}

// Start starts the DLNA service. If duration is provided, then the service
// is stopped after the duration has elapsed.
func (s *Service) Start(duration *time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		go func() {
			logger.Info("Starting DLNA")
			if err := s.server.Serve(); err != nil {
				logger.Fatal(err)
			}
		}()
		s.running = true

		if s.startTimer != nil {
			s.startTimer.Stop()
			s.startTimer = nil
		}
	}

	if duration != nil {
		// clear the existing stop timer
		if s.stopTimer != nil {
			s.stopTimer.Stop()
		}

		if s.stopTimer == nil {
			s.stopTimer = time.AfterFunc(*duration, func() {
				s.Stop(nil)
			})
		}
	}
}

// Stop stops the DLNA service. If duration is provided, then the service
// is started after the duration has elapsed.
func (s *Service) Stop(duration *time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		logger.Info("Stopping DLNA")
		err := s.server.Close()
		if err != nil {
			logger.Fatal(err)
		}
		s.running = false

		if s.stopTimer != nil {
			s.stopTimer.Stop()
			s.stopTimer = nil
		}
	}

	if duration != nil {
		// clear the existing stop timer
		if s.startTimer != nil {
			s.startTimer.Stop()
		}

		if s.startTimer == nil {
			s.startTimer = time.AfterFunc(*duration, func() {
				s.Start(nil)
			})
		}
	}
}

// IsRunning returns true if the DLNA service is running.
func (s *Service) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}
