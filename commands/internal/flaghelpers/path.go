package flaghelpers

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Path string

func (f *Path) UnmarshalFlag(value string) error {
	if value == "" {
		return nil
	}

	matches, err := filepath.Glob(value)
	if err != nil {
		return fmt.Errorf("failed to expand path '%s': %s", value, err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("path '%s' does not exist", value)
	}

	if len(matches) > 1 {
		return fmt.Errorf("path '%s' resolves to multiple entries: %s", value, strings.Join(matches, ", "))
	}

	*f = Path(matches[0])
	return nil
}
