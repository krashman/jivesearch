package instant

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jivesearch/jivesearch/instant/contributors"
	"golang.org/x/text/language"
)

// Temperature is an instant answer
type Temperature struct {
	Answer
}

func (t *Temperature) setQuery(r *http.Request, qv string) answerer {
	t.Answer.setQuery(r, qv)
	return t
}

func (t *Temperature) setUserAgent(r *http.Request) answerer {
	return t
}

func (t *Temperature) setLanguage(lang language.Tag) answerer {
	t.language = lang
	return t
}

func (t *Temperature) setType() answerer {
	t.Type = "temperature"
	return t
}

func (t *Temperature) setContributors() answerer {
	t.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return t
}

func (t *Temperature) setRegex() answerer {
	triggers := []string{
		"celsius to fahrenheit", "fahrenheit to celsius", "c to f", "f to c",
	}

	tr := strings.Join(triggers, "|")
	t.regex = append(t.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<trigger>%s) (?P<remainder>.*)$`, tr)))
	t.regex = append(t.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<remainder>.*) (?P<trigger>%s)$`, tr)))

	return t
}

func (t *Temperature) solve() answerer {
	matches := make(map[string]float64)
	combos := [][]string{
		{"<f>fahrenheit", "<c>celsius"},
		{"<c>celsius", "<f>fahrenheit"},
		{"<c>c", "<f>f"},
		{"<f>f", "<c>c"},
	}

	for _, c := range combos {
		// this seems expensive to compile regexp on each loop...better way???
		re := regexp.MustCompile(fmt.Sprintf(`(?P<temp>-?\d+(\.\d+)?).*?(?P%v).*?(?P%v)`, c[0], c[1]))
		match := re.FindStringSubmatch(t.query)

		if len(match) > 0 {
			for i, name := range re.SubexpNames() {
				if i == 0 {
					continue
				}
				if name == "temp" {
					f, _ := strconv.ParseFloat(match[i], 64)
					matches[name] = f
				} else {
					matches[name] = float64(i)
				}
			}

			var converted float64
			var text string
			if matches["f"] < matches["c"] { // fahrenheit to celsius
				converted = (matches["temp"] - 32) * 5 / 9
				text = "%.1f degrees Fahrenheit is %s degrees Celsius"
			} else { // celsius to fahrenheit
				converted = (matches["temp"] * 9 / 5) + 32
				text = "%.1f degrees Celsius is %s degrees Fahrenheit"
			}

			t.Solution = fmt.Sprintf("%.1f", converted)
			t.Solution = fmt.Sprintf(text, matches["temp"], t.Solution)
			break
		}
	}

	return t
}

func (t *Temperature) setCache() answerer {
	t.Cache = true
	return t
}

func (t *Temperature) tests() []test {
	typ := "temperature"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		{
			query: "17 degrees c to f",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "17.0 degrees Celsius is 62.6 degrees Fahrenheit",
					Cache:        true,
				},
			},
		},
		{
			query: "79.9 f to c",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "79.9 degrees Fahrenheit is 26.6 degrees Celsius",
					Cache:        true,
				},
			},
		},
		{
			query: "107.9 fahrenheit to celsius",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "107.9 degrees Fahrenheit is 42.2 degrees Celsius",
					Cache:        true,
				},
			},
		},
		{
			query: "-9.3 celsius to fahrenheit",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "-9.3 degrees Celsius is 15.3 degrees Fahrenheit",
					Cache:        true,
				},
			},
		},
	}

	return tests
}
