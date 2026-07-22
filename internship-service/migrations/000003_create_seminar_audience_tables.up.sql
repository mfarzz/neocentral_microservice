CREATE TABLE IF NOT EXISTS internship_seminar_audiences (
    id VARCHAR(36) PRIMARY KEY,
    seminar_id VARCHAR(36) NOT NULL,
    student_id VARCHAR(36) NOT NULL,
    status ENUM('PENDING', 'VALIDATED', 'REJECTED') DEFAULT 'PENDING',
    validated_by VARCHAR(36),
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (seminar_id) REFERENCES internship_seminars(id) ON DELETE CASCADE,
    UNIQUE INDEX idx_seminar_student (seminar_id, student_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
