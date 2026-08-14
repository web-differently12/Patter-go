package engine

import (
	"context"
	"encoding/xml"
	"fmt"
)

// CarrierType defines the telecommunication carrier provider types
type CarrierType string

const (
	CarrierTwilio CarrierType = "twilio"
	CarrierTelnyx CarrierType = "telnyx"
	CarrierPlivo  CarrierType = "plivo"
)

// Carrier defines the interface for telecommunication carriers supporting outbound calling and custom markup generation
type Carrier interface {
	GetType() CarrierType
	InitiateCall(ctx context.Context, from, to, callbackURL string) (string, string, error)
	GenerateStreamResponse(wsURL string) (string, error)
}

// Twilio TwiML structural helpers
type TwiMLStream struct {
	XMLName xml.Name `xml:"Response"`
	Connect Connect  `xml:"Connect"`
}

type Connect struct {
	Stream Stream `xml:"Stream"`
}

type Stream struct {
	URL string `xml:"url,attr"`
}

// TwilioCarrier represents the Twilio implementation
type TwilioCarrier struct {
	accountSID string
	authToken  string
}

func NewTwilioCarrier(accountSID, authToken string) *TwilioCarrier {
	return &TwilioCarrier{
		accountSID: accountSID,
		authToken:  authToken,
	}
}

func (tc *TwilioCarrier) GetType() CarrierType {
	return CarrierTwilio
}

func (tc *TwilioCarrier) InitiateCall(ctx context.Context, from, to, callbackURL string) (string, string, error) {
	// Simulated/Mocked or dynamic twilio call initiation via Twilio REST Client API
	return "twilio_call_sid_" + to, "queued", nil
}

func (tc *TwilioCarrier) GenerateStreamResponse(wsURL string) (string, error) {
	twiml := TwiMLStream{
		Connect: Connect{
			Stream: Stream{
				URL: wsURL,
			},
		},
	}
	output, err := xml.Marshal(twiml)
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}

// TelnyxCarrier represents the Telnyx implementation
type TelnyxCarrier struct {
	apiKey string
}

func NewTelnyxCarrier(apiKey string) *TelnyxCarrier {
	return &TelnyxCarrier{apiKey: apiKey}
}

func (tc *TelnyxCarrier) GetType() CarrierType {
	return CarrierTelnyx
}

func (tc *TelnyxCarrier) InitiateCall(ctx context.Context, from, to, callbackURL string) (string, string, error) {
	return "telnyx_call_id_" + to, "queued", nil
}

func (tc *TelnyxCarrier) GenerateStreamResponse(wsURL string) (string, error) {
	texml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Connect>
        <Stream url="%s" />
    </Connect>
</Response>`, wsURL)
	return texml, nil
}

// PlivoCarrier represents the Plivo implementation
type PlivoCarrier struct {
	authID    string
	authToken string
}

func NewPlivoCarrier(authID, authToken string) *PlivoCarrier {
	return &PlivoCarrier{
		authID:    authID,
		authToken: authToken,
	}
}

func (pc *PlivoCarrier) GetType() CarrierType {
	return CarrierPlivo
}

func (pc *PlivoCarrier) InitiateCall(ctx context.Context, from, to, callbackURL string) (string, string, error) {
	return "plivo_call_uuid_" + to, "queued", nil
}

func (pc *PlivoCarrier) GenerateStreamResponse(wsURL string) (string, error) {
	plivoXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Stream keepCallActive="true">%s</Stream>
</Response>`, wsURL)
	return plivoXML, nil
}
