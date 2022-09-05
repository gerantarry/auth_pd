package password

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

const originPsw = "hgyuatwerqhhoi^123"

func TestHashPassword(t *testing.T) {
	hashPsw := HashPassword(originPsw)
	result := bcrypt.CompareHashAndPassword([]byte(hashPsw), []byte(originPsw))
	require.Nil(t, result, "Хэши не сопадают")
}
