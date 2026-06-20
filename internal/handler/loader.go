package handler

import (
	"fmt"

	"github.com/jwart212/protoc-gen-go-mapper/internal/config"
)

// LoadHandlers converts field handler configurations into handler instances.
func LoadHandlers(cfgs []config.FieldHandlerConfig) (*HandlerRegistry, error) {
	registry := NewHandlerRegistry()

	for _, cfg := range cfgs {
		var h FieldHandler

		switch cfg.Type {
		case "skip":
			h = NewSkipHandler(cfg.MatchField)
		case "type_assertion":
			h = NewTypeAssertionHandler(cfg.MatchField, cfg.MatchDBTypes, cfg.AssertType)
		case "default_value":
			h = NewDefaultValueHandler(cfg.MatchField, cfg.MatchDBTypes, cfg.DefaultValue)
		case "field_mapping":
			h = NewFieldMappingHandler(cfg.MatchField, cfg.MatchMessage, cfg.MatchDBTypes, cfg.ToProto, cfg.ToDB)
		default:
			return nil, fmt.Errorf("unknown handler type: %s", cfg.Type)
		}

		registry.Register(h)
	}

	return registry, nil
}

// LoadTypeConversions converts type conversion configurations into handler instances.
func LoadTypeConversions(cfgs []config.TypeConversionConfig, aliases map[string]config.TypeAliasConfig) (*HandlerRegistry, error) {
	registry := NewHandlerRegistry()

	for _, cfg := range cfgs {
		// If alias is specified, resolve it
		protoType := cfg.ProtoType
		dbType := cfg.DBType
		isOptional := cfg.IsOptional
		toProtoExpr := cfg.ToProtoExpr
		toDBExpr := cfg.ToDBExpr

		if cfg.Alias != "" {
			alias, ok := aliases[cfg.Alias]
			if !ok {
				return nil, fmt.Errorf("type alias %q not found", cfg.Alias)
			}
			// Override with alias values if not set in config
			if protoType == "" {
				protoType = alias.ProtoType
			}
			if dbType == "" {
				dbType = alias.DBType
			}
			if !cfg.IsOptional {
				isOptional = alias.IsOptional
			}
			if toProtoExpr == "" {
				toProtoExpr = alias.ToProtoExpr
			}
			if toDBExpr == "" {
				toDBExpr = alias.ToDBExpr
			}
		}

		// Ensure we have the required fields
		if toProtoExpr == "" {
			fmt.Printf("Warning: type conversion missing toProtoExpr, skipping\n")
			continue
		}

		// Create handler with patterns if configured
		if cfg.MatchFieldPattern != "" || cfg.MatchMessagePattern != "" || cfg.PointerStrategy != "" {
			h, err := NewTypeConversionHandlerWithPatterns(
				protoType, dbType, isOptional, toProtoExpr, toDBExpr,
				cfg.Priority, cfg.MatchFieldPattern, cfg.MatchMessagePattern, cfg.PointerStrategy,
			)
			if err != nil {
				return nil, fmt.Errorf("creating type conversion handler with patterns: %w", err)
			}
			registry.Register(h)
		} else {
			h := NewTypeConversionHandler(protoType, dbType, isOptional, toProtoExpr, toDBExpr, cfg.Priority)
			registry.Register(h)
		}
	}

	return registry, nil
}
