package icsreader

// Constants for useful ICS keywords
const (
	eventPropertyBeginV       = "BEGIN"
	eventPropertyEndV         = "END"
	eventPropertyStart        = "DTSTART"
	eventPropertyEnd          = "DTEND"
	eventPropertySummary      = "SUMMARY"
	eventPropertyLocation     = "LOCATION"
	eventPropertyDescription  = "DESCRIPTION"
	eventPropertyLastModified = "LAST-MODIFIED"
)

// Date : basic date struct
type Date struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// Event : ICS VENVENT struct
type Event struct {
	Start        Date   `json:"start"`
	End          Date   `json:"end"`
	Summary      string `json:"summary"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	LastModified Date   `json:"lastModified"`
}

// Events : shortcut to a list of event
type Events []Event

// CalendarError : error for handling calendar not found in file or from source (URL)
type CalendarError struct{}

func (e *CalendarError) Error() string {
	return "Cannot load Calendar"
}
