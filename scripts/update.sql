UPDATE entries
SET in_live = FALSE
WHERE visible_for = (SELECT id FROM entry_privacy WHERE type = 'followers');
