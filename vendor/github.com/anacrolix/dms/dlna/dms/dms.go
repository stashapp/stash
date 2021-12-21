package dms

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/anacrolix/dms/dlna"
	"github.com/anacrolix/dms/soap"
	"github.com/anacrolix/dms/ssdp"
	"github.com/anacrolix/dms/transcode"
	"github.com/anacrolix/dms/upnp"
	"github.com/anacrolix/dms/upnpav"
	"github.com/anacrolix/ffprobe"
)

const (
	serverField                 = "Linux/3.4 DLNADOC/1.50 UPnP/1.0 DMS/1.0"
	rootDeviceType              = "urn:schemas-upnp-org:device:MediaServer:1"
	rootDeviceModelName         = "dms 1.0"
	resPath                     = "/res"
	iconPath                    = "/icon"
	rootDescPath                = "/rootDesc.xml"
	contentDirectoryEventSubURL = "/evt/ContentDirectory"
	serviceControlURL           = "/ctl"
	deviceIconPath              = "/deviceIcon"
)

type transcodeSpec struct {
	mimeType        string
	DLNAProfileName string
	Transcode       func(path string, start, length time.Duration, stderr io.Writer) (r io.ReadCloser, err error)
}

var transcodes = map[string]transcodeSpec{
	"t": {
		mimeType:        "video/mpeg",
		DLNAProfileName: "MPEG_PS_PAL",
		Transcode:       transcode.Transcode,
	},
	"vp8":        {mimeType: "video/webm", Transcode: transcode.VP8Transcode},
	"chromecast": {mimeType: "video/mp4", Transcode: transcode.ChromecastTranscode},
	"web":        {mimeType: "video/mp4", Transcode: transcode.WebTranscode},
}

func makeDeviceUuid(unique string) string {
	h := md5.New()
	if _, err := io.WriteString(h, unique); err != nil {
		log.Panicf("makeDeviceUuid write failed: %s", err)
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
		},
		SCPD: contentDirectoryServiceDescription,
	},
	{
		Service: upnp.Service{
			ServiceType: "urn:schemas-upnp-org:service:ConnectionManager:1",
			ServiceId:   "urn:upnp-org:serviceId:ConnectionManager",
		},
		SCPD: connectionManagerServiceDescription,
	},
	{
		Service: upnp.Service{
			ServiceType: "urn:microsoft.com:service:X_MS_MediaReceiverRegistrar:1",
			ServiceId:   "urn:microsoft.com:serviceId:X_MS_MediaReceiverRegistrar",
		},
		SCPD: mediaReceiverRegistrarDescription,
	},
}

// The control URL for every service is the same. We're able to infer the desired service from the request headers.
func init() {
	for _, s := range services {
		s.ControlURL = serviceControlURL
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
				fmt.Fprintf(os.Stderr, "%s %s\r\n", r.Method, r.RequestURI)
				r.Header.Write(os.Stderr)
				fmt.Fprintln(os.Stderr)
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
		log.Printf("error creating ssdp server on %s: %s", if_.Name, err)
		return
	}
	defer s.Close()
	log.Println("started SSDP on", if_.Name)
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		if err := s.Serve(); err != nil {
			log.Printf("%q: %q\n", if_.Name, err)
		}
	}()
	select {
	case <-me.closed:
		// Returning will close the server.
	case <-stopped:
	}
}

var startTime time.Time

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
	FFProbeCache   Cache
	closed         chan struct{}
	ssdpStopped    chan struct{}
	// The service SOAP handler keyed by service URN.
	services   map[string]UPnPService
	LogHeaders bool
	// Disable transcoding, and the resource elements implied in the CDS.
	NoTranscode bool
	// Force transcoding to certain format of the 'transcodes' map
	ForceTranscodeTo string
	// Disable media probing with ffprobe
	NoProbe bool
	Icons   []Icon
	// Stall event subscription requests until they drop. A workaround for
	// some bad clients.
	StallEventSubscribe bool
	// Time interval between SSPD announces
	NotifyInterval time.Duration
	// Ignore hidden files and directories
	IgnoreHidden bool
	// Ingnore unreadable files and directories
	IgnoreUnreadable bool
	// White list of clients
	AllowedIpNets []*net.IPNet
}

