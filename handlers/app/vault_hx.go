package app

import (
	// "maxwarden/.jet/table"

	"maxwarden/security"
	. "maxwarden/ui"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"maxwarden/database"
	"maxwarden/entries"

	"maxwarden/middleware"

	"net/http"
)

func VaultHxHandler(w http.ResponseWriter, r *http.Request) {
	identity := middleware.GetIdentity(r)

	filter := database.ParseFilterFromRequest(r)
	filter.Pagination.Enabled = true

	entryFilter := entries.EntryFilter{
		Filter:    filter,
		UserId:    identity.UserID,
		MasterKey: identity.MasterKey,
	}

	// fetch entities from filter function
	// this first counts the possible items before pagination
	searchFilter := entries.EntryFilter{
		Filter:    database.NewFilterFromSearch(filter.Search),
		UserId:    identity.UserID,
		MasterKey: identity.MasterKey,
	}

	searchItems, _ := entries.Filter(searchFilter)

	// this query gets the data AFTER pagination
	entryList, _ := entries.Filter(entryFilter)

	// generate page numbers according to total length of data
	entryFilter.Filter.Pagination.GeneratePagination(len(searchItems), len(entryList))

	// Col header names and referenced database col names
	cols := []database.ColInfo{
		{DbName: "Description", DisplayName: "Description", Sortable: true},
		{DisplayName: "Username"},
		{DisplayName: "Password"},
		{DisplayName: "URL"},
		{DisplayName: "Action"},
	}

	// Generate HTML
	elId := "order_table"
	AutoTable(
		elId,
		r.URL.Path,
		cols,
		entryFilter.Filter,
		entryList,
		AutotableSearchGroup(
			AutotableSearch(
				Placeholder("Search Description..."),
				BindSearch(elId, "description"),
				AutoFocus(),
			),
		),
		func(entry entries.Secret) Node {
			return Tr(
				TdLeft(
					Div(
						InlineStyle("$me { overflow-x: auto; }"),
						Text(entry.Description),
					),
				),
				TdLeft(
					IfElse(entry.Username != "",
						Div(
							InlineStyle("$me { overflow-x: auto; }"),
							Flex(
								Span(
									Title("Click to copy username"),
									InlineStyle("$me { cursor: pointer; }"),
									Icon(ICON_COPY, 16),
								),
								P(Text(entry.Username)),
								InlineScript(`
								let btn = me("span", me());
								let text = me("p", me());

								btn.on("click", () => {
									navigator.clipboard.writeText(text.innerHTML);
								});
							`),
							),
						),
						Text("---"),
					),
				),
				TdLeft(
					IfElse(entry.Password != "",
						Flex(
							Span(
								Class("copy"),
								Title("Click to copy password"),
								InlineStyle("$me { cursor: pointer; }"),
								Icon(ICON_COPY, 16),
							),
							P(Text("•••••••")),
							Input(Type("hidden"), Value(entry.Password)),
							InlineScript(`
								let copyBtn = me(".copy", me());
								let password = me("input", me()).value;


								copyBtn.on("click", () => {
									navigator.clipboard.writeText(password);
								});
							`),
						),
						Text("---"),
					),
				),
				TdLeft(
					IfElse(entry.URL != "",
						Div(
							InlineStyle("$me { overflow-x: auto; }"),
							PageLink(security.SanitizationPolicy.Sanitize(entry.URL), Text(entry.URL), true),
						),
						Text("---"),
					),
				),
				TdCenter(A(Href("/app/editor/edit/"+entry.ID), ButtonUIOutline(Icon(ICON_PENCIL, 16)))),
			)
		},
		nil,
		AutoTableOptions{
			Compact:   false,
			Shadow:    true,
			Hover:     false,
			Alternate: false,
			BorderX:   true,
			BorderY:   false,
		},
	).Render(w)
}
