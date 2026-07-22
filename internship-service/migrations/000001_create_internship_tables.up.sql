CREATE TABLE IF NOT EXISTS companies (
    id VARCHAR(36) PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    company_address TEXT NOT NULL,
    alasan TEXT,
    status ENUM('save', 'blacklist', 'diajukan') DEFAULT 'save',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_proposals (
    id VARCHAR(36) PRIMARY KEY,
    coordinator_id VARCHAR(36) NOT NULL,
    proposal_document_id VARCHAR(36) NOT NULL,
    academic_year_id VARCHAR(36) NOT NULL,
    target_company_id VARCHAR(36) NOT NULL,
    status ENUM('PENDING', 'APPROVED_PROPOSAL', 'REJECTED_PROPOSAL', 'WAITING_FOR_VERIFICATION', 'ACCEPTED_BY_COMPANY', 'PARTIALLY_ACCEPTED', 'REJECTED_BY_COMPANY') DEFAULT 'PENDING',
    proposal_sekdep_notes TEXT,
    proposed_start_date DATE NOT NULL,
    proposed_end_date DATE NOT NULL,
    
    app_letter_doc_number VARCHAR(255) UNIQUE,
    app_letter_date_issued DATE,
    start_date_planned DATE,
    end_date_planned DATE,
    app_letter_doc_id VARCHAR(36),
    app_letter_signed_by_id VARCHAR(36),
    app_letter_signed_as_role_id VARCHAR(36),
    
    company_response_doc_id VARCHAR(36),
    company_response_notes TEXT,
    
    assign_letter_doc_number VARCHAR(255) UNIQUE,
    assign_letter_date_issued DATE,
    start_date_actual DATE,
    end_date_actual DATE,
    assign_letter_doc_id VARCHAR(36),
    assign_letter_signed_by_id VARCHAR(36),
    assign_letter_signed_as_role_id VARCHAR(36),
    
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (target_company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_supervisor_letters (
    id VARCHAR(36) PRIMARY KEY,
    document_number VARCHAR(255) UNIQUE NOT NULL,
    date_issued DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    supervisor_id VARCHAR(36) NOT NULL,
    document_id VARCHAR(36),
    signed_by_id VARCHAR(36),
    signed_as_role_id VARCHAR(36),
    status ENUM('ACTIVE', 'SUPERSEDED') DEFAULT 'ACTIVE',
    superseded_at DATETIME(3),
    superseded_reason TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internships (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL,
    proposal_id VARCHAR(36) NOT NULL,
    supervisor_id VARCHAR(36),
    field_supervisor_name VARCHAR(255),
    field_supervisor_email VARCHAR(255),
    field_supervisor_phone VARCHAR(255),
    field_supervisor_nip VARCHAR(255),
    unit_section VARCHAR(255),
    actual_start_date DATE,
    actual_end_date DATE,
    status ENUM('PENDING', 'ACCEPTED', 'REJECTED', 'ACCEPTED_BY_COMPANY', 'REJECTED_BY_COMPANY', 'ONGOING', 'COMPLETED', 'FAILED') DEFAULT 'PENDING',
    is_logbook_locked BOOLEAN DEFAULT FALSE,
    logbook_locked_at DATETIME(3),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    sup_letter_id VARCHAR(36),
    
    report_title VARCHAR(255),
    report_document_id VARCHAR(36),
    report_status ENUM('SUBMITTED', 'APPROVED', 'REVISION_NEEDED') DEFAULT NULL,
    report_notes TEXT,
    report_uploaded_at DATETIME(3),
    report_feedback_document_id VARCHAR(36),
    
    lecturer_assessment_status ENUM('PENDING', 'APPROVED', 'COMPLETED') DEFAULT NULL,
    field_assessment_status ENUM('PENDING', 'APPROVED', 'COMPLETED') DEFAULT NULL,
    field_assessment_notes TEXT,
    field_assessment_doc_id VARCHAR(36),
    completion_certificate_doc_id VARCHAR(36),
    completion_certificate_status ENUM('SUBMITTED', 'APPROVED', 'REVISION_NEEDED') DEFAULT NULL,
    completion_certificate_notes TEXT,
    company_receipt_doc_id VARCHAR(36),
    company_receipt_status ENUM('SUBMITTED', 'APPROVED', 'REVISION_NEEDED') DEFAULT NULL,
    company_receipt_notes TEXT,
    logbook_document_id VARCHAR(36),
    logbook_document_status ENUM('SUBMITTED', 'APPROVED', 'REVISION_NEEDED') DEFAULT NULL,
    logbook_document_notes TEXT,
    company_report_doc_id VARCHAR(36),
    company_report_status ENUM('SUBMITTED', 'APPROVED', 'REVISION_NEEDED') DEFAULT NULL,
    company_report_notes TEXT,
    
    field_assessment_submitted_at DATETIME(3),
    field_assessment_signature_hash VARCHAR(255),
    logbook_field_signature_hash VARCHAR(255),
    logbook_field_signed_at DATETIME(3),
    
    final_numeric_score DOUBLE,
    final_grade VARCHAR(10),

    FOREIGN KEY (proposal_id) REFERENCES internship_proposals(id) ON DELETE CASCADE,
    FOREIGN KEY (sup_letter_id) REFERENCES internship_supervisor_letters(id) ON DELETE SET NULL,
    UNIQUE INDEX idx_student_proposal (student_id, proposal_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_logbooks (
    id VARCHAR(36) PRIMARY KEY,
    internship_id VARCHAR(36) NOT NULL,
    activity_date DATE NOT NULL,
    activity_description TEXT NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_seminars (
    id VARCHAR(36) PRIMARY KEY,
    internship_id VARCHAR(36) NOT NULL,
    room_id VARCHAR(36) NOT NULL,
    seminar_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    link_meeting VARCHAR(255),
    moderator_student_id VARCHAR(36) NOT NULL,
    status ENUM('REQUESTED', 'APPROVED', 'REJECTED', 'COMPLETED', 'FAILED') DEFAULT 'REQUESTED',
    approved_by VARCHAR(36),
    supervisor_notes TEXT,
    berita_acara_document_id VARCHAR(36),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_holidays (
    id VARCHAR(36) PRIMARY KEY,
    holiday_date DATE NOT NULL UNIQUE,
    name VARCHAR(255),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS supervisor_replacement_requests (
    id VARCHAR(36) PRIMARY KEY,
    letter_id VARCHAR(36) NOT NULL,
    internship_id VARCHAR(36) NOT NULL,
    old_supervisor_id VARCHAR(36) NOT NULL,
    new_supervisor_id VARCHAR(36) NOT NULL,
    reason TEXT NOT NULL,
    status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING',
    requested_by_id VARCHAR(36) NOT NULL,
    approved_by_id VARCHAR(36),
    requested_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    resolved_at DATETIME(3),
    rejection_notes TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE,
    FOREIGN KEY (letter_id) REFERENCES internship_supervisor_letters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
