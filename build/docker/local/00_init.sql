CREATE TABLE IF NOT EXISTS tb_user (
    id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE IF NOT EXISTS tb_video_jobs (
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES tb_user(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    video_path VARCHAR(255),
    output_path VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now()
);
