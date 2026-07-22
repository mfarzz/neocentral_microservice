ALTER TABLE `documents` DROP FOREIGN KEY `fk_document_type`;
DROP TABLE IF EXISTS `documents`;
DROP TABLE IF EXISTS `document_templates`;
DROP TABLE IF EXISTS `document_types`;
