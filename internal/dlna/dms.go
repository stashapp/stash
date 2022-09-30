package dlna

// Derived from: https://github.com/anacrolix/dms
// Copyright (c) 2012, Matt Joiner <anacrolix@gmail.com>.
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the <organization> nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/anacrolix/dms/soap"
	"github.com/anacrolix/dms/ssdp"
	"github.com/anacrolix/dms/upnp"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

type SceneFinder interface {
	scene.Queryer
	scene.IDFinder
}

type StudioFinder interface {
	All(ctx context.Context) ([]*models.Studio, error)
}

type TagFinder interface {
	All(ctx context.Context) ([]*models.Tag, error)
}

type PerformerFinder interface {
	All(ctx context.Context) ([]*models.Performer, error)
}

type MovieFinder interface {
	All(ctx context.Context) ([]*models.Movie, error)
}

const (
	serverField                 = "Linux/3.4 DLNADOC/1.50 UPnP/1.0 DMS/1.0"
	rootDeviceType              = "urn:schemas-upnp-org:device:MediaServer:1"
	rootDeviceModelName         = "dms 1.0xb"
	resPath                     = "/res"
	iconPath                    = "/icon"
	rootDescPath                = "/rootDesc.xml"
	contentDirectoryEventSubURL = "/evt/ContentDirectory"
	serviceControlURL           = "/ctl"
	deviceIconPath              = "/deviceIcon"
)

func makeDeviceUuid(unique string) string {
	h := md5.New()
	if _, err := io.WriteString(h, unique); err != nil {
		panic("makeDeviceUuid write failed: " + err.Error())
	}
	buf := h.Sum(nil)
	return upnp.FormatUUID(buf)
}

// Groups the service definition with its XML description.
type service struct {
	upnp.Service
	SCPD string
}

// Exposed UPnP AV services.
var services = []*service{
	{
		Service: upnp.Service{
			ServiceType: "urn:schemas-upnp-org:service:ContentDirectory:1",
			ServiceId:   "urn:upnp-org:serviceId:ContentDirectory",
			EventSubURL: contentDirectoryEventSubURL,
			ControlURL:  serviceControlURL,
		},
		SCPD: contentDirectoryServiceDescription,
	},
	{
		Service: upnp.Service{
			ServiceType: "urn:schemas-upnp-org:service:ConnectionManager:1",
			ServiceId:   "urn:upnp-org:serviceId:ConnectionManager",
			ControlURL:  serviceControlURL,
		},
		SCPD: connectionManagerServiceDescription,
	},
	{
		Service: upnp.Service{
			ServiceType: "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1",
			ServiceId:   "urn:microsoft.com:serviceId:X_MS_MediaReceiverRegistrar",
			ControlURL:  serviceControlURL,
		},
		SCPD: xmsMediaReceiverServiceDescription,
	},
}

func init() {
	for _, s := range services {
		p := path.Join("/scpd", s.ServiceId)
		s.SCPDURL = p
	}
}

func devices() []string {
	return []string{
		"urn:schemas-upnp-org:device:MediaServer:1",
	}
}

func serviceTypes() (ret []string) {
	for _, s := range services {
		ret = append(ret, s.ServiceType)
	}
	return
}
func (me *Server) httpPort() int {
	return me.HTTPConn.Addr().(*net.TCPAddr).Port
}

func (me *Server) serveHTTP() error {
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if me.LogHeaders {
				logger.Debugf("%s %s", r.Method, r.RequestURI)
				for k, v := range r.Header {
					logger.Debugf("%s: %s", k, v)
				}
			}
			w.Header().Set("Ext", "")
			w.Header().Set("Server", serverField)
			me.httpServeMux.ServeHTTP(&mitmRespWriter{
				ResponseWriter: w,
				logHeader:      me.LogHeaders,
			}, r)
		}),
	}
	err := srv.Serve(me.HTTPConn)
	select {
	case <-me.closed:
		return nil
	default:
		return err
	}
}

// An interface with these flags should be valid for SSDP.
const ssdpInterfaceFlags = net.FlagUp | net.FlagMulticast

func (me *Server) doSSDP() {
	active := 0
	stopped := make(chan struct{})
	for _, if_ := range me.Interfaces {
		active++
		go func(if_ net.Interface) {
			defer func() {
				stopped <- struct{}{}
			}()
			me.ssdpInterface(if_)
		}(if_)
	}
	for active > 0 {
		<-stopped
		active--
	}
}

