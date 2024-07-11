package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCutHtml(t *testing.T) {
	req := require.New(t)

	in := "test test test"
	out, cut := CutHtml(in, 1, 20, 0)
	req.False(cut)
	req.Equal(in, out)

	out, cut = CutHtml(in, 1, 10, 0)
	req.True(cut)
	req.Equal("test test…", out)

	in = "<p>test <br>test test </p>"
	out, cut = CutHtml(in, 1, 10, 0)
	req.True(cut)
	req.Equal("<p>test…</p>", out)

	out, cut = CutHtml(in, 2, 10, 0)
	req.Equal(in, out)
	req.False(cut)

	in = "<p><b><i>test <br>test test</i></b></p>"
	out, cut = CutHtml(in, 2, 5, 0)
	req.True(cut)
	req.Equal("<p><b><i>test <br>test…</i></b></p>", out)

	in = "<p>test <br>test</p><p>test test test test</p> <p>test</p>"
	out, cut = CutHtml(in, 4, 7, 0)
	req.True(cut)
	req.Equal("<p>test <br>test</p><p>test test test…</p>", out)

	in = "<p>проверяем кириллические буквы</p>"
	out, cut = CutHtml(in, 1, 10, 0)
	req.True(cut)
	req.Equal("<p>проверяем…</p>", out)

	in = "<h1>test test test</h1><br>"
	out, cut = CutHtml(in, 1, 40, 0)
	req.True(cut)
	req.Equal("<h1>test…</h1>", out)

	in = "<p>test test test</p><img src='link'>"
	out, cut = CutHtml(in, 1, 40, 0)
	req.True(cut)
	req.Equal("<p>test test test…</p>", out)

	in = "<p>test test test</p><img src='link'><p>after image</p>"
	out, cut = CutHtml(in, 3, 40, 0)
	req.True(cut)
	req.Equal("<p>test test test…</p>", out)

	out, cut = CutHtml(in, 3, 40, 1)
	req.True(cut)
	req.Equal("<p>test test test</p><img src='link'>", out)

	out, cut = CutHtml(in, 4, 40, 1)
	req.False(cut)
}

func BenchmarkCutHtml(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		in := "<p>test test test</p><img src='link'><p>after image</p>"
		CutHtml(in, 3, 40, 0)
	}
}
