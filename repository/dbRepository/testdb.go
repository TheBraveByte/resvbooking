package dbRepository

import (
	_ "context"
	"log"
	"time"

	"github.com/dev-ayaa/resvbooking/pkg/models"
	"github.com/pkg/errors"
)

func (tpg *TestPostgresDBRepository) AllUser() bool {
	return true
}

//InsertReservation Insert a Reservation data into the database
func (tpg *TestPostgresDBRepository) InsertReservation(resv models.Reservation) (int, error) {
	if resv.RoomID == 14 {
		return 0, errors.New("can't insert room id")
	}
	return 1, nil
}

//InsertRoomRestriction  the helps to store rooms that are already reservated
func (tpg *TestPostgresDBRepository) InsertRoomRestriction(resv models.RoomRestriction) error {
	if resv.RoomID == 11 {
		return errors.New("Failed to insert room restriction")
	}
	return nil
}

func (tpg *TestPostgresDBRepository) SearchRoomAvailabileByRoomID(roomID int, checkInDate, checkOutDate time.Time) (bool, error) {
	dateLayout := "2006-01-02"
	checkIn := "2025-09-09"
	testDate, err := time.Parse(dateLayout, checkIn)
	if err != nil {
		log.Println(err)

	}

	errdate, err := time.Parse(dateLayout, "2045-09-09")
	if err != nil {
		log.Println(err)

	}

	if errdate == checkInDate {
		return false, errors.New("Room not available")
	}

	if checkInDate.After(testDate) {
		return false, nil
	}

	return true, nil

	// cid := checkInDate.Format(dateLayout)
	// cod := checkOutDate.Format(dateLayout)
	// if cid == "2029-09-09" {
	// 	return false, errors.New("Room not available")
	// }
	// if cod == "2029-09-10" {
	// 	return false, errors.New("Room not available")
	// }
	// return true, nil
}

func (tpg *TestPostgresDBRepository) SearchForAvailableRoom(checkInDate, checkOutDate time.Time) ([]models.Room, error) {

	var rooms []models.Room

	dateLayout := "2006-01-02"
	date := "2025-09-09"

	testDate, err := time.Parse(dateLayout, date)
	if err != nil {
		log.Println(err)
	}

	errdate, err := time.Parse(dateLayout, "2045-09-09")
	if err != nil {
		log.Println(err)
	}

	if checkInDate == errdate {
		return rooms, errors.New("invalid date for reservation")
	}

	if checkInDate.After(testDate) {
		return rooms, nil
	}

	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)
	return rooms, nil

}

func (tpg *TestPostgresDBRepository) GetRooms(room_id int) (models.Room, error) {
	var room models.Room

	if room_id > 4 {
		return room, errors.New("cannot get any rooms")
	}

	return room, nil

}
