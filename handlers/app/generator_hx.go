package app

import (
	"maxwarden/generator"
	. "maxwarden/ui"
	"net/http"
	"strconv"

	. "maxwarden/basic"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	hx "maragu.dev/gomponents-htmx"
)

func GeneratorHxHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	wordCount, _ := strconv.Atoi(r.FormValue("wordCount"))

	if wordCount <= 0 {
		wordCount = 6
	}

	phraseOutput := generator.GeneratePassphrase(wordCount)

	length, _ := strconv.Atoi(r.FormValue("length"))
	if length <= 0 {
		length = 24
	}

	numbers, _ := strconv.Atoi(r.FormValue("numbers"))
	if numbers <= 0 {
		numbers = 5
	}

	symbols, _ := strconv.Atoi(r.FormValue("symbols"))
	if symbols <= 0 {
		symbols = 5
	}

	passwordOutput := generator.GeneratePassword(length, numbers, symbols, false, true)

	customFormInput := func(children ...Node) Node {
		return Div(
			InlineStyle(`
				$me {
					border: 1px solid $color(neutral-300);
					margin-bottom: $3;
				}
			`),
			FormInput(
				Group(children),
			),
		)
	}

	Div(ID("generator"),
		Card(
			Heading("Generate Passphrase"),
			Form(
				hx.Post(r.URL.Path),
				hx.Swap("outerHTML"),
				hx.Target("#generator"),
				FormLabel(Text("Word Count")),
				customFormInput(Type("number"), Name("wordCount"), Value(ToString(wordCount))),

				ButtonUI(Text("Generate"), Type("submit")),
			),

			Br(),

			FormLabel(Text("Output")),
			customFormInput(
				Value(phraseOutput),
			),
		),

		Br(),

		Card(
			Heading("Generate Password"),
			Form(
				hx.Post(r.URL.Path),
				hx.Swap("outerHTML"),
				hx.Target("#generator"),
				Flex(
					Div(
						FormLabel(Text("Password length")),
						customFormInput(Type("number"), Name("length"), Value(ToString(length))),
					),
					Div(
						FormLabel(Text("Number Count")),
						customFormInput(Type("number"), Name("numbers"), Value(ToString(numbers))),
					),
					Div(
						FormLabel(Text("Symbol Count")),
						customFormInput(Type("number"), Name("symbols"), Value(ToString(symbols))),
					),
				),

				ButtonUI(Text("Generate"), Type("submit")),
			),

			Br(),

			FormLabel(Text("Output")),
			customFormInput(
				Value(passwordOutput),
			),
		),
	).Render(w)
}
