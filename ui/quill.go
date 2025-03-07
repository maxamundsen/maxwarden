package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// QuillJS editor components

func Quill() Node {
	return Div(
		InlineStyle(`
			$me {
				width: 100%;
				background: $color(white);
				min-height: $64;
			}
		`),

		InlineScript(`
			let editor = me();

			let quill = new Quill(editor, {
				theme: 'snow',
			});
		`),
	)
}