package compositelogger

// LogConfig is a struct containing the configuration for an individual log.
type LogConfig struct {
	Location string `json:"location" bson:"location" yaml:"location" hcl:"location"` // location of the log
	Level    string `json:"level" bson:"level" yaml:"level" hcl:"level"` // level of the log: INFO, WARN, DEBUG, etc.
	Format   string `json:"format" bson:"format" yaml:"format" hcl:"format"` // text or json
}
