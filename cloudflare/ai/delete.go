package ai

import (
	"github.com/syumai/workers/internal/jsutil"
)

// Delete deletes key-value pair specified by the key.
//   - if a network error happens, returns error.
func (ns *Ai) Delete(key string) error {
	p := ns.instance.Call("delete", key)
	_, err := jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}
