package interview

const (
	getEventIDMapQuery = `SELECT id, event_id FROM interviews where event_id IS NOT NULL ORDER BY id`
	updateEventIDQuery = `UPDATE interviews SET event_id = ? WHERE event_id = ?`
)
