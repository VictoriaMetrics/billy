curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"[0-9]{1,5}|1000000"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"1[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"2[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"3[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"4[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"5[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"6[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"7[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"8[0-9]{5}"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"9[0-9]{5}"}[366d]))' &

wait
