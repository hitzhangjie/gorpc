package config

type Config interface {
	LoadConfig(cfg string) error
	Read(section, property string, dftValue string) string
	ReadInt(section, property string, dftValue int) int
	ReadBool(section, property string, dftValue bool) bool
}
