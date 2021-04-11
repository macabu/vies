package vies

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// ProductionEndpoint is the production endpoint for the VIES VAT service.
	ProductionEndpoint string = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

	// TestEndpoint is the test endpoint for the VIES VAT service.
	TestEndpoint string = "http://ec.europa.eu/taxation_customs/vies/services/checkVatTestService"
)

type service struct {
	endpoint   string
	httpClient *http.Client
}

// NewService creates a VIES service with a default timeout of 10 seconds for the HTTP request.
// You can pass a custom endpoint for testing purposes or proxying, or use vies.ProductionEndpoint
// or vies.TestEndpoint to reach the VIES VAT service.
func NewService(endpoint string) *service {
	return &service{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// NewServiceWithTimeout creates a VIES service with the specified timeout for the HTTP request.
// You can pass a custom endpoint for testing purposes or proxying, or use vies.ProductionEndpoint
// or vies.TestEndpoint to reach the VIES VAT service.
func NewServiceWithTimeout(endpoint string, timeout time.Duration) *service {
	return &service{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type vatData struct {
	CountryCode string
	VatNumber   string
	Valid       bool
	Name        string
	Address     string
}

// CheckVAT given a VAT number.
// This will infer the country based on the two initial digits from the VAT number.
//
// It will not necessarily return a non-nil error if the VAT is not valid.
//
// Example:
//    svc := vies.NewService(vies.TestEndpoint)
//    vat, err := svc.CheckVAT("NL123456789B01")
//    if err != nil {
//        return
//    }
//    fmt.Println("Valid?:", vat.Valid)
func (s *service) CheckVAT(fullVatNumber string) (*vatData, error) {
	if len(fullVatNumber) <= 2 {
		return nil, ErrInvalidInput
	}

	countryCode := fullVatNumber[0:2]
	vatNumber := fullVatNumber[2:]

	return s.checkVat(countryCode, vatNumber)
}

// CheckVAT given a VAT number and a country.
//
// The VAT number can either contain the country code or not, as long as it matches with
// the countryCode provided.
// If you want to use the country inside the VAT number, consider using CheckVAT instead.
//
// It will not necessarily return a non-nil error if the VAT is not valid.
//
// Example:
//    svc := vies.NewService(vies.TestEndpoint)
//    vat, err := svc.CheckVATWithCountry("NL", "123456789B01")
//    if err != nil {
//        return
//    }
//    fmt.Println("Valid?:", vat.Valid)
func (s *service) CheckVATWithCountry(countryCode, vatNumber string) (*vatData, error) {
	if len(vatNumber) <= 2 {
		return nil, ErrInvalidInput
	}

	if countryCode[0] == vatNumber[0] && countryCode[1] == vatNumber[1] {
		return s.checkVat(countryCode, vatNumber[2:])
	}

	return s.checkVat(countryCode, vatNumber)
}

func (s *service) checkVat(countryCode, vatNumber string) (*vatData, error) {
	if _, ok := validCountries[countryCode]; !ok {
		return nil, ErrInvalidInput
	}

	envelope, err := s.request(countryCode, vatNumber)
	if err != nil {
		return nil, err
	}

	return &vatData{
		CountryCode: envelope.Body.CheckVatResponse.CountryCode,
		VatNumber:   envelope.Body.CheckVatResponse.VatNumber,
		Valid:       envelope.Body.CheckVatResponse.Valid,
		Name:        envelope.Body.CheckVatResponse.Name,
		Address:     envelope.Body.CheckVatResponse.Address,
	}, nil
}

type transportResponse struct {
	Body struct {
		CheckVatResponse struct {
			CountryCode string `xml:"countryCode"`
			VatNumber   string `xml:"vatNumber"`
			Valid       bool   `xml:"valid"`
			Name        string `xml:"name"`
			Address     string `xml:"address"`
		} `xml:"checkVatResponse"`
		Fault struct {
			Message string `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

const requestBody string = `<s11:Envelope xmlns:s11="http://schemas.xmlsoap.org/soap/envelope/"><s11:Body><tns1:checkVat xmlns:tns1="urn:ec.europa.eu:taxud:vies:services:checkVat:types"><tns1:countryCode>%s</tns1:countryCode><tns1:vatNumber>%s</tns1:vatNumber></tns1:checkVat></s11:Body></s11:Envelope>`

func (s *service) request(countryCode, vatNumber string) (*transportResponse, error) {
	reqBody := fmt.Sprintf(requestBody, countryCode, vatNumber)
	req, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBufferString(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Add("Content-Type", "charset=utf-8")
	req.Header.Set("SOAPAction", "checkVat")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var envelope transportResponse
	if err := xml.Unmarshal([]byte(resBody), &envelope); err != nil {
		return nil, err
	}

	if envelope.Body.Fault.Message != "" {
		sentinelError, ok := toSentinelError[envelope.Body.Fault.Message]
		if !ok {
			sentinelError = ErrServiceUnavailable
		}

		return nil, fmt.Errorf("request failed with code %w", sentinelError)
	}

	return &envelope, nil
}
