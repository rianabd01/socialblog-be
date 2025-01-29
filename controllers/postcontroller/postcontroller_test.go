package postcontroller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestIndex(t *testing.T) {
	// Membuat mock SQL dan database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Membuat mock gorm.DB dari sql.DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening gorm DB: %v", err)
	}
	models.DB = gormDB // Menggunakan DB yang sudah di-mock

	// Menyiapkan mock data yang akan dikembalikan oleh DB
	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "First Post", "This is the first post").
		AddRow(2, "Second Post", "This is the second post")

	// Menentukan bahwa kita mengharapkan query untuk mengambil posts dengan menggunakan regex
	mock.ExpectQuery(`SELECT \* FROM "posts"`).WillReturnRows(rows)

	// Set up Gin router
	router := gin.Default()
	router.GET("/posts", Index) // Menambahkan route yang menggunakan Index handler

	// Membuat request untuk memanggil endpoint
	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	// Jalankan handler
	router.ServeHTTP(w, req)

	// Pastikan response statusnya OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifikasi bahwa mock DB query telah dipanggil dengan benar
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %s", err)
	}

	// Verifikasi isi response body
	expected := `{"posts":[{"id":1,"title":"First Post","body":"This is the first post"},{"id":2,"title":"Second Post","body":"This is the second post"}]}`
	assert.JSONEq(t, expected, w.Body.String())
}