// Run SSDP server on an interface.
func (me *Server) ssdpInterface(if_ net.Interface) {
	s := ssdp.Server{
		Interface: if_,
		Devices:   devices(),
		Services:  serviceTypes(),
		Location: func(ip net.IP) string {
			return me.location(ip)
		},
		Server:         serverField,
		UUID:           me.rootDeviceUUID,
		NotifyInterval: me.NotifyInterval,
	}
	if err := s.Init(); err != nil {
		if if_.Flags&ssdpInterfaceFlags != ssdpInterfaceFlags {
			// Didn't expect it to work anyway.
			return
		}
		if strings.Contains(err.Error(), "listen") {
			// OSX has a lot of dud interfaces. Failure to create a socket on
			// the interface are what we're expecting if the interface is no
			// good.
			return
		}
		logger.Errorf("error creating ssdp server on %s: %s", if_.Name, err)
		return
	}
	defer s.Close()
	logger.Debugf("started SSDP on %s", if_.Name)
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		if err := s.Serve(); err != nil {
			logger.Errorf("%q: %q\n", if_.Name, err)
		}
	}()
	select {
	case <-me.closed:
		// Returning will close the server.
	case <-stopped:
	}
}

var (
	startTime time.Time
)

type Icon struct {
	Width, Height, Depth int
	Mimetype             string
	io.ReadSeeker
}

type Server struct {
	HTTPConn       net.Listener
	FriendlyName   string
	Interfaces     []net.Interface
	httpServeMux   *http.ServeMux
	RootObjectPath string
	rootDescXML    []byte
	rootDeviceUUID string
	closed         chan struct{}
	ssdpStopped    chan struct{}
	// The service SOAP handler keyed by service URN.
	services   map[string]UPnPService
	LogHeaders bool
	Icons      []Icon
	// Stall event subscription requests until they drop. A workaround for
	// some bad clients.
	StallEventSubscribe bool
	// Time interval between SSPD announces
	NotifyInterval time.Duration

	txnManager         txn.Manager
	repository         Repository
	sceneServer        sceneServer
	ipWhitelistManager *ipWhitelistManager
}

// UPnP SOAP service.
type UPnPService interface {
	Handle(action string, argsXML []byte, r *http.Request) (respArgs map[string]string, err error)
	Subscribe(callback []*url.URL, timeoutSeconds int) (sid string, actualTimeout int, err error)
	Unsubscribe(sid string) error
}

type Cache interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

func init() {
	startTime = time.Now()
}

func xmlMarshalOrPanic(value interface{}) []byte {
	ret, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("xmlMarshalOrPanic failed to marshal %v: %s", value, err))
	}
	return ret
}

// TODO: Document the use of this for debugging.
type mitmRespWriter struct {
	http.ResponseWriter
	loggedHeader bool
	logHeader    bool
}

func (me *mitmRespWriter) WriteHeader(code int) {
	me.doLogHeader(code)
	me.ResponseWriter.WriteHeader(code)
}

func (me *mitmRespWriter) doLogHeader(code int) {
	if !me.logHeader {
		return
	}
	logger.Debugf("Response: %d", code)
	for k, v := range me.Header() {
		logger.Debugf("%s: %s", k, v)
	}
	me.loggedHeader = true
}

func (me *mitmRespWriter) Write(b []byte) (int, error) {
	if !me.loggedHeader {
		me.doLogHeader(200)
	}
	return me.ResponseWriter.Write(b)
}

// Deprecated: the CloseNotifier interface predates Go's context package.
// New code should use Request.Context instead.
func (me *mitmRespWriter) CloseNotify() <-chan bool {
	return me.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Set the SCPD serve paths.
func init() {
	for _, s := range services {
		p := path.Join("/scpd", s.ServiceId)
		s.SCPDURL = p
	}
}

// Install handlers to serve SCPD for each UPnP service.
func handleSCPDs(mux *http.ServeMux) {
	for _, s := range services {
		mux.HandleFunc(s.SCPDURL, func(serviceDesc string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", `text/xml; charset="utf-8"`)
				http.ServeContent(w, r, ".xml", startTime, bytes.NewReader([]byte(serviceDesc)))
			}
		}(s.SCPD))
	}
}

