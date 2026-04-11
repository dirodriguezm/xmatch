package app

import (
	"reflect"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
	"github.com/stretchr/testify/require"
)

func TestLightcurveService_AlwaysAddsZtfDrClient(t *testing.T) {
	cfg := config.Config{
		Service: config.ServiceConfig{
			LightcurveServiceConfig: config.LightcurveServiceConfig{
				NeowiseConfig: config.NeowiseConfig{},
				ZtfDrConfig:   config.ZtfDrConfig{UseIdFilter: false},
			},
		},
	}

	service, err := LightcurveService(cfg, &conesearch.ConesearchService{})
	require.NoError(t, err)

	sources := reflect.ValueOf(service).Elem().FieldByName("sources")
	require.Equal(t, 2, sources.Len())
	require.Equal(t, "ztf", sources.Index(1).FieldByName("Catalog").String())
	require.Equal(t, reflect.ValueOf(lightcurve.DummyLightcurveFilter).Pointer(), sources.Index(1).FieldByName("Filter").Pointer())
}

func TestLightcurveService_AddsZtfDrFilterWhenEnabled(t *testing.T) {
	cfg := config.Config{
		Service: config.ServiceConfig{
			LightcurveServiceConfig: config.LightcurveServiceConfig{
				NeowiseConfig: config.NeowiseConfig{},
				ZtfDrConfig:   config.ZtfDrConfig{UseIdFilter: true},
			},
		},
	}

	service, err := LightcurveService(cfg, &conesearch.ConesearchService{})
	require.NoError(t, err)

	sources := reflect.ValueOf(service).Elem().FieldByName("sources")
	require.Equal(t, 2, sources.Len())
	require.Equal(t, "ztf", sources.Index(1).FieldByName("Catalog").String())
	require.Equal(t, reflect.ValueOf(ztfdr.Filter).Pointer(), sources.Index(1).FieldByName("Filter").Pointer())
}
