package store

const Ink = "ink"

func SetInkCodeBlock(id string, v string) error {
	key := []byte(id)
	val := []byte(v)
	return SealSave(Ink, key, val)
}

func GetInkCodeBlock(id string) (string, error) {
	val, err := SealGet(Ink, []byte(id))
	if err != nil {
		if err.Error() == "key not found" || err.Error() == "bucket not found" {
			return "", nil
		}
		return "", err
	}

	return string(val), nil
}
