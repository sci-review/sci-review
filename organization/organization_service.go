package organization

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
)

type OrganizationService struct {
	OrganizationRepo *OrganizationRepo
}

func NewOrganizationService(organizationRepo *OrganizationRepo) *OrganizationService {
	return &OrganizationService{OrganizationRepo: organizationRepo}
}

func (os *OrganizationService) Create(data OrganizationCreateForm, userId uuid.UUID) (*Organization, error) {
	organization := NewOrganization(data.Name, data.Description)

	tx := os.OrganizationRepo.DB.MustBegin()

	if err := os.OrganizationRepo.Create(organization, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("organization create", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("organization create", "error", err.Error(), "organizationData", data)
		return nil, common.DbInternalError
	}

	member := NewMember(userId, organization.Id, Owner, true)
	if err := os.OrganizationRepo.AddMember(member, tx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.Error("organization create", "error", "error rolling back transaction", "data", data)
			return nil, err
		}
		slog.Error("organization create", "error", err.Error(), "organizationData", data)
		return nil, common.DbInternalError
	}

	if err := tx.Commit(); err != nil {
		slog.Error("login", "error", "error committing transaction", "data", data)
		return nil, err
	}

	slog.Info("organization create", "result", "success", "organization", organization)

	return organization, nil
}
