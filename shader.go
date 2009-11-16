package shader
import "./gyu3d"

type PSColor uint32;
func (p PSColor) RGBA() (r, g, b, a uint32) {
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

func (p PSColor) Mul(f float) PSColor {
	if (f<0) {f=0}
	x := uint32(p);
	a :=  x >> 24;
	r := (x >> 16) & 0xFF;
	g := (x >> 8) & 0xFF;
	b :=  x & 0xFF;
	
	r = uint32( float(r)*f );
	g = uint32( float(g)*f );
	b = uint32( float(b)*f );
	if (r>255) {r = 255}
	if (g>255) {g = 255}
	if (b>255) {b = 255}
	
	return PSColor((a<<24)|(r<<16)|(g<<8)|b);
}

type ShaderData interface {
	GetVec3(index int) *gyu3d.Vec3;
	GetVec4(index int) *gyu3d.Vec4;
	GetTexCoord2D(index int) *gyu3d.TexCoord2D;
}

type PixelShader interface {
	DoPixel(tu, tv, nx, ny, nz float) (PSColor, bool);
}

type VertexShader interface {
	DoVertex(vsdata *ShaderData, g *gyu3d.G3DContext);
}

func TransformVertices(list []ShaderData, sh VertexShader, g *gyu3d.G3DContext) {
	for i := range list {
		d := list[i];
		sh.DoVertex(&d, g);
	}
}