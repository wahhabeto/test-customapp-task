package logger

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts  PrettyHandlerOptions
	attrs []slog.Attr
	l     *stdLog.Logger
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	return &PrettyHandler{
		opts:  opts,
		l:     stdLog.New(out, "", 0),
		attrs: []slog.Attr{},
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.SlogOpts.Level.Level()
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	levelStr := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		levelStr = color.MagentaString(levelStr)
	case slog.LevelInfo:
		levelStr = color.BlueString(levelStr)
	case slog.LevelWarn:
		levelStr = color.YellowString(levelStr)
	case slog.LevelError:
		levelStr = color.RedString(levelStr)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		levelStr,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append(h.attrs, attrs...)
	return &PrettyHandler{
		opts:  h.opts,
		attrs: newAttrs,
		l:     h.l,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		opts:  h.opts,
		attrs: h.attrs,
		l:     h.l,
	}
}
