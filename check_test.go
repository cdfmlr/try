package try

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
)

type mockLoggerWriter struct {
	lines []string
}

func (w *mockLoggerWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("🎭 mockLoggerWriter.Write: %s\n", string(p))
	if w.lines == nil {
		w.lines = make([]string, 0)
	}

	w.lines = append(w.lines, string(p))
	return len(p), nil
}

func (w *mockLoggerWriter) reset() {
	w.lines = make([]string, 0)
}

func TestCheckHandle(t *testing.T) {
	logWriter := new(mockLoggerWriter)

	slog.SetDefault(slog.New(
		slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	))

	customLogWriter := new(mockLoggerWriter)
	customLogger := slog.New(
		slog.NewJSONHandler(customLogWriter, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelWarn,
		}),
	)

	tests := []struct {
		name         string
		f            func() error
		expectedErr  error
		expectPanics bool
	}{
		{
			name: "basicCheckErrHandle",
			f: func() (err error) {
				defer Handle(&err)

				err = fmt.Errorf("there is a error")
				Check(err)

				return nil
			},
			expectedErr:  fmt.Errorf("there is a error"),
			expectPanics: false,
		},
		{
			name: "badUsage_CheckWithoutHandle",
			f: func() (err error) {
				err = fmt.Errorf("there is a error")
				Check(err)

				return nil
			},
			expectPanics: true,
		},
		{
			name: "HandleWithCustomPanic",
			f: func() (err error) {
				defer func() {
					r := recover()
					if r != "custom panic" {
						t.Errorf("expect recover() = %#v, got %#v", "custom panic", r)
					}
					t.Logf("✅ recover() = %#v\n", r)
					panic(r)
				}()
				defer Handle(&err)

				panic("custom panic")

				return nil
			},
			expectPanics: true,
		},
		{
			name: "CheckWithCustomHandler",
			f: func() (err error) {
				logWriter.reset()

				defer func() {
					if len(logWriter.lines) != 1 {
						t.Errorf("expect logWriter.lines = %#v, got %#v", 1, len(logWriter.lines))
					}
					record := map[string]any{}
					err := json.Unmarshal([]byte(logWriter.lines[0]), &record)
					if err != nil {
						t.Errorf("expect json.Unmarshal to pass, got %#v", err)
					}
					if record["level"] != "ERROR" {
						t.Errorf("expect record[\"level\"] = %#v, got %#v", "ERROR", record["level"])
					}
					if record["err"] != "custom handler: there is a error" {
						t.Errorf("expect record[\"err\"] = %#v, got %#v", "custom handler: there is a error", record["err"])
					}
				}()
				defer Handle(&err)

				err = fmt.Errorf("there is a error")
				Check(err,
					func(err error) error {
						return fmt.Errorf("custom handler: %w", err)
					},
					Log(),
					func(err error) error {
						return nil // handled
					},
					func(err error) error { // this is not expected to be called
						return fmt.Errorf("custom handler failed (should early stoped): %w", err)
					},
				)

				return nil
			},
			expectedErr:  nil,
			expectPanics: false,
		},
		{
			name: "logWithMsgWithLevel",
			f: func() (err error) {
				logWriter.reset()
				defer func() { // check log
					if len(logWriter.lines) != 1 {
						t.Errorf("expect logWriter.lines = %#v, got %#v", 1, len(logWriter.lines))
					}
					record := map[string]any{}
					err := json.Unmarshal([]byte(logWriter.lines[0]), &record)
					if err != nil {
						t.Errorf("expect json.Unmarshal to pass, got %#v", err)
					}
					if record["level"] != "WARN" {
						t.Errorf("expect record[\"level\"] = %#v, got %#v", "WARN", record["level"])
					}
					if record["msg"] != "check error failed" {
						t.Errorf("expect record[\"msg\"] = %#v, got %#v", "check error failed", record["msg"])
					}
					if record["err"] != "there is a error" {
						t.Errorf("expect record[\"err\"] = %#v, got %#v", "there is a error", record["err"])
					}
					if !strings.Contains(fmt.Sprint(record["source"]), "try.TestCheckHandle") {
						t.Errorf("expect record[\"source\"] contains %#v, got %#v", "try.TestCheckHandle", record["source"])
					}
					t.Logf("✅ logWriter.lines = %#v\n", logWriter.lines)
				}()
				defer Handle(&err)

				err = fmt.Errorf("there is a error")
				fmt.Printf("err generate: %v\n", err)

				Check(err, Log(WithMsg("check error failed"), WithLevel(slog.LevelWarn)))

				return nil
			},
			expectedErr:  fmt.Errorf("there is a error"),
			expectPanics: false,
		},
		{
			name: "logWithCustomLogger",
			f: func() (err error) {
				customLogWriter.reset()
				logWriter.reset()
				defer func() { // check log
					if len(logWriter.lines) != 0 {
						t.Errorf("expect logWriter.lines = %#v, got %#v", 0, len(logWriter.lines))
					}
					if len(customLogWriter.lines) != 1 {
						t.Errorf("expect logWriter.lines = %#v, got %#v", 1, len(customLogWriter.lines))
					}
					record := map[string]any{}
					e := json.Unmarshal([]byte(customLogWriter.lines[0]), &record)
					if e != nil {
						t.Errorf("expect json.Unmarshal to pass, got %#v", e)
					}
					if record["level"] != "WARN" {
						t.Errorf("expect record[\"level\"] = %#v, got %#v", "WARN", record["level"])
					}
					if record["msg"] != "this record should be written to custom logger" {
						t.Errorf("expect record[\"msg\"] = %#v, got %#v", "check error failed", record["msg"])
					}
					if record["err"] != "there is a error" {
						t.Errorf("expect record[\"err\"] = %#v, got %#v", "there is a error", record["err"])
					}
					if !strings.Contains(fmt.Sprint(record["source"]), "try.TestCheckHandle") {
						t.Errorf("expect record[\"source\"] contains %#v, got %#v", "try.TestCheckHandle", record["source"])
					}
					t.Logf("✅ customLogWriter.lines = %#v\n", customLogWriter.lines)
				}()
				defer Handle(&err)

				err = fmt.Errorf("there is a error")
				fmt.Printf("err generate: %v\n", err)

				Check(err,
					Log(
						WithLogger(customLogger),
						WithMsg("this record should not be enabled"),
						WithLevel(slog.LevelInfo)),
					func(err error) error {
						return nil // continue
					})

				Check(err,
					Log(
						WithLogger(customLogger),
						WithMsg("this record should be written to custom logger"),
						WithLevel(slog.LevelWarn)))

				return nil
			},
			expectedErr:  fmt.Errorf("there is a error"),
			expectPanics: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			defer func() {
				r := recover()
				if (r != nil) != tt.expectPanics {
					t.Errorf("expect panics = %#v, got %#v", tt.expectPanics, r)
				}
				t.Logf("✅ recover() = %#v: %v\n", r, r)
			}()

			err = tt.f()
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("expect err = %#v, got %#v", tt.expectedErr, err)
			}
			t.Logf("✅ err = %#v\n", err)
		})
	}
}

