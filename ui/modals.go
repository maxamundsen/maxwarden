package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// DIALOG / MODAL
func ModalActuator(id string, contents Node) Node {
	return Span(
		InlineScriptf(`
			let act = me();
			let dialog = me("#%s_dialog");
			act.on("click", () => { dialog.showModal(); });
		`, id),
		contents,
	)
}

func ModalCloser() Node {
	return Class("modal-close-el")
}

func Modal(id string, header Node, body Node, closeElements []Node) Node {
	return Dialog(
		ID(id+"_dialog"),
		InlineStyle(`
			$me {
				top: 50%;
				left: 50%;
				translate: -50% -50%;
				z-index: 10;
				background: $color(white);
				box-shadow: var(--shadow-md);
				font-size: var(--text-sm);
			}

			$me::backdrop {
				background: $color(black/30);
			}
		`),

		// Header
		Div(
			InlineStyle(`
				$me {
					display: flex;
					align-items: center;
					justify-content: space-between;
					padding: $6 $8 $1 $8;
					font-size: var(--text-xl);
					font-weight: var(--font-weight-bold);
				}
			`),

			IfElse(header == nil,
				InlineStyle(`
					$me {
						flex-direction: row-reverse;
					}
				`),
				header,
			),


			Div(
				Class("modal-close-el"),
				InlineStyle("$me { cursor: pointer; }"),
				Icon(ICON_X_DIALOG_CLOSE, 24),
			),
		),

		//Modal contents
		Div(
			ID(id), // we use the passed id here so that swapping the content is easier
			InlineStyle(`
				$me {
					padding: $3 $8 $8 $8;
					color: $color(neutral-700);
				}
			`),
			body,
		),

		// footer
		If(len(closeElements) > 0,
			Div(
				InlineStyle(`
					$me {
						padding: $3 $8 $3 $8;
						display: flex;
						flex-direction: row-reverse;
						gap: $2;
					}
				`),

				Group(closeElements),
			),
		),

		InlineScript(`
			let dialog = me();
			let close_buttons = any(".modal-close-el", me())

			dialog.on("click", (ev) => {
				if (ev.target === dialog) {
					dialog.close();
				}
			});

			close_buttons.on("click", () => { dialog.close(); });
		`),
	)
}
