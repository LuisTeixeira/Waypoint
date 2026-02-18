CREATE INDEX idx_active_entity_realization 
ON activity_realizations (entity_id, family_id) 
WHERE status = 'in_progress';