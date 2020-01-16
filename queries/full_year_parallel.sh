curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"[1-3].*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"[4-6].*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"[7-9].*"}[366d]))' &

wait
