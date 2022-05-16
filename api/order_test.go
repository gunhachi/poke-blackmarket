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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/gunhachi/poke-blackmarket/db/mock"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func TestGetOrderAPI(t *testing.T) {
	order := mockRandomOrder()

	testCases := []struct {
		name          string
		ID            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Succes_GetOrder_API_nil_error",
			ID:   order.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonOrderData(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(order, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyOrder(t, recorder.Body, order)
			},
		},
		{
			name: "NotFound_GetOrder_API_with_error",
			ID:   order.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonOrderData(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(db.PokeOrder{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError_GetOrder_API_with_error",
			ID:   order.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonOrderData(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(db.PokeOrder{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID_GetOrder_API_with_error",
			ID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPokemonOrderData(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/order/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.route.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateOrderAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Succes_CreateOrder_API_nil_error",
			body: gin.H{
				"user_id":      11,
				"product_id":   12,
				"quantity":     2,
				"total_price":  4000,
				"order_detail": "selling",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.InsertPokemonOrderDataParams{
					UserID:      11,
					ProductID:   12,
					Quantity:    2,
					TotalPrice:  4000,
					OrderDetail: "selling",
				}

				store.EXPECT().
					InsertPokemonOrderData(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.PokeOrder{
						UserID:      11,
						ProductID:   12,
						Quantity:    2,
						TotalPrice:  4000,
						OrderDetail: "selling",
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyOrder(t, recorder.Body, db.PokeOrder{
					UserID:      11,
					ProductID:   12,
					Quantity:    2,
					TotalPrice:  4000,
					OrderDetail: "selling",
				})
			},
		},
		{
			name: "InternalError_CreateOrder_API_with_error",
			body: gin.H{
				"user_id":      11,
				"product_id":   12,
				"quantity":     2,
				"total_price":  4000,
				"order_detail": "selling",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertPokemonOrderData(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.PokeOrder{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name: "InvalidParamRole_CreateOrder_API_nil_error",
			body: gin.H{
				"user_id":      "11",
				"product_id":   "12",
				"quantity":     2,
				"total_price":  4000,
				"order_detail": "selling",
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					InsertPokemonOrderData(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/order"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.route.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListOrderAccountAPI(t *testing.T) {
	n := 5
	orders := make([]db.PokeOrder, n)
	for i := 0; i < n; i++ {
		orders[i] = mockRandomOrder()
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
			name: "Succes_ListOrder_API_nil_error",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListPokemonOrderDataParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListPokemonOrderData(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orders, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				reqBodyOrders(t, recorder.Body, orders)
			},
		},
		{
			name: "InternalError_ListOrder_API_with_error",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonOrderData(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.PokeOrder{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID_ListOrder_API_with_error",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonOrderData(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize_ListOrder_API_with_error",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListPokemonOrderData(gomock.Any(), gomock.Any()).
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/order"
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

//
// Note to add for update and delete test unit
//

// mockRandomOrder create random user data
func mockRandomOrder() db.PokeOrder {
	return db.PokeOrder{
		ID:          util.RandomInt(1, 200),
		UserID:      util.RandomInt(1, 300),
		ProductID:   util.RandomInt(1, 300),
		Quantity:    int32(util.RandomInt(1, 10)),
		TotalPrice:  util.RandomAmount(),
		OrderDetail: "selling",
	}
}

// reqBodyOrder to check the response given on test
func reqBodyOrder(t *testing.T, body *bytes.Buffer, param db.PokeOrder) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotData db.PokeOrder
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)
	require.Equal(t, param, gotData)
}

// reqBodyOrders check the response in form of list target data
func reqBodyOrders(t *testing.T, body *bytes.Buffer, params []db.PokeOrder) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotData []db.PokeOrder
	err = json.Unmarshal(data, &gotData)
	require.NoError(t, err)
	require.Equal(t, params, gotData)
}
