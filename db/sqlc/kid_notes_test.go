package db

import (
	"birdie/util"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateKidNote(t *testing.T) {
	createRandKidNote(t)
}

func TestGetAllKidNotes(t *testing.T) {
	var kns []KidNote
	for i := 0; i < 5; i++ {
		kn, _ := createRandKidNote(t)
		kns = append(kns, kn)
	}

	allKn, e := testQueries.GetAllKidNotes(context.Background())
	require.NoError(t, e)
	for _, v := range kns {
		require.NotEmpty(t, v.ID)
		require.Contains(t, allKn, v)
	}
}

func TestGetKidNote(t *testing.T) {
	kn, _ := createRandKidNote(t)
	gotKn, e := testQueries.GetKidNote(context.Background(), kn.ID)
	require.NoError(t, e)
	require.Equal(t, kn.KidID, gotKn.KidID)
	require.Equal(t, kn.Note, gotKn.Note)
	require.Equal(t, kn.Date, gotKn.Date)
	require.Equal(t, kn.Presence, gotKn.Presence)
	require.Equal(t, kn.HasMeal, gotKn.HasMeal)
}

func TestGetKidNoteByPeriod(t *testing.T) {
	kn, k := createRandKidNote(t)
	expRes := GetKidNotesByPeriodRow{
		ID:       kn.ID,
		Note:     kn.Note,
		Date:     kn.Date,
		Presence: kn.Presence,
		HasMeal:  kn.HasMeal,
		KidName: sql.NullString{
			String: k.Name,
			Valid:  true,
		},
		KidSurname: sql.NullString{
			String: k.Surname,
			Valid:  true,
		},
		KidID: sql.NullInt64{
			Int64: kn.KidID,
			Valid: true,
		},
	}
	// if it is today i can use the same date
	arg1 := GetKidNotesByPeriodParams{
		Date:   time.Now(),
		Date_2: time.Now(),
	}
	kns, e := testQueries.GetKidNotesByPeriod(context.Background(), arg1)
	require.NoError(t, e)
	require.Contains(t, kns, expRes)

	// Invalid
	arg2 := GetKidNotesByPeriodParams{
		Date:   time.Now().Add(time.Hour * 100000),
		Date_2: time.Now().Add(time.Hour * 200000),
	}
	kns, e = testQueries.GetKidNotesByPeriod(context.Background(), arg2)
	require.Empty(t, kns)

}

func TestUpdateKidNote(t *testing.T) {
	kn, _ := createRandKidNote(t)
	// Update data
	//@ToDo: cannot be multiple of the same, do it in api request validation
	np := []Presence{PresenceAFTERNOON, PresenceMORNING, PresenceEVENING, PresenceABSENT}
	nn := "New Note"
	// .Update data
	upd := kn
	upd.Note = nn
	upd.Presence = np

	arg := UpdateKidNoteParams{
		ID:       upd.ID,
		Note:     upd.Note,
		Presence: upd.Presence,
		HasMeal:  upd.HasMeal,
	}
	updatedKid, e := testQueries.UpdateKidNote(context.Background(), arg)
	require.NoError(t, e)
	require.Equal(t, updatedKid, upd)

}

func createRandKidNote(t *testing.T) (KidNote, Kid) {
	arg1 := CreateKidParams{
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
	presence := []Presence{PresenceAFTERNOON, PresenceEVENING}
	k, e := testQueries.CreateKid(context.Background(), arg1)
	require.NoError(t, e)
	arg2 := CreateKidNoteParams{
		Note:     util.RandomString(200),
		KidID:    k.ID,
		Presence: presence,
		HasMeal:  false,
		Date:     time.Now(),
	}
	kn, e := testQueries.CreateKidNote(context.Background(), arg2)
	require.NoError(t, e)
	require.NotEmpty(t, kn)
	return kn, k
}
