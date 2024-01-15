-- ALTER TABLE sessions
--   ADD COLUMN session_place_id INTEGER DEFAULT 0,
--   ADD CONSTRAINT ak_sessions_session_place_id
--     UNIQUE(session_place_id)
-- ;

ALTER TABLE session_places
  ADD CONSTRAINT ak_session_places_session_id
    UNIQUE(session_id)
;

ALTER TABLE session_places_online
  ADD CONSTRAINT ak_session_places_online_session_place_id
    UNIQUE(session_place_id)
;

ALTER TABLE session_places_given_venue
  ADD CONSTRAINT ak_session_places_given_venue_session_place_id
    UNIQUE(session_place_id)
;

ALTER TABLE session_places_member_home
  ADD CONSTRAINT ak_session_places_member_home_session_place_id
    UNIQUE(session_place_id)
;
