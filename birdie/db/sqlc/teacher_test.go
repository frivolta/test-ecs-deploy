package db

import (
	"birdie/util"
	"context"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
)

func createRandomTeacher(t *testing.T) Teacher {
	arg := CreateTeacherParams{
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
	th, e := testQueries.CreateTeacher(context.Background(), arg)
	require.NoError(t, e)
	require.NotEmpty(t, th)
	require.Equal(t, arg.Name, th.Name)
	require.Equal(t, arg.Surname, th.Surname)
	require.NotEmpty(t, th.ID)
	return th
}

func TestCreateTeacher(t *testing.T) {
	createRandomTeacher(t)
}

func TestGetAllTeachers(t *testing.T) {
	var ths []Teacher
	for i := 0; i < 5; i++ {
		ths = append(ths, createRandomTeacher(t))
	}
	allT, e := testQueries.GetAllTeachers(context.Background())
	require.NoError(t, e)
	for _, v := range ths {
		require.NotEmpty(t, v.ID)
		require.True(t, slices.Contains(allT, v))
	}
}

func TestGetTeacher(t *testing.T) {
	th := createRandomTeacher(t)
	teacher, e := testQueries.GetTeacher(context.Background(), th.ID)
	require.NoError(t, e)
	require.Equal(t, teacher.Name, th.Name)
	require.Equal(t, teacher.Surname, th.Surname)
	require.Equal(t, teacher.ID, th.ID)
}
