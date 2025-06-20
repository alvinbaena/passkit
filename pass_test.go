package passkit

import (
	"testing"
	"time"
)

func getBasicRelevantDate() PassRelevantDate {
	t := time.Now()
	return PassRelevantDate{StartDate: (*time.Time)(&t)}
}

func TestPassRelevantDate_Invalid(t *testing.T) {
	prd := new(PassRelevantDate)

	if prd.IsValid() {
		t.Errorf("PassRelevantDate should be invalid")
	}

	if len(prd.GetValidationErrors()) == 0 {
		t.Errorf("PassRelevantDate should have errors")
	}
}

func TestPassRelevantDate_JSONMarshallingSingleDate(t *testing.T) {
	unixTimeUTC := time.Date(2025, time.June, 19, 01, 23, 45, 0, time.UTC)
	prd := PassRelevantDate{StartDate: (*time.Time)(&unixTimeUTC)}

	prdJSON, err := prd.toJSON()

	if err != nil {
		t.Errorf("PassRelevantDate JSON Marshalling failed. Reason %v", err)
	}

	prdJSONString := string(prdJSON)
	expected := "{\"relevantDate\":\"2025-06-19T01:23:45Z\"}"
	if prdJSONString != expected {
		t.Errorf("PassRelevantDate JSON did not matched expected format")
	}
}

func TestPassRelevantDate_JSONMarshallingDateRange(t *testing.T) {
	unixTimeUTC := time.Date(2025, time.June, 19, 01, 23, 45, 0, time.UTC)
	prd := PassRelevantDate{StartDate: (*time.Time)(&unixTimeUTC), EndDate: (*time.Time)(&unixTimeUTC)}

	prdJSON, err := prd.toJSON()

	if err != nil {
		t.Errorf("PassRelevantDate JSON Marshalling failed. Reason %v", err)
	}

	prdJSONString := string(prdJSON)
	expected := "{\"startDate\":\"2025-06-19T01:23:45Z\",\"endDate\":\"2025-06-19T01:23:45Z\"}"
	if prdJSONString != expected {
		t.Errorf("PassRelevantDate JSON did not matched expected format")
	}
}

func getBasicPersonalization() Personalization {
	return Personalization{
		Description:        "Description for pass",
		TermsAndConditions: "WTF",
		RequiredPersonalizationFields: []PassPersonalizationField{
			PassPersonalizationFieldName,
			PassPersonalizationFieldEmailAddress,
		},
	}
}

func TestPersonalization_GetSet(t *testing.T) {
	prs := getBasicPersonalization()

	if !prs.IsValid() {
		t.Errorf("Personalization should be valid. Reason: %v", prs.GetValidationErrors())
	}

	if len(prs.RequiredPersonalizationFields) == 0 {
		t.Errorf("Personalization should have fields")
	}

	if len(prs.GetValidationErrors()) != 0 {
		t.Errorf("Personalization should not have errors. Have %v", len(prs.GetValidationErrors()))
	}
}

func TestPersonalization_OptionalFields(t *testing.T) {
	prs := getBasicPersonalization()
	prs.TermsAndConditions = ""

	if !prs.IsValid() {
		t.Errorf("Personalization should be valid. Reason: %v", prs.GetValidationErrors())
	}

	if len(prs.RequiredPersonalizationFields) == 0 {
		t.Errorf("Personalization should have fields")
	}

	if len(prs.GetValidationErrors()) != 0 {
		t.Errorf("Personalization should not have errors. Have %v", len(prs.GetValidationErrors()))
	}
}

func TestPersonalization_Invalid(t *testing.T) {
	prs := getBasicPersonalization()
	prs.Description = ""

	t.Logf("%d, %v", len(prs.GetValidationErrors()), prs.GetValidationErrors())
	if prs.IsValid() {
		t.Errorf("Personalization should be invalid")
	}

	if len(prs.GetValidationErrors()) != 1 {
		t.Errorf("Personalization should have only one validation error. Have %v", len(prs.GetValidationErrors()))
	}

	prs.RequiredPersonalizationFields = []PassPersonalizationField{}

	t.Logf("%d, %v", len(prs.GetValidationErrors()), prs.GetValidationErrors())
	if prs.IsValid() {
		t.Errorf("Personalization should be invalid")
	}

	if len(prs.GetValidationErrors()) != 2 {
		t.Errorf("Personalization should have only two validation errors. Have %v", len(prs.GetValidationErrors()))
	}
}

