package mitm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func httpDump(reqDump []byte, resp *http.Response) {

}

func ParseReq(b []byte) (*http.Request, error) {
	// func ReadRequest(b *bufio.Reader) (req *Request, err error) { return readRequest(b, deleteHostHeader) }
	fmt.Println(string(b))
	fmt.Println("-----------------------")
	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	buf.Write(b)
	bufr := bufio.NewReader(buf)
	return http.ReadRequest(bufr)
}
