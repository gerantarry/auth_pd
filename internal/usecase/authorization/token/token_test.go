package token

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	plainText     = "some really really really long plaintext"
	encryptedText = "C�&y.\u0016k�w[��\f��\u0015����\u0014����I�l�'�5"
)

func TestEncrypt(t *testing.T) {
	enc, err := encrypt([]byte(plainText))
	if err != nil {
		require.FailNow(t, "ошибка при шифровании")
	}
	s := string(enc)
	err = os.WriteFile("token_test_result.txt", enc, 0666)

	fmt.Printf("Результат: %v", s)
}

func TestDecrypt(t *testing.T) {
	dec, err := decrypt([]byte(encryptedText))
	if err != nil {
		require.FailNow(t, "ошибка при дешифрации")
	}
	s := string(dec)
	fmt.Printf("Результат: %v", s)
}
