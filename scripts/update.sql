UPDATE entries
SET visible_for = (SELECT id FROM entry_privacy WHERE type = 'followers')
WHERE visible_for = (SELECT id FROM entry_privacy WHERE type = 'some');
