// Create constraints for Story and Scene nodes
CREATE CONSTRAINT story_id_unique IF NOT EXISTS FOR (s:Story) REQUIRE s.id IS UNIQUE;
CREATE CONSTRAINT scene_id_unique IF NOT EXISTS FOR (s:Scene) REQUIRE s.id IS UNIQUE;

// Create indexes for common queries
CREATE INDEX story_status_idx IF NOT EXISTS FOR (s:Story) ON (s.status);
CREATE INDEX story_created_at_idx IF NOT EXISTS FOR (s:Story) ON (s.created_at);
