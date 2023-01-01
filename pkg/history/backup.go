package history

import (
	"fmt"
	"io"
)

type Back struct {
	r io.Reader
	w io.Writer
}

func NewBackup(r io.Reader, w io.Writer) *Back { return &Back{r, w} }

func (b *Back) Backup() error {
	//	buf := make([]byte, 1024)
	//	for {
	//		n, err := b.r.Read(buf)
	//		if err == io.EOF {
	//			// there is no more data to read
	//			break
	//		}
	//		if err != nil {
	//			fmt.Println(err)
	//			continue
	//		}
	//		if n > 0 {
	//			fmt.Print(buf[:n])
	//		}
	//	}
	_, err := io.Copy(b.w, b.r)
	if err != nil {
		return fmt.Errorf("failed to backup file: %w", err)
	}

	return nil
}
