package regources

type KeyValue struct {
	Key         string  `json:"key"`
	Type        Flag    `json:"type,omitempty"`
	Ui32Value   *uint32 `json:"ui32_value,omitempty"`
	Ui64Value   *uint64 `json:"ui64_value,omitempty"`
	StringValue *string `json:"string_value,omitempty"`
}
