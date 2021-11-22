import sys, math

q = [ '\0\0\0\0\0\0\0\0' ] * 930

gamma=2.5
high_cutoff=0x30
low_cutoff=0x08

power_factor = 0xff00 / (0xff**gamma)
lut = [ int(x ** gamma * power_factor) for x in xrange(256) ]
for i, v in enumerate(lut):
	if v < low_cutoff:
		lut[i] = 0

def le16(i):
	return chr(i&0xff)+chr( (i>>8)&0xff )

def rgbw16(rgb):
	rgb = [ lut[ord(x)] for x in rgb ]
	m = min(rgb)
	mx = max(rgb)
	if mx < high_cutoff:
		return  '\0\0\0\0\0\0\0\0'
	else:
		return ''.join( le16(x) for x in rgb )+le16(m)


while True:
	s = sys.stdin.read(900*3)
	if len(s) != 900*3:
		break
	for i in xrange(900):
		si = i*3
		di = i+(i>=450)*30
		#q[di] = s[si:si+3]+'\x00'
		q[di] = rgbw16(s[si:si+3])


	sys.stdout.write(''.join(q)+'\xff\xff')
	
