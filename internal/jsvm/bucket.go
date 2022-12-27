package jsvm

type BucketJs struct {
	Bucket BoltBucket
}

func (b *BucketJs) Get(key, defaultValue string) string {
	v := b.Bucket.GetString(key)
	if len(v) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue
		}
		return ""
	}
	return v
}

func (b *BucketJs) Set(key, value string) {
	_ = b.Bucket.Set(key, value)
}
func (b *BucketJs) Keys() []string {
	var ss []string
	b.Bucket.Foreach(func(k, _ []byte) error {
		ss = append(ss, string(k))
		return nil
	})
	return ss
}
func (b *BucketJs) DeleteAll() {
	var ss []string
	b.Bucket.Foreach(func(k, _ []byte) error {
		ss = append(ss, string(k))
		return nil
	})

	for _, s := range ss {
		_ = b.Bucket.Set(s, "")
	}

}
func (b *BucketJs) Name() string {
	return string(b.Bucket)
}
