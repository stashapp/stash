package dms

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anacrolix/dms/dlna"
	"github.com/anacrolix/dms/misc"
	"github.com/anacrolix/dms/upnp"
	"github.com/anacrolix/dms/upnpav"
	"github.com/anacrolix/ffprobe"
)

type contentDirectoryService struct {
	*Server
	upnp.Eventing
}

func (cds *contentDirectoryService) updateIDString() string {
	return fmt.Sprintf("%d", uint32(os.Getpid()))
}

// Turns the given entry and DMS host into a UPnP object. A nil object is
// returned if the entry is not of interest.
func (me *contentDirectoryService) cdsObjectToUpnpavObject(cdsObject object, fileInfo os.FileInfo, host, userAgent string) (ret interface{}, err error) {
	entryFilePath := cdsObject.FilePath()
	ignored, err := me.IgnorePath(entryFilePath)
	if err != nil {
		return
	}
	if ignored {
		return
	}
	obj := upnpav.Object{
		ID:         cdsObject.ID(),
		Restricted: 1,
		ParentID:   cdsObject.ParentID(),
	}
	if fileInfo.IsDir() {
		obj.Class = "object.container.storageFolder"
		obj.Title = fileInfo.Name()
		ret = upnpav.Container{Object: obj, ChildCount: me.objectChildCount(cdsObject)}
		return
	}
	if !fileInfo.Mode().IsRegular() {
		log.Printf("%s ignored: non-regular file", cdsObject.FilePath())
		return
	}
	mimeType, err := MimeTypeByPath(entryFilePath)
	if err != nil {
		return
	}
	if !mimeType.IsMedia() {
		log.Printf("%s ignored: non-media file (%s)", cdsObject.FilePath(), mimeType)
		return
	}
	iconURI := (&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   iconPath,
		RawQuery: url.Values{
			"path": {cdsObject.Path},
		}.Encode(),
	}).String()
	obj.Icon = iconURI
	// TODO(anacrolix): This might not be necessary due to item res image
	// element.
	obj.AlbumArtURI = iconURI
	obj.Class = "object.item." + mimeType.Type() + "Item"
	var (
		ffInfo        *ffprobe.Info
		nativeBitrate uint
		resDuration   string
	)
	if !me.NoProbe {
		ffInfo, probeErr := me.ffmpegProbe(entryFilePath)
		switch probeErr {
		case nil:
			if ffInfo != nil {
				nativeBitrate, _ = ffInfo.Bitrate()
				if d, err := ffInfo.Duration(); err == nil {
					resDuration = misc.FormatDurationSexagesimal(d)
				}
			}
		case ffprobe.ExeNotFound:
		default:
			log.Printf("error probing %s: %s", entryFilePath, probeErr)
		}
	}
	if obj.Title == "" {
		obj.Title = fileInfo.Name()
	}
	resolution := func() string {
		if ffInfo != nil {
			for _, strm := range ffInfo.Streams {
				if strm["codec_type"] != "video" {
					continue
				}
				width := strm["width"]
				height := strm["height"]
				return fmt.Sprintf("%.0fx%.0f", width, height)
			}
		}
		return ""
	}()
	item := upnpav.Item{
		Object: obj,
		// Capacity: 1 for raw, 1 for icon, plus transcodes.
		Res: make([]upnpav.Resource, 0, 2+len(transcodes)),
	}
	item.Res = append(item.Res, upnpav.Resource{
		URL: (&url.URL{
			Scheme: "http",
			Host:   host,
			Path:   resPath,
			RawQuery: url.Values{
				"path": {cdsObject.Path},
			}.Encode(),
		}).String(),
		ProtocolInfo: fmt.Sprintf("http-get:*:%s:%s", mimeType, dlna.ContentFeatures{
			SupportRange: true,
		}.String()),
		Bitrate:    nativeBitrate,
		Duration:   resDuration,
		Size:       uint64(fileInfo.Size()),
		Resolution: resolution,
	})
	if mimeType.IsVideo() {
		if !me.NoTranscode {
			item.Res = append(item.Res, transcodeResources(host, cdsObject.Path, resolution, resDuration)...)
		}
	}
	if mimeType.IsVideo() || mimeType.IsImage() {
		item.Res = append(item.Res, upnpav.Resource{
			URL: (&url.URL{
				Scheme: "http",
				Host:   host,
				Path:   iconPath,
				RawQuery: url.Values{
					"path": {cdsObject.Path},
					"c":    {"jpeg"},
				}.Encode(),
			}).String(),
			ProtocolInfo: "http-get:*:image/jpeg:DLNA.ORG_PN=JPEG_TN",
		})
	}
	ret = item
	return
}

