package cli

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/tink-crypto/tink-go/v2/streamingaead"
	"github.com/tink-crypto/tink-go/v2/tink"
)

func newStreamingAEAD() *cobra.Command {
	return newGroup("streaming-aead", "Encrypt and decrypt streams with a streaming AEAD keyset",
		newStreamingCommand("encrypt", "Encrypt --in as a stream with the keyset's primary key", false),
		newStreamingCommand("decrypt", "Decrypt the stream --in with the keyset", true),
	)
}

// newStreamingCommand makes streaming-aead encrypt or decrypt. Streams are
// raw bytes, so there is no --encoding.
func newStreamingCommand(use, short string, decrypt bool) *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short}
	flags := addPrimitiveFlags(cmd, true)
	adFlags := addBytesFlags(cmd, "associated-data")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		ad, err := adFlags.get(cmd)
		if err != nil {
			return err
		}
		h, err := flags.handle(cmd, streamingAEADPrimitive, anyKeyset)
		if err != nil {
			return err
		}
		s, err := streamingaead.New(h)
		if err != nil {
			return err
		}
		in, err := flags.openData(cmd)
		if err != nil {
			return err
		}
		defer in.Close()
		return flags.stream(cmd, func(out io.Writer) error {
			if decrypt {
				return decryptStream(s, in, out, ad)
			}
			return encryptStream(s, in, out, ad)
		})
	}
	return cmd
}

func encryptStream(s tink.StreamingAEAD, in io.Reader, out io.Writer, ad []byte) error {
	w, err := s.NewEncryptingWriter(out, ad)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, in); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

// decryptStream copies the plaintext of in to out. An error reading in or
// writing out is an I/O error; any other error is a ciphertext that did not
// decrypt.
func decryptStream(s tink.StreamingAEAD, in io.Reader, out io.Writer, ad []byte) error {
	src, dst := &ioErrReader{r: in}, &ioErrWriter{w: out}
	r, err := s.NewDecryptingReader(src, ad)
	if err == nil {
		_, err = io.Copy(dst, r)
	}
	switch {
	case src.err != nil:
		return src.err
	case dst.err != nil:
		return dst.err
	case err != nil:
		return failed("decryption")
	}
	return nil
}

// ioErrReader keeps the first error of r other than io.EOF.
type ioErrReader struct {
	r   io.Reader
	err error
}

func (e *ioErrReader) Read(p []byte) (int, error) {
	n, err := e.r.Read(p)
	if err != nil && !errors.Is(err, io.EOF) && e.err == nil {
		e.err = err
	}
	return n, err
}

// ioErrWriter keeps the first error of w.
type ioErrWriter struct {
	w   io.Writer
	err error
}

func (e *ioErrWriter) Write(p []byte) (int, error) {
	n, err := e.w.Write(p)
	if err != nil && e.err == nil {
		e.err = err
	}
	return n, err
}