// UPnP SOAP service.
type UPnPService interface {
	Handle(action string, argsXML []byte, r *http.Request) (respArgs [][2]string, err error)
	Subscribe(callback []*url.URL, timeoutSeconds int) (sid string, actualTimeout int, err error)
	Unsubscribe(sid string) error
}

type Cache interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

type dummyFFProbeCache struct{}

func (dummyFFProbeCache) Set(interface{}, interface{}) {}

func (dummyFFProbeCache) Get(interface{}) (interface{}, bool) {
	return nil, false
}

// Public definition so that external modules can persist cache contents.
type FfprobeCacheItem struct {
	Key   ffmpegInfoCacheKey
	Value *ffprobe.Info
}

// update the UPnP object fields from ffprobe data
// priority is given the format section, and then the streams sequentially
func itemExtra(item *upnpav.Object, info *ffprobe.Info) {
	setFromTags := func(m map[string]interface{}) {
		for key, val := range m {
			setIfUnset := func(s *string) {
				if *s == "" {
					*s = val.(string)
				}
			}
			switch strings.ToLower(key) {
			case "tag:artist":
				setIfUnset(&item.Artist)
			case "tag:album":
				setIfUnset(&item.Album)
			case "tag:genre":
				setIfUnset(&item.Genre)
			}
		}
	}
	setFromTags(info.Format)
	for _, m := range info.Streams {
		setFromTags(m)
	}
}

type ffmpegInfoCacheKey struct {
	Path    string
	ModTime int64
}

func transcodeResources(host, path, resolution, duration string) (ret []upnpav.Resource) {
	ret = make([]upnpav.Resource, 0, len(transcodes))
	for k, v := range transcodes {
		ret = append(ret, upnpav.Resource{
			ProtocolInfo: fmt.Sprintf("http-get:*:%s:%s", v.mimeType, dlna.ContentFeatures{
				SupportTimeSeek: true,
				Transcoded:      true,
				ProfileName:     v.DLNAProfileName,
			}.String()),
			URL: (&url.URL{
				Scheme: "http",
				Host:   host,
				Path:   resPath,
				RawQuery: url.Values{
					"path":      {path},
					"transcode": {k},
				}.Encode(),
			}).String(),
			Resolution: resolution,
			Duration:   duration,
		})
	}
	return
}

func parseDLNARangeHeader(val string) (ret dlna.NPTRange, err error) {
	if !strings.HasPrefix(val, "npt=") {
		err = errors.New("bad prefix")
		return
	}
	ret, err = dlna.ParseNPTRange(val[len("npt="):])
	if err != nil {
		return
	}
	return
}

