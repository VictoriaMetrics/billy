curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"1[0-9]{5}"}[366d]))'
