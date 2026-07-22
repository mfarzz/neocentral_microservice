package model

import "time"

type ProposalStatus string

const (
	ProposalPending             ProposalStatus = "PENDING"
	ProposalApproved            ProposalStatus = "APPROVED_PROPOSAL"
	ProposalRejected            ProposalStatus = "REJECTED_PROPOSAL"
	ProposalWaitingVerification ProposalStatus = "WAITING_FOR_VERIFICATION"
	ProposalCompanyAccepted     ProposalStatus = "ACCEPTED_BY_COMPANY"
	ProposalPartiallyAccepted   ProposalStatus = "PARTIALLY_ACCEPTED"
	ProposalCompanyRejected     ProposalStatus = "REJECTED_BY_COMPANY"
)

type InternshipProposal struct {
	ID                 string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CoordinatorID      string         `gorm:"type:varchar(36);not null" json:"coordinatorId"`
	ProposalDocumentID string         `gorm:"type:varchar(36);not null" json:"proposalDocumentId"`
	AcademicYearID     string         `gorm:"type:varchar(36);not null" json:"academicYearId"`
	TargetCompanyID    string         `gorm:"type:varchar(36);not null" json:"targetCompanyId"`
	Status             ProposalStatus `gorm:"type:enum('PENDING', 'APPROVED_PROPOSAL', 'REJECTED_PROPOSAL', 'WAITING_FOR_VERIFICATION', 'ACCEPTED_BY_COMPANY', 'PARTIALLY_ACCEPTED', 'REJECTED_BY_COMPANY');default:'PENDING'" json:"status"`

	ProposalSekdepNotes *string   `gorm:"type:text" json:"proposalSekdepNotes"`
	ProposedStartDate   time.Time `gorm:"type:date;not null" json:"proposedStartDate"`
	ProposedEndDate     time.Time `gorm:"type:date;not null" json:"proposedEndDate"`

	AppLetterDocNumber      *string    `gorm:"type:varchar(255);unique" json:"appLetterDocNumber"`
	AppLetterDateIssued     *time.Time `gorm:"type:date" json:"appLetterDateIssued"`
	StartDatePlanned        *time.Time `gorm:"type:date" json:"startDatePlanned"`
	EndDatePlanned          *time.Time `gorm:"type:date" json:"endDatePlanned"`
	AppLetterDocID          *string    `gorm:"type:varchar(36)" json:"appLetterDocId"`
	AppLetterSignedByID     *string    `gorm:"type:varchar(36)" json:"appLetterSignedById"`
	AppLetterSignedAsRoleID *string    `gorm:"type:varchar(36)" json:"appLetterSignedAsRoleId"`

	CompanyResponseDocID *string `gorm:"type:varchar(36)" json:"companyResponseDocId"`
	CompanyResponseNotes *string `gorm:"type:text" json:"companyResponseNotes"`

	AssignLetterDocNumber      *string    `gorm:"type:varchar(255);unique" json:"assignLetterDocNumber"`
	AssignLetterDateIssued     *time.Time `gorm:"type:date" json:"assignLetterDateIssued"`
	StartDateActual            *time.Time `gorm:"type:date" json:"startDateActual"`
	EndDateActual              *time.Time `gorm:"type:date" json:"endDateActual"`
	AssignLetterDocID          *string    `gorm:"type:varchar(36)" json:"assignLetterDocId"`
	AssignLetterSignedByID     *string    `gorm:"type:varchar(36)" json:"assignLetterSignedById"`
	AssignLetterSignedAsRoleID *string    `gorm:"type:varchar(36)" json:"assignLetterSignedAsRoleId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Company *Company `gorm:"foreignKey:TargetCompanyID" json:"company,omitempty"`
}

func (InternshipProposal) TableName() string {
	return "internship_proposals"
}
