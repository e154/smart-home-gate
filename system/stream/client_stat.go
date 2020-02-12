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

package stream

type Stat struct {
	// total messages received
	received int64
	// total messages sent
	sent       int64
	requestMin float64
	requestMax float64
}

// increment counter
func (s *Stat) receivedInc() {
	s.received++
}

// increment counter
func (s *Stat) sentInc() {
	s.sent++
}

func (s *Stat) getReceived() int64 {
	return s.received
}

func (s *Stat) getSent() int64 {
	return s.sent
}
