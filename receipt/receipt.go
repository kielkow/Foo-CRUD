package receipt

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

// ReceiptDirectory uploads
var ReceiptDirectory string = filepath.Join("uploads")

// Receipt struct
type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

// GetReceipts function
func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	files, err := ioutil.ReadDir(ReceiptDirectory)

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		receipts = append(
			receipts, 
			Receipt{
				ReceiptName: f.Name(), 
				UploadDate: f.ModTime(),
			},
		)
	}

	return receipts, nil
}
