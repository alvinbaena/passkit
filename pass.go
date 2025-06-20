package passkit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/colors.v1"
)

type BarcodeFormat string
type TextAlignment string
type DataDetectorType string
type DateStyle string
type NumberStyle string
type PassPersonalizationField string
type TransitType string
type EventType string

const (
	expectedAuthTokenLen = 16

	TextAlignmentLeft    TextAlignment = "PKTextAlignmentLeft"
	TextAlignmentCenter  TextAlignment = "PKTextAlignmentCenter"
	TextAlignmentRight   TextAlignment = "PKTextAlignmentRight"
	TextAlignmentNatural TextAlignment = "PKTextAlignmentNatural"

	BarcodeFormatQR      BarcodeFormat = "PKBarcodeFormatQR"
	BarcodeFormatPDF417  BarcodeFormat = "PKBarcodeFormatPDF417"
	BarcodeFormatAztec   BarcodeFormat = "PKBarcodeFormatAztec"
	BarcodeFormatCode128 BarcodeFormat = "PKBarcodeFormatCode128"

	DataDetectorTypePhoneNumber   DataDetectorType = "PKDataDetectorTypePhoneNumber"
	DataDetectorTypeLink          DataDetectorType = "PKDataDetectorTypeLink"
	DataDetectorTypeAddress       DataDetectorType = "PKDataDetectorTypeAddress"
	DataDetectorTypeCalendarEvent DataDetectorType = "PKDataDetectorTypeCalendarEvent"

	DateStyleNone   DateStyle = "PKDateStyleNone"
	DateStyleShort  DateStyle = "PKDateStyleShort"
	DateStyleMedium DateStyle = "PKDateStyleMedium"
	DateStyleLong   DateStyle = "PKDateStyleLong"
	DateStyleFull   DateStyle = "PKDateStyleFull"

	NumberStyleDecimal    NumberStyle = "PKNumberStyleDecimal"
	NumberStylePercent    NumberStyle = "PKNumberStylePercent"
	NumberStyleScientific NumberStyle = "PKNumberStyleScientific"
	NumberStyleSpellOut   NumberStyle = "PKNumberStyleSpellOut"

	PassPersonalizationFieldName         PassPersonalizationField = "PKPassPersonalizationFieldName"
	PassPersonalizationFieldPostalCode   PassPersonalizationField = "PKPassPersonalizationFieldPostalCode"
	PassPersonalizationFieldEmailAddress PassPersonalizationField = "PKPassPersonalizationFieldEmailAddress"
	PassPersonalizationFieldPhoneNumber  PassPersonalizationField = "PKPassPersonalizationFieldPhoneNumber"

	TransitTypeAir     TransitType = "PKTransitTypeAir"
	TransitTypeBoat    TransitType = "PKTransitTypeBoat"
	TransitTypeBus     TransitType = "PKTransitTypeBus"
	TransitTypeGeneric TransitType = "PKTransitTypeGeneric"
	TransitTypeTrain   TransitType = "PKTransitTypeTrain"

	EventTypeGeneric         EventType = "PKEventTypeGeneric"
	EventTypeLivePerformance EventType = "PKEventTypeLivePerformance"
	EventTypeMovie           EventType = "PKEventTypeMovie"
	EventTypeSports          EventType = "PKEventTypeSports"
	EventTypeConference      EventType = "PKEventTypeConference"
	EventTypeConvention      EventType = "PKEventTypeConvention"
	EventTypeWorkshop        EventType = "PKEventTypeWorkshop"
	EventTypeSocialGathering EventType = "PKEventTypeSocialGathering"
)

var (
	BarcodeTypesBeforeIos9 = [3]BarcodeFormat{BarcodeFormatQR, BarcodeFormatPDF417, BarcodeFormatAztec}
)

type Validateable interface {
	IsValid() bool
	GetValidationErrors() []string
}

