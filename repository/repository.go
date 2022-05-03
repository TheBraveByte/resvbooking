package repository

import (
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"time"
)

type DatabaseRepository interface {
	AllRoom() ([]models.Room, error)
	InsertReservation(resv models.Reservation) (int, error)
	InsertRoomRestriction(resv models.RoomRestriction) error
	SearchRoomAvailabileByRoomID(roomID int, checkInDate, checkOutDate time.Time) (bool, error)
	SearchForAvailableRoom(checkInDate, checkOutDate time.Time) ([]models.Room, error)
	GetRooms(room_id int) (models.Room, error)

	//Users
	GetUserInfoByID(user_id int) (models.User, error)
	UpdateUserInfo(user models.User) error
	AuthenticateUser(typedPassword, email string) (int, string, error)

	//Admin page
	AllReservation() ([]models.Reservation, error)
	AllNewReservation() ([]models.Reservation, error)
	ShowUserReservation(id int) (models.Reservation, error)
	UpdateUserReservation(resv models.Reservation) error
	ProcessedUpdateReservation(id int, processed int) error
	DeleteUserReservation(id int) error

	GetRestrictionsForRoomByDate(roomID int, checkInDate, checkOutDate time.Time) ([]models.RoomRestriction, error)
	DeleteBlockByID(id int) error
	InsertBlockForRoom(id int, checkInDate time.Time) error
}
