package logger

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLevel_String(t *testing.T) {
	type tcase struct {
		l    Level
		want string
	}
	tests := map[string]tcase{
		"Disabled": {DSB, "Disabled"},
		"Error":    {ERR, "Error"},
		"Info":     {INF, "Info"},
		"Warning":  {WRN, "Warning"},
		"Debug":    {DBG, "Debug"},
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
		got := NewStdLog()
		if got.lvl != INF {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}
		if got.err.Writer() != os.Stderr || got.inf.Writer() != os.Stderr ||
			got.wrn.Writer() != os.Stderr || got.dbg.Writer() != os.Stderr {
			t.Errorf("all default writers should be os.Stderr")
		}
		if !reflect.TypeOf(got).Implements(reflect.TypeOf((*Logger)(nil)).Elem()) {
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
		got := NewStdLog(WithWriter(tb))
		got.Info(strToLog)

		if !strings.Contains(tb.String(), strToLog) {
			t.Errorf("expected %s but got %s", strToLog, tb.String())
		}

		if got.lvl != INF {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}
		if got.err.Writer() != tb || got.inf.Writer() != tb ||
			got.wrn.Writer() != tb || got.dbg.Writer() != tb {
			t.Errorf("all default writers should be os.Stderr")
		}
	})
}

func TestStdLog_Print(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log := NewStdLog(WithWriter(buf), WithLevel(DBG))
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
		lvl Level
		w   *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		exec   func(logger Logger)
		want   string
	}{
		{
			name: "no error",
			fields: fields{
				lvl: DSB,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger Logger) {
				logger.Error("errorLog")
			},
			want: "errorLog",
		},
		{
			name: "no warn",
			fields: fields{
				lvl: ERR,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger Logger) {
				logger.Warn("warnLog")
			},
			want: "warnLog",
		},
		{
			name: "no info",
			fields: fields{
				lvl: WRN,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger Logger) {
				logger.Info("infoLog")
			},
			want: "infoLog",
		},
		{
			name: "no debug",
			fields: fields{
				lvl: INF,
				w:   bytes.NewBuffer(nil),
			},
			exec: func(logger Logger) {
				logger.Debug("debugLog")
			},
			want: "debugLog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewStdLog(WithWriter(tt.fields.w), WithLevel(tt.fields.lvl))
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
		want    Level
		wantErr bool
	}{
		{
			name: "Disabled",
			args: args{
				level: "Disabled",
			},
			want:    DSB,
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				level: "Error",
			},
			want:    ERR,
			wantErr: false,
		},
		{
			name: "Warning",
			args: args{
				level: "Warning",
			},
			want:    WRN,
			wantErr: false,
		},
		{
			name: "Info",
			args: args{
				level: "Info",
			},
			want:    INF,
			wantErr: false,
		},
		{
			name: "Debug",
			args: args{
				level: "Debug",
			},
			want:    DBG,
			wantErr: false,
		},
		{
			name: "Parse unknown level",
			args: args{
				level: "unknown",
			},
			want:    INF,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.args.level)
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
