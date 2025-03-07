package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// BUTTONS
func ButtonUI(children ...Node) Node {
	return Button(
		InlineStyle(`
			$me {
				cursor: pointer;
				background-color: $color(deep-blue);
				color: $color(neutral-50);
				padding: $1.5 $8 $1.5 $8;
				font-size: var(--text-sm);
				border-radius: var(--radius-xs);
				font-weight: var(--font-weight-semibold);
			}

			$me:hover {
				background: $color(indigo-blue);
			}
		`),
		Group(children),
	)
}

func ButtonUIOutline(children ...Node) Node {
	return ButtonUI(
		InlineStyle(`
			$me {
				background: none;
				box-sizing: border-box;
                box-shadow: 0 0 0 0, inset 0 0 0 1px $color(gray-600);
				color: $color(gray-600);
			}

			$me:hover {
				background: none;
                box-shadow: 0 0 0 0, inset 0 0 0 2px $color(gray-600);
			}
		`),
		Group(children),
	)
}

func ButtonUIDanger(children ...Node) Node {
	return ButtonUI(
		InlineStyle(`
			$me {
				background: $color(red-700);
				color: $color(white);
			}

			$me:hover {
				background: $color(red-800);
			}
		`),
		Group(children),
	)
}

func ButtonUISuccess(children ...Node) Node {
	return ButtonUI(
		InlineStyle(`
			$me {
				background: $color(green-700);
				color: $color(white);
			}

			$me:hover {
				background: $color(green-800);
			}
		`),
		Group(children),
	)
}