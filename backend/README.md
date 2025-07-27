### DB起動
```bash
docker-compose up -d
```

### DBデータ設定
```bash
docker-compose exec -T postgres psql -U postgres -d postgres < schema.sql
```