// Determines the time-based range to transcode, and sets the appropriate
// headers. Returns !ok if there was an error and the caller should stop
// handling the request.
func handleDLNARange(w http.ResponseWriter, hs http.Header) (r dlna.NPTRange, partialResponse, ok bool) {
	if len(hs[http.CanonicalHeaderKey(dlna.TimeSeekRangeDomain)]) == 0 {
		ok = true
		return
	}
	partialResponse = true
	h := hs.Get(dlna.TimeSeekRangeDomain)
	r, err := parseDLNARangeHeader(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Passing an exact NPT duration seems to cause trouble pass the "iono"
	// (*) duration instead.
	//
	// TODO: Check that the request range can't already have /.
	w.Header().Set(dlna.TimeSeekRangeDomain, h+"/*")
	ok = true
	return
}

func (me *Server) serveDLNATranscode(w http.ResponseWriter, r *http.Request, path_ string, ts transcodeSpec, tsname string) {
	w.Header().Set(dlna.TransferModeDomain, "Streaming")
	w.Header().Set("content-type", ts.mimeType)
	w.Header().Set(dlna.ContentFeaturesDomain, (dlna.ContentFeatures{
		Transcoded:      true,
		SupportTimeSeek: true,
	}).String())
	// If a range of any kind is given, we have to respond with 206 if we're
	// interpreting that range. Since only the DLNA range is handled in this
	// function, it alone determines if we'll give a partial response.
	range_, partialResponse, ok := handleDLNARange(w, r.Header)
	if !ok {
		return
	}
	ffInfo, _ := me.ffmpegProbe(path_)
	if ffInfo != nil {
		if duration, err := ffInfo.Duration(); err == nil {
			s := fmt.Sprintf("%f", duration.Seconds())
			w.Header().Set("content-duration", s)
			w.Header().Set("x-content-duration", s)
		}
	}
	stderrPath := func() string {
		u, _ := user.Current()
		return filepath.Join(u.HomeDir, ".dms", "log", tsname, filepath.Base(path_))
	}()
	os.MkdirAll(filepath.Dir(stderrPath), 0750)
	logFile, err := os.Create(stderrPath)
	if err != nil {
		log.Printf("couldn't create transcode log file: %s", err)
	} else {
		defer logFile.Close()
		log.Printf("logging transcode to %q", stderrPath)
	}
	p, err := ts.Transcode(path_, range_.Start, range_.End-range_.Start, logFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer p.Close()
	// I recently switched this to returning 200 if no range is specified for
	// pure UPnP clients. It's possible that DLNA clients will *always* expect
	// 206. It appears the HTTP standard requires that 206 only be used if a
	// response is not interpreting any range headers.
	w.WriteHeader(func() int {
		if partialResponse {
			return http.StatusPartialContent
		} else {
			return http.StatusOK
		}
	}())
	io.Copy(w, p)
}

func init() {
	startTime = time.Now()
}

func getDefaultFriendlyName() string {
	return fmt.Sprintf("%s: %s on %s", rootDeviceModelName, func() string {
		user, err := user.Current()
		if err != nil {
			log.Panicf("getDefaultFriendlyName could not get username: %s", err)
		}
		return user.Name
	}(), func() string {
		name, err := os.Hostname()
		if err != nil {
			log.Panicf("getDefaultFriendlyName could not get hostname: %s", err)
		}
		return name
	}())
}

func xmlMarshalOrPanic(value interface{}) []byte {
	ret, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Panicf("xmlMarshalOrPanic failed to marshal %v: %s", value, err)
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
	fmt.Fprintln(os.Stderr, code)
	for k, v := range me.Header() {
		fmt.Fprintln(os.Stderr, k, v)
	}
	fmt.Fprintln(os.Stderr)
	me.loggedHeader = true
}

func (me *mitmRespWriter) Write(b []byte) (int, error) {
	if !me.loggedHeader {
		me.doLogHeader(200)
	}
	return me.ResponseWriter.Write(b)
}

func (me *mitmRespWriter) CloseNotify() <-chan bool {
	return me.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Set the SCPD serve paths.
func init() {
	for _, s := range services {
		lastInd := strings.LastIndex(s.ServiceId, ":")
		p := path.Join("/scpd", s.ServiceId[lastInd+1:])
		s.SCPDURL = p + ".xml"
	}
}

// Install handlers to serve SCPD for each UPnP service.
func handleSCPDs(mux *http.ServeMux) {
	for _, s := range services {
		mux.HandleFunc(s.SCPDURL, func(serviceDesc string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-type", `text/xml; charset="utf-8"`)
				http.ServeContent(w, r, "", startTime, bytes.NewReader([]byte(serviceDesc)))
			}
		}(s.SCPD))
	}
}

// Marshal SOAP response arguments into a response XML snippet.
func marshalSOAPResponse(sa upnp.SoapAction, args [][2]string) []byte {
	soapArgs := make([]soap.Arg, 0, len(args))
	for _, arg := range args {
		argName, value := arg[0], arg[1]
		soapArgs = append(soapArgs, soap.Arg{
			XMLName: xml.Name{Local: argName},
			Value:   value,
		})
	}
	return []byte(fmt.Sprintf(`<u:%[1]sResponse xmlns:u="%[2]s">%[3]s</u:%[1]sResponse>`, sa.Action, sa.ServiceURN.String(), xmlMarshalOrPanic(soapArgs)))
}

// Handle a SOAP request and return the response arguments or UPnP error.
func (me *Server) soapActionResponse(sa upnp.SoapAction, actionRequestXML []byte, r *http.Request) ([][2]string, error) {
	service, ok := me.services[sa.Type]
	if !ok {
		// TODO: What's the invalid service error?!
		return nil, upnp.Errorf(upnp.InvalidActionErrorCode, "Invalid service: %s", sa.Type)
	}
	return service.Handle(sa.Action, actionRequestXML, r)
}

// Handle a service control HTTP request.
func (me *Server) serviceControlHandler(w http.ResponseWriter, r *http.Request) {
	found := false
	clientIp, _, _ := net.SplitHostPort(r.RemoteAddr)
	if zoneDelimiterIdx := strings.Index(clientIp, "%"); zoneDelimiterIdx != -1 {
		// IPv6 addresses may have the form address%zone (e.g. ::1%eth0)
		clientIp = clientIp[:zoneDelimiterIdx]
	}
	for _, ipnet := range me.AllowedIpNets {
		if ipnet.Contains(net.ParseIP(clientIp)) {
			found = true
		}
	}
	if !found {
		log.Printf("not allowed client %s, %+v", clientIp, me.AllowedIpNets)
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
	// log.Println(r.UserAgent())
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
		log.Print(err)
	}
}

func safeFilePath(root, given string) string {
	return filepath.Join(root, filepath.FromSlash(path.Clean("/" + given))[1:])
}

func (s *Server) filePath(_path string) string {
	return safeFilePath(s.RootObjectPath, _path)
}

func (me *Server) serveIcon(w http.ResponseWriter, r *http.Request) {
	filePath := me.filePath(r.URL.Query().Get("path"))
	c := r.URL.Query().Get("c")
	if c == "" {
		c = "png"
	}
	cmd := exec.Command("ffmpegthumbnailer", "-i", filePath, "-o", "/dev/stdout", "-c"+c)
	// cmd.Stderr = os.Stderr
	body, err := cmd.Output()
	if err != nil {
		// serve 1st Icon if no ffmpegthumbnailer
		w.Header().Set("Content-Type", me.Icons[0].Mimetype)
		http.ServeContent(w, r, "", time.Time{}, me.Icons[0].ReadSeeker)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(body))
}

func (server *Server) contentDirectoryInitialEvent(urls []*url.URL, sid string) {
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
	eventingLogger.Print(string(body))
	for _, _url := range urls {
		bodyReader := bytes.NewReader(body)
		req, err := http.NewRequest("NOTIFY", _url.String(), bodyReader)
		if err != nil {
			log.Printf("Could not create a request to notify %s: %s", _url.String(), err)
			continue
		}
		req.Header["CONTENT-TYPE"] = []string{`text/xml; charset="utf-8"`}
		req.Header["NT"] = []string{"upnp:event"}
		req.Header["NTS"] = []string{"upnp:propchange"}
		req.Header["SID"] = []string{sid}
		req.Header["SEQ"] = []string{"0"}
		// req.Header["TRANSFER-ENCODING"] = []string{"chunked"}
		// req.ContentLength = int64(bodyReader.Len())
		eventingLogger.Print(req.Header)
		eventingLogger.Print("starting notify")
		resp, err := http.DefaultClient.Do(req)
		eventingLogger.Print("finished notify")
		if err != nil {
			log.Printf("Could not notify %s: %s", _url.String(), err)
			continue
		}
		eventingLogger.Print(resp)
		b, _ := ioutil.ReadAll(resp.Body)
		eventingLogger.Println(string(b))
		resp.Body.Close()
	}
}

var eventingLogger = log.New(ioutil.Discard, "", 0)

func (server *Server) contentDirectoryEventSubHandler(w http.ResponseWriter, r *http.Request) {
	if server.StallEventSubscribe {
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
		<-w.(http.CloseNotifier).CloseNotify()
		eventingLogger.Printf("stalled subscribe connection went away after %s", time.Since(t))
		return
	}
	// The following code is a work in progress. It partially implements
	// the spec on eventing but hasn't been completed as I have nothing to
	// test it with.
	eventingLogger.Print(r.Header)
	service := server.services["ContentDirectory"]
	eventingLogger.Println(r.RemoteAddr, r.Method, r.Header.Get("SID"))
	if r.Method == "SUBSCRIBE" && r.Header.Get("SID") == "" {
		urls := upnp.ParseCallbackURLs(r.Header.Get("CALLBACK"))
		eventingLogger.Println(urls)
		var timeout int
		fmt.Sscanf(r.Header.Get("TIMEOUT"), "Second-%d", &timeout)
		eventingLogger.Println(timeout, r.Header.Get("TIMEOUT"))
		sid, timeout, _ := service.Subscribe(urls, timeout)
		w.Header()["SID"] = []string{sid}
		w.Header()["TIMEOUT"] = []string{fmt.Sprintf("Second-%d", timeout)}
		// TODO: Shouldn't have to do this to get headers logged.
		w.WriteHeader(http.StatusOK)
		go func() {
			time.Sleep(100 * time.Millisecond)
			server.contentDirectoryInitialEvent(urls, sid)
		}()
	} else if r.Method == "SUBSCRIBE" {
		http.Error(w, "meh", http.StatusPreconditionFailed)
	} else {
		eventingLogger.Printf("unhandled event method: %s", r.Method)
	}
}

func (server *Server) initMux(mux *http.ServeMux) {
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("content-type", "text/html")
		err := rootTmpl.Execute(resp, struct {
			Readonly bool
			Path     string
		}{
			true,
			server.RootObjectPath,
		})
		if err != nil {
			log.Println(err)
		}
	})
	mux.HandleFunc(contentDirectoryEventSubURL, server.contentDirectoryEventSubHandler)
	mux.HandleFunc(iconPath, server.serveIcon)
	mux.HandleFunc(resPath, func(w http.ResponseWriter, r *http.Request) {
		filePath := server.filePath(r.URL.Query().Get("path"))
		if ignored, err := server.IgnorePath(filePath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if ignored {
			http.Error(w, "no such object", http.StatusNotFound)
			return
		}
		var k string
		if server.ForceTranscodeTo != "" {
			k = server.ForceTranscodeTo
		} else {
			k = r.URL.Query().Get("transcode")
		}
		mimeType, err := MimeTypeByPath(filePath)
		if k == "" || mimeType.IsImage() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", string(mimeType))
			http.ServeFile(w, r, filePath)
			return
		}
		if server.NoTranscode {
			http.Error(w, "transcodes disabled", http.StatusNotFound)
			return
		}
		spec, ok := transcodes[k]
		if !ok {
			http.Error(w, fmt.Sprintf("bad transcode spec key: %s", k), http.StatusBadRequest)
			return
		}
		server.serveDLNATranscode(w, r, filePath, spec, k)
	})
	mux.HandleFunc(rootDescPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", `text/xml; charset="utf-8"`)
		w.Header().Set("content-length", fmt.Sprint(len(server.rootDescXML)))
		w.Header().Set("server", serverField)
		w.Write(server.rootDescXML)
	})
	handleSCPDs(mux)
	mux.HandleFunc(serviceControlURL, server.serviceControlHandler)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	// DeviceIcons
	iconHandl := func(w http.ResponseWriter, r *http.Request) {
		idStr := path.Base(r.URL.Path)
		id, _ := strconv.Atoi(idStr)
		if id < 0 || id >= len(server.Icons) {
			id = 0
		}
		di := server.Icons[id]
		w.Header().Set("Content-Type", di.Mimetype)
		http.ServeContent(w, r, "", time.Time{}, di.ReadSeeker)
	}
	for i, _ := range server.Icons {
		mux.HandleFunc(fmt.Sprintf("%s/%d", deviceIconPath, i), iconHandl)
	}
}

