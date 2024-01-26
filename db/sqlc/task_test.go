package db

import (
	"context"
	"database/sql"
	"m1thrandir225/your_time/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)


func createRandomTask(t *testing.T, user User) Task {
	require.NotEmpty(t, user)

	dueDate := util.RandomDate();

	arg := CreateTaskParams {
		UserID: user.ID,
		Title: util.RandomString(6),
		Description: sql.NullString{String: util.RandomString(6), Valid: true},
		ReminderDate: sql.NullTime{Time: util.RandomReminderDate(dueDate), Valid: true},
		DueDate: dueDate,
	}

	task, err := testQueries.CreateTask(context.Background(), arg)

	require.NoError(t, err);

	require.NotEmpty(t, task)

	require.Equal(t, arg.UserID, task.UserID)

	require.Equal(t, arg.Title, task.Title)

	require.Equal(t, arg.Description, task.Description)

	require.WithinDuration(t, arg.ReminderDate.Time, task.ReminderDate.Time, time.Second)

	require.WithinDuration(t, arg.DueDate, task.DueDate, time.Second)

	require.NotZero(t, task.ID)

	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	user := createRandomUser(t)
	createRandomTask(t, user)
}

func TestGetTask(t *testing.T) {
	user := createRandomUser(t)

	task := createRandomTask(t, user)

	task2, err := testQueries.GetTaskByID(context.Background(), task.ID)

	require.NoError(t, err)

	require.NotEmpty(t, task2)

	require.Equal(t, task.ID, task2.ID)

	require.Equal(t, task.UserID, task2.UserID)

	require.Equal(t, task.Title, task2.Title)

	require.Equal(t, task.Description, task2.Description)

	require.WithinDuration(t, task.ReminderDate.Time, task2.ReminderDate.Time, time.Second)

	require.WithinDuration(t, task.DueDate, task2.DueDate, time.Second)
}


func TestGetTasksByUser(t *testing.T) {
	user := createRandomUser(t);

	for i := 0; i < 10; i++ {
		createRandomTask(t, user)
	}

	tasks, err := testQueries.GetTasksByUser(context.Background(), user.ID)

	require.NoError(t, err)

	require.NotEmpty(t, tasks)

	for _, task := range tasks {
		require.NotEmpty(t, task)

		require.Equal(t, user.ID, task.UserID)

		require.NotZero(t, task.ID)

		require.NotZero(t, task.CreatedAt)
	}
}