// Pass Representation of https://developer.apple.com/documentation/walletpasses/pass
type Pass struct {
	FormatVersion              int                    `json:"formatVersion,omitempty"`
	SerialNumber               string                 `json:"serialNumber,omitempty"`
	PassTypeIdentifier         string                 `json:"passTypeIdentifier,omitempty"`
	WebServiceURL              string                 `json:"webServiceURL,omitempty"`
	AuthenticationToken        string                 `json:"authenticationToken,omitempty"`
	Description                string                 `json:"description,omitempty"`
	TeamIdentifier             string                 `json:"teamIdentifier,omitempty"`
	OrganizationName           string                 `json:"organizationName,omitempty"`
	LogoText                   string                 `json:"logoText,omitempty"`
	ForegroundColor            string                 `json:"foregroundColor,omitempty"`
	BackgroundColor            string                 `json:"backgroundColor,omitempty"`
	LabelColor                 string                 `json:"labelColor,omitempty"`
	GroupingIdentifier         string                 `json:"groupingIdentifier,omitempty"`
	Beacons                    []Beacon               `json:"beacons,omitempty"`
	Locations                  []Location             `json:"locations,omitempty"`
	Barcodes                   []Barcode              `json:"barcodes,omitempty"`
	EventTicket                *EventTicket           `json:"eventTicket,omitempty"`
	Coupon                     *Coupon                `json:"coupon,omitempty"`
	StoreCard                  *StoreCard             `json:"storeCard,omitempty"`
	BoardingPass               *BoardingPass          `json:"boardingPass,omitempty"`
	Generic                    *GenericPass           `json:"generic,omitempty"`
	AppLaunchURL               string                 `json:"appLaunchURL,omitempty"`
	AssociatedStoreIdentifiers []int64                `json:"associatedStoreIdentifiers,omitempty"`
	UserInfo                   map[string]interface{} `json:"userInfo,omitempty"`
	MaxDistance                int64                  `json:"maxDistance,omitempty"`
	RelevantDate               *time.Time             `json:"relevantDate,omitempty"`
	ExpirationDate             *time.Time             `json:"expirationDate,omitempty"`
	RelevantDates              []PassRelevantDate     `json:"relevantDates,omitempty"`
	Voided                     bool                   `json:"voided,omitempty"`
	Nfc                        *NFC                   `json:"nfc,omitempty"`
	SharingProhibited          bool                   `json:"sharingProhibited,omitempty"`
	Semantics                  *SemanticTag           `json:"semantics,omitempty"`

	//Private
	associatedApps []PWAssociatedApp
}

func (p *Pass) SetForegroundColorHex(hex string) error {
	h, err := colors.ParseHEX(hex)
	if err != nil {
		return err
	}

	p.ForegroundColor = h.ToRGB().String()
	return nil
}

func (p *Pass) SetForegroundColorRGB(r, g, b uint8) error {
	rgb, _ := colors.RGB(r, g, b)

	p.ForegroundColor = rgb.String()
	return nil
}

func (p *Pass) SetBackgroundColorHex(hex string) error {
	h, err := colors.ParseHEX(hex)
	if err != nil {
		return err
	}

	p.BackgroundColor = h.ToRGB().String()
	return nil
}

func (p *Pass) SetBackgroundColorRGB(r, g, b uint8) error {
	rgb, _ := colors.RGB(r, g, b)

	p.BackgroundColor = rgb.String()
	return nil
}

func (p *Pass) SetLabelColorHex(hex string) error {
	h, err := colors.ParseHEX(hex)
	if err != nil {
		return err
	}

	p.LabelColor = h.ToRGB().String()
	return nil
}

func (p *Pass) SetLabelColorRGB(r, g, b uint8) error {
	rgb, _ := colors.RGB(r, g, b)

	p.LabelColor = rgb.String()
	return nil
}

func (p *Pass) toJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Pass) IsValid() bool {
	return len(p.GetValidationErrors()) == 0
}

