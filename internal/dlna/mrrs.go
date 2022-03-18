package dlna

import (
	"net/http"

	"github.com/anacrolix/dms/upnp"
)

type mediaReceiverRegistrarService struct {
	*Server
	upnp.Eventing
}

func (mrrs *mediaReceiverRegistrarService) Handle(action string, argsXML []byte, r *http.Request) (map[string]string, error) {
	switch action {
	case "IsAuthorized", "IsValidated":
		return map[string]string{
			"Result": "1",
		}, nil
	case "RegisterDevice":
		return map[string]string{
			"RegistrationRespMsg": mrrs.rootDeviceUUID,
		}, nil
	default:
		return nil, upnp.InvalidActionError
	}
}
