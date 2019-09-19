package main

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/support"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/shamil/s3-bucket-permissions-exporter/aws_support"
)

const namespace = "s3_bucket_permissions"

var (
	desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "status"),
		"S3 Bucket permissions status, per region, bucket and permission (1 for true, 0 for false)",
		[]string{"region", "bucket", "permission"},
		nil,
	)

	scrapeError = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "last_scrape_error",
		Help:      "Whether the last scrape of metrics from TrustedAdvisor resulted in an error (1 for error, 0 for success).",
	})

	permissionStatus = map[string]float64{
		"No bucket policy": 0,
		"No":               0,
		"Yes":              1,
	}
)

// Collector collects data from TrustedAdvisor check
type Collector struct {
	awsSupport *aws_support.AwsSupport

	nextRefreshTime time.Time
	mutex           sync.Mutex
}

// NewCollector creates a new instance of Collector with some defaults
func NewCollector() *Collector {
	return &Collector{
		awsSupport:      aws_support.New(),
		nextRefreshTime: time.Now(),
	}
}

func (c *Collector) processCheckDetail(detail *support.TrustedAdvisorResourceDetail) (metrics []prometheus.Metric) {
	region := *detail.Region
	bucket := *detail.Metadata[2]
	aclAllowsList := *detail.Metadata[3]
	aclAllowsUploadDelete := *detail.Metadata[4]
	policyAllowsAccess := *detail.Metadata[6]

	commonLabels := []string{region, bucket}

	metric, err := prometheus.NewConstMetric(
		desc,
		prometheus.GaugeValue,
		permissionStatus[aclAllowsList],
		append(commonLabels, "acl-allows-list")...,
	)

	if err != nil {
		log.With("err", err).Errorf("failed to create metric")
	} else {
		metrics = append(metrics, metric)
	}

	metric, err = prometheus.NewConstMetric(
		desc,
		prometheus.GaugeValue,
		permissionStatus[aclAllowsUploadDelete],
		append(commonLabels, "acl-allows-upload-delete")...,
	)

	if err != nil {
		log.With("err", err).Errorf("failed to create metric")
	} else {
		metrics = append(metrics, metric)
	}

	metric, err = prometheus.NewConstMetric(
		desc,
		prometheus.GaugeValue,
		permissionStatus[policyAllowsAccess],
		append(commonLabels, "policy-allows-access")...,
	)

	if err != nil {
		log.With("err", err).Errorf("failed to create metric")
	} else {
		metrics = append(metrics, metric)
	}

	return metrics
}

func (c *Collector) scrape(ch chan<- prometheus.Metric) {
	scrapeError.Set(0)

	result, err := c.awsSupport.DescribeS3BucketPermissionsCheck()
	if err != nil {
		log.With("err", err).Errorf("failed describing TrustedAdbisor check")
		scrapeError.Set(1)
	}

	for _, d := range result.FlaggedResources {
		for _, m := range c.processCheckDetail(d) {
			ch <- m
		}
	}

	if time.Now().After(c.nextRefreshTime) {
		log.Info("refreshing TrustedAdvisor check...")
		result, err := c.awsSupport.RefreshtS3BucketPermissionsCheck()

		if err != nil {
			log.With("err", err).Warn("failed refreshing TrustedAdvisor check")
			return
		}

		c.nextRefreshTime = time.Now().Add(time.Duration(*result.MillisUntilNextRefreshable) * time.Millisecond)
	}
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	c.mutex.Lock()

	defer func() {
		c.mutex.Unlock()
		log.Infof("finished scraping TrustedAdvisor, took %v.", time.Now().Sub(start))
	}()

	log.Info("scraping TrustedAdvisor...")
	c.scrape(ch)
	ch <- scrapeError
}
