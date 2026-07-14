package domain

import "time"

// ── Enums ────────────────────────────────────────────────────────

type IdentityType string

const (
	IdentityNIM   IdentityType = "NIM"
	IdentityNIP   IdentityType = "NIP"
	IdentityOTHER IdentityType = "OTHER"
)

type RoleStatus string

const (
	RoleStatusActive    RoleStatus = "active"
	RoleStatusNonActive RoleStatus = "nonActive"
)

type StudentStatus string

const (
	StudentActive          StudentStatus = "active"
	StudentDropout         StudentStatus = "dropout"
	StudentBSS             StudentStatus = "bss"
	StudentLulus           StudentStatus = "lulus"
	StudentMengundurkanDiri StudentStatus = "mengundurkan_diri"
)

// ── User ─────────────────────────────────────────────────────────

type User struct {
	ID                string       `gorm:"primaryKey;type:varchar(255)" json:"id"`
	FullName          string       `gorm:"column:full_name;not null" json:"fullName"`
	IdentityNumber    string       `gorm:"column:identity_number;uniqueIndex;not null" json:"identityNumber"`
	IdentityType      IdentityType `gorm:"column:identity_type;type:enum('NIM','NIP','OTHER');not null" json:"identityType"`
	Email             *string      `gorm:"uniqueIndex" json:"email"`
	Password          *string      `json:"-"`
	PhoneNumber       *string      `gorm:"column:phone_number" json:"phoneNumber"`
	IsVerified        bool         `gorm:"column:is_verified;default:false" json:"isVerified"`
	Token             *string      `gorm:"type:text" json:"-"`
	RefreshToken      *string      `gorm:"column:refresh_token;type:text" json:"-"`
	OAuthProvider     *string      `gorm:"column:oauth_provider" json:"-"`
	OAuthID           *string      `gorm:"column:oauth_id" json:"-"`
	OAuthRefreshToken *string      `gorm:"column:oauth_refresh_token;type:text" json:"-"`
	AvatarURL         *string      `gorm:"column:avatar_url" json:"avatarUrl"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`

	// Relations
	UserHasRoles []UserHasRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
	Student      *Student      `gorm:"foreignKey:ID" json:"student,omitempty"`
	Lecturer     *Lecturer     `gorm:"foreignKey:ID" json:"lecturer,omitempty"`
}

func (User) TableName() string { return "users" }

// ── Role ─────────────────────────────────────────────────────────

type UserRole struct {
	ID   string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

func (UserRole) TableName() string { return "user_roles" }

type UserHasRole struct {
	UserID string     `gorm:"column:user_id;primaryKey" json:"userId"`
	RoleID string     `gorm:"column:role_id;primaryKey" json:"roleId"`
	Status RoleStatus `gorm:"type:enum('active','nonActive');not null" json:"status"`

	Role UserRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserHasRole) TableName() string { return "user_has_roles" }

// ── Student ──────────────────────────────────────────────────────

type Student struct {
	ID                        string        `gorm:"column:user_id;primaryKey" json:"id"`
	Status                    StudentStatus `gorm:"column:student_status;type:enum('dropout','bss','lulus','mengundurkan_diri','active');default:active" json:"status"`
	EnrollmentYear            *int          `gorm:"column:enrollment_year" json:"enrollmentYear"`
	SKSCompleted              int           `gorm:"column:skscompleted" json:"sksCompleted"`
	MandatoryCoursesCompleted bool          `gorm:"column:mandatory_courses_completed;default:false" json:"mandatoryCoursesCompleted"`
	MKWUCompleted             bool          `gorm:"column:mkwu_completed;default:false" json:"mkwuCompleted"`
	InternshipCompleted       bool          `gorm:"column:internship_completed;default:false" json:"internshipCompleted"`
	KKNCompleted              bool          `gorm:"column:kkn_completed;default:false" json:"kknCompleted"`
	ResearchMethodCompleted   bool          `gorm:"column:research_method_completed;default:false" json:"researchMethodCompleted"`
	CurrentSemester           *int          `gorm:"column:current_semester" json:"currentSemester"`
	CreatedAt                 time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt                 time.Time     `gorm:"column:updated_at" json:"updatedAt"`
}

func (Student) TableName() string { return "students" }

// ── Lecturer ─────────────────────────────────────────────────────

type Lecturer struct {
	ID             string    `gorm:"column:user_id;primaryKey" json:"id"`
	ScienceGroupID *string   `gorm:"column:science_group_id" json:"scienceGroupId"`
	Data           *string   `gorm:"type:json" json:"data"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Lecturer) TableName() string { return "lecturers" }

// ── Well-known role names ────────────────────────────────────────

const (
	RoleKetuaDepartemen     = "Ketua Departemen"
	RoleSekretarisDepartemen = "Sekretaris Departemen"
	RolePembimbing1         = "Pembimbing 1"
	RolePembimbing2         = "Pembimbing 2"
	RoleAdmin               = "Admin"
	RolePenguji             = "Penguji"
	RoleMahasiswa           = "Mahasiswa"
	RoleGKM                 = "GKM"
	RoleKoordinatorYudisium = "Koordinator Yudisium"
)
