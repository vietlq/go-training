package main

import (
    "io"
    "os"
    "strings"
)

/*
Exercise: rot13Reader
A common pattern is an io.Reader that wraps another io.Reader, modifying the stream in some way.
For example, the gzip.NewReader function takes an io.Reader (a stream of compressed data) and returns a *gzip.Reader that also implements io.Reader (a stream of the decompressed data).
Implement a rot13Reader that implements io.Reader and reads from an io.Reader, modifying the stream by applying the rot13 substitution cipher to all alphabetical characters.
The rot13Reader type is provided for you. Make it an io.Reader by implementing its Read method.
*/

type rot13Reader struct {
    r io.Reader
}

func rot13(c byte) byte {
    if (c >= 'A' && c <= 'M') || (c >= 'a' && c <= 'm') {
        return c + 13;
    } else if (c >= 'N' && c <= 'Z') || (c >= 'n' && c <= 'z') {
        return c - 13;
    }
    return c
}

func (rot rot13Reader) Read(b []byte) (int, error) {
    // Read from the reader chained before
    n, err := rot.r.Read(b)
    // If we got good read, perform ROT13
    if err == nil {
        for i := range b {
            b[i] = rot13(b[i])
        }
    }
    return n, err
}

func main() {
    s := strings.NewReader("Lbh penpxrq gur pbqr!\n")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
}
