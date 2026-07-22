package model

import "time"

type InternshipStatus string

const (
	InternshipPending           InternshipStatus = "PENDING"
	InternshipAccepted          InternshipStatus = "ACCEPTED"
	InternshipRejected          InternshipStatus = "REJECTED"
	InternshipCompanyAccepted   InternshipStatus = "ACCEPTED_BY_COMPANY"
	InternshipCompanyRejected   InternshipStatus = "REJECTED_BY_COMPANY"
	InternshipOngoing           InternshipStatus = "ONGOING"
	InternshipCompleted         InternshipStatus = "COMPLETED"
	InternshipFailed            InternshipStatus = "FAILED"
)

type DocumentSubmitStatus string

const (
	DocSubmitted      DocumentSubmitStatus = "SUBMITTED"
	DocApproved       DocumentSubmitStatus = "APPROVED"
	DocRevisionNeeded DocumentSubmitStatus = "REVISION_NEEDED"
)

type AssessmentStatus string

const (
	AssessmentPending   AssessmentStatus = "PENDING"
	AssessmentApproved  AssessmentStatus = "APPROVED"
	AssessmentCompleted AssessmentStatus = "COMPLETED"
)

type Internship struct {
	ID                   string           `gorm:"primaryKey;type:varchar(36)" json:"id"`
	StudentID            string           `gorm:"type:varchar(36);not null" json:"studentId"`
	ProposalID           string           `gorm:"type:varchar(36);not null" json:"proposalId"`
	SupervisorID         *string          `gorm:"type:varchar(36)" json:"supervisorId"`
	FieldSupervisorName  *string          `gorm:"type:varchar(255)" json:"fieldSupervisorName"`
	FieldSupervisorEmail *string          `gorm:"type:varchar(255)" json:"fieldSupervisorEmail"`
	FieldSupervisorPhone *string          `gorm:"type:varchar(255)" json:"fieldSupervisorPhone"`
	FieldSupervisorNip   *string          `gorm:"type:varchar(255)" json:"fieldSupervisorNip"`
	UnitSection          *string          `gorm:"type:varchar(255)" json:"unitSection"`
	ActualStartDate      *time.Time       `gorm:"type:date" json:"actualStartDate"`
	ActualEndDate        *time.Time       `gorm:"type:date" json:"actualEndDate"`
	Status               InternshipStatus `gorm:"type:enum('PENDING', 'ACCEPTED', 'REJECTED', 'ACCEPTED_BY_COMPANY', 'REJECTED_BY_COMPANY', 'ONGOING', 'COMPLETED', 'FAILED');default:'PENDING'" json:"status"`
	IsLogbookLocked      bool             `gorm:"default:false" json:"isLogbookLocked"`
	LogbookLockedAt      *time.Time       `json:"logbookLockedAt"`

	SupLetterID *string `gorm:"type:varchar(36)" json:"supLetterId"`

	ReportTitle              *string               `gorm:"type:varchar(255)" json:"reportTitle"`
	ReportDocumentID         *string               `gorm:"type:varchar(36)" json:"reportDocumentId"`
	ReportStatus             *DocumentSubmitStatus `gorm:"type:enum('SUBMITTED', 'APPROVED', 'REVISION_NEEDED')" json:"reportStatus"`
	ReportNotes              *string               `gorm:"type:text" json:"reportNotes"`
	ReportUploadedAt         *time.Time            `json:"reportUploadedAt"`
	ReportFeedbackDocumentID *string               `gorm:"type:varchar(36)" json:"reportFeedbackDocumentId"`

	LecturerAssessmentStatus    *AssessmentStatus     `gorm:"type:enum('PENDING', 'APPROVED', 'COMPLETED')" json:"lecturerAssessmentStatus"`
	FieldAssessmentStatus       *AssessmentStatus     `gorm:"type:enum('PENDING', 'APPROVED', 'COMPLETED')" json:"fieldAssessmentStatus"`
	FieldAssessmentNotes        *string               `gorm:"type:text" json:"fieldAssessmentNotes"`
	FieldAssessmentDocID        *string               `gorm:"type:varchar(36)" json:"fieldAssessmentDocId"`
	CompletionCertificateDocID  *string               `gorm:"type:varchar(36)" json:"completionCertificateDocId"`
	CompletionCertificateStatus *DocumentSubmitStatus `gorm:"type:enum('SUBMITTED', 'APPROVED', 'REVISION_NEEDED')" json:"completionCertificateStatus"`
	CompletionCertificateNotes  *string               `gorm:"type:text" json:"completionCertificateNotes"`

	CompanyReceiptDocID  *string               `gorm:"type:varchar(36)" json:"companyReceiptDocId"`
	CompanyReceiptStatus *DocumentSubmitStatus `gorm:"type:enum('SUBMITTED', 'APPROVED', 'REVISION_NEEDED')" json:"companyReceiptStatus"`
	CompanyReceiptNotes  *string               `gorm:"type:text" json:"companyReceiptNotes"`

	LogbookDocumentID     *string               `gorm:"type:varchar(36)" json:"logbookDocumentId"`
	LogbookDocumentStatus *DocumentSubmitStatus `gorm:"type:enum('SUBMITTED', 'APPROVED', 'REVISION_NEEDED')" json:"logbookDocumentStatus"`
	LogbookDocumentNotes  *string               `gorm:"type:text" json:"logbookDocumentNotes"`

	CompanyReportDocID  *string               `gorm:"type:varchar(36)" json:"companyReportDocId"`
	CompanyReportStatus *DocumentSubmitStatus `gorm:"type:enum('SUBMITTED', 'APPROVED', 'REVISION_NEEDED')" json:"companyReportStatus"`
	CompanyReportNotes  *string               `gorm:"type:text" json:"companyReportNotes"`

	FieldAssessmentSubmittedAt   *time.Time `json:"fieldAssessmentSubmittedAt"`
	FieldAssessmentSignatureHash *string    `gorm:"type:varchar(255)" json:"fieldAssessmentSignatureHash"`
	LogbookFieldSignatureHash    *string    `gorm:"type:varchar(255)" json:"logbookFieldSignatureHash"`
	LogbookFieldSignedAt         *time.Time `json:"logbookFieldSignedAt"`

	FinalNumericScore *float64 `json:"finalNumericScore"`
	FinalGrade        *string  `gorm:"type:varchar(10)" json:"finalGrade"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Proposal *InternshipProposal `gorm:"foreignKey:ProposalID" json:"proposal,omitempty"`
}

func (Internship) TableName() string {
	return "internships"
}
