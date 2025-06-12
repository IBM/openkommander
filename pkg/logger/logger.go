package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

var (
	Logger *slog.Logger
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type Config struct {
	Level     LogLevel
	Format    string
	AddColors bool
}

func DefaultConfig() *Config {
	return &Config{
		Level:     LevelInfo,
		Format:    "pretty",
		AddColors: true,
	}
}

// Custom handler for prettier logs
type PrettyHandler struct {
	opts      *slog.HandlerOptions
	writer    io.Writer
	addColors bool
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions, addColors bool) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts:      opts,
		writer:    w,
		addColors: addColors,
	}
}

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

func (h *PrettyHandler) colorize(color, text string) string {
	if !h.addColors {
		return text
	}
	return color + text + ColorReset
}

func (h *PrettyHandler) levelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return ColorGray
	case slog.LevelInfo:
		return ColorBlue
	case slog.LevelWarn:
		return ColorYellow
	case slog.LevelError:
		return ColorRed
	default:
		return ColorReset
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	// Format timestamp
	timestamp := r.Time.Format("15:04:05.000")

	// Format level with color and padding
	levelStr := r.Level.String()
	levelStr = strings.ToUpper(levelStr)
	switch len(levelStr) {
	case 4: // INFO, WARN
		levelStr = levelStr + " "
	case 5: // DEBUG, ERROR
		// no padding needed
	default:
		levelStr = fmt.Sprintf("%-5s", levelStr)
	}

	coloredLevel := h.colorize(h.levelColor(r.Level), levelStr)

	// Format message
	message := r.Message

	// Build the log line
	var logLine strings.Builder

	// Add timestamp
	logLine.WriteString(h.colorize(ColorGray, timestamp))
	logLine.WriteString(" ")

	// Add level
	logLine.WriteString(coloredLevel)
	logLine.WriteString(" ")

	// Add message
	if h.addColors && r.Level >= slog.LevelWarn {
		logLine.WriteString(h.colorize(ColorBold, message))
	} else {
		logLine.WriteString(message)
	}

	// Add attributes in a readable format
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	if len(attrs) > 0 {
		logLine.WriteString(" ")
		logLine.WriteString(h.colorize(ColorCyan, "•"))
		logLine.WriteString(" ")

		var parts []string

		// Handle special attributes first
		if broker, ok := attrs["broker"]; ok {
			parts = append(parts, h.colorize(ColorGreen, fmt.Sprintf("broker=%v", broker)))
			delete(attrs, "broker")
		}

		if method, ok := attrs["method"]; ok {
			parts = append(parts, h.colorize(ColorBlue, fmt.Sprintf("method=%v", method)))
			delete(attrs, "method")
		}

		if statusCode, ok := attrs["status_code"]; ok {
			color := ColorGreen
			if code, ok := statusCode.(int); ok && code >= 400 {
				color = ColorRed
			}
			parts = append(parts, h.colorize(color, fmt.Sprintf("status=%v", statusCode)))
			delete(attrs, "status_code")
		}

		if duration, ok := attrs["duration_ms"]; ok {
			parts = append(parts, h.colorize(ColorYellow, fmt.Sprintf("duration=%vms", duration)))
			delete(attrs, "duration_ms")
		}

		// Add remaining attributes
		for key, value := range attrs {
			if key == "error" && value != nil {
				parts = append(parts, h.colorize(ColorRed, fmt.Sprintf("%s=%v", key, value)))
			} else {
				parts = append(parts, fmt.Sprintf("%s=%v", key, value))
			}
		}

		logLine.WriteString(strings.Join(parts, " "))
	}

	logLine.WriteString("\n")

	_, err := h.writer.Write([]byte(logLine.String()))
	return err
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll return the same handler
	// In a more complex implementation, you'd create a new handler with the attrs
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll return the same handler
	// In a more complex implementation, you'd handle grouping
	return h
}

func Init(config *Config) {
	if config == nil {
		config = DefaultConfig()
	}

	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch config.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "pretty":
		handler = NewPrettyHandler(os.Stdout, opts, config.AddColors)
	default:
		handler = NewPrettyHandler(os.Stdout, opts, config.AddColors)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

func GetLogger() *slog.Logger {
	if Logger == nil {
		Init(DefaultConfig())
	}
	return Logger
}

func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}

func WithGroup(name string) *slog.Logger {
	return GetLogger().WithGroup(name)
}

// Convenience functions for common logging patterns
func HTTP(msg string, method, path string, statusCode int, duration int64, args ...any) {
	allArgs := append([]any{"method", method, "path", path, "status_code", statusCode, "duration_ms", duration}, args...)
	Info(msg, allArgs...)
}

func Kafka(msg string, broker, operation string, args ...any) {
	allArgs := append([]any{"broker", broker, "operation", operation}, args...)
	Info(msg, allArgs...)
}

func KafkaError(msg string, broker, operation string, err error, args ...any) {
	allArgs := append([]any{"broker", broker, "operation", operation, "error", err}, args...)
	Error(msg, allArgs...)
}
