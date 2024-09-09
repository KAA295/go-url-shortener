package sl

import "log/slog"

func Err(err error) slog.Attr {
	//Ключ error, значение - ошибка в текстовом формате
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
