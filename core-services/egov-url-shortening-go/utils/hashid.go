package utils

import (
	"egov-url-shortening-go/config"
	"github.com/speps/go-hashids/v2"
)

// HashIDConverter provides HashID encoding and decoding functionality
type HashIDConverter struct {
	hashids *hashids.HashID
}

// NewHashIDConverter creates a new HashIDConverter instance
func NewHashIDConverter(cfg *config.HashIDsConfig) (*HashIDConverter, error) {
	hd := hashids.NewData()
	hd.Salt = cfg.Salt
	hd.MinLength = cfg.MinLength
	
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	
	return &HashIDConverter{
		hashids: h,
	}, nil
}

// CreateHashStringForID converts an ID to a hash string
func (h *HashIDConverter) CreateHashStringForID(id int64) (string, error) {
	return h.hashids.EncodeInt64([]int64{id})
}

// GetIDForString converts a hash string back to an ID
func (h *HashIDConverter) GetIDForString(hashString string) (int64, error) {
	ids, err := h.hashids.DecodeInt64WithError(hashString)
	if err != nil {
		return 0, err
	}
	
	if len(ids) != 1 {
		return 0, nil // Return 0 for invalid hash (matches Java null behavior)
	}
	
	return ids[0], nil
}