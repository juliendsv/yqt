package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func main() {
	dt := time.Now()
	dtnow := dt.Format("01-02-2006 15:04:05")

	err := yqEval("./examples/example1.yaml", fmt.Sprintf(".my.test.time = \"%s\"", dtnow))
	if err != nil {
		fmt.Println(err)
	}

	err = yqEval("./examples/example1.yaml", fmt.Sprintf(".my.test.count = %d", rand.IntN(100)))
	if err != nil {
		fmt.Println(err)
	}
}

func yqEval(yamlFile, expression string) (Error error) {
	var err error

	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(yamlFile)
	outFile, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}

	printerWriter := yqlib.NewSinglePrinterWriter(outFile)
	defer func() {
		if Error == nil {
			Error = writeInPlaceHandler.FinishWriteInPlace(true)
		}
	}()

	decoder, err := configureDecoder()
	if err != nil {
		return err
	}

	encoder, err := configureEncoder()
	if err != nil {
		return err
	}

	printer := yqlib.NewPrinter(encoder, printerWriter)

	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()

	err = allAtOnceEvaluator.EvaluateFiles(expression, []string{yamlFile}, printer, decoder)
	if err != nil {
		return err
	}

	return err
}

func configureDecoder() (yqlib.Decoder, error) {
	format, err := yqlib.FormatFromString("yaml")
	if err != nil {
		return nil, err
	}

	yqlibDecoder := format.DecoderFactory()
	if yqlibDecoder == nil {
		return nil, fmt.Errorf("no support for %s input format", "yaml")
	}
	return yqlibDecoder, nil
}

func configureEncoder() (yqlib.Encoder, error) {
	yqlibOutputFormat, err := yqlib.FormatFromString("yaml")
	if err != nil {
		return nil, err
	}

	yqlib.ConfiguredYamlPreferences.ColorsEnabled = false

	encoder := yqlibOutputFormat.EncoderFactory()

	if encoder == nil {
		return nil, fmt.Errorf("no support for %s output format", "yaml")
	}
	return encoder, err
}
