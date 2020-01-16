curl -X POST http://localhost:8428/api/v1/import -d '{"metric":{"__name__":"temperature","sensor_id":"12345"},"values":[123],"timestamps":[1564617600000]}'
