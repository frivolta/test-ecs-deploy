package db

import (
	"birdie/util"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
	"time"
)

func createRandTeachNote(t *testing.T) (TeacherNote, Teacher) {
	arg1 := CreateTeacherParams{
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
	th, e := testQueries.CreateTeacher(context.Background(), arg1)
	require.NoError(t, e)
	arg2 := CreateTeacherNoteParams{
		Note: util.RandomString(50),
		TeacherID: sql.NullInt64{
			Int64: th.ID,
			Valid: true,
		},
		Date: time.Now(),
	}
	tn, e := testQueries.CreateTeacherNote(context.Background(), arg2)
	require.NoError(t, e)
	require.NotEmpty(t, tn)
	return tn, th
}

func TestCreateNote(t *testing.T) {
	createRandTeachNote(t)
}

func TestGetAllTeacherNotes(t *testing.T) {
	var tns []TeacherNote
	for i := 0; i < 5; i++ {
		tn, _ := createRandTeachNote(t)
		tns = append(tns, tn)
	}

	allTn, e := testQueries.GetAllTeacherNotes(context.Background())
	require.NoError(t, e)
	for _, v := range tns {
		require.NotEmpty(t, v.ID)
		require.True(t, slices.Contains(allTn, v))
	}
}

func TestGetTeacherNotesByDate(t *testing.T) {
	n, th := createRandTeachNote(t)
	allNotes, e := testQueries.GetTeacherNotesByDate(context.Background(), n.Date)
	require.NoError(t, e)
	expRes := GetTeacherNotesByDateRow{
		ID:   n.ID,
		Note: n.Note,
		Date: n.Date,
		TeacherID: sql.NullInt64{
			Int64: th.ID,
			Valid: true,
		},
		TeacherName: sql.NullString{
			String: th.Name,
			Valid:  true,
		},
		TeacherSurname: sql.NullString{
			String: th.Surname,
			Valid:  true,
		},
	}
	require.Contains(t, allNotes, expRes)
	allNotes, e = testQueries.GetTeacherNotesByDate(context.Background(), n.Date.Add(time.Hour*48))
	require.NotContains(t, allNotes, n)
}

func TestQueries_GetTeacherNotesByPeriod(t *testing.T) {
	n, th := createRandTeachNote(t)
	expRes := GetTeacherNotesByPeriodRow{
		ID:   n.ID,
		Note: n.Note,
		Date: n.Date,
		TeacherID: sql.NullInt64{
			Int64: th.ID,
			Valid: true,
		},
		TeacherName: sql.NullString{
			String: th.Name,
			Valid:  true,
		},
		TeacherSurname: sql.NullString{
			String: th.Surname,
			Valid:  true,
		},
	}
	arg1 := GetTeacherNotesByPeriodParams{
		Date:   time.Now().Add(-time.Hour * 48),
		Date_2: time.Now().Add(time.Hour * 48),
	}
	notes, e := testQueries.GetTeacherNotesByPeriod(context.Background(), arg1)
	require.NoError(t, e)
	require.True(t, slices.Contains(notes, expRes))

	arg2 := GetTeacherNotesByPeriodParams{
		Date:   time.Now().Add(time.Hour * 100000),
		Date_2: time.Now().Add(time.Hour * 200000),
	}
	notes, e = testQueries.GetTeacherNotesByPeriod(context.Background(), arg2)
	require.Empty(t, notes)
}
