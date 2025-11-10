package registry

import (
	"ddd/internal/serde"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"reflect"
)

type ctor func(payload []byte) (any, error)

type typeStore struct {
	serde.Serder
	ermu  sync.RWMutex
	items map[string]ctor
}

func (r *typeStore) Register(item any) {
	t := reflect.TypeOf(item)
	slog.Info("event registered", "type", t.Name())

	if t.Kind() != reflect.Struct && t.Kind() != reflect.Interface {
		panic("register: registered type must be struct or interface")
	}
	ctor := func(payload []byte) (any, error) {

		vt := reflect.New(t).Interface()
		if err := r.Deserialize(payload, vt); err != nil {
			return nil, fmt.Errorf("registry: %w", err)
		}
		return vt, nil
		//	var val V

	}
	r.ermu.Lock()
	r.items[TypeNameFrom(item)] = ctor
	r.ermu.Unlock()

}

func TypeNameFrom(e any) string {
	if strev, ok := e.(fmt.Stringer); ok {
		return strev.String()
	}

	t := reflect.TypeOf(e)
	switch t.Kind() {
	case reflect.Struct:
		return t.Name()
	case reflect.Pointer:
		return t.Elem().Name()
	default:
		panic("unsupported type")

		//	json.Marshal()
	}
}
func (r *typeStore) TypeExists(tname string) bool {
	r.ermu.RLock()
	defer r.ermu.RUnlock()
	if _, ok := r.items[tname]; ok {
		return true
	}
	return false
}

func (r *typeStore) GetType(tname string, b []byte) (any, error) {
	r.ermu.RLock()
	defer r.ermu.RUnlock()
	if ct, ok := r.items[tname]; ok {

		tt, err := ct(b)
		if err != nil {
			return tt, fmt.Errorf("registry: %w", err)
		}
		return tt, nil
	}
	slog.Error("registry: no type found", "type", tname)
	panic("unrecovered")
	//return nil, fmt.Errorf(")
}

func New(s serde.Serder) *typeStore {
	return &typeStore{items: make(map[string]ctor), Serder: s}
}

type TypeRegistry interface {
	serde.Serder
	TypeExists(tname string) bool
	GetType(tname string, b []byte) (any, error)
	Register(item any)
}

func MetaFromType[T any]() (aname string, bctx string) {
	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		panic("T must be a struct")
	}
	aname = t.Name()
	sep := strings.Split(t.PkgPath(), "/")
	bctx = sep[len(sep)-1]
	return
}