func TestPersonalization_JSON(t *testing.T) {
	prs := getBasicPersonalization()
	_, err := prs.toJSON()

	if err != nil {
		t.Errorf("Error marshalling json. %v", err)
	}
}

func getBasicBarcode() Barcode {
	return Barcode{
		Format:          BarcodeFormatQR,
		AltText:         "Alt Text",
		Message:         "message",
		MessageEncoding: "utf-8",
	}
}

func TestBarcode_GetSet(t *testing.T) {
	bar := getBasicBarcode()

	if !bar.IsValid() {
		t.Errorf("Barcode should be valid. Reason: %v", bar.GetValidationErrors())
	}

	if len(bar.GetValidationErrors()) != 0 {
		t.Errorf("Barcode should not have errors. Have %v", bar.GetValidationErrors())
	}
}

func TestBarcode_NoMessage(t *testing.T) {
	bar := getBasicBarcode()
	bar.Message = ""

	t.Logf("%d, %v", len(bar.GetValidationErrors()), bar.GetValidationErrors())
	if bar.IsValid() {
		t.Errorf("Barcode should be invalid")
	}

	if len(bar.GetValidationErrors()) != 1 {
		t.Errorf("Barcode should have only one error")
	}
}

func TestBarcode_NoAltText(t *testing.T) {
	bar := getBasicBarcode()
	bar.AltText = ""

	t.Logf("%d, %v", len(bar.GetValidationErrors()), bar.GetValidationErrors())
	if bar.IsValid() {
		t.Errorf("Barcode should be invalid")
	}

	if len(bar.GetValidationErrors()) != 1 {
		t.Errorf("Barcode should have only one error")
	}
}

func TestBarcode_NoEncoding(t *testing.T) {
	bar := getBasicBarcode()
	bar.MessageEncoding = ""

	t.Logf("%d, %v", len(bar.GetValidationErrors()), bar.GetValidationErrors())
	if bar.IsValid() {
		t.Errorf("Barcode should be invalid")
	}

	if len(bar.GetValidationErrors()) != 1 {
		t.Errorf("Barcode should have only one error")
	}
}

func TestBarcode_NoFormat(t *testing.T) {
	bar := getBasicBarcode()
	bar.Format = ""

	t.Logf("%d, %v", len(bar.GetValidationErrors()), bar.GetValidationErrors())
	if bar.IsValid() {
		t.Errorf("Barcode should be invalid")
	}

	if len(bar.GetValidationErrors()) != 1 {
		t.Errorf("Barcode should have only one error")
	}
}

func getBasicBeacon() Beacon {
	return Beacon{
		Major:         3,
		Minor:         29,
		ProximityUUID: "123456789-abcdefghijklmnopqrstuwxyz",
		RelevantText:  "County of Zadar",
	}
}

func TestBeacon_GetSet(t *testing.T) {
	bea := getBasicBeacon()

	if !bea.IsValid() {
		t.Errorf("Beacon should be valid. Reason: %v", bea.GetValidationErrors())
	}

	if len(bea.GetValidationErrors()) != 0 {
		t.Errorf("Beacon should not have errors. Have %v", len(bea.GetValidationErrors()))
	}
}

func TestBeacon_NoUUID(t *testing.T) {
	bea := getBasicBeacon()
	bea.ProximityUUID = ""

	t.Logf("%d, %v", len(bea.GetValidationErrors()), bea.GetValidationErrors())
	if bea.IsValid() {
		t.Errorf("Beacon should be invalid")
	}

	if len(bea.GetValidationErrors()) != 1 {
		t.Errorf("Beacon should have one error. Have: %v", len(bea.GetValidationErrors()))
	}
}