// Marshal SOAP response arguments into a response XML snippet.
func marshalSOAPResponse(sa upnp.SoapAction, args map[string]string) []byte {
	soapArgs := make([]soap.Arg, 0, len(args))
	for argName, value := range args {
		soapArgs = append(soapArgs, soap.Arg{
			XMLName: xml.Name{Local: argName},
			Value:   value,
		})
	}
	return []byte(fmt.Sprintf(`<u:%[1]sResponse xmlns:u="%[2]s">%[3]s</u:%[1]sResponse>`, sa.Action, sa.ServiceURN.String(), xmlMarshalOrPanic(soapArgs)))
}

// Handle a SOAP request and return the response arguments or UPnP error.
func (me *Server) soapActionResponse(sa upnp.SoapAction, actionRequestXML []byte, r *http.Request) (map[string]string, error) {
	service, ok := me.services[sa.Type]
	if !ok {
		// TODO: What's the invalid service error?!
		return nil, upnp.Errorf(upnp.InvalidActionErrorCode, "Invalid service: %s", sa.Type)
	}

	logger.Tracef("%s::Handle %s - %s", sa.Type, sa.Action, actionRequestXML)
	ret, err := service.Handle(sa.Action, actionRequestXML, r)
	if err == nil {
		logger.Tracef("< %v", ret)
	}

	return ret, err
}

