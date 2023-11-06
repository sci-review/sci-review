package service

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/repo"
)

type OrganizationService struct {
	OrganizationRepo *repo.OrganizationRepo
}

func NewOrganizationService(organizationRepo *repo.OrganizationRepo) *OrganizationService {
	return &OrganizationService{OrganizationRepo: organizationRepo}
}

func (os *OrganizationService) Create(data form.OrganizationCreateForm, userId uuid.UUID) (*model.Organization, error) {
	organization := model.NewOrganization(data.Name, data.Description)

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

	member := model.NewMember(userId, organization.Id, model.MemberOwner, true)
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

func (os *OrganizationService) List(id uuid.UUID) ([]model.Organization, error) {
	organizations, err := os.OrganizationRepo.FindAllByUserId(id)
	if err != nil {
		slog.Error("organization list", "error", err.Error())
		return nil, common.DbInternalError
	}

	slog.Info("organization list", "result", "success", "organizations", organizations)

	return organizations, nil
}

func (os *OrganizationService) Get(id uuid.UUID, userId uuid.UUID) (*model.Organization, error) {
	organization, err := os.OrganizationRepo.GetById(id)
	if err != nil {
		slog.Error("organization get", "error", err.Error())
		return nil, common.DbInternalError
	}

	if !organization.IsActiveMember(userId) {
		slog.Error("organization get", "error", "user is not a active member of the organization")
		return nil, common.ForbiddenError
	}

	slog.Info("organization get", "result", "success", "organization", organization)

	return organization, nil
}

func (os *OrganizationService) Archive(id uuid.UUID, userId uuid.UUID) error {
	organization, err := os.OrganizationRepo.GetById(id)
	if err != nil {
		slog.Error("organization archive", "error", err.Error())
		return common.DbInternalError
	}

	if !organization.IsOwner(userId) {
		slog.Error("organization archive", "error", "user is not a owner of the organization")
		return common.ForbiddenError
	}

	if err := os.OrganizationRepo.Archive(id); err != nil {
		slog.Error("organization archive", "error", err.Error())
		return common.DbInternalError
	}

	slog.Info("organization archive", "result", "success", "organization", organization)

	return nil
}
