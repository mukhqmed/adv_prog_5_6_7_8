package filtering

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetBarbersFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "basic_info", "price", "experience", "status", "image_path"}).
		AddRow(1, "Бекболат", "Мастер Бекболат специализируется на классических и современных стрижках", 3500, "9 лет", "Expert", "https://mukhqmed.github.io/images/bekbolat.jpg").
		AddRow(2, "Ерасыл", "Ерасыл - наш эксперт по стрижке бороды и уходу за ней, поэтому вы будете выглядеть безупречно", 3500, "6 лет", "Senior", "https://mukhqmed.github.io/images/master2.jpeg").
		AddRow(3, "Иван", "Точная техника бритья мастера Ивана обеспечит вам комфортное и тщательное бритье", 3000, "5 лет", "Senior", "https://mukhqmed.github.io/images/master3.jpeg").
		AddRow(4, "Ерсин", "Ерсин - наш самый молодой специалист. Он всегда в курсе последних трендов и техник.", 1500, "1 год", "Junior", "https://mukhqmed.github.io/images/master4.jpeg")

	mock.ExpectQuery("SELECT id, name, basic_info, price, experience, status, image_path FROM barbers").WillReturnRows(rows)

	barbers, err := GetBarbersFromDB(db)
	if err != nil {
		t.Fatalf("error getting barbers from DB: %v", err)
	}

	expected := []Barber{
		{ID: 1, Name: "Бекболат", BasicInfo: "Мастер Бекболат специализируется на классических и современных стрижках", Price: 3500, Experience: "9 лет", Status: "Expert", ImagePath: "https://mukhqmed.github.io/images/bekbolat.jpg"},
		{ID: 2, Name: "Ерасыл", BasicInfo: "Ерасыл - наш эксперт по стрижке бороды и уходу за ней, поэтому вы будете выглядеть безупречно", Price: 3500, Experience: "6 лет", Status: "Senior", ImagePath: "https://mukhqmed.github.io/images/master2.jpeg"},
		{ID: 3, Name: "Иван", BasicInfo: "Точная техника бритья мастера Ивана обеспечит вам комфортное и тщательное бритье", Price: 3000, Experience: "5 лет", Status: "Senior", ImagePath: "https://mukhqmed.github.io/images/master3.jpeg"},
		{ID: 4, Name: "Ерсин", BasicInfo: "Ерсин - наш самый молодой специалист. Он всегда в курсе последних трендов и техник.", Price: 1500, Experience: "1 год", Status: "Junior", ImagePath: "https://mukhqmed.github.io/images/master4.jpeg"},
	}

	if !reflect.DeepEqual(barbers, expected) {
		t.Errorf("unexpected result, got: %v, want: %v", barbers, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
