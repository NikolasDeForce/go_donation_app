package useranddonatesdata

import "donation/db"

type Data struct {
	UsersData   []db.User
	DonatesData []db.Donate
}
