package helper

import "fmt"

func ParseSize(size int64) string {
	const KB = 1024
	const MB = KB * KB
	const GB = MB * KB

	if size < KB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}

	if size > KB && size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}

	if size > KB && size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	}

	if size > GB {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}

	return "0 KB"
}
