package kurento

import (
	"errors"
	"fmt"
	"log"
)

type IWebRtcEndpoint interface {
	GatherCandidates() error
	AddIceCandidate(candidate IceCandidate) error
}

// WebRtcEndpoint interface. This type of "Endpoint" offers media streaming using
// WebRTC.
type WebRtcEndpoint struct {
	BaseRtpEndpoint

	// Address of the STUN server (Only IP address are supported)
	StunServerAddress string

	// Port of the STUN server
	StunServerPort int
}

// Return contructor params to be called by "Create".
func (elem *WebRtcEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Init the gathering of ICE candidates.
// It must be called after SdpEndpoint::generateOffer or SdpEndpoint::processOffer
func (elem *WebRtcEndpoint) GatherCandidates() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "gatherCandidates",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Otherwise we want to wait for the other candidates
	// Returns error or nil
	if response.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil
}

// Provide a remote ICE candidate
func (elem *WebRtcEndpoint) AddIceCandidate(candidate IceCandidate) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "candidate", candidate)
	setIfNotEmpty(params, "__type__", "IceCandidate")
	setIfNotEmpty(params, "__module__", "kurento")
	log.Println(params)

	req["params"] = map[string]interface{}{
		"operation":       "addIceCandidate",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return errors.New(fmt.Sprintf("Error: [%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil

}
