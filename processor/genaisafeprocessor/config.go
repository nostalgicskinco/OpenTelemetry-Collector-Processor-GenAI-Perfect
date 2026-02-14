package genaisafeprocessor

import "time"

// Config holds the full processor configuration.
type Config struct {
  	Redact        RedactConfig        `mapstructure:"redact"`
  	Metrics       MetricsConfig       `mapstructure:"metrics"`
  	LoopDetection LoopDetectionConfig `mapstructure:"loop_detection"`
  }

// RedactConfig controls how span attributes are redacted.
type RedactConfig struct {
  	Mode         string   `mapstructure:"mode"` // drop|hash|hash_and_preview|truncate
  	PreviewChars int      `mapstructure:"preview_chars"`
  	Salt         string   `mapstructure:"salt"`
  	Keys         []string `mapstructure:"keys"`
  	DenylistRe   []string `mapstructure:"denylist_regex"`
  }

// MetricsConfig controls token/cost metric extraction.
type MetricsConfig struct {
  	Enable              bool          `mapstructure:"enable"`
  	EmitInterval        time.Duration `mapstructure:"emit_interval"`
  	TokenAttrCandidates []string      `mapstructure:"token_attr_candidates"`
  	CostAttrCandidates  []string      `mapstructure:"cost_attr_candidates"`
  }

// LoopDetectionConfig controls simple runaway/loop detection.
type LoopDetectionConfig struct {
  	Enable          bool   `mapstructure:"enable"`
  	ToolSpanNameRe  string `mapstructure:"tool_span_name_regex"`
  	RepeatThreshold int    `mapstructure:"repeat_threshold"`
  }
