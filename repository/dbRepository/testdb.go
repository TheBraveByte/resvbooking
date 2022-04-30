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

//SearchRoomAvailabileByRoomID testing to check for all available with a certaion period ogf time
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

//SearchForAvailableRoom Testing to search for all available room in the database within certain date
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

//GetRoom Testing to get the correct room with it id from the database
func (tpg *TestPostgresDBRepository) GetRooms(room_id int) (models.Room, error) {
	var room models.Room

	if room_id > 4 {
		return room, errors.New("cannot get any rooms")
	}

	return room, nil

}

//GetUserInfoByID testing to get user details in the database
func (tpg *TestPostgresDBRepository) GetUserInfoByID(userID int) (models.User, error) {
	var user models.User
	return user, nil
}

//UpdateUserInfo tesing for updating user information in the database
func (tpg *TestPostgresDBRepository) UpdateUserInfo(user models.User) error {
	return nil
}

//AuthenticateUser testing for authenticated user with the database function
func (tpg *TestPostgresDBRepository) AuthenticateUser(testPassword, email string) (int, string, error) {
	var userID int
	var hashedPassword string
	return userID, hashedPassword, nil
}

//AllReservation testing for the database function for all present reservation
func (tpg *TestPostgresDBRepository) AllReservation() ([]models.Reservation, error) {
	var allResv []models.Reservation
	return allResv, nil
}

//AllNewReservation testing for the database function for all new reservation
func (tpg *TestPostgresDBRepository) AllNewReservation() ([]models.Reservation, error) {
	var allNewResv []models.Reservation
	return allNewResv, nil
}

func (tpg *TestPostgresDBRepository) ShowUserReservation(id int) (models.Reservation, error) {
	var userResv models.Reservation
	return userResv, nil
}
func (tpg TestPostgresDBRepository) UpdateUserReservation(resv models.Reservation) error {
	return nil
}

func (tpg *TestPostgresDBRepository) ProcessedUpdateReservation(id int) error {
	return nil
}
