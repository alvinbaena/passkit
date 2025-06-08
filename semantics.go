package passkit

import "time"

// SemanticTag Representation of https://developer.apple.com/documentation/walletpasses/semantictags
type SemanticTag struct {
	AdditionalTicketAttributes     string                           `json:"additionalTicketAttributes,omitempty"`
	AdmissionLevel                 string                           `json:"admissionLevel,omitempty"`
	AdmissionLevelAbbreviation     string                           `json:"admissionLevelAbbreviation,omitempty"`
	AlbumIDs                       []string                         `json:"albumIDs,omitempty"`
	AirlineCode                    string                           `json:"airlineCode,omitempty"`
	ArtistIds                      []string                         `json:"artistIDs,omitempty"`
	AttendeeName                   string                           `json:"attendeeName,omitempty"`
	AwayTeamAbbreviation           string                           `json:"awayTeamAbbreviation,omitempty"`
	AwayTeamLocation               string                           `json:"awayTeamLocation,omitempty"`
	AwayTeamName                   string                           `json:"awayTeamName,omitempty"`
	Balance                        *SemanticTagCurrencyAmount       `json:"balance,omitempty"`
	BoardingGroup                  string                           `json:"boardingGroup,omitempty"`
	BoardingSequenceNumber         string                           `json:"boardingSequenceNumber,omitempty"`
	CarNumber                      string                           `json:"carNumber,omitempty"`
	ConfirmationNumber             string                           `json:"confirmationNumber,omitempty"`
	CurrentArrivalDate             *time.Time                       `json:"currentArrivalDate,omitempty"`
	CurrentBoardingDate            *time.Time                       `json:"currentBoardingDate,omitempty"`
	CurrentDepartureDate           *time.Time                       `json:"currentDepartureDate,omitempty"`
	DepartureAirportCode           string                           `json:"departureAirportCode,omitempty"`
	DepartureAirportName           string                           `json:"departureAirportName,omitempty"`
	DepartureGate                  string                           `json:"departureGate,omitempty"`
	DepartureLocation              *SemanticTagLocation             `json:"departureLocation,omitempty"`
	DepartureLocationDescription   string                           `json:"departureLocationDescription,omitempty"`
	DeparturePlatform              string                           `json:"departurePlatform,omitempty"`
	DepartureStationName           string                           `json:"departureStationName,omitempty"`
	DepartureTerminal              string                           `json:"departureTerminal,omitempty"`
	DestinationAirportCode         string                           `json:"destinationAirportCode,omitempty"`
	DestinationAirportName         string                           `json:"destinationAirportName,omitempty"`
	DestinationGate                string                           `json:"destinationGate,omitempty"`
	DestinationLocation            *SemanticTagLocation             `json:"destinationLocation,omitempty"`
	DestinationLocationDescription string                           `json:"destinationLocationDescription,omitempty"`
	DestinationPlatform            string                           `json:"destinationPlatform,omitempty"`
	DestinationStationName         string                           `json:"destinationStationName,omitempty"`
	DestinationTerminal            string                           `json:"destinationTerminal,omitempty"`
	Duration                       *uint64                          `json:"duration,omitempty"`
	EventEndDate                   *time.Time                       `json:"eventEndDate,omitempty"`
	EventLiveMessage               string                           `json:"eventLiveMessage,omitempty"`
	EventName                      string                           `json:"eventName,omitempty"`
	EventStartDate                 *time.Time                       `json:"eventStartDate,omitempty"`
	EventStartDateInfo             *SemanticTagEventDateInfo        `json:"eventStartDateInfo,omitempty"`
	EventType                      EventType                        `json:"eventType,omitempty"`
	EntranceDescription            string                           `json:"entranceDescription,omitempty"`
	FlightCode                     string                           `json:"flightCode,omitempty"`
	FlightNumber                   string                           `json:"flightNumber,omitempty"`
	Genre                          string                           `json:"genre,omitempty"`
	HomeTeamAbbreviation           string                           `json:"homeTeamAbbreviation,omitempty"`
	HomeTeamLocation               string                           `json:"homeTeamLocation,omitempty"`
	HomeTeamName                   string                           `json:"homeTeamName,omitempty"`
	LeagueAbbreviation             string                           `json:"leagueAbbreviation,omitempty"`
	LeagueName                     string                           `json:"leagueName,omitempty"`
	MembershipProgramName          string                           `json:"membershipProgramName,omitempty"`
	MembershipProgramNumber        string                           `json:"membershipProgramNumber,omitempty"`
	OriginalArrivalDate            *time.Time                       `json:"originalArrivalDate,omitempty"`
	OriginalBoardingDate           *time.Time                       `json:"originalBoardingDate,omitempty"`
	OriginalDepartureDate          *time.Time                       `json:"originalDepartureDate,omitempty"`
	PassengerName                  *SemanticTagPersonNameComponents `json:"passengerName,omitempty"`
	PerformerNames                 []string                         `json:"performerNames,omitempty"`
	PlaylistIDs                    []string                         `json:"playlistIDs,omitempty"`
	PriorityStatus                 string                           `json:"priorityStatus,omitempty"`
	Seats                          []SemanticTagSeat                `json:"seats,omitempty"`
	SecurityScreening              string                           `json:"securityScreening,omitempty"`
	SilenceRequested               bool                             `json:"silenceRequested,omitempty"`
	SportName                      string                           `json:"sportName,omitempty"`
	TailgatingAllowed              bool                             `json:"tailgatingAllowed,omitempty"`
	TotalPrice                     *SemanticTagCurrencyAmount       `json:"totalPrice,omitempty"`
	TransitProvider                string                           `json:"transitProvider,omitempty"`
	TransitStatus                  string                           `json:"transitStatus,omitempty"`
	TransitStatusReason            string                           `json:"transitStatusReason,omitempty"`
	VehicleName                    string                           `json:"vehicleName,omitempty"`
	VehicleNumber                  string                           `json:"vehicleNumber,omitempty"`
	VehicleType                    string                           `json:"vehicleType,omitempty"`
	VenueBoxOfficeOpenDate         *time.Time                       `json:"venueBoxOfficeOpenDate,omitempty"`
	VenueCloseDate                 *time.Time                       `json:"venueCloseDate,omitempty"`
	VenueDoorsOpenDate             *time.Time                       `json:"venueDoorsOpenDate,omitempty"`
	VenueEntrance                  string                           `json:"venueEntrance,omitempty"`
	VenueEntranceDoor              string                           `json:"venueEntranceDoor,omitempty"`
	VenueEntranceGate              string                           `json:"venueEntranceGate,omitempty"`
	VenueEntrancePortal            string                           `json:"venueEntrancePortal,omitempty"`
	VenueFanZoneOpenDate           *time.Time                       `json:"venueFanZoneOpenDate,omitempty"`
	VenueGatesOpenDate             *time.Time                       `json:"venueGatesOpenDate,omitempty"`
	VenueLocation                  *SemanticTagLocation             `json:"venueLocation,omitempty"`
	VenueName                      string                           `json:"venueName,omitempty"`
	VenueOpenDate                  *time.Time                       `json:"venueOpenDate,omitempty"`
	VenueParkingLotsOpenDate       *time.Time                       `json:"venueParkingLotsOpenDate,omitempty"`
	VenuePhoneNumber               string                           `json:"venuePhoneNumber,omitempty"`
	VenueRegionName                string                           `json:"venueRegionName,omitempty"`
	VenueRoom                      string                           `json:"venueRoom,omitempty"`
	WifiAccess                     []SemanticTagWifiNetwork         `json:"wifiAccess,omitempty"`
}

