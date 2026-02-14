package genaisafeprocessor

import (
  	"regexp"

  	"go.opentelemetry.io/collector/pdata/ptrace"
  )

// detectSimpleLoop scans a batch of spans for repeated tool/function calls
// that exceed the configured threshold, signalling a possible runaway agent loop.
func detectSimpleLoop(spans ptrace.SpanSlice, toolRe *regexp.Regexp, threshold int) bool {
  	if threshold <= 0 {
      		threshold = 6
      	}

  	counts := map[string]int{}
  	for i := 0; i < spans.Len(); i++ {
      		name := spans.At(i).Name()
      		if toolRe != nil && toolRe.MatchString(name) {
            			counts[name]++
            			if counts[name] >= threshold {
                    				return true
                    			}
            		}
      	}
  	return false
  }
