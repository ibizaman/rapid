// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "fmt"

// ID returns its argument as is. ID is a helper for use with [SliceOfDistinct] and similar functions.
func ID[V any](v V) V {
	return v
}

func SliceOf[E any](elem *Generator[E]) *Generator[[]E] {
	return SliceOfN(elem, -1, -1)
}

func SliceOfN[E any](elem *Generator[E], minLen int, maxLen int) *Generator[[]E] {
	assertValidRange(minLen, maxLen)

	return newGenerator[[]E](&sliceGen[E, struct{}]{
		minLen: minLen,
		maxLen: maxLen,
		elem:   elem,
	})
}

func SliceOfDistinct[E any, K comparable](elem *Generator[E], keyFn func(E) K) *Generator[[]E] {
	return SliceOfNDistinct(elem, -1, -1, keyFn)
}

func SliceOfNDistinct[E any, K comparable](elem *Generator[E], minLen int, maxLen int, keyFn func(E) K) *Generator[[]E] {
	assertValidRange(minLen, maxLen)

	return newGenerator[[]E](&sliceGen[E, K]{
		minLen: minLen,
		maxLen: maxLen,
		elem:   elem,
		keyFn:  keyFn,
	})
}

type sliceGen[E any, K comparable] struct {
	minLen int
	maxLen int
	elem   *Generator[E]
	keyFn  func(E) K
}

func (g *sliceGen[E, K]) String() string {
	if g.keyFn == nil {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("SliceOf(%v)", g.elem)
		} else {
			return fmt.Sprintf("SliceOfN(%v, minLen=%v, maxLen=%v)", g.elem, g.minLen, g.maxLen)
		}
	} else {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("SliceOfDistinct(%v, key=%T)", g.elem, g.keyFn)
		} else {
			return fmt.Sprintf("SliceOfNDistinct(%v, minLen=%v, maxLen=%v, key=%T)", g.elem, g.minLen, g.maxLen, g.keyFn)
		}
	}
}

func (g *sliceGen[E, K]) value(t *T) []E {
	repeat := newRepeat(g.minLen, g.maxLen, -1, g.elem.String())

	var seen map[K]struct{}
	if g.keyFn != nil {
		seen = make(map[K]struct{}, repeat.avg())
	}

	sl := make([]E, 0, repeat.avg())
	for repeat.more(t.s) {
		e := g.elem.value(t)
		if g.keyFn == nil {
			sl = append(sl, e)
		} else {
			k := g.keyFn(e)
			if _, ok := seen[k]; ok {
				repeat.reject()
			} else {
				seen[k] = struct{}{}
				sl = append(sl, e)
			}
		}
	}

	return sl
}

func MapOf[K comparable, V any](key *Generator[K], val *Generator[V]) *Generator[map[K]V] {
	return MapOfN(key, val, -1, -1)
}

func MapOfN[K comparable, V any](key *Generator[K], val *Generator[V], minLen int, maxLen int) *Generator[map[K]V] {
	assertValidRange(minLen, maxLen)

	return newGenerator[map[K]V](&mapGen[K, V]{
		minLen: minLen,
		maxLen: maxLen,
		key:    key,
		val:    val,
	})
}

func MapOfValues[K comparable, V any](val *Generator[V], keyFn func(V) K) *Generator[map[K]V] {
	return MapOfNValues(val, -1, -1, keyFn)
}

func MapOfNValues[K comparable, V any](val *Generator[V], minLen int, maxLen int, keyFn func(V) K) *Generator[map[K]V] {
	assertValidRange(minLen, maxLen)

	return newGenerator[map[K]V](&mapGen[K, V]{
		minLen: minLen,
		maxLen: maxLen,
		val:    val,
		keyFn:  keyFn,
	})
}

type mapGen[K comparable, V any] struct {
	minLen int
	maxLen int
	key    *Generator[K]
	val    *Generator[V]
	keyFn  func(V) K
}

func (g *mapGen[K, V]) String() string {
	if g.key != nil {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("MapOf(%v, %v)", g.key, g.val)
		} else {
			return fmt.Sprintf("MapOfN(%v, %v, minLen=%v, maxLen=%v)", g.key, g.val, g.minLen, g.maxLen)
		}
	} else {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("MapOfValues(%v, key=%T)", g.val, g.keyFn)
		} else {
			return fmt.Sprintf("MapOfNValues(%v, minLen=%v, maxLen=%v, key=%T)", g.val, g.minLen, g.maxLen, g.keyFn)
		}
	}
}

func (g *mapGen[K, V]) value(t *T) map[K]V {
	label := g.val.String()
	if g.key != nil {
		label = g.key.String() + "," + label
	}

	repeat := newRepeat(g.minLen, g.maxLen, -1, label)

	m := make(map[K]V, repeat.avg())
	for repeat.more(t.s) {
		var k K
		var v V
		if g.key != nil {
			k = g.key.value(t)
			v = g.val.value(t)
		} else {
			v = g.val.value(t)
			k = g.keyFn(v)
		}

		if _, ok := m[k]; ok {
			repeat.reject()
		} else {
			m[k] = v
		}
	}

	return m
}
