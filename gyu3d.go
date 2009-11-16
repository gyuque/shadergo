package gyu3d
import ("fmt"; "math")

type M44 struct {
	_11, _12, _13, _14 float;
	_21, _22, _23, _24 float;
	_31, _32, _33, _34 float;
	_41, _42, _43, _44 float;
}

type TexCoord2D struct {
	U, V float;
}

func (c *TexCoord2D) Set(u, v float) *TexCoord2D {
	c.U = u;
	c.V = v;
	return c;
}

type Vec3 struct {
	X, Y, Z float;
}

func ZeroVec3() *Vec3 {
	v := new(Vec3);
	v.X, v.Y, v.Z = 0, 0, 0;
	return v;
}

func (v *Vec3) Set(x, y, z float) *Vec3 {
	v.X, v.Y, v.Z = x, y, z;
	return v;
}

func (v *Vec3) Sub(x, y, z float) *Vec3 {
	v.X -= x;
	v.Y -= y;
	v.Z -= z;
	return v;
}

func (v *Vec3) Components() (float, float, float) {
	return v.X, v.Y, v.Z;
}

func (v *Vec3) Xp(vx, vy, vz, wx, wy, wz float) *Vec3 {
	v.X = (wy * vz) - (wz * vy);
	v.Y = (wz * vx) - (wx * vz);
	v.Z = (wx * vy) - (wy * vx);

	return v;
}

func (v *Vec3) Dp(vx, vy, vz float) float {
	return v.X*vx + v.Y*vy + v.Z*vz;
}

func (v *Vec3) Normalize() *Vec3 {
	nrm := float(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)));
	v.X /= nrm;
	v.Y /= nrm;
	v.Z /= nrm;

	return v;
}

type Vec4 struct {
	X, Y, Z, W float; 
}

func ZeroVec4() *Vec4 {
	v := new(Vec4);
	v.X, v.Y, v.Z, v.W = 0, 0, 0, 1;
	return v;
}

type Viewport struct {
	originX, originY, centerX, centerY, width, height int;
}

type G3DContext struct {
	viewport *Viewport;
	projectionMatrix *M44;
	ViewMatrix *M44;
	ModelTransMatrix *M44;

	ModelView *M44;
}

func (g *G3DContext) Update() {
	if nil == g.ModelView {g.ModelView = IdentityM44()}
	
	g.ModelView.mul(g.ModelTransMatrix, g.ViewMatrix);
}

func (g *G3DContext) SetViewport(v *Viewport) {
	g.viewport = v;
}

func (g *G3DContext) Viewport() *Viewport {
	return g.viewport;
}

func (g *G3DContext) SetProjection(m *M44) {
	g.projectionMatrix = m;
}

func (g *G3DContext) GetProjection() *M44 {
	return g.projectionMatrix;
}

func (g *G3DContext) ResetTransforms() {
	if nil == g.ModelTransMatrix
	{g.ModelTransMatrix = IdentityM44()}
	else
	{g.ModelTransMatrix.LoadIdentity()}

	if nil == g.ViewMatrix
	{g.ViewMatrix = IdentityM44()}
	else
	{g.ViewMatrix.LoadIdentity()}
}


func StandardViewport(w, h int) *Viewport {
	v := new(Viewport);
	
	v.originX = 0;
	v.originY = 0;

	v.centerX = w/2;
	v.centerY = h/2;
	
	v.width  = w;
	v.height = h;
	
	return v;
}

func (v *Viewport) FromNormalizedCoord(nx, ny float) (float, float) {
	return float(v.width )*0.5 * nx + float(v.centerX),
	       float(v.height)*0.5 *-ny + float(v.centerY);
}

func IdentityM44() *M44 {
	m := new(M44);
	m.LoadIdentity();
	
	return m;
}

func (m *M44) LoadIdentity() *M44 {
	m._11, m._22, m._33, m._44 = 1.0, 1.0, 1.0, 1.0;
	
	       m._12, m._13, m._14 = 0, 0, 0;
	m._21,        m._23, m._24 = 0, 0, 0;
	m._31, m._32,        m._34 = 0, 0, 0;
	m._41, m._42, m._43        = 0, 0, 0;

	return m;
}

func (m *M44) SetPerspective(vw, vh, z_near, z_far float) *M44 {
	m._11 = 2.0*z_near/vw;
	m._12 = 0;
	m._13 = 0;
	m._14 = 0;

	m._21 = 0;
	m._22 = 2*z_near/vh;
	m._23 = 0;
	m._24 = 0;

	m._31 = 0;
	m._32 = 0;
	m._33 = z_far/(z_far-z_near);
	m._34 = 1;

	m._41 = 0;
	m._42 = 0;
	m._43 = z_near*z_far/(z_near-z_far);
	m._44 = 0;
	
	return m;
}

