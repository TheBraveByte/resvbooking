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
func (pg PostgresDBRepository) InsertReservation(resv models.Reservation) (int, error) {

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	stmt := `insert into reservation (first_name, last_name, email, phone_number, check_in_date, check_out_date, room_id, created_at, updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
returning id`

	// to get the id of newly insert reservation we need to querythe database
	var NewID int
	err := pg.DB.QueryRowContext(ctx, stmt,

		resv.FirstName,
		resv.LastName,
		resv.Email,
		resv.PhoneNumber,
		resv.CheckInDate,
		resv.CheckOutDate,
		resv.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&NewID)
	if err != nil {
		return 0, err
	}

	return NewID, nil
}

func (pg PostgresDBRepository) InsertRoomRestriction(resv models.RoomRestriction) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	stmt := `insert into room_restriction (check_in_date, check_out_date, room_id, reservation_id,restriction_id,created_at,updated_at )
values ($1,$2,$3,$4,$5,$6,$7)`
	_, err := pg.DB.ExecContext(ctx, stmt,
		resv.CheckInDate,
		resv.CheckOutDate,
		resv.RoomID,
		resv.ReservationID,
		resv.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