// Handle a service control HTTP request.
func (me *Server) serviceControlHandler(w http.ResponseWriter, r *http.Request) {
	clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)

	ip := net.ParseIP(clientIp).String()
	if !me.ipWhitelistManager.ipAllowed(ip) {
		// only log if we haven't seen it
		if !me.ipWhitelistManager.addRecent(ip) {
			logger.Infof("not allowed client %s", clientIp)
		}

		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	soapActionString := r.Header.Get("SOAPACTION")
	soapAction, err := upnp.ParseActionHTTPHeader(soapActionString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var env soap.Envelope
	if err := xml.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// AwoX/1.1 UPnP/1.0 DLNADOC/1.50
	w.Header().Set("Content-Type", `text/xml; charset="utf-8"`)
	w.Header().Set("Ext", "")
	w.Header().Set("Server", serverField)
	soapRespXML, code := func() ([]byte, int) {
		respArgs, err := me.soapActionResponse(soapAction, env.Body.Action, r)
		if err != nil {
			upnpErr := upnp.ConvertError(err)
			return xmlMarshalOrPanic(soap.NewFault("UPnPError", upnpErr)), 500
		}
		return marshalSOAPResponse(soapAction, respArgs), 200
	}()
	bodyStr := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" standalone="yes"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>%s</s:Body></s:Envelope>`, soapRespXML)
	w.WriteHeader(code)
	if _, err := w.Write([]byte(bodyStr)); err != nil {
		logger.Errorf(err.Error())
	}
}

func (me *Server) serveIcon(w http.ResponseWriter, r *http.Request) {
	sceneId := r.URL.Query().Get("scene")
	if sceneId == "" {
		return
	}

	var scene *models.Scene
	err := txn.WithTxn(r.Context(), me.txnManager, func(ctx context.Context) error {
		idInt, err := strconv.Atoi(sceneId)
		if err != nil {
			return nil
		}
		scene, _ = me.repository.SceneFinder.Find(ctx, idInt)
		return nil
	})
	if err != nil {
		logger.Warnf("failed to execute read transaction while trying to serve an icon: %v", err)
	}

	if scene == nil {
		return
	}

	me.sceneServer.ServeScreenshot(scene, w, r)
}

func (me *Server) contentDirectoryInitialEvent(ctx context.Context, urls []*url.URL, sid string) {
	body := xmlMarshalOrPanic(upnp.PropertySet{
		Properties: []upnp.Property{
			{
				Variable: upnp.Variable{
					XMLName: xml.Name{
						Local: "SystemUpdateID",
					},
					Value: "0",
				},
			},
			// upnp.Property{
			// 	Variable: upnp.Variable{
			// 		XMLName: xml.Name{
			// 			Local: "ContainerUpdateIDs",
			// 		},
			// 	},
			// },
			// upnp.Property{
			// 	Variable: upnp.Variable{
			// 		XMLName: xml.Name{
			// 			Local: "TransferIDs",
			// 		},
			// 	},
			// },
		},
		Space: "urn:schemas-upnp-org:event-1-0",
	})
	body = append([]byte(`<?xml version="1.0"?>`+"\n"), body...)
	for _, _url := range urls {
		bodyReader := bytes.NewReader(body)
		req, err := http.NewRequestWithContext(ctx, "NOTIFY", _url.String(), bodyReader)
		if err != nil {
			logger.Errorf("Could not create a request to notify %s: %s", _url.String(), err)
			continue
		}
		req.Header["CONTENT-TYPE"] = []string{`text/xml; charset="utf-8"`}
		req.Header["NT"] = []string{"upnp:event"}
		req.Header["NTS"] = []string{"upnp:propchange"}
		req.Header["SID"] = []string{sid}
		req.Header["SEQ"] = []string{"0"}
		// req.Header["TRANSFER-ENCODING"] = []string{"chunked"}
		// req.ContentLength = int64(bodyReader.Len())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Errorf("Could not notify %s: %s", _url.String(), err)
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		if len(b) > 0 {
			logger.Debug(string(b))
		}
		resp.Body.Close()
	}
}

func (me *Server) contentDirectoryEventSubHandler(w http.ResponseWriter, r *http.Request) {
	if me.StallEventSubscribe {
		// I have an LG TV that doesn't like my eventing implementation.
		// Returning unimplemented (501?) errors, results in repeat subscribe
		// attempts which hits some kind of error count limit on the TV
		// causing it to forcefully disconnect. It also won't work if the CDS
		// service doesn't include an EventSubURL. The best thing I can do is
		// cause every attempt to subscribe to timeout on the TV end, which
		// reduces the error rate enough that the TV continues to operate
		// without eventing.
		//
		// I've not found a reliable way to identify this TV, since it and
		// others don't seem to include any client-identifying headers on
		// SUBSCRIBE requests.
		//
		// TODO: Get eventing to work with the problematic TV.
		t := time.Now()
		<-r.Context().Done()
		logger.Debugf("stalled subscribe connection went away after %s", time.Since(t))
		return
	}
	// The following code is a work in progress. It partially implements
	// the spec on eventing but hasn't been completed as I have nothing to
	// test it with.
	service := me.services["ContentDirectory"]
	switch {
	case r.Method == "SUBSCRIBE" && r.Header.Get("SID") == "":
		urls := upnp.ParseCallbackURLs(r.Header.Get("CALLBACK"))
		var timeout int
		fmt.Sscanf(r.Header.Get("TIMEOUT"), "Second-%d", &timeout)
		sid, timeout, _ := service.Subscribe(urls, timeout)
		w.Header()["SID"] = []string{sid}
		w.Header()["TIMEOUT"] = []string{fmt.Sprintf("Second-%d", timeout)}
		// TODO: Shouldn't have to do this to get headers logged.
		w.WriteHeader(http.StatusOK)
		go func() {
			time.Sleep(100 * time.Millisecond)
			me.contentDirectoryInitialEvent(r.Context(), urls, sid)
		}()
	case r.Method == "SUBSCRIBE":
		http.Error(w, "meh", http.StatusPreconditionFailed)
	default:
		logger.Debugf("unhandled event method: %s", r.Method)
	}
}

func (me *Server) initMux(mux *http.ServeMux) {
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("content-type", "text/html")
		err := rootTmpl.Execute(resp, struct {
			Readonly bool
			Path     string
		}{
			true,
			me.RootObjectPath,
		})
		if err != nil {
			logger.Errorf(err.Error())
		}
	})
	mux.HandleFunc(contentDirectoryEventSubURL, me.contentDirectoryEventSubHandler)
	mux.HandleFunc(iconPath, me.serveIcon)
	mux.HandleFunc(resPath, func(w http.ResponseWriter, r *http.Request) {
		sceneId := r.URL.Query().Get("scene")
		var scene *models.Scene
		err := txn.WithTxn(r.Context(), me.txnManager, func(ctx context.Context) error {
			sceneIdInt, err := strconv.Atoi(sceneId)
			if err != nil {
				return nil
			}
			scene, _ = me.repository.SceneFinder.Find(ctx, sceneIdInt)
			return nil
		})
		if err != nil {
			logger.Warnf("failed to execute read transaction for scene id (%v): %v", sceneId, err)
		}

		if scene == nil {
			return
		}

		me.sceneServer.StreamSceneDirect(scene, w, r)
	})
	mux.HandleFunc(rootDescPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", `text/xml; charset="utf-8"`)
		w.Header().Set("content-length", fmt.Sprint(len(me.rootDescXML)))
		w.Header().Set("server", serverField)
		if k, err := w.Write(me.rootDescXML); err != nil {
			logger.Warnf("could not write rootDescXML (wrote %v bytes of %v): %v", k, len(me.rootDescXML), err)
		}
	})
	handleSCPDs(mux)
	mux.HandleFunc(serviceControlURL, me.serviceControlHandler)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	for i, di := range me.Icons {
		mux.HandleFunc(fmt.Sprintf("%s/%d", deviceIconPath, i), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", di.Mimetype)
			http.ServeContent(w, r, "", time.Time{}, di.ReadSeeker)
		})
	}
}

