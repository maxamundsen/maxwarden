package ui

import (
	"maxwarden/auth"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const (
	LAYOUT_SECTION_VAULT = iota
	LAYOUT_SECTION_TOOLS  = iota
	LAYOUT_SECTION_ACCOUNT   = iota
	LAYOUT_SECTION_API       = iota
)

type NavGroup struct {
	SectionId int
	Title     string
	URL       string
	SubGroup  []NavGroup
	NewTab    bool
}

var NavGroups = []NavGroup{
	{SectionId: LAYOUT_SECTION_VAULT, Title: "Vault", URL: "/app", SubGroup: nil},
	{
		Title: "Tools",
		SubGroup: []NavGroup{
			{SectionId: LAYOUT_SECTION_TOOLS, Title: "Generator", URL: "/app/generator"},
		},
	},
}

func AppLayout(title string, identity auth.Identity, session map[string]interface{}, children ...Node) Node {
	navbarDropdown := func(dropdownHeader Node, dropdownItems Node) Node {
		return Div(
			InlineStyle("$me{cursor: pointer; position: relative; margin-left: $3;}"),
			Div(
				Class("button"),
				InlineStyle(`
					$me {
						cursor: pointer;
						display: flex;
						position: relative;
						padding-top: $2;
						padding-bottom: $2;
						padding-left: $3;
						font-size: var(--text-sm);
						line-height: var(--text-sm--line-height);
						font-weight: var(--font-weight-medium);
						color: $color(white);
					}

					$me:hover{color: $color(white);}
				`),
				Button(
					Div(InlineStyle("$me{cursor: pointer; display: flex; align-items: center;}"),
						dropdownHeader, Span(Text(" ")),
						Icon(ICON_CHEVRON_DOWN, 16),
					),
				),
			),
			Div(
				Class("dropdown"),
				InlineStyle(`$me{display: none; position: absolute; right: 0; z-index: 10; padding-top: $1; padding-bottom: $1; margin-top: $2; width: $48; background-color: $color(white); transform-origin: top right; box-shadow: var(--shadow-lg);}`),
				TabIndex("-1"),
				dropdownItems,
			),
			InlineScript(`
				let button = me(".button", me());
				let dropdown = me(".dropdown", me());

				button.on("click", ev => { toggleShowHide(dropdown) });
				onClickOutsideOrEscape(me(), () => { hide(dropdown) });
			`),
		)
	}

	navbarDropdownItem := func(name string, url string, newPage bool) Node {
		return A(InlineStyle(`$me{display: block; padding-top: $2; padding-bottom: $2; padding-left: $4; padding-right: $4; font-size: var(--text-sm); line-height: $5; color: $color(neutral-700); } $me:hover{background: $color(neutral-100);}`), Href(url), TabIndex("-1"), Text(name), If(newPage, Target("_blank")))
	}

	navbarLink := func(name string, url string, newPage bool) Node {
		return A(
			InlineStyle(`$me{ padding-left: $3; padding-right: $3; padding-top: $2; padding-bottom: $2; font-size: var(--text-sm); font-weight: var(--font-weight-medium); color: $color(white);}`),
			InlineStyle("$me:hover{color: $color(white);}"),
			Href(url),
			Text(name),
			If(newPage, Target("_blank")),
		)
	}

	return RootLayout(title+" | MaxWarden",
		Body(InlineStyle("$me{background-color: $color(light-grey); height: 100%;}"),
			Div(InlineStyle("$me{min-height: 100%}"),
				Nav(InlineStyle("$me{background-color: $color(deep-blue);}"),
					Div(InlineStyle("$me{margin-left: auto; margin-right: auto; max-width: var(--container-7xl);}"),
						Div(InlineStyle("$me{display: flex; height: $16; align-items: center; justify-content: space-between;}"),
							Div(InlineStyle("$me{align-items: center; display: flex;}"),
								Div(InlineStyle("@media $lg-{ $me{display: block;}}"),
									Div(InlineStyle(`$me{margin-left: $1; display: flex; align-items: baseline;} $me:not(:last-child){ margin-left: $4; }`),
										Map(NavGroups, func(nav NavGroup) Node {
											if len(nav.SubGroup) > 0 {
												return navbarDropdown(
													Text(nav.Title),
													Map(nav.SubGroup, func(sub NavGroup) Node {
														return navbarDropdownItem(sub.Title, sub.URL, sub.NewTab)
													}),
												)
											} else {
												return navbarLink(nav.Title, nav.URL, nav.NewTab)
											}
										}),
									),
								),
							),
							Div(InlineStyle("$me{display: none;} @media $md { $me{ display: block; }}"),
								Div(InlineStyle("$me{ margin-left: $4; display: flex; align-items: center;} @media $md { $me{ margin-left: $6;}}"),
									Div(InlineStyle("$me{position: relative; margin-left: $3;}"),
										navbarDropdown(
											Icon(ICON_USERS, 24),
											Group{
												navbarDropdownItem("My Profile", "/app/account", false),
												navbarDropdownItem("Lock Vault", "/auth/logout", false),
											},
										),
									),
								),
							),
							Div(InlineStyle("$me{margin-right: $2; display: flex;} @media $md{ $me{display: none;}}"),
								Button(
									InlineStyle("$me{position: relative; display: inline-flex; justify-items: center; padding: $2; color: $color(neutral-400)}"),
									InlineStyle("$me:hover{color: $color(white); background-color: $color(neutral-900);}"),
									Type("button"),
									Span(InlineStyle("$me{position: absolute;}")),
									Icon(ICON_MENU, 24),
								),
							),
						),
					),
					Div(InlineStyle("@media $md { $me {display: none; }}"),
						Div(Class("space-y-1 px-2 pb-3 pt-2 sm:px-3"),
							A(Href("/app/dashboard"), Class("block hover:bg-neutral-900 px-3 py-2 text-base font-medium text-white"), Text("Dashboard")),
						),
						Div(Class("border-t border-neutral-700 pb-3 pt-4"),
							Div(Class("flex items-center px-5"),
								Div(Class("flex-shrink-0"),
									Img(Class("h-10 w-10 rounded-full"), Src(""), Alt("profile picture")),
								),
								Div(Class("ml-3"),
									// Div(Class("text-base/5 font-medium text-white"), Text(identity.User.Firstname+" "+identity.User.Lastname)),
									// Div(Class("text-sm font-medium text-neutral-400"), Text(identity.User.Email)),
								),
							),
							Div(Class("mt-3 space-y-1 px-2"),
								A(Href("/auth/logout"), Class("block px-3 py-2 text-base font-medium text-neutral-200 hover:bg-neutral-900 hover:text-white"), Text("Log out")),
							),
						),
					),
				),
				Header(InlineStyle("$me{background-color: $color(white); box-shadow: var(--shadow-sm);}"),
					Div(InlineStyle("$me{margin-left: auto; margin-right: auto; max-width: var(--container-7xl); padding: $4;} @media $lg { $me{ padding-left: $8; padding-right: $8;}}"),
						H1(InlineStyle("$me{font-size: var(--text-3xl); font-weight: var(--font-weight-bold); color: $color(neutral-950); letter-spacing: var(--tracking-tight);}"), Text(title)),
					),
				),
				Main(
					Div(InlineStyle("$me{margin-left: auto; margin-right: auto; max-width: var(--container-7xl); padding: $6 $4;} @media $sm { $me{padding-left: $6; padding-right: $6; }} @media $lg { $me{padding-left: $8; padding-right: $8;}}"),
						Group(children),
					),
				),
			),
		),
	)
}
