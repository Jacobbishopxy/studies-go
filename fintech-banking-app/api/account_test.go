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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mockdb "fintech-banking-app/db/mock"
	db "fintech-banking-app/db/sqlc"
	"fintech-banking-app/util"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	// 为了达成 100% 覆盖率的测试，我们需要声明一个数组的测试用例。
	// 每个测试用例都有一个唯一的名称作为与其它用例的区分。
	testCases := []struct {
		name      string
		accountID int64
		// 以 mock store 作为入参，用于构建 stub 满足每个测试用例
		buildStubs func(store *mockdb.MockStore)
		// 测试 API 的输出
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		// 将每个测试案例作为单元测试的单独子测试运行
		t.Run(tc.name, func(t *testing.T) {
			// ctrl 作为 `mockdb.NewMockStore` 的入参
			// 由 `gomock.NewController` 生成，函数入参为 `testing.T` 对象
			ctrl := gomock.NewController(t)
			// defer 非常重要因为它将检查所有预期被调用的方法是否被调用了
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// 创建 server。
			// 不需要启动一个真实的 HTTP 服务，我们只需要用 `httptest` 包来记录 API 请求的响应。
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// 声明我们希望调用的 API 的 url 路径
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

// 为了让测试更加的健壮，我们也需要检查响应的 body。
// 响应的 body 是存储于 `recorder.Body` 字段中的，实际上就是一个 `bytes.Buffer` 的指针。
// 该函数有三个入参：`testing.T`，响应的 body，以及需要用于比较的 account 对象
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// 读取所有响应 body 的数据并存储于 data 变量中。需要其返回为无错。
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// 声明变量用于存储响应 body 的数据
	var gotAccount db.Account
	// unmarshal 数据成为 gotAccount 对象，需要无错以及与 `account` 相等
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
