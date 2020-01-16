VERSION=v1.32.3
ARCHIVE=victoria-metrics-$VERSION.tar.gz

test -f $ARCHIVE || curl -L https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/$VERSION/$ARCHIVE > $ARCHIVE
tar xzf $ARCHIVE
./victoria-metrics-prod -retentionPeriod=200 -storageDataPath=/mnt/disks/billy/victoria-metrics-data
