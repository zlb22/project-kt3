# Nginx templates for project-kt3

This setup follows: single public domain for App+API via path routing, and a separate domain for MinIO (for stable presigned direct uploads).

- App (Vue, topic-three): `https://__APP_DOMAIN__/`
- API (FastAPI): `https://__APP_DOMAIN__/api/`
- MinIO: `https://__MINIO_DOMAIN__/`

Replace placeholders in `kt3.conf`:
- `__APP_DOMAIN__` → your public domain for the app (e.g., topic-three.example.com)
- `__MINIO_DOMAIN__` → your public domain for MinIO (e.g., minio.topic-three.example.com)
- `__A_LAN_IP__` → the LAN IP of Machine A (where backend and MinIO are running)

## Deploy steps (Machine B)

1) DNS: point `__APP_DOMAIN__` and `__MINIO_DOMAIN__` to Machine B's public IP.

2) TLS:
   - Issue certificates via certbot (example):
     ```bash
     sudo certbot certonly --nginx -d __APP_DOMAIN__ -d __MINIO_DOMAIN__
     ```

3) Upload static files (built artifacts):
   - Vue app (topic-three): copy `sub_project/sub_frontend/dist/` to Machine B, e.g. `/var/www/topic-three`.

4) Install config:
   - Copy `kt3.conf` to `/etc/nginx/conf.d/kt3.conf` (or sites-available + symlink)
   - Replace placeholders in that file with your actual values
   - Test and reload Nginx:
     ```bash
     sudo nginx -t && sudo systemctl reload nginx
     ```

5) Backend and MinIO upstreams (Machine A):
   - Ensure FastAPI listens on `0.0.0.0:8000` and MinIO on `0.0.0.0:9000`
   - Firewall: allow ONLY Machine B's IP to access A:8000 and A:9000

## Backend .env tips (Machine A)

- Set public MinIO URL so presigned URLs point to the public domain:
  ```
  MINIO_URL=https://__MINIO_DOMAIN__
  ```
- Restrict CORS to your app domain:
  ```
  CORS_ORIGINS=https://__APP_DOMAIN__
  ```
- Keti3 redirect target:
  ```
  KETI3_FRONTEND_URL=https://__APP_DOMAIN__
  ```

## Quick verification

- API health:
  ```bash
  curl -k https://__APP_DOMAIN__/api/health
  ```

- MinIO reachability:
  ```bash
  curl -I -k https://__MINIO_DOMAIN__/
  ```

If uploads fail in browser, check:
- Browser can directly reach `https://__MINIO_DOMAIN__/`
- Backend `MINIO_URL` matches that exact public domain (scheme + host)
- Nginx has `client_max_body_size 0` and `proxy_request_buffering off` for MinIO
