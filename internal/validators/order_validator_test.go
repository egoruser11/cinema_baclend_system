package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Вспомогательная функция для создания mock DB
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

// Вспомогательная функция для создания echo.Context с user_data
func setupEchoContext(user *models.User) echo.Context {
	e := echo.New()
	c := e.NewContext(nil, nil)
	c.Set("user_data", user)
	return c
}

func float64Ptr(value float64) *float64 {
	return &value
}

func TestValidatePaidOrder(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(sqlmock.Sqlmock)
		user           *models.User
		req            requests.OrderPaidRequest
		expectedErrors map[string]string
		expectedOk     bool
	}{
		{
			name: "✅ успешная оплата только деньгами",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 1000.0,
				CoinBalance:  100,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   nil,
			},
			expectedErrors: map[string]string{},
			expectedOk:     true,
		},
		{
			name: "✅ успешная оплата деньгами + коины",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)

				coinsRows := sqlmock.NewRows([]string{"coins"}).AddRow(200)
				mock.ExpectQuery(`^SELECT "coins" FROM "users" WHERE id = \$1 AND "users"."deleted_at" IS NULL$`).
					WithArgs(1).
					WillReturnRows(coinsRows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 400.0,
				CoinBalance:  200,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   float64Ptr(100.0),
			},
			expectedErrors: map[string]string{},
			expectedOk:     true,
		},
		{
			name: "❌ заказ просрочен (больше 30 минут)",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now().Add(-40*time.Minute), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)

				mock.ExpectExec(`^DELETE FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			user: &models.User{
				ID: 1,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   nil,
			},
			expectedErrors: map[string]string{
				"orderExpired": "Order is too old, create new",
			},
			expectedOk: false,
		},
		{
			name: "❌ заказ не найден",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			user: &models.User{
				ID: 1,
			},
			req: requests.OrderPaidRequest{
				OrderID: 999,
				Coins:   nil,
			},
			expectedErrors: map[string]string{
				"order": "Order not found",
			},
			expectedOk: false,
		},
		{
			name: "❌ заказ принадлежит другому пользователю",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 2, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 10000,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   nil,
			},
			expectedErrors: map[string]string{
				"order": "User id does not match , dont do this",
			},
			expectedOk: false,
		},
		{
			name: "❌ недостаточно коинов",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)

				coinsRows := sqlmock.NewRows([]string{"coins"}).AddRow(50)
				mock.ExpectQuery(`^SELECT "coins" FROM "users" WHERE id = \$1 AND "users"."deleted_at" IS NULL$`).
					WithArgs(1).
					WillReturnRows(coinsRows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 1000.0,
				CoinBalance:  50,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   float64Ptr(100.0),
			},
			expectedErrors: map[string]string{
				"coins": "Coins 50 < 100 KINGSLAYEEER",
			},
			expectedOk: false,
		},
		{
			name: "❌ коинов больше чем стоимость заказа",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)

				coinsRows := sqlmock.NewRows([]string{"coins"}).AddRow(1000)
				mock.ExpectQuery(`^SELECT "coins" FROM "users" WHERE id = \$1 AND "users"."deleted_at" IS NULL$`).
					WithArgs(1).
					WillReturnRows(coinsRows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 1000.0,
				CoinBalance:  1000,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   float64Ptr(100.0),
			},
			expectedErrors: map[string]string{
				"coins": "Ypu input more then total amount of this order",
			},
			expectedOk: false,
		},
		{
			name: "❌ недостаточно денег на балансе",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "premiere_id", "seats",
					"total_amount", "status", "coins", "created_at", "updated_at",
				}).AddRow(
					1, 1, 10, "1-1,1-2", 500.0, "pending", 0.0,
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`^SELECT \* FROM "orders" WHERE id = \$1$`).
					WithArgs(1).
					WillReturnRows(rows)

				coinsRows := sqlmock.NewRows([]string{"coins"}).AddRow(50)
				mock.ExpectQuery(`^SELECT "coins" FROM "users" WHERE id = \$1 AND "users"."deleted_at" IS NULL$`).
					WithArgs(1).
					WillReturnRows(coinsRows)
			},
			user: &models.User{
				ID:           1,
				MoneyBalance: 100.0,
				CoinBalance:  50,
			},
			req: requests.OrderPaidRequest{
				OrderID: 1,
				Coins:   float64Ptr(100.0),
			},
			expectedErrors: map[string]string{
				"money": "Money is not enough",
			},
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupMockDB(t)
			tt.setupMocks(mock)
			c := setupEchoContext(tt.user)

			errors, ok := ValidatePaidOrder(c, gormDB, tt.req)

			assert.Equal(t, tt.expectedOk, ok)
			assert.Equal(t, tt.expectedErrors, errors)

			err := mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
