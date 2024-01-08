package stats

import (
    "fmt"
    "io/ioutil"
    "encoding/json"

    "go.opentelemetry.io/otel/exporters/prometheus"
    api "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/sdk/metric"

    "github.com/cloudputation/service-factory/packages/config"
    "github.com/cloudputation/service-factory/packages/consul"
)


type FactoryInfo struct {
  FactoryState  string   `json:"factory-state"`
  Services     []string `json:"services"`
}

const meterName = "SF.Metrics"
var (
    ErrorCounter                 api.Int64Counter
    HealthEndpointCounter        api.Int64Counter
    ApplyEndpointCounter         api.Int64Counter
    DestroyEndpointCounter       api.Int64Counter
    RepoStatusEndpointCounter    api.Int64Counter
    ServiceStatusEndpointCounter api.Int64Counter
    SystemStatusEndpointCounter  api.Int64Counter
    SystemMetricsEndpointCounter api.Int64Counter
)


func InitMetrics() error {
    exporter, err := prometheus.New()
    if err != nil {
        return fmt.Errorf("Failed to initialize prometheus client: %v", err)
    }
    provider := metric.NewMeterProvider(metric.WithReader(exporter))
    meter := provider.Meter(meterName)

    ErrorCounter, err = meter.Int64Counter(
        "agent_errors",
        api.WithDescription("Counts the number of errors during agent runtime"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize the error counter: %v", err)
    }

    HealthEndpointCounter, err = meter.Int64Counter(
        "health_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /health endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /health endpoint counter: %v", err)
    }

    ApplyEndpointCounter, err = meter.Int64Counter(
        "apply_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /apply endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /apply endpoint counter: %v", err)
    }

    DestroyEndpointCounter, err = meter.Int64Counter(
        "destroy_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /destroy endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /destroy endpoint counter: %v", err)
    }

    RepoStatusEndpointCounter, err = meter.Int64Counter(
        "repo_status_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /repo/status endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /repo/status endpoint counter: %v", err)
    }

    ServiceStatusEndpointCounter, err = meter.Int64Counter(
        "service_status_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /service/status endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /service/status endpoint counter: %v", err)
    }

    SystemStatusEndpointCounter, err = meter.Int64Counter(
        "system_status_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /system/status endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /system/status endpoint counter: %v", err)
    }

    SystemMetricsEndpointCounter, err = meter.Int64Counter(
        "system_metrics_endpoint_hits",
        api.WithDescription("Counts the number of hits to the /system/metrics endpoint"),
    )
    if err != nil {
        return fmt.Errorf("Failed to initialize /system/metrics endpoint counter: %v", err)
    }


    return nil
}

func GenerateState() error {
  dirName := config.RootDir + "/" + config.AppConfig.DataDir + "/services"
  jsonData, err := GetFactoryInfo(dirName)
	if err != nil {
      return fmt.Errorf("Failed to get factory info: %v", err)
	}

	keyPath := config.ConsulServiceSummaryDataDir
	err = consul.ConsulStorePut(jsonData, keyPath)
	if err != nil {
      return fmt.Errorf("Failed to put data in Consul KV: %v", err)
	}


  return nil
}

func GetFactoryInfo(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
      return "", fmt.Errorf("Failed to read factory files in %s: %v", dirPath, err)
	}

	var dirNames []string
	for _, file := range files {
      if file.IsDir() {
          dirNames = append(dirNames, file.Name())
      }
	}

	factoryInfo := FactoryInfo{
		FactoryState: "running",
		Services:     dirNames,
	}

	jsonData, err := json.MarshalIndent(factoryInfo, "", "  ")
	if err != nil {
      return "", fmt.Errorf("Failed to process factory info: %v", err)
	}


	return string(jsonData), nil
}
