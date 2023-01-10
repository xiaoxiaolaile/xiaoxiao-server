package core

type WatchBolt map[string][]func(old, now, key string) map[string]interface{}

var WatchBoltMap map[string]WatchBolt

type BucketJs struct {
	Bucket BoltBucket
}

func InitWatch() {
	WatchBoltMap = make(map[string]WatchBolt)
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
	old := b.Get(key, "")
	now := value
	var edit interface{}
	var editOk bool

	for _, f := range WatchBoltMap[b.Name()]["*"] {
		result := f(old, now, key)
		if s, ok := result["now"]; ok {
			edit = s
			editOk = ok
		}
	}

	for _, f := range WatchBoltMap[b.Name()][key] {
		result := f(old, now, key)
		if s, ok := result["now"]; ok {
			edit = s
			editOk = ok
		}
	}
	_ = b.Bucket.Set(key, value)
	if editOk {
		_ = b.Bucket.Set(key, edit)
	}
}

func (b *BucketJs) Delete(key string) {
	_ = b.Bucket.Set(key, "")
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

func (b *BucketJs) Watch(key string, f func(old, now, key string) map[string]interface{}) {

	if _, ok := WatchBoltMap[b.Name()]; !ok {
		WatchBoltMap[b.Name()] = make(WatchBolt)
	}

	WatchBoltMap[b.Name()][key] = append(WatchBoltMap[b.Name()][key], f)

}
