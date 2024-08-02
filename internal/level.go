package internal

const (
	LevelErr  = "ERR"
	LevelWarn = "WARN"
	LevelInfo = "INFO"
)

func GetSupportedLevels() []string {
	return []string{LevelErr, LevelWarn, LevelInfo}
}

func GetBreakingLevels() []string {
	return []string{LevelErr, LevelWarn}
}

func GetSupportedLevelsLower() []string {
	return []string{"error", "warn", "info"}
}
