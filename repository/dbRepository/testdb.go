package dbRepository

import (
	_"context"
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"time"
)

func (tpg *TestPostgresDBRepository) AllUser() bool {
	return true
}

//InsertReservation Insert a Reservation data into the database
func (tpg *TestPostgresDBRepository) InsertReservation(resv models.Reservation) (int, error) {
	return 1, nil
}

//InsertRoomRestriction  the helps to store rooms that are already reservated
func (tpg *TestPostgresDBRepository) InsertRoomRestriction(resv models.RoomRestriction) error {
	
	return nil
}

func (tpg *TestPostgresDBRepository) SearchRoomAvailabileByRoomID(roomID int, checkInDate, checkOutDate time.Time) (bool, error) {
	
	return false, nil
}

func (tpg *TestPostgresDBRepository) SearchForAvailableRoom(checkInDate, checkOutDate time.Time) ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil

}

func (tpg *TestPostgresDBRepository) GetRooms(room_id int) (models.Room, error) {
	var room models.Room
	
	return room, nil

}
