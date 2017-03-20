/*  HexWS2811z - an OctoWS2811 fork to drive 16 LED strips at once using a Teensy 3.x
    by Erik Bosman <erik@minemu.org> Zero copy version, similar to the Fadecandy
    implementation.

    Based on:
    OctoWS2811 - High Performance WS2811 LED Display Library
    http://www.pjrc.com/teensy/td_libs_OctoWS2811.html
    Copyright (c) 2013 Paul Stoffregen, PJRC.COM, LLC

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
*/

#include <string.h>
#include "HexWS2811z.h"

#define LATCH_DELAY 200

uint32_t HexWS2811z::bufsize;
DMAChannel HexWS2811z::dma2;
DMAChannel HexWS2811z::dma3;

static uint8_t ones_zeroes[] = {0xFF, 0xFF, 0x00, 0x00};
static volatile uint8_t update_in_progress = 0;
static uint32_t update_completed_at = 0;


HexWS2811z::HexWS2811z(uint32_t bufsize)
{
	this->bufsize = bufsize;
}

// Waveform timing: these set the high time for a 0 and 1 bit, as a fraction of
// the total 800 kHz or 400 kHz clock cycle.  The scale is 0 to 255.  The Worldsemi
// datasheet seems T1H should be 600 ns of a 1250 ns cycle, or 48%.  That may
// erroneous information?  Other sources reason the chip actually samples the
// line close to the center of each bit time, so T1H should be 80% if TOH is 20%.
// The chips appear to work based on a simple one-shot delay triggered by the
// rising edge.  At least 1 chip tested retransmits 0 as a 330 ns pulse (26%) and
// a 1 as a 660 ns pulse (53%).  Perhaps it's actually sampling near 500 ns?
// There doesn't seem to be any advantage to making T1H less, as long as there
// is sufficient low time before the end of the cycle, so the next rising edge
// can be detected.  T0H has been lengthened slightly, because the pulse can
// narrow if the DMA controller has extra latency during bus arbitration.  If you
// have an insight about tuning these parameters AND you have actually tested on
// real LED strips, please contact paul@pjrc.com.  Please do not email based only
// on reading the datasheets and purely theoretical analysis.
#define WS2811_TIMING_T0H  60
#define WS2811_TIMING_T1H  176

// Discussion about timing and flicker & color shift issues:
// http://forum.pjrc.com/threads/23877-WS2812B-compatible-with-OctoWS2811-library?p=38190&viewfull=1#post38190


