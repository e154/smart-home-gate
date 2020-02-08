// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package common

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/jinzhu/copier"
)

type CopyEngine string

const (
	JsonEngine = CopyEngine("json")
	GobEngine  = CopyEngine("gob")
)

func gobCopy(to, from interface{}) (err error) {
	buff := new(bytes.Buffer)
	if err = gob.NewEncoder(buff).Encode(from); err != nil {
		return
	}
	err = gob.NewDecoder(buff).Decode(to)
	return
}

func jsonCopy(to, from interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(from); err != nil {
		return
	}
	err = json.Unmarshal(b, to)
	return
}

func Copy(to, from interface{}, params ...CopyEngine) (err error) {

	if len(params) == 0 {
		err = copier.Copy(to, from)
		return
	}

	switch params[0] {
	case JsonEngine:
		err = jsonCopy(to, from)
	case GobEngine:
		err = gobCopy(to, from)
	}

	return
}
