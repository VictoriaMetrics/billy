VERSION=v1.35.0
ARCHIVE=victoria-metrics-$VERSION.tar.gz

test -f $ARCHIVE || curl -L https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/$VERSION/$ARCHIVE > $ARCHIVE
tar xzf $ARCHIVE
ulimit -n 100000
./victoria-metrics-prod -retentionPeriod=200 -storageDataPath=/mnt/disks/billy/victoria-metrics-data -search.maxQueryDuration=24h -search.maxUniqueTimeseries=2000000
