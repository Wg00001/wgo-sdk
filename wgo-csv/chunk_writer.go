package wgo_csv

import (
	"encoding/csv"
	"errors"
	"fmt"
)

//Provides a function that can be used for fast chunk write into csv.writer.

// CsvRow defines an interface for objects that can be written to a CSV file.
// Types that implement this interface must provide a ToStringSlice method,
// which returns the data of the object as a slice of strings.
type CsvRow interface {
	ToStringSlice() []string
}

// ChunkWrite writes data to a CSV file in chunks to optimize memory usage.
//
// @param writer *csv.Writer The CSV writer to write data.
// @param data []T A slice of data to be written, where T implements CsvRow with a ToStringSlice() method.
// @param size Optional, maximum number of rows per chunk. Default is 2000.
//
// @return err Returns an error if something goes wrong, or nil if successful.
//
// This function writes data in batches to avoid excessive memory usage and flushes the writer after each batch.
func ChunkWrite[T CsvRow](writer *csv.Writer, data []T, size ...int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ChunkWrite Panic: %v", r)
		}
	}()
	if data == nil {
		return errors.New("data slice is NULL")
	}
	if writer == nil {
		return errors.New("csv writer is NULL")
	}
	pageSize := 2000
	if len(size) != 0 {
		pageSize = size[0]
	}
	cur := make([][]string, 0, pageSize+5)
	for i := 0; i < len(data); i++ {
		cur = append(cur, data[i].ToStringSlice())
		if (i+1)%pageSize == 0 {
			err := writer.WriteAll(cur)
			if err != nil {
				return err
			}
			writer.Flush()
			cur = cur[:0]
		}
	}
	if len(cur) > 0 {
		err := writer.WriteAll(cur)
		if err != nil {
			return err
		} //写入剩余数据
	}
	writer.Flush()
	return nil
}
