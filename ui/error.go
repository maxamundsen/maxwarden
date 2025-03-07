package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ErrorPage(status int) Node {
	return RootLayout("An error has occurred.",
		Section(
			InlineStyle(`
				$me {
					background: $color(white);
				}
			`),
			Div(
				InlineStyle(`
					$me {
						padding-top: $8;
						padding-bottom: $8;
						padding-left: $4;
						padding-right: $4;
						margin-right: auto;
						margin-left: auto;
						max-width: var(--container-xl);
					}

					@media $lg {
						$me {
							padding-top: $16;
							padding-bottom: $16;
							padding-right: $6;
							padding-left: $6;
						}
					}
				`),
				Div(
					InlineStyle(`
						$me {
							margin-right: auto;
							margin-left: auto;
							max-width: var(--container-sm);
							text-align: center
						}
					`),
					H1(
						InlineStyle(`
							$me {
								margin-bottom: $4;
								font-size: var(--text-7xl);
								color: $color(indigo-blue);
								letter-spacing: var(--tracking-tight);
							}

							@media $lg {
								$me {
									font-size: var(--text-9xl);
								}
							}
						`),
						ToText(status),
					),
					P(
						InlineStyle(`
							$me {
								margin-bottom: $4;
								font-size: var(--text-lg);
								color: $color(neutral-500);
							}
						`),
						Text("We're sorry, an error has occured."),
					),
					A(
						Href("/"),
						InlineStyle(`
						    $me {
						        display: inline-flex;
						        color: $color(white);
						        background: $color(indigo-blue);
						        font-weight: var(--font-weight-medium);
								border-radius: var(--radius-xs);
						        font-size: var(--text-sm);
						        padding-left: $5;
						        padding-right: $5;
						        padding-top: $3;
						        padding-bottom: $3;
						        text-align: center;
						        margin-top: $4;
						        margin-bottom: $4;
						    }

						    $me:hover {
						    	background: $color(neutral-950);
						    }
						`),
						Text("Back to Homepage"),
					),
				),
			),
		),
	)
}
