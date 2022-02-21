package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keyval-dev/offsets-tracker/target"
	"io/fs"
	"io/ioutil"
	"strings"
)

const fileName = "/tmp/offset_results.json"

func WriteResults(results ...*target.Result) error {
	var offsets TrackedOffsets
	for _, r := range results {
		offsets.Data = append(offsets.Data, convertResult(r))
	}

	jsonData, err := json.Marshal(&offsets)
	if err != nil {
		return err
	}

	var prettyJson bytes.Buffer
	err = json.Indent(&prettyJson, jsonData, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, prettyJson.Bytes(), fs.ModePerm)
}

func convertResult(r *target.Result) TrackedLibrary {
	tl := TrackedLibrary{
		Name: r.ModuleName,
	}

	offsetsMap := make(map[string][]VersionedOffset)
	for _, vr := range r.ResultsByVersion {
		for _, od := range vr.OffsetData.DataMembers {
			key := fmt.Sprintf("%s,%s", od.StructName, od.Field)
			offsetsMap[key] = append(offsetsMap[key], VersionedOffset{
				Offset:  od.Offset,
				Version: vr.Version,
			})
		}
	}

	for key, offsets := range offsetsMap {
		parts := strings.Split(key, ",")
		tl.DataMembers = append(tl.DataMembers, TrackedDataMember{
			Struct:  parts[0],
			Field:   parts[1],
			Offsets: offsets,
		})
	}

	return tl
}
