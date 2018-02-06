package compositelogger

type LogConfig struct {
  Location string `json:"location" bson:"location" yaml:"location" hcl:"location"`
  Level string `json:"level" bson:"level" yaml:"level" hcl:"level"`
  Format string `json:"format" bson:"format" yaml:"format" hcl:"format"`
}