func (s *Server) initServices() (err error) {
	urn, err := upnp.ParseServiceType(services[0].ServiceType)
	if err != nil {
		return
	}
	urn1, err := upnp.ParseServiceType(services[1].ServiceType)
	if err != nil {
		return
	}
	urn2, err := upnp.ParseServiceType(services[2].ServiceType)
	if err != nil {
		return
	}
	s.services = map[string]UPnPService{
		urn.Type: &contentDirectoryService{
			Server: s,
		},
		urn1.Type: &connectionManagerService{
			Server: s,
		},
		urn2.Type: &mediaReceiverRegistrarService{
			Server: s,
		},
	}
	return
}

func (srv *Server) Init() (err error) {
	if err = srv.initServices(); err != nil {
		return
	}
	srv.closed = make(chan struct{})
	if srv.FriendlyName == "" {
		srv.FriendlyName = getDefaultFriendlyName()
	}
	if srv.HTTPConn == nil {
		srv.HTTPConn, err = net.Listen("tcp", "")
		if err != nil {
			return
		}
	}
	if srv.Interfaces == nil {
		ifs, err := net.Interfaces()
		if err != nil {
			log.Print(err)
		}
		var tmp []net.Interface
		for _, if_ := range ifs {
			if if_.Flags&net.FlagUp == 0 || if_.MTU <= 0 {
				continue
			}
			tmp = append(tmp, if_)
		}
		srv.Interfaces = tmp
	}
	if srv.FFProbeCache == nil {
		srv.FFProbeCache = dummyFFProbeCache{}
	}
	srv.httpServeMux = http.NewServeMux()
	srv.rootDeviceUUID = makeDeviceUuid(srv.FriendlyName)
	srv.rootDescXML, err = xml.MarshalIndent(
		upnp.DeviceDesc{
			NSDLNA:      "urn:schemas-dlna-org:device-1-0",
			NSSEC:       "http://www.sec.co.kr/dlna",
			SpecVersion: upnp.SpecVersion{Major: 1, Minor: 0},
			Device: upnp.Device{
				DeviceType:   rootDeviceType,
				FriendlyName: srv.FriendlyName,
				Manufacturer: "Matt Joiner <anacrolix@gmail.com>",
				ModelName:    rootDeviceModelName,
				UDN:          srv.rootDeviceUUID,
				VendorXML: `
     <dlna:X_DLNACAP/>
     <dlna:X_DLNADOC>DMS-1.50</dlna:X_DLNADOC>
     <dlna:X_DLNADOC>M-DMS-1.50</dlna:X_DLNADOC>
     <sec:ProductCap>smi,DCM10,getMediaInfo.sec,getCaptionInfo.sec</sec:ProductCap>
     <sec:X_ProductCap>smi,DCM10,getMediaInfo.sec,getCaptionInfo.sec</sec:X_ProductCap>`,
				ServiceList: func() (ss []upnp.Service) {
					for _, s := range services {
						ss = append(ss, s.Service)
					}
					return
				}(),
				IconList: func() (ret []upnp.Icon) {
					for i, di := range srv.Icons {
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
				PresentationURL: "/",
			},
		},
		" ", "  ")
	if err != nil {
		return
	}
	srv.rootDescXML = append([]byte(`<?xml version="1.0"?>`), srv.rootDescXML...)
	log.Println("HTTP srv on", srv.HTTPConn.Addr())
	srv.initMux(srv.httpServeMux)
	srv.ssdpStopped = make(chan struct{})
	return nil
}

// Deprecated: Use Init and then Run. There's a race calling Close on a Server that's had Serve
// called on it.
func (srv *Server) Serve() (err error) {
	err = srv.Init()
	if err != nil {
		return
	}
	return srv.Run()
}

func (srv *Server) Run() (err error) {
	go func() {
		srv.doSSDP()
		close(srv.ssdpStopped)
	}()
	return srv.serveHTTP()
}

func (srv *Server) Close() (err error) {
	close(srv.closed)
	err = srv.HTTPConn.Close()
	<-srv.ssdpStopped
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

// Can return nil info with nil err if an earlier Probe gave an error.
func (srv *Server) ffmpegProbe(path string) (info *ffprobe.Info, err error) {
	// We don't want relative paths in the cache.
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	key := ffmpegInfoCacheKey{path, fi.ModTime().UnixNano()}
	value, ok := srv.FFProbeCache.Get(key)
	if !ok {
		info, err = ffprobe.Run(path)
		err = suppressFFmpegProbeDataErrors(err)
		srv.FFProbeCache.Set(key, info)
		return
	}
	info = value.(*ffprobe.Info)
	return
}

// IgnorePath detects if a file/directory should be ignored.
func (server *Server) IgnorePath(path string) (bool, error) {
	if !filepath.IsAbs(path) {
		return false, fmt.Errorf("Path must be absolute: %s", path)
	}
	if server.IgnoreHidden {
		if hidden, err := isHiddenPath(path); err != nil {
			return false, err
		} else if hidden {
			log.Print(path, " ignored: hidden")
			return true, nil
		}
	}
	if server.IgnoreUnreadable {
		if readable, err := isReadablePath(path); err != nil {
			return false, err
		} else if !readable {
			log.Print(path, " ignored: unreadable")
			return true, nil
		}
	}
	return false, nil
}

func tryToOpenPath(path string) (bool, error) {
	// Ugly but portable way to check if we can open a file/directory
	if fh, err := os.Open(path); err == nil {
		fh.Close()
		return true, nil
	} else if !os.IsPermission(err) {
		return false, err
	}
	return false, nil
}
