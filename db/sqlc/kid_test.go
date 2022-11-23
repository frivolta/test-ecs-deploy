package db

import (
	"birdie/util"
	"context"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
)

func TestCreateKid(t *testing.T) {
	createRandomKid(t)
}
func TestGetAllKids(t *testing.T) {
	var ks []Kid
	for i := 0; i < 5; i++ {
		ks = append(ks, createRandomKid(t))
	}

	allK, e := testQueries.GetAllKids(context.Background())
	require.NoError(t, e)
	for _, v := range ks {
		require.NotEmpty(t, v.ID)
		require.True(t, slices.Contains(allK, v))
	}
}

func TestGetKid(t *testing.T) {
	k := createRandomKid(t)
	kid, e := testQueries.GetKid(context.Background(), k.ID)
	require.NoError(t, e)
	require.Equal(t, kid.Name, k.Name)
	require.Equal(t, kid.Surname, k.Surname)
	require.Equal(t, kid.ID, k.ID)
}

func createRandomKid(t *testing.T) Kid {
	arg := CreateKidParams{
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}

	k, e := testQueries.CreateKid(context.Background(), arg)
	require.NoError(t, e)
	require.NotEmpty(t, k.ID)
	require.Equal(t, arg.Name, k.Name)
	require.Equal(t, arg.Surname, k.Surname)
	return k
}