func (m *M44) Dump() {
	fmt.Printf(": %f %f %f %f :\n"  , m._11, m._12, m._13, m._14);
	fmt.Printf(": %f %f %f %f :\n"  , m._21, m._22, m._23, m._24);
	fmt.Printf(": %f %f %f %f :\n"  , m._31, m._32, m._33, m._34);
	fmt.Printf(": %f %f %f %f :\n\n", m._41, m._42, m._43, m._44);
}

func (m *M44) TransVec3 (x, y, z float) (_x, _y, _z, _w float) {
	_x = x * m._11 + y * m._21 + z * m._31 + m._41;
	_y = x * m._12 + y * m._22 + z * m._32 + m._42;
	_z = x * m._13 + y * m._23 + z * m._33 + m._43;
	_w = x * m._14 + y * m._24 + z * m._34 + m._44;
	
	return;
}

func (m *M44) SetVectors(x, y, z *Vec3) *M44 {
	m._11, m._12, m._13 = x.X, x.Y, x.Z;
	m._21, m._22, m._23 = y.X, y.Y, y.Z;
	m._31, m._32, m._33 = z.X, z.Y, z.Z;
	
	return m;
}

func (m *M44) Transpose33() *M44 {
	m._12, m._21 = m._21, m._12;
	m._13, m._31 = m._31, m._13;
	m._23, m._32 = m._32, m._23;

	return m;
}

func (m *M44) mul(A, B *M44) *M44 {
	m._11 = A._11*B._11  +  A._12*B._21  +  A._13*B._31  +  A._14*B._41;
	m._12 = A._11*B._12  +  A._12*B._22  +  A._13*B._32  +  A._14*B._42;
	m._13 = A._11*B._13  +  A._12*B._23  +  A._13*B._33  +  A._14*B._43;
	m._14 = A._11*B._14  +  A._12*B._24  +  A._13*B._34  +  A._14*B._44;

	m._21 = A._21*B._11  +  A._22*B._21  +  A._23*B._31  +  A._24*B._41;
	m._22 = A._21*B._12  +  A._22*B._22  +  A._23*B._32  +  A._24*B._42;
	m._23 = A._21*B._13  +  A._22*B._23  +  A._23*B._33  +  A._24*B._43;
	m._24 = A._21*B._14  +  A._22*B._24  +  A._23*B._34  +  A._24*B._44;

	m._31 = A._31*B._11  +  A._32*B._21  +  A._33*B._31  +  A._34*B._41;
	m._32 = A._31*B._12  +  A._32*B._22  +  A._33*B._32  +  A._34*B._42;
	m._33 = A._31*B._13  +  A._32*B._23  +  A._33*B._33  +  A._34*B._43;
	m._34 = A._31*B._14  +  A._32*B._24  +  A._33*B._34  +  A._34*B._44;

	m._41 = A._41*B._11  +  A._42*B._21  +  A._43*B._31  +  A._44*B._41;
	m._42 = A._41*B._12  +  A._42*B._22  +  A._43*B._32  +  A._44*B._42;
	m._43 = A._41*B._13  +  A._42*B._23  +  A._43*B._33  +  A._44*B._43;
	m._44 = A._41*B._14  +  A._42*B._24  +  A._43*B._34  +  A._44*B._44;

	return m;
}

func (m *M44) LookAt(
		upX, upY, upZ,
		fromX, fromY, fromZ,
		atX, atY, atZ float) *M44 {
	aX := new(Vec3);
	aY := new(Vec3);
	aZ := new(Vec3).Set(atX, atY, atZ);
	vFrom := new(Vec3).Set(fromX, fromY, fromZ);

	aZ.Sub(fromX, fromY, fromZ).Normalize();
	aX.Xp(upX, upY, upZ, aZ.X, aZ.Y, aZ.Z).Normalize();
	aY.Xp(aZ.X, aZ.Y, aZ.Z, aX.X, aX.Y, aX.Z);

	m._11 = aX.X;  m._12 = aY.X;  m._13 = aZ.X;  m._14 = 0;
	m._21 = aX.Y;  m._22 = aY.Y;  m._23 = aZ.Y;  m._24 = 0;
	m._31 = aX.Z;  m._32 = aY.Z;  m._33 = aZ.Z;  m._34 = 0;

	m._41 = -vFrom.Dp(aX.Components());
	m._42 = -vFrom.Dp(aY.Components());
	m._43 = -vFrom.Dp(aZ.Components());
	m._44 = 1;

	return m;
}

