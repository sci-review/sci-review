package organization

import "github.com/jmoiron/sqlx"

type OrganizationRepo struct {
	DB *sqlx.DB
}

func NewOrganizationRepo(DB *sqlx.DB) *OrganizationRepo {
	return &OrganizationRepo{DB: DB}
}

func (or *OrganizationRepo) Create(organization *Organization, tx *sqlx.Tx) error {
	query := `
		INSERT INTO organizations (id, name, description, created_at, updated_at)
		VALUES (:id, :name, :description, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, organization)
	if err != nil {
		return err
	}
	return nil
}

func (or *OrganizationRepo) AddMember(member *Member, tx *sqlx.Tx) error {
	query := `
		INSERT INTO members (id, user_id, organization_id, role, active, created_at, updated_at)
		VALUES (:id, :user_id, :organization_id, :role, :active, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, member)
	if err != nil {
		return err
	}
	return nil
}
