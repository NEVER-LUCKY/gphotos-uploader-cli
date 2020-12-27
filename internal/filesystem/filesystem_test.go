package filesystem_test

import (
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filesystem"
)

func TestAbsolutePath(t *testing.T) {

	t.Run("WithAbsolutePaths", func(t *testing.T) {
		var absolutePathInputs = []struct {
			in  string
			out string
		}{
			{"/", "/"},
			{"/xyz", "/xyz"},
			{"/xyz/./abc", "/xyz/abc"},
			{"/xyz/../abc", "/abc"},
			{"/xyz/abc/..", "/xyz"},
			{"/xyz/../abc/..", "/"},
			{"/xyz/../..", "/"},
			{"/xyz///abc/..", "/xyz"},
		}

		for _, test := range absolutePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			if got != test.out {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, test.out, got)
			}
		}
	})

	t.Run("WithRelativePath", func(t *testing.T) {
		var relativePathInputs = []struct {
			in  string
			out string
		}{
			{"", ""},
			{"./", ""},
			{"xyz", "xyz"},
			{"xyz/./abc", "xyz/abc"},
			{"xyz/../abc", "abc"},
			{"xyz/abc/..", "xyz"},
			{"xyz/../abc/..", ""},
			{"xyz/../..", ".."},
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		for _, test := range relativePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			expected := path.Join(cwd, test.out)

			if got != expected {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, expected, got)
			}
		}
	})

	t.Run("WithTildePath", func(t *testing.T) {
		var tildePathInputs = []struct {
			in  string
			out string
		}{
			{"~", ""},
			{"~/", ""},
			{"~/xyz", "xyz"},
			{"~/xyz/./abc", "xyz/abc"},
			{"~/xyz/../abc", "abc"},
			{"~/xyz/abc/..", "xyz"},
			{"~/xyz/../abc/..", ""},
			{"~/xyz/../..", ".."},
			{"~/xyz/~/abc", "xyz/~/abc"},
		}

		usr, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		dir := usr.HomeDir

		for _, test := range tildePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			expected := path.Join(dir, test.out)

			if got != expected {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, expected, got)
			}
		}
	})
}