func (s *SemanticTag) IsValid() bool {
	return len(s.GetValidationErrors()) == 0
}

func (s *SemanticTag) GetValidationErrors() []string {
	var validationErrors []string
	// Only validate what is validatable
	if s.WifiAccess != nil {
		for _, wifiAccess := range s.WifiAccess {
			if !wifiAccess.IsValid() {
				validationErrors = append(validationErrors, wifiAccess.GetValidationErrors()...)
			}
		}
	}
	return validationErrors
}

// SemanticTagEventDateInfo Representation of https://developer.apple.com/documentation/walletpasses/semantictagtype/eventdateinfo-data.dictionary
type SemanticTagEventDateInfo struct {
	DateDescription string     `json:"dateDescription,omitempty"`
	IsTentative     bool       `json:"isTentative,omitempty"`
	OriginalDate    *time.Time `json:"originalDate,omitempty"`
}

// SemanticTagCurrencyAmount Representation of https://developer.apple.com/documentation/walletpasses/semantictagtype/currencyamount
type SemanticTagCurrencyAmount struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

func (s *SemanticTagCurrencyAmount) IsValid() bool {
	return len(s.GetValidationErrors()) == 0
}

func (s *SemanticTagCurrencyAmount) GetValidationErrors() []string {
	return []string{}
}

// SemanticTagLocation Representation of https://developer.apple.com/documentation/walletpasses/semantictagtype/location
type SemanticTagLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *SemanticTagLocation) IsValid() bool {
	return len(l.GetValidationErrors()) == 0
}

func (l *SemanticTagLocation) GetValidationErrors() []string {
	return []string{}
}

// SemanticTagPersonNameComponents Representation of https://developer.apple.com/documentation/walletpasses/semantictagtype/personnamecomponents
type SemanticTagPersonNameComponents struct {
	FamilyName             string `json:"familyName,omitempty"`
	GivenName              string `json:"givenName,omitempty"`
	MiddleName             string `json:"middleName,omitempty"`
	NamePrefix             string `json:"namePrefix,omitempty"`
	NameSuffix             string `json:"nameSuffix,omitempty"`
	Nickname               string `json:"nickname,omitempty"`
	PhoneticRepresentation string `json:"phoneticRepresentation,omitempty"`
}

func (l *SemanticTagPersonNameComponents) IsValid() bool {
	return len(l.GetValidationErrors()) == 0
}

func (l *SemanticTagPersonNameComponents) GetValidationErrors() []string {
	return []string{}
}

// SemanticTagSeat Representation of https://developer.apple.com/documentation/walletpasses/semantictagtype/seat
type SemanticTagSeat struct {
	SeatDescription string `json:"seatDescription,omitempty"`
	SeatIdentifier  string `json:"seatIdentifier,omitempty"`
	SeatNumber      string `json:"seatNumber,omitempty"`
	SeatRow         string `json:"seatRow,omitempty"`
	SeatSection     string `json:"seatSection,omitempty"`
	SeatType        string `json:"seatType,omitempty"`
}

func (l *SemanticTagSeat) IsValid() bool {
	return len(l.GetValidationErrors()) == 0
}

func (l *SemanticTagSeat) GetValidationErrors() []string {
	return []string{}
}

type SemanticTagWifiNetwork struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}

func (w *SemanticTagWifiNetwork) IsValid() bool {
	return len(w.GetValidationErrors()) == 0
}

func (w *SemanticTagWifiNetwork) GetValidationErrors() []string {
	var validationErrors []string
	// Must have both attributes
	if w.SSID == "" || w.Password == "" {
		validationErrors = append(validationErrors, "SemanticTagWifiNetwork: Both ssid and password must be set")
	}
	return validationErrors
}