void HexWS2811z::begin(void)
{
	// configure the 16 output pins
	GPIOC_PCOR = 0xFF;
	pinMode(15, OUTPUT);	// strip #1
	pinMode(22, OUTPUT);	// strip #2
	pinMode(23, OUTPUT);	// strip #3
	pinMode(9, OUTPUT);	// strip #4
	pinMode(10, OUTPUT);	// strip #5
	pinMode(13, OUTPUT);	// strip #6
	pinMode(11, OUTPUT);	// strip #7
	pinMode(12, OUTPUT);	// strip #8

	GPIOD_PCOR = 0xFF;
	pinMode(2, OUTPUT);	// strip #9
	pinMode(14, OUTPUT);	// strip #10
	pinMode(7, OUTPUT);	// strip #11
	pinMode(8, OUTPUT);	// strip #12
	pinMode(6, OUTPUT);	// strip #13
	pinMode(20, OUTPUT);	// strip #14
	pinMode(21, OUTPUT);	// strip #15
	pinMode(5, OUTPUT);	// strip #16

	// create the two waveforms for WS2811 low and high bits
	analogWriteResolution(8);
	analogWriteFrequency(3, 800000);
	analogWriteFrequency(4, 800000);
	analogWrite(3, WS2811_TIMING_T0H);
	analogWrite(4, WS2811_TIMING_T1H);

	// pin 16 triggers DMA(port B) on rising edge (configure for pin 3's waveform)
	CORE_PIN16_CONFIG = PORT_PCR_IRQC(1)|PORT_PCR_MUX(3);
	pinMode(3, INPUT_PULLUP); // pin 3 no longer needed

	// pin 4 triggers DMA(port A) on both edges of high duty waveform
	CORE_PIN4_CONFIG = PORT_PCR_IRQC(3)|PORT_PCR_MUX(3);

	// DMA channel #2 writes the pixel data at 20% of the cycle
//	dma2.TCD->SADDR = frameBuffer;
	dma2.TCD->SADDR = NULL;
	dma2.TCD->SOFF = 1;
	dma2.TCD->ATTR_SRC = 0;
	dma2.TCD->SLAST = -bufsize;

	/* Send data to both PORT C and D in the same minor loop (executed after the same trigger( */

	#define PORT_DELTA ( (uint32_t)&GPIOD_PDOR - (uint32_t)&GPIOC_PDOR )
    dma2.TCD->DADDR = &GPIOC_PDOR;
	dma2.TCD->DOFF = PORT_DELTA;
                         /* loop GPIOC_PDOR, GPIOD_PDOR and back */
	dma2.TCD->ATTR_DST = ((31 - __builtin_clz(PORT_DELTA*2)) << 3) | 0;
	dma2.TCD->DLASTSGA = 0;

	dma2.TCD->NBYTES = 2;
	dma2.TCD->BITER = bufsize / 2;
	dma2.TCD->CITER = bufsize / 2;

	dma2.disableOnCompletion();

	// DMA channel #3 sets WS2811 high at the beginning of each cycle
    // and clears all the pins at 48% of the cycle
	dma3.TCD->SADDR = ones_zeroes;
	dma3.TCD->SOFF = 1;
	dma3.TCD->ATTR_SRC = (2 << 3) | 0; /* loop modulo 4 */
	dma3.TCD->SLAST = 0;

    dma3.TCD->DADDR = &GPIOC_PDOR;
	dma3.TCD->DOFF = PORT_DELTA;
                         /* loop GPIOC_PDOR, GPIOD_PDOR and back */
	dma3.TCD->ATTR_DST = ((31 - __builtin_clz(PORT_DELTA*2)) << 3) | 0;
	dma3.TCD->DLASTSGA = 0;

	dma3.TCD->NBYTES = 2;
	dma3.TCD->BITER = bufsize;
	dma3.TCD->CITER = bufsize;
	dma3.disableOnCompletion();
	dma3.interruptAtCompletion();

#ifdef __MK20DX256__
	MCM_CR = MCM_CR_SRAMLAP(1) | MCM_CR_SRAMUAP(0);
	AXBS_PRS0 = 0x1032;
#endif

	// route the edge detect interrupts to trigger the 3 channels
	dma2.triggerAtHardwareEvent(DMAMUX_SOURCE_PORTB);
	dma3.triggerAtHardwareEvent(DMAMUX_SOURCE_PORTA);

	// enable a done interrupts when channel #3 completes
	dma3.attachInterrupt(isr);
	//pinMode(9, OUTPUT); // testing: oscilloscope trigger
}

void HexWS2811z::isr(void)
{
	dma3.clearInterrupt();
	update_completed_at = micros();
	update_in_progress = 0;
}

int HexWS2811z::busy(void)
{
	if (update_in_progress) return 1;
	// busy for 50 us after the done interrupt, for WS2811 reset
	if (micros() - update_completed_at < LATCH_DELAY) return 1;
	return 0;
}

void HexWS2811z::show(void *frameBuffer)
{
	// wait for any prior DMA operation
	//Serial1.print("1");
	while (update_in_progress) ; 

	dma2.TCD->SADDR = frameBuffer;

	// wait for WS2811 reset
	while (micros() - update_completed_at < LATCH_DELAY) ;

	// ok to start, but we must be very careful to begin
	// without any prior 3 x 800kHz DMA requests pending
	uint32_t sc = FTM1_SC;
	uint32_t cv = FTM1_C1V;
	noInterrupts();
	// CAUTION: this code is timing critical.  Any editing should be
	// tested by verifying the oscilloscope trigger pulse at the end
	// always occurs while both waveforms are still low.  Simply
	// counting CPU cycles does not take into account other complex
	// factors, like flash cache misses and bus arbitration from USB
	// or other DMA.  Testing should be done with the oscilloscope
	// display set at infinite persistence and a variety of other I/O
	// performed to create realistic bus usage.  Even then, you really
	// should not mess with this timing critical code!
	update_in_progress = 1;
	while (FTM1_CNT <= cv) ; 
	while (FTM1_CNT > cv) ; // wait for beginning of an 800 kHz cycle
	while (FTM1_CNT < cv) ;
	FTM1_SC = sc & 0xE7;	// stop FTM1 timer (hopefully before it rolls over)
	PORTB_ISFR = (1<<0);    // clear any prior rising edge
	PORTC_ISFR = (1<<0);	// clear any prior low duty falling edge
	PORTA_ISFR = (1<<13);	// clear any prior high duty falling edge
	dma2.enable();		// enable all 3 DMA channels
	dma3.enable();
	FTM1_SC = sc;		// restart FTM1 timer
	interrupts();
}

