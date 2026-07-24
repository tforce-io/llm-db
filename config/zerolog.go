// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func InitZerolog() zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger
	return logger
}
