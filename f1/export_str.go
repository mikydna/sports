// Code generated by "stringer -type=ExportFormat -linecomment=true -output export_str.go"; DO NOT EDIT.

package f1

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ExportFormatUnknown-0]
	_ = x[ExportFormatJSON-1]
	_ = x[ExportFormatProto-2]
	_ = x[ExportFormatGob-3]
}

const _ExportFormat_name = "unknownjsonlpbgob"

var _ExportFormat_index = [...]uint8{0, 7, 12, 14, 17}

func (i ExportFormat) String() string {
	if i >= ExportFormat(len(_ExportFormat_index)-1) {
		return "ExportFormat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ExportFormat_name[_ExportFormat_index[i]:_ExportFormat_index[i+1]]
}
