package file
import ("os";
        "syscall";
)

type File struct {
	fd int;
	writable bool;
}

func WritableFile(name string) (file *File, err os.Error)
{
	fd, e := syscall.Open(name, syscall.O_CREAT|syscall.O_WRONLY, 0766);
	if e==0 {
		f := new(File);
		f.fd = fd;
		f.writable = true;
		return f, nil;
	}
	
	return nil, os.Errno(e);
}

func ReadableFile(name string) (file *File, err os.Error)
{
	fd, e := syscall.Open(name, syscall.O_RDONLY, 0);
	if e==0 {
		f := new(File);
		f.fd = fd;
		f.writable = false;
		return f, nil;
	}
	
	return nil, os.Errno(e);
}

func (f *File) Read(b []byte) (ret int, err os.Error) {
	if f == nil || f.writable {
		return -1, os.EINVAL
	}
	
	r, e := syscall.Read(f.fd, b);
	if e != 0 {
		err = os.Errno(e);
	}
	
	return int(r), err
}

func (f *File) Write(b []byte) (ret int, err os.Error) {
	if f == nil || !f.writable {
		return -1, os.EINVAL;
	}

	r, e := syscall.Write(f.fd, b);
	if e != 0 {
		err = os.Errno(e);
	}
	return int(r), err
}

func (f *File) Close() os.Error {
	if f == nil {
		return os.EINVAL
	}
	
	e := syscall.Close(f.fd);
	f.fd = -1;
	if e != 0 {
		return os.Errno(e);
	}
	return nil
}
