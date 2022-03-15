package dbRepository

import (
	"context"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"time"
)

func (pg *PostgresDBRepository) AllUser() bool {
	return true
}

//InsertReservation Insert a Reservation data into the database
func (pg PostgresDBRepository) InsertReservation(resv models.Reservation) error {

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	stmt := `insert into reservation (first_name, last_name,email,phone_number, check_in_data, check_out_date,
             room_id,created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := pg.DB.ExecContext(ctx, stmt,
		resv.FirstName,
		resv.LastName,
		resv.Email,
		resv.PhoneNumber,
		resv.CheckInDate,
		resv.CheckOutDate,
		resv.RoomID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
