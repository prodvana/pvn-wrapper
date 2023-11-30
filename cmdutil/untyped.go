package cmdutil

import "github.com/pkg/errors"

func GetOrCreateUntypedMap(m map[string]interface{}, key string) (map[string]interface{}, error) {
	if m[key] == nil {
		m[key] = map[string]interface{}{}
	}
	typed, ok := m[key].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("unexpected type for %s: %T", key, m[key])
	}
	return typed, nil
}
