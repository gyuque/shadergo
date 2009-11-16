package canvas
import ("image";
		"./gyu3d";
        "./shader";
        "fmt";
		)

func minmax(a, b int) (min, max int, reversed bool) {
	if a < b {
		min = a;
		max = b;
		reversed = false;
	} else {
		max = a;
		min = b;
		reversed = true;
	}
	
	return;
}

type Color32 uint32
func (p Color32) RGBA() (r, g, b, a uint32) {
	x := uint32(p);
	a = x >> 24;
	a |= a << 8;
	a |= a << 16;
	r = (x >> 16) & 0xFF;
	r |= r << 8;
	r |= r << 16;
	g = (x >> 8) & 0xFF;
	g |= g << 8;
	g |= g << 16;
	b = x & 0xFF;
	b |= b << 8;
	b |= b << 16;
	return;
}


type CanvasImage struct {
	Pixel [][]Color32;
}

func (m *CanvasImage) At(x, y int) image.Color {
	return m.Pixel[y][x];
}

func (m *CanvasImage) Height() int {
	return len(m.Pixel);
}

func (m *CanvasImage) Width() int {
	if len(m.Pixel) == 0 {
		return 0
	}

	return len(m.Pixel[0]);
}

func EmptyCanvasImage(dx, dy int, fillColor Color32) *CanvasImage {
	pix    := make([][]Color32, dy);
	linear := make([]Color32, dx*dy);

	for i := range pix {
		pix[i] = linear[dx*i : dx*(i+1)];
		for k := range pix[i] {
			pix[i][k] = fillColor;
		}
	}
	return &CanvasImage{pix};
}

func (m *CanvasImage) SetAt(x, y int, c Color32) {
	m.Pixel[y][x] = c;
}

func (m *CanvasImage) DrawVerticalLine(x1, y1, length int, c Color32) {
	for y := 0;y < length;y++ {
		m.SetAt(x1, y1+y, c);
	}
}

func (m *CanvasImage) DrawLine(x1, y1, x2, y2 int, c Color32) {
	x_min, x_max, xrev := minmax(x1, x2);
	y_min, y_max, yrev := minmax(y1, y2);
	
	dx := x_max - x_min;
	dy := y_max - y_min;
	if dx == 0 {
		m.DrawVerticalLine(x_min, y_min, dy, c);
		return;
	}
	e := 0;
	
	if dy < dx {
		var ys, y int;
		if xrev==yrev {y = y_min; ys = 1;} else {y = y_max; ys = -1;}
		
		for x := x_min;x <= x_max;x++ {
			e += dy;
			if (e > dx) {
				e -= dx;
				y += ys;
			}
		
			m.SetAt(x, y, c);
		}
	} else {
		var xs, x int;
		if yrev==xrev {x = x_min; xs = 1;} else {x = x_max; xs = -1;}
			
		for y := y_min;y <= y_max;y++ {
			e += dx;
			if (e > dy) {
				e -= dy;
				x += xs;
			}
		
			m.SetAt(x, y, c);
		}
	}
}

func sort3(a, b, c int) (x, y, z, i1, i2, i3 int) {
	if a < c && b <= c {
		z  = c;
		i3 = 2;
		if a<b {x=a; y=b;  i1=0; i2=1;} else {y=a; x=b;  i2=0; i1=1;}
	} else if a < b && c < b {
		z = b;
		i3 = 1;
		if a<c {x=a; y=c;  i1=0; i2=2;} else {y=a; x=c;  i2=0; i1=2;}
	} else {
		z = a;
		i3 = 0;
		if b<c {x=b; y=c;  i1=1; i2=2;} else {y=b; x=c;  i2=1; i1=2;}
	}
		
	return;
}

