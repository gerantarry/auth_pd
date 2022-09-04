package password

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	originPsw := "hgyuatwerqhhoi^123"
	hashPsw := HashPassword(originPsw)

	result := bcrypt.CompareHashAndPassword([]byte(hashPsw), []byte(originPsw))
	require.Nil(t, result, "Хэши не сопадают")
}
