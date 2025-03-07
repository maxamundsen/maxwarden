package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Grid2x2(children ...Node) Node {
	return Div(
		InlineStyle(`
			$me {
				display: grid;
				grid-template-columns: repeat(2, 1fr);
				gap: $5;
				padding: $5
				max-width: var(--container-lg);
				margin: 0 auto;
			}

			@media $md- {
				$me {
					grid-template-columns: 1fr;
				}
			}
		`),
		Group(children),
	)
}