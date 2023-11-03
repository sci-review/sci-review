package organization

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type OrganizationRepo struct {
	DB *sqlx.DB
}

func NewOrganizationRepo(DB *sqlx.DB) *OrganizationRepo {
	return &OrganizationRepo{DB: DB}
}

type OrgAndMember struct {
	OrgId           uuid.UUID  `db:"org_id"`
	OrgName         string     `db:"org_name"`
	OrgDesc         string     `db:"org_description"`
	OrgArchived     bool       `db:"org_archived"`
	OrgCreatedAt    time.Time  `db:"org_created_at"`
	OrgUpdatedAt    time.Time  `db:"org_updated_at"`
	MemberId        uuid.UUID  `db:"member_id"`
	MemberUserId    uuid.UUID  `db:"member_user_id"`
	MemberRole      MemberRole `db:"member_role"`
	MemberActive    bool       `db:"member_active"`
	MemberCreatedAt time.Time  `db:"member_created_at"`
	MemberUpdatedAt time.Time  `db:"member_updated_at"`
}

func (or *OrganizationRepo) Create(organization *Organization, tx *sqlx.Tx) error {
	query := `
		INSERT INTO organizations (id, name, description, archived, created_at, updated_at)
		VALUES (:id, :name, :description, false, :created_at, :updated_at)
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

func (or *OrganizationRepo) Archive(id uuid.UUID) error {
	query := `UPDATE organizations SET archived = true WHERE id = $1`
	_, err := or.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (or *OrganizationRepo) FindAllByUserId(userId uuid.UUID) ([]Organization, error) {
	var organizationJoinMember []OrgAndMember
	query := `
		SELECT o.id AS org_id, o.name AS org_name, o.description AS org_description, o.created_at AS org_created_at,
		o.updated_at AS org_updated_at, o.archived AS org_archived, m.id AS member_id, m.user_id AS member_user_id, 
		m.role AS member_role, m.active AS member_active, m.created_at AS member_created_at, m.updated_at AS member_updated_at
		FROM organizations AS o
		LEFT JOIN members AS m ON o.id = m.organization_id
		WHERE o.id IN (SELECT organization_id FROM members WHERE user_id = $1 AND active = true); 
	`
	err := or.DB.Select(&organizationJoinMember, query, userId)
	if err != nil {
		return nil, err
	}
	return convertToOrganizations(organizationJoinMember), nil
}

func (or *OrganizationRepo) GetById(id uuid.UUID) (*Organization, error) {
	var orgAndMember []OrgAndMember
	query := `
		SELECT o.id AS org_id, o.name AS org_name, o.description AS org_description, o.created_at AS org_created_at,
		o.updated_at AS org_updated_at, o.archived AS org_archived, m.id AS member_id, m.user_id AS member_user_id, 
		m.role AS member_role, m.active AS member_active, m.created_at AS member_created_at, m.updated_at AS member_updated_at
		FROM organizations AS o
		LEFT JOIN members AS m ON o.id = m.organization_id
		WHERE o.id = $1; 
	`
	err := or.DB.Select(&orgAndMember, query, id)
	if err != nil {
		return nil, err
	}

	organizations := convertToOrganizations(orgAndMember)

	if len(organizations) == 0 {
		return nil, errors.New("organization not found")
	}

	return &organizations[0], nil
}

func convertToOrganizations(organizationMembers []OrgAndMember) []Organization {
	organizationMap := make(map[uuid.UUID]*Organization)
	for _, orgmember := range organizationMembers {
		if _, ok := organizationMap[orgmember.OrgId]; !ok {
			organization := &Organization{}
			organization.Id = orgmember.OrgId
			organization.Name = orgmember.OrgName
			organization.Description = orgmember.OrgDesc
			organization.Archived = orgmember.OrgArchived
			organization.CreatedAt = orgmember.OrgCreatedAt
			organization.UpdatedAt = orgmember.OrgUpdatedAt
			organization.Members = []Member{}
			organizationMap[orgmember.OrgId] = organization
		}
		organization := organizationMap[orgmember.OrgId]
		member := Member{
			Id:             orgmember.MemberId,
			UserId:         orgmember.MemberUserId,
			OrganizationId: orgmember.OrgId,
			Role:           orgmember.MemberRole,
			Active:         orgmember.MemberActive,
			CreatedAt:      orgmember.MemberCreatedAt,
			UpdatedAt:      orgmember.MemberUpdatedAt,
		}
		organization.AddMember(member)
	}

	var organizations []Organization
	for _, org := range organizationMap {
		organizations = append(organizations, *org)
	}

	return organizations
}
