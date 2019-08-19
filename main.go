package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skillian/argparse"
	"github.com/skillian/errors"
	"github.com/skillian/logging"
)

var (
	logger = logging.GetLogger("preify")

	parser = argparse.MustNewArgumentParser(
		argparse.Description(
			"Apply a suffix to the given file or folder with "+
				"the date in it so the same folder can be "+
				"archived more than once into the same "+
				"archive"))

	includeTimeArg = parser.MustAddArgument(
		argparse.OptionStrings("-t", "--include-time"),
		argparse.Dest("includeTime"),
		argparse.Action("store_true"),
		argparse.Type(argparse.Bool),
		argparse.Default(false),
		argparse.Help("Include the time in the filename"))

	modTimeArg = parser.MustAddArgument(
		argparse.OptionStrings("-m", "--mod-time"),
		argparse.Action("store_true"),
		argparse.Help("use the file's modification time instead of "+
			"the current time"))

	logLevelArg = parser.MustAddArgument(
		argparse.OptionStrings("-L", "--log-level"),
		argparse.Dest("logLevel"),
		argparse.Action("store"),
		argparse.Default(logging.WarnLevel),
		argparse.Choices(
			argparse.Choice{Key: "debug", Value: logging.DebugLevel},
			argparse.Choice{Key: "info", Value: logging.InfoLevel},
			argparse.Choice{Key: "warn", Value: logging.WarnLevel},
			argparse.Choice{Key: "error", Value: logging.ErrorLevel}),
		argparse.Help("Set the logging level"))

	printOnlyArg = parser.MustAddArgument(
		argparse.OptionStrings("-p", "--print-only"),
		argparse.Action("store_true"),
		argparse.Help("Only print the result"))

	filenameArg = parser.MustAddArgument(
		argparse.Dest("filename"),
		argparse.Action("store"),
		argparse.Required,
		argparse.Help("filename to \"preify\""))

	ns = parser.MustParseArgs()

	includeTime = ns.MustGet(includeTimeArg).(bool)

	loggingLevel = ns.MustGet(logLevelArg).(logging.Level)

	modTime = ns.MustGet(modTimeArg).(bool)

	printOnly = ns.MustGet(printOnlyArg).(bool)

	filename = ns.MustGet(filenameArg).(string)
)

func init() {
	h := new(logging.ConsoleHandler)
	h.SetFormatter(logging.DefaultFormatter{})
	h.SetLevel(loggingLevel)
	logger.AddHandler(h)
	logger.SetLevel(loggingLevel)
}

func main() {
	fi, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.LogErr(
				errors.ErrorfWithCause(
					err,
					"File %q does not exist", filename))
		} else {
			logger.LogErr(err)
		}
		os.Exit(-1)
	}
	filename, err = filepath.Abs(filename)
	if err != nil {
		logger.LogErr(err)
		os.Exit(-1)
	}
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	base := filepath.Base(filename[:len(filename)-len(ext)])
	var t time.Time
	if modTime {
		t = fi.ModTime()
	} else {
		t = time.Now()
	}
	format := "2006-01-02"
	if includeTime {
		format = "2006-01-02_15-04-05"
	}
	date := t.Format(format)
	newname := strings.Join([]string{base, ".pre-", date, ext}, "")
	if printOnly {
		fmt.Println(filename)
		return
	}
	if err = os.Rename(filename, filepath.Join(dir, newname)); err != nil {
		logger.LogErr(
			errors.ErrorfWithCause(
				err, "error while renaming file"))
		os.Exit(-1)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
