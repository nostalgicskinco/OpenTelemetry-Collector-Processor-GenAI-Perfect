package main

import (
  	"log"

  	"go.opentelemetry.io/collector/component"
  	"go.opentelemetry.io/collector/exporter/otlpexporter"
  	"go.opentelemetry.io/collector/otelcol"
  	"go.opentelemetry.io/collector/processor"
  	"go.opentelemetry.io/collector/receiver/otlpreceiver"

  	"github.com/nostalgicskinco/OpenTelemetry-Collector-Processor-GenAI-Perfect/processor/genaisafeprocessor"
  )

func main() {
  	info := component.BuildInfo{
      		Command:     "otelcol-genai-safe",
      		Description: "OpenTelemetry Collector with GenAI Safe Processor",
      		Version:     "0.1.0",
      	}

  	factories, err := otelcol.Factories(
      		otelcol.WithReceivers(otlpreceiver.NewFactory()),
      		otelcol.WithProcessors(genaisafeprocessor.NewFactory()),
      		otelcol.WithExporters(otlpexporter.NewFactory()),
      	)
  	if err != nil {
      		log.Fatalf("failed to build factories: %v", err)
      	}

  	settings := otelcol.CollectorSettings{
      		Factories: func() (otelcol.Factories, error) { return factories, nil },
      		BuildInfo: info,
      	}

  	cmd := otelcol.NewCommand(settings)
  	if err := cmd.Execute(); err != nil {
      		log.Fatalf("collector exited with error: %v", err)
      	}
  }
