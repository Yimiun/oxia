package oxia

import (
	"oxia/common"
)

type memoryClient struct {
	clock common.Clock
	data  map[string]GetResult
}

func NewMemoryClient() AsyncClient {
	return newMemoryClientWithClock(common.SystemClock())
}

func newMemoryClientWithClock(clock common.Clock) AsyncClient {
	return &memoryClient{
		clock: clock,
		data:  make(map[string]GetResult),
	}
}

func (c *memoryClient) Close() error {
	return nil
}

func (c *memoryClient) Put(key string, payload []byte, expectedVersion *int64) <-chan PutResult {
	ch := make(chan PutResult, 1)
	now := c.clock.NowMillis()
	if value, ok := c.data[key]; ok {
		if expectedVersion != nil && *expectedVersion != value.Stat.Version {
			ch <- PutResult{Err: ErrorBadVersion}
		} else {
			value.Payload = payload
			value.Stat.Version = value.Stat.Version + 1
			value.Stat.ModifiedTimestamp = now
			ch <- PutResult{
				Stat: value.Stat,
			}
		}
	} else {
		if expectedVersion != nil && *expectedVersion != VersionNotExists {
			ch <- PutResult{Err: ErrorBadVersion}
		} else {
			value = GetResult{
				Payload: payload,
				Stat: Stat{
					Version:           1,
					CreatedTimestamp:  now,
					ModifiedTimestamp: now,
				},
			}
			c.data[key] = value
			ch <- PutResult{
				Stat: value.Stat,
			}
		}
	}
	close(ch)
	return ch
}

func (c *memoryClient) Delete(key string, expectedVersion *int64) <-chan error {
	ch := make(chan error, 1)
	if value, ok := c.data[key]; ok {
		if expectedVersion != nil && *expectedVersion != value.Stat.Version {
			ch <- ErrorBadVersion
		} else {
			delete(c.data, key)
		}
	} else {
		ch <- ErrorKeyNotFound
	}
	close(ch)
	return ch
}

func (c *memoryClient) DeleteRange(minKeyInclusive string, maxKeyExclusive string) <-chan error {
	for key := range c.data {
		if minKeyInclusive <= key && key < maxKeyExclusive {
			delete(c.data, key)
		}
	}
	ch := make(chan error, 1)
	close(ch)
	return ch
}

func (c *memoryClient) Get(key string) <-chan GetResult {
	ch := make(chan GetResult, 1)
	if value, ok := c.data[key]; ok {
		ch <- value
	} else {
		ch <- GetResult{
			Err: ErrorKeyNotFound,
		}
	}
	close(ch)
	return ch
}

func (c *memoryClient) GetRange(minKeyInclusive string, maxKeyExclusive string) <-chan GetRangeResult {
	ch := make(chan GetRangeResult, 1)
	result := make([]string, 0)
	for key := range c.data {
		if minKeyInclusive <= key && key < maxKeyExclusive {
			result = append(result, key)
		}
	}
	ch <- GetRangeResult{
		Keys: result,
	}
	close(ch)
	return ch
}