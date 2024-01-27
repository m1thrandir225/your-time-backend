package api

import (
	"bytes"
	"encoding/json"
	"io"
	mockdb "m1thrandir225/your_time/db/mock"
	db "m1thrandir225/your_time/db/sqlc"
	"m1thrandir225/your_time/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUserApi(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name string
		userID string
		build func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			userID: user.ID.String(),
			build: func (store *mockdb.MockStore) {
				store.EXPECT().
				GetUserByID(gomock.Any(), gomock.Eq(user.ID)).
				Times(1).
				Return(user, nil)
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i:= range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			defer controller.Finish()

			store := mockdb.NewMockStore(controller)

			tc.build(store)

			sever := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := "/users/" + tc.userID

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			sever.router.ServeHTTP(recorder, request)
		})
	}
}


func randomUser() db.User {
	return db.User{
		ID: uuid.New(),
		FirstName: util.RandomString(6),
		LastName: util.RandomString(6),
		Email: util.RandomEmail(),
		Password: util.RandomString(6),
		CreatedAt: util.RandomDate(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}