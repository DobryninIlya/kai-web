package influxdb

import (
	"context"
	"fmt"
	"time"

	//influx "github.com/influxdata/influxdb/client/v2"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	context.Context
	client   *influxdb3.Client
	log      *logrus.Logger
	database string
}

func NewMetrics(ctx context.Context, influxToken, influxUrl, database string, log *logrus.Logger) *Metrics {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     influxUrl,
		Token:    influxToken,
		Database: database,
	})
	if err != nil {
		logrus.Fatal("Ошибка инициализации клиента influxdb3", err)
		return nil
	}
	metrics := &Metrics{
		client:   client,
		log:      log,
		database: database,
		Context:  ctx,
	}
	return metrics
}

func (r *Metrics) IncrementAPI(path string) {
	line := []byte(fmt.Sprintf("api_requests,method=%v value=10 %v", path, time.Now().UnixNano()))
	if err := r.client.Write(r.Context, line); err != nil {
		r.log.Error("Ошибка записи в influxdb: ", err)
	}
}

//func initDB(metrics *Metrics) {
//	// Create a new database
//	q := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %s", metrics.database), "", "")
//	if response, err := metrics.client.Query(q); err == nil && response.Error() == nil {
//		fmt.Println("Database created")
//	}
//
//	// Create a new retention policy
//	//rp := influx.NewRetentionPolicy("myrp", metrics.database, "30d", 1, true) precision
//	q = influx.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY rp ON %s DURATION 30d REPLICATION 1d DEFAULT", metrics.database), metrics.database, "")
//	if response, err := metrics.client.Query(q); err == nil && response.Error() == nil {
//		fmt.Println("RetentionPolicy created")
//	}
//
//	q = influx.NewQuery(fmt.Sprintf("CREATE CONTINUOUS QUERY \"cq_requests_per_method\" ON \"%v\" BEGIN"+
//		"  SELECT count(*) INTO \"requests_per_method\""+
//		"  FROM \"http_requests\""+
//		"  GROUP BY time(1m), method"+
//		"END", metrics.database), metrics.database, "")
//	if response, err := metrics.client.Query(q); err == nil && response.Error() == nil {
//		fmt.Println("Table created")
//	}
//}
