-- auth-service/migrations/000001_create_users.up.sql
-- Creates the core user and role tables for auth-service.

CREATE TABLE IF NOT EXISTS users (
    id                  VARCHAR(255) PRIMARY KEY,
    full_name           VARCHAR(255) NOT NULL,
    identity_number     VARCHAR(255) NOT NULL UNIQUE,
    identity_type       ENUM('NIM', 'NIP', 'OTHER') NOT NULL,
    email               VARCHAR(255) UNIQUE,
    password            VARCHAR(255),
    phone_number        VARCHAR(50),
    is_verified         BOOLEAN DEFAULT FALSE,
    token               TEXT,
    refresh_token       TEXT,
    oauth_provider      VARCHAR(100),
    oauth_id            VARCHAR(255),
    oauth_refresh_token TEXT,
    avatar_url          VARCHAR(500),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_users_email (email),
    INDEX idx_users_identity (identity_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS user_roles (
    id   VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS user_has_roles (
    user_id VARCHAR(255) NOT NULL,
    role_id VARCHAR(255) NOT NULL,
    status  ENUM('active', 'nonActive') NOT NULL DEFAULT 'active',

    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES user_roles(id) ON DELETE CASCADE,
    INDEX idx_uhr_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS students (
    user_id                     VARCHAR(255) PRIMARY KEY,
    student_status              ENUM('dropout','bss','lulus','mengundurkan_diri','active') DEFAULT 'active',
    enrollment_year             INT,
    skscompleted                INT NOT NULL DEFAULT 0,
    mandatory_courses_completed BOOLEAN DEFAULT FALSE,
    mkwu_completed              BOOLEAN DEFAULT FALSE,
    internship_completed        BOOLEAN DEFAULT FALSE,
    kkn_completed               BOOLEAN DEFAULT FALSE,
    research_method_completed   BOOLEAN DEFAULT FALSE,
    current_semester            INT,
    created_at                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS lecturers (
    user_id          VARCHAR(255) PRIMARY KEY,
    science_group_id VARCHAR(255),
    data             JSON,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Seed default roles
INSERT INTO user_roles (id, name) VALUES
    ('role-ketua-dept',      'Ketua Departemen'),
    ('role-sekretaris-dept', 'Sekretaris Departemen'),
    ('role-pembimbing1',     'Pembimbing 1'),
    ('role-pembimbing2',     'Pembimbing 2'),
    ('role-admin',           'Admin'),
    ('role-penguji',         'Penguji'),
    ('role-mahasiswa',       'Mahasiswa'),
    ('role-gkm',             'GKM'),
    ('role-koord-yudisium',  'Koordinator Yudisium')
ON DUPLICATE KEY UPDATE name = VALUES(name);


-- ── Seed Users ───────────────────────────────────────────────────
INSERT INTO users (id, full_name, identity_number, identity_type, email, password, is_verified) VALUES
    ('usr-kadep-001', 'Ricky Akbar M.Kom', '198410062012121001', 'NIP', 'kadep_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-sekdep-001', 'Afriyanti Dwi Kartika, M.T', '198904212019032024', 'NIP', 'sekdep_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-pembimbing-001', 'Husnil Kamil, MT', '198201182008121002', 'NIP', 'pembimbing_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-penguji-001', 'Aina Hubby Aziira, M.Eng', '199504302022032013', 'NIP', 'penguji_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-gkm-001', 'Ullya Mega Wahyuni, M.Kom', '199011032019032008', 'NIP', 'gkm_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-admin-001', 'Nindy Malisha, SE', '220199206201501201', 'OTHER', 'admin_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-yudisium-001', 'Koordinator Yudisium', '199203152020121003', 'NIP', 'yudisium_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-cpl-001', 'GKM User', '199107282019031005', 'NIP', 'cpl_si@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-fariz-034', 'Muhammad Fariz', '2211523034', 'NIM', 'fariz_2211523034@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-nabil-018', 'Nabil Rizki Navisa', '2211522018', 'NIM', 'nabil_2211522018@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-khalied-030', 'Khalied Nauly Maturino', '2211523030', 'NIM', 'khalied_2211523030@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-mustafa-036', 'Mustafa Fathur Rahman', '2211522036', 'NIM', 'mustafa_2211522036@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-muhammad-020', 'Muhammad Nouval Habibie', '2211521020', 'NIM', 'muhammad_2211521020@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-daffa-022', 'Daffa Agustian Saadi', '2211523022', 'NIM', 'daffa_2211523022@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-ilham-028', 'Ilham', '2211522028', 'NIM', 'ilham_2211522028@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-syauqi-012', 'Syauqi', '2211523012', 'NIM', 'syauqi_2211523012@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-dimas-026', 'Dimas', '2311523026', 'NIM', 'dimas_2311523026@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-john-001', 'John', '2411522001', 'NIM', 'john_2411522001@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-test-topic-001', 'Test Ganti Topik', '2211522101', 'NIM', 'test_changetopic@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-test-super-001', 'Test Ganti Dospem', '2211522102', 'NIM', 'test_changesupervisor@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE),
    ('usr-test-nothes-001', 'Test Tanpa Thesis', '2211522103', 'NIM', 'test_nothesis@fti.unand.ac.id', '$2a$10$MUFIMgHQX9c1njKExADjDOjtNScX/.vfeeCiit49wl3OnQlYGwJuC', TRUE)
ON DUPLICATE KEY UPDATE full_name = VALUES(full_name);


-- ── Seed Lecturers ───────────────────────────────────────────────
INSERT INTO lecturers (user_id) VALUES
    ('usr-kadep-001'),
    ('usr-sekdep-001'),
    ('usr-pembimbing-001'),
    ('usr-penguji-001'),
    ('usr-gkm-001'),
    ('usr-yudisium-001'),
    ('usr-cpl-001')
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);


-- ── Seed Students ────────────────────────────────────────────────
INSERT INTO students (user_id, student_status, enrollment_year, skscompleted) VALUES
    ('usr-fariz-034', 'active', 2022, 122),
    ('usr-nabil-018', 'active', 2022, 137),
    ('usr-khalied-030', 'active', 2022, 141),
    ('usr-mustafa-036', 'active', 2022, 137),
    ('usr-muhammad-020', 'active', 2022, 137),
    ('usr-daffa-022', 'active', 2022, 137),
    ('usr-ilham-028', 'active', 2022, 137),
    ('usr-syauqi-012', 'active', 2022, 125),
    ('usr-dimas-026', 'active', 2023, 99),
    ('usr-john-001', 'active', 2024, 60),
    ('usr-test-topic-001', 'active', 2022, 130),
    ('usr-test-super-001', 'active', 2022, 130),
    ('usr-test-nothes-001', 'active', 2022, 130)
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);


-- ── Seed User Roles Assignments ──────────────────────────────────
INSERT INTO user_has_roles (user_id, role_id, status) VALUES
    -- kadep: KETUA_DEPARTEMEN, PEMBIMBING_1, PEMBIMBING_2, PENGUJI
    ('usr-kadep-001', 'role-ketua-dept', 'active'),
    ('usr-kadep-001', 'role-pembimbing1', 'active'),
    ('usr-kadep-001', 'role-pembimbing2', 'active'),
    ('usr-kadep-001', 'role-penguji', 'active'),
    -- sekdep: SEKRETARIS_DEPARTEMEN, KOORDINATOR_YUDISIUM, PEMBIMBING_1, PEMBIMBING_2, PENGUJI
    ('usr-sekdep-001', 'role-sekretaris-dept', 'active'),
    ('usr-sekdep-001', 'role-koord-yudisium', 'active'),
    ('usr-sekdep-001', 'role-pembimbing1', 'active'),
    ('usr-sekdep-001', 'role-pembimbing2', 'active'),
    ('usr-sekdep-001', 'role-penguji', 'active'),
    -- pembimbing: PEMBIMBING_1, PEMBIMBING_2, PENGUJI
    ('usr-pembimbing-001', 'role-pembimbing1', 'active'),
    ('usr-pembimbing-001', 'role-pembimbing2', 'active'),
    ('usr-pembimbing-001', 'role-penguji', 'active'),
    -- penguji: PENGUJI, GKM, PEMBIMBING_1, PEMBIMBING_2
    ('usr-penguji-001', 'role-penguji', 'active'),
    ('usr-penguji-001', 'role-gkm', 'active'),
    ('usr-penguji-001', 'role-pembimbing1', 'active'),
    ('usr-penguji-001', 'role-pembimbing2', 'active'),
    -- gkm: GKM, PENGUJI, PEMBIMBING_2
    ('usr-gkm-001', 'role-gkm', 'active'),
    ('usr-gkm-001', 'role-penguji', 'active'),
    ('usr-gkm-001', 'role-pembimbing2', 'active'),
    -- admin: ADMIN
    ('usr-admin-001', 'role-admin', 'active'),
    -- yudisium: KOORDINATOR_YUDISIUM, PEMBIMBING_2, PENGUJI
    ('usr-yudisium-001', 'role-koord-yudisium', 'active'),
    ('usr-yudisium-001', 'role-pembimbing2', 'active'),
    ('usr-yudisium-001', 'role-penguji', 'active'),
    -- cpl: GKM, PEMBIMBING_2, PENGUJI
    ('usr-cpl-001', 'role-gkm', 'active'),
    ('usr-cpl-001', 'role-pembimbing2', 'active'),
    ('usr-cpl-001', 'role-penguji', 'active'),
    -- mahasiswa
    ('usr-fariz-034', 'role-mahasiswa', 'active'),
    ('usr-nabil-018', 'role-mahasiswa', 'active'),
    ('usr-khalied-030', 'role-mahasiswa', 'active'),
    ('usr-mustafa-036', 'role-mahasiswa', 'active'),
    ('usr-muhammad-020', 'role-mahasiswa', 'active'),
    ('usr-daffa-022', 'role-mahasiswa', 'active'),
    ('usr-ilham-028', 'role-mahasiswa', 'active'),
    ('usr-syauqi-012', 'role-mahasiswa', 'active'),
    ('usr-dimas-026', 'role-mahasiswa', 'active'),
    ('usr-john-001', 'role-mahasiswa', 'active'),
    ('usr-test-topic-001', 'role-mahasiswa', 'active'),
    ('usr-test-super-001', 'role-mahasiswa', 'active'),
    ('usr-test-nothes-001', 'role-mahasiswa', 'active')
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);

