package model

type ConfigEntity struct {
	MaxCache       int   `json:"max_cache" yaml:"max_cache" bson:"max_cache" mapstructure:"max_cache"`
	MaxConcurrency int32 `json:"max_concurrency" yaml:"max_concurrency" bson:"max_concurrency" mapstructure:"max_concurrency"`
	MinConcurrency int32 `json:"min_concurrency" yaml:"min_concurrency" bson:"min_concurrency" mapstructure:"min_concurrency"`
}
