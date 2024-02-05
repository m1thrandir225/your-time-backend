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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateTaskApi(t *testing.T) {
	user := randomUser()

	task := randomTask(user);

	testCases := []struct {
		name 	string
		body 	gin.H
		build 	func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H {
				"title": task.Title,
				"description": task.Description.String,
				"reminder_date": "2021-07-13T15:28:51.818095+00:00",
				"due_date": "2021-07-13T15:28:51.818095+00:00",
				"user_id": user.ID.String(),
			},
			build: func(store *mockdb.MockStore) {
				arg := db.CreateTaskParams {
					Title: 	 task.Title,
					Description: task.Description,
					ReminderDate: task.ReminderDate,
					DueDate: task.DueDate,
					UserID: user.ID,
				}

				store.EXPECT().CreateTask(gomock.Any(), gomock.Eq(arg)).Times(1).Return(task, nil)
			},	
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				
			},
		},
		{
			name: "InternalError",
			body: gin.H {
				"title": task.Title,
				"description": task.Description.String,
				"reminder_date": "2021-07-13T15:28:51.818095+00:00",
				"due_date": "2021-07-13T15:28:51.818095+00:00",
				"user_id": user.ID.String(),
			},
			build: func(store *mockdb.MockStore) {
				arg := db.CreateTaskParams {
					Title: 	 task.Title,
					Description: task.Description,
					ReminderDate: task.ReminderDate,
					DueDate: task.DueDate,
					UserID: user.ID,
				}

				store.EXPECT().CreateTask(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H {
				"title": task.Title,
				"description": task.Description.String,
				"reminder_date": "2021-07-13T15:28:51.818095+00:00",
				"due_date": "2021-07-13T15:28:51.818095+00:00",
				"user_id": "invalid",
			},
			build: func(store *mockdb.MockStore) {


				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(0)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidDueDate",
			body: gin.H {
				"title": task.Title,
				"description": task.Description.String,
				"reminder_date": "2021-07-13T15:28:51.818095+00:00",
				"due_date": "invalid",
				"user_id": user.ID.String(),
			},
			build: func(store *mockdb.MockStore) {

				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(0)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "InvalidReminderDate",
			body: gin.H {
				"title": task.Title,
				"description": task.Description.String,
				"reminder_date": "invalid",
				"due_date": "2021-07-13T15:28:51.818095+00:00",
				"user_id": user.ID.String(),
			},
			build: func(store *mockdb.MockStore) {
				store.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Times(0)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.build(store)

			server := newTestServer(t, store)


			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			url := "/tasks"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})

	}
}

func randomTask(user db.User) db.Task {
	dueDate, err := time.Parse(time.RFC3339, "2021-07-13T15:28:51.818095+00:00")

	if err != nil {
		panic(err)
	}


	return db.Task{
		ID:           uuid.New(),
		UserID:       user.ID,
		Title:        util.RandomString(6),
		Description:  sql.NullString{String: util.RandomString(6), Valid: true},
		ReminderDate: sql.NullTime{Time: dueDate, Valid: true},
		DueDate:      dueDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
