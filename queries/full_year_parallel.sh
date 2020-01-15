curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"1.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"2.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"3.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"4.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"5.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"6.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"7.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"8.*"}[366d]))' &
curl -G http://localhost:8428/api/v1/query -d "time=$(date -d '2020-01-01 00:00:00 UTC' +%s)" -d 'query=max(max_over_time(temperature{sensor_id=~"9.*"}[366d]))' &

wait
