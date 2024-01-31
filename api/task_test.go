package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	mockdb "m1thrandir225/your_time/db/mock"
	db "m1thrandir225/your_time/db/sqlc"
	"m1thrandir225/your_time/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)



func TestCreateTaskApi(t *testing.T) {
	user := randomUser();

	testCases := []struct {
		name string
		task db.Task
		build func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(1).Return(randomTask(user), nil)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, randomTask(user))
			},
		},
		{
			name: "Unauthorized",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {},
		},
		{
			name: "InternalError",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(1).Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(1).Return(db.Task{}, sql.ErrNoRows)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDueDate",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(1).Return(db.Task{}, sql.ErrNoRows)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidReminderDate",
			task: randomTask(user),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(1).Return(db.Task{}, sql.ErrNoRows)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}


	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			defer controller.Finish()

			store := mockdb.NewMockStore(controller)

			tc.build(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.task)

			require.NoError(t, err)

			url := "/tasks"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	
	}
}


func randomTask(user db.User) db.Task {
	return db.Task{
		ID: uuid.New(),
		UserID: user.ID,
		Title: util.RandomString(6),
		Description: sql.NullString{String: util.RandomString(6), Valid: true},
		ReminderDate: sql.NullTime{Time: util.RandomReminderDate(util.RandomDate()), Valid: true},
		DueDate: util.RandomDate(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func requireBodyMatchTask(t *testing.T, body *bytes.Buffer, task db.Task) {
	data, err := io.ReadAll(body)

	require.NoError(t, err)

	var gotTask db.Task

	err = json.Unmarshal(data, &gotTask)

	require.NoError(t, err)

	require.Equal(t, task, gotTask)
}