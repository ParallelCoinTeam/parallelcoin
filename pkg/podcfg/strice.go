package podcfg

// Strice is a wrapper around byte slices to enable optional security features and possibly better performance for
// bulk comparison and editing. There isn't any extensive editing primitives for this purpose,
type Strice []byte

// S returns the underlying bytes converted into string
func (s *Strice) S() string {
	return string(*s)
}

// E returns the byte at the requested index in the string
func (s *Strice) E(elem int) byte {
	if s.Len() > elem {
		return (*s)[elem]
	}
	return 0
}

// Len returns the length of the string in bytes
func (s *Strice) Len() int {
	return len(*s)
}

// Equal returns true if two Strices are equal in both length and content
func (s *Strice) Equal(sb *Strice) bool {
	if s.Len() == sb.Len() {
		for i := range *s {
			if s.E(i) != sb.E(i) {
				return false
			}
		}
		return true
	}
	return false
}

// Cat two Strices together
func (s *Strice) Cat(sb *Strice) *Strice {
	*s = append(*s, *sb...)
	return s
}

// Find returns true if a match of a substring is found and if found, the position in the first string that the second
// string starts, the number of matching characters from the start of the search Strice, or -1 if not found.
//
// You specify a minimum length match and it will trawl through it systematically until it finds the first match of the
// minimum length.
func (s *Strice) Find(sb *Strice, minLengthMatch int) (found bool, extent, pos int) {
	// can't be a substring if it's longer
	if sb.Len() > s.Len() {
		return
	}
	for pos = range *s {
		// if we find a match, grab onto it
		if s.E(pos) == sb.E(pos) {
			extent++
			// this exhaustively searches for a match between the two strings, but we do not restrict the match to the
			// minimum, maximising the ways this function can be used for simple position tests and editing
			for srchPos := 1; srchPos < sb.Len() || srchPos+pos < s.Len(); srchPos++ {
				// the first element is skipped
				if s.E(srchPos+pos) != sb.E(srchPos) {
					break
				}
				extent++
			}
			// the above loop ends when the bytes stop matching, then if it is under the minimum length requested, it
			// continues. Note that we are not mutating `i` so it iterates for a match comprehensively.
			if extent < minLengthMatch {
				// reset the extent
				extent = 0
			} else {
				break
			}
		}
	}
	return
}

// HasPrefix returns true if the given string forms the beginning of the current string
func (s *Strice) HasPrefix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == 0 {
			return true
		}
	}
	return false
}

// HasSuffix returns true if the given string forms the ending of the current string
func (s *Strice) HasSuffix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == s.Len()-sb.Len()-1 {
			return true
		}
	}
	return false
}

// Dup copies a string and returns it
func (s *Strice) Dup() *Strice {
	ns := make(Strice, s.Len())
	copy(ns, *s)
	return &ns
}

// Wipe zeroes the bytes of a string
func (s *Strice) Wipe() {
	for i := range *s {
		(*s)[i] = 0
	}
}

// Split the string by a given cutset
func (s *Strice) Split(cutset string) (out []*Strice) {
	// convert immutable string type to Strice bytes
	c := Strice(cutset)
	// need the pointer to call the methods
	cs := &c
	// copy the bytes so we can guarantee the original is unmodified
	cp := s.Dup()
	for {
		// locate the next instance of the cutset
		found, _, pos := s.Find(cp, cp.Len())
		if found {
			// add the found section to the return slice
			before := (*s)[:pos+cp.Len()]
			out = append(out, &before)
			// trim off the prefix and cutslice from the working copy
			*cs = (*cs)[pos+cp.Len():]
			// continue to search for more instances of the cutset
			continue
		} else {
			// once we get not found, the searching is over and whatever we have, we return
			break
		}
	}
	return
}
