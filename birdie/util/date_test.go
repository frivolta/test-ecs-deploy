package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConvertValidDate(t *testing.T) {
	initialDate := "2022-10-26"
	dt, e := ConvertDate(initialDate)
	require.NoError(t, e)
	fmt.Println(dt)
	require.Equal(t, dt.Day(), 26)
	require.Equal(t, dt.Month(), time.Month(10))
	require.Equal(t, dt.Year(), 2022)
}

func TestConvertInvalidDate(t *testing.T) {
	initialDate := "invalid"
	_, e := ConvertDate(initialDate)
	require.Error(t, e)
}
