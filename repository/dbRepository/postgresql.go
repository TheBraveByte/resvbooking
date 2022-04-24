package dbRepository

import (
	"context"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

func (pg *PostgresDBRepository) GetRooms(room_id int) (models.Room, error) {
	var room models.Room
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()
	query := `select id, room_name , created_at, updated_at from rooms where id = $1`

	rooms := pg.DB.QueryRowContext(ctx, query, room_id)

	err := rooms.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil

}

func (pg PostgresDBRepository) GetUserInfoByID(userID int) (models.User, error) {
	var user models.User
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	query := `select id, first_name, last_name, email, password, created_at, updated_at, access_level from user where id = $1`
	row := pg.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.AccessLevel)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

//UpdateUserInfo to Update the users information or details in the database
func (pg PostgresDBRepository) UpdateUserInfo(user models.User) error {
	//var user models.User
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	query := `update users set first_name= $1, last_name=$2,email=$3,updated_at=$4, access_level=$5`
	_, err := pg.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.UpdatedAt, user.AccessLevel)
	if err != nil {
		return err
	}
	return nil
}

//AuthenticateUser to Athenticate the user by verifying the email and the Password
func (pg *PostgresDBRepository) AuthenticateUser(typedPassword, email string) (int, string, error) {
	var userID int
	var hashedPassword string
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()
	query := "select id, password from users where email= $1 "
	row := pg.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&userID, &hashedPassword)
	if err != nil {
		return userID, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(typedPassword))
	//If the password did not match
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password!")
	} else if err != nil {
		return 0, "", err
	}
	return userID, hashedPassword, nil

}

//DataBase Functions for the administration pages

//AllReservation this show all the registered resservations in the database
func (pg PostgresDBRepository) AllReservation() ([]models.Reservation, error) {
	var allResv []models.Reservation
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()
	query := `select r.id,
       r.first_name,
       r.last_name,
       r.email,
       r.phone_number,
       r.room_id,
       r.check_in_date,
       r.check_in_date,
       r.updated_at,
       r.created_at
       from reservation r
         left join rooms rm on (r.room_id = rm.id)
       order by r.check_in_date`
	row, err := pg.DB.QueryContext(ctx, query)
	if err != nil {
		return allResv, err
	}
	for row.Next() {
		var rs models.Reservation
		err = row.Scan(
			&rs.ID,
			&rs.FirstName,
			&rs.LastName,
			&rs.Email,
			&rs.CheckInDate,
			&rs.CheckOutDate,
			&rs.RoomID,
			&rs.UpdatedAt,
			&rs.CreatedAt,
			&rs.Room.RoomName,
			&rs.PhoneNumber,
		)
		if err != nil {
			return allResv, err
		}
		allResv = append(allResv, rs)
	}
	if err = row.Err(); err != nil {
		return allResv, err
	}
	return allResv, nil
}