func (m *CanvasImage) DrawTriangle(
						x1, y1, x2, y2, x3, y3 int,
						u1, v1, u2, v2, u3, v3 float,
						N1, N2, N3 *gyu3d.Vec3,
						sh shader.PixelShader
					) {	
	var ylist [3]int;
	var xmax, xmin int;
	var y_limit = m.Height();
	var x_limit = m.Width();
	
	xlist := [3]int{x1, x2, x3};	
	ulist := [3]float{u1, u2, u3};
	vlist := [3]float{v1, v2, v3};
	Nlist := [3]*gyu3d.Vec3{N1, N2, N3};
	
	{
		var i1, i2, i3 int;
		ylist[0], ylist[1], ylist[2], i1,i2,i3 = sort3(y1, y2, y3);
		xlist[0], xlist[1], xlist[2] = xlist[i1], xlist[i2], xlist[i3];
		ulist[0], ulist[1], ulist[2] = ulist[i1], ulist[i2], ulist[i3];
		vlist[0], vlist[1], vlist[2] = vlist[i1], vlist[i2], vlist[i3];
		Nlist[0], Nlist[1], Nlist[2] = Nlist[i1], Nlist[i2], Nlist[i3];
		xmin, _, xmax, _,_,_ = sort3(x1, x2, x3);
	}
		
	dy1 := ylist[1] - ylist[0];
	dx1 := xlist[1] - xlist[0];
	if (dx1<0) {dx1=-dx1}

	dy2 := ylist[2] - ylist[0];
	dx2 := xlist[2] - xlist[0];
	if (dx2<0) {dx2=-dx2}

	var e1, e2 int;
	var e1s, e2s int;
	
	e1x := xlist[0];
	if xlist[1]<xlist[0] {e1s=-1} else {e1s=1}

	e2x := xlist[0];
	if xlist[2]<xlist[0] {e2s=-1} else {e2s=1}
	  
	e1 = 0;
	e2 = 0;
	x1end := xlist[1];
	var yr1, yr2 float;
	ylen1 := float(ylist[1] - ylist[0]);
	yorg1 := float(ylist[0]);
	uorg1 := ulist[0];
	vorg1 := vlist[0];
	uend1 := ulist[1];
	vend1 := vlist[1];

	Norg1 := Nlist[0];
	Nend1 := Nlist[1];
	
	ylen2 := float(ylist[2] - ylist[0]);
	
	if (e2s<0 && xlist[0]>xlist[1] && dy2!=0) {
		e2 = -dx2;
	}
	
	for y := ylist[0];y <= ylist[2];y++ {
		if y == ylist[1] {
			dy1 = ylist[2] - ylist[1];
			dx1 = xlist[2] - xlist[1];
			if (dx1<0) {dx1=-dx1}
			
			e1x = xlist[1];
			x1end = xlist[2];
			ylen1 = float(ylist[2] - ylist[1]);
			yorg1 = float(ylist[1]);
			uorg1 = ulist[1];
			vorg1 = vlist[1];
			uend1 = ulist[2];
			vend1 = vlist[2];

			Norg1 = Nlist[1];
			Nend1 = Nlist[2];
			
			e1 = 0;
			if xlist[2]<xlist[1] {e1s=-1} else {e1s=1}
		}

		yr1 = (float(y)-yorg1)/(ylen1);
		yr2 = float(y-ylist[0])/(ylen2) ;

		_yr1 := 1.0-yr1;
		_yr2 := 1.0-yr2;
		side1_u  := _yr1*uorg1 + yr1*uend1;
		side1_v  := _yr1*vorg1 + yr1*vend1;
		
		side1_nx := _yr1*Norg1.X + yr1*Nend1.X;
		side1_ny := _yr1*Norg1.Y + yr1*Nend1.Y;
		side1_nz := _yr1*Norg1.Z + yr1*Nend1.Z;
		
		side2_u  := _yr2*ulist[0] + yr2*ulist[2];
		side2_v  := _yr2*vlist[0] + yr2*vlist[2];
		
		side2_nx := _yr2*Nlist[0].X + yr2*Nlist[2].X;
		side2_ny := _yr2*Nlist[0].Y + yr2*Nlist[2].Y;
		side2_nz := _yr2*Nlist[0].Z + yr2*Nlist[2].Z;
		
		e2 += dx2;
		for e2>dy2 && 0 != dy2 {
			e2 -= dy2;
			e2x += e2s;
			
			if (e2s<0 && e2x<=xlist[2]) ||
				(e2s>0 && e2x>=xlist[2])
				{break}
		}

		if y>=0 && y<y_limit && e2x != e1x {
			func(e1x, e2x, xmax, x_limit, y int, 
			side1_u, side1_v, side2_u, side2_v,
			side1_nx, side1_ny, side1_nz,
			side2_nx, side2_ny, side2_nz float) {
				var tu, tv, nx, ny, nz float;
				side1_is_right := false;
				on := false;

				xlen := e2x - e1x;
				if xlen<0 {xlen=-xlen; side1_is_right = true; e1x++;} else {e2x++;}
				
				xt := 0;
				for x := xmin;x <= xmax;x++ {
					if x == e1x {on = !on; }
					if x == e2x {on = !on; }
					if x<0 || x>=x_limit {continue}
					
					if on {
						xxt  := float(xt) / float(xlen);
						_xxt := 1.0-xxt;
						xt++;
						if !side1_is_right {
							tu = side1_u*_xxt + side2_u*xxt;
							tv = side1_v*_xxt + side2_v*xxt;
							
							nx = side1_nx*_xxt + side2_nx*xxt;
							ny = side1_ny*_xxt + side2_ny*xxt;
							nz = side1_nz*_xxt + side2_nz*xxt;
						} else {
							tu = side2_u*_xxt + side1_u*xxt;
							tv = side2_v*_xxt + side1_v*xxt;
							
							nx = side2_nx*_xxt + side1_nx*xxt;
							ny = side2_ny*_xxt + side1_ny*xxt;
							nz = side2_nz*_xxt + side1_nz*xxt;
						}
						
						pxcolor, accept := sh.DoPixel(tu, tv, nx,ny,nz);
						if (accept) {
							m.Pixel[y][x] = Color32(pxcolor);
						}
					}
				}	
			}(e1x, e2x, xmax, x_limit, y,
			side1_u, side1_v, side2_u, side2_v,
			side1_nx, side1_ny, side1_nz,
			side2_nx, side2_ny, side2_nz
			);
		}

		e1 += dx1;
		for e1>dy1 && 0 != dy1 {				
			e1 -= dy1;
			e1x += e1s;
			
			if (e1s<0 && e1x<=x1end) ||
				(e1s>0 && e1x>=x1end)
				{break}
		}
	}
}

func TestFuncs() {
	var list [3]int;
	var ilist [3]int;	

	list[0], list[1], list[2],
	ilist[0], ilist[1], ilist[2] = sort3(1,2,1); 
	fmt.Printf("%d %d %d  %d %d %d\n", list[0], list[1], list[2], ilist[0], ilist[1], ilist[2]);
}

func (m *CanvasImage) ColorModel() image.ColorModel {
	return ColorModel
}

func makeColor(r, g, b, a uint32) Color32 {
	return Color32(a>>24<<24 | r>>24<<16 | g>>24<<8 | b>>24)
}

func toColor(color image.Color) image.Color {
	if c, ok := color.(Color32); ok {
		return c
	}
	return makeColor(color.RGBA());
}

var ColorModel = image.ColorModelFunc(toColor)