func (me *Server) initServices() {
	me.services = map[string]UPnPService{
		"ContentDirectory": &contentDirectoryService{
			Server: me,
		},
		"ConnectionManager": &connectionManagerService{
			Server: me,
		},
		"X_MS_MediaReceiverRegistrar": &mediaReceiverRegistrarService{
			Server: me,
		},
	}
}

func (me *Server) Serve() (err error) {
	me.initServices()
	me.closed = make(chan struct{})
	if me.HTTPConn == nil {
		me.HTTPConn, err = net.Listen("tcp", "")
		if err != nil {
			return
		}
	}
	if me.Interfaces == nil {
		ifs, err := net.Interfaces()
		if err != nil {
			logger.Errorf(err.Error())
		}
		var tmp []net.Interface
		for _, if_ := range ifs {
			if if_.Flags&net.FlagUp == 0 || if_.MTU <= 0 {
				continue
			}
			tmp = append(tmp, if_)
		}
		me.Interfaces = tmp
	}
	me.httpServeMux = http.NewServeMux()
	me.rootDeviceUUID = makeDeviceUuid(me.FriendlyName)
	me.rootDescXML, err = xml.MarshalIndent(
		upnp.DeviceDesc{
			SpecVersion: upnp.SpecVersion{Major: 1, Minor: 0},
			Device: upnp.Device{
				DeviceType:   rootDeviceType,
				FriendlyName: me.FriendlyName,
				Manufacturer: me.FriendlyName,
				ModelName:    rootDeviceModelName,
				UDN:          me.rootDeviceUUID,
				ServiceList: func() (ss []upnp.Service) {
					for _, s := range services {
						ss = append(ss, s.Service)
					}
					return
				}(),
				IconList: func() (ret []upnp.Icon) {
					for i, di := range me.Icons {
						ret = append(ret, upnp.Icon{
							Height:   di.Height,
							Width:    di.Width,
							Depth:    di.Depth,
							Mimetype: di.Mimetype,
							URL:      fmt.Sprintf("%s/%d", deviceIconPath, i),
						})
					}
					return
				}(),
			},
		},
		" ", "  ")
	if err != nil {
		return
	}
	me.rootDescXML = append([]byte(`<?xml version="1.0"?>`), me.rootDescXML...)
	logger.Debug("HTTP srv on", me.HTTPConn.Addr())
	me.initMux(me.httpServeMux)
	me.ssdpStopped = make(chan struct{})
	go func() {
		me.doSSDP()
		close(me.ssdpStopped)
	}()
	return me.serveHTTP()
}

func (me *Server) Close() (err error) {
	close(me.closed)
	err = me.HTTPConn.Close()
	<-me.ssdpStopped
	return
}

func didl_lite(chardata string) string {
	return `<DIDL-Lite` +
		` xmlns:dc="http://purl.org/dc/elements/1.1/"` +
		` xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/"` +
		` xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"` +
		` xmlns:dlna="urn:schemas-dlna-org:metadata-1-0/">` +
		chardata +
		`</DIDL-Lite>`
}

func (me *Server) location(ip net.IP) string {
	url := url.URL{
		Scheme: "http",
		Host: (&net.TCPAddr{
			IP:   ip,
			Port: me.httpPort(),
		}).String(),
		Path: rootDescPath,
	}
	return url.String()
}
