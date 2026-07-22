CREATE TABLE IF NOT EXISTS internship_guidance_questions (
    id VARCHAR(36) PRIMARY KEY,
    week_number INT NOT NULL,
    question_text TEXT NOT NULL,
    order_index INT NOT NULL DEFAULT 0,
    academic_year_id VARCHAR(36) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_guidance_lecturer_criteria (
    id VARCHAR(36) PRIMARY KEY,
    criteria_name VARCHAR(255) NOT NULL,
    week_number INT NOT NULL,
    input_type ENUM('EVALUATION', 'TEXT') NOT NULL,
    order_index INT NOT NULL DEFAULT 0,
    academic_year_id VARCHAR(36) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_guidance_lecturer_criteria_options (
    id VARCHAR(36) PRIMARY KEY,
    criteria_id VARCHAR(36) NOT NULL,
    option_text VARCHAR(255) NOT NULL,
    order_index INT NOT NULL DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (criteria_id) REFERENCES internship_guidance_lecturer_criteria(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_guidance_sessions (
    id VARCHAR(36) PRIMARY KEY,
    internship_id VARCHAR(36) NOT NULL,
    week_number INT NOT NULL,
    status ENUM('SUBMITTED', 'LATE', 'APPROVED') DEFAULT 'SUBMITTED',
    submission_date DATE,
    approved_at DATETIME(3),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_guidance_student_answers (
    guidance_session_id VARCHAR(36) NOT NULL,
    question_id VARCHAR(36) NOT NULL,
    week_number INT NOT NULL,
    answer_text TEXT NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (guidance_session_id, question_id, week_number),
    FOREIGN KEY (guidance_session_id) REFERENCES internship_guidance_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES internship_guidance_questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_guidance_lecturer_answers (
    guidance_session_id VARCHAR(36) NOT NULL,
    criteria_id VARCHAR(36) NOT NULL,
    week_number INT NOT NULL,
    evaluation_value VARCHAR(255),
    answer_text TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (guidance_session_id, criteria_id, week_number),
    FOREIGN KEY (guidance_session_id) REFERENCES internship_guidance_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (criteria_id) REFERENCES internship_guidance_lecturer_criteria(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