func (p *Pass) GetValidationErrors() []string {
	var validationErrors []string

	if strings.TrimSpace(p.SerialNumber) == "" || strings.TrimSpace(p.PassTypeIdentifier) == "" ||
		strings.TrimSpace(p.TeamIdentifier) == "" || strings.TrimSpace(p.Description) == "" ||
		p.FormatVersion == 0 || strings.TrimSpace(p.OrganizationName) == "" {

		validationErrors = append(validationErrors, fmt.Sprintf("Pass: Not all required Fields are set. SerialNumber: %q, PassTypeIdentifier: %q, teamIdentifier: %q, Description: ,%q, FormatVersion: %q, OrganizationName: %q", p.SerialNumber, p.PassTypeIdentifier, p.TeamIdentifier, p.Description, p.FormatVersion, p.OrganizationName))
	}

	if p.EventTicket == nil && p.BoardingPass == nil && p.Coupon == nil && p.StoreCard == nil && p.Generic == nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Pass: No pass was set. EventTicket: %v, BoardingPass: %v, Coupon: %v, StoreCard: %v, Generic: %v", p.EventTicket, p.BoardingPass, p.Coupon, p.StoreCard, p.Generic))
	}

	if p.EventTicket != nil && (p.BoardingPass != nil || p.Coupon != nil || p.StoreCard != nil || p.Generic != nil) {
		validationErrors = append(validationErrors, "Pass: Only one pass should be set")

	} else if p.BoardingPass != nil && (p.EventTicket != nil || p.Coupon != nil || p.StoreCard != nil || p.Generic != nil) {
		validationErrors = append(validationErrors, "Pass: Only one pass should be set")

	} else if p.Coupon != nil && (p.BoardingPass != nil || p.EventTicket != nil || p.StoreCard != nil || p.Generic != nil) {
		validationErrors = append(validationErrors, "Pass: Only one pass should be set")

	} else if p.StoreCard != nil && (p.BoardingPass != nil || p.Coupon != nil || p.EventTicket != nil || p.Generic != nil) {
		validationErrors = append(validationErrors, "Pass: Only one pass should be set")

	} else if p.Generic != nil && (p.BoardingPass != nil || p.Coupon != nil || p.StoreCard != nil || p.EventTicket != nil) {
		validationErrors = append(validationErrors, "Pass: Only one pass should be set")
	}

	if p.WebServiceURL != "" && (len(p.AuthenticationToken) < expectedAuthTokenLen) {
		validationErrors = append(validationErrors,
			"Pass: The authenticationToken needs to be at least "+strconv.Itoa(expectedAuthTokenLen)+" characters long")
	}

	if p.EventTicket != nil && !p.EventTicket.IsValid() {
		validationErrors = append(validationErrors, p.EventTicket.GetValidationErrors()...)
	} else if p.BoardingPass != nil && !p.BoardingPass.IsValid() {
		validationErrors = append(validationErrors, p.BoardingPass.GetValidationErrors()...)
	} else if p.Coupon != nil && !p.Coupon.IsValid() {
		validationErrors = append(validationErrors, p.Coupon.GetValidationErrors()...)
	} else if p.StoreCard != nil && !p.StoreCard.IsValid() {
		validationErrors = append(validationErrors, p.StoreCard.GetValidationErrors()...)
	} else if p.Generic != nil && !p.Generic.IsValid() {
		validationErrors = append(validationErrors, p.Generic.GetValidationErrors()...)
	}

	// If appLaunchURL key is present, the associatedStoreIdentifiers key must also be present
	if p.AppLaunchURL != "" && len(p.AssociatedStoreIdentifiers) == 0 {
		validationErrors = append(validationErrors, "Pass: The appLaunchURL requires associatedStoreIdentifiers to be specified")
	}

	if !(p.EventTicket == nil && p.BoardingPass == nil && p.Coupon == nil && p.StoreCard == nil && p.Generic == nil) {
		// groupingIdentifier key is optional for event tickets and boarding passes; otherwise not allowed
		if strings.TrimSpace(p.GroupingIdentifier) != "" && p.EventTicket == nil && p.BoardingPass == nil {
			validationErrors = append(validationErrors, "Pass: The groupingIdentifier is optional for event tickets and boarding passes, otherwise not allowed")
		}
	}

	if p.Beacons != nil {
		for _, b := range p.Beacons {
			if !b.IsValid() {
				validationErrors = append(validationErrors, b.GetValidationErrors()...)
			}
		}
	}

	if p.Barcodes != nil {
		for _, b := range p.Barcodes {
			if !b.IsValid() {
				validationErrors = append(validationErrors, b.GetValidationErrors()...)
			}
		}
	}

	if p.Semantics != nil && !p.Semantics.IsValid() {
		validationErrors = append(validationErrors, p.Semantics.GetValidationErrors()...)
	}

	if p.RelevantDates != nil {
		for _, prd := range p.RelevantDates {
			if !prd.IsValid() {
				validationErrors = append(validationErrors, prd.GetValidationErrors()...)
			}
		}
	}

	return validationErrors
}

