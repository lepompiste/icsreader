package icsreader

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func isIgnoredKeyword(kw string) bool { // test if keyword is ignored or not
	switch kw {
	case "METHOD", "PRODID", "VERSION", "CALSCALE", "UID", "CREATED", "SEQUENCE", "DTSTAMP":
		return true
	}
	return false
}

// Parse : transform a reader into a list of events
func Parse(r io.Reader) Events {
	var reader *bufio.Reader = bufio.NewReader(r) // Create a bufio.Reader to read the stream line by line

	var idle bool = false   // used to ignore several types of entries (VALARM, VTODO, ...)
	var idleLevel uint8 = 0 // used to get the level of indentation inside ignored entry

	var lastKw string // last ICS keyword, use to deal with multiple lines text

	var res Events // return table

	var e *Event // ongoing event struct

	for {
		line, _, err := reader.ReadLine() // Reading line by line
		if err == io.EOF {
			break // breaking when end of stream
		}

		parts := strings.Split(string(line), ":") // splitting at ':' to get the fisrt keyword (BEGIN, END, ...)

		if !idle || parts[0] == "END" { // if current event is not ignored, continue. if the keyword is "END", then reducing idleLevel and maybe set idle to false

			switch parts[0] { // swicth the keyword
			case eventPropertyBeginV:
				if parts[1] == "VEVENT" { // only VEVENT is insteresting us
					e = &Event{}
				} else if parts[1] != "VCALENDAR" { // VCALENDAR is ignored, but not setting idle to true, because all the data is in it. if not VCALENDAR, then starting idle until the end of the entry
					idle = true
					idleLevel++
					//Debugging
					//println("Idling")
				}

			case eventPropertyStart:
				e.Start = makeDate(parts[1], 1)

			case eventPropertyEnd:
				e.End = makeDate(parts[1], 1)

			case eventPropertySummary:
				lastKw = eventPropertySummary
				e.Summary = strings.Join(parts[1:], ":")

			case eventPropertyLocation:
				lastKw = eventPropertyLocation
				e.Location = strings.Join(parts[1:], ":")

			case eventPropertyDescription:
				lastKw = eventPropertyDescription
				e.Description = strings.Join(parts[1:], ":")

			case eventPropertyLastModified:
				e.LastModified = makeDate(parts[1], 1)

			case eventPropertyEndV:
				if parts[1] == "VEVENT" { // at the end of VEVENT, adding the event struct to the res (Events) table
					res = append(res, *e)
				} else if parts[1] != "VCALENDAR" { // same as BEGIN, ignoring if VCALENDAR, and reducing idleLevel at the end of ignored entry
					idleLevel--

					if idleLevel == 0 {
						idle = false
						//Debugging
						//println("Desidling")
					}
				}
			default:
				if !isIgnoredKeyword(parts[0]) { // many of keywords are ignored, even inside of a VEVENT (SEQUENCE, UID, ...)
					//Debugging
					//println("suite ligne :", string(line))
					switch lastKw { // if a line not strating with a keyword, then the text of this line is referred to the previous line, so we add it to the previous field
					case eventPropertySummary:
						e.Summary = e.Summary + string(line)
					case eventPropertyLocation:
						e.Location = e.Location + string(line)
					case eventPropertyDescription:
						e.Description = e.Description + string(line)
					}
				}
			}
		} /* else {
			println("Ignored line")
		}*/
		// Debugging comment
	}
	return res
}

// make a Date struct from the ICS date format
func makeDate(d string, correction int) (res Date) {
	//Format is YYYYMMDD T HHMMSS Z given at UTC
	year, err1 := strconv.Atoi(d[0:4])
	month, err2 := strconv.Atoi(d[4:6])
	day, err3 := strconv.Atoi(d[6:8])
	hour, err4 := strconv.Atoi(d[9:11])
	min, err5 := strconv.Atoi(d[11:13])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return
	}

	res.Year = year
	res.Month = month
	res.Day = day
	res.Hour = hour + correction
	res.Minute = min
	return
}
