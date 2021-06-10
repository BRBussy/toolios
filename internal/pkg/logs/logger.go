package logs

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Level              zerolog.Level
	LevelString        string
	LogDir             string
	Logfile            string
	DisableFileLogs    bool
	DisableConsoleLogs bool
	EnableJSONLogs     bool
	EnableColorLogs    bool
}

var levels = map[string]zerolog.Level{
	"Debug":    zerolog.DebugLevel,
	"Info":     zerolog.InfoLevel,
	"Warn":     zerolog.WarnLevel,
	"Error":    zerolog.ErrorLevel,
	"Fatal":    zerolog.FatalLevel,
	"Panic":    zerolog.PanicLevel,
	"Disabled": zerolog.Disabled,
}

func Setup(config *Config) {
	if config == nil {
		config = &Config{}
	}

	if config.LevelString != "" {
		zerolog.SetGlobalLevel(levels[config.LevelString])
	} else {
		zerolog.SetGlobalLevel(config.Level)
	}

	if config.DisableFileLogs && config.DisableConsoleLogs {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return
	}

	var jsonFileWriter io.Writer
	var normalFileWriter io.Writer
	var colorFileWriter io.Writer
	var consoleWriter io.Writer

	logDir := "logs"
	if config.LogDir != "" {
		logDir = config.LogDir
	}

	if config.Logfile == "" {
		config.Logfile = "application"
	}

	if !config.DisableConsoleLogs {
		consoleWriter = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		}
	}

	if !config.DisableFileLogs {
		normalLogFileName := filepath.Join(".", logDir, config.Logfile+".log")
		_ = os.MkdirAll(filepath.Dir(normalLogFileName), os.ModePerm)
		fW := &lumberjack.Logger{
			Filename:   normalLogFileName,
			MaxSize:    20, // megabytes
			MaxAge:     5,
			MaxBackups: 5,
			LocalTime:  true,
			Compress:   true,
		}

		normalFileWriter = zerolog.ConsoleWriter{
			Out:        fW,
			NoColor:    true,
			TimeFormat: "2006-01-02 15:04:05",
		}

		if config.EnableColorLogs {
			colorLogFileName := filepath.Join(".", logDir, config.Logfile+"_color.log")
			_ = os.MkdirAll(filepath.Dir(colorLogFileName), os.ModePerm)
			cW := &lumberjack.Logger{
				Filename:   colorLogFileName,
				MaxSize:    5, // megabytes
				MaxAge:     1,
				MaxBackups: 2,
				LocalTime:  true,
				Compress:   true,
			}
			colorFileWriter = zerolog.ConsoleWriter{
				Out:        cW,
				TimeFormat: time.RFC3339,
			}
		}

		if config.EnableJSONLogs {
			jsonLogFileName := filepath.Join(".", logDir, config.Logfile+"_json.log")
			_ = os.MkdirAll(filepath.Dir(jsonLogFileName), os.ModePerm)

			jsonFileWriter = &lumberjack.Logger{
				Filename:  jsonLogFileName,
				MaxSize:   50, // megabytes
				LocalTime: true,
				Compress:  true,
			}
		}
	}

	var writers []io.Writer

	if consoleWriter != nil {
		writers = append(writers, consoleWriter)
	}

	if normalFileWriter != nil {
		writers = append(writers, normalFileWriter)
	}

	if colorFileWriter != nil {
		writers = append(writers, colorFileWriter)
	}

	if jsonFileWriter != nil {
		writers = append(writers, jsonFileWriter)
	}

	log.Logger = log.Output(io.MultiWriter(writers...)).
		With().
		Caller().
		Logger()

}