// Returns all the upnpav objects in a directory.
func (me *contentDirectoryService) readContainer(o object, host, userAgent string) (ret []interface{}, err error) {
	sfis := sortableFileInfoSlice{
		// TODO(anacrolix): Dig up why this special cast was added.
		FoldersLast: strings.Contains(userAgent, `AwoX/1.1`),
	}
	sfis.fileInfoSlice, err = o.readDir()
	if err != nil {
		return
	}
	sort.Sort(sfis)
	for _, fi := range sfis.fileInfoSlice {
		child := object{path.Join(o.Path, fi.Name()), me.RootObjectPath}
		obj, err := me.cdsObjectToUpnpavObject(child, fi, host, userAgent)
		if err != nil {
			log.Printf("error with %s: %s", child.FilePath(), err)
			continue
		}
		if obj != nil {
			ret = append(ret, obj)
		}
	}
	return
}

type browse struct {
	ObjectID       string
	BrowseFlag     string
	Filter         string
	StartingIndex  int
	RequestedCount int
}

// ContentDirectory object from ObjectID.
func (me *contentDirectoryService) objectFromID(id string) (o object, err error) {
	o.Path, err = url.QueryUnescape(id)
	if err != nil {
		return
	}
	if o.Path == "0" {
		o.Path = "/"
	}
	o.Path = path.Clean(o.Path)
	if !path.IsAbs(o.Path) {
		err = fmt.Errorf("bad ObjectID %v", o.Path)
		return
	}
	o.RootObjectPath = me.RootObjectPath
	return
}

func (me *contentDirectoryService) Handle(action string, argsXML []byte, r *http.Request) ([][2]string, error) {
	host := r.Host
	userAgent := r.UserAgent()
	switch action {
	case "GetSystemUpdateID":
		return [][2]string{
			{"Id", me.updateIDString()},
		}, nil
	case "GetSortCapabilities":
		return [][2]string{
			{"SortCaps", "dc:title"},
		}, nil
	case "Browse":
		var browse browse
		if err := xml.Unmarshal([]byte(argsXML), &browse); err != nil {
			return nil, err
		}
		obj, err := me.objectFromID(browse.ObjectID)
		if err != nil {
			return nil, upnp.Errorf(upnpav.NoSuchObjectErrorCode, err.Error())
		}
		switch browse.BrowseFlag {
		case "BrowseDirectChildren":
			objs, err := me.readContainer(obj, host, userAgent)
			if err != nil {
				return nil, upnp.Errorf(upnpav.NoSuchObjectErrorCode, err.Error())
			}
			totalMatches := len(objs)
			objs = objs[func() (low int) {
				low = browse.StartingIndex
				if low > len(objs) {
					low = len(objs)
				}
				return
			}():]
			if browse.RequestedCount != 0 && int(browse.RequestedCount) < len(objs) {
				objs = objs[:browse.RequestedCount]
			}
			result, err := xml.Marshal(objs)
			if err != nil {
				return nil, err
			}
			return [][2]string{
				{"Result", didl_lite(string(result))},
				{"NumberReturned", fmt.Sprint(len(objs))},
				{"TotalMatches", fmt.Sprint(totalMatches)},
				{"UpdateID", me.updateIDString()},
			}, nil
		case "BrowseMetadata":
			fileInfo, err := os.Stat(obj.FilePath())
			if err != nil {
				if os.IsNotExist(err) {
					return nil, &upnp.Error{
						Code: upnpav.NoSuchObjectErrorCode,
						Desc: err.Error(),
					}
				}
				return nil, err
			}
			upnp, err := me.cdsObjectToUpnpavObject(obj, fileInfo, host, userAgent)
			if err != nil {
				return nil, err
			}
			buf, err := xml.Marshal(upnp)
			if err != nil {
				return nil, err
			}
			return [][2]string{
				{"Result", didl_lite(func() string { return string(buf) }())},
				{"NumberReturned", "1"},
				{"TotalMatches", "1"},
				{"UpdateID", me.updateIDString()},
			}, nil
		default:
			return nil, upnp.Errorf(upnp.ArgumentValueInvalidErrorCode, "unhandled browse flag: %v", browse.BrowseFlag)
		}
	case "GetSearchCapabilities":
		return [][2]string{
			{"SearchCaps", ""},
		}, nil
	// Samsung Extensions
	case "X_GetFeatureList":
	// TODO: make it dependable on model
	// https://github.com/1100101/minidlna/blob/ca6dbba18390ad6f8b8d7b7dbcf797dbfd95e2db/upnpsoap.c#L2153-L2199
		return [][2]string{
			{"FeatureList", `<Features xmlns="urn:schemas-upnp-org:av:avs" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:schemas-upnp-org:av:avs http://www.upnp.org/schemas/av/avs.xsd">
	<Feature name="samsung.com_BASICVIEW" version="1">
		<container id="0" type="object.item.audioItem"/> // "A"
		<container id="0" type="object.item.videoItem"/> // "V"
		<container id="0" type="object.item.imageItem"/> // "I"
	</Feature>
</Features>`}}, nil
	case "X_SetBookmark":
		// just ignore
		return [][2]string{}, nil
	default:
		return nil, upnp.InvalidActionError
	}
}

