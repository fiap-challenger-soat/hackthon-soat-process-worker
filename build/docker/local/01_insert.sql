-- Populando tb_user
INSERT INTO tb_user (id, email, password_hash)
VALUES 
  ('51b00f47-5293-4100-8a24-0c5349b5ff47', 'nicolascastanha2002@gmail.com', '19f45c84bc9877f6bdd7b1c3c6c124355283126f');

-- Populando tb_video_jobs
INSERT INTO tb_video_jobs (id, user_id, status, video_path, output_path)
VALUES 
  ('194f2506-3a19-42fb-91a0-50442a1bfcfd', '51b00f47-5293-4100-8a24-0c5349b5ff47', 'queued', 'uploads/video-inicial.mp4', NULL);


INSERT INTO tb_video_jobs (id, user_id, status, video_path, output_path)
VALUES 
  ('33daf232-990b-4411-8146-c5cd7c2e5c86', '51b00f47-5293-4100-8a24-0c5349b5ff47', 'queued', 'uploads/foto.mp4', NULL);
