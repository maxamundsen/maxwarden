package ui

import (
	. "maxwarden/basic"

	"maxwarden/database"

	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

const TABLE_ABOVE_PREFIX = "_tableAbove"
const FORM_BIND_SUFFIX = "_form"
const FORM_PAGINATION_SUFFIX = "_paginationForm"

// FILTER COMPONENTS

// Autotable bind search to col
func BindSearch(elId string, identifier string) Node {
	return Group{
		FormAttr(elId + FORM_BIND_SUFFIX),
		Name(database.SEARCH_URL_KEY_PREFIX + identifier),
	}
}

const (
	AUTOTABLE_HEADER_COLOR_DEFAULT = iota
)

type AutoTableOptions struct {
	Compact       bool
	Shadow        bool
	Hover         bool // highlight rows when hovering over them
	Alternate     bool // highlight alternating rows
	HeaderBorderY bool
	BorderX       bool
	BorderY       bool
	Color         int // "enum"
}

// THE TABLE
// Note that "aboveTable" node is not swapped with HTMX, but "belowTable" is.
func AutoTable[E any](tableId string, url string, cols []database.ColInfo, f database.Filter, entities []E, aboveTable Node, rowComponent func(E) Node, belowTable Node, opts AutoTableOptions) Node {
	paginationButton := func(icon string, page int) Node {
		return Button(
			InlineStyle(`
				$me {
					padding-left: $3;
					padding-right: $3;
					padding-top: $1;
					padding-bottom: $1;
					min-height: $9;
					font-size: var(--text-sm);
					font-weight: var(--font-weight-normal);
					color: $color(neutral-800);
					cursor: pointer;
				}
			`),
			Icon(icon, 16),
			hx.Get(url+database.QueryParamsFromPagenum(page, f)),
			hx.Swap(CSSID(tableId)),
			hx.Target(CSSID(tableId)),
			hx.Select(CSSID(tableId)),
			hx.Trigger("click"),
		)
	}

	return Group{
		If((aboveTable != nil) && (tableId != ""),
			Div(ID(tableId+TABLE_ABOVE_PREFIX),
				aboveTable,
			),
		),
		Div(
			If(tableId != "",
				Group{
					ID(tableId),
					Form(ID(tableId+FORM_BIND_SUFFIX),
						AutoComplete("off"),
						hx.Get(url),
						hx.Trigger("keyup delay:100ms from:(#"+tableId+TABLE_ABOVE_PREFIX+" input), change from:(#"+tableId+TABLE_ABOVE_PREFIX+" input[type=date]), change from:(#"+tableId+TABLE_ABOVE_PREFIX+" input[type=datetime-local]), change from:(#"+tableId+TABLE_ABOVE_PREFIX+" select)"),
						hx.Swap("outerHTML"),
						hx.Target(CSSID(tableId)),
						hx.Select(CSSID(tableId)),
						Input(Type("hidden"), Name(database.ORDER_BY_URL_KEY), Value(f.OrderBy)),
						Input(Type("hidden"), Name(database.ORDER_DESC_URL_KEY), Value(ToString(f.OrderDescending))),
						Input(Type("hidden"), Name(database.ITEMS_PER_PAGE_URL_KEY), Value(ToString(f.Pagination.MaxItemsPerPage))),
					),
				},
			),
			Div(
				If(opts.Shadow,
					InlineStyle(`
						$me {
							background-color: $color(white);
							box-shadow: var(--shadow-md);
						}
					`),
				),
				InlineStyle(`
					$me {
						position: relative;
						display: flex;
						flex-direction: column;
						width: 100%;
						height: 100%;
						color: $color(gray-700);
						background-color: $color(white);
					}
				`),
				Div(
					InlineStyle(`
						$me {
							display: flex;
							flex-direction: column;
							width: 100%;
							height: 100%;
						}
					`),
					Table(
						InlineStyle("$me { table-layout: fixed; width: 100%; }"),
						// ---- TABLE HEADER ----
						THead(
							Tr(
								Map(cols, func(col database.ColInfo) Node {
									return Th(
										If(col.DisplayPosition == database.COL_POS_RIGHT,
											InlineStyle("$me { text-align: right; }"),
										),
										InlineStyle(`
											$me {
												border-bottom: 1px solid $color(neutral-200);
											}

										`),
										If(opts.HeaderBorderY,
											InlineStyle(`
												$me:not(:first-child) {
													border-left: 1px solid $color(neutral-200);
												}
											`),
										),
										IfElse(opts.Compact,
											InlineStyle("$me { padding: $2 $3; }"),
											InlineStyle("$me { padding: $4 $4; }"),
										),
										If(col.Sortable,
											Group{
												hx.Get(url + database.QueryParamsFromOrderBy(col.DbName, !f.OrderDescending && (col.DbName == f.OrderBy), f)),
												hx.Swap(CSSID(tableId)),
												hx.Target(CSSID(tableId)),
												hx.Select(CSSID(tableId)),
												hx.Trigger("click"),
												InlineStyle(`
													$me {
														cursor: pointer;
													}

													$me:hover {
														background-color: $color(neutral-200);
													}
												`),
											},
										),
										P(
											InlineStyle(`
												$me {
													display: flex;
													justify-content: space-between;
													align-items: center;
													gap: $2;
													font-size: var(--text-sm);
													color: $color(neutral-800);
													font-weight: var(--font-weight-bold);
												}
											`),
											If(col.DisplayPosition == database.COL_POS_LEFT,
												InlineStyle("$me { flex-direction: row; }"),
											),
											If(col.DisplayPosition == database.COL_POS_RIGHT,
												InlineStyle("$me { flex-direction: row-reverse; }"),
											),
											Text(col.DisplayName),
											If((f.OrderBy == col.DbName) && (col.DbName != ""),
												Group{
													InlineStyle("$me { color: $color(black); text-decoration: underline; }"),
													IfElse(f.OrderDescending,
														Icon(ICON_ARROW_DOWN, 16),
														Icon(ICON_ARROW_UP, 16),
													),
												},
											),
										),
									)
								}),
							),
						),
						// ---- TABLE BODY ----
						TBody(
							InlineStyle(`
								$me > tr > td {
									overflow-x: auto;
									font-size: var(--text-sm);
									color: $color(neutral-700);
								}
							`),
							If(opts.Alternate,
								InlineStyle("$me > tr:nth-child(even) { background: $color(neutral-100) }"),
							),
							If(opts.Hover && len(entities) > 0,
								InlineStyle("$me > tr:hover {background-color: $color(neutral-100);}"),
							),
							If(opts.BorderX && f.Pagination.Enabled,
								InlineStyle("$me > tr {border-bottom: 1px solid $color(neutral-200); }"),
							),
							If(opts.BorderX && !f.Pagination.Enabled,
								InlineStyle("$me > tr:not(:last-child) {border-bottom: 1px solid $color(neutral-200); }"),
							),
							If(opts.BorderY,
								InlineStyle(`
									$me > tr > td:not(:first-child) {
										border-left: 1px solid $color(neutral-200);
									}
								`),
							),
							IfElse(opts.Compact,
								InlineStyle("$me > tr > td { padding: $1 $3; }"),
								InlineStyle("$me > tr > td { padding: $4; }"),
							),
							IfElse(len(entities) > 0,
								Map(entities, func(e E) Node {
									return rowComponent(e)
								}),
								Tr(
									Td(
										Text("Dataset contains no entries."),
									),
								),
							),
						),
					),
				),
				// ---- PAGINATION ----
				If(f.Pagination.Enabled,
					Div(InlineStyle("$me { display: flex; justify-content: space-between; align-items: center; }"),
						IfElse(opts.Compact,
							InlineStyle("$me { padding: $1 $4; }"),
							InlineStyle("$me { padding: $3 $4; }"),
						),
						Div(
							InlineStyle("$me { display: flex; align-items: center; font-size: var(--text-sm); color: $color(neutral-500); }"),
							B(Icon(ICON_LIST_ORDERED, 16)),
							Span(InlineStyle("$me {margin-right: $3;}")), ToText(f.Pagination.ViewRangeLower), Text("-"), ToText(f.Pagination.ViewRangeUpper), Text(" of "), ToText(f.Pagination.TotalItems),
						),

						Div(InlineStyle("$me { display: flex; align-items: center; }"),
							Div(InlineStyle("$me { margin-right: $3; font-size: var(--text-sm); color: $color(neutral-500);}"),
								Span(Text("Items per page:")),
							),
							Form(
								ID(tableId+FORM_PAGINATION_SUFFIX),
								hx.Get(url),
								hx.Trigger("change from:(#"+tableId+FORM_PAGINATION_SUFFIX+" select)"),
								hx.Target(CSSID(tableId)),
								hx.Select(CSSID(tableId)),
								hx.Swap("outerHTML"),

								MapMapWithKey(f.Search, func(s string, v string) Node {
									return Input(Type("hidden"), Name(database.SEARCH_URL_KEY_PREFIX+s), Value(v))
								}),

								Input(Type("hidden"), Name(database.ORDER_BY_URL_KEY), Value(f.OrderBy)),
								Input(Type("hidden"), Name(database.ORDER_DESC_URL_KEY), Value(ToString(f.OrderDescending))),
								Select(
									IfElse(opts.Compact,
										InlineStyle("$me { padding: $2; }"),
										InlineStyle("$me { padding: $3; }"),
									),
									InlineStyle(`
										$me {
											background-color: $color(gray-50);
											border: 1px solid $color(gray-50);
											font-size: var(--text-sm);
											display: block;
											margin-right: $3;
											box-shadow: var(--shadow-sm);
										}
									`),
									Name(database.ITEMS_PER_PAGE_URL_KEY),
									Option(If(f.Pagination.MaxItemsPerPage == 5, Selected()), Text("5"), Value("5")),
									Option(If(f.Pagination.MaxItemsPerPage == 10, Selected()), Text("10"), Value("10")),
									Option(If(f.Pagination.MaxItemsPerPage == 25, Selected()), Text("25"), Value("25")),
									Option(If(f.Pagination.MaxItemsPerPage == 50, Selected()), Text("50"), Value("50")),
									Option(If(f.Pagination.MaxItemsPerPage == 100, Selected()), Text("100"), Value("100")),
								),
							),

							paginationButton(ICON_CHEVRON_FIRST, 1),
							paginationButton(ICON_CHEVRON_LEFT, f.Pagination.PreviousPage),

							Div(
								InlineStyle("$me { font-size: var(--text-sm); display: flex; justify-content: center; color: $color(neutral-500); padding: $3; min-height: $9; }"),
								Text("Page "),
								ToText(f.Pagination.CurrentPage),
								Text(" of "),
								ToText(f.Pagination.TotalPages),
							),

							paginationButton(ICON_CHEVRON_RIGHT, f.Pagination.NextPage),
							paginationButton(ICON_CHEVRON_LAST, f.Pagination.TotalPages),
						),
					),
				),
			),
			belowTable,
		),
	}
}

func AutotableSearchGroup(children ...Node) Node {
	return Div(InlineStyle("$me { width: 100%; display: flex; justify-content: space-between; margin-bottom: $3; margin-top: $1; }"),
		Div(InlineStyle("$me { width: 100%; position: relative; }"),
			Div(InlineStyle("$me { position: relative; display: flex; flex-direction: row; align-items: center; gap: $1;}"),
				Group(children),
			),
		),
	)
}

func AutotableSearchDropdown(c ...Node) Node {
	return Div(
		InlineStyle(`
			$me {
				width: 100%;
			}
		`),
		FormSelect(Group(c)),
	)
}

func AutotableSearch(c ...Node) Node {
	return Div(
		InlineStyle(`
			$me {
				width: 100%;
			}
		`),
		FormInput(Group(c)),
	)
}

// "premade" autotable for simple tables that don't feature dynamic filtering.
// We can still use the styling and general form of the fancy autotable defined above, but
// for simple datasets because why not. This also gives you the option to "upgrade"
// to the "full" table later on, since you are using the same api
func AutoTableLite[E any](columnNames []string, entities []E, rowComponent func(E) Node, opts AutoTableOptions) Node {
	cols := []database.ColInfo{}

	for _, v := range columnNames {
		cols = append(cols, database.ColInfo{DisplayName: v})
	}

	return AutoTable(
		"",
		"",
		cols,
		database.Filter{},
		entities,
		nil,
		rowComponent,
		nil,
		opts,
	)
}
