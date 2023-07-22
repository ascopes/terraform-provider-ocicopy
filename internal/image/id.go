package image

import (
	"encoding/base64"
	"encoding/json"
)

type idDescriptor struct {
	SourceRegistry   string `json:"sg"`
	SourceRepository string `json:"sp"`
	SourceDigest     string `json:"sd"`
	TargetRegistry   string `json:"tg"`
	TargetRepository string `json:"tp"`
	TargetTag        string `json:"tt"`
}

func (desc idDescriptor) encode() string {
	bytes, err := json.Marshal(&desc)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}