func getBasicField() Field {
	return Field{
		Key:               "key",
		ChangeMessage:     "Changed %@",
		Label:             "Label",
		TextAlignment:     TextAlignmentCenter,
		AttributedValue:   "<a href='http://example.com/customers/123'>Edit my profile</a>",
		DataDetectorTypes: []DataDetectorType{DataDetectorTypeAddress},
	}
}

func TestField_GetSet(t *testing.T) {
	field := getBasicField()
	field.Value = "test"

	if !field.IsValid() {
		t.Errorf("Field should be valid. Reason: %v", field.GetValidationErrors())
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Beacon should have no errors. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_Values(t *testing.T) {
	field := getBasicField()
	field.Value = "Value"

	if !field.IsValid() {
		t.Errorf("Field should be valid. Reason: %v", field.GetValidationErrors())
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = 1

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = 1.0

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int8(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int8(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int16(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int32(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int64(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = float32(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = time.Now()

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

}

func TestField_Currency(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "COP"
	field.Value = 43000

	if !field.IsValid() {
		t.Errorf("Field should be valid. Reason: %v", field.GetValidationErrors())
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should have no errors. Have %v", len(field.GetValidationErrors()))
	}
}

func TestField_CurrencyValues(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "USD"
	field.Value = 1

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = 1.0

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int8(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int8(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int16(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int32(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = int64(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

	field.Value = float32(1)

	if !field.IsValid() {
		t.Errorf("Field should be valid")
	}

	if len(field.GetValidationErrors()) != 0 {
		t.Errorf("Field should not have errors. Have %v", len(field.GetValidationErrors()))
	}

}

func TestField_CurrencyAndNumberFormat(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "COP"
	field.Value = 1
	field.NumberStyle = NumberStyleDecimal

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_CurrencyAndDateFormat(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "COP"
	field.DateStyle = DateStyleFull
	field.Value = 1

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_CurrencyAndTimeFormat(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "COP"
	field.TimeStyle = DateStyleFull
	field.Value = 1

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_CurrencyAndNotNumber(t *testing.T) {
	field := getBasicField()
	field.CurrencyCode = "COP"
	field.Value = "1.2321"

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_NumberAndDateStyleSet(t *testing.T) {
	field := getBasicField()
	field.NumberStyle = NumberStyleDecimal
	field.DateStyle = DateStyleFull
	field.Value = "string"

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_NoKey(t *testing.T) {
	field := getBasicField()
	field.Value = "Value"
	field.Key = ""

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_NoValue(t *testing.T) {
	field := getBasicField()
	field.Value = nil

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error %v", len(field.GetValidationErrors()))
	}
}

func TestField_InvalidValueType(t *testing.T) {
	field := getBasicField()
	field.Value = false

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func TestField_InvalidChangeMessage(t *testing.T) {
	field := getBasicField()
	field.ChangeMessage = "Fake"
	field.Value = "Test"

	t.Logf("%d, %v", len(field.GetValidationErrors()), field.GetValidationErrors())
	if field.IsValid() {
		t.Errorf("Field should be invalid")
	}

	if len(field.GetValidationErrors()) != 1 {
		t.Errorf("Field should have one error. Have: %v", len(field.GetValidationErrors()))
	}
}

func getBasicLocation() Location {
	return Location{
		Longitude:    1.0,
		Latitude:     1.2,
		Altitude:     1.3,
		RelevantText: "text",
	}
}

func TestLocation_GetSet(t *testing.T) {
	location := getBasicLocation()

	t.Logf("%d, %v", len(location.GetValidationErrors()), location.GetValidationErrors())
	if !location.IsValid() {
		t.Errorf("Location should be valid. Reason: %v", location.GetValidationErrors())
	}

	if len(location.GetValidationErrors()) != 0 {
		t.Errorf("Location should have no errors. Have: %v", len(location.GetValidationErrors()))
	}
}

func getBasicPass() Pass {
	exp := time.Now()
	f := getBasicField()
	f.Value = "string"

	return Pass{
		FormatVersion:              1,
		OrganizationName:           "Org",
		Description:                "test",
		AppLaunchURL:               "app://open",
		MaxDistance:                99999,
		Voided:                     false,
		UserInfo:                   map[string]interface{}{"name": "John Doe"},
		ExpirationDate:             &exp,
		Barcodes:                   []Barcode{getBasicBarcode()},
		SerialNumber:               "1234",
		PassTypeIdentifier:         "test",
		TeamIdentifier:             "TEAM1",
		AuthenticationToken:        "asldadilno21o31n41lkasndio123",
		Generic:                    &GenericPass{PrimaryFields: []Field{f}},
		AssociatedStoreIdentifiers: []int64{123},
	}
}

func TestPass_GetSet(t *testing.T) {
	pass := getBasicPass()
	pass.GroupingIdentifier = ""

	if !pass.IsValid() {
		t.Errorf("Pass should be valid. Reason: %v", pass.GetValidationErrors())
	}

	if len(pass.GetValidationErrors()) != 0 {
		t.Errorf("Pass should have no errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_JSON(t *testing.T) {
	prs := getBasicPass()
	_, err := prs.toJSON()

	if err != nil {
		t.Errorf("Error marshalling json. %v", err)
	}
}

func TestPass_MissingSerial(t *testing.T) {
	pass := getBasicPass()
	pass.SerialNumber = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MissingPassTypeId(t *testing.T) {
	pass := getBasicPass()
	pass.PassTypeIdentifier = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}
	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MissingTeamID(t *testing.T) {
	pass := getBasicPass()
	pass.TeamIdentifier = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MissingDescription(t *testing.T) {
	pass := getBasicPass()
	pass.Description = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MissingFormat(t *testing.T) {
	pass := getBasicPass()
	pass.FormatVersion = 0

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MissingOrganization(t *testing.T) {
	pass := getBasicPass()
	pass.OrganizationName = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_NoPass(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = nil
	pass.Coupon = nil
	pass.BoardingPass = nil
	pass.EventTicket = nil
	pass.StoreCard = nil

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidAuthToken(t *testing.T) {
	pass := getBasicPass()
	pass.AuthenticationToken = ""

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_MultiplePass(t *testing.T) {
	pass := getBasicPass()
	f := getBasicField()
	f.Value = "test"

	pass.Coupon = &Coupon{GenericPass: &GenericPass{HeaderFields: []Field{f}}}
	pass.BoardingPass = &BoardingPass{
		GenericPass: &GenericPass{HeaderFields: []Field{f}},
		TransitType: TransitTypeAir,
	}
	pass.EventTicket = &EventTicket{GenericPass: &GenericPass{HeaderFields: []Field{f}}}
	pass.StoreCard = &StoreCard{GenericPass: &GenericPass{HeaderFields: []Field{f}}}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidGeneric(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = &GenericPass{HeaderFields: []Field{getBasicField()}}
	pass.Coupon = nil
	pass.BoardingPass = nil
	pass.EventTicket = nil
	pass.StoreCard = nil

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.Generic.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidBoarding(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = nil
	pass.Coupon = nil
	pass.EventTicket = nil
	pass.StoreCard = nil
	pass.BoardingPass = &BoardingPass{
		GenericPass: &GenericPass{HeaderFields: []Field{getBasicField()}},
		TransitType: TransitTypeAir,
	}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidStoreCard(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = nil
	pass.Coupon = nil
	pass.BoardingPass = nil
	pass.EventTicket = nil
	pass.StoreCard = &StoreCard{GenericPass: &GenericPass{HeaderFields: []Field{getBasicField()}}}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidEventTicket(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = nil
	pass.Coupon = nil
	pass.BoardingPass = nil
	pass.StoreCard = nil
	pass.EventTicket = &EventTicket{GenericPass: &GenericPass{HeaderFields: []Field{getBasicField()}}}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidCoupon(t *testing.T) {
	pass := getBasicPass()
	pass.Generic = nil
	pass.Coupon = &Coupon{GenericPass: &GenericPass{HeaderFields: []Field{getBasicField()}}}
	pass.BoardingPass = nil
	pass.StoreCard = nil
	pass.EventTicket = nil

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_NoStoreId(t *testing.T) {
	pass := getBasicPass()
	pass.AssociatedStoreIdentifiers = []int64{}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_InvalidGroupIdentifier(t *testing.T) {
	pass := getBasicPass()
	pass.GroupingIdentifier = "213131"

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPass_SetForegroundColorHex(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetForegroundColorHex("#000000")

	if pass.ForegroundColor != "rgb(0,0,0)" {
		t.Errorf("Foreground color invalid")
	}
}

func TestPass_SetInvalidForegroundColorHex(t *testing.T) {
	pass := getBasicPass()
	err := pass.SetForegroundColorHex("#PPPPPP")

	t.Logf("Hex error: %v", err)
	if err == nil {
		t.Errorf("Hex parse should fail")
	}
}

func TestPass_SetBackgroundColorHex(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetBackgroundColorHex("#000000")

	if pass.BackgroundColor != "rgb(0,0,0)" {
		t.Errorf("Background color invalid")
	}
}

func TestPass_SetInvalidBackgroundColorHex(t *testing.T) {
	pass := getBasicPass()
	err := pass.SetBackgroundColorHex("#PPPPPP")

	t.Logf("Hex error: %v", err)
	if err == nil {
		t.Errorf("Hex parse should fail")
	}
}

func TestPass_SetLabelColorHex(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetLabelColorHex("#000000")

	if pass.LabelColor != "rgb(0,0,0)" {
		t.Errorf("Label color invalid")
	}
}

func TestPass_SetInvalidLabelColorHex(t *testing.T) {
	pass := getBasicPass()
	err := pass.SetLabelColorHex("#PPPPPP")

	t.Logf("Hex error: %v", err)
	if err == nil {
		t.Errorf("Hex parse should fail")
	}
}

func TestPass_SetForegroundColorRGB(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetForegroundColorRGB(0, 0, 0)

	if pass.ForegroundColor != "rgb(0,0,0)" {
		t.Errorf("Foreground color invalid")
	}
}

func TestPass_SetBackgroundColorRGB(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetBackgroundColorRGB(0, 0, 0)

	if pass.BackgroundColor != "rgb(0,0,0)" {
		t.Errorf("Foreground color invalid")
	}
}

func TestPass_SetLabelColorRGB(t *testing.T) {
	pass := getBasicPass()
	_ = pass.SetLabelColorRGB(0, 0, 0)

	if pass.LabelColor != "rgb(0,0,0)" {
		t.Errorf("Foreground color invalid")
	}
}

func TestBoardingPass_NoTransitType(t *testing.T) {
	pass := getBasicPass()
	f := getBasicField()
	f.Value = "1"

	pass.Generic = nil
	pass.Coupon = nil
	pass.EventTicket = nil
	pass.StoreCard = nil
	pass.BoardingPass = &BoardingPass{
		GenericPass: &GenericPass{HeaderFields: []Field{f}},
	}

	t.Logf("%d, %v", len(pass.GetValidationErrors()), pass.GetValidationErrors())
	if pass.IsValid() {
		t.Errorf("Pass should be invalid")
	}

	if len(pass.GetValidationErrors()) != 1 {
		t.Errorf("Pass should have one errors. Have: %v", len(pass.GetValidationErrors()))
	}
}

func TestPWAssociatedApp_GetSet(t *testing.T) {
	pw := PWAssociatedApp{}

	if !pw.IsValid() {
		t.Errorf("PWAssociatedApp should be valid. Reason: %v", pw.GetValidationErrors())
	}

	if len(pw.GetValidationErrors()) != 0 {
		t.Errorf("PWAssociatedApp should have no errors. Have: %v", len(pw.GetValidationErrors()))
	}
}

func TestSetRelevantDates(t *testing.T) {
	pass := getBasicPass()

	earlyDate := time.Date(2025, time.June, 19, 01, 23, 45, 0, time.UTC)
	pdr := getBasicRelevantDate()
	earlierPDR := PassRelevantDate{StartDate: (*time.Time)(&earlyDate)}

	var dates []PassRelevantDate

	dates = append(dates, pdr)
	dates = append(dates, earlierPDR)
	pass.SetRelevantDates(dates)

	if pass.RelevantDate != &earlyDate {
		t.Error("relevantDate was not set to earliest date within slice")
	}
}
