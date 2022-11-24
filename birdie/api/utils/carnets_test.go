package api

import (
	db "birdie/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var Cs *[]db.Carnet
var k *[]db.Kid
var kn *[]db.KidNote

func initCs() {
	Cs = &[]db.Carnet{
		{
			ID:       10,
			Date:     time.Now(),
			Quantity: 1,
			KidID:    101,
		},
		{
			ID:       20,
			Date:     time.Now(),
			Quantity: 2,
			KidID:    101,
		},
		{
			ID:       30,
			Date:     time.Now(),
			Quantity: 1,
			KidID:    201,
		},
	}
}

func initKids() {
	k = &[]db.Kid{
		{
			ID:      101,
			Name:    "Name",
			Surname: "Surname",
		},
		{
			ID:      102,
			Name:    "Name",
			Surname: "Surname",
		},
		{
			ID:      201,
			Name:    "Name",
			Surname: "Surname",
		},
		{
			ID:      301,
			Name:    "Name",
			Surname: "Surname",
		},
		{
			ID:      401,
			Name:    "Name",
			Surname: "Surname",
		},
	}
}

func initKidNotes() {
	kn = &[]db.KidNote{
		// kid 101 has credit of 1
		{
			ID:       4,
			Note:     "note",
			KidID:    101,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
		{
			ID:       5,
			Note:     "note",
			KidID:    101,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
		// kid 201 has a debit of 2
		{
			ID:       1,
			Note:     "note",
			KidID:    201,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
		{
			ID:       2,
			Note:     "note",
			KidID:    201,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
		{
			ID:       3,
			Note:     "note",
			KidID:    201,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
		// kid 301 has no carnets and remains at 0
		{
			ID:       6,
			Note:     "note",
			KidID:    301,
			Presence: nil,
			HasMeal:  false,
			Date:     time.Time{},
		},
		// kid 401 has no carnets and is in debt of 1
		{
			ID:       7,
			Note:     "note",
			KidID:    401,
			Presence: nil,
			HasMeal:  true,
			Date:     time.Time{},
		},
	}
}

func TestNewCUtil(t *testing.T) {
	initCs()
	u := NewCUtil(*Cs)
	require.Equal(t, u.Cs, *Cs)
	// Not yet initialized
	require.Empty(t, u.Kc)
}

func TestCUtil_ToKidCarnet(t *testing.T) {
	initCs()
	u := NewCUtil(*Cs)
	u.ToKidCarnet()
	Cs := *Cs
	m := map[int64][]db.Carnet{
		Cs[0].KidID: {
			{
				ID:       Cs[0].ID,
				Date:     Cs[0].Date,
				Quantity: Cs[0].Quantity,
				KidID:    Cs[0].KidID,
			},
			{
				ID:       Cs[1].ID,
				Date:     Cs[1].Date,
				Quantity: Cs[1].Quantity,
				KidID:    Cs[1].KidID,
			},
		},
		Cs[2].KidID: {
			{
				ID:       Cs[2].ID,
				Date:     Cs[2].Date,
				Quantity: Cs[2].Quantity,
				KidID:    Cs[2].KidID,
			},
		},
	}
	assert.Equal(t, m, u.Kc)
}

func TestCUtil_ToKidCarnetInfo(t *testing.T) {
	initCs()
	initKids()
	initKidNotes()
	u := NewCUtil(*Cs)
	u.ToKidCarnet()
	u.ToKidCarnetInfo(*kn, *k)

	// kid 101 has credit of 1
	require.Equal(t, u.Kci[101], CarnetInfo{
		KidName:     "Name",
		KidSurname:  "Surname",
		TotalBought: 3,
		TotalUsed:   2,
		TotalLeft:   1,
		TotalDebit:  0,
	})

	// kid 201 has a debit of 2
	require.Equal(t, u.Kci[201], CarnetInfo{
		KidName:     "Name",
		KidSurname:  "Surname",
		TotalBought: 1,
		TotalUsed:   3,
		TotalLeft:   0,
		TotalDebit:  2,
	})

	// kid 301 has no carnets and remains at 0
	require.Equal(t, u.Kci[301], CarnetInfo{
		KidName:     "Name",
		KidSurname:  "Surname",
		TotalBought: 0,
		TotalUsed:   0,
		TotalLeft:   0,
		TotalDebit:  0,
	})

	// kid 401 has no carnets and is in debt of 1
	require.Equal(t, u.Kci[401], CarnetInfo{
		KidName:     "Name",
		KidSurname:  "Surname",
		TotalBought: 0,
		TotalUsed:   1,
		TotalLeft:   0,
		TotalDebit:  1,
	})
}