func NewGenericPass() *GenericPass {
	return &GenericPass{}
}

// GenericPass Representation of https://developer.apple.com/documentation/walletpasses/pass/generic
type GenericPass struct {
	HeaderFields    []Field `json:"headerFields,omitempty"`
	PrimaryFields   []Field `json:"primaryFields,omitempty"`
	SecondaryFields []Field `json:"secondaryFields,omitempty"`
	AuxiliaryFields []Field `json:"auxiliaryFields,omitempty"`
	BackFields      []Field `json:"backFields,omitempty"`
}

func (gp *GenericPass) AddHeaderField(field Field) {
	gp.HeaderFields = append(gp.HeaderFields, field)
}

func (gp *GenericPass) AddPrimaryFields(field Field) {
	gp.PrimaryFields = append(gp.PrimaryFields, field)
}

func (gp *GenericPass) AddSecondaryFields(field Field) {
	gp.SecondaryFields = append(gp.SecondaryFields, field)
}

func (gp *GenericPass) AddAuxiliaryFields(field Field) {
	gp.AuxiliaryFields = append(gp.AuxiliaryFields, field)
}

func (gp *GenericPass) AddBackFields(field Field) {
	gp.BackFields = append(gp.BackFields, field)
}

func (gp *GenericPass) IsValid() bool {
	return len(gp.GetValidationErrors()) == 0
}

func (gp *GenericPass) GetValidationErrors() []string {
	var validationErrors []string

	var fields [][]Field
	fields = append(fields, gp.HeaderFields)
	fields = append(fields, gp.PrimaryFields)
	fields = append(fields, gp.SecondaryFields)
	fields = append(fields, gp.AuxiliaryFields)
	fields = append(fields, gp.BackFields)

	for _, fieldList := range fields {
		for _, field := range fieldList {
			if !field.IsValid() {
				validationErrors = append(validationErrors, field.GetValidationErrors()...)
			}
		}
	}

	return validationErrors
}

// BoardingPass Representation of https://developer.apple.com/documentation/walletpasses/pass/boardingpass
type BoardingPass struct {
	*GenericPass
	TransitType TransitType `json:"transitType,omitempty"`
}

func NewBoardingPass(transitType TransitType) *BoardingPass {
	return &BoardingPass{GenericPass: NewGenericPass(), TransitType: transitType}
}

func (b *BoardingPass) IsValid() bool {
	return len(b.GetValidationErrors()) == 0
}

func (b *BoardingPass) GetValidationErrors() []string {
	var validationErrors []string

	validationErrors = append(validationErrors, b.GenericPass.GetValidationErrors()...)
	if string(b.TransitType) == "" {
		validationErrors = append(validationErrors, "BoardingPass: TransitType is not set")
	}

	return validationErrors
}

// Coupon Representation of https://developer.apple.com/documentation/walletpasses/pass/coupon
type Coupon struct {
	*GenericPass
}

func NewCoupon() *Coupon {
	return &Coupon{GenericPass: NewGenericPass()}
}

// EventTicket Representation of https://developer.apple.com/documentation/walletpasses/pass/eventticket
type EventTicket struct {
	*GenericPass
}

func NewEventTicket() *EventTicket {
	return &EventTicket{GenericPass: NewGenericPass()}
}

// StoreCard Representation of https://developer.apple.com/documentation/walletpasses/pass/storecard
type StoreCard struct {
	*GenericPass
}

