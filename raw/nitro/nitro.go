package nitro

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

type Archive struct {
	Files map[string]File
}

type File struct {
	Name string
	Data []byte
}

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) readShort() (v uint16, err error) {
	var buf [2]byte
	_, err = io.ReadFull(r.r, buf[:])
	if err != nil {
		return
	}
	v = binary.BigEndian.Uint16(buf[:])
	return
}

func (r *Reader) readInt() (v uint32, err error) {
	var buf [4]byte
	_, err = io.ReadFull(r.r, buf[:])
	if err != nil {
		return
	}
	v = binary.BigEndian.Uint32(buf[:])
	return
}

func (r *Reader) readString() (s string, err error) {
	length, err := r.readShort()
	if err != nil {
		return
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(r.r, buf)
	if err == nil {
		s = string(buf)
	}
	return
}

func (r *Reader) ReadArchive() (archive Archive, err error) {
	n, err := r.readShort()
	if err != nil {
		return
	}

	archive = Archive{Files: map[string]File{}}
	for range n {
		var file File
		file, err = r.ReadFile()
		if err != nil {
			return
		}
		archive.Files[file.Name] = file
	}

	return
}

func (r *Reader) ReadFile() (file File, err error) {
	name, err := r.readString()
	if err != nil {
		return
	}

	length, err := r.readInt()
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0, length*3/2))
	z, err := zlib.NewReader(io.LimitReader(r.r, int64(length)))
	if err != nil {
		return
	}

	_, err = io.Copy(buffer, z)
	if err != nil {
		return
	}

	file = File{name, buffer.Bytes()}
	return
}
