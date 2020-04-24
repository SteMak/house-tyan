CREATE TABLE rewards (
    award_id BIGINT NOT NULL,
    user_id varchar(18) NOT NULL,
    amount bigint NOT NULL,
    paid boolean NOT NULL DEFAULT FALSE
);

CREATE INDEX ids_rewards_award_id ON rewards (award_id);
CREATE INDEX ids_rewards_user_id ON rewards (user_id);
CREATE INDEX ids_rewards_award_id_user_id ON rewards (award_id, user_id);

