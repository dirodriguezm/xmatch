package app

import (
	"reflect"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
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
	require.Equal(t, 2, reflect.ValueOf(service).Elem().FieldByName("externalClients").Len())
	require.Equal(t, 1, reflect.ValueOf(service).Elem().FieldByName("lightcurveFilters").Len())
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
	require.Equal(t, 2, reflect.ValueOf(service).Elem().FieldByName("externalClients").Len())
	require.Equal(t, 1, reflect.ValueOf(service).Elem().FieldByName("lightcurveFilters").Len())
}