func NewStoreCard() *StoreCard {
	return &StoreCard{GenericPass: NewGenericPass()}
}

// Field Representation of https://developer.apple.com/documentation/walletpasses/passfieldcontent
type Field struct {
	Key               string             `json:"key,omitempty"`
	Label             string             `json:"label,omitempty"`
	Value             interface{}        `json:"value,omitempty"`
	AttributedValue   interface{}        `json:"attributedValue,omitempty"`
	ChangeMessage     string             `json:"changeMessage,omitempty"`
	TextAlignment     TextAlignment      `json:"textAlignment,omitempty"`
	DataDetectorTypes []DataDetectorType `json:"dataDetectorTypes,omitempty"`
	CurrencyCode      string             `json:"currencyCode,omitempty"`
	NumberStyle       NumberStyle        `json:"numberStyle,omitempty"`
	DateStyle         DateStyle          `json:"dateStyle,omitempty"`
	TimeStyle         DateStyle          `json:"timeStyle,omitempty"`
	IsRelative        bool               `json:"isRelative,omitempty"`
	IgnoreTimeZone    bool               `json:"ignoresTimeZone,omitempty"`
	Semantics         *SemanticTag       `json:"semantics,omitempty"`
	Row               int                `json:"row,omitempty"`
}

func (f *Field) IsValid() bool {
	return len(f.GetValidationErrors()) == 0
}

func (f *Field) GetValidationErrors() []string {
	var validationErrors []string

	if f.Value == nil || f.Key == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("Field: Not all required Fields are set. Key: %v Value: %v", f.Key, f.Value))
	}

	if f.Value != nil {
		switch f.Value.(type) {
		case string:
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
		case float32:
		case float64:
		case time.Time:
		default:
			validationErrors = append(validationErrors, "Field: Invalid value type. Allowed: string, int, float, time.Time")
		}
	}

	if strings.TrimSpace(f.CurrencyCode) != "" && string(f.NumberStyle) != "" {
		validationErrors = append(validationErrors, "Field: CurrencyCode and numberStyle are both set")
	}

	if (strings.TrimSpace(f.CurrencyCode) != "" || string(f.NumberStyle) != "") && (string(f.DateStyle) != "" || string(f.TimeStyle) != "") {
		validationErrors = append(validationErrors, "Field: Can't be number/currency and date at the same time")
	}

	if strings.TrimSpace(f.ChangeMessage) != "" && !strings.Contains(f.ChangeMessage, "%@") {
		validationErrors = append(validationErrors, "Field: ChangeMessage needs to contain %@ placeholder")
	}

	if strings.TrimSpace(f.CurrencyCode) != "" {
		switch f.Value.(type) {
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
		case float32:
		case float64:
		default:
			validationErrors = append(validationErrors, "Field: When using currencies, the values have to be numbers")
		}
	}

	if f.Semantics != nil && !f.Semantics.IsValid() {
		validationErrors = append(validationErrors, f.Semantics.GetValidationErrors()...)
	}

	if f.Row != 0 && f.Row != 1 {
		validationErrors = append(validationErrors, "Row must be 0 or 1")
	}

	return validationErrors
}

// Beacon Representation of https://developer.apple.com/documentation/walletpasses/pass/beacons
type Beacon struct {
	Major         int    `json:"major,omitempty"`
	Minor         int    `json:"minor,omitempty"`
	ProximityUUID string `json:"proximityUUID,omitempty"`
	RelevantText  string `json:"relevantText,omitempty"`
}

func (b *Beacon) IsValid() bool {
	return len(b.GetValidationErrors()) == 0
}

func (b *Beacon) GetValidationErrors() []string {
	var validationErrors []string

	if strings.TrimSpace(b.ProximityUUID) == "" {
		validationErrors = append(validationErrors, "Beacon: Not all required Fields are set: proximityUUID")
	}

	return validationErrors
}

// Location Representation of https://developer.apple.com/documentation/walletpasses/pass/locations
type Location struct {
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	Altitude     float64 `json:"altitude,omitempty"`
	RelevantText string  `json:"relevantText,omitempty"`
}

