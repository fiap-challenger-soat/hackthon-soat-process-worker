INSERT INTO users (id, email, password_hash)
VALUES ('51b00f47-5293-4100-8a24-0c5349b5ff47', 'pedrinho@email.com', '19f45c84bc9877f6bdd7b1c3c6c124355283126f');


INSERT INTO video_jobs (id, user_id, status, video_path, output_path)
VALUES ('194f2506-3a19-42fb-91a0-50442a1bfcfd', '51b00f47-5293-4100-8a24-0c5349b5ff47', 'queued', 'uploads/video-inicial.mp4', NULL);


INSERT INTO job_status_history (id, job_id, status)
VALUES ('a38e12c0-c36d-4ef3-9c3f-3027d86091b3', '194f2506-3a19-42fb-91a0-50442a1bfcfd', 'queued');
