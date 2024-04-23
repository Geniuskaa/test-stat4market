SELECT eventType, count(1) as count FROM test.events GROUP BY eventType HAVING count > 1000;

SELECT * FROM test.events WHERE toStartOfMonth(eventTime) = toDate(eventTime);

WITH ranked_messages AS (SELECT eventType, userID, row_number() OVER (PARTITION BY
    (userID)) as rn FROM test.events GROUP BY eventType, userID) SELECT
    userID FROM ranked_messages WHERE rn > 3;
