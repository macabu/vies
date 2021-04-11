package vies

import "errors"

const (
	invalidInput               = "INVALID_INPUT"
	invalidRequesterInfo       = "INVALID_REQUESTER_INFO"
	serviceUnavailable         = "SERVICE_UNAVAILABLE"
	msUnavailable              = "MS_UNAVAILABLE"
	timeout                    = "TIMEOUT"
	vatBlocked                 = "VAT_BLOCKED"
	ipBlocked                  = "IP_BLOCKED"
	globalMaxConcurrentReq     = "GLOBAL_MAX_CONCURRENT_REQ"
	globalMaxConcurrentReqTime = "GLOBAL_MAX_CONCURRENT_REQ_TIME"
	msMaxConcurrentReq         = "MS_MAX_CONCURRENT_REQ"
	msMaxConcurrentReqTime     = "MS_MAX_CONCURRENT_REQ_TIME"
)

var (
	ErrInvalidInput               = errors.New(invalidInput)
	ErrInvalidRequesterInfo       = errors.New(invalidRequesterInfo)
	ErrServiceUnavailable         = errors.New(serviceUnavailable)
	ErrMsUnavailable              = errors.New(msUnavailable)
	ErrTimeout                    = errors.New(timeout)
	ErrVATBlocked                 = errors.New(vatBlocked)
	ErrIpBlocked                  = errors.New(ipBlocked)
	ErrGlobalMaxConcurrentReq     = errors.New(globalMaxConcurrentReq)
	ErrGlobalMaxConcurrentReqTime = errors.New(globalMaxConcurrentReqTime)
	ErrMsMaxConcurrentReq         = errors.New(msMaxConcurrentReq)
	ErrMsMaxConcurrentReqTime     = errors.New(msMaxConcurrentReqTime)
)

var toSentinelError = map[string]error{
	invalidInput:               ErrInvalidInput,
	invalidRequesterInfo:       ErrInvalidRequesterInfo,
	serviceUnavailable:         ErrServiceUnavailable,
	msUnavailable:              ErrMsUnavailable,
	timeout:                    ErrTimeout,
	vatBlocked:                 ErrVATBlocked,
	ipBlocked:                  ErrIpBlocked,
	globalMaxConcurrentReq:     ErrGlobalMaxConcurrentReq,
	globalMaxConcurrentReqTime: ErrGlobalMaxConcurrentReqTime,
	msMaxConcurrentReq:         ErrMsMaxConcurrentReq,
	msMaxConcurrentReqTime:     ErrMsMaxConcurrentReqTime,
}
