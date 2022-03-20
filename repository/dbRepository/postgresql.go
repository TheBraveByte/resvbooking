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
func (pg *PostgresDBRepository) InsertReservation(resv models.Reservation) (int, error) {

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	stmt := `insert into reservation (first_name, last_name, email, phone_number, check_in_date, 
                         check_out_date, room_id, created_at, updated_at) 
              values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

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

//InsertRoomRestriction  the helps to store rooms that are already reservated
func (pg *PostgresDBRepository) InsertRoomRestriction(resv models.RoomRestriction) error {
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

//SearchRoomAvailabile To check if a certain room is available within or at certain
//period of time..... if the function return true (that means the room is available) if
//false (the room is not available )
func (pg *PostgresDBRepository) SearchRoomAvailabileByRoomID(roomID int, checkInDate, checkOutDate time.Time) (bool, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	var rowCount int
	queryStmt := `select count(id) from room_restriction where room_id = $1 and $2 < check_out_date and $3 > check_in_date`
	row := pg.DB.QueryRowContext(ctx, queryStmt, roomID, checkInDate, checkOutDate)
	err := row.Scan(&rowCount)
	if err != nil {
		return false, err
	}
	if rowCount == 0 {
		return true, nil
	}
	return false, nil
}

func (pg *PostgresDBRepository) SearchForAvailableRoom(checkInDate, checkOutDate time.Time) ([]models.Room, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelCtx()
	var rooms []models.Room

	queryStmt := `select r.id, r.room_name from rooms r
                  where r.id not in (select rr.room_id from room_restriction rr
                  where $1 < rr.check_out_date and $2 > rr.check_in_date);`

	rows, err := pg.DB.QueryContext(ctx, queryStmt, checkInDate, checkOutDate)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil

}
