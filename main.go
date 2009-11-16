package main
import ("./canvas";
        "./shader";
        "./file";
        "./gyu3d";
        "math";
        "image";
        "image/png";
        "fmt";
		"time";
		"flag";
		"strconv";
)

func main() {
	screen_width   := 1280;
	screen_height  := 720;
	benchmark_mode := false; //  true -> loop 50 times to benchmark
	
	canvas := canvas.EmptyCanvasImage(screen_width, screen_height, canvas.Color32(0xff000000));
	
	g := new(gyu3d.G3DContext);
	g.SetViewport(gyu3d.StandardViewport(screen_width, screen_height));
	g.SetProjection(
		gyu3d.IdentityM44().SetPerspective(float(screen_width)*0.006, float(screen_height)*0.006, 1.8, 100.0)
	);
	
	lx, ly := .1, -0.3;
	if len(flag.Args()) == 2 {
		lx, _ = strconv.Atof(flag.Args()[0]);
		ly, _ = strconv.Atof(flag.Args()[1]);
	}
	
	g.ResetTransforms();
	g.ViewMatrix.LookAt(0,1,0, -1.2,1.9,7, -0.5,0.9,0);
	g.Update();

	// draw
	drawTestScene(canvas, g, lx, ly, benchmark_mode);

	// export image
	outfile,e := file.WritableFile("./out.png");
	if e == nil {
		png.Encode(outfile, canvas);
		outfile.Close();
	}
	 
	fmt.Printf("done\n");
}

func createNormalTexture() image.Image {
	size := 512;

	tex := canvas.EmptyCanvasImage(size, size, canvas.Color32(0xff000000));
	for y := 0;y < size;y++ {
		for x := 0;x < size;x++ {
			sy  := math.Cos(float64(y)*0.18);
			
			v := gyu3d.ZeroVec3().Set(float((x%11)-5)*0.2, float(sy*sy), 7.5).Normalize();
			r := uint32(v.X * 127.0 + 128.0);
			g := uint32(v.Y * 127.0 + 128.0);
			b := uint32(v.Z * 127.0 + 128.0);
			
			tex.SetAt(x, y, canvas.Color32(0xff000000 | (r<<16) | (g<<8) | b));
		}
	}
	
	return tex;
}

func createTestMesh() *Mesh {
	m := newMesh(81, 128);
	
	vi := 0;
	for y := 0;y < 9;y++ {	
		for x := 0;x <= 9;x++ {
			
			th0 := (float64(x ) - 4.0)*0.2;
			th1 := (float64(-y) + 4.0)*0.2;
			
			nx, ny, nz := float(math.Sin(th0)), float(math.Sin(th1)), float(math.Cos(th1)*math.Cos(th0));
			
			if (x>0) {
				vx := nx*11.0;
				vy := ny*11.0;
				vz := -9.0 + nz*11.0;
				m.Vertices[vi-1].BN.Set( 
					vx - m.Vertices[vi-1].pos.X,
					vy - m.Vertices[vi-1].pos.Y,
					vz - m.Vertices[vi-1].pos.Z
					).Normalize();
				m.Vertices[vi-1].TN.Xp(
					m.Vertices[vi-1].N.X, m.Vertices[vi-1].N.Y, m.Vertices[vi-1].N.Z,
					m.Vertices[vi-1].BN.X, m.Vertices[vi-1].BN.Y, m.Vertices[vi-1].BN.Z).Normalize();
			}
			
			if (x<9) {
				m.Vertices[vi].N.X, m.Vertices[vi].N.Y, m.Vertices[vi].N.Z = nx, ny, nz;
				vx := nx*11.0;
				vy := ny*11.0;
				vz := -9.0 + nz*11.0;
				m.Vertices[vi].N.Normalize();
				
				m.Vertices[vi].pos.Set(vx, vy, vz);
				m.Vertices[vi].uv1.Set(float(x)/8.0, float(y)/8.0);
				
				vi++;
			}
		}
	}

	ii := 0;
	for y := 0;y < 8;y++ {	
		for x := 0;x < 8;x++ {
			m.Indices[ii]   = uint(x + y*9);
			m.Indices[ii+1] = m.Indices[ii]+1;
			m.Indices[ii+2] = m.Indices[ii]+9;
			ii+= 3;
			
			m.Indices[ii]   = m.Indices[ii-1];
			m.Indices[ii+1] = m.Indices[ii-2];
			m.Indices[ii+2] = m.Indices[ii-1]+1;

			ii+= 3;
		}
	}
		
	return m;
}

