package logger_test

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ivch/dynasty/common/logger"
)

func TestLevel_String(t *testing.T) {
	type tcase struct {
		l    logger.Level
		want string
	}
	tests := map[string]tcase{
		"Disabled": {logger.DSB, "Disabled"},
		"Error":    {logger.ERR, "Error"},
		"Info":     {logger.INF, "Info"},
		"Warning":  {logger.WRN, "Warning"},
		"Debug":    {logger.DBG, "Debug"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.l.String(); got != tc.want {
				t.Errorf("String() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("NewStdLog()", func(t *testing.T) {
		got := logger.NewStdLog()
		if got.Lvl != logger.INF {
			t.Errorf("default log lvl should be INF but got %s", got.Lvl.String())
		}
		if got.Err.Writer() != os.Stderr || got.Inf.Writer() != os.Stderr ||
			got.Wrn.Writer() != os.Stderr || got.Dbg.Writer() != os.Stderr {
			t.Errorf("all default writers should be os.Stderr")
		}
		if !reflect.TypeOf(got).Implements(reflect.TypeOf((*logger.Logger)(nil)).Elem()) {
			t.Errorf("type does't implement logger.Logger interface")
		}
		if reflect.TypeOf(*got).Name() != "StdLog" {
			t.Errorf("type name should be StdLog but got: %s", reflect.TypeOf(*got).Name())
		}
		if reflect.TypeOf(got).String() != "*logger.StdLog" {
			t.Errorf("struct type returned by NewStdLog() should have name StdLog but got: %s", reflect.TypeOf(got).String())
		}
		if reflect.TypeOf(*got).Kind() != reflect.Struct {
			t.Errorf("type kind returned by NewStdLog() should be struct but got: %s", reflect.TypeOf(got).Kind())
		}
	})
	t.Run("NewStdLog(WithWriter)", func(t *testing.T) {
		tb := &testWriter{byf: make([]byte, 0, 20)}
		strToLog := "test string"
		got := logger.NewStdLog(logger.WithWriter(tb))
		got.Info(strToLog)

		if !strings.Contains(tb.String(), strToLog) {
			t.Errorf("expected %s but got %s", strToLog, tb.String())
		}

		if got.Lvl != logger.INF {
			t.Errorf("default log lvl should be INF but got %s", got.Lvl.String())
		}
		if got.Err.Writer() != tb || got.Inf.Writer() != tb ||
			got.Wrn.Writer() != tb || got.Dbg.Writer() != tb {
			t.Errorf("all default writers should be os.Stderr")
		}
	})
}

func TestStdLog_Print(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log := logger.NewStdLog(logger.WithWriter(buf), logger.WithLevel(logger.DBG))
	log.Error("hello, %s", "error")
	log.Warn("hello, %s", "warn")
	log.Info("hello, %s", "info")
	log.Debug("hello, %s", "debug")
	wVal := buf.String()
	if !strings.Contains(wVal, "hello, debug") {
		t.Error("invalid debug log output")
	}
	if !strings.Contains(wVal, "hello, warn") {
		t.Error("invalid warn log output")
	}
	if !strings.Contains(wVal, "hello, info") {
		t.Error("invalid info log output")
	}
	if !strings.Contains(wVal, "hello, error") {
		t.Error("invalid error log output")
	}
}

type testWriter struct {
	n   int
	byf []byte
}

func (w *testWriter) String() string {
	return string(w.byf)
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	for i, b := range p {
		w.byf = append(w.byf, b)
		w.n = i
	}
	return w.n, nil
}

func TestStdLog_Debug(t *testing.T) {
	type fields struct {
		lvl logger.Level
		w   *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		exec   func(logger logger.Logger)
		want   string
	}{
		{
			name: "no error",
			fields: fields{
				lvl: logger.DSB,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger logger.Logger) {
				logger.Error("errorLog")
			},
			want: "errorLog",
		},
		{
			name: "no warn",
			fields: fields{
				lvl: logger.ERR,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger logger.Logger) {
				logger.Warn("warnLog")
			},
			want: "warnLog",
		},
		{
			name: "no info",
			fields: fields{
				lvl: logger.WRN,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger logger.Logger) {
				logger.Info("infoLog")
			},
			want: "infoLog",
		},
		{
			name: "no debug",
			fields: fields{
				lvl: logger.INF,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger logger.Logger) {
				logger.Debug("debugLog")
			},
			want: "debugLog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.NewStdLog(logger.WithWriter(tt.fields.w), logger.WithLevel(tt.fields.lvl))
			tt.exec(l)
			if strings.Contains(tt.fields.w.String(), tt.want) {
				t.Error("out put contains unwanted message")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		want    logger.Level
		wantErr bool
	}{
		{
			name: "Disabled",
			args: args{
				level: "Disabled",
			},
			want:    logger.DSB,
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				level: "Error",
			},
			want:    logger.ERR,
			wantErr: false,
		},
		{
			name: "Warning",
			args: args{
				level: "Warning",
			},
			want:    logger.WRN,
			wantErr: false,
		},
		{
			name: "Info",
			args: args{
				level: "Info",
			},
			want:    logger.INF,
			wantErr: false,
		},
		{
			name: "Debug",
			args: args{
				level: "Debug",
			},
			want:    logger.DBG,
			wantErr: false,
		},
		{
			name: "Parse unknown level",
			args: args{
				level: "unknown",
			},
			want:    logger.INF,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logger.ParseLevel(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLevel() got = %v, want %v", got, tt.want)
			}
		})
	}
}
