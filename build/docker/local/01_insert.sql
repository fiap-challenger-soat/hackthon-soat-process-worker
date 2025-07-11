-- Populando tb_user
INSERT INTO tb_user (id, email, created_at)
SELECT '3fa85f64-5717-4562-b3fc-2c963f66afa6', 'usuario@example.com', NOW()
WHERE NOT EXISTS (SELECT 1 FROM tb_user);

-- Populando tb_video_jobs
INSERT INTO tb_video_jobs (id, user_id, status, video_path, output_path)
VALUES 
  ('194f2506-3a19-42fb-91a0-50442a1bfcfd', '3fa85f64-5717-4562-b3fc-2c963f66afa6', 'queued', 'uploads/video-inicial.mp4', NULL);


INSERT INTO tb_video_jobs (id, user_id, status, video_path, output_path)
VALUES 
  ('33daf232-990b-4411-8146-c5cd7c2e5c86', '3fa85f64-5717-4562-b3fc-2c963f66afa6', 'queued', 'uploads/foto.mp4', NULL);
