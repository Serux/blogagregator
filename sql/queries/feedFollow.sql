
-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feeds_follows
    (
        id, 
        created_at, 
        updated_at,
        user_id,
        feed_id 
    )
    VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
inserted_feed_follow.*, 
users.name as username,
feeds.name as feedname 
FROM inserted_feed_follow 
LEFT JOIN users ON inserted_feed_follow.user_id = users.id
LEFT JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
;


-- name: GetFeedFollowsForUser :many

SELECT 
feeds_follows.*, 
feeds.name AS feedname, 
users.name AS username
FROM feeds_follows
LEFT JOIN feeds ON feeds_follows.feed_id = feeds.id
LEFT JOIN users ON feeds_follows.user_id = users.id
WHERE feeds_follows.user_id = $1;

-- name: DeleteFeedFollowsUserIDAndURL :exec

DELETE FROM feeds_follows
WHERE feeds_follows.user_id = $1 
AND  feeds_follows.feed_id = (
    SELECT id FROM feeds WHERE feeds.url = $2
);
