package vies

import "sync"

var (
	once           sync.Once
	defaultService *service
)

// CheckVAT given a VAT number.
// This will infer the country based on the two initial digits from the VAT number.
//
// This is a default implementation and uses the default singleton client and production endpoint.
// The client is only created and allocated if you use these default implementations.
//
// It will not necessarily return a non-nil error if the VAT is not valid.
//
// Example:
//    vat, err := vies.CheckVAT("NL123456789B01")
//    if err != nil {
//        return
//    }
//    fmt.Println("Valid?:", vat.Valid)
func CheckVAT(fullVatNumber string) (*vatData, error) {
	once.Do(func() {
		defaultService = NewService(ProductionEndpoint)
	})

	return defaultService.CheckVAT(fullVatNumber)
}

// CheckVAT given a VAT number and a country.
//
// The VAT number can either contain the country code or not, as long as it matches with
// the countryCode provided.
// If you want to use the country inside the VAT number, consider using CheckVAT instead.
//
// This is a default implementation and uses the default singleton client and production endpoint.
// The client is only created and allocated if you use these default implementations.
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
func CheckVATWithCountry(countryCode, vatNumber string) (*vatData, error) {
	once.Do(func() {
		defaultService = NewService(ProductionEndpoint)
	})

	return defaultService.CheckVATWithCountry(countryCode, vatNumber)
}
