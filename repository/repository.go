package repository

import (
	"github.com/dev-ayaa/resvbooking/pkg/models"
	"time"
)

type DatabaseRepository interface {
	AllUser() bool
	InsertReservation(resv models.Reservation) (int, error)
	InsertRoomRestriction(resv models.RoomRestriction) error
	SearchRoomAvailabileByRoomID(roomID int, checkInDate, checkOutDate time.Time) (bool, error)
	SearchForAvailableRoom(checkInDate, checkOutDate time.Time) ([]models.Room, error)
	GetRooms(room_id int) (models.Room, error)

	GetUserInfoByID(user_id int) (models.User, error)
	UpdateUserInfo(user models.User) error
	AuthenticateUser(typedPassword, email string) (int, string, error)
}