// Represents a ContentDirectory object.
type object struct {
	Path           string // The cleaned, absolute path for the object relative to the server.
	RootObjectPath string
}

// Returns the number of children this object has, such as for a container.
func (cds *contentDirectoryService) objectChildCount(me object) int {
	objs, err := cds.readContainer(me, "", "")
	if err != nil {
		log.Printf("error reading container: %s", err)
	}
	return len(objs)
}

func (cds *contentDirectoryService) objectHasChildren(obj object) bool {
	return cds.objectChildCount(obj) != 0
}

// Returns the actual local filesystem path for the object.
func (o *object) FilePath() string {
	return filepath.Join(o.RootObjectPath, filepath.FromSlash(o.Path))
}

// Returns the ObjectID for the object. This is used in various ContentDirectory actions.
func (o object) ID() string {
	if !path.IsAbs(o.Path) {
		log.Panicf("Relative object path: %s", o.Path)
	}
	if len(o.Path) == 1 {
		return "0"
	}
	return url.QueryEscape(o.Path)
}

func (o *object) IsRoot() bool {
	return o.Path == "/"
}

// Returns the object's parent ObjectID. Fortunately it can be deduced from the
// ObjectID (for now).
func (o object) ParentID() string {
	if o.IsRoot() {
		return "-1"
	}
	o.Path = path.Dir(o.Path)
	return o.ID()
}

// This function exists rather than just calling os.(*File).Readdir because I
// want to stat(), not lstat() each entry.
func (o *object) readDir() (fis []os.FileInfo, err error) {
	dirPath := o.FilePath()
	dirFile, err := os.Open(dirPath)
	if err != nil {
		return
	}
	defer dirFile.Close()
	var dirContent []string
	dirContent, err = dirFile.Readdirnames(-1)
	if err != nil {
		return
	}
	fis = make([]os.FileInfo, 0, len(dirContent))
	for _, file := range dirContent {
		fi, err := os.Stat(filepath.Join(dirPath, file))
		if err != nil {
			continue
		}
		fis = append(fis, fi)
	}
	return
}

type sortableFileInfoSlice struct {
	fileInfoSlice []os.FileInfo
	FoldersLast   bool
}

func (me sortableFileInfoSlice) Len() int {
	return len(me.fileInfoSlice)
}

func (me sortableFileInfoSlice) Less(i, j int) bool {
	if me.fileInfoSlice[i].IsDir() && !me.fileInfoSlice[j].IsDir() {
		return !me.FoldersLast
	}
	if !me.fileInfoSlice[i].IsDir() && me.fileInfoSlice[j].IsDir() {
		return me.FoldersLast
	}
	return strings.ToLower(me.fileInfoSlice[i].Name()) < strings.ToLower(me.fileInfoSlice[j].Name())
}

func (me sortableFileInfoSlice) Swap(i, j int) {
	me.fileInfoSlice[i], me.fileInfoSlice[j] = me.fileInfoSlice[j], me.fileInfoSlice[i]
}
