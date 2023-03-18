package dredis

func NewMutex(key string, options ...redsync.Option) *redsync.Mutex {
	key = RealKey(key)
	return cli.RedsSync.NewMutex(key, options...)
}