func drawTestScene(canvas *canvas.CanvasImage, g *gyu3d.G3DContext, lx, ly float, benchmark_mode bool) {
	mesh := createTestMesh();

	// load texture
	infile,_ := file.ReadableFile("./texture.png");
	tex, _ := png.Decode(infile);
	infile.Close();
	
	// create shader
	sh := new(NMapShader);
	sh.sampler0 = newSampler(tex);
	sh.sampler1 = newSampler(createNormalTexture());
	sh.light = gyu3d.ZeroVec3().Set(lx, ly, -1).Normalize();
		
	shader.TransformVertices(mesh.ShDataList, sh, g);
	vlist := mesh.Vertices;
	ilist := mesh.Indices;
	for i := range vlist {
		vlist[i].CalcViewportPosition(g.Viewport());
	}
	
	fmt.Printf(" start\n");
	starttime := time.Nanoseconds();
	
	var loops int;
	if benchmark_mode {loops=50} else {loops=1}
	for k := 0;k < loops;k++{
		ii := 0;
		flen := mesh.FacesCount;
		ch := make(chan uint, flen);
		for fi := uint(0);fi < flen;fi++ {
			go DrawFace(canvas, vlist[ ilist[ii] ], vlist[ ilist[ii+1] ], vlist[ ilist[ii+2] ], sh, ch, fi);
			ii += 3;
		}
		
		for i := uint(0);i < flen;i++ {		
			<- ch;
		}

/*	
		fmt.Printf("[LOOP %d]\n", k);
		for i := uint(0);i < flen;i++ {		
			fmt.Printf("%d  ", <- ch);
		}
		fmt.Printf("\n-------------------------------------\n");
*/
	}
	et := time.Nanoseconds() - starttime;
	fmt.Printf(" %d ms\n", et/1000000);

}

func DrawFace(cv *canvas.CanvasImage, v1, v2, v3 *ModelVertex, sh shader.PixelShader, ch chan uint, index uint) {
	cv.DrawTriangle(
		int(v1.viewportPos.X),
		int(v1.viewportPos.Y),

		int(v2.viewportPos.X),
		int(v2.viewportPos.Y),

		int(v3.viewportPos.X),
		int(v3.viewportPos.Y),
		
		v1.uv1.U, v1.uv1.V,
		v2.uv1.U, v2.uv1.V,
		v3.uv1.U, v3.uv1.V,
		
		v1.localLight, v2.localLight, v3.localLight,
		
		sh
	);
	
	ch <- index;
}

// Vertex
type ModelVertex struct {
	pos  *gyu3d.Vec3;
	N    *gyu3d.Vec3;
	BN   *gyu3d.Vec3;
	TN   *gyu3d.Vec3;
	localLight *gyu3d.Vec3;
	uv1  *gyu3d.TexCoord2D;
	wpos *gyu3d.Vec3;
	spos *gyu3d.Vec4;
	viewportPos *gyu3d.Vec4;
}

func newModelVertex() *ModelVertex {
	v := new(ModelVertex);
	v.pos         = gyu3d.ZeroVec3();
	v.N           = gyu3d.ZeroVec3();
	v.BN          = gyu3d.ZeroVec3();
	v.TN          = gyu3d.ZeroVec3();
	v.localLight  = gyu3d.ZeroVec3();
	v.uv1         = new(gyu3d.TexCoord2D);
	v.wpos        = gyu3d.ZeroVec3();
	v.spos        = gyu3d.ZeroVec4();
	v.viewportPos = gyu3d.ZeroVec4();
	
	return v;
}

func (v *ModelVertex) CalcViewportPosition(vp *gyu3d.Viewport) {
	v.viewportPos.X, v.viewportPos.Y = vp.FromNormalizedCoord(v.spos.X, v.spos.Y);
	v.viewportPos.Z, v.viewportPos.W = v.spos.Z, v.spos.W;
}

