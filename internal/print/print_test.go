package print

import (
	"bytes"
	"errors"
	"testing"
)

type StubPrintable struct {
	items []struct {
		id   int
		name string
	}
	returnErr error
}

func (s *StubPrintable) ForEach(fn func(id int, name string) error) error {
	if s.returnErr != nil {
		return s.returnErr
	}
	for _, item := range s.items {
		if err := fn(item.id, item.name); err != nil {
			return err
		}
	}
	return nil
}

type errWriter struct{}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("Ошибка записи")
}

func TestPrinter_Print(t *testing.T) {
	tests := []struct {
		name      string
		source    *StubPrintable
		useErrW   bool
		wantErr   bool
		wantPrint string
	}{
		{
			name: "Успешная печать",
			source: &StubPrintable{
				items: []struct {
					id   int
					name string
				}{
					{1, "шлем"},
					{2, "самокат"},
				},
			},
			useErrW:   false,
			wantErr:   false,
			wantPrint: "Шлем в наличии, инвентарный номер: 1\nСамокат в наличии, инвентарный номер: 2\n",
		},
		{
			name: "Ошибка источника данных (ForEach)",
			source: &StubPrintable{
				returnErr: errors.New("Сбой базы данных"),
			},
			useErrW: false,
			wantErr: true,
		},
		{
			name: "Ошибка записи (io.Writer)",
			source: &StubPrintable{
				items: []struct {
					id   int
					name string
				}{{1, "шлем"}},
			},
			useErrW: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := NewPrinter(tt.source)

			var buf bytes.Buffer
			var w interface {
				Write(p []byte) (n int, err error)
			} = &buf
			if tt.useErrW {
				w = errWriter{}
			}

			err := printer.Print(w)

			if (err != nil) != tt.wantErr {
				t.Errorf("Ожидалось наличие ошибки: %v, получена ошибка: %v", tt.wantErr, err)
			}

			if !tt.wantErr && !tt.useErrW {
				if got := buf.String(); got != tt.wantPrint {
					t.Errorf("Ожидался вывод:\n%q\nПолучено:\n%q", tt.wantPrint, got)
				}
			}
		})
	}
}
