package dms

import (
	"net/http"

	"github.com/anacrolix/dms/upnp"
)

type mediaReceiverRegistrarService struct {
	*Server
	upnp.Eventing
}

func (mrrs *mediaReceiverRegistrarService) Handle(action string, argsXML []byte, r *http.Request) ([][2]string, error) {
	switch action {
	case "IsAuthorized", "IsValidated":
		return [][2]string{
			{"Result", "1"},
		}, nil
	case "RegisterDevice":
		return [][2]string{
			{"RegistrationRespMsg", mrrs.rootDeviceUUID},
		}, nil
		//		return nil, nil
	default:
		return nil, upnp.InvalidActionError
	}
}
