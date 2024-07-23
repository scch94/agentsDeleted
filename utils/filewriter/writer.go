package fileWriter

import (
	"context"
	"fmt"
	"os"

	"github.com/scch94/ins_log"
)

func WriteInAfile(ctx context.Context, text string, filename string, comment string) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "fileWriter")
	ins_log.Infof(ctx, "starting to write in the file %s", filename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening log file: %v", err)
		return err
	}
	defer file.Close()

	comment = "-- " + comment

	_, err = file.WriteString(fmt.Sprintf("%v \n\n%v \n", comment, text))
	if err != nil {
		ins_log.Errorf(ctx, "error writing to file: %v", err)
		return err
	}
	ins_log.Tracef(ctx, "end to wrtite in a file: ")

	return nil
}
