package upload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/facebookgo/symwalk"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// ScanFolder return the list of Items{} to be uploaded. It scans the folder and skip
// non allowed files (includePatterns & excludePattens).
func (job *UploadFolderJob) ScanFolder(logger log.Logger) ([]FileItem, error) {
	var result []FileItem
	err := symwalk.Walk(job.SourceFolder, job.getItemToUploadFn(&result, logger))
	return result, err
}

func (job *UploadFolderJob) getItemToUploadFn(reqs *[]FileItem, logger log.Logger) filepath.WalkFunc {
	return func(fp string, fi os.FileInfo, errP error) error {
		if fi == nil {
			return nil
		}

		relativePath := RelativePath(job.SourceFolder, fp)

		// If a directory is excluded, skip it!
		if fi.IsDir() {
			if job.Filter.IsExcluded(relativePath) {
				logger.Infof("Not allowed by config: %s: skipping directory...", fp)
				return filepath.SkipDir
			}
			return nil
		}

		// check if the item should be uploaded given the include and exclude patterns in the
		// configuration file. It uses relative Path from the source folder Path to facilitate
		// then set up of includePatterns and excludePatterns.

		if !job.Filter.IsAllowed(relativePath) {
			logger.Infof("Not allowed by config: %s: skipping file...", fp)
			return nil
		}

		// check completed uploads db for previous uploads
		isAlreadyUploaded, err := job.FileTracker.IsAlreadyUploaded(fp)
		if err != nil {
			logger.Error(err)
		} else if isAlreadyUploaded {
			logger.Debugf("Already uploaded: %s: skipping file...", fp)
			return nil
		}

		logger.Infof("File '%s' will be uploaded to album '%s'.", fp, job.albumName(relativePath))

		// set file upload Options depending on folder upload Options
		*reqs = append(*reqs, FileItem{
			Path:      fp,
			AlbumName: job.albumName(relativePath),
		})
		return nil
	}
}

// RelativePath returns a path relative to the base.
// If a relative path could not be calculated or it contains ' ../`,
// returns the original path.
func RelativePath(base string, path string) string {
	rp, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	if strings.HasPrefix(rp, "../") {
		return path
	}
	return rp
}
