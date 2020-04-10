CREATE TABLE user_awards (
    award_id varchar(18) NOT NULL,
    user_id varchar(18) NOT NULL,
    reward bigint NOT NULL,
    paid boolean NOT NULL DEFAULT false,
    PRIMARY KEY (award_id, user_id)
);