func (l *Location) IsValid() bool {
	return len(l.GetValidationErrors()) == 0
}

func (l *Location) GetValidationErrors() []string {
	return []string{}
}

// Barcode Representation of https://developer.apple.com/documentation/walletpasses/pass/barcodes
type Barcode struct {
	Format          BarcodeFormat `json:"format,omitempty"`
	AltText         string        `json:"altText,omitempty"`
	Message         string        `json:"message,omitempty"`
	MessageEncoding string        `json:"messageEncoding,omitempty"`
}

func (b *Barcode) IsValid() bool {
	return len(b.GetValidationErrors()) == 0
}

func (b *Barcode) GetValidationErrors() []string {
	var validationErrors []string

	if string(b.Format) == "" || strings.TrimSpace(b.Message) == "" || strings.TrimSpace(b.MessageEncoding) == "" {
		validationErrors = append(validationErrors, fmt.Sprintf("Barcode: Not all required Fields are set. Format: %v, Message: %v, MessageEncoding: %v, AltText: %v", b.Format, b.Message, b.MessageEncoding, b.AltText))
	}

	return validationErrors
}

type PWAssociatedApp struct {
	Title        string
	IdGooglePlay string
	IdAmazon     string
}

func (a *PWAssociatedApp) IsValid() bool {
	return len(a.GetValidationErrors()) == 0
}

func (a *PWAssociatedApp) GetValidationErrors() []string {
	return []string{}
}

// NFC Representation of https://developer.apple.com/documentation/walletpasses/pass/nfc
type NFC struct {
	Message                string `json:"message,omitempty"`
	EncryptionPublicKey    string `json:"encryptionPublicKey,omitempty"`
	RequiresAuthentication bool   `json:"requiresAuthentication,omitempty"`
}

// Personalization Representation of https://developer.apple.com/documentation/walletpasses/personalize
type Personalization struct {
	RequiredPersonalizationFields []PassPersonalizationField `json:"requiredPersonalizationFields"`
	Description                   string                     `json:"description"`
	TermsAndConditions            string                     `json:"termsAndConditions"`
}

func (pz *Personalization) toJSON() ([]byte, error) {
	return json.Marshal(pz)
}

func (pz *Personalization) IsValid() bool {
	return len(pz.GetValidationErrors()) == 0
}

func (pz *Personalization) GetValidationErrors() []string {
	var validationErrors []string

	if len(pz.RequiredPersonalizationFields) == 0 {
		validationErrors = append(validationErrors, "Personalization: You need to provide at least one requiredPersonalizationField")
	}

	if strings.TrimSpace(pz.Description) == "" {
		validationErrors = append(validationErrors, "Personalization: You need to provide a description")
	}

	return validationErrors
}

type PassRelevantDate struct {
	StartDate *time.Time
	EndDate   *time.Time
}

func (prd *PassRelevantDate) toJSON() ([]byte, error) {
	if prd.EndDate != nil {
		return json.Marshal(&struct {
			StartDate *time.Time `json:"startDate"`
			EndDate   *time.Time `json:"endDate"`
		}{
			StartDate: prd.StartDate,
			EndDate:   prd.EndDate,
		})
	} else {
		return json.Marshal(&struct {
			RelevantDate *time.Time `json:"relevantDate"`
		}{
			RelevantDate: prd.StartDate,
		})
	}
}

func (prd *PassRelevantDate) IsValid() bool {
	return len(prd.GetValidationErrors()) == 0
}

func (prd *PassRelevantDate) GetValidationErrors() []string {
	var validationErrors []string

	if prd.StartDate == nil {
		validationErrors = append(validationErrors, "PassRelevantDate: Not all required Fields are set: startDate")
	}

	return validationErrors
}

func (p *Pass) SetRelevantDates(d []PassRelevantDate) {

	if len(d) == 0 {
		return
	}
	minDate := d[0].StartDate

	for _, pdr := range d {
		if pdr.StartDate.Before(*minDate) {
			minDate = pdr.StartDate
		}
	}

	p.RelevantDates = d
	p.RelevantDate = minDate
}
