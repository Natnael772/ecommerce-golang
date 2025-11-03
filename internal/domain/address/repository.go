package address

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/pkg/database"
	"ecommerce-app/internal/pkg/database/sqlc"
	"ecommerce-app/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context, userID string, addr CreateAddressRequest) (Address, error)
	GetByID(ctx context.Context, id string) (Address, error)
	GetByUserID(ctx context.Context, userID string) ([]Address, error)
	Update(ctx context.Context, id string, addr UpdateAddressRequest) (Address, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) Create(ctx context.Context, userID string, addr CreateAddressRequest) (Address, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return Address{}, err
	}

	params := sqlc.CreateAddressParams{
		UserID:     userUUID,
		Label:      pgtype.Text{String: addr.Label, Valid: true},
		Line1:      addr.Line1,
		Line2:      pgtype.Text{String: addr.Line2, Valid: addr.Line2 != ""},
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
		IsDefault:  pgtype.Bool{Bool: addr.IsDefault, Valid: true},
	}

	row,  err := r.q.CreateAddress(ctx, params)
	if err != nil {
		return Address{}, err
	}

	return mapAddress(row), nil
}

func (r *repository) GetByID(ctx context.Context, id string) (Address, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Address{}, err
	}

	row, err := r.q.GetAddress(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Address{}, errs.ErrNotFound
		}
		
		return Address{}, err
	}

	return mapAddress(row), nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) ([]Address, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return nil, err
	}

	rows,  err := r.q.GetAddressesByUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Address{}, errs.ErrNotFound
		}

		return nil, err
	}

	addresses := make([]Address, len(rows))
	for i, row := range rows {
		addresses[i] = mapAddress(row)
	}

	return addresses, nil
}

func (r *repository) Update(ctx context.Context, id string, addr UpdateAddressRequest) (Address, error){
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return Address{}, err
	}

	params := sqlc.UpdateAddressParams{
		ID: uuid,
	}

	// Update only provided fields
	if addr.Label != nil {
		params.Label = pgtype.Text{String: *addr.Label, Valid: true}
	}
	
	if addr.Line2 != nil {
		params.Line2 = pgtype.Text{String: *addr.Line2, Valid: true}
	}

	if addr.IsDefault != nil {
		params.IsDefault = pgtype.Bool{Bool: *addr.IsDefault, Valid: true}
	}
	
	database.AssignIfProvided(&params.Line1, addr.Line1)
	database.AssignIfProvided(&params.City, addr.City)
	database.AssignIfProvided(&params.State, addr.State)
	database.AssignIfProvided(&params.PostalCode, addr.PostalCode)
	database.AssignIfProvided(&params.Country, addr.Country)
	
	row, err := r.q.UpdateAddress(ctx, params)
	if err != nil {
		return Address{}, err
	}

	return mapAddress(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}

	if err := r.q.DeleteAddress(ctx, uuid); err != nil {
		return err
	}
	return nil
}

func mapAddress(row sqlc.Address) Address {
	return Address{
		ID:         uuid.UUID(row.ID.Bytes),
		UserID:     uuid.UUID(row.UserID.Bytes),
		Label:      row.Label.String,
		Line1:      row.Line1,
		Line2:      row.Line2.String,
		City:       row.City,
		State:      row.State,
		PostalCode: row.PostalCode,
		Country:    row.Country,
		IsDefault:  row.IsDefault.Bool,
		IsDeleted:  row.IsDeleted.Bool,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}
