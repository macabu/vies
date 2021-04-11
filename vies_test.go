package vies_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/macabu/vies"
)

const (
	validMockedResponse   string = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><checkVatResponse xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types"><countryCode>NL</countryCode><vatNumber>100</vatNumber><requestDate>2021-04-10+02:00</requestDate><valid>true</valid><name>John Doe</name><address>123 Main St, Anytown, UK</address></checkVatResponse></soap:Body></soap:Envelope>`
	invalidMockedResponse string = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><soap:Fault><faultcode>soap:Server</faultcode><faultstring>INVALID_REQUESTER_INFO</faultstring></soap:Fault></soap:Body></soap:Envelope>`
)

func TestCheckVAT(t *testing.T) {
	validReq := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validReq {
			fmt.Fprintln(w, validMockedResponse)
		} else {
			fmt.Fprintln(w, invalidMockedResponse)
		}
	}))
	defer ts.Close()

	svc := vies.NewService(ts.URL)
	vat, err := svc.CheckVAT("NL1234")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err.Error())
	}
	if vat == nil {
		t.Error("expected vat to not be nil, got nil")
	}

	validReq = false
	svc = vies.NewServiceWithTimeout(ts.URL, time.Minute)
	vat, err = svc.CheckVAT("NL1234")
	if err == nil {
		t.Errorf("expected err to not be nil, got nil")
	}
	if !errors.Is(err, vies.ErrInvalidRequesterInfo) {
		t.Errorf("expected err to be %v, got %v", vies.ErrInvalidRequesterInfo, err.Error())
	}
	if vat != nil {
		t.Errorf("expected vat to be nil, got %+v", vat)
	}
}

func TestCheckVATWithCountry(t *testing.T) {
	validReq := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validReq {
			fmt.Fprintln(w, validMockedResponse)
		} else {
			fmt.Fprintln(w, invalidMockedResponse)
		}
	}))
	defer ts.Close()

	svc := vies.NewService(ts.URL)
	vat, err := svc.CheckVATWithCountry("NL", "NL1234")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err.Error())
	}
	if vat == nil {
		t.Error("expected vat to not be nil, got nil")
	}

	validReq = false
	svc = vies.NewServiceWithTimeout(ts.URL, time.Minute)
	vat, err = svc.CheckVATWithCountry("NL", "NL1234")
	if err == nil {
		t.Errorf("expected err to not be nil, got nil")
	}
	if !errors.Is(err, vies.ErrInvalidRequesterInfo) {
		t.Errorf("expected err to be %v, got %v", vies.ErrInvalidRequesterInfo, err.Error())
	}
	if vat != nil {
		t.Errorf("expected vat to be nil, got %+v", vat)
	}
}

func TestCheckVATWithTestEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	svc := vies.NewService(vies.TestEndpoint)

	testCases := []struct {
		vatNumber     string
		expectedValid bool
		expectedError error
	}{
		{vatNumber: "NL100", expectedValid: true},
		{vatNumber: "NL201", expectedError: vies.ErrInvalidInput},
		{vatNumber: "NL202", expectedError: vies.ErrInvalidRequesterInfo},
		{vatNumber: "NL301", expectedError: vies.ErrMsUnavailable},
		{vatNumber: "NL302", expectedError: vies.ErrTimeout},
		{vatNumber: "NL400", expectedError: vies.ErrVATBlocked},
		{vatNumber: "NL401", expectedError: vies.ErrIpBlocked},
		{vatNumber: "NL500", expectedError: vies.ErrGlobalMaxConcurrentReq},
		{vatNumber: "NL501", expectedError: vies.ErrGlobalMaxConcurrentReqTime},
		{vatNumber: "NL600", expectedError: vies.ErrMsMaxConcurrentReq},
		{vatNumber: "NL601", expectedError: vies.ErrMsMaxConcurrentReqTime},
		{vatNumber: "NLanythingelse", expectedError: vies.ErrServiceUnavailable},
	}

	for _, tc := range testCases {
		t.Run("with vat "+tc.vatNumber, func(t *testing.T) {
			vat, err := svc.CheckVAT(tc.vatNumber)
			if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
				t.Errorf("expected err to be %v, got %v", tc.expectedError, err.Error())
			}
			if vat != nil && vat.Valid != tc.expectedValid {
				t.Errorf("expected vat validation flag to be %v, got %v", tc.expectedValid, vat.Valid)
			}

			<-time.After(500 * time.Millisecond) // don't stress VIES
		})
	}
}

func TestCheckVATWithCountryWithTestEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	svc := vies.NewService(vies.TestEndpoint)

	testCases := []struct {
		countryCode   string
		vatNumber     string
		expectedValid bool
		expectedError error
	}{
		{countryCode: "NL", vatNumber: "100", expectedValid: true},
		{countryCode: "NL", vatNumber: "201", expectedError: vies.ErrInvalidInput},
		{countryCode: "NL", vatNumber: "202", expectedError: vies.ErrInvalidRequesterInfo},
		{countryCode: "NL", vatNumber: "301", expectedError: vies.ErrMsUnavailable},
		{countryCode: "NL", vatNumber: "302", expectedError: vies.ErrTimeout},
		{countryCode: "NL", vatNumber: "400", expectedError: vies.ErrVATBlocked},
		{countryCode: "NL", vatNumber: "401", expectedError: vies.ErrIpBlocked},
		{countryCode: "NL", vatNumber: "500", expectedError: vies.ErrGlobalMaxConcurrentReq},
		{countryCode: "NL", vatNumber: "501", expectedError: vies.ErrGlobalMaxConcurrentReqTime},
		{countryCode: "NL", vatNumber: "600", expectedError: vies.ErrMsMaxConcurrentReq},
		{countryCode: "NL", vatNumber: "601", expectedError: vies.ErrMsMaxConcurrentReqTime},
		{countryCode: "NL", vatNumber: "anythingelse", expectedError: vies.ErrServiceUnavailable},
	}

	for _, tc := range testCases {
		t.Run("with country "+tc.countryCode+" and vat "+tc.vatNumber, func(t *testing.T) {
			vat, err := svc.CheckVATWithCountry(tc.countryCode, tc.vatNumber)
			if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
				t.Errorf("expected err to be %v, got %v", tc.expectedError, err.Error())
			}
			if vat != nil && vat.Valid != tc.expectedValid {
				t.Errorf("expected vat validation flag to be %v, got %v", tc.expectedValid, vat.Valid)
			}

			<-time.After(500 * time.Millisecond) // don't stress VIES
		})
	}
}