func benchmarkCheckHandleWithLog() (err error) {
	defer Handle(&err)

	err = fmt.Errorf("there is a error")
	Check(err, Log())

	return nil
}

func benchmarkCheckHandleWithoutLog() (err error) {
	defer Handle(&err)

	err = fmt.Errorf("there is a error")
	Check(err)

	return nil
}

func benchmarkIfErrNeNilWithLog() (err error) {
	defer Handle(&err)

	err = fmt.Errorf("there is a error")
	if err != nil {
		slog.Error("if got an error", "err", err)
		return err
	}

	return nil
}

func benchmarkIfErrNeNilWithoutLog() (err error) {
	defer Handle(&err)

	err = fmt.Errorf("there is a error")
	if err != nil {
		return err
	}

	return nil
}

type NullWriter struct{}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkCheckHandle(b *testing.B) {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(&NullWriter{}, &slog.HandlerOptions{}),
	))

	b.Run("CheckHandleWithLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = benchmarkCheckHandleWithLog()
		}
	})
	b.Run("IfErrNeNilWithLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = benchmarkIfErrNeNilWithLog()
		}
	})
	b.Run("CheckHandleWithoutLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = benchmarkCheckHandleWithoutLog()
		}
	})
	b.Run("IfErrNeNilWithoutLog", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = benchmarkIfErrNeNilWithoutLog()
		}
	})
}
