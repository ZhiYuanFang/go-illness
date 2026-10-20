-- 库：illness
-- 教学完成标记（OpenSpec: add-tutorial-phase-protocol）
ALTER TABLE `user`
  ADD COLUMN `tutorial_completed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '教学是否已完成' AFTER `use_live_day`;