func (v *ModelVertex) GetTexCoord2D(index int) *gyu3d.TexCoord2D {
	return v.uv1;
}

func (v *ModelVertex) GetVec3(index int) *gyu3d.Vec3 {
	switch index {
		case 0: return v.pos;
		case 1: return v.N;
		case 2: return v.BN;
		case 3: return v.TN;
		case 4: return v.localLight;
	}

	return v.wpos;
}

func (v *ModelVertex) GetVec4(index int) *gyu3d.Vec4 {
	return v.spos;
}

// My Texture Sampler & Pixel Shader
type Sampler struct {
	texture image.Image;
	tw, th int;
} 

func newSampler(img image.Image) *Sampler{
	s := new(Sampler);
	s.texture = img;
	s.tw = img.Width();
	s.th = img.Height();
	return s;
}

func (smp *Sampler) Nearest(tu, tv float) image.Color {
	x := int(tu*float(smp.tw)+0.5);
	y := int(tv*float(smp.th)+0.5);
	if x>=smp.tw {x=smp.tw-1}
	else if x<0  {x=0};
	if y>=smp.th {y=smp.th-1}
	else if y<0  {y=0};

	return smp.texture.At(x,y);
//	r,g,b,a := smp.texture.At(x,y).RGBA();
//	return shader.PSColor(a>>24<<24 | r>>24<<16 | g>>24<<8 | b>>24);
}


type NMapShader struct {
	sampler0, sampler1 *Sampler;
	light *gyu3d.Vec3;
}

func (sh *NMapShader) DoVertex(vsdata *shader.ShaderData, g *gyu3d.G3DContext) {
	pos  := vsdata.GetVec3(0);
	wpos := vsdata.GetVec3(1);
	spos := vsdata.GetVec4(0);

	N  := vsdata.GetVec3(1);
	BN := vsdata.GetVec3(2);
	TN := vsdata.GetVec3(3);
	L  := vsdata.GetVec3(4);
	
	mT := g.ModelView;
	mP := g.GetProjection();
	
	planeM := gyu3d.IdentityM44().SetVectors(BN, TN, N).Transpose33();
	L.X, L.Y, L.Z, _ = planeM.TransVec3(sh.light.Components()); 
	
	wpos.X, wpos.Y, wpos.Z, _      = mT.TransVec3(pos.X, pos.Y, pos.Z);
	spos.X, spos.Y, spos.Z, spos.W = mP.TransVec3(wpos.X, wpos.Y, wpos.Z);
	spos.X, spos.Y, spos.Z = spos.X/spos.W, spos.Y/spos.W, spos.Z/spos.W;
}

func (sh *NMapShader) DoPixel(tu, tv, nx, ny, nz float) (shader.PSColor, bool) {
	L := gyu3d.ZeroVec3().Set(nx, ny, nz).Normalize();
	
	c0 := sh.sampler0.Nearest(tu, tv);
	cN := sh.sampler1.Nearest(tu, tv);
	
	ntx, nty, ntz, a := cN.RGBA();
	fx, fy, fz := float(int(ntx>>24)-128)/127.0, float(int(nty>>24)-128)/127.0, float(int(ntz>>24)-128)/127.0;
	r, g, b, a := c0.RGBA(); 
	f := L.Dp(fx, fy, fz);
	
	return shader.PSColor(a>>24<<24 | r>>24<<16 | g>>24<<8 | b>>24).Mul(-f*f*f), true;
}

// Mesh

type Mesh struct {
	Vertices   []*ModelVertex;
	ShDataList []shader.ShaderData;
	Indices    []uint;
	
	FacesCount uint;
}

func newMesh(nVerts, nFaces uint) *Mesh {
	m := new(Mesh);
	m.FacesCount = nFaces;
	m.Indices  = make([]uint, nFaces*3);
	m.Vertices = make([]*ModelVertex, nVerts);
	m.ShDataList = make([]shader.ShaderData, nVerts);
	
	for i := range m.Vertices {
		m.Vertices[i]   = newModelVertex();
		m.ShDataList[i] = m.Vertices[i];
	}
	
	return m;
}
