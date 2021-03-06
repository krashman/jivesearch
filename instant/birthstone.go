package instant

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jivesearch/jivesearch/instant/contributors"
	"golang.org/x/text/language"
)

// BirthStone is an instant answer
type BirthStone struct {
	Answer
}

func (b *BirthStone) setQuery(r *http.Request, qv string) answerer {
	b.Answer.setQuery(r, qv)
	return b
}

func (b *BirthStone) setUserAgent(r *http.Request) answerer {
	return b
}

func (b *BirthStone) setLanguage(lang language.Tag) answerer {
	b.language = lang
	return b
}

func (b *BirthStone) setType() answerer {
	b.Type = "birthstone"
	return b
}

func (b *BirthStone) setContributors() answerer {
	b.Contributors = contributors.Load(
		[]string{
			"brentadamson",
		},
	)
	return b
}

func (b *BirthStone) setRegex() answerer {
	triggers := []string{
		"birthstones",
		"birth stones",
		"birthstone",
		"birth stone",
	}

	t := strings.Join(triggers, "|")
	b.regex = append(b.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<trigger>%s) (?P<remainder>.*)$`, t)))
	b.regex = append(b.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<remainder>.*) (?P<trigger>%s)$`, t)))

	return b
}

func (b *BirthStone) solve() answerer {
	switch b.remainder {
	case "january":
		b.Solution = "Garnet"
	case "february":
		b.Solution = "Amethyst"
	case "march":
		b.Solution = "Aquamarine, Bloodstone"
	case "april":
		b.Solution = "Diamond"
	case "may":
		b.Solution = "Emerald"
	case "june":
		b.Solution = "Pearl, Moonstone, Alexandrite"
	case "july":
		b.Solution = "Ruby"
	case "august":
		b.Solution = "Peridot, Spinel"
	case "september":
		b.Solution = "Sapphire"
	case "october":
		b.Solution = "Opal, Tourmaline"
	case "november":
		b.Solution = "Topaz, Citrine"
	case "december":
		b.Solution = "Turquoise, Zircon, Tanzanite"
	}

	return b
}

func (b *BirthStone) setCache() answerer {
	b.Cache = true
	return b
}

func (b *BirthStone) tests() []test {
	typ := "birthstone"

	contrib := contributors.Load([]string{"brentadamson"})

	tests := []test{
		{
			query: "January birthstone",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Garnet",
					Cache:        true,
				},
			},
		},
		{
			query: "birthstone february",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Amethyst",
					Cache:        true,
				},
			},
		},
		{
			query: "march birth stone",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Aquamarine, Bloodstone",
					Cache:        true,
				},
			},
		},
		{
			query: "birth stone April",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Diamond",
					Cache:        true,
				},
			},
		},
		{
			query: "birth stones may",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Emerald",
					Cache:        true,
				},
			},
		},
		{
			query: "birthstones June",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Pearl, Moonstone, Alexandrite",
					Cache:        true,
				},
			},
		},
		{
			query: "July Birth Stones",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Ruby",
					Cache:        true,
				},
			},
		},
		{
			query: "birthstones August",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Peridot, Spinel",
					Cache:        true,
				},
			},
		},
		{
			query: "september birthstones",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Sapphire",
					Cache:        true,
				},
			},
		},
		{
			query: "October birthstone",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Opal, Tourmaline",
					Cache:        true,
				},
			},
		},
		{
			query: "birthstone November",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Topaz, Citrine",
					Cache:        true,
				},
			},
		},
		{
			query: "December birthstone",
			expected: []Data{
				{
					Type:         typ,
					Triggered:    true,
					Contributors: contrib,
					Solution:     "Turquoise, Zircon, Tanzanite",
					Cache:        true,
				},
			},
		},
	}

	return tests
}
