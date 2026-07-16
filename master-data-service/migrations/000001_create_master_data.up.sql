-- ============================================
-- Master Data Service - Initial Schema
-- ============================================

-- ── Academic Years ─────────────────────────────
CREATE TABLE IF NOT EXISTS `academic_years` (
  `id`         VARCHAR(255) NOT NULL,
  `semester`   ENUM('ganjil','genap') NOT NULL DEFAULT 'ganjil',
  `year`       VARCHAR(20) DEFAULT NULL,
  `start_date` DATETIME DEFAULT NULL,
  `end_date`   DATETIME DEFAULT NULL,
  `is_active`  TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ── Rooms ──────────────────────────────────────
CREATE TABLE IF NOT EXISTS `rooms` (
  `id`         VARCHAR(255) NOT NULL,
  `name`       VARCHAR(255) NOT NULL,
  `location`   VARCHAR(255) DEFAULT NULL,
  `capacity`   INT DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ── Science Groups (Kelompok Keahlian) ────────
CREATE TABLE IF NOT EXISTS `science_groups` (
  `id`         VARCHAR(255) NOT NULL,
  `name`       VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ── Thesis Topics ─────────────────────────────
CREATE TABLE IF NOT EXISTS `thesis_topics` (
  `id`         VARCHAR(255) NOT NULL,
  `name`       VARCHAR(255) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ── Thesis Status ─────────────────────────────
CREATE TABLE IF NOT EXISTS `thesis_status` (
  `id`   VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Seed Data
-- ============================================

-- Seed Academic Years
INSERT INTO `academic_years` (`id`, `semester`, `year`, `start_date`, `end_date`, `is_active`) VALUES
  (UUID(), 'ganjil', '2024', '2024-08-01', '2025-01-31', 0),
  (UUID(), 'genap',  '2024', '2025-02-01', '2025-07-31', 0),
  (UUID(), 'ganjil', '2025', '2025-08-01', '2026-01-31', 1);

-- Seed Thesis Statuses
INSERT INTO `thesis_status` (`id`, `name`) VALUES
  (UUID(), 'Pengajuan Judul'),
  (UUID(), 'Bimbingan'),
  (UUID(), 'Acc Seminar'),
  (UUID(), 'Seminar Proposal'),
  (UUID(), 'Revisi Seminar'),
  (UUID(), 'Sidang'),
  (UUID(), 'Revisi Sidang'),
  (UUID(), 'Selesai'),
  (UUID(), 'Gagal');
