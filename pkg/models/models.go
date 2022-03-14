package models

import "time"

//used to store models in the Database

//Reservation reservation data
type Reservation struct {
	ID              int
	FirstName       string
	LastName        string
	Email           string
	PhoneNumber     string
	Password        string
	ConfirmPassword string
	RoomID          int
	checkInDate     time.Time
	checkOutDate    time.Time
}

//Room rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//User user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Restriction restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Room Restriction model
type RoomRestriction struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	checkInDate   time.Time
	checkOutDate  time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
