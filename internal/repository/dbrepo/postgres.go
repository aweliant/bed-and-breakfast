package dbrepo

import (
	"context"
	"github.com/aweliant/bed-and-breakfast/internal/models"
	"log"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//
func (m *postgresDBRepo) InsertReservation(r models.Reservation) (int, error) {
	//why use context? Here because want to tell the function that it should abandon its work after the timeout elapses.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`
	err := m.DB.QueryRowContext(ctx, stmt,
		r.FirstName, r.LastName, r.Email, r.Phone, r.StartDate, r.EndDate, r.RoomID, time.Now(), time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, 
			created_at, updated_at, restriction_id)
			values ($1, $2, $3, $4, $5, $6, $7);`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate, r.EndDate, r.RoomID, r.ReservationID, time.Now(), time.Now(), r.RestrictionID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res := 0
	query := `select 
    		 	 count(id)
             from 
                 room_restrictions 
             where roomID=$1 and $2 > start_date and  $3 < end_date;`

	err := m.DB.QueryRowContext(ctx, query, roomID, end, start).Scan(res)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return res == 0, nil
}
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rooms := []models.Room{}
	query := `select 
    		 	 r.id, r.room_name
             from 
                 rooms r
             where r.id not in
            (select
                 rr.room_id
			from
			    room_restrictions rr
			where
			    $1 < rr.end_date and $2 > rr.start_date);`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		room := models.Room{}
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roomName string
	query := `
		select room_name from rooms where id =$1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&roomName,
	)
	if err != nil {
		return roomName, err
	}
	return roomName, nil
}
