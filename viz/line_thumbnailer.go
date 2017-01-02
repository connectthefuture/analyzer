package viz

import (
	"github.com/gonum/plot/vg/draw"
)

type LineThumbnailer struct {
	draw.LineStyle
}

func (l *LineThumbnailer) Thumbnail(da *draw.Canvas) {
	y := da.Center().Y
	da.StrokeLine2(l.LineStyle, da.Min.X, y, da.Max.X, y)
}
