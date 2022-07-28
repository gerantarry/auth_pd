package router

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestSetLogger(t *testing.T) {

}

//TODO сделать тест параметризованным
const expectedFormat = "{    \"text\": \"some text\"}"
const expectedFormat2 = "{\r\n    \"text\": \"some text\"}"

func TestFormatData(t *testing.T) {
	//TODO убрать чтение из файла ?
	f, err := os.Open("C:\\Users\\Anton\\GolandProjects\\auth_pd\\resources\\reader_test.txt")
	require.Nil(t, err)
	defer f.Close()

	s, err := formatReaderData(strings.NewReader(expectedFormat2))
	require.Nil(t, err)
	assert.Equal(t,
		expectedFormat,
		s)
	fmt.Println(s)
}
