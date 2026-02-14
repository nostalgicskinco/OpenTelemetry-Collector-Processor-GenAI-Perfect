package genaisafeprocessor

import (
  	"crypto/sha256"
  	"encoding/hex"
  	"strings"

  	"go.opentelemetry.io/collector/pdata/pcommon"
  	"go.opentelemetry.io/collector/pdata/ptrace"
  )

// redactSpan applies redaction policies to a single span.
func (p *genAIProc) redactSpan(s ptrace.Span) {
  	attrs := s.Attributes()

  	// 1) Redact configured keys (prompt/completion text).
  	for _, key := range p.cfg.Redact.Keys {
      		if v, ok := attrs.Get(key); ok && v.Type() == pcommon.ValueTypeStr {
            			orig := v.Str()
            			applyRedaction(attrs, key, orig, p.cfg.Redact.Mode, p.cfg.Redact.PreviewChars, p.cfg.Redact.Salt)
            		}
      	}

  	// 2) Generic denylist scrub across all string attributes.
  	attrs.Range(func(k string, v pcommon.Value) bool {
      		if v.Type() != pcommon.ValueTypeStr {
            			return true
            		}
      		val := v.Str()
      		for _, re := range p.denylist {
            			if re.MatchString(k) || re.MatchString(val) {
                    				attrs.PutStr(k, "[REDACTED]")
                    				break
                    			}
            		}
      		return true
      	})
  }

// applyRedaction replaces or hashes an attribute value based on the mode.
func applyRedaction(attrs pcommon.Map, key, orig, mode string, previewChars int, salt string) {
  	trim := strings.TrimSpace(orig)
  	if trim == "" {
      		return
      	}

  	hash := sha256.Sum256([]byte(salt + "::" + trim))
  	h := hex.EncodeToString(hash[:])

  	switch mode {
      	case "drop":
      		attrs.Remove(key)
      		attrs.PutStr(key+".hash", h)
      	case "hash":
      		attrs.PutStr(key, "[HASHED]")
      		attrs.PutStr(key+".hash", h)
      	case "truncate":
      		attrs.PutStr(key, truncate(trim, previewChars))
      		attrs.PutStr(key+".hash", h)
      	default: // hash_and_preview
      		attrs.PutStr(key, truncate(trim, previewChars))
      		attrs.PutStr(key+".hash", h)
      	}
  }

// truncate returns the first n characters of s, appending an ellipsis if truncated.
func truncate(s string, n int) string {
  	if n <= 0 || len(s) <= n {
      		return s
      	}
  	return s[:n] + "..."
  }
