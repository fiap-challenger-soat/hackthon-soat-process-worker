CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE IF NOT EXISTS video_jobs (
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    video_path VARCHAR(255),
    output_path VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE IF NOT EXISTS job_status_history (
    id uuid PRIMARY KEY,
    job_id uuid REFERENCES video_jobs(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
