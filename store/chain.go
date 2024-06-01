package store

import "encoding/binary"

const ChainBucket = "worker"

func SetChainUrl(id string) error {
	key := []byte("ChainUrl")
	val := []byte(id)
	return SealSave(ChainBucket, key, val)
}

func GetChainUrl() (string, error) {
	val, err := SealGet(ChainBucket, []byte("ChainUrl"))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func SetChainBlock(id uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)

	key := []byte("ChainBlock")
	val := []byte(b)
	return SealSave(ChainBucket, key, val)
}

func GetChainBlock() (uint64, error) {
	val, err := SealGet(ChainBucket, []byte("ChainBlock"))
	if err != nil {
		if err.Error() == "key not found" || err.Error() == "bucket not found" {
			return 0, nil
		}
		return 0, err
	}

	num := binary.LittleEndian.Uint64(val)
	return num, nil
}
