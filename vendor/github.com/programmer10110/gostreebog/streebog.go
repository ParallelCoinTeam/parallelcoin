package gostreebog

func s(x []byte) []byte {
	for i := 0; i < 64; i++ {
		x[i] = Sbox[x[i]]
	}
	return x
}

func p(x []byte) []byte {
	b := make([]byte, 64)
	for i := 0; i < 64; i++ {
		b[i] = x[Tau[i]]
	}
	return b
}

func l(x []byte) []byte {
	var v uint64
	for i := 0; i < 8; i++ {
		v = 0
		for k := 0; k < 8; k++ {
			for j := 0; j < 8; j++ {
				if (x[i*8+k])&(1<<uint(7-j)) != 0 {
					v ^= A[k*8+j]
				}
			}
		}
		for k := 0; k < 8; k++ {
			x[i*8+k] = byte((v & uint64(0xFF<<uint((7-k)*8))) >> uint((7-k)*8))
		}
	}
	return x
}

func xor512(a, b []byte) []byte {
	c := make([]byte, 64)
	for i := 0; i < 64; i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func k(K []byte, i int) []byte {
	K = xor512(K, C[i])
	K = s(K)
	K = p(K)
	K = l(K)
	return K
}

func e(K, m []byte) []byte {
	b := xor512(K, m)
	for i := 0; i < 12; i++ {
		b = s(b)
		b = p(b)
		b = l(b)
		K = k(K, i)
		b = xor512(b, K)
	}
	return b
}

func gN(N, h, m []byte) []byte {
	K := xor512(h, N)
	K = s(K)
	K = p(K)
	K = l(K)
	t := e(K, m)
	t = xor512(h, t)
	G := xor512(t, m)
	return G
}

func addModulo512(a, b []byte) []byte {
	var t uint
	t = 0
	c := make([]byte, 64)
	for i := 63; i >= 0; i-- {
		t = uint(a[i]) + uint(b[i]) + (t >> 8)
		c[i] = byte(t & 0xFF)
	}
	return c
}

func hash(data, IV []byte, t string) []byte {
	//Stage 1
	h := IV
	N := make([]byte, 64)
	Sigma := make([]byte, 64)
	v512 := make([]byte, 64)
	v512[62] = 0x2
	v0 := make([]byte, 64)

	//Stage 2
	length := len(data) * 8
	for length >= 512 {
		m := data[length/8-64:]
		h = gN(N, h, m)
		N = addModulo512(N, v512)
		Sigma = addModulo512(Sigma, m)
		data = data[:length/8-64]
		length -= 512
	}

	//Stage 3
	m := make([]byte, 64-length/8)
	m = append(m, data...)
	m[63-length/8] |= (1 << uint(length&0x7))
	h = gN(N, h, m)
	v512[63] = byte(length)
	v512[62] = byte(length >> 8)
	N = addModulo512(N, v512)
	Sigma = addModulo512(Sigma, m)
	h = gN(v0, h, N)
	h = gN(v0, h, Sigma)
	if t == "256" {
		return h[:32]
	}
	return h
}

func Hash(data []byte, t string) []byte {
	IV := make([]byte, 64)
	if t == "256" {
		for i := 0; i < 64; i++ {
			IV[i] = 0x1
		}
	}
	return hash(data, IV, t)
}
