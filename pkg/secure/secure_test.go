package secure

import "testing"

func TestEncryptAES(t *testing.T) {
	info, err := EncryptAES([]byte("ccccc"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("aes info: %v", info)
}

func TestDecryptAES(t *testing.T) {
	info, err := DecryptAES([]byte{94, 134, 140, 161, 41, 127, 190, 117, 190, 145, 44, 157, 121, 29, 251, 229, 81, 45, 208, 247, 216})
	if err != nil {
		t.Error(err)
	}
	t.Logf("des info: %v", info)
}
