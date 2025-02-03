package sops

import "github.com/getsops/sops/v3/decrypt"

type sopsCLientWrapper struct{}

func (w *sopsCLientWrapper) DecryptFile(data []byte, format string) (cleartext []byte, err error) {
	return decrypt.Data(data, format)
}
