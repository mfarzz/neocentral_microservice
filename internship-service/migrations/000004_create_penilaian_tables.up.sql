CREATE TABLE IF NOT EXISTS internship_cpmks (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name TEXT NOT NULL,
    weight DOUBLE NOT NULL,
    assessor_type VARCHAR(20) NOT NULL,
    academic_year_id VARCHAR(36) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_assessment_rubrics (
    id VARCHAR(36) PRIMARY KEY,
    cpmk_id VARCHAR(36) NOT NULL,
    level_name VARCHAR(255) NOT NULL,
    rubric_level_description TEXT NOT NULL,
    min_score DOUBLE NOT NULL,
    max_score DOUBLE NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (cpmk_id) REFERENCES internship_cpmks(id) ON DELETE CASCADE,
    UNIQUE INDEX idx_cpmk_min_score (cpmk_id, min_score),
    UNIQUE INDEX idx_cpmk_max_score (cpmk_id, max_score)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_lecturer_scores (
    internship_id VARCHAR(36) NOT NULL,
    chosen_rubric_id VARCHAR(36) NOT NULL,
    score DOUBLE NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE,
    FOREIGN KEY (chosen_rubric_id) REFERENCES internship_assessment_rubrics(id) ON DELETE CASCADE,
    PRIMARY KEY (internship_id, chosen_rubric_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS internship_field_scores (
    internship_id VARCHAR(36) NOT NULL,
    chosen_rubric_id VARCHAR(36) NOT NULL,
    score DOUBLE NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (internship_id) REFERENCES internships(id) ON DELETE CASCADE,
    FOREIGN KEY (chosen_rubric_id) REFERENCES internship_assessment_rubrics(id) ON DELETE CASCADE,
    PRIMARY KEY (internship_id, chosen_rubric_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
