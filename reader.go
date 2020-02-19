package icsreader

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func isIgnoredKeyword(kw string) bool {
	switch kw {
	case "METHOD", "PRODID", "VERSION", "CALSCALE", "UID", "CREATED", "SEQUENCE", "DTSTAMP":
		return true
	}
	return false
}

// Parse : transform a reader into a list of events
func Parse(r io.Reader) Events {
	var reader *bufio.Reader = bufio.NewReader(r)
	var idle bool = false
	var idleLevel uint8 = 0
	var lastKw string

	var res Events
	var e *Event

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		parts := strings.Split(string(line), ":")

		if !idle || parts[0] == "END" {
			switch parts[0] {
			case eventPropertyBeginV:
				if parts[1] == "VEVENT" {
					e = &Event{}
				} else if parts[1] != "VCALENDAR" {
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
				if parts[1] == "VEVENT" {
					res = append(res, *e)
				} else if parts[1] != "VCALENDAR" {
					idleLevel--
					if idleLevel == 0 {
						idle = false
						//Debugging
						//println("Desidling")
					}
				}
			default:
				if !isIgnoredKeyword(parts[0]) {
					//Debugging
					//println("suite ligne :", string(line))
					switch lastKw {
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

func makeDate(d string, correction int) (res Date) {
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
