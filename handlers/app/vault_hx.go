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
				TdLeft(Text(entry.Description)),
				TdLeft(Text(entry.Username)),
				TdLeft(Text("********")),
				TdLeft(PageLink(security.SanitizationPolicy.Sanitize(entry.URL), Text(entry.URL), false)),
				TdCenter(A(Href("/app/editor/edit/" + entry.ID), ButtonUIOutline(Icon(ICON_PENCIL, 16)))),
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
