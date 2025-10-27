package repositories

import (
	"clean_architecture_fiber/data/db/generated"
	"context"
)

type RoleRepository interface {
	GetByValue(ctx context.Context, value string) (*generated.GetRoleByValueRow, error)
}

type roleRepository struct {
	query *generated.Queries
}

func NewRoleRepository(query *generated.Queries) RoleRepository {
	return &roleRepository{query: query}
}

func (r *roleRepository) GetByValue(ctx context.Context, value string) (*generated.GetRoleByValueRow, error) {
	roleSQLC, err := r.query.GetRoleByValue(ctx, value)
	if err != nil {
		return nil, err
	}
	return &roleSQLC, nil
}
