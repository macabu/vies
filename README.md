# VIES
An idiomatic Go API wrapper for the [VIES VAT service](https://ec.europa.eu/taxation_customs/vies/vatRequest.html), used to validate VAT numbers.  
  
- No external dependencies 
- Compatible with Go 1.11+

## Usage
You can create a service through `vies.NewService` or `vies.NewServiceWithTimeout`.  
In both cases you can pass a custom `endpoint`, if you want to use it for testing or proxying for example.  
Otherwise if you need to call the VIES services, there are two public variables that you can use:
- `vies.ProductionEndpoint`
- `vies.TestEndpoint` -> this goes to VIES' testing API.

### Example
```go
package main

import (
    "fmt"

    "github.com/macabu/vies"
)

func main() {
    testSvc := vies.NewService(vies.TestEndpoint)
    vat, err := testSvc.CheckVAT("NL810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)

    prodSvc := vies.NewService(vies.ProductionEndpoint)
    vat, err = prodSvc.CheckVAT("NL810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)
}
```

### Default Client
There is a default (singleton) client provided that calls the production endpoint of VIES.
You can use it in such a way:
```go
package main

import (
    "fmt"

    "github.com/macabu/vies"
)

func main() {
    vat, err := vies.CheckVAT("NL810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)

    vat, err = vies.CheckVATWithCountry("NL", "810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)

    // You can also use a variable for the country if you don't know its country code.
    vat, err = vies.CheckVATWithCountry(vies.TheNetherlands, "810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)

    // As long as the countryCode matches the prefixed country code on the VAT number, it is fine.
    // If they are different, and the prefix is not removed, it will most likely fail.
    // Consider using vies.CheckVAT if you have only the full VAT number with the country code.
    vat, err = vies.CheckVATWithCountry("NL", "NL810060255B01")
    if err != nil {
        panic(err)
    }
    fmt.Println("Valid?:", vat.Valid)
}
```

## Running Tests
### E2E (with VIES Test Service)
```sh
go test -count=1 -cover .
```

### Unit (only)
```sh
go test -short -count=1 -cover .
```
