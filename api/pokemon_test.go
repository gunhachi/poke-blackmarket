package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/gunhachi/poke-blackmarket/db/mock"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/gunhachi/poke-blackmarket/token"
	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonAPI(t *testing.T) {
	poke := mockRandomPoke()

	testCases := []struct {
		name          string
		pokeID        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "Succes_GetPokemon_API_nil_error",
			pokeID: poke.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, gomock.Any().String(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonData(gomock.Any(), gomock.Eq(poke.ID)).
					Times(1).
					Return(poke, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyPoke(t, recorder.Body, poke)
			},
		},
		{
			name:   "NotFound_GetPokemon_API_with_error",
			pokeID: poke.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, gomock.Any().String(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonData(gomock.Any(), gomock.Eq(poke.ID)).
					Times(1).
					Return(db.PokeProduct{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError_GetPokemon_API_with_error",
			pokeID: poke.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, gomock.Any().String(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonData(gomock.Any(), gomock.Eq(poke.ID)).
					Times(1).
					Return(db.PokeProduct{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID_GetPokemon_API_with_error",
			pokeID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, gomock.Any().String(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonData(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/pokemon/%d", tc.pokeID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.route.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreatePokemonAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Succes_CreatePokemon_API_nil_error",
			body: gin.H{
				"poke_name":  "Alala",
				"status":     "good",
				"poke_price": 2000,
				"poke_stock": 2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreatePokemonDataParams{
					PokeName:  "Alala",
					Status:    "good",
					PokePrice: 2000,
					PokeStock: 2,
				}

				store.EXPECT().
					CreatePokemonData(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.PokeProduct{
						PokeName:  "Alala",
						Status:    "good",
						PokePrice: 2000,
						PokeStock: 2,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyPoke(t, recorder.Body, db.PokeProduct{
					PokeName:  "Alala",
					Status:    "good",
					PokePrice: 2000,
					PokeStock: 2,
				})
			},
		},
		{
			name: "InternalError_CreatePokemon_API_with_error",
			body: gin.H{
				"poke_name":  "Alala",
				"status":     "good",
				"poke_price": 2000,
				"poke_stock": 2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreatePokemonData(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.PokeProduct{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name: "InvalidParam_CreatePokemon_API_nil_error",
			body: gin.H{
				"poke_name":  "Alala",
				"poke_stock": 2,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreatePokemonData(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/pokemon"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.route.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListPokemonAPI(t *testing.T) {
	n := 5
	pokes := make([]db.PokeProduct, n)
	for i := 0; i < n; i++ {
		pokes[i] = mockRandomPoke()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Succes_ListPokemon_API_nil_error",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListPokemonDataParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListPokemonData(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(pokes, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyPokeList(t, recorder.Body, pokes)
			},
		},
		{
			name: "InternalError_ListPokemon_API_with_error",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonData(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.PokeProduct{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID_ListPokemon_API_with_error",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonData(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize_ListPokemon_API_with_error",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonData(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/pokemon"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.route.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestUpdatePokemonAPI(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		pokeID        int64
// 		body          gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:   "Succes_UpdatePokemon_API_nil_error",
// 			pokeID: 11,
// 			body: gin.H{
// 				"status":     "good",
// 				"poke_price": 2000,
// 				"poke_stock": 1,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.UpdatePokemonDataParams{
// 					ID:        11,
// 					Status:    "good",
// 					PokePrice: 2000,
// 					PokeStock: 1,
// 				}

// 				store.EXPECT().
// 					UpdatePokemonData(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(db.PokeProduct{
// 						ID:        11,
// 						PokeName:  "Alala",
// 						Status:    "good",
// 						PokePrice: 2000,
// 						PokeStock: 1,
// 					}, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				reqBodyPoke(t, recorder.Body, db.PokeProduct{
// 					PokeName:  "Alala",
// 					Status:    "good",
// 					PokePrice: 2000,
// 					PokeStock: 2,
// 				})
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := fmt.Sprintf("/pokemon/%d", tc.pokeID)
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.route.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

//
// Note to add for update test unit
//

// mockRandomPoke create random user data
func mockRandomPoke() db.PokeProduct {
	return db.PokeProduct{
		ID:        util.RandomInt(1, 200),
		PokeName:  util.RandomString(8),
		Status:    util.RandomString(5),
		PokePrice: util.RandomAmount(),
		PokeStock: util.RandomInt(1, 15),
	}
}

// reqBodyPoke to check the response given on test
func reqBodyPoke(t *testing.T, body *bytes.Buffer, poke db.PokeProduct) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotData db.PokeProduct
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)
	require.Equal(t, poke, gotData)
}

// reqBodyPokeList check the response in form of list target data
func reqBodyPokeList(t *testing.T, body *bytes.Buffer, pokes []db.PokeProduct) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotData []db.PokeProduct
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)
	require.Equal(t, pokes, gotData)
}
