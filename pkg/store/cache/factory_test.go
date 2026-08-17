// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package storecache

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewIndexCacheTracingEnabled(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		tracingConfig          string
		expectedTracingEnabled bool
	}{
		{
			name:                   "disabled by default",
			expectedTracingEnabled: false,
		},
		{
			name:                   "explicitly enabled",
			tracingConfig:          "tracing_enabled: true\n",
			expectedTracingEnabled: true,
		},
		{
			name:                   "explicitly disabled",
			tracingConfig:          "tracing_enabled: false\n",
			expectedTracingEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conf := []byte("type: IN-MEMORY\nconfig:\n  max_size: 1GB\n" + tc.tracingConfig)
			cache, err := NewIndexCache(log.NewNopLogger(), conf, prometheus.NewRegistry())
			testutil.Ok(t, err)

			_, tracingEnabled := cache.(*TracingIndexCache)
			testutil.Equals(t, tc.expectedTracingEnabled, tracingEnabled)
		})
	}
}

func TestIndexCacheMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	commonMetrics := NewCommonMetrics(reg)

	memcached := newMockedMemcachedClient(nil)
	_, err := NewRemoteIndexCache(log.NewNopLogger(), memcached, commonMetrics, reg, memcachedDefaultTTL)
	testutil.Ok(t, err)
	conf := []byte(`
max_size: 10MB
max_item_size: 1MB
`)
	// Make sure that the in memory cache does not register the same metrics of the remote index cache.
	// If so, we should move those metrics to the `commonMetrics`
	_, err = NewInMemoryIndexCache(log.NewNopLogger(), commonMetrics, reg, conf)
	testutil.Ok(t, err)
}
