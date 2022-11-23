package db

import (
	"birdie/util"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomCarnet(t *testing.T) (Carnet, Kid) {
	d, _ := util.ConvertDate("2022-10-01")
	arg1 := CreateKidParams{
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
	k, e := testQueries.CreateKid(context.Background(), arg1)
	require.NoError(t, e)

	arg2 := CreateCarnetParams{
		Date:     d,
		Quantity: int32(3),
		KidID:    k.ID,
	}
	cn, e := testQueries.CreateCarnet(context.Background(), arg2)
	require.NoError(t, e)
	return cn, k
}

func TestCreateCarnet(t *testing.T) {
	createRandomCarnet(t)
}

func TestGetAllCarnets(t *testing.T) {
	var cns []GetAllCarnetsRow
	for i := 0; i < 5; i++ {
		cn, k := createRandomCarnet(t)
		c := GetAllCarnetsRow{
			ID:       cn.ID,
			Date:     cn.Date,
			Quantity: cn.Quantity,
			KidName: sql.NullString{
				String: k.Name,
				Valid:  true,
			},
			KidSurname: sql.NullString{
				String: k.Surname,
				Valid:  true,
			},
			KidID: sql.NullInt64{
				Int64: k.ID,
				Valid: true,
			},
		}
		cns = append(cns, c)
	}

	allCns, e := testQueries.GetAllCarnets(context.Background())
	require.NoError(t, e)
	for _, v := range cns {
		require.NotEmpty(t, v.ID)
		require.Contains(t, allCns, v)
	}
}

func TestGetCarnet(t *testing.T) {
	c, _ := createRandomCarnet(t)
	g, e := testQueries.GetCarnet(context.Background(), c.ID)
	require.NoError(t, e)
	require.Equal(t, c.ID, g.ID)
	require.Equal(t, g.KidID, c.KidID)
	require.Equal(t, g.Date, c.Date)
	require.Equal(t, g.Quantity, c.Quantity)
}
