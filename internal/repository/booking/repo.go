package repository

import (
	"context"
	booking "fewoserv/internal/domain/booking"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/logger"
	"fewoserv/pkg/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	collectionName       = "bookings"
	cachedRepo     *Repo = nil
	log                  = logger.New("REPO")
)

type (
	Repo struct {
		repository                *mongodb.Repository[booking.Booking]
		repositoryBookingOverview *mongodb.Repository[booking.BookingOverview]
	}
)

func New(dbClient mongodb.IClient) *Repo {
	if cachedRepo == nil {
		cachedRepo = &Repo{
			repository:                mongodb.NewRepository[booking.Booking](dbClient, collectionName),
			repositoryBookingOverview: mongodb.NewRepository[booking.BookingOverview](dbClient, collectionName),
		}
	}

	return cachedRepo
}

func (r *Repo) LoadByID(id string) (*booking.Booking, error) {
	ctx := context.Background()

	query := bson.M{"_id": id}
	return r.repository.FindOne(ctx, query)
}

func (r *Repo) ValidateRecordExists(booking *booking.Booking) bool {
	ctx := context.Background()

	query := bson.M{"_id": booking.ID}
	foundRecord, err := r.repository.FindOne(ctx, query)

	doesRecordExists := err == nil && foundRecord != nil
	return doesRecordExists
}

func (r *Repo) Insert(booking *booking.Booking) error {
	ctx := context.Background()

	_, err := r.repository.InsertOne(ctx, booking)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(updaterID string, booking *booking.Booking) error {
	ctx := context.Background()

	booking.Edited.Update(updaterID)
	return r.repository.UpdateOne(ctx, booking.ID, booking)
}

func (r *Repo) DeleteByID(id string) error {
	ctx := context.Background()

	return r.repository.DeleteByID(ctx, id)
}

func (r *Repo) FindMany(name string, sort *common.Sort, skip, limit int64) ([]*booking.Booking, error) {
	ctx := context.Background()

	query := bson.M{
		"$or": []bson.M{
			{"name.deDE": bson.M{"$regex": name, "$options": "i"}},
			{"name.enGB": bson.M{"$regex": name, "$options": "i"}},
			{"name.frFR": bson.M{"$regex": name, "$options": "i"}},
			{"name.itIT": bson.M{"$regex": name, "$options": "i"}},
		},
	}

	return r.repository.Find(ctx, &query, nil, sort.ToBson(), &skip, &limit)
}

func (r *Repo) FindManyBy(apartmentIDs []*string, date time.Time) ([]*booking.BookingOverview, error) {
	ctx := context.Background()

	matchConditions := bson.M{
		"fromDate": bson.M{
			"$gte": time.Date(date.Year(), date.Month()-6, date.Day(), 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(date.Year(), date.Month()+6, date.Day(), 0, 0, 0, 0, time.UTC),
		},
		// we ignore bookings with the status canceled
		"status": bson.M{"$ne": 3},
	}

	if apartmentIDs != nil {
		matchConditions["apartmentID"] = bson.M{"$in": apartmentIDs}
	}

	aggregate := []bson.M{
		{
			"$match": matchConditions,
		},
		{
			"$lookup": bson.M{
				"from":         "admin_users",
				"localField":   "userID",
				"foreignField": "_id",
				"as":           "adminUser",
			},
		},
		{
			"$unwind": bson.M{"path": "$adminUser", "preserveNullAndEmptyArrays": true},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userID",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true},
		},
		{
			"$addFields": bson.M{
				"userName": bson.M{
					"$concat": []interface{}{
						bson.M{"$ifNull": []interface{}{"$user.firstName", "$adminUser.firstName"}},
						" ",
						bson.M{"$ifNull": []interface{}{"$user.lastName", "$adminUser.lastName"}},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$apartmentID",
				"summaries": bson.M{
					"$push": bson.M{
						"bookingID":     "$_id",
						"bookingNumber": "$bookingNumber",
						"status":        "$status",
						"fromDate":      "$fromDate",
						"toDate":        "$toDate",
						"stayDays":      "$stayDays",
						"guestInfo":     "$guestInfo",
						"priceSummary":  bson.M{"total": "$priceSummary.total"},
						"messages":      "$messages",
						"userName":      "$userName",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":         0,
				"apartmentID": "$_id",
				"summaries":   1,
			},
		},
		{
			"$sort": bson.M{"apartmentName": -1},
		},
	}

	return r.repositoryBookingOverview.FindByAggregate(ctx, aggregate)
}

func (r *Repo) FindManyBetweenDate(apartmentID string, bookingID *string, fromDate, toDate time.Time) ([]string, error) {
	ctx := context.Background()

	usedFromDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 6, 0, 0, 0, time.UTC)
	usedToDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 6, 0, 0, 0, time.UTC)

	query := bson.M{
		"apartmentID": apartmentID,
		// we ignore canceled bookings
		"status": bson.M{"$ne": 3},

		"$or": []bson.M{
			{
				"fromDate": bson.M{
					"$lte": usedFromDate,
				},
				"toDate": bson.M{
					"$gt": usedFromDate,
				},
			},
			{
				"$and": []bson.M{
					{"fromDate": bson.M{"$gt": usedFromDate}},
					{"fromDate": bson.M{"$lt": usedToDate}},
				},
			},
			{
				"$and": []bson.M{
					{"toDate": bson.M{"$gt": usedFromDate}},
					{"toDate": bson.M{"$lt": usedToDate}},
				},
			},
		},
	}

	if bookingID != nil {
		query["_id"] = bson.M{
			"$ne": bookingID,
		}
	}

	///////

	projection := bson.M{
		"_id": 1,
	}

	result, err := r.repository.Find(ctx, &query, &projection, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	recordIDs := []string{}
	for _, record := range result {
		recordIDs = append(recordIDs, record.ID)
	}

	return recordIDs, nil
}

func (r *Repo) GetBlockedDates(apartmentID string, fromDate, toDate time.Time) ([]time.Time, error) {
	ctx := context.Background()

	query := bson.M{
		"apartmentID": apartmentID,
		"status":      bson.M{"$ne": 3},
		"fromDate":    bson.M{"$gte": fromDate, "$lte": toDate},
	}

	result, err := r.repository.Find(ctx, &query, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	blockedDates := []time.Time{}
	for _, record := range result {
		currentDate := record.FromDate
		endDate := record.ToDate

		for currentDate.Before(endDate) || currentDate.Equal(endDate) {
			blockedDates = append(blockedDates, currentDate)
			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	return blockedDates, nil
}
